package payment

import "context"

type Service struct {
	store              Store
	invoiceStore       InvoicePaymentStore
	walletService      WalletPaymentService
	walletPaymentStore WalletInvoicePaymentStore
}

func NewService(store Store) *Service {
	service := &Service{store: store}
	if walletPaymentStore, ok := store.(WalletInvoicePaymentStore); ok {
		service.walletPaymentStore = walletPaymentStore
	}
	return service
}

func NewServiceWithBilling(store Store, invoiceStore InvoicePaymentStore, walletService WalletPaymentService) *Service {
	service := NewService(store)
	service.invoiceStore = invoiceStore
	service.walletService = walletService
	return service
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

func (service *Service) ListPaymentReconciliations(ctx context.Context, filter ReconciliationFilter) ([]PaymentReconciliation, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	reconciliationStore, err := service.reconciliationStore()
	if err != nil {
		return nil, err
	}
	filter = normalizeReconciliationFilter(filter)
	if err := validateReconciliationFilter(filter); err != nil {
		return nil, err
	}
	return reconciliationStore.ListPaymentReconciliations(ctx, filter)
}

func (service *Service) GetPaymentReconciliation(ctx context.Context, lookup ReconciliationLookup) (PaymentReconciliation, error) {
	if err := service.ready(); err != nil {
		return PaymentReconciliation{}, err
	}
	reconciliationStore, err := service.reconciliationStore()
	if err != nil {
		return PaymentReconciliation{}, err
	}
	if err := validateReconciliationLookup(lookup); err != nil {
		return PaymentReconciliation{}, err
	}
	return reconciliationStore.GetPaymentReconciliation(ctx, lookup)
}

func (service *Service) PayInvoiceFromWallet(ctx context.Context, input PayInvoiceFromWalletInput) (WalletInvoicePayment, error) {
	if err := service.ready(); err != nil {
		return WalletInvoicePayment{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return WalletInvoicePayment{}, err
	}
	if service.walletPaymentStore != nil {
		return service.walletPaymentStore.PayInvoiceFromWallet(ctx, input)
	}
	if service.invoiceStore == nil || service.walletService == nil {
		return WalletInvoicePayment{}, ErrBillingDependencyMissing
	}
	return payInvoiceFromWallet(ctx, service.store, service.invoiceStore, service.walletService, input)
}

func (service *Service) reconciliationStore() (ReconciliationStore, error) {
	reconciliationStore, ok := service.store.(ReconciliationStore)
	if !ok {
		return nil, ErrBillingDependencyMissing
	}
	return reconciliationStore, nil
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}
