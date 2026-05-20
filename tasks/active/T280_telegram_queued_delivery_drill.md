# T280 - Queued Telegram notification delivery drill

Status: TODO
Owner: -
Branch: codex/t280-telegram-queued-delivery-drill
PR: -
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
