# T173 - Add agent behavior guidelines

Status: REVIEW
Owner: Codex
Branch: codex/t173-agent-behavior-guidelines
PR: https://github.com/Chinsusu/Billing-V2/pull/377
Risk: repository workflow documentation
Created: 2026-04-27
Updated: 2026-04-27

## Summary

Add the requested behavioral guidelines to `AGENTS.md` so future coding agents follow clearer expectations around correctness, verification, minimal diffs, and communication.

## Scope

- Add the requested guideline content to `AGENTS.md`.
- Preserve existing project-specific workflow and safety instructions.
- Do not change runtime code.

## Acceptance Criteria

- `AGENTS.md` includes the requested behavioral guidance.
- `AGENTS.md` remains below the 500-line file limit.
- Taskguard and diff check pass.

## Notes

- User requested adding the guideline block directly to `AGENTS.md`.

## Agent Log

- 2026-04-27: Task created and claimed by Codex.
- 2026-04-27: Added behavioral guidelines to `AGENTS.md`; taskguard and diff check pass.
- 2026-04-27: Opened PR #377 and moved task to REVIEW.
