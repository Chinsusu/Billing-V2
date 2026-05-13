CREATE TYPE notification_channel AS ENUM (
    'email',
    'dashboard',
    'telegram',
    'webhook'
);

CREATE TYPE notification_status AS ENUM (
    'queued',
    'sent',
    'failed',
    'cancelled'
);

CREATE TYPE notification_priority AS ENUM (
    'low',
    'normal',
    'high',
    'critical'
);

CREATE SEQUENCE notifications_display_id_seq AS BIGINT START WITH 10000;

CREATE TABLE notifications (
    notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('notifications_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    recipient_user_id UUID REFERENCES users(user_id),
    recipient_group TEXT,
    channel notification_channel NOT NULL,
    template_key TEXT NOT NULL,
    event_type TEXT NOT NULL,
    priority notification_priority NOT NULL DEFAULT 'normal',
    payload_redacted JSONB NOT NULL DEFAULT '{}'::jsonb,
    reference_type TEXT,
    reference_id UUID,
    dedupe_key TEXT NOT NULL,
    status notification_status NOT NULL DEFAULT 'queued',
    last_error_code TEXT,
    last_error_message_redacted TEXT,
    correlation_id UUID NOT NULL,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT notifications_display_id_unique UNIQUE (display_id),
    CONSTRAINT notifications_recipient_present CHECK (recipient_user_id IS NOT NULL OR btrim(COALESCE(recipient_group, '')) <> ''),
    CONSTRAINT notifications_template_key_not_blank CHECK (btrim(template_key) <> ''),
    CONSTRAINT notifications_event_type_not_blank CHECK (btrim(event_type) <> ''),
    CONSTRAINT notifications_dedupe_key_not_blank CHECK (btrim(dedupe_key) <> ''),
    CONSTRAINT notifications_tenant_channel_dedupe_unique UNIQUE (tenant_id, channel, dedupe_key)
);

ALTER SEQUENCE notifications_display_id_seq OWNED BY notifications.display_id;

CREATE INDEX idx_notifications_tenant_status_created ON notifications(tenant_id, status, created_at);
CREATE INDEX idx_notifications_recipient_created ON notifications(recipient_user_id, created_at);
CREATE INDEX idx_notifications_event_created ON notifications(event_type, created_at);
CREATE INDEX idx_notifications_reference ON notifications(reference_type, reference_id);
CREATE INDEX idx_notifications_correlation_id ON notifications(correlation_id);
