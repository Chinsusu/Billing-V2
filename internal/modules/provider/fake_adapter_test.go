package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func validOperation() OperationContext {
	return OperationContext{
		OperationID:    "op_1",
		TenantID:       tenant.ID("tenant_a"),
		SourceID:       "source_1",
		IdempotencyKey: "tenant_a:order_1:item_1",
		CorrelationID:  "11111111-1111-1111-1111-111111111111",
		AttemptNumber:  1,
	}
}

func TestOperationContextValidate(t *testing.T) {
	if err := validOperation().Validate(); err != nil {
		t.Fatalf("expected valid operation, got %v", err)
	}
}

func TestOperationContextValidateRequiresTenant(t *testing.T) {
	operation := validOperation()
	operation.TenantID = ""

	if err := operation.Validate(); !errors.Is(err, ErrTenantMissing) {
		t.Fatalf("expected tenant error, got %v", err)
	}
}

func TestFakeAdapterDefaultsToSuccess(t *testing.T) {
	adapter := NewFakeAdapter(TypeManual)

	result, err := adapter.Provision(context.Background(), validOperation(), ProvisionRequest{PlanKey: "vps-small"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if result.Status != OperationStatusSuccess {
		t.Fatalf("expected success status, got %s", result.Status)
	}
	if len(adapter.Calls) != 1 || adapter.Calls[0] != OperationProvision {
		t.Fatalf("expected provision call, got %#v", adapter.Calls)
	}
}

func TestFakeAdapterReturnsConfiguredResult(t *testing.T) {
	adapter := NewFakeAdapter(TypeManual)
	observedAt := time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC)
	adapter.SetResult(OperationProvision, OperationResult{
		Status:             OperationStatusSuccess,
		ExternalResourceID: "external_1",
		ServiceIdentifier:  "vps-100",
		ObservedAt:         observedAt,
	})

	result, err := adapter.Provision(context.Background(), validOperation(), ProvisionRequest{})
	if err != nil {
		t.Fatalf("expected configured success, got %v", err)
	}
	if result.ExternalResourceID != "external_1" {
		t.Fatalf("expected configured resource id, got %q", result.ExternalResourceID)
	}
}

func TestFakeAdapterReturnsAdapterErrorResult(t *testing.T) {
	adapter := NewFakeAdapter(TypeManual)
	adapter.SetError(OperationProvision, NewError(ErrorTimeoutUnknown, "provider request timed out"))

	result, err := adapter.Provision(context.Background(), validOperation(), ProvisionRequest{})
	if err == nil {
		t.Fatal("expected provider error")
	}
	if result.Status != OperationStatusUnknown {
		t.Fatalf("expected unknown status, got %s", result.Status)
	}
	if result.RetrySafety != RetrySafetyUnsafeRetry {
		t.Fatalf("expected unsafe retry, got %s", result.RetrySafety)
	}
}

func TestFakeAdapterHealthMapsAuthFailureDown(t *testing.T) {
	adapter := NewFakeAdapter(TypeManual)
	adapter.SetError(OperationCheckHealth, NewError(ErrorAuthFailed, "credential rejected"))

	result, err := adapter.CheckHealth(context.Background(), validOperation(), HealthRequest{})
	if err == nil {
		t.Fatal("expected health error")
	}
	if result.HealthStatus != HealthStatusDown {
		t.Fatalf("expected down health, got %s", result.HealthStatus)
	}
	if result.Result.RetrySafety != RetrySafetyDoNotRetry {
		t.Fatalf("expected do not retry, got %s", result.Result.RetrySafety)
	}
}

func TestFakeAdapterStockDefaultAvailable(t *testing.T) {
	adapter := NewFakeAdapter(TypeManual)

	result, err := adapter.CheckStock(context.Background(), validOperation(), StockRequest{PlanKey: "proxy"})
	if err != nil {
		t.Fatalf("expected stock success, got %v", err)
	}
	if result.StockStatus != StockStatusAvailable {
		t.Fatalf("expected available stock, got %s", result.StockStatus)
	}
}
