package catalog

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

type ProviderSourceStatus string

const (
	ProviderSourceStatusActive      ProviderSourceStatus = "active"
	ProviderSourceStatusDisabled    ProviderSourceStatus = "disabled"
	ProviderSourceStatusMaintenance ProviderSourceStatus = "maintenance"
	ProviderSourceStatusOutOfStock  ProviderSourceStatus = "out_of_stock"
)

func (status ProviderSourceStatus) Valid() bool {
	switch status {
	case ProviderSourceStatusActive, ProviderSourceStatusDisabled, ProviderSourceStatusMaintenance, ProviderSourceStatusOutOfStock:
		return true
	default:
		return false
	}
}

type InventoryMode string

const (
	InventoryModeFiniteStock     InventoryMode = "finite_stock"
	InventoryModeProviderLive    InventoryMode = "provider_live"
	InventoryModeManualUnlimited InventoryMode = "manual_unlimited"
	InventoryModePreloadedList   InventoryMode = "preloaded_list"
)

func (mode InventoryMode) Valid() bool {
	switch mode {
	case InventoryModeFiniteStock, InventoryModeProviderLive, InventoryModeManualUnlimited, InventoryModePreloadedList:
		return true
	default:
		return false
	}
}

type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

func (level RiskLevel) Valid() bool {
	switch level {
	case RiskLevelLow, RiskLevelMedium, RiskLevelHigh:
		return true
	default:
		return false
	}
}

type ProviderSource struct {
	ID                ProviderSourceID
	DisplayID         int64
	Type              provider.Type
	Name              string
	ProviderAccountID provider.AccountID
	Location          string
	Status            ProviderSourceStatus
	CapabilityProfile provider.CapabilityProfile
	InventoryMode     InventoryMode
	RiskLevel         RiskLevel
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateProviderSourceInput struct {
	Type              provider.Type
	Name              string
	ProviderAccountID provider.AccountID
	Location          string
	Status            ProviderSourceStatus
	CapabilityProfile provider.CapabilityProfile
	InventoryMode     InventoryMode
	RiskLevel         RiskLevel
}

func providerTypeValid(providerType provider.Type) bool {
	switch providerType {
	case provider.TypeManual, provider.TypeProxmox, provider.TypeOVH, provider.TypeHetzner,
		provider.TypeProxyUpstream, provider.TypePreloadedProxyPool, provider.TypeCustomAPI:
		return true
	default:
		return false
	}
}

func (input CreateProviderSourceInput) Normalize() CreateProviderSourceInput {
	output := input
	output.Name = trim(output.Name)
	output.Location = trim(output.Location)
	if output.Status == "" {
		output.Status = ProviderSourceStatusDisabled
	}
	if output.RiskLevel == "" {
		output.RiskLevel = RiskLevelMedium
	}
	if output.CapabilityProfile == (provider.CapabilityProfile{}) && output.Type != "" {
		output.CapabilityProfile = provider.DefaultCapabilityProfile(output.Type)
	}
	return output
}

func (input CreateProviderSourceInput) Validate() error {
	if !providerTypeValid(input.Type) {
		return ErrSourceTypeInvalid
	}
	if input.Name == "" {
		return ErrSourceNameMissing
	}
	if !input.Status.Valid() {
		return ErrSourceStatusInvalid
	}
	if !input.InventoryMode.Valid() {
		return ErrInventoryModeInvalid
	}
	if !input.RiskLevel.Valid() {
		return ErrRiskLevelInvalid
	}
	return nil
}

type PlanSourceStatus string

const (
	PlanSourceStatusActive   PlanSourceStatus = "active"
	PlanSourceStatusDisabled PlanSourceStatus = "disabled"
)

func (status PlanSourceStatus) Valid() bool {
	switch status {
	case PlanSourceStatusActive, PlanSourceStatusDisabled:
		return true
	default:
		return false
	}
}

type PlanSource struct {
	ID                 PlanSourceID
	DisplayID          int64
	PlanID             PlanID
	SourceID           ProviderSourceID
	Priority           int
	CostOverrideMinor  int64
	CapacityPolicy     json.RawMessage
	CapabilityOverride json.RawMessage
	Status             PlanSourceStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type CreatePlanSourceInput struct {
	PlanID             PlanID
	SourceID           ProviderSourceID
	Priority           int
	CostOverrideMinor  int64
	CapacityPolicy     json.RawMessage
	CapabilityOverride json.RawMessage
	Status             PlanSourceStatus
}

func (input CreatePlanSourceInput) Normalize() CreatePlanSourceInput {
	output := input
	output.CapacityPolicy = defaultJSON(output.CapacityPolicy)
	output.CapabilityOverride = defaultJSON(output.CapabilityOverride)
	if output.Status == "" {
		output.Status = PlanSourceStatusDisabled
	}
	return output
}

func (input CreatePlanSourceInput) Validate() error {
	if input.PlanID.Empty() {
		return ErrPlanIDMissing
	}
	if input.SourceID.Empty() {
		return ErrSourceIDMissing
	}
	if input.Priority <= 0 {
		return ErrPlanSourcePriority
	}
	if err := validateMinorAmount(input.CostOverrideMinor); err != nil {
		return err
	}
	if !input.Status.Valid() {
		return ErrPlanSourceStatus
	}
	return nil
}
