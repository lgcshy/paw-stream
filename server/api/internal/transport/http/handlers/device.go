package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/lgc/pawstream/api/internal/domain/device"
	"github.com/lgc/pawstream/api/internal/pkg/errors"
)

// DeviceHandler handles device management requests
type DeviceHandler struct {
	deviceService *device.Service
}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler(deviceService *device.Service) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// AuthDevice handles POST /api/device/auth - for edge client authentication
func (h *DeviceHandler) AuthDevice(c *fiber.Ctx) error {
	ctx := c.Context()
	
	// Parse request
	var req struct {
		DeviceID string `json:"device_id"`
		Secret   string `json:"secret"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "Invalid request body",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Validate input
	if req.DeviceID == "" || req.Secret == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "validation_error",
			Message:   "Device ID and secret are required",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Verify device credentials
	valid, err := h.deviceService.VerifySecret(ctx, req.DeviceID, req.Secret)
	if err != nil {
		if errors.Is(err, errors.ErrDeviceNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
				Error:     "unauthorized",
				Message:   "Invalid device credentials",
				RequestID: c.Locals("request_id").(string),
			})
		}
		if errors.Is(err, errors.ErrDeviceDisabled) {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error:     "forbidden",
				Message:   "Device is disabled",
				RequestID: c.Locals("request_id").(string),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to authenticate device",
			RequestID: c.Locals("request_id").(string),
		})
	}

	if !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "Invalid device credentials",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get device info
	dev, err := h.deviceService.GetByID(ctx, req.DeviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to get device info",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return device info
	return c.JSON(fiber.Map{
		"id":           dev.ID,
		"name":         dev.Name,
		"location":     dev.Location,
		"publish_path": dev.PublishPath,
		"disabled":     dev.Disabled,
		"created_at":   dev.CreatedAt,
		"updated_at":   dev.UpdatedAt,
	})
}

// Create handles POST /api/devices
func (h *DeviceHandler) Create(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "User not authenticated",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Parse request
	var req CreateDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "Invalid request body",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Validate input
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "validation_error",
			Message:   "Device name is required",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Create device
	newDevice, deviceSecret, err := h.deviceService.Create(c.Context(), device.CreateDeviceInput{
		OwnerUserID: userID,
		Name:        req.Name,
		Location:    req.Location,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to create device",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return device info with secret (only once!)
	return c.Status(fiber.StatusCreated).JSON(CreateDeviceResponse{
		Device: &DeviceInfo{
			ID:          newDevice.ID,
			Name:        newDevice.Name,
			Location:    newDevice.Location,
			PublishPath: newDevice.PublishPath,
			Disabled:    newDevice.Disabled,
			CreatedAt:   newDevice.CreatedAt,
			UpdatedAt:   newDevice.UpdatedAt,
		},
		Secret: deviceSecret.Secret,
	})
}

// List handles GET /api/devices
func (h *DeviceHandler) List(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "User not authenticated",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get user's devices
	devices, err := h.deviceService.ListByOwner(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to list devices",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Convert to response format (exclude secrets)
	deviceInfos := make([]*DeviceInfo, len(devices))
	for i, d := range devices {
		deviceInfos[i] = &DeviceInfo{
			ID:          d.ID,
			Name:        d.Name,
			Location:    d.Location,
			PublishPath: d.PublishPath,
			Disabled:    d.Disabled,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		}
	}

	return c.JSON(deviceInfos)
}

// Get handles GET /api/devices/:id
func (h *DeviceHandler) Get(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "User not authenticated",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get device ID from URL
	deviceID := c.Params("id")
	if deviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "Device ID is required",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get device
	dev, err := h.deviceService.GetByID(c.Context(), deviceID)
	if err != nil {
		if errors.Is(err, errors.ErrDeviceNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:     "device_not_found",
				Message:   "Device not found",
				RequestID: c.Locals("request_id").(string),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to get device",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Check ownership
	if dev.OwnerUserID != userID {
		// Return 404 to not reveal device existence
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:     "device_not_found",
			Message:   "Device not found",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return device info (exclude secret)
	return c.JSON(DeviceInfo{
		ID:          dev.ID,
		Name:        dev.Name,
		Location:    dev.Location,
		PublishPath: dev.PublishPath,
		Disabled:    dev.Disabled,
		CreatedAt:   dev.CreatedAt,
		UpdatedAt:   dev.UpdatedAt,
	})
}

// Update handles PUT /api/devices/:id
func (h *DeviceHandler) Update(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "User not authenticated",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get device ID from URL
	deviceID := c.Params("id")
	if deviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "Device ID is required",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Check ownership first
	dev, err := h.deviceService.GetByID(c.Context(), deviceID)
	if err != nil {
		if errors.Is(err, errors.ErrDeviceNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:     "device_not_found",
				Message:   "Device not found",
				RequestID: c.Locals("request_id").(string),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to get device",
			RequestID: c.Locals("request_id").(string),
		})
	}

	if dev.OwnerUserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:     "device_not_found",
			Message:   "Device not found",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Parse request
	var req UpdateDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "Invalid request body",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Update device
	updatedDevice, err := h.deviceService.Update(c.Context(), deviceID, device.UpdateDeviceInput{
		Name:     req.Name,
		Location: req.Location,
		Disabled: req.Disabled,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to update device",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return updated device info
	return c.JSON(DeviceInfo{
		ID:          updatedDevice.ID,
		Name:        updatedDevice.Name,
		Location:    updatedDevice.Location,
		PublishPath: updatedDevice.PublishPath,
		Disabled:    updatedDevice.Disabled,
		CreatedAt:   updatedDevice.CreatedAt,
		UpdatedAt:   updatedDevice.UpdatedAt,
	})
}

// Delete handles DELETE /api/devices/:id
func (h *DeviceHandler) Delete(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "User not authenticated",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get device ID from URL
	deviceID := c.Params("id")
	if deviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "Device ID is required",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Check ownership first
	dev, err := h.deviceService.GetByID(c.Context(), deviceID)
	if err != nil {
		if errors.Is(err, errors.ErrDeviceNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:     "device_not_found",
				Message:   "Device not found",
				RequestID: c.Locals("request_id").(string),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to get device",
			RequestID: c.Locals("request_id").(string),
		})
	}

	if dev.OwnerUserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:     "device_not_found",
			Message:   "Device not found",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Delete device
	if err := h.deviceService.Delete(c.Context(), deviceID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to delete device",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return 204 No Content
	return c.SendStatus(fiber.StatusNoContent)
}

// RotateSecret handles POST /api/devices/:id/rotate-secret
func (h *DeviceHandler) RotateSecret(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "User not authenticated",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get device ID from URL
	deviceID := c.Params("id")
	if deviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "Device ID is required",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Check ownership first
	dev, err := h.deviceService.GetByID(c.Context(), deviceID)
	if err != nil {
		if errors.Is(err, errors.ErrDeviceNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:     "device_not_found",
				Message:   "Device not found",
				RequestID: c.Locals("request_id").(string),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to get device",
			RequestID: c.Locals("request_id").(string),
		})
	}

	if dev.OwnerUserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:     "device_not_found",
			Message:   "Device not found",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Rotate secret
	newSecret, err := h.deviceService.RotateSecret(c.Context(), deviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to rotate secret",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return new secret (only once!)
	return c.JSON(RotateSecretResponse{
		Secret:        newSecret.Secret,
		SecretVersion: dev.SecretVersion + 1,
	})
}
