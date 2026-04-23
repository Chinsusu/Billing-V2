package payment

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateTransactionInputNormalizeValidate(t *testing.T) {
	input := CreateTransactionInput{
		TenantID:       tenant.ID("tenant-1"),
		AccountUserID:  identity.UserID("account-1"),
		Type:           TransactionTypeCharge,
		Currency:       " usd ",
		AmountMinor:    1500,
		Description:    "  VPS invoice payment  ",
		IdempotencyKey: " txn-key-1 ",
		Metadata:       json.RawMessage(`{"channel":"wallet"}`),
	}.Normalize()

	if input.Currency != "USD" {
		t.Fatalf("expected normalized currency, got %q", input.Currency)
	}
	if input.Description != "VPS invoice payment" {
		t.Fatalf("expected trimmed description, got %q", input.Description)
	}
	if input.Status != TransactionStatusPosted {
		t.Fatalf("expected default posted status, got %q", input.Status)
	}
	if input.IdempotencyKey != IdempotencyKey("txn-key-1") {
		t.Fatalf("expected trimmed idempotency key, got %q", input.IdempotencyKey)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid transaction input: %v", err)
	}

	input.AmountMinor = 0
	if err := input.Validate(); !errors.Is(err, ErrAmountInvalid) {
		t.Fatalf("expected amount error, got %v", err)
	}
}

func TestCreateTransactionInputRequiresType(t *testing.T) {
	err := CreateTransactionInput{
		TenantID:       tenant.ID("tenant-1"),
		AccountUserID:  identity.UserID("account-1"),
		Currency:       "USD",
		AmountMinor:    100,
		IdempotencyKey: "txn-key-1",
	}.Normalize().Validate()
	if !errors.Is(err, ErrTypeInvalid) {
		t.Fatalf("expected type error, got %v", err)
	}
}

func TestTransactionEnums(t *testing.T) {
	if !TransactionTypeCharge.Valid() || !TransactionTypeRefund.Valid() || !TransactionTypeAdjustment.Valid() {
		t.Fatal("expected core transaction types to be valid")
	}
	if TransactionType("bad").Valid() {
		t.Fatal("unexpected valid transaction type")
	}
	if !TransactionStatusPending.Valid() || !TransactionStatusPosted.Valid() {
		t.Fatal("expected core transaction statuses to be valid")
	}
	if TransactionStatus("bad").Valid() {
		t.Fatal("unexpected valid transaction status")
	}
}
