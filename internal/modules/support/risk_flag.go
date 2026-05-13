package support

import (
	"errors"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrRiskFlagTypeInvalid   = errors.New("risk flag type invalid")
	ErrRiskFlagStatusInvalid = errors.New("risk flag status invalid")
	ErrRiskFlagTargetMissing = errors.New("risk flag target missing")
)

type RiskFlagID string

type RiskFlagType string

const (
	RiskFlagTypeNewAccountHighValue RiskFlagType = "new_account_high_value"
	RiskFlagTypePaymentMismatch     RiskFlagType = "payment_mismatch"
	RiskFlagTypeAbuseHistory        RiskFlagType = "abuse_history"
	RiskFlagTypeManualBlacklist     RiskFlagType = "manual_blacklist"
	RiskFlagTypeProviderRisk        RiskFlagType = "provider_risk"
)

func (flagType RiskFlagType) Valid() bool {
	switch flagType {
	case RiskFlagTypeNewAccountHighValue, RiskFlagTypePaymentMismatch, RiskFlagTypeAbuseHistory,
		RiskFlagTypeManualBlacklist, RiskFlagTypeProviderRisk:
		return true
	default:
		return false
	}
}

type RiskFlagStatus string

const (
	RiskFlagStatusOpen      RiskFlagStatus = "open"
	RiskFlagStatusReviewing RiskFlagStatus = "reviewing"
	RiskFlagStatusCleared   RiskFlagStatus = "cleared"
	RiskFlagStatusConfirmed RiskFlagStatus = "confirmed"
)

func (status RiskFlagStatus) Valid() bool {
	switch status {
	case RiskFlagStatusOpen, RiskFlagStatusReviewing, RiskFlagStatusCleared, RiskFlagStatusConfirmed:
		return true
	default:
		return false
	}
}

type RiskFlag struct {
	ID            RiskFlagID
	DisplayID     int64
	TenantID      tenant.ID
	UserID        identity.UserID
	ServiceID     order.ServiceID
	OrderID       order.OrderID
	FlagType      RiskFlagType
	Severity      AbuseSeverity
	Status        RiskFlagStatus
	NoteRedacted  string
	CreatedBy     identity.UserID
	CorrelationID CorrelationID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateRiskFlagInput struct {
	Actor         identity.Actor
	TenantID      tenant.ID
	UserID        identity.UserID
	ServiceID     order.ServiceID
	OrderID       order.OrderID
	FlagType      RiskFlagType
	Severity      AbuseSeverity
	NoteRedacted  string
	CorrelationID CorrelationID
}

func (input CreateRiskFlagInput) Normalize() CreateRiskFlagInput {
	output := input
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.UserID = identity.UserID(trim(string(output.UserID)))
	output.ServiceID = order.ServiceID(trim(string(output.ServiceID)))
	output.OrderID = order.OrderID(trim(string(output.OrderID)))
	output.FlagType = RiskFlagType(trim(string(output.FlagType)))
	output.Severity = AbuseSeverity(trim(string(output.Severity)))
	output.NoteRedacted = trim(output.NoteRedacted)
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	return output
}

func (input CreateRiskFlagInput) Validate() error {
	if err := input.Actor.Validate(); err != nil {
		return err
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.UserID == "" && input.ServiceID.Empty() && input.OrderID.Empty() {
		return ErrRiskFlagTargetMissing
	}
	if !input.FlagType.Valid() {
		return ErrRiskFlagTypeInvalid
	}
	if !input.Severity.Valid() {
		return ErrAbuseSeverityInvalid
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}
