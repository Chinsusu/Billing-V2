# 73 - Cloudmini Idempotency Evidence Runbook

**Date:** 2026-05-18
**Scope:** Operator instructions for the T248 Cloudmini duplicate-create and timeout-after-send evidence smoke.
**Decision:** Tooling only until both scenarios are run against an approved non-production provider account and redacted evidence is recorded.

## Boundary

Run only after T226 mutating preflight still passes and Cloudmini runtime env is loaded from a protected path outside git.

The command refuses production, requires owner fields, and writes raw provider cleanup references only to `CLOUDMINI_IDEMPOTENCY_RAW_EVIDENCE_PATH`, which must be an absolute path outside the repo. Standard output is redacted and intentionally excludes raw DSNs, tokens, provider IDs, provider payloads, and proxy credentials.

## Required Env

Both scenarios require these owner and approval fields:

```text
APP_ENV=dev
BILLING_CLOUDMINI_IDEMPOTENCY_EVIDENCE_APPROVED=yes
CLOUDMINI_SOURCE_ACCOUNT_OWNER=Admin
CLOUDMINI_ENGINEERING_OWNER=Admin
CLOUDMINI_OPS_OWNER=Admin
CLOUDMINI_SECURITY_OWNER=Admin
CLOUDMINI_CLEANUP_OWNER=Admin
CLOUDMINI_FINANCE_QUOTA_OWNER=Admin
CLOUDMINI_REVIEWER_SIGNOFF=Admin
CLOUDMINI_PILOT_CLEANUP_DEADLINE=<same-session-deadline>
CLOUDMINI_PILOT_STOP_CONDITION=<approved-stop-condition>
CLOUDMINI_PILOT_READONLY_EVIDENCE_REF=<redacted-readonly-evidence-ref>
CLOUDMINI_PILOT_CLEANUP_PROCEDURE_REF=<redacted-cleanup-procedure-ref>
```

Both scenarios require Cloudmini runtime env from the protected credential path:

```text
CLOUDMINI_V3_BASE_URL
CLOUDMINI_V3_API_TOKEN
CLOUDMINI_V3_SOURCE_ID
CLOUDMINI_V3_KIND
CLOUDMINI_V3_GROUP_ID
CLOUDMINI_V3_PROTOCOL
ENCRYPTION_KEY
```

## Duplicate Create

Use this scenario to prove two create attempts with the same idempotency key do not create two provider resources.

```bash
APP_ENV=dev \
BILLING_CLOUDMINI_IDEMPOTENCY_EVIDENCE_APPROVED=yes \
CLOUDMINI_IDEMPOTENCY_SCENARIO=duplicate-create \
CLOUDMINI_IDEMPOTENCY_PILOT_ID=<approved-pilot-id> \
CLOUDMINI_IDEMPOTENCY_MAX_CREATE_ATTEMPTS=2 \
CLOUDMINI_IDEMPOTENCY_MAX_ACTIVE_RESOURCES=1 \
CLOUDMINI_IDEMPOTENCY_PROVIDER_RATE_LIMIT=no-parallel-mutating-calls \
CLOUDMINI_IDEMPOTENCY_MAX_SPEND_EXPOSURE=single-dev-resource \
CLOUDMINI_IDEMPOTENCY_RAW_EVIDENCE_PATH=/tmp/<approved-private-raw-evidence>.json \
go run ./cmd/smoke cloudmini-idempotency-evidence
```

Expected redacted pass evidence:

- `mutating_routes_called=yes`
- `create_attempts=2`
- `distinct_resource_count=1`
- `duplicate_same_resource=true`
- `cleanup_attempts=1`
- `sensitive_values_printed=no`
- `raw_provider_ids_printed=no`

## Timeout After Send

Use this scenario to prove Billing maps an accepted create whose status wait times out to request-known timeout/manual-review behavior, then cleans up the created provider resource.

```bash
APP_ENV=dev \
BILLING_CLOUDMINI_IDEMPOTENCY_EVIDENCE_APPROVED=yes \
CLOUDMINI_IDEMPOTENCY_SCENARIO=timeout-after-send \
CLOUDMINI_IDEMPOTENCY_PILOT_ID=<approved-pilot-id> \
CLOUDMINI_IDEMPOTENCY_MAX_CREATE_ATTEMPTS=1 \
CLOUDMINI_IDEMPOTENCY_MAX_ACTIVE_RESOURCES=1 \
CLOUDMINI_IDEMPOTENCY_PROVIDER_RATE_LIMIT=no-parallel-mutating-calls \
CLOUDMINI_IDEMPOTENCY_MAX_SPEND_EXPOSURE=single-dev-resource \
CLOUDMINI_IDEMPOTENCY_RAW_EVIDENCE_PATH=/tmp/<approved-private-raw-evidence>.json \
CLOUDMINI_V3_POLL_TIMEOUT=<short-approved-timeout> \
go run ./cmd/smoke cloudmini-idempotency-evidence
```

Expected redacted pass evidence:

- `create_1_error_code=PROVIDER_TIMEOUT_REQUEST_KNOWN`
- `create_1_retry_safety=manual_review_required`
- `timeout_after_send_manual_review=true`
- `cleanup_attempts=1`
- `sensitive_values_printed=no`
- `raw_provider_ids_printed=no`

## Evidence Handling

Record only stdout redacted evidence in repo docs. Keep the raw cleanup reference file outside git with mode `0600` until provider owner confirms no manual cleanup is needed.

If cleanup fails, do not retry blindly. Use the raw cleanup reference file, provider console/API, and the cleanup owner path from doc 71.
