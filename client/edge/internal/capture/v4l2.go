package capture

import (
	"fmt"
	"os"
	"runtime"
)

// V4L2Source captures video from a V4L2 device (Linux only)
type V4L2Source struct {
	Device    string
	Width     int
	Height    int
	Framerate int
}

// NewV4L2Source creates a new V4L2 source
func NewV4L2Source(device string, width, height, framerate int) *V4L2Source {
	return &V4L2Source{
		Device:    device,
		Width:     width,
		Height:    height,
		Framerate: framerate,
	}
}

// Type returns the input type
func (s *V4L2Source) Type() InputType {
	return InputTypeV4L2
}

// FFmpegArgs returns the FFmpeg input arguments for V4L2 source
func (s *V4L2Source) FFmpegArgs() []string {
	return []string{
		"-f", "v4l2",
		"-video_size", fmt.Sprintf("%dx%d", s.Width, s.Height),
		"-framerate", fmt.Sprintf("%d", s.Framerate),
		"-i", s.Device,
	}
}

// Validate validates the V4L2 source configuration
func (s *V4L2Source) Validate() error {
	// Check if running on Linux
	if runtime.GOOS != "linux" {
		return fmt.Errorf("v4l2 is only supported on Linux (current OS: %s)", runtime.GOOS)
	}
	
	if s.Device == "" {
		return fmt.Errorf("device path is required")
	}
	
	// Check if device exists
	if _, err := os.Stat(s.Device); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("device not found: %s", s.Device)
		}
		return fmt.Errorf("failed to access device: %w", err)
	}
	
	if s.Width <= 0 || s.Height <= 0 {
		return fmt.Errorf("invalid video dimensions: %dx%d", s.Width, s.Height)
	}
	
	if s.Framerate <= 0 {
		return fmt.Errorf("invalid framerate: %d", s.Framerate)
	}
	
	return nil
}

// String returns a string representation
func (s *V4L2Source) String() string {
	return fmt.Sprintf("v4l2 source: %s (%dx%d@%dfps)", s.Device, s.Width, s.Height, s.Framerate)
}
