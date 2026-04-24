package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

const walletPaymentReferenceInvoice wallet.ReferenceType = "invoice"

type InvoicePaymentStore interface {
	GetInvoice(ctx context.Context, lookup invoice.InvoiceLookup) (invoice.InvoiceDetail, error)
	MarkInvoicePaid(ctx context.Context, input invoice.MarkInvoicePaidInput) (invoice.InvoiceDetail, error)
}

type WalletPaymentService interface {
	GetWallet(ctx context.Context, lookup wallet.WalletLookup) (wallet.Wallet, error)
	PostLedgerEntry(ctx context.Context, input wallet.PostLedgerEntryInput) (wallet.LedgerEntry, error)
}

type OrderPaymentFinalizer interface {
	FinalizePayment(ctx context.Context, input order.FinalizePaymentInput) (order.Order, error)
}

type WalletInvoicePaymentStore interface {
	PayInvoiceFromWallet(ctx context.Context, input PayInvoiceFromWalletInput) (WalletInvoicePayment, error)
}

type PayInvoiceFromWalletInput struct {
	TenantID       tenant.ID
	InvoiceID      invoice.InvoiceID
	WalletID       wallet.WalletID
	ActorID        identity.UserID
	IdempotencyKey IdempotencyKey
}

type WalletInvoicePayment struct {
	Invoice               invoice.InvoiceDetail
	Transaction           Transaction
	LedgerEntry           wallet.LedgerEntry
	Order                 order.Order
	PreviousInvoiceStatus invoice.Status
}

func (input PayInvoiceFromWalletInput) Normalize() PayInvoiceFromWalletInput {
	output := input
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	return output
}

func (input PayInvoiceFromWalletInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.InvoiceID.Empty() {
		return invoice.ErrInvoiceIDMissing
	}
	if input.WalletID.Empty() {
		return wallet.ErrWalletIDMissing
	}
	if input.ActorID == "" {
		return identity.ErrActorIDMissing
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	return nil
}

func payInvoiceFromWallet(
	ctx context.Context,
	transactionStore Store,
	invoiceStore InvoicePaymentStore,
	walletService WalletPaymentService,
	orderFinalizer OrderPaymentFinalizer,
	input PayInvoiceFromWalletInput,
) (WalletInvoicePayment, error) {
	detail, err := invoiceStore.GetInvoice(ctx, invoice.InvoiceLookup{ID: input.InvoiceID, TenantID: input.TenantID})
	if err != nil {
		return WalletInvoicePayment{}, err
	}
	if detail.Invoice.TenantID != input.TenantID {
		return WalletInvoicePayment{}, tenant.ErrAccessDenied
	}
	if detail.Invoice.Status == invoice.StatusPaid {
		transaction, err := transactionStore.GetTransaction(ctx, TransactionLookup{
			TenantID:       input.TenantID,
			AccountUserID:  detail.Invoice.BuyerUserID,
			IdempotencyKey: input.IdempotencyKey,
		})
		if err != nil {
			return WalletInvoicePayment{}, ErrInvoiceNotPayable
		}
		if err := ensureTransactionMatchesInvoice(transaction, detail.Invoice); err != nil {
			return WalletInvoicePayment{}, err
		}
		paidOrder, err := finalizeWalletPaymentOrder(ctx, orderFinalizer, detail.Invoice)
		if err != nil {
			return WalletInvoicePayment{}, err
		}
		return WalletInvoicePayment{Invoice: detail, Transaction: transaction, Order: paidOrder}, nil
	}
	if !invoicePayableStatus(detail.Invoice.Status) {
		return WalletInvoicePayment{}, ErrInvoiceNotPayable
	}

	accountWallet, err := walletService.GetWallet(ctx, wallet.WalletLookup{
		ID:        input.WalletID,
		TenantID:  input.TenantID,
		OwnerType: wallet.OwnerTypeUser,
		OwnerID:   wallet.UserOwnerID(detail.Invoice.BuyerUserID),
	})
	if err != nil {
		return WalletInvoicePayment{}, err
	}
	if err := validateWalletForInvoice(accountWallet, detail.Invoice); err != nil {
		return WalletInvoicePayment{}, err
	}

	transaction, err := transactionStore.GetTransaction(ctx, TransactionLookup{
		TenantID:       input.TenantID,
		AccountUserID:  detail.Invoice.BuyerUserID,
		IdempotencyKey: input.IdempotencyKey,
	})
	if err != nil && !errors.Is(err, ErrTransactionNotFound) {
		return WalletInvoicePayment{}, err
	}
	if err == nil {
		if err := ensureTransactionMatchesInvoice(transaction, detail.Invoice); err != nil {
			return WalletInvoicePayment{}, err
		}
	}

	ledgerEntry, err := walletService.PostLedgerEntry(ctx, walletLedgerPaymentInput(input, detail.Invoice))
	if err != nil {
		return WalletInvoicePayment{}, err
	}
	if err := ensureLedgerMatchesInvoice(ledgerEntry, detail.Invoice); err != nil {
		return WalletInvoicePayment{}, err
	}

	if transaction.ID.Empty() {
		transaction, err = transactionStore.CreateTransaction(ctx, createWalletPaymentTransactionInput(input, detail.Invoice, ledgerEntry.ID))
		if err != nil {
			return WalletInvoicePayment{}, err
		}
		if err := ensureTransactionMatchesInvoice(transaction, detail.Invoice); err != nil {
			return WalletInvoicePayment{}, err
		}
	}

	paidInvoice, err := invoiceStore.MarkInvoicePaid(ctx, invoice.MarkInvoicePaidInput{
		ID:                   detail.Invoice.ID,
		TenantID:             detail.Invoice.TenantID,
		PaymentTransactionID: string(transaction.ID),
		WalletID:             string(accountWallet.ID),
		LedgerEntryID:        string(ledgerEntry.ID),
		IdempotencyKey:       invoice.IdempotencyKey(input.IdempotencyKey),
	})
	if err != nil {
		return WalletInvoicePayment{}, err
	}
	paidOrder, err := finalizeWalletPaymentOrder(ctx, orderFinalizer, paidInvoice.Invoice)
	if err != nil {
		return WalletInvoicePayment{}, err
	}
	return WalletInvoicePayment{
		Invoice:               paidInvoice,
		Transaction:           transaction,
		LedgerEntry:           ledgerEntry,
		Order:                 paidOrder,
		PreviousInvoiceStatus: detail.Invoice.Status,
	}, nil
}

func finalizeWalletPaymentOrder(
	ctx context.Context,
	finalizer OrderPaymentFinalizer,
	record invoice.Invoice,
) (order.Order, error) {
	if finalizer == nil || record.OrderID.Empty() {
		return order.Order{}, nil
	}
	return finalizer.FinalizePayment(ctx, order.FinalizePaymentInput{
		ID:          record.OrderID,
		TenantID:    record.TenantID,
		BuyerUserID: record.BuyerUserID,
	})
}

func invoicePayableStatus(status invoice.Status) bool {
	return status == invoice.StatusIssued || status == invoice.StatusOverdue
}

func validateWalletForInvoice(accountWallet wallet.Wallet, record invoice.Invoice) error {
	if accountWallet.TenantID != record.TenantID ||
		accountWallet.OwnerType != wallet.OwnerTypeUser ||
		accountWallet.OwnerID != wallet.UserOwnerID(record.BuyerUserID) {
		return tenant.ErrAccessDenied
	}
	if accountWallet.Currency != record.Currency {
		return ErrWalletCurrencyMismatch
	}
	return nil
}

func walletLedgerPaymentInput(input PayInvoiceFromWalletInput, record invoice.Invoice) wallet.PostLedgerEntryInput {
	return wallet.PostLedgerEntryInput{
		WalletID:       input.WalletID,
		TenantID:       input.TenantID,
		Direction:      wallet.DirectionDebit,
		AmountMinor:    record.TotalMinor,
		Currency:       record.Currency,
		EntryType:      wallet.EntryTypePurchase,
		ReferenceType:  walletPaymentReferenceInvoice,
		ReferenceID:    wallet.ReferenceID(record.ID),
		IdempotencyKey: walletPaymentLedgerIdempotency(input),
		CreatedBy:      input.ActorID,
		Reason:         fmt.Sprintf("Invoice %d wallet payment", record.DisplayID),
		CorrelationID:  wallet.CorrelationID(record.ID),
	}
}

func createWalletPaymentTransactionInput(
	input PayInvoiceFromWalletInput,
	record invoice.Invoice,
	ledgerEntryID wallet.LedgerEntryID,
) CreateTransactionInput {
	return CreateTransactionInput{
		TenantID:       input.TenantID,
		AccountUserID:  record.BuyerUserID,
		OrderID:        record.OrderID,
		InvoiceID:      record.ID,
		Type:           TransactionTypeCharge,
		Status:         TransactionStatusPosted,
		Currency:       record.Currency,
		AmountMinor:    record.TotalMinor,
		Description:    fmt.Sprintf("Invoice %d wallet payment", record.DisplayID),
		IdempotencyKey: input.IdempotencyKey,
		Metadata:       walletPaymentMetadata(input.WalletID, ledgerEntryID),
	}
}

func walletPaymentLedgerIdempotency(input PayInvoiceFromWalletInput) wallet.IdempotencyKey {
	return wallet.IdempotencyKey(fmt.Sprintf("invoice-payment:%s:%s", input.InvoiceID, input.IdempotencyKey))
}

func walletPaymentMetadata(walletID wallet.WalletID, ledgerEntryID wallet.LedgerEntryID) json.RawMessage {
	payload := struct {
		Source        string               `json:"source"`
		WalletID      wallet.WalletID      `json:"wallet_id"`
		LedgerEntryID wallet.LedgerEntryID `json:"ledger_entry_id"`
	}{
		Source:        "wallet",
		WalletID:      walletID,
		LedgerEntryID: ledgerEntryID,
	}
	value, err := json.Marshal(payload)
	if err != nil {
		return json.RawMessage(`{"source":"wallet"}`)
	}
	return value
}

func ensureTransactionMatchesInvoice(transaction Transaction, record invoice.Invoice) error {
	if transaction.TenantID != record.TenantID ||
		transaction.AccountUserID != record.BuyerUserID ||
		transaction.OrderID != record.OrderID ||
		transaction.InvoiceID != record.ID ||
		transaction.Type != TransactionTypeCharge ||
		transaction.Status != TransactionStatusPosted ||
		transaction.Currency != record.Currency ||
		transaction.AmountMinor != record.TotalMinor {
		return ErrIdempotencyConflict
	}
	return nil
}

func ensureLedgerMatchesInvoice(entry wallet.LedgerEntry, record invoice.Invoice) error {
	if entry.TenantID != record.TenantID ||
		entry.Direction != wallet.DirectionDebit ||
		entry.EntryType != wallet.EntryTypePurchase ||
		entry.Status != wallet.LedgerStatusPosted ||
		entry.Currency != record.Currency ||
		entry.AmountMinor != record.TotalMinor ||
		entry.ReferenceType != walletPaymentReferenceInvoice ||
		entry.ReferenceID != wallet.ReferenceID(record.ID) {
		return ErrIdempotencyConflict
	}
	return nil
}
