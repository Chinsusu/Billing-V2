# T141 - Hide request trace IDs

Status: REVIEW
Owner: Codex
Branch: codex/t141-hide-internal-ui-ids
PR: https://github.com/Chinsusu/Billing-V2/pull/313
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Stop showing internal request/correlation IDs in frontend audit and provisioning UI.

## Scope

- Hide live audit `correlation_id` values from admin audit rows.
- Hide provisioning job `correlation_id` values from the job detail panel.
- Hide demo audit `req-*` request IDs in fallback audit/report rows.
- Keep public log/job/service/provider IDs visible.

## Acceptance Criteria

- Admin audit rows no longer show raw or shortened request trace IDs.
- Job detail panel no longer shows raw or shortened correlation IDs.
- Demo audit/report fallback rows no longer show `req-*` request IDs.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- Request/correlation IDs are still present in API data and can be used internally; this task only changes display labels.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T140 was marked done; starting frontend request trace ID cleanup.
- 2026-04-26: Hid live audit correlation labels, provisioning job correlation labels, recovery audit request labels, and demo `req-*` request IDs.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/313 for review.
- 2026-04-26: CI frontend gate failed because smoke still expected `req-smoke`; updated smoke to expect `Request not shown` and assert `req-smoke` stays hidden.
- 2026-04-26: Re-validation passed: `npm --prefix frontend run smoke:admin:ci`, `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
