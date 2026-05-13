package order

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type ServiceLifecycleAction string

const (
	ServiceLifecycleActionRenew     ServiceLifecycleAction = "renew"
	ServiceLifecycleActionExpire    ServiceLifecycleAction = "expire"
	ServiceLifecycleActionGrace     ServiceLifecycleAction = "grace"
	ServiceLifecycleActionSuspend   ServiceLifecycleAction = "suspend"
	ServiceLifecycleActionUnsuspend ServiceLifecycleAction = "unsuspend"
	ServiceLifecycleActionTerminate ServiceLifecycleAction = "terminate"
)

func (action ServiceLifecycleAction) Valid() bool {
	switch action {
	case ServiceLifecycleActionRenew,
		ServiceLifecycleActionExpire,
		ServiceLifecycleActionGrace,
		ServiceLifecycleActionSuspend,
		ServiceLifecycleActionUnsuspend,
		ServiceLifecycleActionTerminate:
		return true
	default:
		return false
	}
}

type TransitionServiceLifecycleInput struct {
	ID                       ServiceID
	TenantID                 tenant.ID
	BuyerUserID              identity.UserID
	ActorID                  audit.ActorID
	ActorType                audit.ActorType
	Action                   ServiceLifecycleAction
	FromStatus               ServiceStatus
	ToStatus                 ServiceStatus
	BillingStatus            BillingStatus
	SuspensionReason         SuspensionReason
	Reason                   string
	TermEnd                  time.Time
	ExpectedTermEnd          time.Time
	ExpectedBillingStatus    BillingStatus
	ExpectedSuspensionReason SuspensionReason
}

func (input TransitionServiceLifecycleInput) Normalize() TransitionServiceLifecycleInput {
	output := input
	output.ID = ServiceID(trim(string(output.ID)))
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.BuyerUserID = identity.UserID(trim(string(output.BuyerUserID)))
	output.ActorID = audit.ActorID(trim(string(output.ActorID)))
	if output.ActorType == "" {
		output.ActorType = audit.ActorTypeUser
	}
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

func (input TransitionServiceLifecycleInput) Validate() error {
	if input.ID.Empty() {
		return ErrServiceIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if !input.ActorType.Valid() {
		return audit.ErrActorTypeInvalid
	}
	if input.ActorType == audit.ActorTypeUser && input.ActorID.Empty() {
		return identity.ErrActorIDMissing
	}
	if !input.Action.Valid() {
		return ErrServiceLifecycleActionInvalid
	}
	if !input.FromStatus.Valid() || !input.ToStatus.Valid() {
		return ErrServiceStatusInvalid
	}
	if input.BillingStatus != "" && !input.BillingStatus.Valid() {
		return ErrBillingStatusInvalid
	}
	if input.SuspensionReason != "" && !input.SuspensionReason.Valid() {
		return ErrSuspensionReasonInvalid
	}
	if input.ExpectedBillingStatus != "" && !input.ExpectedBillingStatus.Valid() {
		return ErrBillingStatusInvalid
	}
	if input.ExpectedSuspensionReason != "" && !input.ExpectedSuspensionReason.Valid() {
		return ErrSuspensionReasonInvalid
	}
	if !CanTransitionService(input.FromStatus, input.ToStatus) {
		return ErrServiceStatusTransitionInvalid
	}
	return validateServiceLifecycleAction(input)
}

func validateServiceLifecycleAction(input TransitionServiceLifecycleInput) error {
	switch input.Action {
	case ServiceLifecycleActionRenew:
		if input.ToStatus != ServiceStatusActive || input.BillingStatus != BillingStatusPaid || input.TermEnd.IsZero() {
			return ErrServiceStatusTransitionInvalid
		}
		return requireServiceLifecycleReason(input)
	case ServiceLifecycleActionExpire:
		if input.FromStatus != ServiceStatusActive || input.ToStatus != ServiceStatusExpired || input.BillingStatus != BillingStatusOverdue {
			return ErrServiceStatusTransitionInvalid
		}
	case ServiceLifecycleActionGrace:
		if input.FromStatus != ServiceStatusExpired ||
			input.ToStatus != ServiceStatusSuspended ||
			input.BillingStatus != BillingStatusOverdue ||
			input.SuspensionReason != SuspensionReasonExpiry {
			return ErrServiceStatusTransitionInvalid
		}
	case ServiceLifecycleActionSuspend:
		if input.FromStatus != ServiceStatusActive || input.ToStatus != ServiceStatusSuspended || input.SuspensionReason == "" {
			return ErrServiceStatusTransitionInvalid
		}
		if input.SuspensionReason == SuspensionReasonExpiry {
			return ErrSuspensionReasonInvalid
		}
		return requireServiceLifecycleReason(input)
	case ServiceLifecycleActionUnsuspend:
		if input.FromStatus != ServiceStatusSuspended || input.ToStatus != ServiceStatusActive || input.BillingStatus != BillingStatusPaid {
			return ErrServiceStatusTransitionInvalid
		}
		return requireServiceLifecycleReason(input)
	case ServiceLifecycleActionTerminate:
		if input.FromStatus == ServiceStatusTerminated || input.ToStatus != ServiceStatusTerminated {
			return ErrServiceStatusTransitionInvalid
		}
		return requireServiceLifecycleReason(input)
	}
	return nil
}

func requireServiceLifecycleReason(input TransitionServiceLifecycleInput) error {
	if input.Reason == "" {
		return ErrServiceLifecycleReasonMissing
	}
	return nil
}

type RenewServiceTermInput struct {
	ID          ServiceID
	TenantID    tenant.ID
	BuyerUserID identity.UserID
	ActorID     audit.ActorID
	ActorType   audit.ActorType
	FromStatus  ServiceStatus
	Cycle       ServiceRenewalCycle
	Reason      string
}

func (input RenewServiceTermInput) Normalize() RenewServiceTermInput {
	output := input
	output.ID = ServiceID(trim(string(output.ID)))
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.BuyerUserID = identity.UserID(trim(string(output.BuyerUserID)))
	output.ActorID = audit.ActorID(trim(string(output.ActorID)))
	if output.ActorType == "" {
		output.ActorType = audit.ActorTypeUser
	}
	output.FromStatus = ServiceStatus(trim(string(output.FromStatus)))
	output.Reason = trim(output.Reason)
	return output
}

func (input RenewServiceTermInput) Validate() error {
	if input.ID.Empty() {
		return ErrServiceIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if !input.ActorType.Valid() {
		return audit.ErrActorTypeInvalid
	}
	if input.ActorType == audit.ActorTypeUser && input.ActorID.Empty() {
		return identity.ErrActorIDMissing
	}
	if !input.FromStatus.Valid() {
		return ErrServiceStatusInvalid
	}
	if input.Reason == "" {
		return ErrServiceLifecycleReasonMissing
	}
	return input.Cycle.Validate()
}

type ServiceRenewalCycle struct {
	Type  catalog.BillingCycleType
	Value int
}

func (cycle ServiceRenewalCycle) Validate() error {
	if !cycle.Type.Valid() {
		return catalog.ErrBillingCycleInvalid
	}
	if cycle.Type == catalog.BillingCycleCustom {
		return catalog.ErrBillingCycleInvalid
	}
	if cycle.Value <= 0 {
		return catalog.ErrBillingCycleValue
	}
	return nil
}

func CalculateRenewedTermEnd(service ServiceInstance, cycle ServiceRenewalCycle) (time.Time, error) {
	if err := cycle.Validate(); err != nil {
		return time.Time{}, err
	}
	switch service.Status {
	case ServiceStatusActive, ServiceStatusSuspended, ServiceStatusExpired:
	default:
		return time.Time{}, ErrServiceStatusTransitionInvalid
	}
	if service.Status == ServiceStatusSuspended && service.SuspensionReason != SuspensionReasonExpiry {
		return time.Time{}, ErrServiceStatusTransitionInvalid
	}
	if service.TermEnd.IsZero() {
		return time.Time{}, ErrTermWindowInvalid
	}
	return addRenewalCycle(service.TermEnd, cycle), nil
}

func addRenewalCycle(base time.Time, cycle ServiceRenewalCycle) time.Time {
	switch cycle.Type {
	case catalog.BillingCycleDay:
		return base.AddDate(0, 0, cycle.Value)
	case catalog.BillingCycleMonth30Days:
		return base.Add(time.Duration(cycle.Value*30) * 24 * time.Hour)
	case catalog.BillingCycleCalendarMonth:
		return addCalendarMonths(base, cycle.Value)
	default:
		return base.AddDate(0, 0, cycle.Value)
	}
}

func addCalendarMonths(base time.Time, months int) time.Time {
	year, month, day := base.Date()
	hour, minute, second := base.Clock()
	nanosecond := base.Nanosecond()
	targetMonth := month + time.Month(months)
	lastDay := time.Date(year, targetMonth+1, 0, hour, minute, second, nanosecond, base.Location()).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(year, targetMonth, day, hour, minute, second, nanosecond, base.Location())
}
