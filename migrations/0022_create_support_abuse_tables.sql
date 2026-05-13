CREATE TYPE support_ticket_category AS ENUM (
    'billing',
    'topup',
    'order',
    'provisioning',
    'service_access',
    'credential',
    'renewal_expiry',
    'suspension_termination',
    'provider_issue',
    'abuse_takedown',
    'account_login',
    'reseller_setup',
    'feature_request',
    'other'
);

CREATE TYPE support_ticket_priority AS ENUM (
    'p0',
    'p1',
    'p2',
    'p3',
    'p4'
);

CREATE TYPE support_ticket_status AS ENUM (
    'open',
    'waiting_on_customer',
    'waiting_on_support',
    'resolved',
    'closed'
);

CREATE TYPE support_note_visibility AS ENUM (
    'public',
    'internal'
);

CREATE TYPE risk_flag_type AS ENUM (
    'new_account_high_value',
    'payment_mismatch',
    'abuse_history',
    'manual_blacklist',
    'provider_risk'
);

CREATE TYPE risk_flag_status AS ENUM (
    'open',
    'reviewing',
    'cleared',
    'confirmed'
);

CREATE TYPE abuse_case_type AS ENUM (
    'spam',
    'phishing',
    'malware',
    'botnet',
    'brute_force',
    'port_scanning',
    'ddos',
    'copyright',
    'credential_theft',
    'proxy_scraping_violation',
    'payment_fraud',
    'chargeback',
    'illegal_content',
    'aup_violation',
    'provider_takedown',
    'other'
);

CREATE TYPE abuse_case_severity AS ENUM (
    'low',
    'medium',
    'high',
    'critical'
);

CREATE TYPE abuse_case_status AS ENUM (
    'new',
    'triaging',
    'awaiting_client_action',
    'suspended',
    'resolved',
    'rejected_false_positive',
    'terminated',
    'escalated',
    'closed'
);

CREATE TYPE abuse_report_source AS ENUM (
    'provider',
    'datacenter',
    'email_abuse_desk',
    'legal',
    'payment_processor',
    'internal_monitoring',
    'client',
    'reseller',
    'third_party',
    'other'
);

CREATE SEQUENCE support_tickets_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE support_ticket_notes_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE risk_flags_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE abuse_cases_display_id_seq AS BIGINT START WITH 10000;

ALTER TABLE orders
    ADD CONSTRAINT orders_id_tenant_unique UNIQUE (order_id, tenant_id);

CREATE TABLE support_tickets (
    support_ticket_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('support_tickets_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    requester_user_id UUID NOT NULL REFERENCES users(user_id),
    created_by UUID NOT NULL REFERENCES users(user_id),
    assigned_user_id UUID REFERENCES users(user_id),
    category support_ticket_category NOT NULL,
    priority support_ticket_priority NOT NULL,
    status support_ticket_status NOT NULL DEFAULT 'open',
    subject TEXT NOT NULL,
    reference_type TEXT,
    reference_id UUID,
    correlation_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT support_tickets_display_id_unique UNIQUE (display_id),
    CONSTRAINT support_tickets_subject_not_blank CHECK (btrim(subject) <> '')
);

ALTER SEQUENCE support_tickets_display_id_seq OWNED BY support_tickets.display_id;

CREATE INDEX idx_support_tickets_tenant_status_created ON support_tickets(tenant_id, status, created_at);
CREATE INDEX idx_support_tickets_requester_created ON support_tickets(requester_user_id, created_at);
CREATE INDEX idx_support_tickets_reference ON support_tickets(reference_type, reference_id);
CREATE INDEX idx_support_tickets_correlation_id ON support_tickets(correlation_id);

CREATE TABLE support_ticket_notes (
    support_ticket_note_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('support_ticket_notes_display_id_seq'),
    support_ticket_id UUID NOT NULL REFERENCES support_tickets(support_ticket_id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    author_user_id UUID NOT NULL REFERENCES users(user_id),
    visibility support_note_visibility NOT NULL,
    body_redacted TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT support_ticket_notes_display_id_unique UNIQUE (display_id),
    CONSTRAINT support_ticket_notes_body_not_blank CHECK (btrim(body_redacted) <> '')
);

ALTER SEQUENCE support_ticket_notes_display_id_seq OWNED BY support_ticket_notes.display_id;

CREATE INDEX idx_support_ticket_notes_ticket_created ON support_ticket_notes(support_ticket_id, created_at);
CREATE INDEX idx_support_ticket_notes_tenant_created ON support_ticket_notes(tenant_id, created_at);

CREATE TABLE risk_flags (
    risk_flag_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('risk_flags_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    user_id UUID REFERENCES users(user_id),
    service_instance_id UUID,
    order_id UUID,
    flag_type risk_flag_type NOT NULL,
    severity abuse_case_severity NOT NULL,
    status risk_flag_status NOT NULL DEFAULT 'open',
    note_redacted TEXT,
    created_by UUID NOT NULL REFERENCES users(user_id),
    correlation_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT risk_flags_display_id_unique UNIQUE (display_id),
    CONSTRAINT risk_flags_target_required CHECK (
        user_id IS NOT NULL OR service_instance_id IS NOT NULL OR order_id IS NOT NULL
    ),
    CONSTRAINT risk_flags_service_tenant_fk FOREIGN KEY (service_instance_id, tenant_id)
        REFERENCES service_instances(service_instance_id, tenant_id),
    CONSTRAINT risk_flags_order_tenant_fk FOREIGN KEY (order_id, tenant_id)
        REFERENCES orders(order_id, tenant_id)
);

ALTER SEQUENCE risk_flags_display_id_seq OWNED BY risk_flags.display_id;

CREATE INDEX idx_risk_flags_tenant_status_created ON risk_flags(tenant_id, status, created_at);
CREATE INDEX idx_risk_flags_user_created ON risk_flags(user_id, created_at);
CREATE INDEX idx_risk_flags_service_created ON risk_flags(service_instance_id, created_at);
CREATE INDEX idx_risk_flags_order_created ON risk_flags(order_id, created_at);
CREATE INDEX idx_risk_flags_correlation_id ON risk_flags(correlation_id);

CREATE TABLE abuse_cases (
    abuse_case_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('abuse_cases_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    user_id UUID REFERENCES users(user_id),
    service_instance_id UUID,
    provider_source_id UUID REFERENCES provider_sources(source_id),
    case_type abuse_case_type NOT NULL,
    severity abuse_case_severity NOT NULL,
    report_source abuse_report_source NOT NULL,
    status abuse_case_status NOT NULL DEFAULT 'new',
    evidence_summary_redacted TEXT NOT NULL,
    deadline_at TIMESTAMPTZ,
    assigned_owner_id UUID REFERENCES users(user_id),
    action_taken TEXT,
    final_resolution TEXT,
    created_by UUID NOT NULL REFERENCES users(user_id),
    correlation_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT abuse_cases_display_id_unique UNIQUE (display_id),
    CONSTRAINT abuse_cases_evidence_not_blank CHECK (btrim(evidence_summary_redacted) <> ''),
    CONSTRAINT abuse_cases_service_tenant_fk FOREIGN KEY (service_instance_id, tenant_id)
        REFERENCES service_instances(service_instance_id, tenant_id)
);

ALTER SEQUENCE abuse_cases_display_id_seq OWNED BY abuse_cases.display_id;

CREATE INDEX idx_abuse_cases_tenant_status_created ON abuse_cases(tenant_id, status, created_at);
CREATE INDEX idx_abuse_cases_service_created ON abuse_cases(service_instance_id, created_at);
CREATE INDEX idx_abuse_cases_user_created ON abuse_cases(user_id, created_at);
CREATE INDEX idx_abuse_cases_provider_source_created ON abuse_cases(provider_source_id, created_at);
CREATE INDEX idx_abuse_cases_correlation_id ON abuse_cases(correlation_id);
