package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type fakeOrderStore struct {
	createOrderInput           CreateOrderInput
	createReservationInput     CreateReservationInput
	createProvisioningJobInput CreateProvisioningJobInput
	createServiceInstanceInput CreateServiceInstanceInput
}

func (store *fakeOrderStore) CreateOrder(_ context.Context, input CreateOrderInput) (Order, error) {
	store.createOrderInput = input
	return Order{TenantID: input.TenantID, OrderStatus: input.OrderStatus, BillingStatus: input.BillingStatus}, nil
}

func (store *fakeOrderStore) CreateReservation(_ context.Context, input CreateReservationInput) (Reservation, error) {
	store.createReservationInput = input
	return Reservation{OrderID: input.OrderID, Status: input.Status}, nil
}

func (store *fakeOrderStore) CreateProvisioningJob(_ context.Context, input CreateProvisioningJobInput) (ProvisioningJob, error) {
	store.createProvisioningJobInput = input
	return ProvisioningJob{OrderID: input.OrderID, Status: input.Status}, nil
}

func (store *fakeOrderStore) CreateServiceInstance(_ context.Context, input CreateServiceInstanceInput) (ServiceInstance, error) {
	store.createServiceInstanceInput = input
	return ServiceInstance{OrderID: input.OrderID, Status: input.Status, BillingStatus: input.BillingStatus}, nil
}

func TestServiceRejectsMissingStore(t *testing.T) {
	_, err := NewService(nil).CreateOrder(context.Background(), CreateOrderInput{})
	if !errors.Is(err, ErrServiceStoreMissing) {
		t.Fatalf("expected missing store error, got %v", err)
	}
}

func TestServiceCreateOrderNormalizesBeforeStore(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	_, err := service.CreateOrder(context.Background(), CreateOrderInput{
		TenantID:       tenant.ID("tenant-1"),
		BuyerUserID:    identity.UserID("buyer-1"),
		TenantPlanID:   catalog.TenantPlanID("tenant-plan-1"),
		Currency:       " usd ",
		UnitPriceMinor: 1000,
		TotalMinor:     1000,
		IdempotencyKey: " order-key-1 ",
	})
	if err != nil {
		t.Fatalf("expected order create: %v", err)
	}
	if store.createOrderInput.Currency != "USD" || store.createOrderInput.Quantity != 1 {
		t.Fatalf("expected normalized order input, got %+v", store.createOrderInput)
	}
	if store.createOrderInput.OrderStatus != OrderStatusPendingPayment || store.createOrderInput.BillingStatus != BillingStatusUnpaid {
		t.Fatalf("expected default statuses, got %+v", store.createOrderInput)
	}
}

func TestServiceCreateProvisioningJobRejectsBadInputBeforeStore(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	_, err := service.CreateProvisioningJob(context.Background(), CreateProvisioningJobInput{
		OrderID:          OrderID("order-1"),
		TenantID:         tenant.ID("tenant-1"),
		ProviderSourceID: catalog.ProviderSourceID("source-1"),
		IdempotencyKey:   "provision-key-1",
	})
	if !errors.Is(err, ErrProviderOperationIDMissing) {
		t.Fatalf("expected provider operation error, got %v", err)
	}
	if store.createProvisioningJobInput.OrderID != "" {
		t.Fatalf("store should not be called, got %+v", store.createProvisioningJobInput)
	}
}

func TestServiceCreateServiceInstanceDelegates(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)
	now := time.Now()

	_, err := service.CreateServiceInstance(context.Background(), CreateServiceInstanceInput{
		TenantID:           tenant.ID("tenant-1"),
		OrderID:            OrderID("order-1"),
		TenantPlanID:       catalog.TenantPlanID("tenant-plan-1"),
		ProviderSourceID:   catalog.ProviderSourceID("source-1"),
		ExternalResourceID: provider.ExternalResourceID(" resource-1 "),
		TermStart:          now,
		TermEnd:            now.Add(30 * 24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("expected service create: %v", err)
	}
	if store.createServiceInstanceInput.ExternalResourceID != provider.ExternalResourceID("resource-1") {
		t.Fatalf("expected trimmed external resource id, got %+v", store.createServiceInstanceInput)
	}
	if store.createServiceInstanceInput.Status != ServiceStatusActive || store.createServiceInstanceInput.BillingStatus != BillingStatusPaid {
		t.Fatalf("expected default statuses, got %+v", store.createServiceInstanceInput)
	}
}
