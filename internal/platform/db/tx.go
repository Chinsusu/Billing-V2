package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TxFunc func(ctx context.Context, tx *sql.Tx) error

func WithTx(ctx context.Context, conn *sql.DB, run TxFunc) error {
	if conn == nil {
		return fmt.Errorf("database connection is required")
	}
	if run == nil {
		return fmt.Errorf("transaction function is required")
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := run(ctx, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	committed = true
	return nil
}
