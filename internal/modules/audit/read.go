package audit

import (
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const defaultLogListLimit = 100
const maxLogListLimit = 500

func normalizeFilter(filter Filter) Filter {
	output := filter
	output.Action = stringsTrim(output.Action)
	output.TargetType = stringsTrim(output.TargetType)
	if output.Limit <= 0 {
		output.Limit = defaultLogListLimit
	}
	if output.Limit > maxLogListLimit {
		output.Limit = maxLogListLimit
	}
	return output
}

func validateFilter(filter Filter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.ActorType != "" && !filter.ActorType.Valid() {
		return ErrActorTypeInvalid
	}
	if !filter.CreatedFrom.IsZero() && !filter.CreatedTo.IsZero() && filter.CreatedTo.Before(filter.CreatedFrom) {
		return ErrCreatedWindowInvalid
	}
	return nil
}

func validateLookup(lookup Lookup) error {
	if lookup.ID.Empty() {
		return ErrAuditLogNotFound
	}
	if lookup.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	return nil
}

func stringsTrim(value string) string {
	return strings.TrimSpace(value)
}
