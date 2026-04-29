package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/lgc/pawstream/api/internal/domain/device"
	"github.com/lgc/pawstream/api/internal/domain/user"
)

// AdminHandler handles admin dashboard requests
type AdminHandler struct {
	deviceService *device.Service
	userService   *user.Service
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(deviceService *device.Service, userService *user.Service) *AdminHandler {
	return &AdminHandler{
		deviceService: deviceService,
		userService:   userService,
	}
}

// Dashboard handles GET /api/admin/dashboard
func (h *AdminHandler) Dashboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
	}

	// Check if user is admin (first registered user)
	isAdmin, err := h.userService.IsFirstUser(c.Context(), userID)
	if err != nil || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
			Error:   "forbidden",
			Message: "Admin access required",
		})
	}

	// Get all devices
	devices, err := h.deviceService.ListAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to list devices",
		})
	}

	// Build response
	onlineCount := 0
	adminDevices := make([]AdminDeviceInfo, 0, len(devices))
	for _, d := range devices {
		if d.IsOnline {
			onlineCount++
		}
		adminDevices = append(adminDevices, AdminDeviceInfo{
			ID:          d.ID,
			Name:        d.Name,
			Location:    d.Location,
			PublishPath: d.PublishPath,
			OwnerUserID: d.OwnerUserID,
			Disabled:    d.Disabled,
			IsOnline:    d.IsOnline,
			LastSeenAt:  d.LastSeenAt,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return c.JSON(AdminDashboard{
		TotalDevices:  len(devices),
		OnlineDevices: onlineCount,
		Devices:       adminDevices,
	})
}

// Heartbeat handles POST /api/devices/:id/heartbeat
func (h *AdminHandler) Heartbeat(c *fiber.Ctx) error {
	deviceID := c.Params("id")

	var req HeartbeatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid request body",
		})
	}

	// Verify device exists
	dev, err := h.deviceService.GetByID(c.Context(), deviceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "not_found",
			Message: "Device not found",
		})
	}

	// Update online status via publish path
	if err := h.deviceService.SetOnlineStatus(c.Context(), dev.PublishPath, true); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update status",
		})
	}

	return c.JSON(fiber.Map{"status": "ok"})
}
