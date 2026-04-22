CREATE TYPE outbox_event_status AS ENUM (
    'pending',
    'processing',
    'published',
    'failed_retryable',
    'failed_terminal',
    'discarded'
);

CREATE TYPE job_status AS ENUM (
    'queued',
    'claimed',
    'running',
    'succeeded',
    'failed_retryable',
    'failed_terminal',
    'manual_review',
    'cancelled'
);

CREATE TYPE job_attempt_result AS ENUM (
    'succeeded',
    'failed_retryable',
    'failed_terminal',
    'manual_review',
    'cancelled'
);

CREATE TABLE outbox_events (
    outbox_event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id),
    aggregate_type TEXT NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type TEXT NOT NULL,
    payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    status outbox_event_status NOT NULL DEFAULT 'pending',
    dedupe_key TEXT NOT NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 10,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    locked_by TEXT,
    locked_until TIMESTAMPTZ,
    last_error_code TEXT,
    last_error_message_redacted TEXT,
    correlation_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    CONSTRAINT outbox_events_attempt_count_non_negative CHECK (attempt_count >= 0),
    CONSTRAINT outbox_events_max_attempts_positive CHECK (max_attempts > 0),
    CONSTRAINT outbox_events_unique_dedupe_key UNIQUE (dedupe_key)
);

CREATE INDEX idx_outbox_events_claim ON outbox_events(status, next_attempt_at, created_at);
CREATE INDEX idx_outbox_events_tenant_event_created ON outbox_events(tenant_id, event_type, created_at);
CREATE INDEX idx_outbox_events_correlation_id ON outbox_events(correlation_id);
CREATE INDEX idx_outbox_events_locked_until ON outbox_events(locked_until) WHERE locked_until IS NOT NULL;

CREATE TABLE jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id),
    job_type TEXT NOT NULL,
    reference_type TEXT NOT NULL,
    reference_id UUID NOT NULL,
    source_id UUID,
    payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    status job_status NOT NULL DEFAULT 'queued',
    priority INT NOT NULL DEFAULT 100,
    idempotency_key TEXT NOT NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    locked_by TEXT,
    locked_until TIMESTAMPTZ,
    last_error_code TEXT,
    last_error_message_redacted TEXT,
    manual_review_reason TEXT,
    correlation_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    CONSTRAINT jobs_attempt_count_non_negative CHECK (attempt_count >= 0),
    CONSTRAINT jobs_max_attempts_positive CHECK (max_attempts > 0),
    CONSTRAINT jobs_priority_non_negative CHECK (priority >= 0)
);

CREATE UNIQUE INDEX idx_jobs_tenant_type_idempotency ON jobs(tenant_id, job_type, idempotency_key) WHERE tenant_id IS NOT NULL;
CREATE UNIQUE INDEX idx_jobs_global_type_idempotency ON jobs(job_type, idempotency_key) WHERE tenant_id IS NULL;
CREATE INDEX idx_jobs_claim ON jobs(status, next_attempt_at, priority, created_at);
CREATE INDEX idx_jobs_tenant_type_created ON jobs(tenant_id, job_type, created_at);
CREATE INDEX idx_jobs_reference ON jobs(reference_type, reference_id);
CREATE INDEX idx_jobs_correlation_id ON jobs(correlation_id);
CREATE INDEX idx_jobs_locked_until ON jobs(locked_until) WHERE locked_until IS NOT NULL;

CREATE TABLE job_attempts (
    job_attempt_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(job_id),
    worker_id TEXT NOT NULL,
    attempt_number INT NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ,
    result job_attempt_result NOT NULL,
    error_code TEXT,
    error_message_redacted TEXT,
    duration_ms INT,
    correlation_id UUID NOT NULL,
    CONSTRAINT job_attempts_attempt_number_positive CHECK (attempt_number > 0),
    CONSTRAINT job_attempts_duration_non_negative CHECK (duration_ms IS NULL OR duration_ms >= 0),
    CONSTRAINT job_attempts_unique_attempt_number UNIQUE (job_id, attempt_number)
);
