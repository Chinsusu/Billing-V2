package order

import "context"

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) CreateOrder(ctx context.Context, input CreateOrderInput) (Order, error) {
	if err := service.ready(); err != nil {
		return Order{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Order{}, err
	}
	return service.store.CreateOrder(ctx, input)
}

func (service *Service) CreateReservation(ctx context.Context, input CreateReservationInput) (Reservation, error) {
	if err := service.ready(); err != nil {
		return Reservation{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Reservation{}, err
	}
	return service.store.CreateReservation(ctx, input)
}

func (service *Service) CreateProvisioningJob(ctx context.Context, input CreateProvisioningJobInput) (ProvisioningJob, error) {
	if err := service.ready(); err != nil {
		return ProvisioningJob{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ProvisioningJob{}, err
	}
	return service.store.CreateProvisioningJob(ctx, input)
}

func (service *Service) CreateServiceInstance(ctx context.Context, input CreateServiceInstanceInput) (ServiceInstance, error) {
	if err := service.ready(); err != nil {
		return ServiceInstance{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ServiceInstance{}, err
	}
	return service.store.CreateServiceInstance(ctx, input)
}

func (service *Service) ListOrders(ctx context.Context, filter OrderFilter) ([]Order, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeOrderFilter(filter)
	if err := validateOrderFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListOrders(ctx, filter)
}

func (service *Service) GetOrder(ctx context.Context, lookup OrderLookup) (Order, error) {
	if err := service.ready(); err != nil {
		return Order{}, err
	}
	if err := validateOrderLookup(lookup); err != nil {
		return Order{}, err
	}
	return service.store.GetOrder(ctx, lookup)
}

func (service *Service) TransitionOrderStatus(ctx context.Context, input TransitionOrderStatusInput) (Order, error) {
	if err := service.ready(); err != nil {
		return Order{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Order{}, err
	}
	return service.store.TransitionOrderStatus(ctx, input)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}
