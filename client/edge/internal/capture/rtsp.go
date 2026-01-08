package capture

import (
	"fmt"
	"strings"
)

// RTSPSource captures video from an RTSP stream
type RTSPSource struct {
	URL       string
	Transport string // "tcp" or "udp"
}

// NewRTSPSource creates a new RTSP source
func NewRTSPSource(url, transport string) *RTSPSource {
	if transport == "" {
		transport = "tcp" // Default to TCP for reliability
	}
	return &RTSPSource{
		URL:       url,
		Transport: transport,
	}
}

// Type returns the input type
func (s *RTSPSource) Type() InputType {
	return InputTypeRTSP
}

// FFmpegArgs returns the FFmpeg input arguments for RTSP source
func (s *RTSPSource) FFmpegArgs() []string {
	args := []string{}
	
	// Set RTSP transport protocol
	if s.Transport != "" {
		args = append(args, "-rtsp_transport", s.Transport)
	}
	
	// Add input URL
	args = append(args, "-i", s.URL)
	
	return args
}

// Validate validates the RTSP source configuration
func (s *RTSPSource) Validate() error {
	if s.URL == "" {
		return fmt.Errorf("RTSP URL is required")
	}
	
	// Basic URL validation
	if !strings.HasPrefix(s.URL, "rtsp://") && !strings.HasPrefix(s.URL, "rtsps://") {
		return fmt.Errorf("invalid RTSP URL: must start with rtsp:// or rtsps://")
	}
	
	// Validate transport
	if s.Transport != "" && s.Transport != "tcp" && s.Transport != "udp" {
		return fmt.Errorf("invalid transport: %s (must be tcp or udp)", s.Transport)
	}
	
	return nil
}

// String returns a string representation
func (s *RTSPSource) String() string {
	return fmt.Sprintf("rtsp source: %s (transport: %s)", s.URL, s.Transport)
}
