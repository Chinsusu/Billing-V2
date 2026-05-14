package order

import (
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

func TestClientServiceRenewalInputNormalizesAndValidates(t *testing.T) {
	input := ClientServiceRenewalInput{
		TenantID:       tenant.ID(" tenant-1 "),
		BuyerUserID:    identity.UserID(" buyer-1 "),
		ServiceID:      " service-1 ",
		WalletID:       wallet.WalletID(" wallet-1 "),
		ActorID:        identity.UserID(" buyer-1 "),
		FromStatus:     ServiceStatusActive,
		IdempotencyKey: " renew-key-1 ",
	}.Normalize()

	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid input: %v", err)
	}
	if input.TenantID != tenant.ID("tenant-1") ||
		input.BuyerUserID != identity.UserID("buyer-1") ||
		input.WalletID != wallet.WalletID("wallet-1") ||
		input.Reason != clientServiceRenewalDefaultReason ||
		input.IdempotencyKey != IdempotencyKey("renew-key-1") {
		t.Fatalf("unexpected normalized input: %+v", input)
	}
}

func TestRenewalLedgerIdempotencyScopesServiceAndRequestKey(t *testing.T) {
	key := renewalLedgerIdempotency(ClientServiceRenewalInput{
		ServiceID:      "service-1",
		IdempotencyKey: "renew-key-1",
	})

	if key != wallet.IdempotencyKey("service-renewal:service-1:renew-key-1") {
		t.Fatalf("unexpected ledger idempotency key: %s", key)
	}
}

func TestClientServiceRenewalSQLLocksServiceAndCreatesStandaloneInvoice(t *testing.T) {
	required := []string{
		"FOR UPDATE OF svc",
		"INSERT INTO invoices (tenant_id, buyer_user_id, status, currency",
		"INSERT INTO invoice_items (invoice_id, tenant_id, service_instance_id",
		"INSERT INTO payment_transactions (tenant_id, account_user_id, invoice_id",
		"ON CONFLICT (tenant_id, idempotency_key) DO NOTHING",
	}
	combined := strings.Join([]string{
		clientServiceRenewalContextSQL,
		createRenewalInvoiceSQL,
		createRenewalInvoiceItemSQL,
		createRenewalPaymentTransactionSQL,
	}, "\n")
	for _, clause := range required {
		if !strings.Contains(combined, clause) {
			t.Fatalf("expected %q in renewal SQL: %s", clause, combined)
		}
	}
	if strings.Contains(createRenewalPaymentTransactionSQL, "order_id") {
		t.Fatalf("renewal payment should not attach an order_id or trigger order finalization: %s", createRenewalPaymentTransactionSQL)
	}
}
