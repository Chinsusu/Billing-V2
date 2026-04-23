package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

func TestPayInvoiceFromWalletPaysIssuedInvoice(t *testing.T) {
	invoiceStore := &fakePaymentInvoiceStore{detail: walletPaymentInvoice(invoice.StatusIssued)}
	walletService := &fakePaymentWalletService{record: walletPaymentWallet(2500, "USD")}
	transactionStore := &fakeWalletPaymentTransactionStore{}
	auditLog := &fakePaymentAuditAppender{}
	service := NewServiceWithBillingAndAudit(transactionStore, invoiceStore, walletService, auditLog)

	result, err := service.PayInvoiceFromWallet(context.Background(), walletPaymentInput())
	if err != nil {
		t.Fatalf("expected payment result: %v", err)
	}
	if walletService.lookup.OwnerID != wallet.OwnerID("buyer-1") || walletService.lookup.OwnerType != wallet.OwnerTypeUser {
		t.Fatalf("expected buyer wallet lookup, got %+v", walletService.lookup)
	}
	if walletService.postCalls != 1 || walletService.postInput.Direction != wallet.DirectionDebit {
		t.Fatalf("expected one debit ledger entry, got %+v", walletService.postInput)
	}
	if transactionStore.createCalls != 1 || transactionStore.createInput.InvoiceID != invoice.InvoiceID("invoice-1") {
		t.Fatalf("expected one payment transaction, got %+v", transactionStore.createInput)
	}
	if invoiceStore.markPaidCalls != 1 || invoiceStore.markPaidInput.PaymentTransactionID != "txn-1" {
		t.Fatalf("expected invoice marked paid, got %+v", invoiceStore.markPaidInput)
	}
	if result.Invoice.Invoice.Status != invoice.StatusPaid || result.LedgerEntry.ID != wallet.LedgerEntryID("ledger-1") {
		t.Fatalf("unexpected result: %+v", result)
	}
	if auditLog.calls != 1 ||
		auditLog.input.Action != walletPaymentAuditAction ||
		auditLog.input.TargetID != audit.TargetID("invoice-1") ||
		auditLog.input.ActorID != audit.ActorID("buyer-1") {
		t.Fatalf("unexpected audit input: %+v", auditLog.input)
	}
}

func TestPayInvoiceFromWalletReturnsPaidDuplicateByIdempotency(t *testing.T) {
	invoiceStore := &fakePaymentInvoiceStore{detail: walletPaymentInvoice(invoice.StatusPaid)}
	transactionStore := &fakeWalletPaymentTransactionStore{existing: walletPaymentTransaction()}
	auditLog := &fakePaymentAuditAppender{}
	service := NewServiceWithBillingAndAudit(transactionStore, invoiceStore, &fakePaymentWalletService{}, auditLog)

	result, err := service.PayInvoiceFromWallet(context.Background(), walletPaymentInput())
	if err != nil {
		t.Fatalf("expected duplicate result: %v", err)
	}
	if transactionStore.lookupCalls != 1 || transactionStore.createCalls != 0 {
		t.Fatalf("expected existing transaction lookup only, got lookups=%d creates=%d", transactionStore.lookupCalls, transactionStore.createCalls)
	}
	if invoiceStore.markPaidCalls != 0 {
		t.Fatalf("expected no invoice update for already paid duplicate, got %d", invoiceStore.markPaidCalls)
	}
	if result.Transaction.ID != TransactionID("txn-1") {
		t.Fatalf("expected existing transaction, got %+v", result.Transaction)
	}
	if auditLog.calls != 0 {
		t.Fatalf("expected no duplicate audit event, got %d", auditLog.calls)
	}
}

func TestPayInvoiceFromWalletReturnsAuditError(t *testing.T) {
	invoiceStore := &fakePaymentInvoiceStore{detail: walletPaymentInvoice(invoice.StatusIssued)}
	walletService := &fakePaymentWalletService{record: walletPaymentWallet(2500, "USD")}
	transactionStore := &fakeWalletPaymentTransactionStore{}
	auditLog := &fakePaymentAuditAppender{err: audit.ErrActionMissing}
	service := NewServiceWithBillingAndAudit(transactionStore, invoiceStore, walletService, auditLog)

	_, err := service.PayInvoiceFromWallet(context.Background(), walletPaymentInput())
	if !errors.Is(err, audit.ErrActionMissing) {
		t.Fatalf("expected audit error, got %v", err)
	}
	if invoiceStore.markPaidCalls != 1 || auditLog.calls != 1 {
		t.Fatalf("expected payment mutation and audit attempt, got paid=%d audit=%d", invoiceStore.markPaidCalls, auditLog.calls)
	}
}

func TestPayInvoiceFromWalletRejectsAlreadyPaidWithoutMatchingKey(t *testing.T) {
	invoiceStore := &fakePaymentInvoiceStore{detail: walletPaymentInvoice(invoice.StatusPaid)}
	transactionStore := &fakeWalletPaymentTransactionStore{getErr: ErrTransactionNotFound}
	service := NewServiceWithBilling(transactionStore, invoiceStore, &fakePaymentWalletService{})

	_, err := service.PayInvoiceFromWallet(context.Background(), walletPaymentInput())
	if !errors.Is(err, ErrInvoiceNotPayable) {
		t.Fatalf("expected not payable error, got %v", err)
	}
}

func TestPayInvoiceFromWalletRejectsInsufficientBalance(t *testing.T) {
	invoiceStore := &fakePaymentInvoiceStore{detail: walletPaymentInvoice(invoice.StatusIssued)}
	walletService := &fakePaymentWalletService{record: walletPaymentWallet(1200, "USD")}
	transactionStore := &fakeWalletPaymentTransactionStore{}
	service := NewServiceWithBilling(transactionStore, invoiceStore, walletService)

	_, err := service.PayInvoiceFromWallet(context.Background(), walletPaymentInput())
	if !errors.Is(err, wallet.ErrInsufficientBalance) {
		t.Fatalf("expected insufficient balance, got %v", err)
	}
	if walletService.postCalls != 1 || transactionStore.createCalls != 0 || invoiceStore.markPaidCalls != 0 {
		t.Fatalf("expected no payment transaction or invoice write after balance failure")
	}
}

func TestPayInvoiceFromWalletRejectsCurrencyMismatch(t *testing.T) {
	invoiceStore := &fakePaymentInvoiceStore{detail: walletPaymentInvoice(invoice.StatusIssued)}
	walletService := &fakePaymentWalletService{record: walletPaymentWallet(2500, "VND")}
	service := NewServiceWithBilling(&fakeWalletPaymentTransactionStore{}, invoiceStore, walletService)

	_, err := service.PayInvoiceFromWallet(context.Background(), walletPaymentInput())
	if !errors.Is(err, ErrWalletCurrencyMismatch) {
		t.Fatalf("expected currency mismatch, got %v", err)
	}
	if walletService.postCalls != 0 {
		t.Fatalf("expected no debit after currency mismatch")
	}
}

func TestPayInvoiceFromWalletRejectsCrossTenantInvoice(t *testing.T) {
	detail := walletPaymentInvoice(invoice.StatusIssued)
	detail.Invoice.TenantID = tenant.ID("tenant-2")
	service := NewServiceWithBilling(
		&fakeWalletPaymentTransactionStore{},
		&fakePaymentInvoiceStore{detail: detail},
		&fakePaymentWalletService{record: walletPaymentWallet(2500, "USD")},
	)

	_, err := service.PayInvoiceFromWallet(context.Background(), walletPaymentInput())
	if !errors.Is(err, tenant.ErrAccessDenied) {
		t.Fatalf("expected access denied, got %v", err)
	}
}

func TestPayInvoiceFromWalletRejectsNonPayableStatus(t *testing.T) {
	service := NewServiceWithBilling(
		&fakeWalletPaymentTransactionStore{},
		&fakePaymentInvoiceStore{detail: walletPaymentInvoice(invoice.StatusDraft)},
		&fakePaymentWalletService{record: walletPaymentWallet(2500, "USD")},
	)

	_, err := service.PayInvoiceFromWallet(context.Background(), walletPaymentInput())
	if !errors.Is(err, ErrInvoiceNotPayable) {
		t.Fatalf("expected not payable error, got %v", err)
	}
}

func walletPaymentInput() PayInvoiceFromWalletInput {
	return PayInvoiceFromWalletInput{
		TenantID:       tenant.ID("tenant-1"),
		InvoiceID:      invoice.InvoiceID("invoice-1"),
		WalletID:       wallet.WalletID("wallet-1"),
		ActorID:        identity.UserID("buyer-1"),
		IdempotencyKey: IdempotencyKey("pay-key-1"),
	}
}

func walletPaymentInvoice(status invoice.Status) invoice.InvoiceDetail {
	return invoice.InvoiceDetail{Invoice: invoice.Invoice{
		ID:            invoice.InvoiceID("invoice-1"),
		DisplayID:     10001,
		TenantID:      tenant.ID("tenant-1"),
		BuyerUserID:   identity.UserID("buyer-1"),
		OrderID:       order.OrderID("order-1"),
		Status:        status,
		Currency:      "USD",
		SubtotalMinor: 1800,
		TotalMinor:    1800,
	}}
}

func walletPaymentWallet(balance int64, currency string) wallet.Wallet {
	return wallet.Wallet{
		ID:                    wallet.WalletID("wallet-1"),
		TenantID:              tenant.ID("tenant-1"),
		OwnerType:             wallet.OwnerTypeUser,
		OwnerID:               wallet.OwnerID("buyer-1"),
		Currency:              currency,
		Status:                wallet.StatusActive,
		AvailableBalanceMinor: balance,
	}
}

func walletPaymentTransaction() Transaction {
	return Transaction{
		ID:             TransactionID("txn-1"),
		TenantID:       tenant.ID("tenant-1"),
		AccountUserID:  identity.UserID("buyer-1"),
		OrderID:        order.OrderID("order-1"),
		InvoiceID:      invoice.InvoiceID("invoice-1"),
		Type:           TransactionTypeCharge,
		Status:         TransactionStatusPosted,
		Currency:       "USD",
		AmountMinor:    1800,
		IdempotencyKey: IdempotencyKey("pay-key-1"),
	}
}

type fakeWalletPaymentTransactionStore struct {
	existing    Transaction
	getErr      error
	createInput CreateTransactionInput
	createCalls int
	lookupCalls int
}

func (store *fakeWalletPaymentTransactionStore) CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error) {
	store.createCalls++
	store.createInput = input.Normalize()
	return walletPaymentTransaction(), nil
}

func (store *fakeWalletPaymentTransactionStore) ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	return nil, nil
}

func (store *fakeWalletPaymentTransactionStore) GetTransaction(ctx context.Context, lookup TransactionLookup) (Transaction, error) {
	store.lookupCalls++
	if store.getErr != nil {
		return Transaction{}, store.getErr
	}
	if store.existing.ID != "" {
		return store.existing, nil
	}
	return Transaction{}, ErrTransactionNotFound
}

type fakePaymentInvoiceStore struct {
	detail        invoice.InvoiceDetail
	markPaidInput invoice.MarkInvoicePaidInput
	markPaidCalls int
}

func (store *fakePaymentInvoiceStore) GetInvoice(ctx context.Context, lookup invoice.InvoiceLookup) (invoice.InvoiceDetail, error) {
	return store.detail, nil
}

func (store *fakePaymentInvoiceStore) MarkInvoicePaid(ctx context.Context, input invoice.MarkInvoicePaidInput) (invoice.InvoiceDetail, error) {
	store.markPaidCalls++
	store.markPaidInput = input.Normalize()
	store.detail.Invoice.Status = invoice.StatusPaid
	return store.detail, nil
}

type fakePaymentWalletService struct {
	record    wallet.Wallet
	lookup    wallet.WalletLookup
	postInput wallet.PostLedgerEntryInput
	postCalls int
}

func (service *fakePaymentWalletService) GetWallet(ctx context.Context, lookup wallet.WalletLookup) (wallet.Wallet, error) {
	service.lookup = lookup
	return service.record, nil
}

func (service *fakePaymentWalletService) PostLedgerEntry(ctx context.Context, input wallet.PostLedgerEntryInput) (wallet.LedgerEntry, error) {
	service.postCalls++
	service.postInput = input.Normalize()
	if service.record.AvailableBalanceMinor < input.AmountMinor {
		return wallet.LedgerEntry{}, wallet.ErrInsufficientBalance
	}
	return wallet.LedgerEntry{
		ID:            wallet.LedgerEntryID("ledger-1"),
		WalletID:      input.WalletID,
		TenantID:      input.TenantID,
		Direction:     wallet.DirectionDebit,
		AmountMinor:   input.AmountMinor,
		Currency:      input.Currency,
		EntryType:     wallet.EntryTypePurchase,
		Status:        wallet.LedgerStatusPosted,
		ReferenceType: input.ReferenceType,
		ReferenceID:   input.ReferenceID,
	}, nil
}

type fakePaymentAuditAppender struct {
	input audit.AppendInput
	err   error
	calls int
}

func (appender *fakePaymentAuditAppender) Append(ctx context.Context, input audit.AppendInput) (audit.Log, error) {
	appender.calls++
	appender.input = input
	if appender.err != nil {
		return audit.Log{}, appender.err
	}
	return audit.Log{Action: input.Action}, nil
}
