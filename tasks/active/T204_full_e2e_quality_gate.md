# T204 - Full E2E quality gate

Status: TODO
Owner: -
Branch: codex/t204-full-e2e-quality-gate
PR: -
Risk: CI quality gates, database smoke, billing flow, provisioning, and release safety
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Create a repeatable full end-to-end quality gate for the launch-critical billing and provisioning flow.

## Scope

- Run or wire local/CI parity for DB migration, seed, API, worker, frontend build, and smoke flows.
- Include `dev-db`, `dev-api`, and `dev-billing` smoke where prerequisites are available.
- Document any environment prerequisites and blocked checks.
- Do not point any smoke to production services.

## Acceptance Criteria

- A single documented validation path proves top-up, checkout, payment, job creation, worker provisioning, service activation, and frontend smoke.
- CI or documented local equivalent reports clear pass/fail output.
- Missing prerequisites are explicit blockers with remediation steps.
- Relevant validation and CI pass.

## Notes

- This task may depend on T189-T202 for full launch coverage, but it can start by documenting and automating current coverage.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
