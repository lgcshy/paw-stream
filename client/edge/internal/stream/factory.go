package stream

import (
	"fmt"

	"github.com/lgc/pawstream/edge-client/internal/capture"
	"github.com/rs/zerolog"
)

// EngineType represents the type of streaming engine
type EngineType string

const (
	FFmpegEngineType     EngineType = "ffmpeg"
	GStreamerEngineType  EngineType = "gstreamer"
)

// EngineConfig holds configuration for creating a stream engine
type EngineConfig struct {
	Type           EngineType
	Input          capture.InputSource
	Output         string
	VideoCodec     string
	VideoBitrate   int
	VideoWidth     int
	VideoHeight    int
	VideoFramerate int
	FFmpegPreset   string
	FFmpegTune     string
	FFmpegHWAccel  string
	GStreamerLatencyMs  int
	GStreamerUseHardware bool
	GStreamerBufferSize int
}

// NewStreamEngine creates a new stream engine based on the engine type
func NewStreamEngine(config EngineConfig, logger zerolog.Logger) (StreamEngine, error) {
	switch config.Type {
	case FFmpegEngineType, "":
		return newFFmpegEngineFromConfig(config, logger), nil
		
	case GStreamerEngineType:
		// Check if GStreamer is installed
		if !IsGStreamerInstalled() {
			logger.Warn().Msg("GStreamer not installed, falling back to FFmpeg")
			return newFFmpegEngineFromConfig(config, logger), nil
		}
		return newGStreamerEngineFromConfig(config, logger), nil
		
	default:
		return nil, fmt.Errorf("unsupported engine type: %s", config.Type)
	}
}

// newFFmpegEngineFromConfig creates FFmpeg engine from engine config
func newFFmpegEngineFromConfig(config EngineConfig, logger zerolog.Logger) *FFmpegEngine {
	inputArgs := config.Input.FFmpegArgs()
	
	outputArgs := []string{
		"-c:v", config.VideoCodec,
		"-preset", config.FFmpegPreset,
		"-tune", config.FFmpegTune,
		"-b:v", fmt.Sprintf("%d", config.VideoBitrate),
		"-f", "rtsp",
	}
	
	// Add hardware acceleration if specified
	if config.FFmpegHWAccel != "none" && config.FFmpegHWAccel != "" {
		if config.FFmpegHWAccel == "auto" {
			// Try to detect hardware encoder
			if hwEncoder := detectFFmpegHWEncoder(); hwEncoder != "" {
				outputArgs[1] = hwEncoder // Replace codec
			}
		} else {
			// Use specified hardware encoder
			outputArgs[1] = getFFmpegHWEncoder(config.FFmpegHWAccel)
		}
	}
	
	// Add output URL
	outputArgs = append(outputArgs, config.Output)
	
	return NewFFmpegEngine(inputArgs, outputArgs, logger)
}

// newGStreamerEngineFromConfig creates GStreamer engine from engine config
func newGStreamerEngineFromConfig(config EngineConfig, logger zerolog.Logger) *GStreamerEngine {
	gstConfig := GStreamerConfig{
		LatencyMs:      config.GStreamerLatencyMs,
		UseHardware:    config.GStreamerUseHardware,
		BufferSize:     config.GStreamerBufferSize,
		VideoCodec:     config.VideoCodec,
		VideoBitrate:   config.VideoBitrate,
		VideoWidth:     config.VideoWidth,
		VideoHeight:    config.VideoHeight,
		VideoFramerate: config.VideoFramerate,
	}
	
	return NewGStreamerEngine(config.Input, config.Output, gstConfig, logger)
}

// detectFFmpegHWEncoder detects available hardware encoder for FFmpeg
func detectFFmpegHWEncoder() string {
	// Check NVENC (NVIDIA)
	if isFFmpegEncoderAvailable("h264_nvenc") {
		return "h264_nvenc"
	}
	
	// Check VAAPI (Intel)
	if isFFmpegEncoderAvailable("h264_vaapi") {
		return "h264_vaapi"
	}
	
	// Check QSV (Intel Quick Sync)
	if isFFmpegEncoderAvailable("h264_qsv") {
		return "h264_qsv"
	}
	
	// Check VideoToolbox (macOS)
	if isFFmpegEncoderAvailable("h264_videotoolbox") {
		return "h264_videotoolbox"
	}
	
	return "" // No hardware encoder found
}

// getFFmpegHWEncoder returns the FFmpeg encoder name for a hardware acceleration type
func getFFmpegHWEncoder(hwaccel string) string {
	switch hwaccel {
	case "nvenc":
		return "h264_nvenc"
	case "vaapi":
		return "h264_vaapi"
	case "qsv":
		return "h264_qsv"
	case "videotoolbox":
		return "h264_videotoolbox"
	default:
		return "libx264" // Fallback to software
	}
}

// isFFmpegEncoderAvailable checks if an FFmpeg encoder is available
func isFFmpegEncoderAvailable(encoder string) bool {
	output, err := execCommand("ffmpeg", "-hide_banner", "-encoders")
	if err != nil {
		return false
	}
	return contains(output, encoder)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && stringContains(s, substr)
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
