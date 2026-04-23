package order

import (
	"context"
	"fmt"
)

func (store *PostgresStore) ListOrders(ctx context.Context, filter OrderFilter) ([]Order, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListOrdersQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()
	orders := make([]Order, 0)
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read orders: %w", err)
	}
	return orders, nil
}

func (store *PostgresStore) GetOrder(ctx context.Context, lookup OrderLookup) (Order, error) {
	if err := store.ready(); err != nil {
		return Order{}, err
	}
	query, args, err := buildGetOrderQuery(lookup)
	if err != nil {
		return Order{}, err
	}
	return scanOrder(store.executor.QueryRowContext(ctx, query, args...))
}

func buildListOrdersQuery(filter OrderFilter) (string, []interface{}, error) {
	filter = normalizeOrderFilter(filter)
	if err := validateOrderFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + orderColumns + `
FROM orders
WHERE tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.BuyerUserID != "" {
		args = append(args, filter.BuyerUserID)
		query += fmt.Sprintf("\n  AND buyer_user_id = $%d", len(args))
	}
	if filter.OrderStatus != "" {
		args = append(args, filter.OrderStatus)
		query += fmt.Sprintf("\n  AND order_status = $%d", len(args))
	}
	if filter.BillingStatus != "" {
		args = append(args, filter.BillingStatus)
		query += fmt.Sprintf("\n  AND billing_status = $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetOrderQuery(lookup OrderLookup) (string, []interface{}, error) {
	if err := validateOrderLookup(lookup); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + orderColumns + `
FROM orders
WHERE order_id = $1
  AND tenant_id = $2`
	args := []interface{}{lookup.ID, lookup.TenantID}
	if lookup.BuyerUserID != "" {
		args = append(args, lookup.BuyerUserID)
		query += fmt.Sprintf("\n  AND buyer_user_id = $%d", len(args))
	}
	return query, args, nil
}
