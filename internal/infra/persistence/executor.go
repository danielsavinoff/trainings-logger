package persistence

import (
	"context"
	"database/sql"
)

type Executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// ExecutorFromContext returns the active transaction, or db when none is present.
func ExecutorFromContext(ctx context.Context, db *sql.DB) Executor {
	if tx, ok := FromContext(ctx); ok {
		return tx
	}

	return db
}
