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

type Store interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (Order, error)
	CreateReservation(ctx context.Context, input CreateReservationInput) (Reservation, error)
	CreateProvisioningJob(ctx context.Context, input CreateProvisioningJobInput) (ProvisioningJob, error)
	CreateServiceInstance(ctx context.Context, input CreateServiceInstanceInput) (ServiceInstance, error)
	ListOrders(ctx context.Context, filter OrderFilter) ([]Order, error)
	GetOrder(ctx context.Context, lookup OrderLookup) (Order, error)
}
