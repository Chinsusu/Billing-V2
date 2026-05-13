package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const reserveInventorySQL = `
WITH existing AS (
    SELECT ` + reservationColumns + `
    FROM order_reservations reservation
    WHERE reservation.order_id = $1
      AND reservation.tenant_id = $2
      AND reservation.provider_source_id = $3
      AND reservation.status = '` + string(ReservationStatusReserved) + `'
    ORDER BY reservation.created_at ASC
    LIMIT 1
), updated_inventory AS (
    UPDATE provider_inventory inventory
    SET reserved_count = inventory.reserved_count + $4,
        available_count_cache = inventory.capacity_total - (inventory.reserved_count + $4) - inventory.allocated_count,
        status = CASE
            WHEN inventory.capacity_total - (inventory.reserved_count + $4) - inventory.allocated_count = 0 THEN 'out_of_stock'::provider_inventory_status
            ELSE inventory.status
        END,
        updated_at = NOW()
    WHERE inventory.source_id = $3
      AND inventory.status = 'active'
      AND inventory.capacity_total IS NOT NULL
      AND inventory.capacity_total - inventory.reserved_count - inventory.allocated_count >= $4
      AND NOT EXISTS (SELECT 1 FROM existing)
    RETURNING inventory.source_id
), inserted AS (
    INSERT INTO order_reservations (order_id, tenant_id, provider_source_id, quantity, status, expires_at)
    SELECT $1, $2, $3, $4, '` + string(ReservationStatusReserved) + `', $5
    FROM updated_inventory
    RETURNING ` + reservationColumns + `
)
SELECT ` + reservationColumns + ` FROM inserted
UNION ALL
SELECT ` + reservationColumns + ` FROM existing
LIMIT 1`

const expireReservationsSQL = `
WITH expired AS (
    UPDATE order_reservations reservation
    SET status = '` + string(ReservationStatusExpired) + `',
        updated_at = NOW()
    WHERE reservation.tenant_id = $1
      AND reservation.status = '` + string(ReservationStatusReserved) + `'
      AND reservation.expires_at < $2
    RETURNING reservation.provider_source_id, reservation.quantity
), expired_totals AS (
    SELECT provider_source_id, SUM(quantity)::int AS quantity
    FROM expired
    GROUP BY provider_source_id
), released_inventory AS (
    UPDATE provider_inventory inventory
    SET reserved_count = GREATEST(0, inventory.reserved_count - expired_totals.quantity),
        available_count_cache = CASE
            WHEN inventory.capacity_total IS NULL THEN inventory.available_count_cache
            ELSE inventory.capacity_total - GREATEST(0, inventory.reserved_count - expired_totals.quantity) - inventory.allocated_count
        END,
        status = CASE
            WHEN inventory.status = 'out_of_stock'
             AND inventory.capacity_total IS NOT NULL
             AND inventory.capacity_total - GREATEST(0, inventory.reserved_count - expired_totals.quantity) - inventory.allocated_count > 0
                THEN 'active'::provider_inventory_status
            ELSE inventory.status
        END,
        updated_at = NOW()
    FROM expired_totals
    WHERE inventory.source_id = expired_totals.provider_source_id
    RETURNING inventory.source_id
)
SELECT COUNT(*)::int FROM expired`

func (store *PostgresStore) ReserveInventory(ctx context.Context, input ReserveInventoryInput) (Reservation, error) {
	if err := store.ready(); err != nil {
		return Reservation{}, err
	}
	args, err := reserveInventoryArgs(input)
	if err != nil {
		return Reservation{}, err
	}
	record, err := scanReservation(store.executor.QueryRowContext(ctx, reserveInventorySQL, args...))
	if err == nil {
		return record, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return Reservation{}, ErrReservationOutOfStock
	}
	return Reservation{}, err
}

func (store *PostgresStore) ExpireReservations(ctx context.Context, input ExpireReservationsInput) (int, error) {
	if err := store.ready(); err != nil {
		return 0, err
	}
	args, err := expireReservationsArgs(input)
	if err != nil {
		return 0, err
	}
	var expiredCount int
	if err := store.executor.QueryRowContext(ctx, expireReservationsSQL, args...).Scan(&expiredCount); err != nil {
		return 0, fmt.Errorf("expire order reservations: %w", err)
	}
	return expiredCount, nil
}

func reserveInventoryArgs(input ReserveInventoryInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{input.OrderID, input.TenantID, input.ProviderSourceID, input.Quantity, input.ExpiresAt}, nil
}

func expireReservationsArgs(input ExpireReservationsInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{input.TenantID, input.Now}, nil
}
