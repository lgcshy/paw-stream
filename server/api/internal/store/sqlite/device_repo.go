package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/lgc/pawstream/api/internal/domain/device"
	"github.com/lgc/pawstream/api/internal/pkg/errors"
)

const deviceColumns = `id, owner_user_id, name, location, publish_path,
	secret_hash, secret_cipher, secret_version,
	disabled, is_online, last_seen_at, created_at, updated_at`

// DeviceRepository implements device.Repository for SQLite
type DeviceRepository struct {
	db *DB
}

// NewDeviceRepository creates a new SQLite device repository
func NewDeviceRepository(db *DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func scanDevice(row interface{ Scan(...any) error }) (*device.Device, error) {
	var d device.Device
	err := row.Scan(
		&d.ID, &d.OwnerUserID, &d.Name, &d.Location, &d.PublishPath,
		&d.SecretHash, &d.SecretCipher, &d.SecretVersion,
		&d.Disabled, &d.IsOnline, &d.LastSeenAt, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Create creates a new device
func (r *DeviceRepository) Create(ctx context.Context, d *device.Device) error {
	query := `
		INSERT INTO devices (
			id, owner_user_id, name, location, publish_path,
			secret_hash, secret_cipher, secret_version,
			disabled, is_online, last_seen_at, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		d.ID, d.OwnerUserID, d.Name, d.Location, d.PublishPath,
		d.SecretHash, d.SecretCipher, d.SecretVersion,
		d.Disabled, d.IsOnline, d.LastSeenAt, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create device")
	}
	return nil
}

// GetByID retrieves a device by ID
func (r *DeviceRepository) GetByID(ctx context.Context, id string) (*device.Device, error) {
	query := `SELECT ` + deviceColumns + ` FROM devices WHERE id = ?`
	d, err := scanDevice(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get device by ID")
	}
	return d, nil
}

// GetByPublishPath retrieves a device by its publish path
func (r *DeviceRepository) GetByPublishPath(ctx context.Context, path string) (*device.Device, error) {
	query := `SELECT ` + deviceColumns + ` FROM devices WHERE publish_path = ?`
	d, err := scanDevice(r.db.QueryRowContext(ctx, query, path))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get device by path")
	}
	return d, nil
}

// ListByOwner retrieves all devices owned by a user
func (r *DeviceRepository) ListByOwner(ctx context.Context, ownerUserID string) ([]*device.Device, error) {
	query := `SELECT ` + deviceColumns + ` FROM devices WHERE owner_user_id = ? ORDER BY created_at DESC`
	return r.queryDevices(ctx, query, ownerUserID)
}

// ListAll retrieves all devices (admin use)
func (r *DeviceRepository) ListAll(ctx context.Context) ([]*device.Device, error) {
	query := `SELECT ` + deviceColumns + ` FROM devices ORDER BY created_at DESC`
	return r.queryDevices(ctx, query)
}

// ListSharedWith retrieves all devices shared with a user
func (r *DeviceRepository) ListSharedWith(ctx context.Context, userID string) ([]*device.Device, error) {
	query := `SELECT ` + deviceColumns + `
		FROM devices d
		INNER JOIN device_shares s ON d.id = s.device_id
		WHERE s.shared_to_user_id = ?
		ORDER BY d.created_at DESC`
	return r.queryDevices(ctx, query, userID)
}

func (r *DeviceRepository) queryDevices(ctx context.Context, query string, args ...any) ([]*device.Device, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query devices")
	}
	defer rows.Close()

	var devices []*device.Device
	for rows.Next() {
		d, err := scanDevice(rows)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan device")
		}
		if d != nil {
			devices = append(devices, d)
		}
	}
	return devices, rows.Err()
}

// Update updates an existing device
func (r *DeviceRepository) Update(ctx context.Context, d *device.Device) error {
	query := `
		UPDATE devices
		SET name = ?, location = ?, secret_hash = ?, secret_cipher = ?,
			secret_version = ?, disabled = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		d.Name, d.Location, d.SecretHash, d.SecretCipher,
		d.SecretVersion, d.Disabled, d.UpdatedAt, d.ID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update device")
	}
	return nil
}

// SetOnlineStatus updates the device online status by publish path
func (r *DeviceRepository) SetOnlineStatus(ctx context.Context, publishPath string, online bool) error {
	now := time.Now()
	query := `UPDATE devices SET is_online = ?, last_seen_at = ?, updated_at = ? WHERE publish_path = ?`
	_, err := r.db.ExecContext(ctx, query, online, now, now, publishPath)
	if err != nil {
		return errors.Wrap(err, "failed to update device online status")
	}
	return nil
}

// Delete deletes a device by ID
func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM devices WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete device")
	}
	return nil
}
