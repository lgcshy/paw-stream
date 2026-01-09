package stream

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/lgc/pawstream/edge-client/internal/capture"
)

// buildPipeline builds the GStreamer pipeline based on input source and configuration
func (g *GStreamerEngine) buildPipeline() ([]string, error) {
	var elements []string

	// Build input source elements
	inputElements, err := g.buildInputElements()
	if err != nil {
		return nil, fmt.Errorf("failed to build input elements: %w", err)
	}
	elements = append(elements, inputElements...)

	// Add video conversion and scaling if needed
	elements = append(elements,
		"videoconvert",
		fmt.Sprintf("video/x-raw,width=%d,height=%d,framerate=%d/1",
			g.config.VideoWidth, g.config.VideoHeight, g.config.VideoFramerate),
	)

	// Build encoder elements
	encoderElements, err := g.buildEncoderElements()
	if err != nil {
		return nil, fmt.Errorf("failed to build encoder elements: %w", err)
	}
	elements = append(elements, encoderElements...)

	// Add H.264 parser
	elements = append(elements, "h264parse")

	// Build output elements
	outputElements := g.buildOutputElements()
	elements = append(elements, outputElements...)

	// Join elements with "!"
	pipeline := make([]string, 0, len(elements)*2-1)
	for i, elem := range elements {
		if i > 0 {
			pipeline = append(pipeline, "!")
		}
		pipeline = append(pipeline, elem)
	}

	return pipeline, nil
}

// buildInputElements builds pipeline elements for the input source
func (g *GStreamerEngine) buildInputElements() ([]string, error) {
	switch g.input.Type() {
	case capture.InputTypeV4L2:
		return g.buildV4L2Input()
	case capture.InputTypeRTSP:
		return g.buildRTSPInput()
	case capture.InputTypeFile:
		return g.buildFileInput()
	case capture.InputTypeTest:
		return g.buildTestInput()
	default:
		return nil, fmt.Errorf("unsupported input type: %s", g.input.Type())
	}
}

// buildV4L2Input builds pipeline for V4L2 camera input
func (g *GStreamerEngine) buildV4L2Input() ([]string, error) {
	source := g.input.String()
	// Extract device path from source (e.g., "v4l2:///dev/video0")
	device := strings.TrimPrefix(source, "v4l2://")
	if device == "" {
		device = "/dev/video0"
	}

	return []string{
		fmt.Sprintf("v4l2src device=%s", device),
	}, nil
}

// buildRTSPInput builds pipeline for RTSP stream input
func (g *GStreamerEngine) buildRTSPInput() ([]string, error) {
	// Get RTSP URL from RTSP source
	rtspSrc, ok := g.input.(*capture.RTSPSource)
	if !ok {
		return nil, fmt.Errorf("input is not an RTSP source")
	}

	return []string{
		fmt.Sprintf("rtspsrc location=%s latency=%d", rtspSrc.URL, g.config.LatencyMs),
		"rtph264depay",
		"h264parse",
		"avdec_h264",
	}, nil
}

// buildFileInput builds pipeline for file input
func (g *GStreamerEngine) buildFileInput() ([]string, error) {
	source := g.input.String()
	// Extract file path from source (e.g., "file:///path/to/video.mp4")
	filePath := strings.TrimPrefix(source, "file://")
	if filePath == "" {
		return nil, fmt.Errorf("invalid file path")
	}

	return []string{
		fmt.Sprintf("filesrc location=%s", filePath),
		"decodebin",
	}, nil
}

// buildTestInput builds pipeline for test pattern input
func (g *GStreamerEngine) buildTestInput() ([]string, error) {
	return []string{
		"videotestsrc pattern=smpte",
	}, nil
}

// buildEncoderElements builds pipeline elements for video encoding
func (g *GStreamerEngine) buildEncoderElements() ([]string, error) {
	if g.config.UseHardware {
		// Try hardware encoders first
		if encoder := g.detectHardwareEncoder(); encoder != "" {
			g.stats.Encoder = encoder
			return g.buildHardwareEncoder(encoder)
		}
		g.logger.Warn().Msg("Hardware encoder not available, falling back to software")
	}

	// Fallback to software encoder
	g.stats.Encoder = "x264enc"
	return g.buildSoftwareEncoder()
}

// detectHardwareEncoder detects available hardware encoders
func (g *GStreamerEngine) detectHardwareEncoder() string {
	// Check VAAPI (Intel)
	if isGStreamerElementAvailable("vaapih264enc") {
		return "vaapih264enc"
	}

	// Check NVENC (NVIDIA)
	if isGStreamerElementAvailable("nvh264enc") {
		return "nvh264enc"
	}

	// Check VideoToolbox (macOS)
	if isGStreamerElementAvailable("vtenc_h264") {
		return "vtenc_h264"
	}

	// Check OMX (Raspberry Pi)
	if isGStreamerElementAvailable("omxh264enc") {
		return "omxh264enc"
	}

	return ""
}

// buildHardwareEncoder builds hardware encoder element
func (g *GStreamerEngine) buildHardwareEncoder(encoder string) ([]string, error) {
	bitrateKbps := g.config.VideoBitrate / 1000

	switch encoder {
	case "vaapih264enc":
		return []string{
			fmt.Sprintf("vaapih264enc bitrate=%d rate-control=cbr", bitrateKbps),
		}, nil

	case "nvh264enc":
		return []string{
			fmt.Sprintf("nvh264enc bitrate=%d rc-mode=cbr", bitrateKbps),
		}, nil

	case "vtenc_h264":
		return []string{
			fmt.Sprintf("vtenc_h264 bitrate=%d", bitrateKbps),
		}, nil

	case "omxh264enc":
		return []string{
			fmt.Sprintf("omxh264enc target-bitrate=%d", g.config.VideoBitrate),
		}, nil

	default:
		return nil, fmt.Errorf("unknown hardware encoder: %s", encoder)
	}
}

// buildSoftwareEncoder builds software encoder element (x264)
func (g *GStreamerEngine) buildSoftwareEncoder() ([]string, error) {
	bitrateKbps := g.config.VideoBitrate / 1000

	return []string{
		fmt.Sprintf("x264enc bitrate=%d tune=zerolatency speed-preset=ultrafast", bitrateKbps),
	}, nil
}

// buildOutputElements builds pipeline elements for RTSP output
// Note: GStreamer 1.0 doesn't have rtspclientsink by default
// We use flvmux + rtmpsink to push to MediaMTX via RTMP, which it converts to RTSP
func (g *GStreamerEngine) buildOutputElements() []string {
	// Convert RTSP URL to RTMP URL for MediaMTX
	// rtsp://user:pass@host:8554/path -> rtmp://host:1935/path?user=user&pass=pass
	rtmpURL := convertRTSPtoRTMP(g.output)
	
	return []string{
		"flvmux streamable=true",
		fmt.Sprintf("rtmpsink location=%s", rtmpURL),
	}
}

// convertRTSPtoRTMP converts RTSP URL to RTMP URL for MediaMTX
func convertRTSPtoRTMP(rtspURL string) string {
	// Parse rtsp://user:pass@host:port/path
	url := strings.TrimPrefix(rtspURL, "rtsp://")
	
	var auth, hostPort, path string
	// Split auth and host
	if idx := strings.Index(url, "@"); idx >= 0 {
		auth = url[:idx]
		url = url[idx+1:]
	}
	
	// Split host and path
	if idx := strings.Index(url, "/"); idx >= 0 {
		hostPort = url[:idx]
		path = url[idx+1:]
	} else {
		hostPort = url
		path = ""
	}
	
	// Remove :8554 if present, MediaMTX RTMP is on 1935
	host := strings.Split(hostPort, ":")[0]
	
	// Build RTMP URL
	rtmpURL := fmt.Sprintf("rtmp://%s:1935/%s", host, path)
	
	// Add auth as query parameters if present
	if auth != "" {
		parts := strings.Split(auth, ":")
		if len(parts) == 2 {
			rtmpURL += fmt.Sprintf("?user=%s&pass=%s", parts[0], parts[1])
		}
	}
	
	return rtmpURL
}

// isGStreamerElementAvailable checks if a GStreamer element is available
func isGStreamerElementAvailable(element string) bool {
	cmd := fmt.Sprintf("gst-inspect-1.0 %s", element)
	output, err := execCommand("sh", "-c", cmd)
	if err != nil {
		return false
	}
	// If element exists, gst-inspect will output its details
	return len(output) > 0 && !strings.Contains(output, "No such element")
}

// execCommand is a helper to execute shell commands
func execCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
