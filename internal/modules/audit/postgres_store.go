package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

var ErrAuditStoreExecutorMissing = errors.New("audit store executor missing")

type PostgresStore struct {
	executor platformdb.Executor
}

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const auditColumns = `audit_id, display_id, tenant_id, actor_id, actor_type, action, target_type, target_id, before_snapshot_redacted, after_snapshot_redacted, metadata_redacted, ip_address, user_agent, correlation_id, created_at`
const auditReadColumns = auditColumns + `,
(SELECT actor.display_id FROM users actor WHERE actor.user_id = audit_logs.actor_id AND actor.tenant_id = audit_logs.tenant_id) AS actor_display_id,
CASE
  WHEN audit_logs.target_type = 'invoice' THEN (
    SELECT inv.display_id FROM invoices inv WHERE inv.invoice_id = audit_logs.target_id AND inv.tenant_id = audit_logs.tenant_id
  )
  WHEN audit_logs.target_type = 'order' THEN (
    SELECT ord.display_id FROM orders ord WHERE ord.order_id = audit_logs.target_id AND ord.tenant_id = audit_logs.tenant_id
  )
  WHEN audit_logs.target_type = 'job' THEN (
    SELECT job.display_id FROM jobs job WHERE job.job_id = audit_logs.target_id AND job.tenant_id = audit_logs.tenant_id
  )
  WHEN audit_logs.target_type = 'topup_request' THEN (
    SELECT topup.display_id FROM topup_requests topup WHERE topup.topup_request_id = audit_logs.target_id AND topup.tenant_id = audit_logs.tenant_id
  )
  WHEN audit_logs.target_type IN ('service', 'service_instance') THEN (
    SELECT svc.display_id FROM service_instances svc WHERE svc.service_instance_id = audit_logs.target_id AND svc.tenant_id = audit_logs.tenant_id
  )
  WHEN audit_logs.target_type IN ('provider_source', 'provider') THEN (
    SELECT source.display_id FROM provider_sources source WHERE source.source_id = audit_logs.target_id
  )
  ELSE NULL
END AS target_display_id`

func (store *PostgresStore) Append(ctx context.Context, input AppendInput) (Log, error) {
	if err := store.ready(); err != nil {
		return Log{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Log{}, err
	}

	row := store.executor.QueryRowContext(ctx, `
INSERT INTO audit_logs (tenant_id, actor_id, actor_type, action, target_type, target_id, before_snapshot_redacted, after_snapshot_redacted, metadata_redacted, ip_address, user_agent, correlation_id)
VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8::jsonb, $9::jsonb, $10, $11, $12)
RETURNING `+auditColumns,
		nullableTenantID(input.TenantID), nullableString(string(input.ActorID)), input.ActorType, input.Action, input.TargetType,
		input.TargetID, nullableJSON(input.BeforeSnapshotRedacted), nullableJSON(input.AfterSnapshotRedacted), string(input.MetadataRedacted),
		nullableString(input.IPAddress), nullableString(input.UserAgent), input.CorrelationID)
	return scanLog(row)
}

func (store *PostgresStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrAuditStoreExecutorMissing
	}
	return nil
}

type logScanner interface {
	Scan(dest ...interface{}) error
}

func scanLog(row logScanner) (Log, error) {
	return scanLogFields(row, false)
}

func scanLogRead(row logScanner) (Log, error) {
	return scanLogFields(row, true)
}

func scanLogFields(row logScanner, includeRelatedDisplayIDs bool) (Log, error) {
	var record Log
	var id, tenantID, actorID, actorType, targetID, correlationID string
	var tenantNull, actorNull, ipNull, userAgentNull sql.NullString
	var actorDisplayID, targetDisplayID sql.NullInt64
	var beforeSnapshot, afterSnapshot, metadata []byte

	destinations := []interface{}{
		&id, &record.DisplayID, &tenantNull, &actorNull, &actorType, &record.Action, &record.TargetType, &targetID,
		&beforeSnapshot, &afterSnapshot, &metadata, &ipNull, &userAgentNull, &correlationID, &record.CreatedAt,
	}
	if includeRelatedDisplayIDs {
		destinations = append(destinations, &actorDisplayID, &targetDisplayID)
	}
	if err := row.Scan(destinations...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Log{}, ErrAuditLogNotFound
		}
		return Log{}, fmt.Errorf("scan audit log: %w", err)
	}
	record.ID = ID(id)
	tenantID = tenantNull.String
	record.TenantID = tenant.ID(tenantID)
	actorID = actorNull.String
	record.ActorID = ActorID(actorID)
	if actorDisplayID.Valid {
		record.ActorDisplayID = actorDisplayID.Int64
	}
	record.ActorType = ActorType(actorType)
	record.TargetID = TargetID(targetID)
	if targetDisplayID.Valid {
		record.TargetDisplayID = targetDisplayID.Int64
	}
	record.BeforeSnapshotRedacted = append(record.BeforeSnapshotRedacted, beforeSnapshot...)
	record.AfterSnapshotRedacted = append(record.AfterSnapshotRedacted, afterSnapshot...)
	record.MetadataRedacted = append(record.MetadataRedacted, metadata...)
	record.IPAddress = ipNull.String
	record.UserAgent = userAgentNull.String
	record.CorrelationID = CorrelationID(correlationID)
	return record, nil
}

func nullableTenantID(id tenant.ID) sql.NullString {
	return nullableString(string(id))
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func nullableJSON(value json.RawMessage) interface{} {
	if len(value) == 0 {
		return nil
	}
	return string(value)
}
