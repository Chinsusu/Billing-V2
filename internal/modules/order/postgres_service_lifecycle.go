package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const transitionServiceLifecycleSQL = `
WITH updated AS (
UPDATE service_instances
SET status = $4,
    billing_status = COALESCE($5::order_billing_status, billing_status),
    suspension_reason = $6::order_service_suspension_reason,
    term_end = COALESCE($7::timestamptz, term_end),
    updated_at = NOW()
WHERE service_instance_id = $1
  AND tenant_id = $2
  AND status = $3
RETURNING ` + serviceInstanceColumns + `
), event AS (
    INSERT INTO outbox_events (tenant_id, aggregate_type, aggregate_id, event_type, payload_json, dedupe_key, correlation_id)
    SELECT
        tenant_id,
        '` + ServiceAggregateType + `',
        service_instance_id,
        $8,
        jsonb_build_object(
            'service_id', service_instance_id,
            'display_id', display_id,
            'tenant_id', tenant_id,
            'order_id', order_id,
            'from_status', $3::text,
            'to_status', status,
            'billing_status', billing_status,
            'suspension_reason', suspension_reason,
            'term_end', term_end
        ),
        $8 || ':' || service_instance_id::text || ':' || $3::text || ':' || status::text || ':' || EXTRACT(EPOCH FROM updated_at)::text,
        service_instance_id
    FROM updated
    ON CONFLICT (dedupe_key) DO NOTHING
)
SELECT ` + serviceInstanceColumns + ` FROM updated`

const getServiceStatusSQL = `
SELECT status
FROM service_instances
WHERE service_instance_id = $1
  AND tenant_id = $2`

func (store *PostgresStore) TransitionServiceLifecycle(ctx context.Context, input TransitionServiceLifecycleInput) (ServiceInstance, error) {
	if err := store.ready(); err != nil {
		return ServiceInstance{}, err
	}
	input = input.Normalize()
	args, err := transitionServiceLifecycleArgs(input)
	if err != nil {
		return ServiceInstance{}, err
	}
	service, err := scanServiceInstance(store.executor.QueryRowContext(ctx, transitionServiceLifecycleSQL, args...))
	if err == nil {
		return service, nil
	}
	if !errors.Is(err, ErrServiceNotFound) {
		return ServiceInstance{}, err
	}
	currentStatus, lookupErr := store.currentServiceStatus(ctx, input.ID, input.TenantID)
	if lookupErr != nil {
		return ServiceInstance{}, lookupErr
	}
	if currentStatus != input.FromStatus {
		return ServiceInstance{}, ErrServiceStatusConflict
	}
	return ServiceInstance{}, ErrServiceStatusConflict
}

func transitionServiceLifecycleArgs(input TransitionServiceLifecycleInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.ID,
		input.TenantID,
		input.FromStatus,
		input.ToStatus,
		nullableString(string(input.BillingStatus)),
		nullableString(string(input.SuspensionReason)),
		nullableTime(input.TermEnd),
		serviceLifecycleEventType(input.Action),
	}, nil
}

func (store *PostgresStore) currentServiceStatus(ctx context.Context, serviceID ServiceID, tenantID tenant.ID) (ServiceStatus, error) {
	var status string
	if err := store.executor.QueryRowContext(ctx, getServiceStatusSQL, serviceID, tenantID).Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrServiceNotFound
		}
		return "", fmt.Errorf("read service status: %w", err)
	}
	return ServiceStatus(status), nil
}

func nullableTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: value, Valid: true}
}
