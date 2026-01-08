package webui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// SSEManager manages Server-Sent Events connections
type SSEManager struct {
	clients   map[string]chan []byte
	mu        sync.RWMutex
	logger    zerolog.Logger
	closeOnce sync.Once
	closed    bool
}

// EventType represents SSE event types
type EventType string

const (
	EventTypeStatus EventType = "status"
	EventTypeLog    EventType = "log"
	EventTypeConfig EventType = "config"
)

// SSEEvent represents a Server-Sent Event
type SSEEvent struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}

// NewSSEManager creates a new SSE manager
func NewSSEManager(logger zerolog.Logger) *SSEManager {
	return &SSEManager{
		clients: make(map[string]chan []byte),
		logger:  logger,
	}
}

// Handler handles SSE connections
func (m *SSEManager) Handler(c *fiber.Ctx) error {
	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	// Generate client ID
	clientID := fmt.Sprintf("%s-%d", c.IP(), time.Now().UnixNano())

	// Create channel for this client
	msgChan := make(chan []byte, 100)

	// Register client
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return c.Status(fiber.StatusServiceUnavailable).SendString("SSE service is shutting down")
	}
	m.clients[clientID] = msgChan
	m.mu.Unlock()

	m.logger.Info().Str("client", clientID).Msg("SSE client connected")

	// Cleanup on disconnect
	defer func() {
		m.mu.Lock()
		delete(m.clients, clientID)
		close(msgChan)
		m.mu.Unlock()
		m.logger.Info().Str("client", clientID).Msg("SSE client disconnected")
	}()

	// Use SetBodyStreamWriter to handle streaming
	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Send initial connection event
		initialEvent := SSEEvent{
			Type: "connected",
			Data: map[string]interface{}{
				"timestamp": time.Now(),
				"message":   "Connected to PawStream Edge Client",
			},
		}
		if data, err := json.Marshal(initialEvent); err == nil {
			w.WriteString(fmt.Sprintf("data: %s\n\n", string(data)))
			if err := w.Flush(); err != nil {
				m.logger.Error().Err(err).Str("client", clientID).Msg("Failed to send initial event")
				return
			}
		}

		// Ping ticker
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		// Event loop
		for {
			select {
			case msg, ok := <-msgChan:
				if !ok {
					m.logger.Info().Str("client", clientID).Msg("Message channel closed")
					return
				}
				// Send event
				w.WriteString(fmt.Sprintf("data: %s\n\n", string(msg)))
				if err := w.Flush(); err != nil {
					m.logger.Warn().Err(err).Str("client", clientID).Msg("Failed to send event")
					return
				}

			case <-ticker.C:
				// Send keep-alive comment
				w.WriteString(": keepalive\n\n")
				if err := w.Flush(); err != nil {
					m.logger.Warn().Err(err).Str("client", clientID).Msg("Failed to send keepalive")
					return
				}
			}
		}
	})

	return nil
}

// Broadcast sends an event to all connected clients
func (m *SSEManager) Broadcast(event SSEEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to marshal SSE event")
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return
	}

	for clientID, ch := range m.clients {
		select {
		case ch <- data:
		default:
			m.logger.Warn().Str("client", clientID).Msg("SSE client buffer full, dropping event")
		}
	}
}

// BroadcastStatus broadcasts a status update
func (m *SSEManager) BroadcastStatus(status interface{}) {
	m.Broadcast(SSEEvent{
		Type: EventTypeStatus,
		Data: status,
	})
}

// BroadcastLog broadcasts a log entry
func (m *SSEManager) BroadcastLog(entry LogEntry) {
	m.Broadcast(SSEEvent{
		Type: EventTypeLog,
		Data: entry,
	})
}

// BroadcastConfigChange broadcasts a config change notification
func (m *SSEManager) BroadcastConfigChange() {
	m.Broadcast(SSEEvent{
		Type: EventTypeConfig,
		Data: map[string]interface{}{
			"timestamp": time.Now(),
			"message":   "Configuration changed",
		},
	})
}

// ClientCount returns the number of connected clients
func (m *SSEManager) ClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

// Close closes all SSE connections
func (m *SSEManager) Close() {
	m.closeOnce.Do(func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		m.closed = true
		for clientID, ch := range m.clients {
			close(ch)
			delete(m.clients, clientID)
		}

		m.logger.Info().Msg("SSE manager closed")
	})
}
