package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnect(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	const op string = "database.postgres.NewConnect"

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s - unable to open the database: %w", op, err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s - failed to connect to postgres server: %w", op, err)
	}

	return db, nil
}
