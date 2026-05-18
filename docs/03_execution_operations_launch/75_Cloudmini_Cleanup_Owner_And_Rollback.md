# 75 - Cloudmini Cleanup Owner And Rollback

**Task:** T253  
**Date:** 2026-05-18  
**Scope:** Cleanup ownership and rollback procedure for selected Cloudmini V3 pilot runs.  
**Decision:** cleanup owner/procedure approved for the selected pilot scope only; broader launch remains blocked by the separate provider/security evidence gaps in docs 66, 69, and 70.

## Owner Record

| Role | Owner | Responsibility |
| --- | --- | --- |
| Cleanup Owner | Admin | Owns same-session cleanup, residual risk decision, and incident creation if cleanup fails. |
| Provider Owner | Admin | Owns provider access, provider support escalation, and source disable approval. |
| Ops Lead | Admin | Owns worker/source rollback and service state verification. |
| Security Owner | Admin | Owns redaction boundary and secret exposure stop condition. |

Owner assignment source: T241 records the user-provided assignment that `Admin` owns launch-day Provider, Ops, and Security roles.

## Cleanup Hierarchy

Use this order for every approved Cloudmini pilot cleanup:

1. Prefer Billing lifecycle-worker provider-backed cleanup when the service is eligible for lifecycle termination and the Cloudmini provider registry is configured.
2. Confirm the Billing service is no longer active/paid after cleanup completes.
3. Confirm the provider resource is deleted, disabled, or otherwise no longer billable.
4. If lifecycle cleanup is not applicable for the current dev/test pilot state, use an owner-approved direct Cloudmini V3 delete as the fallback cleanup exception.
5. If direct cleanup fails or status cannot be confirmed, disable the affected provider source, keep launch `NO-GO`, and open an incident/follow-up before any further create attempt.

## Evidence Required Per Run

Every future approved Cloudmini pilot run must record only redacted evidence:

- run ID and UTC timestamp;
- cleanup owner;
- cleanup path used: lifecycle worker or approved direct-provider fallback;
- redacted Billing service display ID if one exists;
- redacted provider resource reference hash;
- cleanup operation result;
- final provider status evidence such as deleted/disabled/not-found/no-longer-billable;
- final Billing service/provider mapping state;
- residual risk and owner decision.

Never record raw provider IDs, raw response bodies, provider payloads, proxy credentials, DSNs, API keys, cookies, or file contents.

## Stop Conditions

Stop all additional Cloudmini create attempts if any of these occurs:

- cleanup/delete does not complete;
- cleanup status is unknown;
- a duplicate provider resource is suspected;
- provider source cannot be disabled;
- Billing still shows an active paid resource after cleanup;
- any raw secret, provider payload, proxy credential, or customer data is exposed.

## Current Evidence

- T228 completed same-session direct provider cleanup for one dev Billing-path create/delete pilot.
- T232 completed same-session fallback cleanup after a non-usable `creating` status.
- T233 proved lifecycle-worker provider cleanup on the approved target dev/test server for one Cloudmini resource.
- T249 proved cleanup success for duplicate-create and timeout-after-send evidence scenarios.
- T253 assigns Admin as cleanup owner and records the required cleanup/rollback procedure for selected pilot runs.

## Boundaries

This cleanup owner packet does not approve:

- production/shared credential storage;
- production customer data or production provider accounts;
- new provider create/delete runs without separate owner-approved run evidence;
- increasing create limits, active-resource limits, worker concurrency, or provider rate limits;
- broader provider launch approval.

The launch decision remains NO-GO until the remaining P0 blockers in docs 69 and 70 are closed or explicitly accepted where policy allows.
