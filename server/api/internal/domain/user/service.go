package user

import (
	"context"
	"time"

	"github.com/lgc/pawstream/api/internal/pkg/errors"
	"github.com/lgc/pawstream/api/internal/pkg/idgen"
	"github.com/lgc/pawstream/api/internal/pkg/password"
)

// Service handles user business logic
type Service struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, input CreateUserInput) (*User, error) {
	// Check if username already exists
	existing, err := s.repo.GetByUsername(ctx, input.Username)
	if err == nil && existing != nil {
		return nil, errors.ErrDuplicateUsername
	}

	// Hash password
	passwordHash, err := password.Hash(input.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}

	// Create user
	user := &User{
		ID:           idgen.NewUUID(),
		Username:     input.Username,
		Nickname:     input.Nickname,
		PasswordHash: passwordHash,
		Disabled:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	return user, nil
}

// Login authenticates a user by username and password
func (s *Service) Login(ctx context.Context, username, pwd string) (*User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Check if user is disabled
	if user.Disabled {
		return nil, errors.ErrUserDisabled
	}

	// Verify password
	if !password.Verify(pwd, user.PasswordHash) {
		return nil, errors.ErrInvalidCredentials
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// Update updates a user's information
func (s *Service) Update(ctx context.Context, id string, input UpdateUserInput) (*User, error) {
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.Nickname != nil {
		user.Nickname = *input.Nickname
	}

	if input.Password != nil {
		passwordHash, err := password.Hash(*input.Password)
		if err != nil {
			return nil, errors.Wrap(err, "failed to hash password")
		}
		user.PasswordHash = passwordHash
	}

	if input.Disabled != nil {
		user.Disabled = *input.Disabled
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	return user, nil
}
