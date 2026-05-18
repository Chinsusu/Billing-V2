# 74 - Cloudmini Usable Status Sign-Off

**Task:** T252  
**Date:** 2026-05-18  
**Scope:** Cloudmini V3 usable-status semantics for the approved dev/test and selected pilot provisioning scope.  
**Decision:** approved for the current Cloudmini pilot semantics only; broader launch remains blocked by the separate provider/security evidence gaps in docs 66, 69, and 70.

## Owner Record

| Role | Owner | Sign-off |
| --- | --- | --- |
| Provider Owner | Admin | Approved for the semantics below. |
| Engineering Lead | Admin | Approved for the semantics below. |
| Ops Lead | Admin | Approved for the semantics below. |
| Security Owner | Admin | Approved for the fail-closed and no-secret-output boundaries below. |

Owner assignment source: T241 records the user-provided assignment that `Admin` owns the launch-day Provider, Engineering, Ops, and Security roles.

## Approved Semantics

Cloudmini operation success is not sufficient by itself to activate a Billing service.

Billing may create an active service only when the Cloudmini proxy status is one of:

- `running`
- `active`
- `ready`
- `available`

Billing must fail closed for all other statuses, including empty, unknown, pending, `creating`, unrecognized, or provider-changed values:

- Do not create an active service.
- Do not return a credential to the customer.
- Return provider partial-success/manual-review behavior.
- Preserve the provider request/resource reference for manual cleanup without printing raw provider IDs in evidence.

After a successful Cloudmini create operation, Billing may poll `GET /api/v3/proxies/:id` within the configured bounded timeout. If the proxy does not become usable before the timeout, Billing must keep the result in manual review.

If the proxy becomes usable but credential fields are missing, Billing must fail closed with credential-missing handling and must not create an active credential.

For service termination, Billing must call provider delete before marking the Billing service terminated when a provider-backed lifecycle runner is configured. Timeout, unknown cleanup status, or provider cleanup error must stop the lifecycle transition and require manual review.

## Evidence References

- T229 adds fail-closed repo behavior and tests for non-usable Cloudmini statuses and provider-backed lifecycle cleanup.
- T230 deploys and build-tests the hardening on the approved test server without provider mutations.
- T232 proves a `creating` resource is manual-reviewed and cleaned up by fallback.
- T233 proves bounded status polling, successful activation only after a usable status, encrypted credential storage, lifecycle-worker cleanup, and final provider `404`.
- T249 proves timeout-after-send maps to `PROVIDER_TIMEOUT_REQUEST_KNOWN` and `manual_review_required` with cleanup.
- T251 stabilizes the provider hardening test so CI consistently verifies the non-usable-status manual-review behavior.

## Boundaries

This sign-off does not approve:

- production/shared credential storage;
- production customer data or production provider accounts;
- provider-controlled permission-denied, rate-limited, out-of-capacity, provider-5xx, or cancel-rejected evidence gaps;
- increasing create limits, active-resource limits, worker concurrency, or provider rate limits;
- broader provider launch approval outside the selected pilot scope.

The launch decision remains NO-GO until the remaining P0 blockers in docs 69 and 70 are closed or explicitly accepted where policy allows.
