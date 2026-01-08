package capture

import "fmt"

// TestSource generates test video using FFmpeg's testsrc
type TestSource struct {
	Width     int
	Height    int
	Framerate int
}

// NewTestSource creates a new test source
func NewTestSource(width, height, framerate int) *TestSource {
	return &TestSource{
		Width:     width,
		Height:    height,
		Framerate: framerate,
	}
}

// Type returns the input type
func (s *TestSource) Type() InputType {
	return InputTypeTest
}

// FFmpegArgs returns the FFmpeg input arguments for test source
func (s *TestSource) FFmpegArgs() []string {
	return []string{
		"-f", "lavfi",
		"-i", fmt.Sprintf("testsrc=size=%dx%d:rate=%d", s.Width, s.Height, s.Framerate),
	}
}

// Validate validates the test source configuration
func (s *TestSource) Validate() error {
	if s.Width <= 0 || s.Height <= 0 {
		return fmt.Errorf("invalid video dimensions: %dx%d", s.Width, s.Height)
	}
	if s.Framerate <= 0 {
		return fmt.Errorf("invalid framerate: %d", s.Framerate)
	}
	return nil
}

// String returns a string representation
func (s *TestSource) String() string {
	return fmt.Sprintf("test source (%dx%d@%dfps)", s.Width, s.Height, s.Framerate)
}
