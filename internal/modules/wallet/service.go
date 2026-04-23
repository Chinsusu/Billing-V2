package wallet

import "context"

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) ListWallets(ctx context.Context, filter WalletFilter) ([]Wallet, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeWalletFilter(filter)
	if err := validateWalletFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListWallets(ctx, filter)
}

func (service *Service) GetWallet(ctx context.Context, lookup WalletLookup) (Wallet, error) {
	if err := service.ready(); err != nil {
		return Wallet{}, err
	}
	if err := validateWalletLookup(lookup); err != nil {
		return Wallet{}, err
	}
	return service.store.GetWallet(ctx, lookup)
}

func (service *Service) CreateLedgerEntry(ctx context.Context, input CreateLedgerEntryInput) (LedgerEntry, error) {
	if err := service.ready(); err != nil {
		return LedgerEntry{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return LedgerEntry{}, err
	}
	return service.store.CreateLedgerEntry(ctx, input)
}

func (service *Service) ListLedgerEntries(ctx context.Context, filter LedgerEntryFilter) ([]LedgerEntry, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeLedgerEntryFilter(filter)
	if err := validateLedgerEntryFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListLedgerEntries(ctx, filter)
}

func (service *Service) GetLedgerEntry(ctx context.Context, lookup LedgerEntryLookup) (LedgerEntry, error) {
	if err := service.ready(); err != nil {
		return LedgerEntry{}, err
	}
	if err := validateLedgerEntryLookup(lookup); err != nil {
		return LedgerEntry{}, err
	}
	return service.store.GetLedgerEntry(ctx, lookup)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}
