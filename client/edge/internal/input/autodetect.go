package input

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// AutoDetectInput automatically detects the best available input source
// Priority: V4L2 devices > Test pattern
func AutoDetectInput() (inputType, inputSource string, err error) {
	// Try to detect V4L2 devices first
	devices, err := detectV4L2Devices()
	if err == nil && len(devices) > 0 {
		return "v4l2", devices[0], nil
	}
	
	// Fallback to test pattern
	return "test", "pattern=smpte", nil
}

// detectV4L2Devices scans for available V4L2 video devices
func detectV4L2Devices() ([]string, error) {
	var devices []string
	
	// Check /dev/video* devices
	for i := 0; i < 10; i++ {
		device := fmt.Sprintf("/dev/video%d", i)
		if _, err := os.Stat(device); err == nil {
			// Verify it's a real camera using v4l2-ctl
			if isValidV4L2Device(device) {
				devices = append(devices, device)
			}
		}
	}
	
	if len(devices) == 0 {
		return nil, fmt.Errorf("no V4L2 devices found")
	}
	
	return devices, nil
}

// isValidV4L2Device checks if a device path is a valid V4L2 camera
func isValidV4L2Device(device string) bool {
	// Try v4l2-ctl first
	cmd := exec.Command("v4l2-ctl", "--device="+device, "--info")
	output, err := cmd.Output()
	if err != nil {
		// v4l2-ctl not available or device not valid
		// Try a simpler check: see if we can open the device
		file, err := os.Open(device)
		if err != nil {
			return false
		}
		defer file.Close()
		return true
	}
	
	// Check if output contains "Capability" which indicates it's a valid device
	return strings.Contains(string(output), "Driver Info") ||
		strings.Contains(string(output), "Capability")
}
