package segments

import (
	"encoding/json"
	"errors"
)

var (
	ErrMissingAppID         = errors.New("app_id is required")
	ErrMissingName          = errors.New("name is required")
	ErrMissingRule          = errors.New("rule is required")
	ErrInvalidRule          = errors.New("rule must specify exactly one of: platform/last_seen_before, or cohort_id")
	ErrCohortNotImplemented = errors.New("cohort-based targeting is not implemented yet")
)

type Rule struct {
	Platform string
	LastSeen string
	CohortID string
}

func parseRule(rule json.RawMessage) (Rule, error) {
	var r Rule
	if err := json.Unmarshal(rule, &r); err != nil {
		return Rule{}, err
	}
	isCohort := r.CohortID != ""
	isLiteral := r.Platform != "" || r.LastSeen != ""

	if isCohort == isLiteral { // both true or both false
		return Rule{}, ErrInvalidRule
	}
	return r, nil
}

func (r Rule) isCohort() bool {
	return r.CohortID != ""
}
