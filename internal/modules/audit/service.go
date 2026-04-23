package audit

import "context"

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) Append(ctx context.Context, input AppendInput) (Log, error) {
	if err := service.ready(); err != nil {
		return Log{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Log{}, err
	}
	return service.store.Append(ctx, input)
}

func (service *Service) ListLogs(ctx context.Context, filter Filter) ([]Log, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeFilter(filter)
	if err := validateFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListLogs(ctx, filter)
}

func (service *Service) GetLog(ctx context.Context, lookup Lookup) (Log, error) {
	if err := service.ready(); err != nil {
		return Log{}, err
	}
	if err := validateLookup(lookup); err != nil {
		return Log{}, err
	}
	return service.store.GetLog(ctx, lookup)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}
