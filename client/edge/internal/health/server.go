package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lgc/pawstream/edge-client/internal/stream"
	"github.com/rs/zerolog"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status     string `json:"status"`       // streaming, stopped, error, reconnecting
	Uptime     int64  `json:"uptime"`       // seconds since start
	StreamURL  string `json:"stream_url"`
	InputType  string `json:"input_type"`
	ErrorCount int    `json:"error_count"`
	LastError  string `json:"last_error,omitempty"`
}

// Server provides HTTP health check endpoint
type Server struct {
	address   string
	streamMgr *stream.Manager
	startTime time.Time
	server    *http.Server
	logger    zerolog.Logger
}

// NewServer creates a new health check server
func NewServer(address string, streamMgr *stream.Manager, logger zerolog.Logger) *Server {
	return &Server{
		address:   address,
		streamMgr: streamMgr,
		startTime: time.Now(),
		logger:    logger.With().Str("component", "health-server").Logger(),
	}
}

// Start starts the health check server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)

	s.server = &http.Server{
		Addr:    s.address,
		Handler: mux,
	}

	s.logger.Info().Str("address", s.address).Msg("Health check server starting")

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop stops the health check server
func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	s.logger.Info().Msg("Stopping health check server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// handleHealth handles GET /health requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := s.streamMgr.Status()

	resp := HealthResponse{
		Status:     status.State,
		Uptime:     int64(time.Since(s.startTime).Seconds()),
		StreamURL:  status.OutputURL,
		InputType:  status.InputType,
		ErrorCount: status.ErrorCount,
		LastError:  status.LastError,
	}

	// Set HTTP status code based on stream state
	statusCode := http.StatusOK
	if resp.Status != "streaming" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.logger.Error().Err(err).Msg("Failed to encode health response")
	}
}
