package config

import (
	"fmt"
)

// ApplyPreset applies a predefined configuration preset
func (c *Config) ApplyPreset(preset string) error {
	switch preset {
	case "low-latency":
		return c.applyLowLatencyPreset()
	case "high-quality":
		return c.applyHighQualityPreset()
	case "balanced":
		return c.applyBalancedPreset()
	case "power-save":
		return c.applyPowerSavePreset()
	default:
		return fmt.Errorf("unknown preset: %s", preset)
	}
}

// applyLowLatencyPreset optimizes for minimal latency (100-200ms)
// Best for: monitoring, real-time interaction
func (c *Config) applyLowLatencyPreset() error {
	// Use GStreamer for lowest latency
	c.Stream.Engine = "gstreamer"
	c.Stream.GStreamer.LatencyMs = 100
	c.Stream.GStreamer.UseHardware = true
	c.Stream.GStreamer.BufferSize = 500 // Smaller buffer for lower latency

	// Video settings optimized for low latency
	c.Video.Bitrate = 2000000  // 2 Mbps
	c.Video.Framerate = 30
	c.Video.Width = 1280
	c.Video.Height = 720

	return nil
}

// applyHighQualityPreset optimizes for maximum quality
// Best for: recording, archiving
func (c *Config) applyHighQualityPreset() error {
	// Use FFmpeg with slow preset for best quality
	c.Stream.Engine = "ffmpeg"
	c.Stream.FFmpeg.Preset = "slow"
	c.Stream.FFmpeg.Tune = "film"
	c.Stream.FFmpeg.HWAccel = "none" // Software encoding for best quality

	// Video settings optimized for quality
	c.Video.Bitrate = 5000000  // 5 Mbps
	c.Video.Framerate = 60
	c.Video.Width = 1920
	c.Video.Height = 1080

	return nil
}

// applyBalancedPreset provides balance between quality and latency
// Best for: general purpose streaming
func (c *Config) applyBalancedPreset() error {
	// Use FFmpeg with medium settings
	c.Stream.Engine = "ffmpeg"
	c.Stream.FFmpeg.Preset = "medium"
	c.Stream.FFmpeg.Tune = "zerolatency"
	c.Stream.FFmpeg.HWAccel = "auto"

	// Video settings balanced
	c.Video.Bitrate = 2500000  // 2.5 Mbps
	c.Video.Framerate = 30
	c.Video.Width = 1280
	c.Video.Height = 720

	return nil
}

// applyPowerSavePreset optimizes for minimal resource usage
// Best for: edge devices, battery-powered devices
func (c *Config) applyPowerSavePreset() error {
	// Prefer hardware encoding to save CPU
	c.Stream.Engine = "ffmpeg"
	c.Stream.FFmpeg.Preset = "ultrafast"
	c.Stream.FFmpeg.Tune = "zerolatency"
	c.Stream.FFmpeg.HWAccel = "auto" // Use hardware if available

	// Video settings optimized for low resource usage
	c.Video.Bitrate = 1000000  // 1 Mbps
	c.Video.Framerate = 15     // Lower framerate
	c.Video.Width = 854
	c.Video.Height = 480       // Lower resolution

	return nil
}

// GetPresets returns a list of available presets with descriptions
func GetPresets() map[string]string {
	return map[string]string{
		"low-latency":  "优化低延迟（100-200ms）- 适合监控、实时互动",
		"high-quality": "优化高质量 - 适合录制、存档",
		"balanced":     "平衡质量和延迟 - 适合通用场景",
		"power-save":   "优化资源占用 - 适合边缘设备、省电模式",
	}
}
