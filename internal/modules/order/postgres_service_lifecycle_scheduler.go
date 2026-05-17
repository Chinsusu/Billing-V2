package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const listDueServiceLifecycleActionsSQL = `
SELECT
    svc.service_instance_id,
    svc.tenant_id,
    svc.provider_source_id,
    source.source_type,
    svc.external_resource_id,
    svc.status,
    CASE
        WHEN svc.status = 'active' THEN 'expire'
        WHEN svc.status = 'expired' THEN 'grace'
        ELSE 'terminate'
    END AS action,
    CASE
        WHEN svc.status = 'active' THEN 'expired'
        WHEN svc.status = 'expired' THEN 'suspended'
        ELSE 'terminated'
    END AS to_status,
    'overdue' AS billing_status,
    CASE
        WHEN svc.status = 'expired' THEN 'expiry'
        WHEN svc.status = 'suspended' THEN svc.suspension_reason::text
        ELSE NULL
    END AS suspension_reason,
    CASE
        WHEN svc.status = 'active' THEN 'paid'
        ELSE 'overdue'
    END AS expected_billing_status,
    CASE
        WHEN svc.status = 'suspended' THEN 'expiry'
        ELSE NULL
    END AS expected_suspension_reason,
    CASE
        WHEN svc.status = 'active' THEN 'service term expired'
        WHEN svc.status = 'expired' THEN 'service entered expiry grace'
        ELSE 'service expired beyond grace period'
    END AS reason,
    svc.term_end
FROM service_instances svc
JOIN provider_sources source
  ON source.source_id = svc.provider_source_id
WHERE (
        svc.status = 'active'
        AND svc.billing_status = 'paid'
        AND svc.term_end <= $1
    )
   OR (
        svc.status = 'expired'
        AND svc.billing_status = 'overdue'
        AND svc.term_end <= $1
    )
   OR (
        svc.status = 'suspended'
        AND svc.billing_status = 'overdue'
        AND svc.suspension_reason = 'expiry'
        AND svc.term_end <= $2
    )
ORDER BY svc.term_end ASC, svc.created_at ASC
LIMIT $3`

func (store *PostgresStore) ListDueServiceLifecycleActions(ctx context.Context, input ListDueServiceLifecycleActionsInput) ([]ServiceLifecycleDueAction, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	args, err := listDueServiceLifecycleActionsArgs(input)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, listDueServiceLifecycleActionsSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("list due service lifecycle actions: %w", err)
	}
	defer rows.Close()
	return scanDueServiceLifecycleActions(rows)
}

func listDueServiceLifecycleActionsArgs(input ListDueServiceLifecycleActionsInput) ([]interface{}, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	graceCutoff := input.Now.Add(-input.GracePeriod)
	return []interface{}{input.Now, graceCutoff, input.Limit}, nil
}

func scanDueServiceLifecycleActions(rows *sql.Rows) ([]ServiceLifecycleDueAction, error) {
	actions := make([]ServiceLifecycleDueAction, 0)
	for rows.Next() {
		action, err := scanDueServiceLifecycleAction(rows)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read due service lifecycle actions: %w", err)
	}
	return actions, nil
}

func scanDueServiceLifecycleAction(row interface {
	Scan(dest ...interface{}) error
}) (ServiceLifecycleDueAction, error) {
	var action ServiceLifecycleDueAction
	var serviceID, tenantID, providerSourceID, providerType, externalResourceID, fromStatus, lifecycleAction, toStatus, billingStatus string
	var suspensionReason, expectedSuspensionReason sql.NullString
	var expectedBillingStatus string
	var termEnd time.Time
	if err := row.Scan(
		&serviceID,
		&tenantID,
		&providerSourceID,
		&providerType,
		&externalResourceID,
		&fromStatus,
		&lifecycleAction,
		&toStatus,
		&billingStatus,
		&suspensionReason,
		&expectedBillingStatus,
		&expectedSuspensionReason,
		&action.Reason,
		&termEnd,
	); err != nil {
		return ServiceLifecycleDueAction{}, fmt.Errorf("scan due service lifecycle action: %w", err)
	}
	action.ServiceID = ServiceID(serviceID)
	action.TenantID = tenant.ID(tenantID)
	action.ProviderSourceID = catalog.ProviderSourceID(providerSourceID)
	action.ProviderType = provider.Type(providerType)
	action.ExternalResourceID = provider.ExternalResourceID(externalResourceID)
	action.FromStatus = ServiceStatus(fromStatus)
	action.Action = ServiceLifecycleAction(lifecycleAction)
	action.ToStatus = ServiceStatus(toStatus)
	action.BillingStatus = BillingStatus(billingStatus)
	action.SuspensionReason = SuspensionReason(suspensionReason.String)
	action.ExpectedBillingStatus = BillingStatus(expectedBillingStatus)
	action.ExpectedSuspensionReason = SuspensionReason(expectedSuspensionReason.String)
	action.TermEnd = termEnd
	return action.Normalize(), nil
}
