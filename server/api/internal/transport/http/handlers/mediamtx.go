package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/lgc/pawstream/api/internal/domain/acl"
)

// MediaMTXHandler handles MediaMTX authentication callbacks
type MediaMTXHandler struct {
	aclService *acl.Service
	log        zerolog.Logger
}

// NewMediaMTXHandler creates a new MediaMTX handler
func NewMediaMTXHandler(aclService *acl.Service, log zerolog.Logger) *MediaMTXHandler {
	return &MediaMTXHandler{
		aclService: aclService,
		log:        log,
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

	// Log auth attempt
	h.log.Info().
		Str("action", req.Action).
		Str("path", req.Path).
		Str("protocol", req.Protocol).
		Str("ip", req.IP).
		Str("user", req.User).
		Str("password", req.Password).
		Str("token", req.Token).
		Msg("mediamtx auth request")

	// Convert to ACL request
	aclReq := acl.AuthRequest{
		Action:   acl.Action(req.Action),
		Path:     req.Path,
		Protocol: req.Protocol,
		IP:       req.IP,
		User:     req.User,
		Password: req.Password,
		Token:    req.Token,
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
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error":   "forbidden",
		"message": result.Reason,
	})
}
