package notifications

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const createQuery = `
INSERT INTO notifications (app_id, title, body, device_id, segment_id, broadcast, schedule_time)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, app_id, title, body, device_id, segment_id, broadcast, schedule_time, status, created_at, updated_at`

func (r *Repository) CreateNotification(ctx context.Context, req Notification) (Notification, error) {
	var out Notification
	err := r.pool.QueryRow(ctx, createQuery,
		req.AppID, req.Title, req.Body, req.DeviceID, req.SegmentID, req.Broadcast, req.ScheduleTime,
	).Scan(
		&out.ID, &out.AppID, &out.Title, &out.Body, &out.DeviceID, &out.SegmentID,
		&out.Broadcast, &out.ScheduleTime, &out.Status, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return Notification{}, err
	}
	return out, nil
}

const GetNotificationByIDQuery = `
SELECT id, app_id, title, body, device_id, segment_id, broadcast, schedule_time, status, created_at, updated_at
FROM notifications
WHERE id = $1
`

func (r *Repository) GetNotificationByID(ctx context.Context, id uuid.UUID) (Notification, error) {
	var out Notification
	err := r.pool.QueryRow(ctx, GetNotificationByIDQuery, id).Scan(
		&out.ID, &out.AppID, &out.Title, &out.Body, &out.DeviceID, &out.SegmentID,
		&out.Broadcast, &out.ScheduleTime, &out.Status, &out.CreatedAt, &out.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Notification{}, ErrNotFound
	}
	return out, nil
}
