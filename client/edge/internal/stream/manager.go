package stream

import (
	"context"
	"fmt"
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
	VideoCodec           string
	VideoBitrate         int
	Preset               string
	ReconnectInterval    time.Duration
	MaxReconnectAttempts int
}

// Manager manages video streaming
type Manager struct {
	input  capture.InputSource
	output string
	config Config
	logger zerolog.Logger

	ffmpeg         *FFmpegManager
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

	// Build FFmpeg arguments
	inputArgs := m.input.FFmpegArgs()
	outputArgs := m.buildOutputArgs()

	// Create FFmpeg manager
	m.ffmpeg = NewFFmpegManager(inputArgs, outputArgs, m.logger)

	// Start FFmpeg
	if err := m.ffmpeg.Start(m.ctx); err != nil {
		m.state = "error"
		m.lastError = err.Error()
		m.errorCount.Add(1)
		return fmt.Errorf("failed to start streaming: %w", err)
	}

	m.state = "streaming"
	m.logger.Info().
		Str("input", m.input.String()).
		Str("output", m.output).
		Msg("Streaming started")

	// Monitor FFmpeg process
	go m.monitorFFmpeg()

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

	// Stop FFmpeg
	if m.ffmpeg != nil {
		if err := m.ffmpeg.Stop(); err != nil {
			m.logger.Error().Err(err).Msg("Failed to stop FFmpeg")
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

// monitorFFmpeg monitors the FFmpeg process and handles errors
func (m *Manager) monitorFFmpeg() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case err := <-m.ffmpeg.ErrorCh():
			m.handleError(err)
		}
	}
}

// handleError handles FFmpeg errors and implements reconnect logic
func (m *Manager) handleError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.errorCount.Add(1)
	m.lastError = err.Error()

	m.logger.Error().Err(err).Msg("FFmpeg error occurred")

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
	m.logger.Info().
		Int32("attempt", currentCount).
		Dur("interval", m.config.ReconnectInterval).
		Msg("Reconnecting...")

	// Wait before reconnecting
	time.Sleep(m.config.ReconnectInterval)

	// Attempt to restart
	inputArgs := m.input.FFmpegArgs()
	outputArgs := m.buildOutputArgs()
	m.ffmpeg = NewFFmpegManager(inputArgs, outputArgs, m.logger)

	if err := m.ffmpeg.Start(m.ctx); err != nil {
		m.logger.Error().Err(err).Msg("Failed to reconnect")
		// Will retry on next error
		return
	}

	m.state = "streaming"
	m.logger.Info().Msg("Reconnected successfully")
	m.reconnectCount.Store(0) // Reset counter on success
}

// buildOutputArgs builds FFmpeg output arguments
func (m *Manager) buildOutputArgs() []string {
	args := []string{
		"-c:v", m.config.VideoCodec,
		"-preset", m.config.Preset,
		"-b:v", fmt.Sprintf("%d", m.config.VideoBitrate),
		"-tune", "zerolatency",
		"-f", "rtsp",
	}

	// Add output URL
	args = append(args, m.output)

	return args
}
