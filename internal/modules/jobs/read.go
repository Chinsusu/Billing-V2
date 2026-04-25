package jobs

import (
	"context"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const defaultJobListLimit = 100
const maxJobListLimit = 500
const defaultJobSummaryType Type = "provider.provision"

type Filter struct {
	TenantID      tenant.ID
	DisplayID     int64
	Type          Type
	Status        Status
	ReferenceType ReferenceType
	ReferenceID   ReferenceID
	SourceID      SourceID
	SourceDisplayID int64
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

type SummaryFilter struct {
	TenantID tenant.ID
	Type     Type
}

type JobStatusCounts struct {
	Queued          int
	Claimed         int
	Running         int
	Succeeded       int
	FailedRetryable int
	FailedTerminal  int
	ManualReview    int
	Cancelled       int
}

func (counts JobStatusCounts) AttentionCount() int {
	return counts.FailedRetryable + counts.FailedTerminal + counts.ManualReview
}

type JobFailureContext struct {
	ID                       ID
	DisplayID                int64
	Status                   Status
	LastErrorCode            string
	LastErrorMessageRedacted string
	ManualReviewReason       string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type JobSummary struct {
	TenantID       tenant.ID
	Type           Type
	Total          int
	AttentionCount int
	Counts         JobStatusCounts
	OldestQueuedAt time.Time
	GeneratedAt    time.Time
	LatestFailure  *JobFailureContext
}

type ReadStore interface {
	ListJobs(ctx context.Context, filter Filter) ([]Job, error)
	GetJob(ctx context.Context, lookup Lookup) (Job, error)
	ListAttempts(ctx context.Context, filter AttemptFilter) ([]Attempt, error)
}

type SummaryStore interface {
	SummarizeJobs(ctx context.Context, filter SummaryFilter) (JobSummary, error)
}

type Service struct {
	store    ReadStore
	recovery RecoveryStore
	audit    AuditAppender
}

func NewService(store ReadStore) *Service {
	service := &Service{store: store}
	if recovery, ok := store.(RecoveryStore); ok {
		service.recovery = recovery
	}
	return service
}

func NewServiceWithAudit(store ReadStore, audit AuditAppender) *Service {
	service := NewService(store)
	service.audit = audit
	return service
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

func (service *Service) SummarizeJobs(ctx context.Context, filter SummaryFilter) (JobSummary, error) {
	if err := service.ready(); err != nil {
		return JobSummary{}, err
	}
	summaryStore, ok := service.store.(SummaryStore)
	if !ok {
		return JobSummary{}, ErrServiceStoreMissing
	}
	filter = normalizeSummaryFilter(filter)
	if err := validateSummaryFilter(filter); err != nil {
		return JobSummary{}, err
	}
	return summaryStore.SummarizeJobs(ctx, filter)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrServiceStoreMissing
	}
	return nil
}

func (service *Service) readyRecovery() error {
	if err := service.ready(); err != nil {
		return err
	}
	if service.recovery == nil {
		return ErrServiceStoreMissing
	}
	return nil
}

func normalizeSummaryFilter(filter SummaryFilter) SummaryFilter {
	filter.TenantID = tenant.ID(trimJobString(string(filter.TenantID)))
	filter.Type = Type(trimJobString(string(filter.Type)))
	if filter.Type == "" {
		filter.Type = defaultJobSummaryType
	}
	return filter
}

func validateSummaryFilter(filter SummaryFilter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
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
