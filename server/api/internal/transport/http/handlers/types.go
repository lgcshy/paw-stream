package handlers

import "time"

// Auth API types

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserInfo `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Device API types

type CreateDeviceRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type CreateDeviceResponse struct {
	Device *DeviceInfo `json:"device"`
	Secret string      `json:"secret"` // Only returned once
}

type UpdateDeviceRequest struct {
	Name     *string `json:"name"`
	Location *string `json:"location"`
	Disabled *bool   `json:"disabled"`
}

type DeviceInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Location    string     `json:"location"`
	PublishPath string     `json:"publish_path"`
	Disabled    bool       `json:"disabled"`
	IsOnline    bool       `json:"is_online"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type RotateSecretResponse struct {
	Secret        string `json:"secret"` // New secret, only returned once
	SecretVersion int    `json:"secret_version"`
}

// Path API types

type PathInfo struct {
	PublishPath    string `json:"publish_path"`
	DeviceID       string `json:"device_id"`
	DeviceName     string `json:"device_name"`
	DeviceLocation string `json:"device_location"`
}

// Admin API types

type AdminDeviceInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Location    string     `json:"location"`
	PublishPath string     `json:"publish_path"`
	OwnerUserID string     `json:"owner_user_id"`
	Disabled    bool       `json:"disabled"`
	IsOnline    bool       `json:"is_online"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type AdminDashboard struct {
	TotalDevices  int               `json:"total_devices"`
	OnlineDevices int               `json:"online_devices"`
	Devices       []AdminDeviceInfo `json:"devices"`
}

// Heartbeat API types

type HeartbeatRequest struct {
	DeviceID string  `json:"device_id"`
	Status   string  `json:"status"`
	FPS      float64 `json:"fps,omitempty"`
	Bitrate  int     `json:"bitrate,omitempty"`
}

// Error response types

type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	RequestID string                 `json:"request_id,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
