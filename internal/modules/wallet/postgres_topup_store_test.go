package wallet

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateTopupRequestArgsNormalizeValidate(t *testing.T) {
	args, err := createTopupRequestArgs(CreateTopupRequestInput{
		TenantID:         tenant.ID("tenant-1"),
		WalletID:         WalletID("wallet-1"),
		RequestedBy:      identity.UserID("account-1"),
		AmountMinor:      1000,
		Currency:         " usd ",
		PaymentMethod:    PaymentMethodBankTransfer,
		PaymentReference: " bank-ref ",
		IdempotencyKey:   IdempotencyKey(" key-1 "),
	})
	if err != nil {
		t.Fatalf("expected args: %v", err)
	}
	reference, ok := args[6].(sql.NullString)
	if !ok || !reference.Valid || reference.String != "bank-ref" {
		t.Fatalf("unexpected payment reference arg: %#v", args[6])
	}
	if args[4] != "USD" || args[7] != TopupStatusSubmitted || args[8] != IdempotencyKey("key-1") {
		t.Fatalf("unexpected normalized args: %#v", args)
	}
}

func TestCreateTopupRequestArgsRejectsMissingIdempotencyKey(t *testing.T) {
	_, err := createTopupRequestArgs(CreateTopupRequestInput{
		TenantID:      tenant.ID("tenant-1"),
		WalletID:      WalletID("wallet-1"),
		RequestedBy:   identity.UserID("account-1"),
		AmountMinor:   1000,
		Currency:      "USD",
		PaymentMethod: PaymentMethodBankTransfer,
	})
	if !errors.Is(err, ErrIdempotencyKeyMissing) {
		t.Fatalf("expected idempotency error, got %v", err)
	}
}
