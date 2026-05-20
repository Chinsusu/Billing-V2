# T279 - Telegram live delivery evidence

Status: REVIEW
Owner: Codex
Branch: codex/t279-telegram-live-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/589
Risk: notifications, secrets, launch evidence
Created: 2026-05-20
Updated: 2026-05-20

## Summary

Record the selected-host Telegram notification preflight evidence after the owner configured the approved Telegram channel secret file.

## Scope

- Capture only redacted evidence for the selected host Telegram preflight.
- Update the notification runbook, Go/No-Go record, and launch evidence packet.
- Keep token, chat ID, raw Telegram request/response body, notification payloads, DSNs, provider payloads, cookies, and customer data out of git.
- Do not change notification runtime code or production launch scope.

## Acceptance Criteria

- Evidence states that `notification-telegram-preflight` passed with `telegram_api_called=yes`, `message_payload_redacted=yes`, and `secrets_printed=no`.
- Evidence records only secret-file metadata and key names, not values.
- Evidence records process argv secret check result without printing command lines.
- Docs distinguish selected-host Telegram preflight proof from broader production notification approval.
- Task board remains consistent and required checks pass.

## Notes

- The owner corrected `TELEGRAM_CHAT_ID` on the selected host before the final preflight.
- The selected-host secret file was verified as mode `600`, owner `root:root`.
- Preflight ran with `APP_ENV=staging` and did not use DB access or queued customer notifications.

## Agent Log

- 2026-05-20: Task created and claimed by Codex from Billing `origin/main`.
- 2026-05-20: Verified redacted runtime evidence: Telegram preflight returned `PASS`; Telegram API was called; message payload was redacted; no secrets were printed; process argv secret match count excluding checker was `0`; `taskguard` passed.
- 2026-05-20: Opened PR #589. Local validation passed: `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `git diff --check`, diff raw-secret pattern scan, and changed-file line-count check.
