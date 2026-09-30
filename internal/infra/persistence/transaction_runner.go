package persistence

import (
	"context"
	"database/sql"
)

type transactionKey struct{}

type TransactionRunner[T any] struct {
	db *sql.DB
}

func NewTransactionRunner[T any](db *sql.DB) *TransactionRunner[T] {
	return &TransactionRunner[T]{db: db}
}

func (r *TransactionRunner[T]) Run(
	ctx context.Context,
	fn func(context.Context) (T, error),
) (T, error) {
	var zero T

	if _, ok := FromContext(ctx); ok {
		result, err := fn(ctx)
		if err != nil {
			return zero, err
		}

		return result, nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return zero, err
	}

	// Roll back on panics and commit failures; no-op after a commit.
	defer tx.Rollback()

	ctx = WithContext(ctx, tx)

	result, err := fn(ctx)
	if err != nil {
		return zero, err
	}

	if err := tx.Commit(); err != nil {
		return zero, err
	}

	return result, nil
}

func WithContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, transactionKey{}, tx)
}

func FromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(transactionKey{}).(*sql.Tx)

	return tx, ok
}
