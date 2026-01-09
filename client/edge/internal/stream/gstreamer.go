package stream

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/lgc/pawstream/edge-client/internal/capture"
	"github.com/rs/zerolog"
)

// GStreamerEngine implements StreamEngine interface using GStreamer
type GStreamerEngine struct {
	cmd     *exec.Cmd
	input   capture.InputSource
	output  string
	config  GStreamerConfig
	logger  zerolog.Logger
	ctx     context.Context
	cancel  context.CancelFunc

	// State
	running atomic.Bool
	mu      sync.Mutex
	errorCh chan error
	doneCh  chan struct{}

	// Stats
	stats      EngineStats
	statsMu    sync.RWMutex
	fpsRegexp  *regexp.Regexp
	bpsRegexp  *regexp.Regexp
}

// GStreamerConfig holds GStreamer-specific configuration
type GStreamerConfig struct {
	LatencyMs   int
	UseHardware bool
	BufferSize  int
	VideoCodec  string
	VideoBitrate int
	VideoWidth  int
	VideoHeight int
	VideoFramerate int
}

// NewGStreamerEngine creates a new GStreamer engine
func NewGStreamerEngine(input capture.InputSource, output string, config GStreamerConfig, logger zerolog.Logger) *GStreamerEngine {
	return &GStreamerEngine{
		input:     input,
		output:    output,
		config:    config,
		logger:    logger.With().Str("component", "gstreamer").Str("engine", "gstreamer").Logger(),
		errorCh:   make(chan error, 1),
		doneCh:    make(chan struct{}),
		fpsRegexp: regexp.MustCompile(`current:\s+(\d+\.\d+)`),
		bpsRegexp: regexp.MustCompile(`bitrate:\s+(\d+)`),
		stats: EngineStats{
			Encoder: "vaapih264enc", // Default, will be detected
		},
	}
}

// Start starts the GStreamer pipeline
func (g *GStreamerEngine) Start(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.running.Load() {
		return fmt.Errorf("gstreamer already running")
	}

	// Check if gstreamer is installed
	if !IsGStreamerInstalled() {
		return fmt.Errorf("gstreamer is not installed")
	}

	// Create context for GStreamer process
	g.ctx, g.cancel = context.WithCancel(ctx)

	// Build GStreamer pipeline
	pipeline, err := g.buildPipeline()
	if err != nil {
		return fmt.Errorf("failed to build pipeline: %w", err)
	}

	// Join pipeline elements into a command string
	pipelineStr := strings.Join(pipeline, " ")
	
	// Use shell to execute the command, which properly handles quotes and special characters
	cmdStr := fmt.Sprintf("gst-launch-1.0 -q %s", pipelineStr)
	g.cmd = exec.CommandContext(g.ctx, "sh", "-c", cmdStr)

	// Capture stdout and stderr
	stdout, err := g.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := g.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the process
	if err := g.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start gstreamer: %w", err)
	}

	g.running.Store(true)

	g.logger.Info().
		Strs("pipeline", pipeline).
		Int("pid", g.cmd.Process.Pid).
		Msg("GStreamer pipeline started")

	// Monitor stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			g.logger.Debug().Str("stdout", line).Msg("GStreamer output")
		}
	}()

	// Buffer to capture stderr for error reporting
	var stderrBuf strings.Builder
	
	// Monitor stderr (GStreamer outputs logs to stderr)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			g.parseStatsFromLog(line)
			g.logger.Debug().Str("stderr", line).Msg("GStreamer log")
			
			// Buffer stderr for error reporting
			stderrBuf.WriteString(line)
			stderrBuf.WriteString("\n")
		}
	}()

	// Wait for process to exit
	go func() {
		err := g.cmd.Wait()
		g.running.Store(false)

		if err != nil {
			// Log GStreamer stderr output on error
			if stderrOutput := stderrBuf.String(); stderrOutput != "" {
				g.logger.Error().
					Err(err).
					Str("gstreamer_output", stderrOutput).
					Msg("GStreamer process exited with error")
			} else {
				g.logger.Error().Err(err).Msg("GStreamer process exited with error")
			}
			select {
			case g.errorCh <- err:
			default:
			}
		} else {
			g.logger.Info().Msg("GStreamer process exited normally")
		}

		close(g.doneCh)
	}()

	return nil
}

// Stop stops the GStreamer pipeline gracefully
func (g *GStreamerEngine) Stop() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.running.Load() {
		return nil
	}

	g.logger.Info().Msg("Stopping GStreamer pipeline...")

	// Cancel context to stop the process
	if g.cancel != nil {
		g.cancel()
	}

	// Wait for process to exit
	<-g.doneCh

	g.logger.Info().Msg("GStreamer pipeline stopped")
	return nil
}

// IsRunning returns whether the GStreamer pipeline is running
func (g *GStreamerEngine) IsRunning() bool {
	return g.running.Load()
}

// ErrorCh returns the error channel
func (g *GStreamerEngine) ErrorCh() <-chan error {
	return g.errorCh
}

// Name returns the engine name
func (g *GStreamerEngine) Name() string {
	return "gstreamer"
}

// Stats returns the current engine statistics
func (g *GStreamerEngine) Stats() EngineStats {
	g.statsMu.RLock()
	defer g.statsMu.RUnlock()
	return g.stats
}

// parseStatsFromLog parses statistics from GStreamer log output
func (g *GStreamerEngine) parseStatsFromLog(line string) {
	g.statsMu.Lock()
	defer g.statsMu.Unlock()

	// Parse FPS: "current: 29.97"
	if matches := g.fpsRegexp.FindStringSubmatch(line); len(matches) > 1 {
		if fps, err := strconv.ParseFloat(matches[1], 64); err == nil {
			g.stats.FPS = fps
		}
	}

	// Parse bitrate: "bitrate: 2000000"
	if matches := g.bpsRegexp.FindStringSubmatch(line); len(matches) > 1 {
		if bps, err := strconv.Atoi(matches[1]); err == nil {
			g.stats.Bitrate = bps
		}
	}
}

// buildPipeline is implemented in gstreamer_pipeline.go

// IsGStreamerInstalled checks if GStreamer is installed on the system
func IsGStreamerInstalled() bool {
	cmd := exec.Command("gst-launch-1.0", "--version")
	return cmd.Run() == nil
}
