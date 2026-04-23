package order

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateOrderInputNormalizeValidate(t *testing.T) {
	input := CreateOrderInput{
		TenantID:       tenant.ID("tenant-1"),
		BuyerUserID:    identity.UserID("buyer-1"),
		TenantPlanID:   catalog.TenantPlanID("tenant-plan-1"),
		Currency:       " usd ",
		UnitPriceMinor: 1000,
		TotalMinor:     1000,
		IdempotencyKey: " order-key-1 ",
		PlanSnapshot:   json.RawMessage(`{"name":"VPS"}`),
	}.Normalize()

	if input.Quantity != 1 {
		t.Fatalf("expected default quantity 1, got %d", input.Quantity)
	}
	if input.Currency != "USD" {
		t.Fatalf("expected normalized currency, got %q", input.Currency)
	}
	if input.OrderStatus != OrderStatusPendingPayment {
		t.Fatalf("expected pending payment status, got %q", input.OrderStatus)
	}
	if input.BillingStatus != BillingStatusUnpaid {
		t.Fatalf("expected unpaid billing status, got %q", input.BillingStatus)
	}
	if string(input.ProductSnapshot) != "{}" || string(input.PriceSnapshot) != "{}" {
		t.Fatalf("expected default snapshots, got product=%s price=%s", input.ProductSnapshot, input.PriceSnapshot)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid order input: %v", err)
	}

	input.TotalMinor = -1
	if err := input.Validate(); !errors.Is(err, ErrAmountInvalid) {
		t.Fatalf("expected amount error, got %v", err)
	}
}

func TestCreateOrderInputRequiresTenantPlan(t *testing.T) {
	err := CreateOrderInput{
		TenantID:       tenant.ID("tenant-1"),
		BuyerUserID:    identity.UserID("buyer-1"),
		Currency:       "USD",
		Quantity:       1,
		UnitPriceMinor: 1000,
		TotalMinor:     1000,
		IdempotencyKey: "order-key-1",
	}.Normalize().Validate()
	if !errors.Is(err, ErrTenantPlanIDMissing) {
		t.Fatalf("expected tenant plan error, got %v", err)
	}
}

func TestCreateReservationInputValidate(t *testing.T) {
	input := CreateReservationInput{
		OrderID:          OrderID("order-1"),
		TenantID:         tenant.ID("tenant-1"),
		ProviderSourceID: catalog.ProviderSourceID("source-1"),
		ExpiresAt:        time.Now().Add(5 * time.Minute),
	}.Normalize()

	if input.Status != ReservationStatusPendingReserve {
		t.Fatalf("expected default pending reserve, got %q", input.Status)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid reservation input: %v", err)
	}

	input.ExpiresAt = time.Time{}
	if err := input.Validate(); !errors.Is(err, ErrReservationExpiryMissing) {
		t.Fatalf("expected expiry error, got %v", err)
	}
}

func TestCreateProvisioningJobInputValidate(t *testing.T) {
	input := CreateProvisioningJobInput{
		OrderID:             OrderID("order-1"),
		TenantID:            tenant.ID("tenant-1"),
		ProviderSourceID:    catalog.ProviderSourceID("source-1"),
		ProviderOperationID: " operation-1 ",
		IdempotencyKey:      " provision-key-1 ",
	}.Normalize()

	if input.Status != ProvisioningStatusQueued {
		t.Fatalf("expected queued status, got %q", input.Status)
	}
	if input.AttemptNumber != 1 {
		t.Fatalf("expected first attempt, got %d", input.AttemptNumber)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid provisioning input: %v", err)
	}

	input.ProviderOperationID = ""
	if err := input.Validate(); !errors.Is(err, ErrProviderOperationIDMissing) {
		t.Fatalf("expected provider operation error, got %v", err)
	}
}

func TestCreateServiceInstanceInputValidate(t *testing.T) {
	now := time.Now()
	input := CreateServiceInstanceInput{
		TenantID:           tenant.ID("tenant-1"),
		OrderID:            OrderID("order-1"),
		TenantPlanID:       catalog.TenantPlanID("tenant-plan-1"),
		ProviderSourceID:   catalog.ProviderSourceID("source-1"),
		ExternalResourceID: provider.ExternalResourceID(" resource-1 "),
		TermStart:          now,
		TermEnd:            now.Add(30 * 24 * time.Hour),
	}.Normalize()

	if input.Status != ServiceStatusActive {
		t.Fatalf("expected active status, got %q", input.Status)
	}
	if input.BillingStatus != BillingStatusPaid {
		t.Fatalf("expected paid billing status, got %q", input.BillingStatus)
	}
	if input.ExternalResourceID != provider.ExternalResourceID("resource-1") {
		t.Fatalf("expected trimmed external resource id, got %q", input.ExternalResourceID)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid service input: %v", err)
	}

	input.Status = ServiceStatusSuspended
	if err := input.Validate(); !errors.Is(err, ErrSuspensionReasonInvalid) {
		t.Fatalf("expected suspension reason error, got %v", err)
	}
}
