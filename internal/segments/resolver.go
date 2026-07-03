package segments

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pageSize = 500

type Resolver struct {
	pool *pgxpool.Pool
}

func NewResolver(pool *pgxpool.Pool) *Resolver {
	return &Resolver{pool: pool}
}

func (r *Resolver) ResolveDevices(ctx context.Context, appID uuid.UUID, rule Rule, onPage func([]uuid.UUID) error) error {
	if rule.isCohort() {
		return ErrCohortNotImplemented
	}

	var cursor uuid.UUID
	for {
		query, args, err := buildLiteralQuery(appID, rule, cursor, pageSize)
		if err != nil {
			return err
		}

		rows, err := r.pool.Query(ctx, query, args...)
		if err != nil {
			return err
		}

		var ids []uuid.UUID
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan device ID: %v", err)
			}
			ids = append(ids, id)
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			return err
		}

		if len(ids) == 0 {
			return nil
		}

		if err := onPage(ids); err != nil {
			return err
		}
		cursor = ids[len(ids)-1]
		if len(ids) < pageSize {
			return nil
		}

	}
}

func buildLiteralQuery(appID uuid.UUID, rule Rule, cursor uuid.UUID, limit int) (string, []any, error) {
	var b strings.Builder
	b.WriteString("SELECT id FROM devices WHERE app_id = $1 AND status = 'active' AND id > $2")
	args := []any{appID, cursor}
	argN := 3

	if rule.Platform != "" {
		fmt.Fprintf(&b, " AND platform = $%d", argN)
		args = append(args, rule.Platform)
		argN++
	}

	if rule.LastSeen != "" {
		dur, err := parseRelativeDuration(rule.LastSeen)
		if err != nil {
			return "", nil, fmt.Errorf("invalid last_seen duration: %v", err)
		}
		fmt.Fprintf(&b, " AND last_seen < $%d", argN)
		args = append(args, time.Now().UTC().Add(-dur))
		argN++
	}
	fmt.Fprintf(&b, " ORDER BY id LIMIT %d", limit)
	return b.String(), args, nil
}

func parseRelativeDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}
