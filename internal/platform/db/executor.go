package db

import (
	"context"
	"database/sql"
)

// Executor is the small database surface store implementations need.
// Both *sql.DB and *sql.Tx satisfy it, so services can choose transaction boundaries.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
