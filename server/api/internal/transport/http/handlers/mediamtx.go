package handlers

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/lgc/pawstream/api/internal/domain/acl"
	"github.com/lgc/pawstream/api/internal/domain/device"
)

// MediaMTXHandler handles MediaMTX authentication callbacks
type MediaMTXHandler struct {
	aclService    *acl.Service
	deviceService *device.Service
	log           zerolog.Logger
}

// NewMediaMTXHandler creates a new MediaMTX handler
func NewMediaMTXHandler(aclService *acl.Service, deviceService *device.Service, log zerolog.Logger) *MediaMTXHandler {
	return &MediaMTXHandler{
		aclService:    aclService,
		deviceService: deviceService,
		log:           log,
	}
}

// AuthRequest represents the MediaMTX auth callback request
type AuthRequest struct {
	Action   string `json:"action"`
	Path     string `json:"path"`
	Protocol string `json:"protocol"`
	IP       string `json:"ip"`
	User     string `json:"user"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Query    string `json:"query"`
}

// Auth handles POST /mediamtx/auth
func (h *MediaMTXHandler) Auth(c *fiber.Ctx) error {
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error().Err(err).Msg("failed to parse mediamtx auth request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "bad_request",
			"message": "Invalid request format",
		})
	}

	// Extract JWT from query parameters if token field is empty
	token := req.Token
	if token == "" && req.Query != "" {
		if params, err := url.ParseQuery(req.Query); err == nil {
			if jwt := params.Get("jwt"); jwt != "" {
				token = jwt
			}
		}
	}

	// Log auth attempt
	h.log.Info().
		Str("action", req.Action).
		Str("path", req.Path).
		Str("protocol", req.Protocol).
		Str("ip", req.IP).
		Str("user", req.User).
		Str("password", req.Password).
		Str("token", token).
		Str("query", req.Query).
		Msg("mediamtx auth request")

	// Convert to ACL request
	aclReq := acl.AuthRequest{
		Action:   acl.Action(req.Action),
		Path:     req.Path,
		Protocol: req.Protocol,
		IP:       req.IP,
		User:     req.User,
		Password: req.Password,
		Token:    token,
	}

	// Check authorization
	result, err := h.aclService.Authorize(c.Context(), aclReq)
	if err != nil {
		h.log.Error().Err(err).Msg("authorization check failed")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "Authorization check failed",
		})
	}

	// Log result
	h.log.Info().
		Str("action", req.Action).
		Str("path", req.Path).
		Bool("allowed", result.Allowed).
		Str("reason", result.Reason).
		Msg("authorization result")

	// Return result
	if result.Allowed {
		// Update device online status on publish
		if req.Action == "publish" {
			if err := h.deviceService.SetOnlineStatus(c.Context(), req.Path, true); err != nil {
				h.log.Error().Err(err).Str("path", req.Path).Msg("failed to set device online")
			}
		}
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error":   "forbidden",
		"message": result.Reason,
	})
}

// StreamClosed handles POST /internal/stream-closed (called when a stream ends)
func (h *MediaMTXHandler) StreamClosed(c *fiber.Ctx) error {
	var req struct {
		Path string `json:"path"`
	}
	if err := c.BodyParser(&req); err != nil || req.Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad_request",
		})
	}

	if err := h.deviceService.SetOnlineStatus(c.Context(), req.Path, false); err != nil {
		h.log.Error().Err(err).Str("path", req.Path).Msg("failed to set device offline")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal_error",
		})
	}

	h.log.Info().Str("path", req.Path).Msg("device stream closed, set offline")
	return c.SendStatus(fiber.StatusOK)
}
