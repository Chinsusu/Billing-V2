package order

import "context"

type Service struct {
	store Store
	audit AuditAppender
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func NewServiceWithAudit(store Store, audit AuditAppender) *Service {
	return &Service{store: store, audit: audit}
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
	record, err := service.store.TransitionOrderStatus(ctx, input)
	if err != nil {
		return Order{}, err
	}
	if err := service.appendOrderStatusAudit(ctx, input, record); err != nil {
		return Order{}, err
	}
	return record, nil
}

func (service *Service) ListServiceInstances(ctx context.Context, filter ServiceInstanceFilter) ([]ServiceInstance, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeServiceInstanceFilter(filter)
	if err := validateServiceInstanceFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListServiceInstances(ctx, filter)
}

func (service *Service) GetServiceInstance(ctx context.Context, lookup ServiceInstanceLookup) (ServiceInstance, error) {
	if err := service.ready(); err != nil {
		return ServiceInstance{}, err
	}
	if err := validateServiceInstanceLookup(lookup); err != nil {
		return ServiceInstance{}, err
	}
	return service.store.GetServiceInstance(ctx, lookup)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}
