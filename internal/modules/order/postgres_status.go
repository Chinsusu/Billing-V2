package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const transitionOrderStatusSQL = `
UPDATE orders
SET order_status = $4,
    billing_status = $5,
    updated_at = NOW()
WHERE order_id = $1
  AND tenant_id = $2
  AND order_status = $3
RETURNING ` + orderColumns

const getOrderStatusSQL = `
SELECT order_status
FROM orders
WHERE order_id = $1
  AND tenant_id = $2`

func (store *PostgresStore) TransitionOrderStatus(ctx context.Context, input TransitionOrderStatusInput) (Order, error) {
	if err := store.ready(); err != nil {
		return Order{}, err
	}
	input = input.Normalize()
	args, err := transitionOrderStatusArgs(input)
	if err != nil {
		return Order{}, err
	}
	order, err := scanOrder(store.executor.QueryRowContext(ctx, transitionOrderStatusSQL, args...))
	if err == nil {
		return order, nil
	}
	if !errors.Is(err, ErrOrderNotFound) {
		return Order{}, err
	}
	currentStatus, lookupErr := store.currentOrderStatus(ctx, input.ID, input.TenantID)
	if lookupErr != nil {
		return Order{}, lookupErr
	}
	if currentStatus != input.FromStatus {
		return Order{}, ErrOrderStatusConflict
	}
	return Order{}, ErrOrderStatusConflict
}

func transitionOrderStatusArgs(input TransitionOrderStatusInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{input.ID, input.TenantID, input.FromStatus, input.ToStatus, input.BillingStatus}, nil
}

func (store *PostgresStore) currentOrderStatus(ctx context.Context, orderID OrderID, tenantID tenant.ID) (OrderStatus, error) {
	var status string
	if err := store.executor.QueryRowContext(ctx, getOrderStatusSQL, orderID, tenantID).Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrOrderNotFound
		}
		return "", fmt.Errorf("read order status: %w", err)
	}
	return OrderStatus(status), nil
}
