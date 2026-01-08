package idgen

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

// NewUUID generates a new UUID string
func NewUUID() string {
	return uuid.New().String()
}

// NewSecret generates a cryptographically secure random secret
// Returns a base64-encoded string of the specified byte length
func NewSecret(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// NewDeviceSecret generates a device secret (32 bytes, base64 encoded)
func NewDeviceSecret() (string, error) {
	return NewSecret(32)
}
