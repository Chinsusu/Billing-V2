package wallet

import (
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateTopupRequestInputNormalizeValidate(t *testing.T) {
	input := CreateTopupRequestInput{
		TenantID:         tenant.ID("tenant-1"),
		WalletID:         WalletID("wallet-1"),
		RequestedBy:      identity.UserID("account-1"),
		AmountMinor:      1000,
		Currency:         " usd ",
		PaymentMethod:    PaymentMethodBankTransfer,
		PaymentReference: " bank-ref ",
		IdempotencyKey:   IdempotencyKey(" key-1 "),
	}.Normalize()

	if input.Currency != "USD" || input.PaymentReference != "bank-ref" {
		t.Fatalf("expected normalized request, got %+v", input)
	}
	if input.Status != TopupStatusSubmitted {
		t.Fatalf("expected submitted status, got %q", input.Status)
	}
	if input.IdempotencyKey != IdempotencyKey("key-1") {
		t.Fatalf("expected trimmed idempotency key, got %q", input.IdempotencyKey)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid topup input: %v", err)
	}

	input.AmountMinor = 0
	if err := input.Validate(); !errors.Is(err, ErrAmountInvalid) {
		t.Fatalf("expected amount error, got %v", err)
	}
}

func TestTopupEnums(t *testing.T) {
	for _, status := range []TopupStatus{
		TopupStatusDraft, TopupStatusSubmitted, TopupStatusUnderReview, TopupStatusApproved,
		TopupStatusRejected, TopupStatusExpired, TopupStatusCancelled,
	} {
		if !status.Valid() {
			t.Fatalf("expected valid status %q", status)
		}
	}
	if TopupStatus("bad").Valid() {
		t.Fatal("unexpected valid topup status")
	}
	for _, method := range []PaymentMethod{PaymentMethodBankTransfer, PaymentMethodCrypto, PaymentMethodManual, PaymentMethodOther} {
		if !method.Valid() {
			t.Fatalf("expected valid payment method %q", method)
		}
	}
	if PaymentMethod("bad").Valid() {
		t.Fatal("unexpected valid payment method")
	}
}
