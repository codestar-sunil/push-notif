package db

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"context"
)


func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
    config, err:=pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
    

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}
	return pool, nil
}

