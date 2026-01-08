package user

import "context"

// Repository defines the interface for user data access
type Repository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error

	// List retrieves all users (for admin purposes)
	List(ctx context.Context, limit, offset int) ([]*User, error)
}
