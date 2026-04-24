package payment

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

type transactionResponse struct {
	ID            TransactionID     `json:"id"`
	DisplayID     int64             `json:"display_id"`
	TenantID      tenant.ID         `json:"tenant_id"`
	AccountUserID identity.UserID   `json:"account_user_id"`
	OrderID       order.OrderID     `json:"order_id,omitempty"`
	InvoiceID     invoice.InvoiceID `json:"invoice_id,omitempty"`
	Type          TransactionType   `json:"type"`
	Status        TransactionStatus `json:"status"`
	Currency      string            `json:"currency"`
	AmountMinor   int64             `json:"amount_minor"`
	Description   string            `json:"description,omitempty"`
	Metadata      json.RawMessage   `json:"metadata"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

func newTransactionResponse(transaction Transaction) transactionResponse {
	return transactionResponse{
		ID:            transaction.ID,
		DisplayID:     transaction.DisplayID,
		TenantID:      transaction.TenantID,
		AccountUserID: transaction.AccountUserID,
		OrderID:       transaction.OrderID,
		InvoiceID:     transaction.InvoiceID,
		Type:          transaction.Type,
		Status:        transaction.Status,
		Currency:      transaction.Currency,
		AmountMinor:   transaction.AmountMinor,
		Description:   transaction.Description,
		Metadata:      transaction.Metadata,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func newTransactionResponses(transactions []Transaction) []transactionResponse {
	responses := make([]transactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		responses = append(responses, newTransactionResponse(transaction))
	}
	return responses
}

type createInvoiceWalletPaymentRequest struct {
	InvoiceID invoice.InvoiceID `json:"invoice_id"`
	WalletID  wallet.WalletID   `json:"wallet_id"`
}

func (request createInvoiceWalletPaymentRequest) toInput(
	tenantID tenant.ID,
	actorID identity.UserID,
	idempotencyKey IdempotencyKey,
) PayInvoiceFromWalletInput {
	return PayInvoiceFromWalletInput{
		TenantID:       tenantID,
		InvoiceID:      request.InvoiceID,
		WalletID:       request.WalletID,
		ActorID:        actorID,
		IdempotencyKey: idempotencyKey,
	}
}

type invoiceWalletPaymentResponse struct {
	Invoice     invoiceWalletPaymentInvoiceResponse `json:"invoice"`
	Transaction transactionResponse                 `json:"transaction"`
	Ledger      *invoiceWalletPaymentLedgerResponse `json:"ledger,omitempty"`
	Order       *invoiceWalletPaymentOrderResponse  `json:"order,omitempty"`
}

type invoiceWalletPaymentInvoiceResponse struct {
	ID         invoice.InvoiceID `json:"id"`
	DisplayID  int64             `json:"display_id"`
	Status     invoice.Status    `json:"status"`
	TotalMinor int64             `json:"total_minor"`
	Currency   string            `json:"currency"`
	PaidAt     time.Time         `json:"paid_at,omitempty"`
}

type invoiceWalletPaymentLedgerResponse struct {
	ID                wallet.LedgerEntryID `json:"id"`
	DisplayID         int64                `json:"display_id"`
	WalletID          wallet.WalletID      `json:"wallet_id"`
	Direction         wallet.Direction     `json:"direction"`
	EntryType         wallet.EntryType     `json:"entry_type"`
	Status            wallet.LedgerStatus  `json:"status"`
	Currency          string               `json:"currency"`
	AmountMinor       int64                `json:"amount_minor"`
	BalanceAfterMinor int64                `json:"balance_after_minor"`
}

type invoiceWalletPaymentOrderResponse struct {
	ID            order.OrderID       `json:"id"`
	DisplayID     int64               `json:"display_id"`
	OrderStatus   order.OrderStatus   `json:"order_status"`
	BillingStatus order.BillingStatus `json:"billing_status"`
}

func newInvoiceWalletPaymentResponse(result WalletInvoicePayment) invoiceWalletPaymentResponse {
	record := result.Invoice.Invoice
	response := invoiceWalletPaymentResponse{
		Invoice: invoiceWalletPaymentInvoiceResponse{
			ID:         record.ID,
			DisplayID:  record.DisplayID,
			Status:     record.Status,
			TotalMinor: record.TotalMinor,
			Currency:   record.Currency,
			PaidAt:     record.PaidAt,
		},
		Transaction: newTransactionResponse(result.Transaction),
	}
	if !result.LedgerEntry.ID.Empty() {
		response.Ledger = &invoiceWalletPaymentLedgerResponse{
			ID:                result.LedgerEntry.ID,
			DisplayID:         result.LedgerEntry.DisplayID,
			WalletID:          result.LedgerEntry.WalletID,
			Direction:         result.LedgerEntry.Direction,
			EntryType:         result.LedgerEntry.EntryType,
			Status:            result.LedgerEntry.Status,
			Currency:          result.LedgerEntry.Currency,
			AmountMinor:       result.LedgerEntry.AmountMinor,
			BalanceAfterMinor: result.LedgerEntry.BalanceAfterMinor,
		}
	}
	if !result.Order.ID.Empty() {
		response.Order = &invoiceWalletPaymentOrderResponse{
			ID:            result.Order.ID,
			DisplayID:     result.Order.DisplayID,
			OrderStatus:   result.Order.OrderStatus,
			BillingStatus: result.Order.BillingStatus,
		}
	}
	return response
}

type paymentReconciliationResponse struct {
	Transaction transactionResponse                   `json:"transaction"`
	Provider    string                                `json:"provider,omitempty"`
	Invoice     *paymentReconciliationInvoiceResponse `json:"invoice,omitempty"`
	Ledger      *paymentReconciliationLedgerResponse  `json:"ledger,omitempty"`
}

type paymentReconciliationInvoiceResponse struct {
	ID         invoice.InvoiceID `json:"id"`
	DisplayID  int64             `json:"display_id"`
	Status     invoice.Status    `json:"status"`
	TotalMinor int64             `json:"total_minor"`
	PaidAt     time.Time         `json:"paid_at,omitempty"`
}

type paymentReconciliationLedgerResponse struct {
	ID                wallet.LedgerEntryID `json:"id"`
	DisplayID         int64                `json:"display_id"`
	WalletID          wallet.WalletID      `json:"wallet_id"`
	WalletDisplayID   int64                `json:"wallet_display_id,omitempty"`
	Direction         wallet.Direction     `json:"direction"`
	EntryType         wallet.EntryType     `json:"entry_type"`
	Status            wallet.LedgerStatus  `json:"status"`
	BalanceAfterMinor int64                `json:"balance_after_minor"`
}

func newPaymentReconciliationResponse(record PaymentReconciliation) paymentReconciliationResponse {
	response := paymentReconciliationResponse{
		Transaction: newTransactionResponse(record.Transaction),
		Provider:    record.Provider,
	}
	if !record.Invoice.Empty() {
		response.Invoice = &paymentReconciliationInvoiceResponse{
			ID:         record.Invoice.ID,
			DisplayID:  record.Invoice.DisplayID,
			Status:     record.Invoice.Status,
			TotalMinor: record.Invoice.TotalMinor,
			PaidAt:     record.Invoice.PaidAt,
		}
	}
	if !record.Ledger.Empty() {
		response.Ledger = &paymentReconciliationLedgerResponse{
			ID:                record.Ledger.ID,
			DisplayID:         record.Ledger.DisplayID,
			WalletID:          record.Ledger.WalletID,
			WalletDisplayID:   record.Ledger.WalletDisplayID,
			Direction:         record.Ledger.Direction,
			EntryType:         record.Ledger.EntryType,
			Status:            record.Ledger.Status,
			BalanceAfterMinor: record.Ledger.BalanceAfterMinor,
		}
	}
	return response
}

func newPaymentReconciliationResponses(records []PaymentReconciliation) []paymentReconciliationResponse {
	responses := make([]paymentReconciliationResponse, 0, len(records))
	for _, record := range records {
		responses = append(responses, newPaymentReconciliationResponse(record))
	}
	return responses
}
