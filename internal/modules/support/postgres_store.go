package support

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

var ErrStoreExecutorMissing = errors.New("support store executor missing")

type PostgresStore struct {
	executor platformdb.Executor
}

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const supportTicketColumns = `support_ticket_id, display_id, tenant_id, requester_user_id, created_by, assigned_user_id, category, priority, status, subject, reference_type, reference_id, correlation_id, created_at, updated_at`

const supportTicketNoteColumns = `support_ticket_note_id, display_id, support_ticket_id, tenant_id, author_user_id, visibility, body_redacted, created_at`

const riskFlagColumns = `risk_flag_id, display_id, tenant_id, user_id, service_instance_id, order_id, flag_type, severity, status, note_redacted, created_by, correlation_id, created_at, updated_at`

const abuseCaseColumns = `abuse_case_id, display_id, tenant_id, user_id, service_instance_id, provider_source_id, case_type, severity, report_source, status, evidence_summary_redacted, deadline_at, assigned_owner_id, action_taken, final_resolution, created_by, correlation_id, created_at, updated_at`

const createSupportTicketSQL = `
INSERT INTO support_tickets (tenant_id, requester_user_id, created_by, category, priority, subject, reference_type, reference_id, correlation_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8::uuid, $9::uuid)
RETURNING ` + supportTicketColumns

func (store *PostgresStore) CreateSupportTicket(ctx context.Context, input CreateTicketInput) (Ticket, error) {
	if err := store.ready(); err != nil {
		return Ticket{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Ticket{}, err
	}
	return scanSupportTicket(store.executor.QueryRowContext(ctx, createSupportTicketSQL,
		input.TenantID,
		input.RequesterUserID,
		input.Actor.ID,
		input.Category,
		input.Priority,
		input.Subject,
		nullableString(string(input.ReferenceType)),
		nullableString(string(input.ReferenceID)),
		input.CorrelationID,
	))
}

const createSupportTicketNoteSQL = `
INSERT INTO support_ticket_notes (support_ticket_id, tenant_id, author_user_id, visibility, body_redacted)
VALUES ($1, $2, $3, $4, $5)
RETURNING ` + supportTicketNoteColumns

func (store *PostgresStore) CreateSupportTicketNote(ctx context.Context, input CreateTicketNoteInput) (TicketNote, error) {
	if err := store.ready(); err != nil {
		return TicketNote{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return TicketNote{}, err
	}
	return scanSupportTicketNote(store.executor.QueryRowContext(ctx, createSupportTicketNoteSQL,
		input.TicketID,
		input.TenantID,
		input.AuthorID,
		input.Visibility,
		input.BodyRedacted,
	))
}

const createRiskFlagSQL = `
INSERT INTO risk_flags (tenant_id, user_id, service_instance_id, order_id, flag_type, severity, note_redacted, created_by, correlation_id)
VALUES ($1, $2::uuid, $3::uuid, $4::uuid, $5, $6, $7, $8, $9::uuid)
RETURNING ` + riskFlagColumns

func (store *PostgresStore) CreateRiskFlag(ctx context.Context, input CreateRiskFlagInput) (RiskFlag, error) {
	if err := store.ready(); err != nil {
		return RiskFlag{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return RiskFlag{}, err
	}
	return scanRiskFlag(store.executor.QueryRowContext(ctx, createRiskFlagSQL,
		input.TenantID,
		nullableString(string(input.UserID)),
		nullableString(string(input.ServiceID)),
		nullableString(string(input.OrderID)),
		input.FlagType,
		input.Severity,
		nullableString(input.NoteRedacted),
		input.Actor.ID,
		input.CorrelationID,
	))
}

const createAbuseCaseSQL = `
INSERT INTO abuse_cases (tenant_id, user_id, service_instance_id, provider_source_id, case_type, severity, report_source, evidence_summary_redacted, deadline_at, assigned_owner_id, created_by, correlation_id)
VALUES ($1, $2::uuid, $3::uuid, $4::uuid, $5, $6, $7, $8, $9, $10::uuid, $11, $12::uuid)
RETURNING ` + abuseCaseColumns

func (store *PostgresStore) CreateAbuseCase(ctx context.Context, input CreateAbuseCaseInput) (AbuseCase, error) {
	if err := store.ready(); err != nil {
		return AbuseCase{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return AbuseCase{}, err
	}
	return scanAbuseCase(store.executor.QueryRowContext(ctx, createAbuseCaseSQL,
		input.TenantID,
		nullableString(string(input.UserID)),
		nullableString(string(input.ServiceID)),
		nullableString(string(input.ProviderSourceID)),
		input.CaseType,
		input.Severity,
		input.ReportSource,
		input.EvidenceSummaryRedacted,
		nullableTime(input.DeadlineAt),
		nullableString(string(input.AssignedOwnerID)),
		input.Actor.ID,
		input.CorrelationID,
	))
}

const getAbuseCaseSQL = `
SELECT ` + abuseCaseColumns + `
FROM abuse_cases
WHERE abuse_case_id = $1
  AND tenant_id = $2`

func (store *PostgresStore) GetAbuseCase(ctx context.Context, tenantID tenant.ID, id AbuseCaseID) (AbuseCase, error) {
	if err := store.ready(); err != nil {
		return AbuseCase{}, err
	}
	id = AbuseCaseID(trim(string(id)))
	tenantID = tenant.ID(trim(string(tenantID)))
	if id.Empty() {
		return AbuseCase{}, ErrAbuseCaseIDMissing
	}
	if tenantID.Empty() {
		return AbuseCase{}, tenant.ErrTenantIDMissing
	}
	return scanAbuseCase(store.executor.QueryRowContext(ctx, getAbuseCaseSQL, id, tenantID))
}

const markAbuseCaseSuspendedSQL = `
UPDATE abuse_cases
SET status = 'suspended',
    action_taken = $3,
    updated_at = NOW()
WHERE abuse_case_id = $1
  AND tenant_id = $2
RETURNING ` + abuseCaseColumns

func (store *PostgresStore) MarkAbuseCaseSuspended(ctx context.Context, input MarkAbuseCaseSuspendedInput) (AbuseCase, error) {
	if err := store.ready(); err != nil {
		return AbuseCase{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return AbuseCase{}, err
	}
	return scanAbuseCase(store.executor.QueryRowContext(ctx, markAbuseCaseSuspendedSQL,
		input.ID,
		input.TenantID,
		input.ActionTaken,
	))
}

func (store *PostgresStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrStoreExecutorMissing
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanSupportTicket(row rowScanner) (Ticket, error) {
	var record Ticket
	var assignedUserID, referenceType, referenceID sql.NullString
	err := row.Scan(
		&record.ID,
		&record.DisplayID,
		&record.TenantID,
		&record.RequesterUserID,
		&record.CreatedBy,
		&assignedUserID,
		&record.Category,
		&record.Priority,
		&record.Status,
		&record.Subject,
		&referenceType,
		&referenceID,
		&record.CorrelationID,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Ticket{}, ErrSupportTicketNotFound
	}
	if err != nil {
		return Ticket{}, fmt.Errorf("scan support ticket: %w", err)
	}
	record.AssignedUserID = identity.UserID(assignedUserID.String)
	record.ReferenceType = ReferenceType(referenceType.String)
	record.ReferenceID = ReferenceID(referenceID.String)
	return record, nil
}

func scanSupportTicketNote(row rowScanner) (TicketNote, error) {
	var record TicketNote
	err := row.Scan(
		&record.ID,
		&record.DisplayID,
		&record.TicketID,
		&record.TenantID,
		&record.AuthorID,
		&record.Visibility,
		&record.BodyRedacted,
		&record.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return TicketNote{}, ErrSupportTicketNotFound
	}
	if err != nil {
		return TicketNote{}, fmt.Errorf("scan support ticket note: %w", err)
	}
	return record, nil
}

func scanRiskFlag(row rowScanner) (RiskFlag, error) {
	var record RiskFlag
	var userID, serviceID, orderID, note sql.NullString
	err := row.Scan(
		&record.ID,
		&record.DisplayID,
		&record.TenantID,
		&userID,
		&serviceID,
		&orderID,
		&record.FlagType,
		&record.Severity,
		&record.Status,
		&note,
		&record.CreatedBy,
		&record.CorrelationID,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return RiskFlag{}, ErrRiskFlagTargetMissing
	}
	if err != nil {
		return RiskFlag{}, fmt.Errorf("scan risk flag: %w", err)
	}
	record.UserID = identity.UserID(userID.String)
	record.ServiceID = order.ServiceID(serviceID.String)
	record.OrderID = order.OrderID(orderID.String)
	record.NoteRedacted = note.String
	return record, nil
}

func scanAbuseCase(row rowScanner) (AbuseCase, error) {
	var record AbuseCase
	var userID, serviceID, providerSourceID, assignedOwnerID, actionTaken, finalResolution sql.NullString
	var deadlineAt sql.NullTime
	err := row.Scan(
		&record.ID,
		&record.DisplayID,
		&record.TenantID,
		&userID,
		&serviceID,
		&providerSourceID,
		&record.CaseType,
		&record.Severity,
		&record.ReportSource,
		&record.Status,
		&record.EvidenceSummaryRedacted,
		&deadlineAt,
		&assignedOwnerID,
		&actionTaken,
		&finalResolution,
		&record.CreatedBy,
		&record.CorrelationID,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return AbuseCase{}, ErrAbuseCaseNotFound
	}
	if err != nil {
		return AbuseCase{}, fmt.Errorf("scan abuse case: %w", err)
	}
	record.UserID = identity.UserID(userID.String)
	record.ServiceID = order.ServiceID(serviceID.String)
	record.ProviderSourceID = catalog.ProviderSourceID(providerSourceID.String)
	record.DeadlineAt = deadlineAt.Time
	record.AssignedOwnerID = identity.UserID(assignedOwnerID.String)
	record.ActionTaken = actionTaken.String
	record.FinalResolution = finalResolution.String
	return record, nil
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func nullableTime(value time.Time) interface{} {
	if value.IsZero() {
		return nil
	}
	return value
}
