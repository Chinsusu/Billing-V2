package order

import (
	"context"
	"net/http"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

func registerOrderTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakeOrderHTTPService struct {
	createOrderCalls                int
	createOrderInput                CreateOrderInput
	listOrderCalls                  int
	orderFilter                     OrderFilter
	getOrderCalls                   int
	orderLookup                     OrderLookup
	transitionOrderStatusCalls      int
	transitionOrderStatusInput      TransitionOrderStatusInput
	transitionOrderStatusError      error
	transitionServiceLifecycleCalls int
	transitionServiceLifecycleInput TransitionServiceLifecycleInput
	transitionServiceLifecycleError error
	renewClientServiceCalls         int
	renewClientServiceInput         ClientServiceRenewalInput
	renewClientServiceResult        ClientServiceRenewal
	renewClientServiceError         error
	listServiceCalls                int
	serviceFilter                   ServiceInstanceFilter
	getServiceCalls                 int
	serviceLookup                   ServiceInstanceLookup
	order                           Order
	orders                          []Order
	service                         ServiceInstance
	services                        []ServiceInstance
}

func (service *fakeOrderHTTPService) CreateOrder(ctx context.Context, input CreateOrderInput) (Order, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Order{}, err
	}
	service.createOrderCalls++
	service.createOrderInput = input
	if service.order.ID != "" {
		return service.order, nil
	}
	return Order{
		ID:             "order_1",
		TenantID:       input.TenantID,
		BuyerUserID:    input.BuyerUserID,
		TenantPlanID:   catalog.TenantPlanID(input.TenantPlanID),
		Quantity:       input.Quantity,
		Currency:       input.Currency,
		UnitPriceMinor: input.UnitPriceMinor,
		DiscountMinor:  input.DiscountMinor,
		TotalMinor:     input.TotalMinor,
		OrderStatus:    input.OrderStatus,
		BillingStatus:  input.BillingStatus,
	}, nil
}

func (service *fakeOrderHTTPService) ListOrders(ctx context.Context, filter OrderFilter) ([]Order, error) {
	service.listOrderCalls++
	service.orderFilter = filter
	return service.orders, nil
}

func (service *fakeOrderHTTPService) GetOrder(ctx context.Context, lookup OrderLookup) (Order, error) {
	service.getOrderCalls++
	service.orderLookup = lookup
	return service.order, nil
}

func (service *fakeOrderHTTPService) TransitionOrderStatus(ctx context.Context, input TransitionOrderStatusInput) (Order, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Order{}, err
	}
	service.transitionOrderStatusCalls++
	service.transitionOrderStatusInput = input
	if service.transitionOrderStatusError != nil {
		return Order{}, service.transitionOrderStatusError
	}
	if service.order.ID != "" {
		return service.order, nil
	}
	return Order{
		ID:            input.ID,
		TenantID:      input.TenantID,
		OrderStatus:   input.ToStatus,
		BillingStatus: input.BillingStatus,
	}, nil
}

func (service *fakeOrderHTTPService) TransitionServiceLifecycle(ctx context.Context, input TransitionServiceLifecycleInput) (ServiceInstance, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ServiceInstance{}, err
	}
	service.transitionServiceLifecycleCalls++
	service.transitionServiceLifecycleInput = input
	if service.transitionServiceLifecycleError != nil {
		return ServiceInstance{}, service.transitionServiceLifecycleError
	}
	if service.service.ID != "" {
		return service.service, nil
	}
	return ServiceInstance{
		ID:               input.ID,
		TenantID:         input.TenantID,
		Status:           input.ToStatus,
		BillingStatus:    input.BillingStatus,
		SuspensionReason: input.SuspensionReason,
	}, nil
}

func (service *fakeOrderHTTPService) RenewClientService(ctx context.Context, input ClientServiceRenewalInput) (ClientServiceRenewal, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ClientServiceRenewal{}, err
	}
	service.renewClientServiceCalls++
	service.renewClientServiceInput = input
	if service.renewClientServiceError != nil {
		return ClientServiceRenewal{}, service.renewClientServiceError
	}
	if service.renewClientServiceResult.Service.ID != "" {
		return service.renewClientServiceResult, nil
	}
	return ClientServiceRenewal{
		Service: ServiceInstance{
			ID:            input.ServiceID,
			TenantID:      input.TenantID,
			Status:        ServiceStatusActive,
			BillingStatus: BillingStatusPaid,
		},
		InvoiceID:                 "invoice_1",
		InvoiceDisplayID:          10001,
		PaymentTransactionID:      "payment_1",
		PaymentTransactionDisplay: 10002,
		WalletID:                  input.WalletID,
		LedgerEntryID:             wallet.LedgerEntryID("ledger_1"),
		LedgerEntryDisplayID:      10003,
		AmountMinor:               2500,
		Currency:                  "USD",
		Renewed:                   true,
		PreviousStatus:            input.FromStatus,
	}, nil
}

func (service *fakeOrderHTTPService) ListServiceInstances(ctx context.Context, filter ServiceInstanceFilter) ([]ServiceInstance, error) {
	service.listServiceCalls++
	service.serviceFilter = filter
	return service.services, nil
}

func (service *fakeOrderHTTPService) GetServiceInstance(ctx context.Context, lookup ServiceInstanceLookup) (ServiceInstance, error) {
	service.getServiceCalls++
	service.serviceLookup = lookup
	return service.service, nil
}
