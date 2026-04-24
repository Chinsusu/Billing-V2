package catalog

import "github.com/Chinsusu/Billing-V2/internal/modules/provider"

type ProviderSourceReadinessState string

const (
	ProviderSourceReadinessReady                 ProviderSourceReadinessState = "ready"
	ProviderSourceReadinessInactiveSource        ProviderSourceReadinessState = "inactive_source"
	ProviderSourceReadinessMissingPlanSource     ProviderSourceReadinessState = "missing_plan_source"
	ProviderSourceReadinessUnsupportedCapability ProviderSourceReadinessState = "unsupported_capability"
	ProviderSourceReadinessFakeProviderOnly      ProviderSourceReadinessState = "fake_provider_only"
)

type ProviderSourceReadiness struct {
	PlanDisplayID       int64
	PlanCode            string
	PlanName            string
	ProductType         ProductType
	PlanStatus          PlanStatus
	PlanSourceDisplayID int64
	PlanSourceStatus    PlanSourceStatus
	SourceDisplayID     int64
	SourceName          string
	SourceType          provider.Type
	SourceStatus        ProviderSourceStatus
	InventoryMode       InventoryMode
	State               ProviderSourceReadinessState
	Reason              string
	capabilityProfile   provider.CapabilityProfile
}

func (record ProviderSourceReadiness) withReadinessState() ProviderSourceReadiness {
	switch {
	case record.PlanSourceDisplayID == 0 || record.SourceDisplayID == 0:
		record.State = ProviderSourceReadinessMissingPlanSource
		record.Reason = "Plan has no provider source link."
	case record.PlanSourceStatus != PlanSourceStatusActive || record.SourceStatus != ProviderSourceStatusActive:
		record.State = ProviderSourceReadinessInactiveSource
		record.Reason = "Plan source or provider source is not active."
	case record.SourceType == provider.TypeManual:
		record.State = ProviderSourceReadinessFakeProviderOnly
		record.Reason = "Manual source only works with the local fake provider path."
	case !supportsProductAutoProvision(record.ProductType, record.capabilityProfile):
		record.State = ProviderSourceReadinessUnsupportedCapability
		record.Reason = "Source does not support automatic provisioning for this product type."
	default:
		record.State = ProviderSourceReadinessReady
		record.Reason = "Source is active and supports automatic provisioning."
	}
	return record
}

func supportsProductAutoProvision(productType ProductType, profile provider.CapabilityProfile) bool {
	if !profile.SupportsAutoProvision {
		return false
	}
	switch productType {
	case ProductTypeVPS:
		return profile.VPS.SupportsOSTemplateSelection ||
			profile.VPS.SupportsCustomHostname ||
			profile.VPS.SupportsIPv6 ||
			profile.VPS.SupportsResize ||
			profile.VPS.SupportsVNCConsole
	case ProductTypeProxy:
		return profile.Proxy.SupportsHTTPProtocol ||
			profile.Proxy.SupportsSOCKS5Protocol ||
			profile.Proxy.SupportsRotatingProxy ||
			profile.Proxy.SupportsStaticProxy
	case ProductTypeServiceAddon:
		return true
	default:
		return false
	}
}
