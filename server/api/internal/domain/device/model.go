package device

import "time"

// Device represents a streaming device (e.g., dog cam)
type Device struct {
	ID            string     `json:"id"`
	OwnerUserID   string     `json:"owner_user_id"`
	Name          string     `json:"name"`
	Location      string     `json:"location"`
	PublishPath   string     `json:"publish_path"` // e.g., "dogcam/<device_id>"
	SecretHash    string     `json:"-"`            // bcrypt hash for authentication
	SecretCipher  string     `json:"-"`            // AES encrypted for retrieval
	SecretVersion int        `json:"secret_version"`
	Disabled      bool       `json:"disabled"`
	IsOnline      bool       `json:"is_online"`
	LastSeenAt    *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// CreateDeviceInput represents input for creating a new device
type CreateDeviceInput struct {
	OwnerUserID string
	Name        string
	Location    string
}

// UpdateDeviceInput represents input for updating a device
type UpdateDeviceInput struct {
	Name     *string
	Location *string
	Disabled *bool
}

// DeviceSecret contains the plain-text secret (only returned once on creation)
type DeviceSecret struct {
	DeviceID string `json:"device_id"`
	Secret   string `json:"secret"`
}
