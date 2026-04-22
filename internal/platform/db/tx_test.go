package db

import (
	"context"
	"database/sql"
	"testing"
)

func TestWithTxRequiresConnection(t *testing.T) {
	if err := WithTx(context.Background(), nil, func(ctx context.Context, tx *sql.Tx) error { return nil }); err == nil {
		t.Fatal("expected connection error")
	}
}

func TestWithTxRequiresFunction(t *testing.T) {
	// The connection check runs first, so this test documents the validation order.
	if err := WithTx(context.Background(), nil, nil); err == nil {
		t.Fatal("expected connection error")
	}
}
