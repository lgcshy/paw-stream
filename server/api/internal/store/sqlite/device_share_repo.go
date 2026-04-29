package sqlite

import (
	"context"
	"database/sql"
	"time"
)

// DeviceShare represents a device sharing record
type DeviceShare struct {
	ID              string
	DeviceID        string
	SharedByUserID  string
	SharedToUserID  string
	CreatedAt       time.Time
}

// DeviceShareRepository handles device share persistence
type DeviceShareRepository struct {
	db *DB
}

// NewDeviceShareRepository creates a new device share repository
func NewDeviceShareRepository(db *DB) *DeviceShareRepository {
	return &DeviceShareRepository{db: db}
}

// Create stores a new device share
func (r *DeviceShareRepository) Create(ctx context.Context, share *DeviceShare) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO device_shares (id, device_id, shared_by_user_id, shared_to_user_id, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		share.ID, share.DeviceID, share.SharedByUserID, share.SharedToUserID, share.CreatedAt,
	)
	return err
}

// Delete removes a device share
func (r *DeviceShareRepository) Delete(ctx context.Context, deviceID, sharedToUserID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM device_shares WHERE device_id = ? AND shared_to_user_id = ?`,
		deviceID, sharedToUserID,
	)
	return err
}

// ListByDevice lists all shares for a device
func (r *DeviceShareRepository) ListByDevice(ctx context.Context, deviceID string) ([]*DeviceShare, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, device_id, shared_by_user_id, shared_to_user_id, created_at
		 FROM device_shares WHERE device_id = ?`, deviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []*DeviceShare
	for rows.Next() {
		var s DeviceShare
		if err := rows.Scan(&s.ID, &s.DeviceID, &s.SharedByUserID, &s.SharedToUserID, &s.CreatedAt); err != nil {
			return nil, err
		}
		shares = append(shares, &s)
	}
	return shares, rows.Err()
}

// IsSharedWith checks if a device is shared with a specific user
func (r *DeviceShareRepository) IsSharedWith(ctx context.Context, deviceID, userID string) (bool, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM device_shares WHERE device_id = ? AND shared_to_user_id = ?`,
		deviceID, userID,
	).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteAllByDevice removes all shares for a device
func (r *DeviceShareRepository) DeleteAllByDevice(ctx context.Context, deviceID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM device_shares WHERE device_id = ?`, deviceID,
	)
	return err
}
