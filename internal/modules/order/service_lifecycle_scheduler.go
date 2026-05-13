package order

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	ServiceLifecycleJobType       jobs.Type          = "service.lifecycle"
	ServiceLifecycleReferenceType jobs.ReferenceType = ServiceAggregateType

	defaultServiceLifecyclePriority = 60
	defaultServiceLifecycleLimit    = 50

	DefaultServiceLifecycleGracePeriod = 72 * time.Hour
)

var (
	ErrServiceLifecycleDueStoreMissing     = errors.New("service lifecycle due store missing")
	ErrServiceLifecycleQueueMissing        = errors.New("service lifecycle queue missing")
	ErrServiceLifecycleTransitionerMissing = errors.New("service lifecycle transitioner missing")
	ErrServiceLifecyclePayloadInvalid      = errors.New("service lifecycle job payload invalid")
	ErrServiceLifecycleLimitInvalid        = errors.New("service lifecycle limit invalid")
	ErrServiceLifecycleGracePeriodInvalid  = errors.New("service lifecycle grace period invalid")
)

type ServiceLifecycleDueStore interface {
	ListDueServiceLifecycleActions(ctx context.Context, input ListDueServiceLifecycleActionsInput) ([]ServiceLifecycleDueAction, error)
}

type ServiceLifecycleTransitioner interface {
	TransitionServiceLifecycle(ctx context.Context, input TransitionServiceLifecycleInput) (ServiceInstance, error)
}

type ServiceLifecycleScheduler struct {
	Store       ServiceLifecycleDueStore
	Queue       jobs.QueueStore
	Now         func() time.Time
	Limit       int
	GracePeriod time.Duration
}

type ListDueServiceLifecycleActionsInput struct {
	Now         time.Time
	Limit       int
	GracePeriod time.Duration
}

type ServiceLifecycleDueAction struct {
	ServiceID                ServiceID
	TenantID                 tenant.ID
	Action                   ServiceLifecycleAction
	FromStatus               ServiceStatus
	ToStatus                 ServiceStatus
	BillingStatus            BillingStatus
	SuspensionReason         SuspensionReason
	ExpectedBillingStatus    BillingStatus
	ExpectedSuspensionReason SuspensionReason
	Reason                   string
	TermEnd                  time.Time
}

type ServiceLifecycleScheduleSummary struct {
	Due       int
	Scheduled int
}

type ServiceLifecycleJobPayload struct {
	ServiceID                ServiceID              `json:"service_id"`
	TenantID                 tenant.ID              `json:"tenant_id"`
	Action                   ServiceLifecycleAction `json:"action"`
	FromStatus               ServiceStatus          `json:"from_status"`
	ToStatus                 ServiceStatus          `json:"to_status"`
	BillingStatus            BillingStatus          `json:"billing_status,omitempty"`
	SuspensionReason         SuspensionReason       `json:"suspension_reason,omitempty"`
	ExpectedBillingStatus    BillingStatus          `json:"expected_billing_status,omitempty"`
	ExpectedSuspensionReason SuspensionReason       `json:"expected_suspension_reason,omitempty"`
	Reason                   string                 `json:"reason,omitempty"`
	TermEnd                  time.Time              `json:"term_end"`
}

type ServiceLifecycleHandler struct {
	Transitioner ServiceLifecycleTransitioner
	Now          func() time.Time
}

func NewServiceLifecycleScheduler(store ServiceLifecycleDueStore, queue jobs.QueueStore) *ServiceLifecycleScheduler {
	return &ServiceLifecycleScheduler{Store: store, Queue: queue}
}

func (scheduler *ServiceLifecycleScheduler) ScheduleDue(ctx context.Context, input ListDueServiceLifecycleActionsInput) (ServiceLifecycleScheduleSummary, error) {
	if err := scheduler.ready(); err != nil {
		return ServiceLifecycleScheduleSummary{}, err
	}
	input = scheduler.normalizeInput(input)
	if err := input.Validate(); err != nil {
		return ServiceLifecycleScheduleSummary{}, err
	}
	actions, err := scheduler.Store.ListDueServiceLifecycleActions(ctx, input)
	if err != nil {
		return ServiceLifecycleScheduleSummary{}, err
	}

	summary := ServiceLifecycleScheduleSummary{Due: len(actions)}
	for _, action := range actions {
		action = action.Normalize()
		if err := action.Validate(); err != nil {
			return summary, err
		}
		payload, err := action.PayloadJSON()
		if err != nil {
			return summary, err
		}
		if _, err := scheduler.Queue.CreateJob(ctx, jobs.CreateJobInput{
			TenantID:       action.TenantID,
			Type:           ServiceLifecycleJobType,
			ReferenceType:  ServiceLifecycleReferenceType,
			ReferenceID:    jobs.ReferenceID(action.ServiceID),
			PayloadJSON:    payload,
			Priority:       defaultServiceLifecyclePriority,
			IdempotencyKey: serviceLifecycleJobKey(action),
			MaxAttempts:    5,
			CorrelationID:  jobs.CorrelationID(action.ServiceID),
		}); err != nil {
			return summary, err
		}
		summary.Scheduled++
	}
	return summary, nil
}

func NewServiceLifecycleHandler(transitioner ServiceLifecycleTransitioner) *ServiceLifecycleHandler {
	return &ServiceLifecycleHandler{Transitioner: transitioner}
}

func NewServiceLifecycleRunner(store jobs.Store, transitioner ServiceLifecycleTransitioner, workerID jobs.WorkerID) jobs.Runner {
	return jobs.Runner{
		Store:     store,
		Handler:   NewServiceLifecycleHandler(transitioner),
		WorkerID:  workerID,
		BatchSize: defaultServiceLifecycleLimit,
		Types:     []jobs.Type{ServiceLifecycleJobType},
	}
}

func (handler *ServiceLifecycleHandler) Handle(ctx context.Context, job jobs.Job) (jobs.Completion, error) {
	if handler == nil || handler.Transitioner == nil {
		return jobs.Completion{}, ErrServiceLifecycleTransitionerMissing
	}
	payload, err := DecodeServiceLifecycleJobPayload(job.PayloadJSON)
	if err != nil || !payload.MatchesJob(job) {
		return jobs.Completion{
			Status:                   jobs.StatusFailedTerminal,
			RetrySafety:              jobs.RetrySafetyDoNotRetry,
			LastErrorCode:            "service_lifecycle_payload_invalid",
			LastErrorMessageRedacted: "service lifecycle job payload is invalid",
			FinishedAt:               handler.now(),
		}, nil
	}
	_, err = handler.Transitioner.TransitionServiceLifecycle(ctx, payload.TransitionInput())
	if err == nil || errors.Is(err, ErrServiceStatusConflict) || errors.Is(err, ErrServiceNotFound) {
		return jobs.Completion{Status: jobs.StatusSucceeded, FinishedAt: handler.now()}, nil
	}
	return jobs.Completion{}, err
}

func (input ListDueServiceLifecycleActionsInput) Validate() error {
	if input.Now.IsZero() {
		return ErrTermWindowInvalid
	}
	if input.Limit <= 0 {
		return ErrServiceLifecycleLimitInvalid
	}
	if input.GracePeriod <= 0 {
		return ErrServiceLifecycleGracePeriodInvalid
	}
	return nil
}

func (action ServiceLifecycleDueAction) Normalize() ServiceLifecycleDueAction {
	output := action
	output.ServiceID = ServiceID(trim(string(output.ServiceID)))
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.Action = ServiceLifecycleAction(trim(string(output.Action)))
	output.FromStatus = ServiceStatus(trim(string(output.FromStatus)))
	output.ToStatus = ServiceStatus(trim(string(output.ToStatus)))
	output.BillingStatus = BillingStatus(trim(string(output.BillingStatus)))
	output.SuspensionReason = SuspensionReason(trim(string(output.SuspensionReason)))
	output.ExpectedBillingStatus = BillingStatus(trim(string(output.ExpectedBillingStatus)))
	output.ExpectedSuspensionReason = SuspensionReason(trim(string(output.ExpectedSuspensionReason)))
	output.Reason = trim(output.Reason)
	return output
}

func (action ServiceLifecycleDueAction) Validate() error {
	if action.TermEnd.IsZero() {
		return ErrTermWindowInvalid
	}
	input := action.TransitionInput()
	if err := input.Validate(); err != nil {
		return err
	}
	return validateServiceLifecycleJobGuards(input)
}

func (action ServiceLifecycleDueAction) TransitionInput() TransitionServiceLifecycleInput {
	return TransitionServiceLifecycleInput{
		ID:                       action.ServiceID,
		TenantID:                 action.TenantID,
		ActorType:                audit.ActorTypeSystem,
		Action:                   action.Action,
		FromStatus:               action.FromStatus,
		ToStatus:                 action.ToStatus,
		BillingStatus:            action.BillingStatus,
		SuspensionReason:         action.SuspensionReason,
		Reason:                   action.Reason,
		TermEnd:                  action.TermEnd,
		ExpectedTermEnd:          action.TermEnd,
		ExpectedBillingStatus:    action.ExpectedBillingStatus,
		ExpectedSuspensionReason: action.ExpectedSuspensionReason,
	}
}

func (action ServiceLifecycleDueAction) PayloadJSON() (json.RawMessage, error) {
	body, err := json.Marshal(ServiceLifecycleJobPayload{
		ServiceID:                action.ServiceID,
		TenantID:                 action.TenantID,
		Action:                   action.Action,
		FromStatus:               action.FromStatus,
		ToStatus:                 action.ToStatus,
		BillingStatus:            action.BillingStatus,
		SuspensionReason:         action.SuspensionReason,
		ExpectedBillingStatus:    action.ExpectedBillingStatus,
		ExpectedSuspensionReason: action.ExpectedSuspensionReason,
		Reason:                   action.Reason,
		TermEnd:                  action.TermEnd,
	})
	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

func DecodeServiceLifecycleJobPayload(body json.RawMessage) (ServiceLifecycleJobPayload, error) {
	var payload ServiceLifecycleJobPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return ServiceLifecycleJobPayload{}, ErrServiceLifecyclePayloadInvalid
	}
	payload = payload.Normalize()
	if err := payload.Validate(); err != nil {
		return ServiceLifecycleJobPayload{}, ErrServiceLifecyclePayloadInvalid
	}
	return payload, nil
}

func (payload ServiceLifecycleJobPayload) Normalize() ServiceLifecycleJobPayload {
	output := payload
	output.ServiceID = ServiceID(trim(string(output.ServiceID)))
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.Action = ServiceLifecycleAction(trim(string(output.Action)))
	output.FromStatus = ServiceStatus(trim(string(output.FromStatus)))
	output.ToStatus = ServiceStatus(trim(string(output.ToStatus)))
	output.BillingStatus = BillingStatus(trim(string(output.BillingStatus)))
	output.SuspensionReason = SuspensionReason(trim(string(output.SuspensionReason)))
	output.ExpectedBillingStatus = BillingStatus(trim(string(output.ExpectedBillingStatus)))
	output.ExpectedSuspensionReason = SuspensionReason(trim(string(output.ExpectedSuspensionReason)))
	output.Reason = trim(output.Reason)
	return output
}

func (payload ServiceLifecycleJobPayload) Validate() error {
	if payload.TermEnd.IsZero() {
		return ErrTermWindowInvalid
	}
	input := payload.TransitionInput()
	if err := input.Validate(); err != nil {
		return err
	}
	return validateServiceLifecycleJobGuards(input)
}

func (payload ServiceLifecycleJobPayload) MatchesJob(job jobs.Job) bool {
	return job.TenantID == payload.TenantID &&
		job.Type == ServiceLifecycleJobType &&
		job.ReferenceType == ServiceLifecycleReferenceType &&
		job.ReferenceID == jobs.ReferenceID(payload.ServiceID)
}

func (payload ServiceLifecycleJobPayload) TransitionInput() TransitionServiceLifecycleInput {
	return TransitionServiceLifecycleInput{
		ID:                       payload.ServiceID,
		TenantID:                 payload.TenantID,
		ActorType:                audit.ActorTypeSystem,
		Action:                   payload.Action,
		FromStatus:               payload.FromStatus,
		ToStatus:                 payload.ToStatus,
		BillingStatus:            payload.BillingStatus,
		SuspensionReason:         payload.SuspensionReason,
		Reason:                   payload.Reason,
		TermEnd:                  payload.TermEnd,
		ExpectedTermEnd:          payload.TermEnd,
		ExpectedBillingStatus:    payload.ExpectedBillingStatus,
		ExpectedSuspensionReason: payload.ExpectedSuspensionReason,
	}
}

func (scheduler *ServiceLifecycleScheduler) ready() error {
	if scheduler == nil || scheduler.Store == nil {
		return ErrServiceLifecycleDueStoreMissing
	}
	if scheduler.Queue == nil {
		return ErrServiceLifecycleQueueMissing
	}
	return nil
}

func (scheduler *ServiceLifecycleScheduler) normalizeInput(input ListDueServiceLifecycleActionsInput) ListDueServiceLifecycleActionsInput {
	output := input
	if output.Now.IsZero() {
		output.Now = scheduler.now()
	}
	if output.Limit == 0 {
		if scheduler.Limit > 0 {
			output.Limit = scheduler.Limit
		} else {
			output.Limit = defaultServiceLifecycleLimit
		}
	}
	if output.GracePeriod == 0 {
		if scheduler.GracePeriod > 0 {
			output.GracePeriod = scheduler.GracePeriod
		} else {
			output.GracePeriod = DefaultServiceLifecycleGracePeriod
		}
	}
	return output
}

func (scheduler *ServiceLifecycleScheduler) now() time.Time {
	if scheduler.Now == nil {
		return time.Now().UTC()
	}
	return scheduler.Now()
}

func (handler *ServiceLifecycleHandler) now() time.Time {
	if handler.Now == nil {
		return time.Now().UTC()
	}
	return handler.Now()
}

func serviceLifecycleJobKey(action ServiceLifecycleDueAction) string {
	return "service_lifecycle:" +
		string(action.TenantID) + ":" +
		string(action.ServiceID) + ":" +
		string(action.Action) + ":" +
		string(action.FromStatus) + ":" +
		string(action.ToStatus) + ":" +
		strconv.FormatInt(action.TermEnd.UTC().UnixNano(), 10)
}

func validateServiceLifecycleJobGuards(input TransitionServiceLifecycleInput) error {
	switch input.Action {
	case ServiceLifecycleActionExpire:
		if input.ExpectedBillingStatus != BillingStatusPaid || input.ExpectedSuspensionReason != "" {
			return ErrServiceStatusTransitionInvalid
		}
	case ServiceLifecycleActionGrace:
		if input.ExpectedBillingStatus != BillingStatusOverdue || input.ExpectedSuspensionReason != "" {
			return ErrServiceStatusTransitionInvalid
		}
	case ServiceLifecycleActionTerminate:
		if input.ExpectedBillingStatus != BillingStatusOverdue || input.ExpectedSuspensionReason != SuspensionReasonExpiry {
			return ErrServiceStatusTransitionInvalid
		}
	default:
		return ErrServiceLifecycleActionInvalid
	}
	return nil
}
