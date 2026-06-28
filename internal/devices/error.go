package devices

import "errors"

var (
	ErrMissingAppID          = errors.New("app_id is required")
	ErrMissingInstallationID = errors.New("installation_id is required")
	ErrMissingPushToken      = errors.New("push_token is required")
	ErrInvalidPlatform       = errors.New("platform must be ios or android")
	ErrDeviceNotFound        = errors.New("device not found")
)