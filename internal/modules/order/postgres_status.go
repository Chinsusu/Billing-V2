package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const transitionOrderStatusSQL = `
WITH updated AS (
UPDATE orders
SET order_status = $4,
    billing_status = $5,
    updated_at = NOW()
WHERE order_id = $1
  AND tenant_id = $2
  AND order_status = $3
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
            'from_status', $3::text,
            'to_status', order_status,
            'billing_status', billing_status
        ),
        '` + OrderEventStatusChanged + `:' || order_id::text || ':' || $3::text || ':' || order_status::text || ':' || EXTRACT(EPOCH FROM updated_at)::text,
        order_id
    FROM updated
    ON CONFLICT (dedupe_key) DO NOTHING
)
SELECT ` + orderColumns + ` FROM updated`

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
