# T134 - Client service public identifier

Status: DONE
Owner: Codex
Branch: codex/t134-client-service-public-identifier
PR: https://github.com/Chinsusu/Billing-V2/pull/299
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Stop showing raw provider resource IDs in the client service table and keep the visible service/order references on public display IDs.

## Scope

- Replace the client service `Identifier` column with a plan label derived from snapshots.
- Keep the first service column on `SVC-` public IDs.
- Keep order references on `ORD-` public IDs.
- Leave backend IDs only in API action bodies or internal joins.

## Acceptance Criteria

- Client service rows no longer render `external_resource_id` as a visible identifier.
- Client service plan/source/order labels use public or redacted labels.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- `external_resource_id` can still be used internally for category detection, but it should not be shown as the client-facing service identifier.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T133 was marked done; starting client service visible identifier cleanup.
- 2026-04-26: Replaced client service visible raw provider identifier with snapshot plan labels while keeping service/order public labels.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/299 for review.
- 2026-04-26: PR https://github.com/Chinsusu/Billing-V2/pull/299 merged into `main`; marking task done.
