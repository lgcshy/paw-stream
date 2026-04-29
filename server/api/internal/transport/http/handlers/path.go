package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/lgc/pawstream/api/internal/domain/device"
)

// PathHandler handles path query requests
type PathHandler struct {
	deviceService *device.Service
}

// NewPathHandler creates a new path handler
func NewPathHandler(deviceService *device.Service) *PathHandler {
	return &PathHandler{
		deviceService: deviceService,
	}
}

// List handles GET /api/paths
func (h *PathHandler) List(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "User not authenticated",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get user's own devices
	devices, err := h.deviceService.ListByOwner(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to list paths",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Convert to path info (only enabled devices)
	paths := make([]*PathInfo, 0)
	for _, d := range devices {
		if !d.Disabled {
			paths = append(paths, &PathInfo{
				PublishPath:    d.PublishPath,
				DeviceID:       d.ID,
				DeviceName:     d.Name,
				DeviceLocation: d.Location,
			})
		}
	}

	// Also include shared devices
	sharedDevices, err := h.deviceService.ListSharedWith(c.Context(), userID)
	if err == nil {
		for _, d := range sharedDevices {
			if !d.Disabled {
				paths = append(paths, &PathInfo{
					PublishPath:    d.PublishPath,
					DeviceID:       d.ID,
					DeviceName:     d.Name,
					DeviceLocation: d.Location,
				})
			}
		}
	}

	return c.JSON(paths)
}
