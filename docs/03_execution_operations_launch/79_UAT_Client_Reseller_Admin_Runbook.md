# 79 - UAT Client Reseller Admin Runbook

**Date:** 2026-05-21
**Scope:** User acceptance testing for the client, reseller, and admin portals in the selected bounded non-production pilot environment.
**Decision boundary:** This runbook can approve test/UAT continuation only. It does not approve production launch, production customer data, broad private beta, or unbounded real-provider provisioning.

## Purpose

UAT verifies that the three user-facing roles can complete realistic work without breaking the P0 invariants: money correctness, tenant isolation, RBAC, credential safety, provider cleanup, notifications, and audit evidence.

This is not a marketing demo. Any P0 issue stops UAT and requires investigation before continuing.

## UAT Scope

Environment:

```text
Selected non-production/test runtime only.
No production customer data.
No production payment rail.
Approved Cloudmini selected test scope only.
```

Portals:

- Client portal.
- Reseller portal.
- Admin portal.

Roles:

- Client user under an approved test tenant.
- Reseller owner or staff under an approved test reseller tenant.
- Admin with 2FA enabled.

## Entry Criteria

Start UAT only when all fields below are true or explicitly waived by the UAT owner:

```text
Main branch CI status: pass
Task board status: TODO=0, IN_PROGRESS=0, REVIEW=0
Target health/readiness: pass
Portal domain mapping: verified and active for each tested domain
Admin 2FA enrolled and enforced: pass
Test accounts created: client, reseller, admin
Test wallet/top-up path ready: pass
Provider mode: approved selected non-production Cloudmini scope or fake/manual provider
Provider quota: within owner-approved test limit
Notification path: Telegram test path or manual fallback available
Secret handling: no raw secret in process argv or evidence
Evidence storage path: redacted non-repo or safe docs reference
```

## Required Test Data

Use safe display IDs and redacted references in evidence:

```text
Client test account:
Reseller test account:
Admin test account:
Client tenant display/reference:
Reseller tenant display/reference:
Sellable test plan:
Wallet/top-up fixture:
Provider source mode:
Notification channel/fallback owner:
```

Do not paste raw UUIDs, cookies, session tokens, DSNs, provider IDs, provider payloads, service credentials, Telegram tokens, chat IDs, TOTP secrets, or customer data into UAT evidence.

## Portal Test Matrix

| Area | Client UAT | Reseller UAT | Admin UAT |
| --- | --- | --- | --- |
| Login/session | Login, logout, expired/invalid session denied | Login, logout, scoped session | Login, 2FA gate, invalid session denied |
| Dashboard | Shows only own tenant data | Shows reseller-scope summary only | Shows platform overview without leaking secrets |
| Catalog | Sellable plan visible with safe labels | Reseller storefront/catalog scope visible | Master catalog/source visibility and safe labels |
| Top-up | Create top-up request if supported | Review reseller/client wallet scope if supported | Approve and reject test top-ups |
| Checkout | Buy approved test plan | Assist/create scoped order if supported | Inspect order/provisioning safely |
| Service | See own service, status, expiry | See reseller/client services in scope only | See service state and lifecycle evidence |
| Renewal | Renew own service using test wallet | Renew/assist scoped client if supported | Verify renewal invoice, ledger, audit |
| Credential reveal | Reveal own credential only, no-store headers | Cannot reveal outside allowed scope | Audit reveal without plaintext in evidence |
| Invoice/payment | See own invoices/payments | See reseller-scope invoices/payments | Finance reconciliation and payment detail |
| Notification | Receive or verify safe notification/fallback | Verify scoped notification visibility | Verify Telegram/fallback ops evidence |
| Support/abuse | Create/read own ticket if available | Manage scoped support if available | Support/abuse handling without sensitive leak |
| RBAC negative | Cannot access reseller/admin routes | Cannot access admin/cross-tenant routes | Low-permission admin route denied |
| Audit | No sensitive audit access | Scoped audit only if supported | Required audit events present |

## Execution Sequence

Run UAT in this order to keep evidence coherent and cleanup possible:

1. Baseline health: target health/readiness, app version, task board clean, CI pass.
2. Auth smoke: login/logout for client, reseller, and admin; verify admin 2FA gate.
3. Client happy path: top-up request, admin approval, checkout, provisioning, service detail, credential reveal, renewal.
4. Reseller path: reseller login, storefront/catalog scope, client/service visibility, wallet/invoice scope, negative cross-tenant checks.
5. Admin operations path: top-up approve/reject, order/provisioning inspect, service lifecycle inspect, finance reconciliation, notification/fallback check, audit review.
6. Negative/security pass: invalid session, low permission, cross-tenant request, credential reveal denied, direct backend ID probing if safe.
7. Cleanup: provider resource cleanup, no claimable stuck jobs, notification queue safe, finance reconciliation balanced.
8. Closeout: record evidence packet, bug list, residual risks, and owner sign-off.

## Client UAT Checklist

Required checks:

- Client can log in and log out.
- Client sees only own dashboard, services, invoices, payments, and notifications.
- Client can submit a test top-up request if the flow is enabled.
- Client can buy the approved test plan after wallet funding.
- Client can see the resulting service and lifecycle state.
- Client can reveal only own service credential through the approved reveal flow.
- Credential reveal response uses no-store headers and does not expose credentials in logs/evidence.
- Client can renew an eligible test service.
- Client cannot access reseller/admin routes.
- Client cannot access another tenant's service, invoice, wallet, ticket, or credential.

## Reseller UAT Checklist

Required checks:

- Reseller can log in and log out.
- Reseller dashboard shows reseller-scope data only.
- Reseller storefront/catalog labels are safe and do not show backend UUIDs.
- Reseller can view only clients/services/invoices/payments within reseller scope.
- Reseller cannot access platform admin routes.
- Reseller cannot access another reseller tenant.
- Reseller cannot reveal credentials outside explicitly allowed scope.
- Reseller support or client-management actions, if enabled, stay tenant-scoped.

## Admin UAT Checklist

Required checks:

- Admin login requires 2FA where configured.
- Admin can approve and reject test top-ups exactly once.
- Approved top-up creates one posted ledger credit and correct audit evidence.
- Rejected top-up creates no ledger credit and records audit evidence.
- Admin can inspect order, provisioning, service, invoice, payment, and audit records with safe display IDs.
- Finance reconciliation remains balanced after UAT mutations.
- Provider resource count stays within approved quota and cleanup is recorded.
- Notification or manual fallback path is available and redacted.
- Low-permission/admin-negative checks are denied.
- No admin view or evidence reveals secrets, plaintext credentials, raw provider payloads, cookies, reset tokens, or DSNs.

## Negative Checks

Run these before sign-off:

| ID | Check | Expected |
| --- | --- | --- |
| UAT-NEG-001 | Client attempts admin route | Denied with no data leak |
| UAT-NEG-002 | Client attempts another client's service | Denied or not found; no credential |
| UAT-NEG-003 | Reseller attempts another reseller tenant | Denied or not found |
| UAT-NEG-004 | Low-permission staff attempts top-up approval | Denied; no ledger entry |
| UAT-NEG-005 | Unauthorized credential reveal | Denied; no plaintext |
| UAT-NEG-006 | Duplicate top-up approval attempt | No double credit |
| UAT-NEG-007 | Duplicate checkout/renewal click if safe to test | No duplicate debit/resource |
| UAT-NEG-008 | Invalid/expired session request | Denied |

## Evidence Rules

Allowed evidence:

- Timestamps.
- Environment labels.
- Safe domain names.
- Public display IDs.
- Worker/job summary counts.
- HTTP status categories where safe.
- Redacted error codes.
- Owner names and sign-off decisions.

Forbidden evidence:

- Raw customer data.
- Raw payment proof.
- Raw provider request/response payload.
- Provider credential, token, API key, or auth header.
- DB DSN or database password.
- Cookies, session tokens, reset tokens, TOTP secrets, or TOTP codes.
- Plaintext service credentials.
- Backend UUIDs unless Security explicitly approves the redacted reference.

## Bug Severity

| Severity | Definition | UAT decision |
| --- | --- | --- |
| P0 | Tenant leak, auth bypass, money mismatch, credential leak, duplicate provider resource, unsafe cleanup failure | Stop UAT; NO-GO for pilot continuation |
| P1 | Checkout, top-up, provisioning, renewal, login, or finance reconciliation broken | Block sign-off until fixed or explicitly waived |
| P2 | Important UI/API issue with safe workaround and no P0/P1 risk | Can continue with owner acceptance |
| P3 | Copy, layout, or minor usability issue | Track follow-up |

## Exit Criteria

UAT can pass only when:

- P0 bugs are zero.
- P1 bugs are zero or have explicit owner-approved deferral that does not weaken P0 controls.
- Tenant/RBAC negative checks pass.
- Finance reconciliation is balanced.
- Credential reveal audit/no-store behavior passes.
- Provider resource cleanup is verified or no real provider resource was created.
- Notification or manual fallback path is verified for the test scope.
- Evidence packet is complete and redacted.
- Client, reseller, admin, QA, Ops, Finance, Security, and Support sign-off fields are filled.

## UAT Evidence Packet Template

```text
UAT ID:
Date/time UTC:
Environment:
Domains tested:
Evidence collector:
Client tester:
Reseller tester:
Admin tester:
QA owner:
Ops owner:
Finance owner:
Security owner:
Support owner:
Provider mode:
Provider resource mutation: yes/no
Provider cleanup result:
Notification path tested:
Finance reconciliation result:
Client UAT result: PASS/FAIL
Reseller UAT result: PASS/FAIL
Admin UAT result: PASS/FAIL
Negative checks result: PASS/FAIL
Open P0 bugs:
Open P1 bugs:
Deferred P2/P3 bugs:
Redaction review:
Decision: PASS / FAIL / BLOCKED
Residual risk:
Sign-off:
```

## Follow-Up Tasks

Recommended next tasks:

- Run client UAT evidence against the selected test environment.
- Run reseller UAT evidence against the selected test environment.
- Run admin UAT evidence against the selected test environment.
- Create bug-fix tasks for any P0/P1 findings before broader pilot continuation.
