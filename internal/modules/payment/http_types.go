package payment

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
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
