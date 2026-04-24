# T104 - Provider readiness seed scenarios

Status: TODO
Owner: -
Branch: codex/t104-provider-readiness-seed-scenarios
PR: -
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
