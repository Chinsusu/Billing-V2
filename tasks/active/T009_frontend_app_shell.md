# T009 - Frontend App Shell

Status: TODO
Owner: -
Branch: feat/frontend-app-shell
PR: -
Risk: frontend
Created: 2026-04-22
Updated: 2026-04-22

## Summary

Build a runnable Next.js/React/TypeScript frontend app shell with package scripts, working navigation, screen registry, mock data, and build validation.

## Scope

- Create `frontend/package.json`.
- Use Next.js App Router, React, TypeScript, and Tailwind CSS.
- Add runnable `dev`, `build`, and `preview` scripts.
- Add app shell, shared layout, navigation, screen registry, mock data, and initial screens.
- Keep Node.js as frontend toolchain only.

## Acceptance Criteria

- Static HTML alone is not accepted.
- No Node.js backend, Express/Nest/Fastify service, Next API route, or Next Server Action is used for Billing business logic.
- Mock data is separated from components.
- Navigation works between screens.
- `npm run build` passes from `frontend/`.
- Backend `make test` and `make build` still pass unless frontend-only CI is added separately.

## Notes

- Follow `docs/05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md`.
- Do not wire production backend routes in the app-shell phase.

## Agent Log

- 2026-04-22: Task file created from `TASKS.md`.
