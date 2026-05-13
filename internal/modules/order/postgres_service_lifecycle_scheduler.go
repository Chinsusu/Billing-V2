package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const listDueServiceLifecycleActionsSQL = `
SELECT
    service_instance_id,
    tenant_id,
    status,
    CASE
        WHEN status = 'active' THEN 'expire'
        WHEN status = 'expired' THEN 'grace'
        ELSE 'terminate'
    END AS action,
    CASE
        WHEN status = 'active' THEN 'expired'
        WHEN status = 'expired' THEN 'suspended'
        ELSE 'terminated'
    END AS to_status,
    'overdue' AS billing_status,
    CASE
        WHEN status = 'expired' THEN 'expiry'
        WHEN status = 'suspended' THEN suspension_reason::text
        ELSE NULL
    END AS suspension_reason,
    CASE
        WHEN status = 'active' THEN 'paid'
        ELSE 'overdue'
    END AS expected_billing_status,
    CASE
        WHEN status = 'suspended' THEN 'expiry'
        ELSE NULL
    END AS expected_suspension_reason,
    CASE
        WHEN status = 'active' THEN 'service term expired'
        WHEN status = 'expired' THEN 'service entered expiry grace'
        ELSE 'service expired beyond grace period'
    END AS reason,
    term_end
FROM service_instances
WHERE (
        status = 'active'
        AND billing_status = 'paid'
        AND term_end <= $1
    )
   OR (
        status = 'expired'
        AND billing_status = 'overdue'
        AND term_end <= $1
    )
   OR (
        status = 'suspended'
        AND billing_status = 'overdue'
        AND suspension_reason = 'expiry'
        AND term_end <= $2
    )
ORDER BY term_end ASC, created_at ASC
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
	var serviceID, tenantID, fromStatus, lifecycleAction, toStatus, billingStatus string
	var suspensionReason, expectedSuspensionReason sql.NullString
	var expectedBillingStatus string
	var termEnd time.Time
	if err := row.Scan(
		&serviceID,
		&tenantID,
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
