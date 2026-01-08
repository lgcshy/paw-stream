package errors

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateUsername  = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled       = errors.New("user is disabled")

	// Device errors
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDuplicateDevicePath = errors.New("device path already exists")
	ErrDeviceDisabled      = errors.New("device is disabled")
	ErrInvalidSecret       = errors.New("invalid device secret")

	// Auth errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
	ErrMissingToken = errors.New("missing token")

	// Database errors
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")
	ErrDatabaseMigration  = errors.New("database migration error")

	// General errors
	ErrInvalidInput   = errors.New("invalid input")
	ErrInternalError  = errors.New("internal server error")
	ErrNotImplemented = errors.New("not implemented")
)

// AppError wraps an error with additional context
type AppError struct {
	Err     error
	Message string
	Code    string
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Wrap wraps an error with a message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &AppError{
		Err:     err,
		Message: message,
	}
}

// Is checks if an error matches a target error
func Is(err, target error) bool {
	return errors.Is(err, target)
}
