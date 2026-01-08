package capture

import (
	"fmt"
	"os"
)

// FileSource reads video from a file
type FileSource struct {
	Path string
	Loop bool
}

// NewFileSource creates a new file source
func NewFileSource(path string, loop bool) *FileSource {
	return &FileSource{
		Path: path,
		Loop: loop,
	}
}

// Type returns the input type
func (s *FileSource) Type() InputType {
	return InputTypeFile
}

// FFmpegArgs returns the FFmpeg input arguments for file source
func (s *FileSource) FFmpegArgs() []string {
	args := []string{"-re"} // Read input at native frame rate
	
	if s.Loop {
		args = append(args, "-stream_loop", "-1") // Loop indefinitely
	}
	
	args = append(args, "-i", s.Path)
	return args
}

// Validate validates the file source configuration
func (s *FileSource) Validate() error {
	if s.Path == "" {
		return fmt.Errorf("file path is required")
	}
	
	// Check if file exists
	if _, err := os.Stat(s.Path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", s.Path)
		}
		return fmt.Errorf("failed to access file: %w", err)
	}
	
	return nil
}

// String returns a string representation
func (s *FileSource) String() string {
	loopStr := ""
	if s.Loop {
		loopStr = " (loop)"
	}
	return fmt.Sprintf("file source: %s%s", s.Path, loopStr)
}
