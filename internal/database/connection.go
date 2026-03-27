package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

func NewPostgresPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	if dbURL == "" {
		dbURL = "postgres://user:pass@localhost:5432/api_db?sslmode=disable"
	}

	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("Unable to connect to database", "err", err)
		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		slog.Error("Database ping failed", "err", err)
		return nil, err
	}

	return dbPool, nil
}

func NewRedisClient(redisAddr string) *goredis.Client {
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := goredis.NewClient(&goredis.Options{
		Addr: redisAddr,
	})

	return client
}
