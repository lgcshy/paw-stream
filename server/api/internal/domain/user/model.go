package user

import "time"

// User represents a business user in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	AvatarPath   string    `json:"avatar_path,omitempty"`
	Disabled     bool      `json:"disabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateUserInput represents input for creating a new user
type CreateUserInput struct {
	Username string
	Nickname string
	Password string
}

// UpdateUserInput represents input for updating a user
type UpdateUserInput struct {
	Nickname *string
	Password *string
	Disabled *bool
}
