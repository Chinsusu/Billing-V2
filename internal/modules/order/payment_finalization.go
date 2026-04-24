package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type FinalizePaymentInput struct {
	ID          OrderID
	TenantID    tenant.ID
	BuyerUserID identity.UserID
}

func (input FinalizePaymentInput) Normalize() FinalizePaymentInput {
	return FinalizePaymentInput{
		ID:          OrderID(trim(string(input.ID))),
		TenantID:    tenant.ID(trim(string(input.TenantID))),
		BuyerUserID: identity.UserID(trim(string(input.BuyerUserID))),
	}
}

func (input FinalizePaymentInput) Validate() error {
	if input.ID.Empty() {
		return ErrOrderIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.BuyerUserID == "" {
		return ErrBuyerIDMissing
	}
	return nil
}

const finalizePaymentSQL = `
WITH updated AS (
UPDATE orders
SET order_status = '` + string(OrderStatusPaid) + `',
    billing_status = '` + string(BillingStatusPaid) + `',
    updated_at = NOW()
WHERE order_id = $1
  AND tenant_id = $2
  AND buyer_user_id = $3
  AND order_status = '` + string(OrderStatusPendingPayment) + `'
  AND billing_status = '` + string(BillingStatusUnpaid) + `'
RETURNING ` + orderColumns + `
), event AS (
    INSERT INTO outbox_events (tenant_id, aggregate_type, aggregate_id, event_type, payload_json, dedupe_key, correlation_id)
    SELECT
        tenant_id,
        '` + OrderAggregateType + `',
        order_id,
        '` + OrderEventStatusChanged + `',
        jsonb_build_object(
            'order_id', order_id,
            'display_id', display_id,
            'tenant_id', tenant_id,
            'from_status', '` + string(OrderStatusPendingPayment) + `',
            'to_status', order_status,
            'billing_status', billing_status
        ),
        '` + OrderEventStatusChanged + `:' || order_id::text || ':payment_paid',
        order_id
    FROM updated
    ON CONFLICT (dedupe_key) DO NOTHING
), already_paid AS (
    SELECT ` + orderColumns + `
    FROM orders
    WHERE order_id = $1
      AND tenant_id = $2
      AND buyer_user_id = $3
      AND order_status = '` + string(OrderStatusPaid) + `'
      AND billing_status = '` + string(BillingStatusPaid) + `'
)
SELECT ` + orderColumns + ` FROM updated
UNION ALL
SELECT ` + orderColumns + ` FROM already_paid
LIMIT 1`

const getOrderPaymentStateSQL = `
SELECT order_status, billing_status
FROM orders
WHERE order_id = $1
  AND tenant_id = $2
  AND buyer_user_id = $3`

func (store *PostgresStore) FinalizePayment(ctx context.Context, input FinalizePaymentInput) (Order, error) {
	if err := store.ready(); err != nil {
		return Order{}, err
	}
	args, err := finalizePaymentArgs(input)
	if err != nil {
		return Order{}, err
	}
	record, err := scanOrder(store.executor.QueryRowContext(ctx, finalizePaymentSQL, args...))
	if err == nil {
		return record, nil
	}
	if !errors.Is(err, ErrOrderNotFound) {
		return Order{}, err
	}
	state, lookupErr := store.currentOrderPaymentState(ctx, input.Normalize())
	if lookupErr != nil {
		return Order{}, lookupErr
	}
	if state.OrderStatus == OrderStatusPaid && state.BillingStatus == BillingStatusPaid {
		return Order{}, ErrOrderStatusConflict
	}
	return Order{}, ErrOrderStatusConflict
}

func finalizePaymentArgs(input FinalizePaymentInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{input.ID, input.TenantID, input.BuyerUserID}, nil
}

type orderPaymentState struct {
	OrderStatus   OrderStatus
	BillingStatus BillingStatus
}

func (store *PostgresStore) currentOrderPaymentState(ctx context.Context, input FinalizePaymentInput) (orderPaymentState, error) {
	var state orderPaymentState
	var orderStatus, billingStatus string
	if err := store.executor.QueryRowContext(ctx, getOrderPaymentStateSQL, input.ID, input.TenantID, input.BuyerUserID).Scan(&orderStatus, &billingStatus); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orderPaymentState{}, ErrOrderNotFound
		}
		return orderPaymentState{}, fmt.Errorf("read order payment state: %w", err)
	}
	state.OrderStatus = OrderStatus(orderStatus)
	state.BillingStatus = BillingStatus(billingStatus)
	return state, nil
}
