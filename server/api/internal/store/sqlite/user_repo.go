package sqlite

import (
	"context"
	"database/sql"

	"github.com/lgc/pawstream/api/internal/domain/user"
	"github.com/lgc/pawstream/api/internal/pkg/errors"
)

// UserRepository implements user.Repository for SQLite
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new SQLite user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (id, username, nickname, password_hash, avatar_path, disabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Username, u.Nickname, u.PasswordHash, u.AvatarPath,
		u.Disabled, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT id, username, nickname, password_hash, avatar_path, disabled, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	var u user.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Nickname, &u.PasswordHash, &u.AvatarPath,
		&u.Disabled, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by ID")
	}
	return &u, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
		SELECT id, username, nickname, password_hash, avatar_path, disabled, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	var u user.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Nickname, &u.PasswordHash, &u.AvatarPath,
		&u.Disabled, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by username")
	}
	return &u, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
		SET nickname = ?, password_hash = ?, avatar_path = ?, disabled = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		u.Nickname, u.PasswordHash, u.AvatarPath, u.Disabled, u.UpdatedAt, u.ID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}
	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}
	return nil
}

// List retrieves all users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	query := `
		SELECT id, username, nickname, password_hash, avatar_path, disabled, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list users")
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Nickname, &u.PasswordHash, &u.AvatarPath,
			&u.Disabled, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan user")
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating users")
	}

	return users, nil
}
