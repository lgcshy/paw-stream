package stream

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"
)

// FFmpegManager manages an FFmpeg process for streaming
type FFmpegManager struct {
	cmd        *exec.Cmd
	inputArgs  []string
	outputArgs []string
	logger     zerolog.Logger
	ctx        context.Context
	cancel     context.CancelFunc

	// State
	running atomic.Bool
	mu      sync.Mutex
	errorCh chan error
	doneCh  chan struct{}
}

// NewFFmpegManager creates a new FFmpeg manager
func NewFFmpegManager(inputArgs, outputArgs []string, logger zerolog.Logger) *FFmpegManager {
	return &FFmpegManager{
		inputArgs:  inputArgs,
		outputArgs: outputArgs,
		logger:     logger.With().Str("component", "ffmpeg").Logger(),
		errorCh:    make(chan error, 1),
		doneCh:     make(chan struct{}),
	}
}

// Start starts the FFmpeg process
func (m *FFmpegManager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running.Load() {
		return fmt.Errorf("ffmpeg already running")
	}

	// Create context for FFmpeg process
	m.ctx, m.cancel = context.WithCancel(ctx)

	// Build FFmpeg command
	args := append([]string{}, m.inputArgs...)
	args = append(args, m.outputArgs...)

	m.cmd = exec.CommandContext(m.ctx, "ffmpeg", args...)

	// Capture stdout and stderr
	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := m.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the process
	if err := m.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	m.running.Store(true)

	m.logger.Info().
		Strs("args", args).
		Int("pid", m.cmd.Process.Pid).
		Msg("FFmpeg process started")

	// Monitor stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			m.logger.Debug().Str("stdout", line).Msg("FFmpeg output")
		}
	}()

	// Monitor stderr (FFmpeg outputs logs to stderr)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			m.logger.Debug().Str("stderr", line).Msg("FFmpeg log")
		}
	}()

	// Wait for process to exit
	go func() {
		err := m.cmd.Wait()
		m.running.Store(false)

		if err != nil {
			m.logger.Error().Err(err).Msg("FFmpeg process exited with error")
			select {
			case m.errorCh <- err:
			default:
			}
		} else {
			m.logger.Info().Msg("FFmpeg process exited normally")
		}

		close(m.doneCh)
	}()

	return nil
}

// Stop stops the FFmpeg process gracefully
func (m *FFmpegManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running.Load() {
		return nil
	}

	m.logger.Info().Msg("Stopping FFmpeg process...")

	// Cancel context to stop the process
	if m.cancel != nil {
		m.cancel()
	}

	// Wait for process to exit (with timeout handled by context)
	<-m.doneCh

	m.logger.Info().Msg("FFmpeg process stopped")
	return nil
}

// Wait waits for the FFmpeg process to exit
func (m *FFmpegManager) Wait() error {
	<-m.doneCh
	return nil
}

// IsRunning returns whether the FFmpeg process is running
func (m *FFmpegManager) IsRunning() bool {
	return m.running.Load()
}

// ErrorCh returns the error channel
func (m *FFmpegManager) ErrorCh() <-chan error {
	return m.errorCh
}
