package notifications

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID           uuid.UUID  `json:"id"`
	AppID        uuid.UUID  `json:"app_id"`
	Title        string     `json:"title"`
	Body         string     `json:"body"`
	DeviceID     *uuid.UUID `json:"device_id,omitempty"`
	SegmentID    *uuid.UUID `json:"segment_id,omitempty"`
	Broadcast    bool       `json:"broadcast"`
	ScheduleTime *time.Time `json:"schedule_time,omitempty"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateNotificationRequest struct {
	AppID        uuid.UUID  `json:"app_id"`
	Title        string     `json:"title"`
	Body         string     `json:"body"`
	DeviceID     *uuid.UUID `json:"device_id,omitempty"`
	SegmentID    *uuid.UUID `json:"segment_id,omitempty"`
	Broadcast    bool       `json:"broadcast"`
	ScheduleTime *time.Time `json:"schedule_time,omitempty"`
}

func (r CreateNotificationRequest) Validate() error {
	switch {
	case r.AppID == uuid.Nil:
		return ErrMissingAppID
	case r.Title == "":
		return ErrMissingTitle
	case r.Body == "":
		return ErrMissingBody
	}

	// Exactly one target — mirrors the DB's exactly_one_target CHECK,
	// but checking here lets us return a clean 400 instead of relying
	// on a Postgres constraint violation bubbling up as a generic 500.
	targetCount := 0
	if r.DeviceID != nil {
		targetCount++
	}
	if r.SegmentID != nil {
		targetCount++
	}
	if r.Broadcast {
		targetCount++
	}
	if targetCount != 1 {
		return ErrInvalidTarget
	}

	if r.ScheduleTime != nil && r.ScheduleTime.Before(time.Now()) {
		return ErrScheduleInPast
	}

	return nil
}


func (r CreateNotificationRequest) ToNotification() Notification {
	return Notification{
		AppID:        r.AppID,
		Title:        r.Title,
		Body:         r.Body,
		DeviceID:     r.DeviceID,
		SegmentID:    r.SegmentID,
		Broadcast:    r.Broadcast,
		ScheduleTime: r.ScheduleTime,
		Status:       "pending",
	}
}