package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/lgc/pawstream/api/internal/domain/user"
	"github.com/lgc/pawstream/api/internal/pkg/errors"
	"github.com/lgc/pawstream/api/internal/pkg/jwtutil"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userService *user.Service
	jwtSecret   string
	jwtExpiry   time.Duration
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService *user.Service, jwtSecret string, jwtExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
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
	return c.Status(fiber.StatusCreated).JSON(UserInfo{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Nickname:  newUser.Nickname,
		Disabled:  newUser.Disabled,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	})
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

	// Generate JWT token
	token, err := jwtutil.GenerateToken(authenticatedUser.ID, authenticatedUser.Username, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:     "internal_error",
			Message:   "Failed to generate token",
			RequestID: c.Locals("request_id").(string),
		})
	}

	// Return token and user info
	return c.JSON(LoginResponse{
		Token: token,
		User: &UserInfo{
			ID:        authenticatedUser.ID,
			Username:  authenticatedUser.Username,
			Nickname:  authenticatedUser.Nickname,
			Disabled:  authenticatedUser.Disabled,
			CreatedAt: authenticatedUser.CreatedAt,
			UpdatedAt: authenticatedUser.UpdatedAt,
		},
	})
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
	return c.JSON(UserInfo{
		ID:        currentUser.ID,
		Username:  currentUser.Username,
		Nickname:  currentUser.Nickname,
		Disabled:  currentUser.Disabled,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	})
}
