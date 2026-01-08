package capture

import "fmt"

// InputType represents the type of video input source
type InputType string

const (
	InputTypeTest  InputType = "test"
	InputTypeFile  InputType = "file"
	InputTypeV4L2  InputType = "v4l2"
	InputTypeRTSP  InputType = "rtsp"
)

// InputSource represents a video input source
type InputSource interface {
	// Type returns the input type
	Type() InputType

	// FFmpegArgs returns the FFmpeg input arguments
	FFmpegArgs() []string

	// Validate validates the input source configuration
	Validate() error

	// String returns a string representation of the input source
	String() string
}

// ValidateInputType checks if the input type is valid
func ValidateInputType(t string) error {
	switch InputType(t) {
	case InputTypeTest, InputTypeFile, InputTypeV4L2, InputTypeRTSP:
		return nil
	default:
		return fmt.Errorf("invalid input type: %s (must be test, file, v4l2, or rtsp)", t)
	}
}
