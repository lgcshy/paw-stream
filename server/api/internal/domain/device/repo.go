package device

import "context"

// Repository defines the interface for device data access
type Repository interface {
	// Create creates a new device
	Create(ctx context.Context, device *Device) error

	// GetByID retrieves a device by ID
	GetByID(ctx context.Context, id string) (*Device, error)

	// GetByPublishPath retrieves a device by its publish path
	GetByPublishPath(ctx context.Context, path string) (*Device, error)

	// ListByOwner retrieves all devices owned by a user
	ListByOwner(ctx context.Context, ownerUserID string) ([]*Device, error)

	// Update updates an existing device
	Update(ctx context.Context, device *Device) error

	// Delete deletes a device by ID
	Delete(ctx context.Context, id string) error
}
