package notifications

import "errors"

var (
	ErrMissingAppID  = errors.New("app_id is required")
	ErrMissingTitle  = errors.New("title is required")
	ErrMissingBody   = errors.New("body is required")
	ErrInvalidTarget = errors.New("exactly one of device_id, segment_id, or broadcast must be set")
	ErrScheduleInPast = errors.New("schedule_time must be in the future")
	ErrNotFound      = errors.New("notification not found")
)