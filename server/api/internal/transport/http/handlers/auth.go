package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/lgc/pawstream/api/internal/domain/user"
	"github.com/lgc/pawstream/api/internal/pkg/errors"
	"github.com/lgc/pawstream/api/internal/pkg/idgen"
	"github.com/lgc/pawstream/api/internal/pkg/jwtutil"
	"github.com/lgc/pawstream/api/internal/store/sqlite"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userService    *user.Service
	refreshRepo    *sqlite.RefreshTokenRepository
	jwtSecret      string
	jwtExpiry      time.Duration
	refreshExpiry  time.Duration
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService *user.Service, refreshRepo *sqlite.RefreshTokenRepository, jwtSecret string, jwtExpiry, refreshExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		userService:   userService,
		refreshRepo:   refreshRepo,
		jwtSecret:     jwtSecret,
		jwtExpiry:     jwtExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// Register handles POST /api/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "Invalid request body",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "validation_error",
			Message:   "Username and password are required",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Validate username format
	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "validation_error",
			Message:   "Username must be at least 3 characters",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Validate password strength
	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "validation_error",
			Message:   "Password must be at least 6 characters",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Set default nickname if not provided
	if req.Nickname == "" {
		req.Nickname = req.Username
	}

	// Register user
	newUser, err := h.userService.Register(c.Context(), user.CreateUserInput{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, errors.ErrDuplicateUsername) {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error:     "duplicate_username",
				Message:   "Username already exists",
				RequestID: c.Locals("request_id").(string),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to register user",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return user info (exclude password)
	return c.Status(fiber.StatusCreated).JSON(toUserInfo(newUser))
}

// Login handles POST /api/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "Invalid request body",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "validation_error",
			Message:   "Username and password are required",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Authenticate user
	authenticatedUser, err := h.userService.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, errors.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
				Error:     "invalid_credentials",
				Message:   "Invalid username or password",
				RequestID: c.Locals("request_id").(string),
			})
		}
		if errors.Is(err, errors.ErrUserDisabled) {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error:     "user_disabled",
				Message:   "Your account has been disabled",
				RequestID: c.Locals("request_id").(string),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Login failed",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Generate JWT access token
	token, err := jwtutil.GenerateToken(authenticatedUser.ID, authenticatedUser.Username, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to generate token",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Generate refresh token
	refreshToken, err := h.createRefreshToken(c, authenticatedUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to generate refresh token",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return token and user info
	info := toUserInfo(authenticatedUser)
	return c.JSON(LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         &info,
	})
}

// Refresh handles POST /api/refresh
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "validation_error",
			Message: "Refresh token is required",
		})
	}

	// Hash the token to look it up
	tokenHash := hashToken(req.RefreshToken)
	stored, err := h.refreshRepo.GetByTokenHash(c.Context(), tokenHash)
	if err != nil || stored == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "invalid_token",
			Message: "Invalid refresh token",
		})
	}

	// Check if revoked or expired
	if stored.Revoked || stored.ExpiresAt.Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "invalid_token",
			Message: "Refresh token expired or revoked",
		})
	}

	// Revoke the old refresh token (single use)
	_ = h.refreshRepo.Revoke(c.Context(), stored.ID)

	// Get user
	usr, err := h.userService.GetByID(c.Context(), stored.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:   "invalid_token",
			Message: "User not found",
		})
	}

	if usr.Disabled {
		return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
			Error:   "user_disabled",
			Message: "Your account has been disabled",
		})
	}

	// Generate new access token
	accessToken, err := jwtutil.GenerateToken(usr.ID, usr.Username, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate token",
		})
	}

	// Generate new refresh token
	newRefreshToken, err := h.createRefreshToken(c, usr.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate refresh token",
		})
	}

	info := toUserInfo(usr)
	return c.JSON(LoginResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
		User:         &info,
	})
}

// createRefreshToken generates and stores a new refresh token
func (h *AuthHandler) createRefreshToken(c *fiber.Ctx, userID string) (string, error) {
	raw, err := idgen.NewSecret(32)
	if err != nil {
		return "", err
	}

	tokenHash := hashToken(raw)
	rt := &sqlite.RefreshToken{
		ID:        idgen.NewUUID(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(h.refreshExpiry),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if err := h.refreshRepo.Create(c.Context(), rt); err != nil {
		return "", err
	}

	return raw, nil
}

// hashToken returns SHA-256 hex hash of a token string
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// GetMe handles GET /api/me
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	// Get user ID from JWT context (set by auth middleware)
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "User not authenticated",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get user info
	currentUser, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, errors.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:     "user_not_found",
				Message:   "User not found",
				RequestID: c.Locals("request_id").(string),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to get user info",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return user info
	return c.JSON(toUserInfo(currentUser))
}

// toUserInfo converts a domain user to a UserInfo response
func toUserInfo(u *user.User) UserInfo {
	info := UserInfo{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Disabled:  u.Disabled,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.AvatarPath != "" {
		info.AvatarURL = fmt.Sprintf("/api/avatars/%s", u.ID)
	}
	return info
}

const (
	avatarDir     = "data/avatars"
	maxAvatarSize = 2 * 1024 * 1024 // 2MB
)

var allowedAvatarTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// UploadAvatar handles POST /api/me/avatar
func (h *AuthHandler) UploadAvatar(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Error:     "unauthorized",
			Message:   "User not authenticated",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Get uploaded file
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "bad_request",
			Message:   "No avatar file provided",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Check file size
	if fileHeader.Size > maxAvatarSize {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "file_too_large",
			Message:   "Avatar file must be less than 2MB",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Check file type
	contentType := fileHeader.Header.Get("Content-Type")
	ext, ok := allowedAvatarTypes[contentType]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:     "invalid_file_type",
			Message:   "Avatar must be JPEG, PNG, or WebP",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Save file
	if err := os.MkdirAll(avatarDir, 0o755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to save avatar",
			RequestID: c.Locals("request_id").(string),
		})
	}

	avatarPath := filepath.Join(avatarDir, userID+ext)

	// Remove old avatar files (different extensions)
	for _, e := range allowedAvatarTypes {
		os.Remove(filepath.Join(avatarDir, userID+e))
	}

	if err := saveUploadedFile(fileHeader, avatarPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to save avatar",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Update database
	updatedUser, err := h.userService.UpdateAvatar(c.Context(), userID, avatarPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to update avatar",
			RequestID: c.Locals("request_id").(string),
		})
	}

	return c.JSON(toUserInfo(updatedUser))
}

// GetAvatar handles GET /api/avatars/:id
func (h *AuthHandler) GetAvatar(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Find avatar file
	for contentType, ext := range allowedAvatarTypes {
		avatarPath := filepath.Join(avatarDir, userID+ext)
		if _, err := os.Stat(avatarPath); err == nil {
			c.Set("Content-Type", contentType)
			c.Set("Cache-Control", "public, max-age=3600")
			return c.SendFile(avatarPath)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		Error:   "not_found",
		Message: "Avatar not found",
	})
}

// saveUploadedFile saves a multipart file to disk
func saveUploadedFile(fh *multipart.FileHeader, dst string) error {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
