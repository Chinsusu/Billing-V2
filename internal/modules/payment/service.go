package payment

import "context"

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error) {
	if err := service.ready(); err != nil {
		return Transaction{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Transaction{}, err
	}
	return service.store.CreateTransaction(ctx, input)
}

func (service *Service) ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeTransactionFilter(filter)
	if err := validateTransactionFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListTransactions(ctx, filter)
}

func (service *Service) GetTransaction(ctx context.Context, lookup TransactionLookup) (Transaction, error) {
	if err := service.ready(); err != nil {
		return Transaction{}, err
	}
	if err := validateTransactionLookup(lookup); err != nil {
		return Transaction{}, err
	}
	return service.store.GetTransaction(ctx, lookup)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}
