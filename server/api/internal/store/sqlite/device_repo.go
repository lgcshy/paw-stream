package sqlite

import (
	"context"
	"database/sql"

	"github.com/lgc/pawstream/api/internal/domain/device"
	"github.com/lgc/pawstream/api/internal/pkg/errors"
)

// DeviceRepository implements device.Repository for SQLite
type DeviceRepository struct {
	db *DB
}

// NewDeviceRepository creates a new SQLite device repository
func NewDeviceRepository(db *DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// Create creates a new device
func (r *DeviceRepository) Create(ctx context.Context, d *device.Device) error {
	query := `
		INSERT INTO devices (
			id, owner_user_id, name, location, publish_path,
			secret_hash, secret_cipher, secret_version,
			disabled, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		d.ID, d.OwnerUserID, d.Name, d.Location, d.PublishPath,
		d.SecretHash, d.SecretCipher, d.SecretVersion,
		d.Disabled, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create device")
	}
	return nil
}

// GetByID retrieves a device by ID
func (r *DeviceRepository) GetByID(ctx context.Context, id string) (*device.Device, error) {
	query := `
		SELECT id, owner_user_id, name, location, publish_path,
			   secret_hash, secret_cipher, secret_version,
			   disabled, created_at, updated_at
		FROM devices
		WHERE id = ?
	`
	var d device.Device
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.OwnerUserID, &d.Name, &d.Location, &d.PublishPath,
		&d.SecretHash, &d.SecretCipher, &d.SecretVersion,
		&d.Disabled, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get device by ID")
	}
	return &d, nil
}

// GetByPublishPath retrieves a device by its publish path
func (r *DeviceRepository) GetByPublishPath(ctx context.Context, path string) (*device.Device, error) {
	query := `
		SELECT id, owner_user_id, name, location, publish_path,
			   secret_hash, secret_cipher, secret_version,
			   disabled, created_at, updated_at
		FROM devices
		WHERE publish_path = ?
	`
	var d device.Device
	err := r.db.QueryRowContext(ctx, query, path).Scan(
		&d.ID, &d.OwnerUserID, &d.Name, &d.Location, &d.PublishPath,
		&d.SecretHash, &d.SecretCipher, &d.SecretVersion,
		&d.Disabled, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get device by path")
	}
	return &d, nil
}

// ListByOwner retrieves all devices owned by a user
func (r *DeviceRepository) ListByOwner(ctx context.Context, ownerUserID string) ([]*device.Device, error) {
	query := `
		SELECT id, owner_user_id, name, location, publish_path,
			   secret_hash, secret_cipher, secret_version,
			   disabled, created_at, updated_at
		FROM devices
		WHERE owner_user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, ownerUserID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list devices")
	}
	defer rows.Close()

	var devices []*device.Device
	for rows.Next() {
		var d device.Device
		if err := rows.Scan(
			&d.ID, &d.OwnerUserID, &d.Name, &d.Location, &d.PublishPath,
			&d.SecretHash, &d.SecretCipher, &d.SecretVersion,
			&d.Disabled, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan device")
		}
		devices = append(devices, &d)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating devices")
	}

	return devices, nil
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

// Delete deletes a device by ID
func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM devices WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete device")
	}
	return nil
}
