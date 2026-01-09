package stream

import (
	"context"
)

// StreamEngine defines the interface for streaming engines (FFmpeg, GStreamer, etc.)
type StreamEngine interface {
	// Start starts the streaming engine
	Start(ctx context.Context) error

	// Stop stops the streaming engine gracefully
	Stop() error

	// IsRunning returns whether the engine is currently running
	IsRunning() bool

	// ErrorCh returns a channel that receives errors from the engine
	ErrorCh() <-chan error

	// Name returns the engine name (e.g., "ffmpeg", "gstreamer")
	Name() string

	// Stats returns the current engine statistics
	Stats() EngineStats
}

// EngineStats holds streaming statistics
type EngineStats struct {
	// FPS is the current frames per second
	FPS float64 `json:"fps"`

	// Bitrate is the current bitrate in bits per second
	Bitrate int `json:"bitrate"`

	// DroppedFrames is the total number of dropped frames
	DroppedFrames int `json:"dropped_frames"`

	// Encoder is the encoder being used (e.g., "x264", "vaapih264enc")
	Encoder string `json:"encoder"`
}
