# T104 - Provider readiness seed scenarios

Status: DONE
Owner: Codex
Branch: codex/t104-provider-readiness-seed-scenarios
PR: https://github.com/Chinsusu/Billing-V2/pull/236
Risk: seed/local-dev
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Make local seed data demonstrate provider readiness states clearly without breaking the green-path billing smoke.

## Scope

- Review current seed catalog/source/plan-source data.
- Ensure at least one ready plan remains available for paid-order provisioning smoke.
- Add safe local-only examples for non-ready readiness states where practical.
- Document which seeded display IDs are useful for readiness checks.
- Do not add real provider credentials or production provider setup.
- Keep each file under 500 lines.

## Acceptance Criteria

- Local seed data can show a ready source and at least one non-ready readiness state.
- Green-path billing smoke still uses a ready provider source.
- Seed docs identify example display IDs for operators and agents.
- Backend and frontend validation commands pass.

## Notes

- Follows T100 and T102.
- If adding missing-plan-source seed data would confuse catalog UI, document the reason and choose safer inactive/unsupported examples instead.

## Agent Log

- 2026-04-24: Task created in the provider readiness follow-up batch.
- 2026-04-24: Codex claimed the task on `codex/t104-provider-readiness-seed-scenarios`.
- 2026-04-24: Added local seed readiness examples for `ready`, `fake_provider_only`, `unsupported_capability`, and `inactive_source`; kept paid-order smoke on a ready fake provider source and documented fresh local display IDs.
- 2026-04-24: Validation passed: `go test ./internal/seed`, `go test ./internal/modules/catalog`, `go test ./cmd/smoke`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `npm ci`, `npm audit --omit=dev`, `npm run lint`, and `npm run build`.
- 2026-04-24: Opened PR #236 for review.
- 2026-04-24: CI passed on PR #236 and merged to main at `c3db616`.
