package order

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type OrderFilter struct {
	TenantID      tenant.ID
	BuyerUserID   identity.UserID
	OrderStatus   OrderStatus
	BillingStatus BillingStatus
	Limit         int
}

type OrderLookup struct {
	ID          OrderID
	TenantID    tenant.ID
	BuyerUserID identity.UserID
}

type ServiceInstanceFilter struct {
	TenantID    tenant.ID
	BuyerUserID identity.UserID
	OrderID     OrderID
	Status      ServiceStatus
	Limit       int
}

type ServiceInstanceLookup struct {
	ID          ServiceID
	TenantID    tenant.ID
	BuyerUserID identity.UserID
}

type Store interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (Order, error)
	CreateReservation(ctx context.Context, input CreateReservationInput) (Reservation, error)
	CreateProvisioningJob(ctx context.Context, input CreateProvisioningJobInput) (ProvisioningJob, error)
	CreateServiceInstance(ctx context.Context, input CreateServiceInstanceInput) (ServiceInstance, error)
	ListOrders(ctx context.Context, filter OrderFilter) ([]Order, error)
	GetOrder(ctx context.Context, lookup OrderLookup) (Order, error)
	TransitionOrderStatus(ctx context.Context, input TransitionOrderStatusInput) (Order, error)
	ListServiceInstances(ctx context.Context, filter ServiceInstanceFilter) ([]ServiceInstance, error)
	GetServiceInstance(ctx context.Context, lookup ServiceInstanceLookup) (ServiceInstance, error)
}
