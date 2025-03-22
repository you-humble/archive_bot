package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewConnect(ctx context.Context, opts *redis.Options) (*redis.Client, error) {
	const op string = "database.redis.NewConnect"

	db := redis.NewClient(opts)
	if err := db.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%s - failed to connect to redis server: %w", op, err)
	}

	return db, nil
}


