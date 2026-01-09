package stream

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"
)

// FFmpegEngine implements StreamEngine interface using FFmpeg
type FFmpegEngine struct {
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
	
	// Stats
	stats EngineStats
}

// NewFFmpegEngine creates a new FFmpeg engine
func NewFFmpegEngine(inputArgs, outputArgs []string, logger zerolog.Logger) *FFmpegEngine {
	return &FFmpegEngine{
		inputArgs:  inputArgs,
		outputArgs: outputArgs,
		logger:     logger.With().Str("component", "ffmpeg").Str("engine", "ffmpeg").Logger(),
		errorCh:    make(chan error, 1),
		doneCh:     make(chan struct{}),
		stats: EngineStats{
			Encoder: "libx264", // Default, will be detected
		},
	}
}

// NewFFmpegManager is deprecated, use NewFFmpegEngine instead
// Kept for backward compatibility
func NewFFmpegManager(inputArgs, outputArgs []string, logger zerolog.Logger) *FFmpegEngine {
	return NewFFmpegEngine(inputArgs, outputArgs, logger)
}

// Start starts the FFmpeg process
func (m *FFmpegEngine) Start(ctx context.Context) error {
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

	// Buffer to capture stderr for error reporting
	var stderrBuf strings.Builder

	// Monitor stderr (FFmpeg outputs logs to stderr)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			m.logger.Debug().Str("stderr", line).Msg("FFmpeg log")
			
			// Buffer stderr for error reporting
			stderrBuf.WriteString(line)
			stderrBuf.WriteString("\n")
		}
	}()

	// Wait for process to exit
	go func() {
		err := m.cmd.Wait()
		m.running.Store(false)

		if err != nil {
			// Log FFmpeg stderr output on error
			if stderrOutput := stderrBuf.String(); stderrOutput != "" {
				// Get last 50 lines of stderr
				lines := strings.Split(strings.TrimSpace(stderrOutput), "\n")
				if len(lines) > 50 {
					lines = lines[len(lines)-50:]
				}
				lastOutput := strings.Join(lines, "\n")
				
				m.logger.Error().
					Err(err).
					Str("ffmpeg_output", lastOutput).
					Msg("FFmpeg process exited with error")
			} else {
				m.logger.Error().Err(err).Msg("FFmpeg process exited with error")
			}
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
func (m *FFmpegEngine) Stop() error {
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
func (m *FFmpegEngine) Wait() error {
	<-m.doneCh
	return nil
}

// IsRunning returns whether the FFmpeg process is running
func (m *FFmpegEngine) IsRunning() bool {
	return m.running.Load()
}

// ErrorCh returns the error channel
func (m *FFmpegEngine) ErrorCh() <-chan error {
	return m.errorCh
}

// Name returns the engine name
func (m *FFmpegEngine) Name() string {
	return "ffmpeg"
}

// Stats returns the current engine statistics
func (m *FFmpegEngine) Stats() EngineStats {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.stats
}
