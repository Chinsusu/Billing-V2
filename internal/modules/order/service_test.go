package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type fakeOrderStore struct {
	createOrderInput                CreateOrderInput
	createReservationInput          CreateReservationInput
	createProvisioningJobInput      CreateProvisioningJobInput
	createServiceInstanceInput      CreateServiceInstanceInput
	listOrdersFilter                OrderFilter
	getOrderLookup                  OrderLookup
	transitionOrderStatusInput      TransitionOrderStatusInput
	transitionServiceLifecycleInput TransitionServiceLifecycleInput
	renewClientServiceInput         ClientServiceRenewalInput
	renewClientServiceResult        ClientServiceRenewal
	listServicesFilter              ServiceInstanceFilter
	getServiceLookup                ServiceInstanceLookup
	service                         ServiceInstance
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

func (store *fakeOrderStore) ListOrders(_ context.Context, filter OrderFilter) ([]Order, error) {
	store.listOrdersFilter = filter
	return []Order{{ID: OrderID("order-1"), TenantID: filter.TenantID, BuyerUserID: filter.BuyerUserID}}, nil
}

func (store *fakeOrderStore) GetOrder(_ context.Context, lookup OrderLookup) (Order, error) {
	store.getOrderLookup = lookup
	return Order{ID: lookup.ID, TenantID: lookup.TenantID, BuyerUserID: lookup.BuyerUserID}, nil
}

func (store *fakeOrderStore) TransitionOrderStatus(_ context.Context, input TransitionOrderStatusInput) (Order, error) {
	store.transitionOrderStatusInput = input
	return Order{ID: input.ID, DisplayID: 30001, TenantID: input.TenantID, OrderStatus: input.ToStatus, BillingStatus: input.BillingStatus}, nil
}

func (store *fakeOrderStore) TransitionServiceLifecycle(_ context.Context, input TransitionServiceLifecycleInput) (ServiceInstance, error) {
	store.transitionServiceLifecycleInput = input
	return ServiceInstance{ID: input.ID, DisplayID: 50001, TenantID: input.TenantID, Status: input.ToStatus, BillingStatus: input.BillingStatus, SuspensionReason: input.SuspensionReason, TermEnd: input.TermEnd}, nil
}

func (store *fakeOrderStore) RenewClientService(_ context.Context, input ClientServiceRenewalInput) (ClientServiceRenewal, error) {
	store.renewClientServiceInput = input
	if store.renewClientServiceResult.Service.ID != "" {
		return store.renewClientServiceResult, nil
	}
	return ClientServiceRenewal{
		Service:        ServiceInstance{ID: input.ServiceID, TenantID: input.TenantID, Status: ServiceStatusActive, BillingStatus: BillingStatusPaid},
		InvoiceID:      "invoice-1",
		WalletID:       input.WalletID,
		AmountMinor:    1000,
		Currency:       "USD",
		Renewed:        true,
		PreviousStatus: input.FromStatus,
	}, nil
}

func (store *fakeOrderStore) ListServiceInstances(_ context.Context, filter ServiceInstanceFilter) ([]ServiceInstance, error) {
	store.listServicesFilter = filter
	return []ServiceInstance{{ID: ServiceID("service-1"), TenantID: filter.TenantID}}, nil
}

func (store *fakeOrderStore) GetServiceInstance(_ context.Context, lookup ServiceInstanceLookup) (ServiceInstance, error) {
	store.getServiceLookup = lookup
	if store.service.ID != "" {
		return store.service, nil
	}
	return ServiceInstance{ID: lookup.ID, TenantID: lookup.TenantID}, nil
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

func TestServiceListOrdersValidatesAndDelegates(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	_, err := service.ListOrders(context.Background(), OrderFilter{
		TenantID:    tenant.ID("tenant-1"),
		BuyerUserID: identity.UserID("buyer-1"),
		Limit:       0,
	})
	if err != nil {
		t.Fatalf("expected list orders: %v", err)
	}
	if store.listOrdersFilter.Limit != defaultOrderListLimit {
		t.Fatalf("expected default limit, got %+v", store.listOrdersFilter)
	}
}

func TestServiceGetOrderRequiresOrderID(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	_, err := service.GetOrder(context.Background(), OrderLookup{TenantID: tenant.ID("tenant-1")})
	if !errors.Is(err, ErrOrderIDMissing) {
		t.Fatalf("expected order id error, got %v", err)
	}
	if store.getOrderLookup.TenantID != "" {
		t.Fatalf("store should not be called, got %+v", store.getOrderLookup)
	}
}

func TestServiceTransitionOrderStatusDelegatesAllowedChange(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	order, err := service.TransitionOrderStatus(context.Background(), TransitionOrderStatusInput{
		ID:            " order-1 ",
		TenantID:      tenant.ID(" tenant-1 "),
		ActorID:       " admin-1 ",
		FromStatus:    OrderStatusPendingPayment,
		ToStatus:      OrderStatusPaid,
		BillingStatus: BillingStatusPaid,
	})
	if err != nil {
		t.Fatalf("expected status transition: %v", err)
	}
	if order.OrderStatus != OrderStatusPaid || order.BillingStatus != BillingStatusPaid {
		t.Fatalf("unexpected order result: %+v", order)
	}
	if store.transitionOrderStatusInput.ID != OrderID("order-1") ||
		store.transitionOrderStatusInput.TenantID != tenant.ID("tenant-1") ||
		store.transitionOrderStatusInput.ActorID != identity.UserID("admin-1") {
		t.Fatalf("expected normalized transition input, got %+v", store.transitionOrderStatusInput)
	}
}

func TestServiceTransitionOrderStatusWritesAudit(t *testing.T) {
	store := &fakeOrderStore{}
	auditLog := &fakeOrderAuditAppender{}
	service := NewServiceWithAudit(store, auditLog)

	_, err := service.TransitionOrderStatus(context.Background(), TransitionOrderStatusInput{
		ID:            "order-1",
		TenantID:      tenant.ID("tenant-1"),
		ActorID:       identity.UserID("admin-1"),
		FromStatus:    OrderStatusPendingPayment,
		ToStatus:      OrderStatusPaid,
		BillingStatus: BillingStatusPaid,
	})
	if err != nil {
		t.Fatalf("expected status transition: %v", err)
	}
	if auditLog.calls != 1 ||
		auditLog.input.Action != orderAuditActionStatusChanged ||
		auditLog.input.TargetID != audit.TargetID("order-1") ||
		auditLog.input.ActorID != audit.ActorID("admin-1") {
		t.Fatalf("unexpected audit input: %+v", auditLog.input)
	}
}

func TestServiceTransitionOrderStatusRejectsBadChangeBeforeStore(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	_, err := service.TransitionOrderStatus(context.Background(), TransitionOrderStatusInput{
		ID:            "order-1",
		TenantID:      tenant.ID("tenant-1"),
		ActorID:       "admin-1",
		FromStatus:    OrderStatusPendingPayment,
		ToStatus:      OrderStatusRefunded,
		BillingStatus: BillingStatusRefunded,
	})
	if !errors.Is(err, ErrStatusTransitionInvalid) {
		t.Fatalf("expected transition error, got %v", err)
	}
	if store.transitionOrderStatusInput.ID != "" {
		t.Fatalf("store should not be called, got %+v", store.transitionOrderStatusInput)
	}
}

func TestServiceTransitionServiceLifecycleDelegatesAllowedChange(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	record, err := service.TransitionServiceLifecycle(context.Background(), TransitionServiceLifecycleInput{
		ID:               " service-1 ",
		TenantID:         tenant.ID(" tenant-1 "),
		ActorID:          " admin-1 ",
		Action:           ServiceLifecycleActionSuspend,
		FromStatus:       ServiceStatusActive,
		ToStatus:         ServiceStatusSuspended,
		SuspensionReason: SuspensionReasonManualAdmin,
		Reason:           "abuse ticket AB-1",
	})
	if err != nil {
		t.Fatalf("expected service lifecycle transition: %v", err)
	}
	if record.Status != ServiceStatusSuspended {
		t.Fatalf("unexpected service result: %+v", record)
	}
	if store.transitionServiceLifecycleInput.ID != ServiceID("service-1") ||
		store.transitionServiceLifecycleInput.TenantID != tenant.ID("tenant-1") ||
		store.transitionServiceLifecycleInput.ActorID != audit.ActorID("admin-1") {
		t.Fatalf("expected normalized lifecycle input, got %+v", store.transitionServiceLifecycleInput)
	}
}

func TestServiceTransitionServiceLifecycleWritesAudit(t *testing.T) {
	store := &fakeOrderStore{}
	auditLog := &fakeOrderAuditAppender{}
	service := NewServiceWithAudit(store, auditLog)

	_, err := service.TransitionServiceLifecycle(context.Background(), TransitionServiceLifecycleInput{
		ID:               "service-1",
		TenantID:         tenant.ID("tenant-1"),
		ActorID:          "admin-1",
		Action:           ServiceLifecycleActionSuspend,
		FromStatus:       ServiceStatusActive,
		ToStatus:         ServiceStatusSuspended,
		SuspensionReason: SuspensionReasonManualAdmin,
		Reason:           "abuse ticket AB-1",
	})
	if err != nil {
		t.Fatalf("expected service lifecycle transition: %v", err)
	}
	if auditLog.calls != 1 ||
		auditLog.input.Action != ServiceEventSuspended ||
		auditLog.input.TargetID != audit.TargetID("service-1") ||
		auditLog.input.ActorID != audit.ActorID("admin-1") {
		t.Fatalf("unexpected audit input: %+v", auditLog.input)
	}
}

func TestServiceTransitionServiceLifecycleRejectsBadChangeBeforeStore(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	_, err := service.TransitionServiceLifecycle(context.Background(), TransitionServiceLifecycleInput{
		ID:         "service-1",
		TenantID:   tenant.ID("tenant-1"),
		ActorID:    "admin-1",
		Action:     ServiceLifecycleActionSuspend,
		FromStatus: ServiceStatusTerminated,
		ToStatus:   ServiceStatusSuspended,
		Reason:     "bad transition",
	})
	if !errors.Is(err, ErrServiceStatusTransitionInvalid) {
		t.Fatalf("expected service transition error, got %v", err)
	}
	if store.transitionServiceLifecycleInput.ID != "" {
		t.Fatalf("store should not be called, got %+v", store.transitionServiceLifecycleInput)
	}
}

func TestServiceRenewTermExtendsFromOldTermEndAndActivates(t *testing.T) {
	termStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	termEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	store := &fakeOrderStore{
		service: ServiceInstance{
			ID:               "service-1",
			TenantID:         "tenant-1",
			Status:           ServiceStatusSuspended,
			BillingStatus:    BillingStatusOverdue,
			SuspensionReason: SuspensionReasonExpiry,
			TermStart:        termStart,
			TermEnd:          termEnd,
		},
	}
	service := NewService(store)

	record, err := service.RenewServiceTerm(context.Background(), RenewServiceTermInput{
		ID:         "service-1",
		TenantID:   "tenant-1",
		ActorID:    "admin-1",
		FromStatus: ServiceStatusSuspended,
		Cycle:      ServiceRenewalCycle{Type: catalog.BillingCycleMonth30Days, Value: 1},
		Reason:     "invoice paid",
	})
	if err != nil {
		t.Fatalf("expected renew service term: %v", err)
	}
	if record.Status != ServiceStatusActive ||
		record.BillingStatus != BillingStatusPaid ||
		!record.TermEnd.Equal(termEnd.Add(30*24*time.Hour)) {
		t.Fatalf("unexpected renewed service: %+v", record)
	}
}

type fakeOrderAuditAppender struct {
	input audit.AppendInput
	calls int
}

func (appender *fakeOrderAuditAppender) Append(ctx context.Context, input audit.AppendInput) (audit.Log, error) {
	appender.calls++
	appender.input = input
	return audit.Log{Action: input.Action}, nil
}

func TestServiceListServiceInstancesValidatesAndDelegates(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	_, err := service.ListServiceInstances(context.Background(), ServiceInstanceFilter{
		TenantID:    tenant.ID("tenant-1"),
		BuyerUserID: identity.UserID("buyer-1"),
	})
	if err != nil {
		t.Fatalf("expected list services: %v", err)
	}
	if store.listServicesFilter.Limit != defaultServiceInstanceListLimit ||
		store.listServicesFilter.BuyerUserID != identity.UserID("buyer-1") {
		t.Fatalf("expected normalized service filter, got %+v", store.listServicesFilter)
	}
}

func TestServiceGetServiceInstanceRequiresID(t *testing.T) {
	store := &fakeOrderStore{}
	service := NewService(store)

	_, err := service.GetServiceInstance(context.Background(), ServiceInstanceLookup{TenantID: tenant.ID("tenant-1")})
	if !errors.Is(err, ErrServiceIDMissing) {
		t.Fatalf("expected service id error, got %v", err)
	}
	if store.getServiceLookup.TenantID != "" {
		t.Fatalf("store should not be called, got %+v", store.getServiceLookup)
	}
}
