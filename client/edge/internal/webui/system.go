package webui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemInfo contains system information
type SystemInfo struct {
	Timestamp   time.Time      `json:"timestamp"`
	Host        HostInfo       `json:"host"`
	CPU         CPUInfo        `json:"cpu"`
	Memory      MemoryInfo     `json:"memory"`
	Disk        DiskInfo       `json:"disk"`
	Process     ProcessInfo    `json:"process"`
	InputSources []InputSource `json:"input_sources,omitempty"`
}

// HostInfo contains host information
type HostInfo struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	Uptime          uint64 `json:"uptime"`
	UptimeFormatted string `json:"uptime_formatted"`
}

// CPUInfo contains CPU information
type CPUInfo struct {
	Cores      int     `json:"cores"`
	ModelName  string  `json:"model_name"`
	UsageTotal float64 `json:"usage_total"`
	UsagePerCPU []float64 `json:"usage_per_cpu,omitempty"`
}

// MemoryInfo contains memory information
type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

// DiskInfo contains disk information
type DiskInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

// ProcessInfo contains current process information
type ProcessInfo struct {
	PID         int32   `json:"pid"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryMB    uint64  `json:"memory_mb"`
	NumThreads  int32   `json:"num_threads"`
	NumGoroutine int    `json:"num_goroutine"`
}

// InputSource represents a video input source
type InputSource struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Device      string `json:"device"`
	Description string `json:"description"`
}

// GetSystemInfo collects and returns system information
func GetSystemInfo() SystemInfo {
	info := SystemInfo{
		Timestamp: time.Now(),
	}

	// Host info
	if hostInfo, err := host.Info(); err == nil {
		info.Host = HostInfo{
			Hostname:        hostInfo.Hostname,
			OS:              hostInfo.OS,
			Platform:        hostInfo.Platform,
			PlatformVersion: hostInfo.PlatformVersion,
			Uptime:          hostInfo.Uptime,
			UptimeFormatted: formatUptime(hostInfo.Uptime),
		}
	}

	// CPU info
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		info.CPU.ModelName = cpuInfo[0].ModelName
	}
	info.CPU.Cores = runtime.NumCPU()
	
	if cpuPercent, err := cpu.Percent(0, false); err == nil && len(cpuPercent) > 0 {
		info.CPU.UsageTotal = cpuPercent[0]
	}

	// Memory info
	if memInfo, err := mem.VirtualMemory(); err == nil {
		info.Memory = MemoryInfo{
			Total:       memInfo.Total,
			Available:   memInfo.Available,
			Used:        memInfo.Used,
			UsedPercent: memInfo.UsedPercent,
		}
	}

	// Disk info (root partition)
	if diskInfo, err := disk.Usage("/"); err == nil {
		info.Disk = DiskInfo{
			Total:       diskInfo.Total,
			Free:        diskInfo.Free,
			Used:        diskInfo.Used,
			UsedPercent: diskInfo.UsedPercent,
		}
	}

	// Process info
	if proc, err := process.NewProcess(int32(os.Getpid())); err == nil {
		info.Process.PID = proc.Pid

		if cpuPercent, err := proc.CPUPercent(); err == nil {
			info.Process.CPUPercent = cpuPercent
		}

		if memInfo, err := proc.MemoryInfo(); err == nil {
			info.Process.MemoryMB = memInfo.RSS / 1024 / 1024
		}

		if numThreads, err := proc.NumThreads(); err == nil {
			info.Process.NumThreads = numThreads
		}
	}
	info.Process.NumGoroutine = runtime.NumGoroutine()

	return info
}

// formatUptime formats uptime in seconds to human-readable string
func formatUptime(seconds uint64) string {
	d := time.Duration(seconds) * time.Second
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// DetectInputSources detects available video input sources
func DetectInputSources() []InputSource {
	var sources []InputSource

	// Always add test source
	sources = append(sources, InputSource{
		Type:        "testsrc",
		Name:        "Test Pattern",
		Device:      "testsrc",
		Description: "Built-in test pattern generator",
	})

	// Detect V4L2 devices (Linux)
	if runtime.GOOS == "linux" {
		v4l2Devices := detectV4L2Devices()
		sources = append(sources, v4l2Devices...)
	}

	// TODO: Add detection for other platforms (Windows, macOS)
	// TODO: Add RTSP source detection/validation

	return sources
}

// detectV4L2Devices detects V4L2 video devices on Linux
func detectV4L2Devices() []InputSource {
	var sources []InputSource

	// Check if v4l2-ctl is available
	if _, err := exec.LookPath("v4l2-ctl"); err != nil {
		return sources
	}

	// List V4L2 devices
	cmd := exec.Command("v4l2-ctl", "--list-devices")
	output, err := cmd.Output()
	if err != nil {
		return sources
	}

	// Parse output
	lines := strings.Split(string(output), "\n")
	var currentName string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "/dev/video") {
			if currentName != "" {
				sources = append(sources, InputSource{
					Type:        "v4l2",
					Name:        currentName,
					Device:      line,
					Description: fmt.Sprintf("V4L2 device: %s", currentName),
				})
			}
		} else if !strings.Contains(line, "(") {
			currentName = line
		}
	}

	return sources
}
