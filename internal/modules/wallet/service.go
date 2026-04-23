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

func (service *Service) PostLedgerEntry(ctx context.Context, input PostLedgerEntryInput) (LedgerEntry, error) {
	if err := service.ready(); err != nil {
		return LedgerEntry{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return LedgerEntry{}, err
	}
	return service.store.PostLedgerEntry(ctx, input)
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

func (service *Service) CreateTopupRequest(ctx context.Context, input CreateTopupRequestInput) (TopupRequest, error) {
	if err := service.ready(); err != nil {
		return TopupRequest{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return TopupRequest{}, err
	}
	return service.store.CreateTopupRequest(ctx, input)
}

func (service *Service) ListTopupRequests(ctx context.Context, filter TopupRequestFilter) ([]TopupRequest, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeTopupRequestFilter(filter)
	if err := validateTopupRequestFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListTopupRequests(ctx, filter)
}

func (service *Service) GetTopupRequest(ctx context.Context, lookup TopupRequestLookup) (TopupRequest, error) {
	if err := service.ready(); err != nil {
		return TopupRequest{}, err
	}
	if err := validateTopupRequestLookup(lookup); err != nil {
		return TopupRequest{}, err
	}
	return service.store.GetTopupRequest(ctx, lookup)
}

func (service *Service) ApproveTopupRequest(ctx context.Context, input ApproveTopupRequestInput) (TopupRequest, error) {
	if err := service.ready(); err != nil {
		return TopupRequest{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return TopupRequest{}, err
	}
	request, err := service.store.GetTopupRequest(ctx, TopupRequestLookup{ID: input.ID, TenantID: input.TenantID})
	if err != nil {
		return TopupRequest{}, err
	}
	if request.Status == TopupStatusApproved {
		return request, nil
	}
	if !reviewableTopupStatus(request.Status) {
		return TopupRequest{}, ErrTopupStatusConflict
	}
	entry, err := service.store.PostLedgerEntry(ctx, approveLedgerInput(request, input.ReviewedBy))
	if err != nil {
		return TopupRequest{}, err
	}
	return service.store.ApproveTopupRequest(ctx, input, entry.ID)
}

func (service *Service) RejectTopupRequest(ctx context.Context, input RejectTopupRequestInput) (TopupRequest, error) {
	if err := service.ready(); err != nil {
		return TopupRequest{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return TopupRequest{}, err
	}
	request, err := service.store.GetTopupRequest(ctx, TopupRequestLookup{ID: input.ID, TenantID: input.TenantID})
	if err != nil {
		return TopupRequest{}, err
	}
	if request.Status == TopupStatusRejected {
		return request, nil
	}
	if !reviewableTopupStatus(request.Status) {
		return TopupRequest{}, ErrTopupStatusConflict
	}
	return service.store.RejectTopupRequest(ctx, input)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}
