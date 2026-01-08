package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// ConfigHandler handles configuration requests
type ConfigHandler struct {
	mediamtxWebRTCURL string
	mediamtxRTSPURL   string
}

// NewConfigHandler creates a new config handler
func NewConfigHandler(mediamtxWebRTCURL, mediamtxRTSPURL string) *ConfigHandler {
	return &ConfigHandler{
		mediamtxWebRTCURL: mediamtxWebRTCURL,
		mediamtxRTSPURL:   mediamtxRTSPURL,
	}
}

// ConfigResponse represents the client configuration
type ConfigResponse struct {
	MediaMTX MediaMTXConfig `json:"mediamtx"`
}

// MediaMTXConfig contains MediaMTX service URLs
type MediaMTXConfig struct {
	WebRTCURL string `json:"webrtc_url"`
	RTSPURL   string `json:"rtsp_url"`
}

// GetConfig handles GET /api/config
// Returns client configuration including MediaMTX URLs
func (h *ConfigHandler) GetConfig(c *fiber.Ctx) error {
	return c.JSON(ConfigResponse{
		MediaMTX: MediaMTXConfig{
			WebRTCURL: h.mediamtxWebRTCURL,
			RTSPURL:   h.mediamtxRTSPURL,
		},
	})
}
