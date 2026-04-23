package order

import "context"

type Store interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (Order, error)
	CreateReservation(ctx context.Context, input CreateReservationInput) (Reservation, error)
	CreateProvisioningJob(ctx context.Context, input CreateProvisioningJobInput) (ProvisioningJob, error)
	CreateServiceInstance(ctx context.Context, input CreateServiceInstanceInput) (ServiceInstance, error)
}
