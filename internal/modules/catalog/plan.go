package catalog

import (
	"encoding/json"
	"time"
)

type BillingCycleType string

const (
	BillingCycleDay           BillingCycleType = "day"
	BillingCycleMonth30Days   BillingCycleType = "month_30d"
	BillingCycleCalendarMonth BillingCycleType = "calendar_month"
	BillingCycleCustom        BillingCycleType = "custom"
)

func (cycleType BillingCycleType) Valid() bool {
	switch cycleType {
	case BillingCycleDay, BillingCycleMonth30Days, BillingCycleCalendarMonth, BillingCycleCustom:
		return true
	default:
		return false
	}
}

type BillingCycle struct {
	Type  BillingCycleType
	Value int
}

func (cycle BillingCycle) Validate() error {
	if !cycle.Type.Valid() {
		return ErrBillingCycleInvalid
	}
	if cycle.Value <= 0 {
		return ErrBillingCycleValue
	}
	return nil
}

type PlanStatus string

const (
	PlanStatusDraft    PlanStatus = "draft"
	PlanStatusActive   PlanStatus = "active"
	PlanStatusDisabled PlanStatus = "disabled"
	PlanStatusArchived PlanStatus = "archived"
)

func (status PlanStatus) Valid() bool {
	switch status {
	case PlanStatusDraft, PlanStatusActive, PlanStatusDisabled, PlanStatusArchived:
		return true
	default:
		return false
	}
}

type Plan struct {
	ID                    PlanID
	DisplayID             int64
	ProductID             ProductID
	Code                  string
	Name                  string
	Specs                 json.RawMessage
	BillingCycle          BillingCycle
	BaseCostMinor         int64
	SuggestedPriceMinor   int64
	ResellerMinPriceMinor int64
	Currency              string
	Status                PlanStatus
	Version               int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type CreatePlanInput struct {
	ProductID             ProductID
	Code                  string
	Name                  string
	Specs                 json.RawMessage
	BillingCycle          BillingCycle
	BaseCostMinor         int64
	SuggestedPriceMinor   int64
	ResellerMinPriceMinor int64
	Currency              string
	Status                PlanStatus
	Version               int
}

func (input CreatePlanInput) Normalize() CreatePlanInput {
	output := input
	output.Code = trim(output.Code)
	output.Name = trim(output.Name)
	output.Currency = upperTrim(output.Currency)
	output.Specs = defaultJSON(output.Specs)
	if output.Status == "" {
		output.Status = PlanStatusDraft
	}
	if output.Version == 0 {
		output.Version = 1
	}
	return output
}

func (input CreatePlanInput) Validate() error {
	if input.ProductID.Empty() {
		return ErrProductIDMissing
	}
	if input.Code == "" {
		return ErrPlanCodeMissing
	}
	if input.Name == "" {
		return ErrPlanNameMissing
	}
	if err := input.BillingCycle.Validate(); err != nil {
		return err
	}
	if err := validateMinorAmount(input.BaseCostMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.SuggestedPriceMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.ResellerMinPriceMinor); err != nil {
		return err
	}
	if err := validateCurrency(input.Currency); err != nil {
		return err
	}
	if !input.Status.Valid() {
		return ErrPlanStatusInvalid
	}
	if input.Version <= 0 {
		return ErrVersionInvalid
	}
	return nil
}
