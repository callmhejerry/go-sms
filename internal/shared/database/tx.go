package database

import (
	"context"
	"fmt"

	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func WithTx(ctx context.Context, pool *pgxpool.Pool, fn func(q *store.Queries) error) error {
	tx, err := pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("Begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := store.New(tx)

	if err := fn(q); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		fmt.Errorf("Commit transaction: %w", err)
	}
	return nil
}
