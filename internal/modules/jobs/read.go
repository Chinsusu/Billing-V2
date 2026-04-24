package jobs

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const defaultJobListLimit = 100
const maxJobListLimit = 500

type Filter struct {
	TenantID      tenant.ID
	DisplayID     int64
	Type          Type
	Status        Status
	ReferenceType ReferenceType
	ReferenceID   ReferenceID
	SourceID      SourceID
	Limit         int
}

type Lookup struct {
	ID       ID
	TenantID tenant.ID
}

type AttemptFilter struct {
	JobID    ID
	TenantID tenant.ID
	Limit    int
}

type ReadStore interface {
	ListJobs(ctx context.Context, filter Filter) ([]Job, error)
	GetJob(ctx context.Context, lookup Lookup) (Job, error)
	ListAttempts(ctx context.Context, filter AttemptFilter) ([]Attempt, error)
}

type Service struct {
	store ReadStore
}

func NewService(store ReadStore) *Service {
	return &Service{store: store}
}

func (service *Service) ListJobs(ctx context.Context, filter Filter) ([]Job, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeFilter(filter)
	if err := validateFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListJobs(ctx, filter)
}

func (service *Service) GetJob(ctx context.Context, lookup Lookup) (Job, error) {
	if err := service.ready(); err != nil {
		return Job{}, err
	}
	lookup = normalizeLookup(lookup)
	if err := validateLookup(lookup); err != nil {
		return Job{}, err
	}
	return service.store.GetJob(ctx, lookup)
}

func (service *Service) ListAttempts(ctx context.Context, filter AttemptFilter) ([]Attempt, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	filter = normalizeAttemptFilter(filter)
	if err := validateAttemptFilter(filter); err != nil {
		return nil, err
	}
	return service.store.ListAttempts(ctx, filter)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}

func normalizeFilter(filter Filter) Filter {
	filter.TenantID = tenant.ID(trimJobString(string(filter.TenantID)))
	filter.Type = Type(trimJobString(string(filter.Type)))
	filter.Status = Status(trimJobString(string(filter.Status)))
	filter.ReferenceType = ReferenceType(trimJobString(string(filter.ReferenceType)))
	filter.ReferenceID = ReferenceID(trimJobString(string(filter.ReferenceID)))
	filter.SourceID = SourceID(trimJobString(string(filter.SourceID)))
	if filter.Limit <= 0 {
		filter.Limit = defaultJobListLimit
	}
	if filter.Limit > maxJobListLimit {
		filter.Limit = maxJobListLimit
	}
	return filter
}

func validateFilter(filter Filter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrStatusInvalid
	}
	return nil
}

func normalizeAttemptFilter(filter AttemptFilter) AttemptFilter {
	filter.JobID = ID(trimJobString(string(filter.JobID)))
	filter.TenantID = tenant.ID(trimJobString(string(filter.TenantID)))
	if filter.Limit <= 0 {
		filter.Limit = defaultJobListLimit
	}
	if filter.Limit > maxJobListLimit {
		filter.Limit = maxJobListLimit
	}
	return filter
}

func validateAttemptFilter(filter AttemptFilter) error {
	if filter.JobID == "" {
		return ErrJobIDMissing
	}
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	return nil
}

func normalizeLookup(lookup Lookup) Lookup {
	return Lookup{
		ID:       ID(trimJobString(string(lookup.ID))),
		TenantID: tenant.ID(trimJobString(string(lookup.TenantID))),
	}
}

func validateLookup(lookup Lookup) error {
	if lookup.ID == "" {
		return ErrJobIDMissing
	}
	if lookup.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	return nil
}
