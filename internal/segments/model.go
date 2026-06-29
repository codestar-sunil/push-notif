package segments

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Segment struct {
	ID        uuid.UUID       `json:"id"`
	AppID     uuid.UUID       `json:"app_id"`
	Name      string          `json:"name"`
	Rule      json.RawMessage `json:"rule"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreateSegmentRequest struct {
	AppID uuid.UUID       `json:"app_id"`
	Name  string          `json:"name"`
	Rule  json.RawMessage `json:"rule"`
}

func (r CreateSegmentRequest) Validate() error {
	if r.AppID == uuid.Nil {
		return ErrMissingAppID
	}
	if r.Name == "" {
		return ErrMissingName
	}
	if len(r.Rule) == 0 {
		return ErrMissingRule
	}
	if _, err := parseRule(r.Rule); err != nil {
		return err
	}
	return nil
}

func (r CreateSegmentRequest) ToSegment() Segment {
	return Segment{
		AppID: r.AppID,
		Name:  r.Name,
		Rule:  r.Rule,
	}
}
