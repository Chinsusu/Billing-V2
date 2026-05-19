# T278 - Telegram notification delivery

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t278-telegram-delivery
PR: -
Risk: notifications, secrets, external delivery, launch readiness
Created: 2026-05-19
Updated: 2026-05-19

## Summary

Add a Telegram notification delivery path so an owner-provided Telegram bot/channel can be used for redacted launch-critical notification evidence.

## Scope

- Add a Telegram delivery handler for queued `telegram` notifications.
- Add `notification-telegram-once` and `notification-telegram-loop` worker commands.
- Read Telegram bot token and chat ID from environment without logging values.
- Add a safe redacted Telegram preflight command for channel evidence.
- Update notification launch docs to describe the Telegram proof path.
- Keep SMTP delivery, customer-facing email delivery, production customer data, and broader launch GO out of scope unless real evidence passes.

## Acceptance Criteria

- Telegram worker commands require `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID`.
- Telegram delivery sends only a safe rendered summary, not raw `payload_redacted`.
- Delivery failures mark notifications failed with redacted error metadata.
- Tests cover success, HTTP failure, missing config, and worker command guards.
- If server secrets are unavailable, task records that live Telegram evidence was not run and why.
- Docs continue to distinguish Telegram proof from selected-pilot manual fallback.
- Task board remains consistent and required checks pass.

## Notes

- A Telegram channel alone is not enough; the worker needs a bot token and chat/channel ID through the approved secret path.
- Current key-only scan did not find `TELEGRAM_*`, `BOT_TOKEN`, or `CHAT_ID` keys under `/etc/billing/secrets`, `/etc`, or `/opt`.
- No Telegram token, chat ID, notification payload, DSN, provider payload, or customer data may be printed or committed.

## Agent Log

- 2026-05-19: Task created and claimed by Codex from Billing `origin/main`; user confirmed a Telegram channel exists, but key-only local scan found no Telegram bot/chat secret configured on the server.
- 2026-05-19: Added channel-specific Telegram delivery job type so Telegram worker does not claim dashboard/email notification jobs; local/dev worker can still process generic and Telegram jobs.
- 2026-05-19: Added Telegram delivery handler, sender, safe message rendering, retry/terminal failure classification, and redacted failure metadata.
- 2026-05-19: Added `notification-telegram-once`, `notification-telegram-loop`, and `notification-telegram-preflight` commands. Production runs require `BILLING_TELEGRAM_DELIVERY_PRODUCTION_APPROVED=yes`.
- 2026-05-19: Updated `.env.example` and notification runbook with Telegram config keys and redaction/evidence rules.
- 2026-05-19: Live Telegram preflight was attempted against current server env and failed closed with missing `TELEGRAM_BOT_TOKEN`; no token, chat ID, payload, DSN, provider payload, or customer data was printed.
- 2026-05-19: Validation passed: `go test ./internal/modules/notification ./cmd/worker`, `GOFLAGS=-buildvcs=false make test`, `GOFLAGS=-buildvcs=false go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `git diff --check`, diff raw-secret pattern scan, and changed-file line-count check.
