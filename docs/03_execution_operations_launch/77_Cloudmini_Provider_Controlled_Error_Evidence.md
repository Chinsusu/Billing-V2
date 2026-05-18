# 77 - Cloudmini Provider-Controlled Error Evidence

**Task:** T255  
**Date:** 2026-05-18  
**Scope:** Remaining Cloudmini V3 provider-controlled error cases for launch evidence.  
**Decision:** not closed. T255 records source-inspection evidence and a safe execution plan only. No new provider create/delete run was performed.

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
| Missing `proxy_crud` returns HTTP `403`. | `/opt/proxy-cloudmini/internal/api/auth_handler.go` `RequirePermission` returns forbidden when permissions do not include the required permission or wildcard. | Needs temporary read-only API key evidence. |
| Global V2/API limiter returns `RATE_LIMITED`; no V3-specific low-limit fixture was found. | `/opt/proxy-cloudmini/internal/api/router.go` has limiter responses with code `RATE_LIMITED`. | Do not induce this on the shared provider; needs isolated fixture or low-limit test config. |
| Capacity exhaustion exists before reservation persistence. | `/opt/proxy-cloudmini/internal/api/handler/v3_handler.go` `CreateReservation` checks group inventory and returns `CAPACITY_EXHAUSTED` before creating a reservation when no allocatable units exist. | Needs owner-approved exhausted-group reservation probe or fixture. |
| V3 returns `INTERNAL_ERROR` on repository/service failures. | `/opt/proxy-cloudmini/internal/api/handler/v3_handler.go` returns `INTERNAL_ERROR` for inventory, reservation, operation, create, delete, and action storage/service failures. | No safe shared-provider trigger; needs provider-side fixture. |
| Delete/action rejection is asynchronous for real proxies; immediate proxy-delete not-found is already covered by `PROXY_NOT_FOUND`. | `/opt/proxy-cloudmini/internal/api/handler/v3_handler.go` starts async delete/action operations and records `DELETE_FAILED` or `ACTION_FAILED` only after service failure. | Needs provider-side fixture or owner-approved controlled resource state. |

## Case Matrix

| Case | Required Billing mapping | Safe evidence path | T255 status |
| --- | --- | --- | --- |
| Permission denied | HTTP `403` -> `PROVIDER_PERMISSION_DENIED`, retry `do_not_retry`. | Create a temporary non-production API key with read-only permission, call a `proxy_crud` route such as `GET /api/v3/proxies`, record redacted `403`, then revoke and verify active key count returns to the previous value. | Blocked: no temporary low-scope key evidence was created in this task. |
| Rate limited | HTTP `429`/`RATE_LIMITED` -> `PROVIDER_RATE_LIMITED`, retry `safe_retry`. | Use provider-owned isolated low-limit fixture or test route. Do not trip the shared 1000 req/min limiter. | Blocked: no safe low-limit fixture exists in current Billing evidence. |
| Out of capacity | `CAPACITY_EXHAUSTED` -> `PROVIDER_OUT_OF_STOCK`, retry `do_not_retry`. | Use an owner-approved exhausted group reservation probe with max attempt `1`, TTL no more than `60s`, and cleanup/verification if a reservation is unexpectedly created. | Blocked: no bounded reservation probe was run in this task. |
| Provider 5xx | HTTP `5xx`/`INTERNAL_ERROR` -> `PROVIDER_TEMPORARY_ERROR`, retry `safe_retry`. | Use a provider-owned non-production fixture that returns the normal V3 error envelope without breaking the shared service. | Blocked: no safe fixture exists in current Billing evidence. |
| Cancel/delete rejected | Provider delete/action failure -> `PROVIDER_PARTIAL_SUCCESS` or provider-specific failed operation mapping with manual review as applicable. | Use a provider-owned non-production fixture or a controlled resource state that rejects delete/action without deleting a sellable resource. | Blocked: no safe fixture or controlled resource-state evidence exists in current Billing evidence. |

## Provider-Side Support Needed

Before the blocker can close, the provider owner must supply one of:

- temporary low-scope credentials plus revoke evidence for permission-denied;
- a non-production error fixture that returns V3 envelopes for `RATE_LIMITED`, `INTERNAL_ERROR`, and delete/action failure without side effects;
- an owner-approved exhausted-group reservation probe with hard bounds and cleanup verification.

The fixture or probe must be disabled or inaccessible in production-like customer paths unless explicitly approved by Security and Provider owners.

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

## T255 Non-Actions

T255 did not:

- create, delete, start, stop, or change a Cloudmini proxy;
- create or revoke a provider API key;
- intentionally hit provider rate limits;
- force a provider 5xx;
- create a reservation probe;
- print or record any secret value, raw provider payload, raw provider ID, proxy credential, DSN, cookie, or file content.

The launch decision remains NO-GO until the blocked cases above are executed safely or a policy-allowed owner exception is recorded in docs 69 and 70.
