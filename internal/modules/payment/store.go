package payment

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type TransactionFilter struct {
	TenantID      tenant.ID
	AccountUserID identity.UserID
	OrderID       order.OrderID
	InvoiceID     invoice.InvoiceID
	Type          TransactionType
	Status        TransactionStatus
	Limit         int
}

type TransactionLookup struct {
	ID             TransactionID
	TenantID       tenant.ID
	AccountUserID  identity.UserID
	IdempotencyKey IdempotencyKey
}

type Store interface {
	CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error)
	ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error)
	GetTransaction(ctx context.Context, lookup TransactionLookup) (Transaction, error)
}
