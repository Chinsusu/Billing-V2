package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

func TestServiceRenewClientServiceDelegatesAndWritesAudits(t *testing.T) {
	store := &fakeOrderStore{
		renewClientServiceResult: ClientServiceRenewal{
			Service:                   ServiceInstance{ID: "service-1", TenantID: "tenant-1", Status: ServiceStatusActive, BillingStatus: BillingStatusPaid},
			InvoiceID:                 "invoice-1",
			InvoiceDisplayID:          10001,
			PaymentTransactionID:      "payment-1",
			PaymentTransactionDisplay: 10002,
			WalletID:                  "wallet-1",
			LedgerEntryID:             "ledger-1",
			LedgerEntryDisplayID:      10003,
			AmountMinor:               1500,
			Currency:                  "USD",
			Renewed:                   true,
			PreviousStatus:            ServiceStatusExpired,
			PreviousTermEnd:           time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	auditLog := &fakeOrderAuditAppender{}
	service := NewServiceWithAudit(store, auditLog)

	result, err := service.RenewClientService(context.Background(), ClientServiceRenewalInput{
		TenantID:       tenant.ID(" tenant-1 "),
		BuyerUserID:    identity.UserID(" buyer-1 "),
		ServiceID:      " service-1 ",
		WalletID:       wallet.WalletID(" wallet-1 "),
		ActorID:        identity.UserID(" buyer-1 "),
		FromStatus:     ServiceStatusExpired,
		IdempotencyKey: " renew-key-1 ",
	})
	if err != nil {
		t.Fatalf("expected client service renewal: %v", err)
	}
	if result.InvoiceID != "invoice-1" {
		t.Fatalf("unexpected renewal result: %+v", result)
	}
	if store.renewClientServiceInput.TenantID != tenant.ID("tenant-1") ||
		store.renewClientServiceInput.BuyerUserID != identity.UserID("buyer-1") ||
		store.renewClientServiceInput.Reason != clientServiceRenewalDefaultReason ||
		store.renewClientServiceInput.IdempotencyKey != IdempotencyKey("renew-key-1") {
		t.Fatalf("expected normalized renewal input, got %+v", store.renewClientServiceInput)
	}
	if auditLog.calls != 2 || auditLog.input.Action != "invoice.wallet_paid" {
		t.Fatalf("expected lifecycle and payment audits, got calls=%d input=%+v", auditLog.calls, auditLog.input)
	}
}

func TestServiceRenewClientServiceDoesNotAuditIdempotentReplay(t *testing.T) {
	store := &fakeOrderStore{
		renewClientServiceResult: ClientServiceRenewal{
			Service:   ServiceInstance{ID: "service-1", TenantID: "tenant-1", Status: ServiceStatusActive, BillingStatus: BillingStatusPaid},
			InvoiceID: "invoice-1",
			WalletID:  "wallet-1",
			Currency:  "USD",
			Renewed:   false,
		},
	}
	auditLog := &fakeOrderAuditAppender{}
	service := NewServiceWithAudit(store, auditLog)

	_, err := service.RenewClientService(context.Background(), ClientServiceRenewalInput{
		TenantID:       "tenant-1",
		BuyerUserID:    "buyer-1",
		ServiceID:      "service-1",
		WalletID:       "wallet-1",
		ActorID:        "buyer-1",
		FromStatus:     ServiceStatusActive,
		IdempotencyKey: "renew-key-1",
	})
	if err != nil {
		t.Fatalf("expected idempotent renewal replay: %v", err)
	}
	if auditLog.calls != 0 {
		t.Fatalf("expected no replay audit, got %d", auditLog.calls)
	}
}

func TestServiceRenewClientServiceRejectsMissingWalletBeforeStore(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	_, err := service.RenewClientService(context.Background(), ClientServiceRenewalInput{
		TenantID:       "tenant-1",
		BuyerUserID:    "buyer-1",
		ServiceID:      "service-1",
		ActorID:        "buyer-1",
		FromStatus:     ServiceStatusActive,
		IdempotencyKey: "renew-key-1",
	})
	if !errors.Is(err, wallet.ErrWalletIDMissing) {
		t.Fatalf("expected wallet id error, got %v", err)
	}
	if store.renewClientServiceInput.ServiceID != "" {
		t.Fatalf("store should not be called, got %+v", store.renewClientServiceInput)
	}
}
