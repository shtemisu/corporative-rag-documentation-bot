package postgres

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}
	cfg.MaxConns = 20
	cfg.MinConns = 1
	cfg.MaxConnLifetime = time.Hour * 1
	cfg.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, time.Second)

	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return pool, nil
}

func mapPgErr(err error) error {
	if pgError, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgError.Code {
		case "23505":
			return fmt.Errorf("%w: %s", domain.ErrUserNotFound, pgError.Detail)
		case "23503":
			return fmt.Errorf("related entity missing: %w", err)
		}
	}

	return fmt.Errorf("db: %w", err)
}
