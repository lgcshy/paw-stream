package device

import (
	"context"
	"fmt"
	"time"

	"github.com/lgc/pawstream/api/internal/pkg/crypto"
	"github.com/lgc/pawstream/api/internal/pkg/errors"
	"github.com/lgc/pawstream/api/internal/pkg/idgen"
	"github.com/lgc/pawstream/api/internal/pkg/password"
)

// Service handles device business logic
type Service struct {
	repo          Repository
	encryptionKey string
}

// NewService creates a new device service
func NewService(repo Repository, encryptionKey string) *Service {
	return &Service{
		repo:          repo,
		encryptionKey: encryptionKey,
	}
}

// Create creates a new device with a unique secret
func (s *Service) Create(ctx context.Context, input CreateDeviceInput) (*Device, *DeviceSecret, error) {
	// Generate device ID
	deviceID := idgen.NewUUID()

	// Generate device secret
	secret, err := idgen.NewDeviceSecret()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate device secret")
	}

	// Hash secret for authentication
	secretHash, err := password.Hash(secret)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to hash device secret")
	}

	// Encrypt secret for storage
	secretCipher := secret
	if s.encryptionKey != "" {
		encrypted, err := crypto.Encrypt(secret, s.encryptionKey)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to encrypt device secret")
		}
		secretCipher = encrypted
	}

	// Create device
	device := &Device{
		ID:            deviceID,
		OwnerUserID:   input.OwnerUserID,
		Name:          input.Name,
		Location:      input.Location,
		PublishPath:   fmt.Sprintf("dogcam/%s", deviceID),
		SecretHash:    secretHash,
		SecretCipher:  secretCipher,
		SecretVersion: 1,
		Disabled:      false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, nil, errors.Wrap(err, "failed to create device")
	}

	// Return device and plain-text secret (only time it's available)
	deviceSecret := &DeviceSecret{
		DeviceID: deviceID,
		Secret:   secret,
	}

	return device, deviceSecret, nil
}

// GetByID retrieves a device by ID
func (s *Service) GetByID(ctx context.Context, id string) (*Device, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get device")
	}

	if device == nil {
		return nil, errors.ErrDeviceNotFound
	}

	return device, nil
}

// GetByPublishPath retrieves a device by its publish path
func (s *Service) GetByPublishPath(ctx context.Context, path string) (*Device, error) {
	device, err := s.repo.GetByPublishPath(ctx, path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get device by path")
	}

	if device == nil {
		return nil, errors.ErrDeviceNotFound
	}

	return device, nil
}

// ListByOwner retrieves all devices owned by a user
func (s *Service) ListByOwner(ctx context.Context, ownerUserID string) ([]*Device, error) {
	devices, err := s.repo.ListByOwner(ctx, ownerUserID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list devices")
	}

	return devices, nil
}

// Update updates a device's information
func (s *Service) Update(ctx context.Context, id string, input UpdateDeviceInput) (*Device, error) {
	device, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.Name != nil {
		device.Name = *input.Name
	}

	if input.Location != nil {
		device.Location = *input.Location
	}

	if input.Disabled != nil {
		device.Disabled = *input.Disabled
	}

	device.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, errors.Wrap(err, "failed to update device")
	}

	return device, nil
}

// VerifySecret verifies a device secret for authentication
func (s *Service) VerifySecret(ctx context.Context, deviceID, secret string) (bool, error) {
	device, err := s.GetByID(ctx, deviceID)
	if err != nil {
		return false, err
	}

	if device.Disabled {
		return false, errors.ErrDeviceDisabled
	}

	// Verify secret
	if !password.Verify(secret, device.SecretHash) {
		return false, nil
	}

	return true, nil
}

// RotateSecret generates a new secret for a device
func (s *Service) RotateSecret(ctx context.Context, deviceID string) (*DeviceSecret, error) {
	device, err := s.GetByID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// Generate new secret
	secret, err := idgen.NewDeviceSecret()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new secret")
	}

	// Hash new secret
	secretHash, err := password.Hash(secret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash new secret")
	}

	// Encrypt new secret
	secretCipher := secret
	if s.encryptionKey != "" {
		encrypted, err := crypto.Encrypt(secret, s.encryptionKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encrypt device secret")
		}
		secretCipher = encrypted
	}

	// Update device
	device.SecretHash = secretHash
	device.SecretCipher = secretCipher
	device.SecretVersion++
	device.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, errors.Wrap(err, "failed to update device")
	}

	deviceSecret := &DeviceSecret{
		DeviceID: deviceID,
		Secret:   secret,
	}

	return deviceSecret, nil
}

// Delete deletes a device by ID
func (s *Service) Delete(ctx context.Context, deviceID string) error {
	// Verify device exists
	device, err := s.GetByID(ctx, deviceID)
	if err != nil {
		return err
	}

	// Delete device
	if err := s.repo.Delete(ctx, device.ID); err != nil {
		return errors.Wrap(err, "failed to delete device")
	}

	return nil
}

// SetOnlineStatus updates a device's online/offline status by publish path
func (s *Service) SetOnlineStatus(ctx context.Context, publishPath string, online bool) error {
	return s.repo.SetOnlineStatus(ctx, publishPath, online)
}

// ListSharedWith retrieves all devices shared with a user
func (s *Service) ListSharedWith(ctx context.Context, userID string) ([]*Device, error) {
	return s.repo.ListSharedWith(ctx, userID)
}

// ListAll retrieves all devices (admin use)
func (s *Service) ListAll(ctx context.Context) ([]*Device, error) {
	return s.repo.ListAll(ctx)
}
