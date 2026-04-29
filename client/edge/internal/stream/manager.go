package stream

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lgc/pawstream/edge-client/internal/capture"
	"github.com/rs/zerolog"
)

// StreamStatus represents the current streaming status
type StreamStatus struct {
	State      string // "streaming", "stopped", "error", "reconnecting"
	InputType  string
	OutputURL  string
	ErrorCount int
	LastError  string
	StartTime  time.Time
}

// Config holds stream manager configuration
type Config struct {
	// Engine configuration
	Engine         EngineType
	VideoCodec     string
	VideoBitrate   int
	VideoWidth     int
	VideoHeight    int
	VideoFramerate int
	
	// FFmpeg specific
	FFmpegPreset   string
	FFmpegTune     string
	FFmpegHWAccel  string
	
	// GStreamer specific
	GStreamerLatencyMs   int
	GStreamerUseHardware bool
	GStreamerBufferSize  int
	
	// Reconnect configuration
	ReconnectInterval    time.Duration
	MaxReconnectAttempts int
}

// Manager manages video streaming
type Manager struct {
	input  capture.InputSource
	output string
	config Config
	logger zerolog.Logger

	engine         StreamEngine
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.RWMutex
	reconnectCount atomic.Int32
	errorCount     atomic.Int32
	lastError      string
	startTime      time.Time
	state          string
}

// NewManager creates a new stream manager
func NewManager(input capture.InputSource, output string, config Config, logger zerolog.Logger) *Manager {
	return &Manager{
		input:  input,
		output: output,
		config: config,
		logger: logger.With().Str("component", "stream-manager").Logger(),
		state:  "stopped",
	}
}

// Start starts the streaming
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state == "streaming" {
		return fmt.Errorf("already streaming")
	}

	m.ctx, m.cancel = context.WithCancel(ctx)
	m.startTime = time.Now()

	// Create stream engine using factory
	engineConfig := EngineConfig{
		Type:                 m.config.Engine,
		Input:                m.input,
		Output:               m.output,
		VideoCodec:           m.config.VideoCodec,
		VideoBitrate:         m.config.VideoBitrate,
		VideoWidth:           m.config.VideoWidth,
		VideoHeight:          m.config.VideoHeight,
		VideoFramerate:       m.config.VideoFramerate,
		FFmpegPreset:         m.config.FFmpegPreset,
		FFmpegTune:           m.config.FFmpegTune,
		FFmpegHWAccel:        m.config.FFmpegHWAccel,
		GStreamerLatencyMs:   m.config.GStreamerLatencyMs,
		GStreamerUseHardware: m.config.GStreamerUseHardware,
		GStreamerBufferSize:  m.config.GStreamerBufferSize,
	}
	
	var err error
	m.engine, err = NewStreamEngine(engineConfig, m.logger)
	if err != nil {
		m.state = "error"
		m.lastError = err.Error()
		m.errorCount.Add(1)
		return fmt.Errorf("failed to create stream engine: %w", err)
	}

	// Start engine
	if err := m.engine.Start(m.ctx); err != nil {
		m.state = "error"
		m.lastError = err.Error()
		m.errorCount.Add(1)
		return fmt.Errorf("failed to start streaming: %w", err)
	}

	m.state = "streaming"
	m.logger.Info().
		Str("input", m.input.String()).
		Str("output", m.output).
		Str("engine", m.engine.Name()).
		Msg("Streaming started")

	// Monitor engine process
	go m.monitorEngine()

	return nil
}

// Stop stops the streaming
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state == "stopped" {
		return nil
	}

	m.logger.Info().Msg("Stopping streaming...")

	// Cancel context
	if m.cancel != nil {
		m.cancel()
	}

	// Stop engine
	if m.engine != nil {
		if err := m.engine.Stop(); err != nil {
			m.logger.Error().Err(err).Str("engine", m.engine.Name()).Msg("Failed to stop engine")
		}
	}

	m.state = "stopped"
	m.logger.Info().Msg("Streaming stopped")
	return nil
}

// Status returns the current streaming status
func (m *Manager) Status() StreamStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return StreamStatus{
		State:      m.state,
		InputType:  string(m.input.Type()),
		OutputURL:  m.output,
		ErrorCount: int(m.errorCount.Load()),
		LastError:  m.lastError,
		StartTime:  m.startTime,
	}
}

// IsRunning returns whether streaming is active
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state == "streaming"
}

// monitorEngine monitors the stream engine process and handles errors
func (m *Manager) monitorEngine() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case err := <-m.engine.ErrorCh():
			m.handleError(err)
		}
	}
}

// handleError handles stream engine errors and implements reconnect logic
func (m *Manager) handleError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.errorCount.Add(1)
	m.lastError = err.Error()

	engineName := "unknown"
	if m.engine != nil {
		engineName = m.engine.Name()
	}

	m.logger.Error().Err(err).Str("engine", engineName).Msg("Stream engine error occurred")

	// Check if we should reconnect
	currentCount := m.reconnectCount.Add(1)
	if m.config.MaxReconnectAttempts > 0 && int(currentCount) > m.config.MaxReconnectAttempts {
		m.logger.Error().
			Int32("attempts", currentCount).
			Int("max", m.config.MaxReconnectAttempts).
			Msg("Max reconnect attempts reached, giving up")
		m.state = "error"
		return
	}

	m.state = "reconnecting"

	// Exponential backoff: base * 2^(attempt-1), capped at 5 minutes, with jitter
	backoff := float64(m.config.ReconnectInterval) * math.Pow(2, float64(currentCount-1))
	maxBackoff := float64(5 * time.Minute)
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	// Add random jitter: 0.5x to 1.0x of the computed backoff
	jitter := 0.5 + rand.Float64()*0.5
	delay := time.Duration(backoff * jitter)

	m.logger.Info().
		Int32("attempt", currentCount).
		Dur("delay", delay).
		Msg("Reconnecting...")

	// Wait before reconnecting
	time.Sleep(delay)

	// Attempt to restart with same engine type
	engineConfig := EngineConfig{
		Type:                 m.config.Engine,
		Input:                m.input,
		Output:               m.output,
		VideoCodec:           m.config.VideoCodec,
		VideoBitrate:         m.config.VideoBitrate,
		VideoWidth:           m.config.VideoWidth,
		VideoHeight:          m.config.VideoHeight,
		VideoFramerate:       m.config.VideoFramerate,
		FFmpegPreset:         m.config.FFmpegPreset,
		FFmpegTune:           m.config.FFmpegTune,
		FFmpegHWAccel:        m.config.FFmpegHWAccel,
		GStreamerLatencyMs:   m.config.GStreamerLatencyMs,
		GStreamerUseHardware: m.config.GStreamerUseHardware,
		GStreamerBufferSize:  m.config.GStreamerBufferSize,
	}
	
	var createErr error
	m.engine, createErr = NewStreamEngine(engineConfig, m.logger)
	if createErr != nil {
		m.logger.Error().Err(createErr).Msg("Failed to create engine for reconnect")
		return
	}

	if err := m.engine.Start(m.ctx); err != nil {
		m.logger.Error().Err(err).Msg("Failed to reconnect")
		// Will retry on next error
		return
	}

	m.state = "streaming"
	m.logger.Info().Msg("Reconnected successfully")
	m.reconnectCount.Store(0) // Reset counter on success
}
