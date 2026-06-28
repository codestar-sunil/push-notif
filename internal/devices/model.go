package devices

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID             uuid.UUID
	AppID          uuid.UUID
	InstallationID string
	PushToken      string
	Platform       string
	UserID         *string
	LastSeen       time.Time
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// RegisterRequest is the wire payload for POST /devices/register.
// AppID comes from the API key/auth middleware in a real deployment;
// for now it's accepted on the body since auth isn't built yet.
type RegisterRequest struct {
	AppID          uuid.UUID `json:"app_id"`
	InstallationID string    `json:"installation_id"`
	PushToken      string    `json:"push_token"`
	Platform       string    `json:"platform"`
	UserID         *string   `json:"user_id,omitempty"`
}

func (r RegisterRequest) Validate() error {
	switch {
	case r.AppID == uuid.Nil:
		return ErrMissingAppID
	case r.InstallationID == "":
		return ErrMissingInstallationID
	case r.PushToken == "":
		return ErrMissingPushToken
	case r.Platform != "ios" && r.Platform != "android":
		return ErrInvalidPlatform
	}
	return nil
}

// ToDevice builds a Device for insertion. ID/timestamps are left zero —
// the DB generates ID via gen_random_uuid() and timestamps via DEFAULT now(),
// so this is really just a field-mapping step, not a full row construction.
func (r RegisterRequest) ToDevice() Device {
	return Device{
		AppID:          r.AppID,
		InstallationID: r.InstallationID,
		PushToken:      r.PushToken,
		Platform:       r.Platform,
		UserID:         r.UserID,
		Status:         "active",
	}
}
