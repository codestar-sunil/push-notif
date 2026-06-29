package segments

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateSegment(ctx context.Context, segment Segment) (Segment, error) {
	query := `
		INSERT INTO segments (app_id, name, rule)
		VALUES ($1, $2, $3)
		RETURNING id, app_id, name, rule, created_at, updated_at
		`
	var out Segment
	err := r.pool.QueryRow(ctx, query, segment.AppID, segment.Name, segment.Rule).Scan(&out.ID, &out.AppID, &out.Name, &out.Rule, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return Segment{}, err
	}
	return out, nil
}
