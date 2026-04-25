# API Contract Drift Guard

**Scope:** Lightweight guard for keeping implemented billing routes aligned with the operational API reference.

## Command

Run:

```bash
make contract-guard
```

This runs:

```bash
go run ./cmd/contractguard
```

CI runs the same command on pull requests and pushes to `main`.

## What It Checks

The guard tracks selected stable billing API routes. For each tracked route it checks:

- backend route wiring or action constants still exist;
- expected RBAC permission constants are still wired;
- the route is documented in `docs/05_development_standards/56_Billing_API_Operational_Reference.md`;
- documented permissions, query names, and important response/redaction notes have not drifted.

It is not a full OpenAPI generator. Add only stable route groups that frontend, smoke tests, or operations depend on.

## When To Update

Update `cmd/contractguard/main.go` in the same PR when you intentionally:

- add, remove, or rename a tracked billing route;
- change a tracked route permission;
- add, remove, or rename a query parameter used by frontend or operations;
- change response redaction expectations such as hiding provider credentials, `payload_json`, or `idempotency_key`.

Also update the operational reference in the same PR. The guard output names the route and the missing token so another agent can fix the docs quickly.

## Failure Handling

If `make contract-guard` fails:

1. Check whether backend behavior changed intentionally.
2. If yes, update the operational reference and the guard manifest together.
3. If no, restore the route wiring, permission, query name, or redaction note that drifted.
