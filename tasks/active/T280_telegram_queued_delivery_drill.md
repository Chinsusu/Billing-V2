# T280 - Queued Telegram notification delivery drill

Status: REVIEW
Owner: Codex
Branch: codex/t280-telegram-queued-delivery-drill
PR: https://github.com/Chinsusu/Billing-V2/pull/592
Risk: notifications, secrets, external delivery, DB/job mutation, launch readiness
Created: 2026-05-20
Updated: 2026-05-20

## Summary

Run a bounded selected-host drill that sends one queued `telegram` notification through the real Telegram delivery worker, extending T279 from preflight reachability to queued delivery evidence.

## Scope

- Use the existing selected-host Telegram secret file and record metadata/key names only.
- Use a non-production/dev-test notification record with a safe event type, template key, display ID, and correlation ID.
- Enqueue or select exactly one queued `telegram` delivery job through the existing app/store path where possible.
- Run `notification-telegram-once` with a protected DB DSN source and a dedicated worker ID.
- Capture only worker summary counts, safe notification display/event/template/correlation fields, final notification status, and redacted evidence that the message reached the approved Telegram channel.
- Keep Telegram token, chat ID, raw Telegram request/response body, raw `payload_redacted`, DB DSN, cookies, provider payloads, credentials, reset tokens, and customer data out of stdout, docs, PRs, and task logs.
- Do not approve broader production notification delivery unless failure/retry evidence or a signed owner exception is recorded.

## Acceptance Criteria

- Secret metadata check passes for the Telegram env file: mode `600`, owner `root:root`, required keys present, values not printed.
- Process argv check reports no Telegram token/chat ID exposure outside the checker process.
- One non-production queued Telegram notification is claimed by `notification-telegram-once`.
- Worker result shows the queued job succeeded and does not claim dashboard/email jobs.
- Notification evidence includes only safe display/event/template/correlation/status fields.
- The approved Telegram channel receives the redacted queued-drill message.
- If retry/failure behavior is not drilled in this task, the task records the residual gap and keeps Telegram primary-path promotion out of scope.
- Required local checks pass before PR: `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `git diff --check`, raw-secret diff scan, and changed-file line-count check.

## Notes

- T278 implemented the Telegram delivery worker and preflight command.
- T279 proved selected-host Telegram preflight reachability with redacted payload and no secret exposure.
- This task should not use production customer data or production notification payloads.
- Prefer a temporary dev/test record or existing safe seed fixture. Clean up or clearly mark any drill artifact according to the existing notification/job data rules.

## Agent Log

- 2026-05-20: Task created by Codex as a follow-up after T279. Keep the scope bounded to one queued dev/test Telegram delivery and redacted evidence only.
- 2026-05-20: Claimed by Codex from latest `origin/main`; starting with notification/job schema inspection and selected-host secret metadata checks before any queued delivery.
- 2026-05-20: Read-only preflight passed: selected-host Telegram secret file mode `600` owner `root:root`, required keys present, DB reachable through protected service-file handling, tenant count `2`, claimable Telegram jobs `0`, and claimable generic notification jobs `0`; no DSN, token, chat ID, payload, UUID, or command line was printed.
- 2026-05-20: Created one bounded dev/test queued Telegram notification and job with safe fields only: notification display `10000`, job display `10000`, event `service.lifecycle`, template `t280.telegram.queued_drill`, channel `telegram`, recipient group `ops`, correlation present but not printed.
- 2026-05-20: Ran `APP_ENV=staging GOFLAGS=-buildvcs=false go run ./cmd/worker notification-telegram-once -worker-id t280-telegram-drill -batch-size 1 -timeout 60s`; result `claimed=1 succeeded=1 retried=0 manual_review=0 terminal_failed=0 cancelled=0`.
- 2026-05-20: Post-run DB verification passed: notification `10000` status `sent`, `sent_at` present, no notification error code, job `10000` status `succeeded`, attempt count `1`, one succeeded attempt row, claimable Telegram jobs `0`, and claimable generic notification jobs `0`.
- 2026-05-20: Process argv secret checks before and after worker reported `0` Telegram token/chat ID matches excluding the checker process; secrets and command lines were not printed.
- 2026-05-20: Failure/retry behavior was not drilled in this task; Telegram primary-path promotion still requires a failure/retry drill or owner-signed exception.
- 2026-05-20: Opened PR #592. Local validation passed: `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `GOFLAGS=-buildvcs=false go test ./internal/modules/notification ./cmd/worker`, `git diff --check`, diff raw-secret pattern scan, and changed-file line-count check.
