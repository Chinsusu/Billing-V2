package invoice

import "context"

type Service struct {
	store       Store
	orderReader OrderReader
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func NewServiceWithOrderReader(store Store, orderReader OrderReader) *Service {
	return &Service{store: store, orderReader: orderReader}
}

func (service *Service) ListInvoices(ctx context.Context, filter InvoiceFilter) ([]Invoice, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeInvoiceFilter(filter)
	if err := validateInvoiceFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListInvoices(ctx, filter)
}

func (service *Service) GetInvoice(ctx context.Context, lookup InvoiceLookup) (InvoiceDetail, error) {
	if err := service.ready(); err != nil {
		return InvoiceDetail{}, err
	}
	if err := validateInvoiceLookup(lookup); err != nil {
		return InvoiceDetail{}, err
	}
	return service.store.GetInvoice(ctx, lookup)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}
