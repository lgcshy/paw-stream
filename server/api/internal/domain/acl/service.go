package acl

import (
	"context"
	"fmt"

	"github.com/lgc/pawstream/api/internal/domain/device"
	"github.com/lgc/pawstream/api/internal/domain/user"
	"github.com/lgc/pawstream/api/internal/pkg/errors"
	"github.com/lgc/pawstream/api/internal/pkg/jwtutil"
	"github.com/lgc/pawstream/api/internal/store/sqlite"
)

// Service handles access control logic
type Service struct {
	userService   *user.Service
	deviceService *device.Service
	shareRepo     *sqlite.DeviceShareRepository
	jwtSecret     string
}

// NewService creates a new ACL service
func NewService(userService *user.Service, deviceService *device.Service, shareRepo *sqlite.DeviceShareRepository, jwtSecret string) *Service {
	return &Service{
		userService:   userService,
		deviceService: deviceService,
		shareRepo:     shareRepo,
		jwtSecret:     jwtSecret,
	}
}

// Authorize checks if an action is allowed
func (s *Service) Authorize(ctx context.Context, req AuthRequest) (*AuthResult, error) {
	switch req.Action {
	case ActionPublish:
		return s.authorizePublish(ctx, req)
	case ActionRead, ActionPlayback:
		return s.authorizeRead(ctx, req)
	default:
		return &AuthResult{
			Allowed: false,
			Reason:  fmt.Sprintf("unknown action: %s", req.Action),
		}, nil
	}
}

// authorizePublish checks if a device can publish to a path
func (s *Service) authorizePublish(ctx context.Context, req AuthRequest) (*AuthResult, error) {
	// Get device by publish path
	dev, err := s.deviceService.GetByPublishPath(ctx, req.Path)
	if err != nil {
		if errors.Is(err, errors.ErrDeviceNotFound) {
			return &AuthResult{
				Allowed: false,
				Reason:  "device not found for path",
			}, nil
		}
		return nil, err
	}

	// Check if device is disabled
	if dev.Disabled {
		return &AuthResult{
			Allowed: false,
			Reason:  "device is disabled",
		}, nil
	}

	// Verify device secret (from password field or user field)
	secret := req.Password
	if secret == "" {
		secret = req.User
	}

	if secret == "" {
		return &AuthResult{
			Allowed: false,
			Reason:  "missing device secret",
		}, nil
	}

	// Verify secret
	valid, err := s.deviceService.VerifySecret(ctx, dev.ID, secret)
	if err != nil {
		return nil, err
	}

	if !valid {
		return &AuthResult{
			Allowed: false,
			Reason:  "invalid device secret",
		}, nil
	}

	return &AuthResult{
		Allowed: true,
		Reason:  "authorized",
	}, nil
}

// authorizeRead checks if a user can read/playback a stream
func (s *Service) authorizeRead(ctx context.Context, req AuthRequest) (*AuthResult, error) {
	// Extract token (from Token field or User field)
	token := req.Token
	if token == "" {
		token = req.User
	}

	if token == "" {
		return &AuthResult{
			Allowed: false,
			Reason:  "missing user token",
		}, nil
	}

	// Validate JWT token
	claims, err := jwtutil.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return &AuthResult{
			Allowed: false,
			Reason:  "invalid token",
		}, nil
	}

	// Get user
	usr, err := s.userService.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, errors.ErrUserNotFound) {
			return &AuthResult{
				Allowed: false,
				Reason:  "user not found",
			}, nil
		}
		return nil, err
	}

	// Check if user is disabled
	if usr.Disabled {
		return &AuthResult{
			Allowed: false,
			Reason:  "user is disabled",
		}, nil
	}

	// Check if user owns a device with this path
	devices, err := s.deviceService.ListByOwner(ctx, usr.ID)
	if err != nil {
		return nil, err
	}

	for _, dev := range devices {
		if dev.PublishPath == req.Path && !dev.Disabled {
			return &AuthResult{
				Allowed: true,
				Reason:  "authorized",
			}, nil
		}
	}

	// Check if device is shared with user
	dev, err := s.deviceService.GetByPublishPath(ctx, req.Path)
	if err == nil && dev != nil && !dev.Disabled && s.shareRepo != nil {
		shared, err := s.shareRepo.IsSharedWith(ctx, dev.ID, usr.ID)
		if err == nil && shared {
			return &AuthResult{
				Allowed: true,
				Reason:  "authorized via share",
			}, nil
		}
	}

	return &AuthResult{
		Allowed: false,
		Reason:  "user does not own device for this path",
	}, nil
}
