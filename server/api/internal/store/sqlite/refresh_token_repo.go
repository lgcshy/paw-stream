package sqlite

import (
	"context"
	"database/sql"
	"time"
)

// RefreshToken represents a stored refresh token
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

// RefreshTokenRepository handles refresh token persistence
type RefreshTokenRepository struct {
	db *DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create stores a new refresh token
func (r *RefreshTokenRepository) Create(ctx context.Context, token *RefreshToken) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.Revoked, token.CreatedAt,
	)
	return err
}

// GetByTokenHash retrieves a refresh token by its hash
func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	var t RefreshToken
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token_hash, expires_at, revoked, created_at
		 FROM refresh_tokens WHERE token_hash = ?`, tokenHash,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.Revoked, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Revoke marks a refresh token as revoked
func (r *RefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked = 1 WHERE id = ?`, id,
	)
	return err
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked = 1 WHERE user_id = ?`, userID,
	)
	return err
}

// DeleteExpired removes expired refresh tokens
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM refresh_tokens WHERE expires_at < ?`, time.Now(),
	)
	return err
}
