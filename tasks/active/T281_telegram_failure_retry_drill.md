# T281 - Telegram failure/retry drill

Status: DONE
Owner: Codex
Branch: codex/t281-telegram-failure-retry-drill
PR: https://github.com/Chinsusu/Billing-V2/pull/594
Risk: notifications, secrets, DB/job mutation, external delivery
Created: 2026-05-20
Updated: 2026-05-20

## Summary

Prove Telegram worker failure classification for retryable and terminal outcomes using a local fake Telegram API endpoint, without calling Telegram or exposing secrets.

## Scope

- In scope: create controlled dev/test notification jobs, run `notification-telegram-once` against a local fake API returning retryable and terminal HTTP statuses, record redacted evidence, and cleanup retryable queue artifacts.
- In scope: update notification launch/runbook evidence and Go/No-Go docs.
- Out of scope: production Telegram activation, real Telegram failure injection, customer-facing messages, and new notification worker behavior.

## Acceptance Criteria

- Retryable case records `failed_retryable` job behavior from a controlled fake Telegram API response.
- Terminal case records `failed_terminal` job behavior from a controlled fake Telegram API response.
- No raw Telegram token, chat ID, DB DSN, notification payload, customer data, UUID, cookie, provider payload, or credential is printed or committed.
- Retryable drill artifact is cleaned up so no claimable Telegram notification job remains after the run.
- Relevant docs record the evidence, limits, and residual risk.
- Required local validation and CI pass before merge.

## Notes

- Use fake `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID` values with a local `127.0.0.1` API base URL.
- Read real secret files only for metadata/argv safety checks; do not print values.

## Agent Log

- 2026-05-20: Task created and claimed on `codex/t281-telegram-failure-retry-drill`.
- 2026-05-20: Added a focused unit test for Telegram HTTP 400 terminal classification.
- 2026-05-20: Ran selected-host drill with a local fake Telegram API endpoint. HTTP 500 produced `retried=1`, notification `10001` failed with `telegram_http_500`, job `10001` `failed_retryable`, and one failed_retryable attempt row. HTTP 400 produced `terminal_failed=1`, notification `10002` failed with `telegram_http_400`, job `10002` `failed_terminal`, and one failed_terminal attempt row.
- 2026-05-20: Cancelled the retryable drill artifact after evidence capture; post-cleanup claimable Telegram and generic notification jobs were `0`.
- 2026-05-20: Verified process argv checks before and after the drill showed `0` Telegram token/chat ID/DB_DSN matches, excluding the checker process.
- 2026-05-20: Validation passed: `GOFLAGS=-buildvcs=false go test ./internal/modules/notification ./cmd/worker`, `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `git diff --check`, diff secret scan, `GOFLAGS=-buildvcs=false make test`, and `GOFLAGS=-buildvcs=false make build`.
- 2026-05-20: Opened PR https://github.com/Chinsusu/Billing-V2/pull/594 and moved task to REVIEW.
- 2026-05-20: PR https://github.com/Chinsusu/Billing-V2/pull/594 merged after GitHub checks passed; task marked DONE.
