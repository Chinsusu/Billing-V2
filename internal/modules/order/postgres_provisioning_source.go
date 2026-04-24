package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

const resolveOrderProvisioningSourceSQL = `
SELECT source.source_id, source.source_type
FROM tenant_plans tenant_plan
JOIN plan_sources plan_source ON plan_source.plan_id = tenant_plan.master_plan_id
JOIN provider_sources source ON source.source_id = plan_source.source_id
WHERE tenant_plan.tenant_plan_id = $1
  AND tenant_plan.tenant_id = $2
  AND plan_source.status = 'active'
  AND source.status = 'active'
ORDER BY plan_source.priority ASC, plan_source.created_at ASC
LIMIT 1`

func (store *PostgresStore) ResolveOrderProvisioningSource(
	ctx context.Context,
	input ResolveOrderProvisioningSourceInput,
) (ProvisioningSource, error) {
	if err := store.ready(); err != nil {
		return ProvisioningSource{}, err
	}
	args, err := resolveOrderProvisioningSourceArgs(input)
	if err != nil {
		return ProvisioningSource{}, err
	}
	var sourceID, providerType string
	if err := store.executor.QueryRowContext(ctx, resolveOrderProvisioningSourceSQL, args...).Scan(&sourceID, &providerType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProvisioningSource{}, ErrProvisioningSourceNotFound
		}
		return ProvisioningSource{}, fmt.Errorf("resolve order provisioning source: %w", err)
	}
	return ProvisioningSource{
		ProviderSourceID: catalog.ProviderSourceID(sourceID),
		ProviderType:     provider.Type(providerType),
	}, nil
}

func resolveOrderProvisioningSourceArgs(input ResolveOrderProvisioningSourceInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{input.TenantPlanID, input.TenantID}, nil
}
