package order

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestServiceReserveInventoryDefaultsTTLAndQuantity(t *testing.T) {
	store := &fakeInventoryReservationStore{capacity: 2}
	now := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
	service := NewServiceWithOptions(ServiceOptions{Store: store, Now: func() time.Time { return now }})

	reservation, err := service.ReserveInventory(context.Background(), ReserveInventoryInput{
		OrderID:          "order-1",
		TenantID:         tenant.ID("tenant-1"),
		ProviderSourceID: catalog.ProviderSourceID("source-1"),
	})
	if err != nil {
		t.Fatalf("expected reservation: %v", err)
	}
	if reservation.Status != ReservationStatusReserved || reservation.Quantity != 1 {
		t.Fatalf("unexpected reservation: %+v", reservation)
	}
	if store.lastReserveInput.Quantity != 1 || !store.lastReserveInput.ExpiresAt.Equal(now.Add(DefaultReservationTTL)) {
		t.Fatalf("expected default quantity and ttl, got %+v", store.lastReserveInput)
	}
}

func TestServiceReserveInventoryRejectsMissingReservationStore(t *testing.T) {
	service := NewService(&fakeOrderStore{})
	_, err := service.ReserveInventory(context.Background(), validReserveInventoryInput("order-1"))
	if !errors.Is(err, ErrServiceStoreMissing) {
		t.Fatalf("expected missing reservation store error, got %v", err)
	}
}

func TestServiceReserveInventoryConcurrentLimitedStock(t *testing.T) {
	store := &fakeInventoryReservationStore{capacity: 1}
	service := NewServiceWithOptions(ServiceOptions{
		Store: store,
		Now:   func() time.Time { return time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC) },
	})
	const attempts = 10
	results := make(chan error, attempts)
	var group sync.WaitGroup
	for index := 0; index < attempts; index++ {
		group.Add(1)
		go func(index int) {
			defer group.Done()
			_, err := service.ReserveInventory(context.Background(), validReserveInventoryInput(OrderID(fmt.Sprintf("order-%d", index))))
			results <- err
		}(index)
	}
	group.Wait()
	close(results)

	successes := 0
	outOfStock := 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrReservationOutOfStock):
			outOfStock++
		default:
			t.Fatalf("unexpected reservation error: %v", err)
		}
	}
	if successes != 1 || outOfStock != attempts-1 || store.reserved != 1 {
		t.Fatalf("expected one reservation and no oversell, successes=%d out_of_stock=%d reserved=%d", successes, outOfStock, store.reserved)
	}
}

func TestServiceExpireReservationsDefaultsNow(t *testing.T) {
	store := &fakeInventoryReservationStore{capacity: 1}
	now := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
	service := NewServiceWithOptions(ServiceOptions{Store: store, Now: func() time.Time { return now }})

	expired, err := service.ExpireReservations(context.Background(), ExpireReservationsInput{TenantID: tenant.ID("tenant-1")})
	if err != nil {
		t.Fatalf("expected expiration: %v", err)
	}
	if expired != 0 || !store.lastExpireInput.Now.Equal(now) {
		t.Fatalf("unexpected expire input/result: %+v expired=%d", store.lastExpireInput, expired)
	}
}

func validReserveInventoryInput(orderID OrderID) ReserveInventoryInput {
	return ReserveInventoryInput{
		OrderID:          orderID,
		TenantID:         tenant.ID("tenant-1"),
		ProviderSourceID: catalog.ProviderSourceID("source-1"),
		ExpiresAt:        time.Date(2026, 5, 13, 10, 5, 0, 0, time.UTC),
	}
}

type fakeInventoryReservationStore struct {
	mu               sync.Mutex
	capacity         int
	reserved         int
	lastReserveInput ReserveInventoryInput
	lastExpireInput  ExpireReservationsInput
}

func (store *fakeInventoryReservationStore) CreateOrder(_ context.Context, input CreateOrderInput) (Order, error) {
	return Order{TenantID: input.TenantID, BuyerUserID: identity.UserID("buyer-1")}, nil
}

func (store *fakeInventoryReservationStore) CreateReservation(_ context.Context, input CreateReservationInput) (Reservation, error) {
	return Reservation{OrderID: input.OrderID, Status: input.Status}, nil
}

func (store *fakeInventoryReservationStore) CreateProvisioningJob(_ context.Context, input CreateProvisioningJobInput) (ProvisioningJob, error) {
	return ProvisioningJob{OrderID: input.OrderID, Status: input.Status}, nil
}

func (store *fakeInventoryReservationStore) CreateServiceInstance(_ context.Context, input CreateServiceInstanceInput) (ServiceInstance, error) {
	return ServiceInstance{OrderID: input.OrderID, Status: input.Status}, nil
}

func (store *fakeInventoryReservationStore) ListOrders(_ context.Context, filter OrderFilter) ([]Order, error) {
	return []Order{{TenantID: filter.TenantID}}, nil
}

func (store *fakeInventoryReservationStore) GetOrder(_ context.Context, lookup OrderLookup) (Order, error) {
	return Order{ID: lookup.ID, TenantID: lookup.TenantID}, nil
}

func (store *fakeInventoryReservationStore) TransitionOrderStatus(_ context.Context, input TransitionOrderStatusInput) (Order, error) {
	return Order{ID: input.ID, TenantID: input.TenantID}, nil
}

func (store *fakeInventoryReservationStore) TransitionServiceLifecycle(_ context.Context, input TransitionServiceLifecycleInput) (ServiceInstance, error) {
	return ServiceInstance{ID: input.ID, TenantID: input.TenantID, Status: input.ToStatus}, nil
}

func (store *fakeInventoryReservationStore) ListServiceInstances(_ context.Context, filter ServiceInstanceFilter) ([]ServiceInstance, error) {
	return []ServiceInstance{{TenantID: filter.TenantID}}, nil
}

func (store *fakeInventoryReservationStore) GetServiceInstance(_ context.Context, lookup ServiceInstanceLookup) (ServiceInstance, error) {
	return ServiceInstance{ID: lookup.ID, TenantID: lookup.TenantID}, nil
}

func (store *fakeInventoryReservationStore) ReserveInventory(_ context.Context, input ReserveInventoryInput) (Reservation, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.lastReserveInput = input
	if store.reserved+input.Quantity > store.capacity {
		return Reservation{}, ErrReservationOutOfStock
	}
	store.reserved += input.Quantity
	return Reservation{
		OrderID:          input.OrderID,
		TenantID:         input.TenantID,
		ProviderSourceID: input.ProviderSourceID,
		Quantity:         input.Quantity,
		Status:           ReservationStatusReserved,
		ExpiresAt:        input.ExpiresAt,
	}, nil
}

func (store *fakeInventoryReservationStore) ExpireReservations(_ context.Context, input ExpireReservationsInput) (int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.lastExpireInput = input
	return 0, nil
}
