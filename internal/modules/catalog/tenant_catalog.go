package catalog

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type TenantProductStatus string

const (
	TenantProductStatusActive   TenantProductStatus = "active"
	TenantProductStatusHidden   TenantProductStatus = "hidden"
	TenantProductStatusDisabled TenantProductStatus = "disabled"
)

func (status TenantProductStatus) Valid() bool {
	switch status {
	case TenantProductStatusActive, TenantProductStatusHidden, TenantProductStatusDisabled:
		return true
	default:
		return false
	}
}

type TenantProduct struct {
	ID                  TenantProductID
	DisplayID           int64
	TenantID            tenant.ID
	MasterProductID     ProductID
	NameOverride        string
	DescriptionOverride string
	Status              TenantProductStatus
	CloneVersion        int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CreateTenantProductInput struct {
	TenantID            tenant.ID
	MasterProductID     ProductID
	NameOverride        string
	DescriptionOverride string
	Status              TenantProductStatus
	CloneVersion        int
}

func (input CreateTenantProductInput) Normalize() CreateTenantProductInput {
	output := input
	output.NameOverride = trim(output.NameOverride)
	output.DescriptionOverride = trim(output.DescriptionOverride)
	if output.Status == "" {
		output.Status = TenantProductStatusHidden
	}
	if output.CloneVersion == 0 {
		output.CloneVersion = 1
	}
	return output
}

func (input CreateTenantProductInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.MasterProductID.Empty() {
		return ErrProductIDMissing
	}
	if !input.Status.Valid() {
		return ErrTenantProductStatus
	}
	if input.CloneVersion <= 0 {
		return ErrVersionInvalid
	}
	return nil
}

type TenantPlanVisibility string

const (
	TenantPlanVisibilityPublic  TenantPlanVisibility = "public"
	TenantPlanVisibilityHidden  TenantPlanVisibility = "hidden"
	TenantPlanVisibilityPrivate TenantPlanVisibility = "private"
)

func (visibility TenantPlanVisibility) Valid() bool {
	switch visibility {
	case TenantPlanVisibilityPublic, TenantPlanVisibilityHidden, TenantPlanVisibilityPrivate:
		return true
	default:
		return false
	}
}

type TenantPlanStatus string

const (
	TenantPlanStatusActive     TenantPlanStatus = "active"
	TenantPlanStatusDisabled   TenantPlanStatus = "disabled"
	TenantPlanStatusMarginRisk TenantPlanStatus = "margin_risk"
	TenantPlanStatusArchived   TenantPlanStatus = "archived"
)

func (status TenantPlanStatus) Valid() bool {
	switch status {
	case TenantPlanStatusActive, TenantPlanStatusDisabled, TenantPlanStatusMarginRisk, TenantPlanStatusArchived:
		return true
	default:
		return false
	}
}

type TenantPlan struct {
	ID                 TenantPlanID
	DisplayID          int64
	TenantID           tenant.ID
	TenantProductID    TenantProductID
	MasterPlanID       PlanID
	SellingPriceMinor  int64
	ResellerCostMinor  int64
	Currency           string
	MarginPolicy       json.RawMessage
	Visibility         TenantPlanVisibility
	Status             TenantPlanStatus
	CloneVersion       int
	ProductSnapshot    json.RawMessage
	PlanSnapshot       json.RawMessage
	PriceSnapshot      json.RawMessage
	CapabilitySnapshot json.RawMessage
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type CreateTenantPlanInput struct {
	TenantID           tenant.ID
	TenantProductID    TenantProductID
	MasterPlanID       PlanID
	SellingPriceMinor  int64
	ResellerCostMinor  int64
	Currency           string
	MarginPolicy       json.RawMessage
	Visibility         TenantPlanVisibility
	Status             TenantPlanStatus
	CloneVersion       int
	ProductSnapshot    json.RawMessage
	PlanSnapshot       json.RawMessage
	PriceSnapshot      json.RawMessage
	CapabilitySnapshot json.RawMessage
}

func (input CreateTenantPlanInput) Normalize() CreateTenantPlanInput {
	output := input
	output.Currency = upperTrim(output.Currency)
	output.MarginPolicy = defaultJSON(output.MarginPolicy)
	output.ProductSnapshot = defaultJSON(output.ProductSnapshot)
	output.PlanSnapshot = defaultJSON(output.PlanSnapshot)
	output.PriceSnapshot = defaultJSON(output.PriceSnapshot)
	output.CapabilitySnapshot = defaultJSON(output.CapabilitySnapshot)
	if output.Visibility == "" {
		output.Visibility = TenantPlanVisibilityHidden
	}
	if output.Status == "" {
		output.Status = TenantPlanStatusDisabled
	}
	if output.CloneVersion == 0 {
		output.CloneVersion = 1
	}
	return output
}

func (input CreateTenantPlanInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.TenantProductID.Empty() {
		return ErrTenantProductIDMissing
	}
	if input.MasterPlanID.Empty() {
		return ErrPlanIDMissing
	}
	if err := validateMinorAmount(input.SellingPriceMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.ResellerCostMinor); err != nil {
		return err
	}
	if err := validateCurrency(input.Currency); err != nil {
		return err
	}
	if !input.Visibility.Valid() {
		return ErrTenantPlanVisibility
	}
	if !input.Status.Valid() {
		return ErrTenantPlanStatus
	}
	if input.CloneVersion <= 0 {
		return ErrVersionInvalid
	}
	return nil
}
