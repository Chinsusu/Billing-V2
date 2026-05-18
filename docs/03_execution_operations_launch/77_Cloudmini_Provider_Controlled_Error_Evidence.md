# 77 - Cloudmini Provider-Controlled Error Evidence

**Tasks:** T255, T256, T257, T258, T259
**Date:** 2026-05-18
**Scope:** Remaining Cloudmini V3 provider-controlled error cases for launch evidence.
**Decision:** partially closed. T255 records source-inspection evidence and a safe execution plan. T256 closes the permission-denied runtime case with a temporary low-scope key and same-run revoke. T257 closes the out-of-capacity runtime case with one bounded exhausted-group reservation probe. T258 adds a guarded Billing runner for a provider-owned rate-limit fixture. T259 adds and merges the provider-side fixture code in Cloudmini PR #9, but live rate-limit runtime evidence remains blocked until that provider code is deployed to the approved dev manager with the fixture env enabled. No provider proxy create/delete run was performed.

## Boundary

This packet does not authorize broad provider provisioning. It must not be used as GO evidence until every provider-controlled case below has redacted runtime output from an owner-approved non-production fixture or bounded test run.

Do not collect these cases by:

- spamming the shared provider until a global rate limiter trips;
- breaking the live provider database or service to force a 5xx;
- creating or deleting a proxy only to make an error path happen;
- committing raw API keys, provider IDs, response bodies, proxy credentials, DSNs, cookies, or file contents.

## Source Inspection Evidence

Provider source inspected from `/opt/proxy-cloudmini` without reading provider secrets:

| Provider behavior | Source-read evidence | Current evidence status |
| --- | --- | --- |
| V3 routes use `AuthMiddleware`; mutating and proxy read routes require `proxy_crud`. | `/opt/proxy-cloudmini/internal/api/router.go` defines `/api/v3` auth and `RequirePermission("proxy_crud")` on reservations, proxies, operations, and actions. | Supports a permission-denied test only with a valid low-scope API key. |
| Missing `proxy_crud` returns HTTP `403`. | `/opt/proxy-cloudmini/internal/api/auth_handler.go` `RequirePermission` returns forbidden when permissions do not include the required permission or wildcard. | T256 captured temporary low-scope API key runtime evidence and same-run revoke. |
| Global V2/API limiter returns `RATE_LIMITED`; T259 adds a disabled-by-default V3 fixture path. | `/opt/proxy-cloudmini/internal/api/router.go` has v1/v2 limiter responses with code `RATE_LIMITED`. Cloudmini PR #9 adds `GET /api/v3/error-fixtures/rate-limited`, guarded by `VPM_BILLING_ERROR_FIXTURES_ENABLED` and `X-Cloudmini-Error-Fixture: rate_limited`. | Runtime evidence still needs deployment to the approved dev manager and one guarded Billing smoke run. Do not induce the shared limiter. |
| Capacity exhaustion exists before reservation persistence. | `/opt/proxy-cloudmini/internal/api/handler/v3_handler.go` `CreateReservation` checks group inventory and returns `CAPACITY_EXHAUSTED` before creating a reservation when no allocatable units exist. | T257 captured owner-approved exhausted-group reservation probe evidence with no reservation created. |
| V3 returns `INTERNAL_ERROR` on repository/service failures. | `/opt/proxy-cloudmini/internal/api/handler/v3_handler.go` returns `INTERNAL_ERROR` for inventory, reservation, operation, create, delete, and action storage/service failures. | No safe shared-provider trigger; needs provider-side fixture. |
| Delete/action rejection is asynchronous for real proxies; immediate proxy-delete not-found is already covered by `PROXY_NOT_FOUND`. | `/opt/proxy-cloudmini/internal/api/handler/v3_handler.go` starts async delete/action operations and records `DELETE_FAILED` or `ACTION_FAILED` only after service failure. | Needs provider-side fixture or owner-approved controlled resource state. |

## Case Matrix

| Case | Required Billing mapping | Safe evidence path | Current status |
| --- | --- | --- | --- |
| Permission denied | HTTP `403` -> `PROVIDER_PERMISSION_DENIED`, retry `do_not_retry`. | Create a temporary non-production API key with read-only permission, call a `proxy_crud` route such as `GET /api/v3/proxies`, record redacted `403`, then revoke and verify active key count returns to the previous value. | Closed by T256: runtime `403`, normalized `PROVIDER_PERMISSION_DENIED`, `do_not_retry`, temporary key revoked, active key count restored. |
| Rate limited | HTTP `429`/`RATE_LIMITED` -> `PROVIDER_RATE_LIMITED`, retry `safe_retry`. | Use provider-owned isolated low-limit fixture or test route. Do not trip the shared 1000 req/min limiter. | Blocked after T259: provider fixture code is merged in Cloudmini PR #9, but it has not been deployed/run on the approved dev manager. |
| Out of capacity | `CAPACITY_EXHAUSTED` -> `PROVIDER_OUT_OF_STOCK`, retry `do_not_retry`. | Use an owner-approved exhausted group reservation probe with max attempt `1`, TTL no more than `60s`, and cleanup/verification if a reservation is unexpectedly created. | Closed by T257: runtime `409`, provider code `CAPACITY_EXHAUSTED`, normalized `PROVIDER_OUT_OF_STOCK`, `do_not_retry`, exhausted group selected, reservation created `false`. |
| Provider 5xx | HTTP `5xx`/`INTERNAL_ERROR` -> `PROVIDER_TEMPORARY_ERROR`, retry `safe_retry`. | Use a provider-owned non-production fixture that returns the normal V3 error envelope without breaking the shared service. | Blocked: no safe fixture exists in current Billing evidence. |
| Cancel/delete rejected | Provider delete/action failure -> `PROVIDER_PARTIAL_SUCCESS` or provider-specific failed operation mapping with manual review as applicable. | Use a provider-owned non-production fixture or a controlled resource state that rejects delete/action without deleting a sellable resource. | Blocked: no safe fixture or controlled resource-state evidence exists in current Billing evidence. |

## Provider-Side Support Needed

Before the remaining blockers can close, the provider owner must supply:

- deployment of the merged rate-limit fixture from Cloudmini PR #9 on the approved dev manager with `VPM_BILLING_ERROR_FIXTURES_ENABLED=yes`;
- a non-production error fixture that returns V3 envelopes for `INTERNAL_ERROR` and delete/action failure without side effects.

The fixture must be disabled or inaccessible in production-like customer paths unless explicitly approved by Security and Provider owners.

For the rate-limit case, the fixture contract needed by the Billing runner is:

```text
method=GET
path=/api/v3/<provider-owned-side-effect-free-fixture-path-containing-fixture-and-rate>
http_status=429
provider_error_code=RATE_LIMITED
response_shape=V3 success:false error envelope with message
side_effects=none
shared_limiter_tripped=no
required_header=X-Cloudmini-Error-Fixture: rate_limited
```

Billing will only call this fixture when these guardrails are set:

```text
CLOUDMINI_ERROR_EVIDENCE_ALLOW_RATE_LIMITED=yes
CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_APPROVED=yes
CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_MAX_REQUESTS=1
CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_FIXTURE_PATH=/api/v3/...
```

## Required Runtime Evidence Format

Each completed case must record only:

```text
case=<case_name>
http_status=<status>
provider_error_code=<stable_code_or_none>
normalized_error_code=<billing_provider_code>
retry_safety=<retry_safety>
error_envelope_present=yes|no
mutating_route_called=yes|no
side_effect_created=no|cleaned_up|not_applicable
raw_response_body_printed=no
sensitive_values_printed=no
raw_provider_ids_printed=no
```

## T256 Permission-Denied Runtime Evidence

Approved dev/test command:

```text
APP_ENV=dev ... go run ./cmd/smoke cloudmini-error-evidence
```

The command sourced Cloudmini dev/test credentials from the protected local credential file without printing file contents or secret values. It enabled only the permission-denied guardrails:

```text
CLOUDMINI_ERROR_EVIDENCE_ALLOW_PERMISSION_DENIED=yes
CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MANAGEMENT_APPROVED=yes
CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MAX_CREATE=1
```

T256 redacted stdout excerpt:

```text
cloudmini_error_evidence result=PASS
pilot_environment=dev
approval_fields_present=yes
owner_fields_present=yes
example_count=4
mutating_routes_called=true
example_4_name=permission_denied_proxy_list
example_4_http_status=403
example_4_provider_error_code=none
example_4_normalized_error_code=PROVIDER_PERMISSION_DENIED
example_4_retry_safety=do_not_retry
example_4_error_envelope_present=true
example_4_error_message_field_present=true
example_4_error_details_field_present=false
example_4_side_effect_created=cleaned_up
example_4_temporary_api_key_created=true
example_4_temporary_api_key_revoked=true
example_4_active_key_count_restored=true
raw_response_body_printed=no
sensitive_values_printed=no
raw_provider_ids_printed=no
provider_payloads_printed=no
remaining_provider_controlled_examples=rate_limited,out_of_capacity,provider_5xx,cancel_rejected
```

This evidence called provider API-key management routes to create and revoke one temporary low-scope key. It did not call provider proxy create, proxy delete, proxy action, reservation, Billing checkout, Billing payment, or provisioning worker mutation routes.

## T257 Out-Of-Capacity Runtime Evidence

Approved dev/test command:

```text
APP_ENV=dev ... go run ./cmd/smoke cloudmini-error-evidence
```

The command sourced Cloudmini dev/test credentials from the protected local credential file without printing file contents or secret values. It enabled only the out-of-capacity guardrails:

```text
CLOUDMINI_ERROR_EVIDENCE_ALLOW_OUT_OF_CAPACITY=yes
CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_APPROVED=yes
CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_MAX_RESERVATIONS=1
CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_KIND=residential
CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_TTL_SECONDS=60
```

T257 redacted stdout excerpt:

```text
cloudmini_error_evidence result=PASS
pilot_environment=dev
approval_fields_present=yes
owner_fields_present=yes
example_count=4
mutating_routes_called=true
example_4_name=out_of_capacity_reservation
example_4_http_status=409
example_4_provider_error_code=CAPACITY_EXHAUSTED
example_4_normalized_error_code=PROVIDER_OUT_OF_STOCK
example_4_retry_safety=do_not_retry
example_4_error_envelope_present=true
example_4_error_message_field_present=true
example_4_error_details_field_present=true
example_4_side_effect_created=no
example_4_reservation_probe_attempted=true
example_4_exhausted_group_selected=true
example_4_reservation_created=false
example_4_reservation_cleaned_up=false
example_4_reservation_max_attempts=1
example_4_reservation_ttl_seconds=60
raw_response_body_printed=no
sensitive_values_printed=no
raw_provider_ids_printed=no
provider_payloads_printed=no
```

This evidence called one provider reservation route against an exhausted non-production group. The provider returned before reservation persistence, so no reservation was created and no cleanup was required. It did not call provider proxy create, proxy delete, proxy action, Billing checkout, Billing payment, or provisioning worker mutation routes.

## T258 Rate-Limit Fixture Runner And Blocker

T258 did not run a live provider rate-limit case. It inspected `/opt/proxy-cloudmini/internal/api/router.go` and confirmed the currently inspected provider source exposes `RATE_LIMITED` through v1/v2 global limiters but does not expose a V3 side-effect-free low-limit fixture route. Intentionally tripping the shared 1000 req/min limiter remains prohibited.

T258 adds Billing smoke support for exactly one owner-approved fixture request. Local `httptest` coverage proves the Billing runner maps HTTP `429` with provider code `RATE_LIMITED` to `PROVIDER_RATE_LIMITED`, retry safety `safe_retry`, and redacted output:

```text
example_4_name=rate_limited_fixture
example_4_http_status=429
example_4_provider_error_code=RATE_LIMITED
example_4_normalized_error_code=PROVIDER_RATE_LIMITED
example_4_retry_safety=safe_retry
example_4_side_effect_created=not_applicable
example_4_rate_limit_fixture_called=true
example_4_rate_limit_max_requests=1
raw_response_body_printed=no
sensitive_values_printed=no
raw_provider_ids_printed=no
provider_payloads_printed=no
```

This is readiness for future fixture evidence, not provider runtime evidence. The rate-limit case remains blocked until the provider fixture above is deployed and an approved non-production run records redacted stdout.

## T259 Provider Fixture Implementation

Cloudmini provider PR #9 merged the rate-limit fixture implementation:

```text
provider_pr=https://github.com/Chinsusu/proxy-cloudmini/pull/9
method=GET
path=/api/v3/error-fixtures/rate-limited
required_provider_env=VPM_BILLING_ERROR_FIXTURES_ENABLED=yes
required_header=X-Cloudmini-Error-Fixture: rate_limited
http_status=429
provider_error_code=RATE_LIMITED
side_effects=none
```

Provider validation run before PR #9:

```text
go test -buildvcs=false ./internal/api ./internal/api/handler ./internal/api/middleware ./pkg/dto ./pkg/errcode
go build -buildvcs=false ... ./cmd/manager
go build -buildvcs=false ... ./cmd/agent
git diff --check
```

Provider full `go test -buildvcs=false ./...` was not used as pass evidence because the current provider base already has unrelated compile failures in root `test_list.go` and `cmd/test_ssh`.

No live Billing rate-limit run was performed in T259. This Billing host did not have a local `vpm-manager`/Cloudmini manager service running, so the merged provider fixture still needs deployment to the approved dev manager before Billing can collect redacted runtime evidence.

## T255 Non-Actions

T255 did not:

- create, delete, start, stop, or change a Cloudmini proxy;
- create or revoke a provider API key;
- intentionally hit provider rate limits;
- force a provider 5xx;
- create a reservation probe;
- print or record any secret value, raw provider payload, raw provider ID, proxy credential, DSN, cookie, or file content.

The launch decision remains NO-GO until the remaining blocked cases above are executed safely or a policy-allowed owner exception is recorded in docs 69 and 70. After T256 through T259, the remaining provider-controlled runtime cases are rate limited, provider 5xx, and cancel/delete rejected; T259 does not close the rate-limit runtime case until the fixture is deployed and run.
