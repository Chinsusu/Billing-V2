package payment

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type TransactionFilter struct {
	TenantID         tenant.ID
	AccountUserID    identity.UserID
	AccountDisplayID int64
	DisplayID        int64
	OrderID          order.OrderID
	OrderDisplayID   int64
	InvoiceID        invoice.InvoiceID
	InvoiceDisplayID int64
	Type             TransactionType
	Status           TransactionStatus
	AmountMinMinor   *int64
	AmountMaxMinor   *int64
	Limit            int
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
