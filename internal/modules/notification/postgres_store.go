package notification

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

var ErrStoreExecutorMissing = errors.New("notification store executor missing")

type PostgresStore struct {
	executor platformdb.Executor
}

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const notificationColumns = `notification_id, display_id, tenant_id, recipient_user_id, recipient_group, channel, template_key, event_type, priority, payload_redacted, reference_type, reference_id, dedupe_key, status, last_error_code, last_error_message_redacted, correlation_id, sent_at, created_at, updated_at`

const createNotificationSQL = `
INSERT INTO notifications (tenant_id, recipient_user_id, recipient_group, channel, template_key, event_type, priority, payload_redacted, reference_type, reference_id, dedupe_key, correlation_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb, $9, $10, $11, $12)
ON CONFLICT (tenant_id, channel, dedupe_key)
DO UPDATE SET updated_at = notifications.updated_at
RETURNING ` + notificationColumns

func (store *PostgresStore) CreateNotification(ctx context.Context, input CreateInput) (Notification, error) {
	if err := store.ready(); err != nil {
		return Notification{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Notification{}, err
	}
	return scanNotification(store.executor.QueryRowContext(ctx, createNotificationSQL,
		input.TenantID,
		nullableString(string(input.RecipientUserID)),
		nullableString(string(input.RecipientGroup)),
		input.Channel,
		input.TemplateKey,
		input.EventType,
		input.Priority,
		string(input.PayloadRedacted),
		nullableString(string(input.ReferenceType)),
		nullableString(string(input.ReferenceID)),
		input.DedupeKey,
		input.CorrelationID,
	))
}

func (store *PostgresStore) MarkNotificationSent(ctx context.Context, id ID, sentAt time.Time) (Notification, error) {
	if err := store.ready(); err != nil {
		return Notification{}, err
	}
	if id == "" {
		return Notification{}, ErrNotificationIDMissing
	}
	if sentAt.IsZero() {
		sentAt = time.Now().UTC()
	}
	return scanNotification(store.executor.QueryRowContext(ctx, `
UPDATE notifications
SET status = 'sent',
    sent_at = $2,
    last_error_code = NULL,
    last_error_message_redacted = NULL,
    updated_at = NOW()
WHERE notification_id = $1
RETURNING `+notificationColumns, id, sentAt))
}

func (store *PostgresStore) MarkNotificationFailed(ctx context.Context, id ID, failedAt time.Time, errorCode string, errorMessageRedacted string) (Notification, error) {
	if err := store.ready(); err != nil {
		return Notification{}, err
	}
	if id == "" {
		return Notification{}, ErrNotificationIDMissing
	}
	return scanNotification(store.executor.QueryRowContext(ctx, `
UPDATE notifications
SET status = 'failed',
    last_error_code = $2,
    last_error_message_redacted = $3,
    updated_at = COALESCE($4, NOW())
WHERE notification_id = $1
RETURNING `+notificationColumns, id, nullableString(errorCode), nullableString(errorMessageRedacted), nullableTime(failedAt)))
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

func scanNotification(row rowScanner) (Notification, error) {
	var record Notification
	var recipientUserID, recipientGroup, referenceType, referenceID, lastErrorCode, lastErrorMessage sql.NullString
	var sentAt sql.NullTime
	err := row.Scan(
		&record.ID,
		&record.DisplayID,
		&record.TenantID,
		&recipientUserID,
		&recipientGroup,
		&record.Channel,
		&record.TemplateKey,
		&record.EventType,
		&record.Priority,
		&record.PayloadRedacted,
		&referenceType,
		&referenceID,
		&record.DedupeKey,
		&record.Status,
		&lastErrorCode,
		&lastErrorMessage,
		&record.CorrelationID,
		&sentAt,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Notification{}, ErrNotificationNotFound
	}
	if err != nil {
		return Notification{}, fmt.Errorf("scan notification: %w", err)
	}
	record.RecipientUserID = identity.UserID(recipientUserID.String)
	record.RecipientGroup = RecipientGroup(recipientGroup.String)
	record.ReferenceType = ReferenceType(referenceType.String)
	record.ReferenceID = ReferenceID(referenceID.String)
	record.LastErrorCode = lastErrorCode.String
	record.LastErrorMessageRedacted = lastErrorMessage.String
	record.SentAt = sentAt.Time
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
