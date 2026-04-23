package order

import (
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestPostgresStoreRequiresExecutor(t *testing.T) {
	_, err := NewPostgresStore(nil).CreateOrder(nil, CreateOrderInput{})
	if !errors.Is(err, ErrStoreExecutorMissing) {
		t.Fatalf("expected missing executor error, got %v", err)
	}
}

func TestCreateOrderArgsNormalizeAndValidate(t *testing.T) {
	args, err := createOrderArgs(CreateOrderInput{
		TenantID:       tenant.ID("tenant-1"),
		BuyerUserID:    identity.UserID("buyer-1"),
		TenantPlanID:   catalog.TenantPlanID("tenant-plan-1"),
		Currency:       " usd ",
		UnitPriceMinor: 1000,
		TotalMinor:     1000,
		IdempotencyKey: " order-key-1 ",
	})
	if err != nil {
		t.Fatalf("expected order args: %v", err)
	}
	if len(args) != 14 {
		t.Fatalf("expected 14 args, got %d", len(args))
	}
	if args[3] != 1 || args[4] != "USD" || args[8] != OrderStatusPendingPayment || args[9] != BillingStatusUnpaid {
		t.Fatalf("unexpected normalized args: %#v", args)
	}
	if args[11] != "{}" || args[12] != "{}" || args[13] != "{}" {
		t.Fatalf("expected default JSON snapshots, got %#v", args[11:])
	}
}

func TestCreateOrderArgsRejectsInvalidInput(t *testing.T) {
	_, err := createOrderArgs(CreateOrderInput{})
	if !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant error, got %v", err)
	}
}

func TestCreateOrderSQLEmitsCreatedOutboxEvent(t *testing.T) {
	for _, clause := range []string{
		"WITH created AS",
		"INSERT INTO outbox_events",
		OrderEventCreated,
		"'display_id', display_id",
		"FROM created",
	} {
		if !strings.Contains(createOrderSQL, clause) {
			t.Fatalf("expected %q in create order SQL: %s", clause, createOrderSQL)
		}
	}
}

func TestCreateReservationArgs(t *testing.T) {
	expiresAt := time.Now().Add(5 * time.Minute)
	args, err := createReservationArgs(CreateReservationInput{
		OrderID:          OrderID("order-1"),
		TenantID:         tenant.ID("tenant-1"),
		ProviderSourceID: catalog.ProviderSourceID("source-1"),
		ExpiresAt:        expiresAt,
	})
	if err != nil {
		t.Fatalf("expected reservation args: %v", err)
	}
	if len(args) != 5 || args[3] != ReservationStatusPendingReserve || args[4] != expiresAt {
		t.Fatalf("unexpected reservation args: %#v", args)
	}
}

func TestCreateProvisioningJobArgs(t *testing.T) {
	args, err := createProvisioningJobArgs(CreateProvisioningJobInput{
		OrderID:             OrderID("order-1"),
		TenantID:            tenant.ID("tenant-1"),
		ProviderSourceID:    catalog.ProviderSourceID("source-1"),
		ProviderOperationID: " operation-1 ",
		IdempotencyKey:      " provision-key-1 ",
	})
	if err != nil {
		t.Fatalf("expected provisioning args: %v", err)
	}
	if len(args) != 7 || args[3] != ProviderOperationID("operation-1") ||
		args[4] != ProvisioningStatusQueued || args[6] != 1 {
		t.Fatalf("unexpected provisioning args: %#v", args)
	}
}

func TestRecordProvisioningResultArgs(t *testing.T) {
	args, err := recordProvisioningResultArgs(RecordProvisioningResultInput{
		OrderID:             OrderID("order-1"),
		TenantID:            tenant.ID("tenant-1"),
		ProviderSourceID:    catalog.ProviderSourceID("source-1"),
		ProviderOperationID: " operation-1 ",
		Status:              ProvisioningStatusFailed,
		IdempotencyKey:      " provision-key-1 ",
		LastErrorCode:       " PROVIDER_TEMPORARY_ERROR ",
		LastErrorMessage:    " temporary provider error ",
	})
	if err != nil {
		t.Fatalf("expected provisioning result args: %v", err)
	}
	if len(args) != 9 || args[3] != ProviderOperationID("operation-1") ||
		args[4] != ProvisioningStatusFailed || args[6] != 1 {
		t.Fatalf("unexpected provisioning result args: %#v", args)
	}
}

func TestRecordProvisioningResultSQLUpsertsByIdempotencyKey(t *testing.T) {
	for _, clause := range []string{
		"INSERT INTO order_provisioning_jobs",
		"ON CONFLICT (tenant_id, idempotency_key)",
		"last_error_code = EXCLUDED.last_error_code",
		"RETURNING",
	} {
		if !strings.Contains(recordProvisioningResultSQL, clause) {
			t.Fatalf("expected %q in provisioning result SQL: %s", clause, recordProvisioningResultSQL)
		}
	}
}

func TestCreateServiceInstanceArgsUsesNullableSuspensionReason(t *testing.T) {
	now := time.Now()
	args, err := createServiceInstanceArgs(CreateServiceInstanceInput{
		TenantID:           tenant.ID("tenant-1"),
		OrderID:            OrderID("order-1"),
		TenantPlanID:       catalog.TenantPlanID("tenant-plan-1"),
		ProviderSourceID:   catalog.ProviderSourceID("source-1"),
		ExternalResourceID: provider.ExternalResourceID("resource-1"),
		TermStart:          now,
		TermEnd:            now.Add(30 * 24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("expected service args: %v", err)
	}
	reason, ok := args[7].(sql.NullString)
	if !ok {
		t.Fatalf("expected nullable suspension reason, got %T", args[7])
	}
	if reason.Valid {
		t.Fatalf("expected empty suspension reason to be null, got %#v", reason)
	}
	if len(args) != 10 || args[5] != ServiceStatusActive || args[6] != BillingStatusPaid {
		t.Fatalf("unexpected service args: %#v", args)
	}
}
