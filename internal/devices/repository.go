package devices

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("device not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Upsert inserts a new device, or updates an existing one if
// (app_id, installation_id) already exists — this is what makes
// re-registration and token rotation safe instead of creating duplicates.
func (r *Repository) Upsert(ctx context.Context, d Device) (Device, error) {
	query := `
		INSERT INTO devices (app_id, installation_id, push_token, platform, user_id, last_seen, status)
		VALUES ($1, $2, $3, $4, $5, now(), 'active')
		ON CONFLICT (app_id, installation_id) DO UPDATE SET
			push_token = EXCLUDED.push_token,
			platform = EXCLUDED.platform,
			user_id = EXCLUDED.user_id,
			last_seen = now(),
			status = 'active'
		RETURNING id, app_id, installation_id, push_token, platform, user_id, last_seen, status, created_at, updated_at
	`
	var out Device
	err := r.pool.QueryRow(ctx, query, d.AppID, d.InstallationID, d.PushToken, d.Platform, d.UserID).Scan(
		&out.ID,
		&out.AppID,
		&out.InstallationID,
		&out.PushToken,
		&out.Platform,
		&out.UserID,
		&out.LastSeen,
		&out.Status,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return Device{}, err
	}
	return out, nil
}

func (r *Repository) Heartbeat(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE devices
		SET last_seen = now()
		WHERE id = $1
	`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (Device, error) {
	query := `
		SELECT id, app_id, installation_id, push_token, platform, user_id, last_seen, status, created_at, updated_at
		FROM devices
		WHERE id = $1
	`
	var out Device
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&out.ID,
		&out.AppID,
		&out.InstallationID,
		&out.PushToken,
		&out.Platform,
		&out.UserID,
		&out.LastSeen,
		&out.Status,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Device{}, ErrDeviceNotFound
	}
	return out, err
}


func (r *Repository) MarkUnregistered(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE devices
		SET status = 'unregistered', last_seen = now()
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id)
	return err
	
}