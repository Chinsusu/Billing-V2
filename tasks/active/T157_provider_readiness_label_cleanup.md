# T157 - Provider readiness label cleanup

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t157-provider-readiness-labels
PR: -
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Show readable product and provider source labels in the admin provider readiness panel instead of raw API keys.

## Scope

- Use shared display label helpers for provider readiness product and source badges.
- Keep API filter values and query payloads unchanged.
- Update lightweight frontend validation if existing smoke coverage can catch the visible labels.

## Acceptance Criteria

- Provider readiness rows show labels such as VPS, Proxy, Hetzner, Manual pool, and Proxmox instead of raw keys.
- Readiness filters still submit raw product and source values to the backend.
- Frontend lint, sensitive-text check, production build, taskguard, and diff check pass.

## Notes

- This is frontend-only and should not change API contracts or backend behavior.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T156 was marked done; starting provider readiness label cleanup.
- 2026-04-26: Applied shared display helpers to readiness product/source badges and added browser smoke coverage for the readable product label.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, and taskguard.
