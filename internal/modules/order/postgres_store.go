package order

import (
	"context"
	"database/sql"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

type PostgresStore struct {
	executor platformdb.Executor
}

var _ Store = (*PostgresStore)(nil)

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const orderColumns = `order_id, display_id, tenant_id, buyer_user_id, tenant_plan_id, quantity, currency, unit_price_minor, discount_minor, total_minor, order_status, billing_status, idempotency_key, product_snapshot, plan_snapshot, price_snapshot, created_at, updated_at`
const reservationColumns = `reservation_id, display_id, order_id, tenant_id, provider_source_id, quantity, status, expires_at, created_at, updated_at`
const provisioningJobColumns = `provisioning_job_id, display_id, order_id, tenant_id, provider_source_id, provider_operation_id, status, idempotency_key, attempt_number, last_error_code, last_error_message, created_at, updated_at`
const serviceInstanceColumns = `service_instance_id, display_id, tenant_id, order_id, tenant_plan_id, provider_source_id, external_resource_id, status, billing_status, suspension_reason, term_start, term_end, created_at, updated_at`

const createOrderSQL = `
WITH created AS (
INSERT INTO orders (tenant_id, buyer_user_id, tenant_plan_id, quantity, currency, unit_price_minor, discount_minor, total_minor, order_status, billing_status, idempotency_key, product_snapshot, plan_snapshot, price_snapshot)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12::jsonb, $13::jsonb, $14::jsonb)
RETURNING ` + orderColumns + `
), event AS (
    INSERT INTO outbox_events (tenant_id, aggregate_type, aggregate_id, event_type, payload_json, dedupe_key, correlation_id)
    SELECT
        tenant_id,
        '` + OrderAggregateType + `',
        order_id,
        '` + OrderEventCreated + `',
        jsonb_build_object(
            'order_id', order_id,
            'display_id', display_id,
            'tenant_id', tenant_id,
            'buyer_user_id', buyer_user_id,
            'tenant_plan_id', tenant_plan_id,
            'order_status', order_status,
            'billing_status', billing_status,
            'total_minor', total_minor,
            'currency', currency
        ),
        '` + OrderEventCreated + `:' || order_id::text,
        order_id
    FROM created
    ON CONFLICT (dedupe_key) DO NOTHING
)
SELECT ` + orderColumns + ` FROM created`

const createReservationSQL = `
INSERT INTO order_reservations (order_id, tenant_id, provider_source_id, status, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING ` + reservationColumns

const createProvisioningJobSQL = `
INSERT INTO order_provisioning_jobs (order_id, tenant_id, provider_source_id, provider_operation_id, status, idempotency_key, attempt_number)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING ` + provisioningJobColumns

const recordProvisioningResultSQL = `
INSERT INTO order_provisioning_jobs (order_id, tenant_id, provider_source_id, provider_operation_id, status, idempotency_key, attempt_number, last_error_code, last_error_message)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (tenant_id, idempotency_key)
DO UPDATE SET
    status = EXCLUDED.status,
    attempt_number = EXCLUDED.attempt_number,
    last_error_code = EXCLUDED.last_error_code,
    last_error_message = EXCLUDED.last_error_message,
    updated_at = NOW()
RETURNING ` + provisioningJobColumns

const createServiceInstanceSQL = `
INSERT INTO service_instances (tenant_id, order_id, tenant_plan_id, provider_source_id, external_resource_id, status, billing_status, suspension_reason, term_start, term_end)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (order_id)
DO UPDATE SET
    tenant_plan_id = EXCLUDED.tenant_plan_id,
    provider_source_id = EXCLUDED.provider_source_id,
    external_resource_id = EXCLUDED.external_resource_id,
    status = EXCLUDED.status,
    billing_status = EXCLUDED.billing_status,
    suspension_reason = EXCLUDED.suspension_reason,
    term_start = EXCLUDED.term_start,
    term_end = EXCLUDED.term_end,
    updated_at = NOW()
RETURNING ` + serviceInstanceColumns

func (store *PostgresStore) CreateOrder(ctx context.Context, input CreateOrderInput) (Order, error) {
	if err := store.ready(); err != nil {
		return Order{}, err
	}
	args, err := createOrderArgs(input)
	if err != nil {
		return Order{}, err
	}
	return scanOrder(store.executor.QueryRowContext(ctx, createOrderSQL, args...))
}

func (store *PostgresStore) CreateReservation(ctx context.Context, input CreateReservationInput) (Reservation, error) {
	if err := store.ready(); err != nil {
		return Reservation{}, err
	}
	args, err := createReservationArgs(input)
	if err != nil {
		return Reservation{}, err
	}
	return scanReservation(store.executor.QueryRowContext(ctx, createReservationSQL, args...))
}

func (store *PostgresStore) CreateProvisioningJob(ctx context.Context, input CreateProvisioningJobInput) (ProvisioningJob, error) {
	if err := store.ready(); err != nil {
		return ProvisioningJob{}, err
	}
	args, err := createProvisioningJobArgs(input)
	if err != nil {
		return ProvisioningJob{}, err
	}
	return scanProvisioningJob(store.executor.QueryRowContext(ctx, createProvisioningJobSQL, args...))
}

func (store *PostgresStore) RecordProvisioningResult(ctx context.Context, input RecordProvisioningResultInput) (ProvisioningJob, error) {
	if err := store.ready(); err != nil {
		return ProvisioningJob{}, err
	}
	args, err := recordProvisioningResultArgs(input)
	if err != nil {
		return ProvisioningJob{}, err
	}
	return scanProvisioningJob(store.executor.QueryRowContext(ctx, recordProvisioningResultSQL, args...))
}

func (store *PostgresStore) CreateServiceInstance(ctx context.Context, input CreateServiceInstanceInput) (ServiceInstance, error) {
	if err := store.ready(); err != nil {
		return ServiceInstance{}, err
	}
	args, err := createServiceInstanceArgs(input)
	if err != nil {
		return ServiceInstance{}, err
	}
	return scanServiceInstance(store.executor.QueryRowContext(ctx, createServiceInstanceSQL, args...))
}

func createOrderArgs(input CreateOrderInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.TenantID, input.BuyerUserID, input.TenantPlanID, input.Quantity, input.Currency,
		input.UnitPriceMinor, input.DiscountMinor, input.TotalMinor, input.OrderStatus, input.BillingStatus,
		input.IdempotencyKey, string(input.ProductSnapshot), string(input.PlanSnapshot), string(input.PriceSnapshot),
	}, nil
}

func createReservationArgs(input CreateReservationInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{input.OrderID, input.TenantID, input.ProviderSourceID, input.Status, input.ExpiresAt}, nil
}

func createProvisioningJobArgs(input CreateProvisioningJobInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.OrderID, input.TenantID, input.ProviderSourceID, input.ProviderOperationID,
		input.Status, input.IdempotencyKey, input.AttemptNumber,
	}, nil
}

func recordProvisioningResultArgs(input RecordProvisioningResultInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.OrderID, input.TenantID, input.ProviderSourceID, input.ProviderOperationID,
		input.Status, input.IdempotencyKey, input.AttemptNumber,
		nullableString(input.LastErrorCode), nullableString(input.LastErrorMessage),
	}, nil
}

func createServiceInstanceArgs(input CreateServiceInstanceInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.TenantID, input.OrderID, input.TenantPlanID, input.ProviderSourceID, input.ExternalResourceID,
		input.Status, input.BillingStatus, nullableString(string(input.SuspensionReason)), input.TermStart, input.TermEnd,
	}, nil
}

func (store *PostgresStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrStoreExecutorMissing
	}
	return nil
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}
