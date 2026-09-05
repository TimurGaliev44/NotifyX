package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*Database, error) {
	dbPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping Postgres: %w", err)
	}

	return &Database{Pool: dbPool}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
