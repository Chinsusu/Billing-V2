# Provider Sandbox Contract Checklist

**Scope:** Readiness checklist for enabling a non-fake provider in sandbox before real provider implementation work starts.

## Read First

References:

- Provider adapter spec: `docs/02_technical_handoff/18_Provider_Adapter_Technical_Spec.md`
- Runtime and error taxonomy: `docs/04_architecture_deep_dive/40_Provider_Adapter_Runtime_And_Error_Taxonomy.md`
- Provisioning ops checklist: `docs/05_development_standards/58_Provisioning_Ops_Readiness_Checklist.md`
- Billing operations runbook: `docs/05_development_standards/57_Billing_Operations_Runbook.md`

Do not paste real provider credentials, production account IDs, production DSNs, raw API responses, or customer data into this checklist, task files, PRs, logs, or docs.

## Readiness Levels

| Level | Meaning | Allowed work | Not allowed |
|---|---|---|---|
| Local fake provider | Uses in-repo fake provider behavior and local seed data. | Unit tests, local smoke, UI/demo flows. | Claiming sandbox or production provider readiness. |
| Sandbox provider | Uses provider sandbox account, sandbox auth material, and non-production resources. | Adapter implementation, sandbox contract tests, retry/error mapping checks. | Production traffic, production credentials, real customer orders. |
| Production provider | Uses approved production account and production runbook. | Production launch only after separate go/no-go approval. | Bypassing sandbox evidence or rollback plan. |

## Required Intake

Before a provider task is `READY`, record these in the task or linked docs:

- provider name, sandbox account owner, and support contact;
- exact sandbox API base URL and docs version;
- product types supported in sandbox, such as VPS or proxy;
- locations/regions available in sandbox;
- required auth method and least-privilege scopes;
- rate limits, quota limits, and concurrency limits;
- idempotency support for create, retry, cancel, and status read operations;
- provider timeout guidance and recommended client timeout;
- webhook availability, signature scheme, and replay behavior if webhooks are in scope;
- expected resource lifecycle states and cancellation behavior;
- provider error examples mapped to internal retry/error taxonomy;
- rollback path if sandbox provisioning creates unwanted resources.

Use display IDs in human notes when available. Use UUIDs only when an API path requires them.

## Credential Handling

Sandbox auth material must:

- live in a local `.env`, approved secret store, or CI secret when CI is explicitly in scope;
- be scoped to sandbox only;
- be rotated or revoked after shared testing if the provider supports it;
- never be committed, pasted into PR descriptions, or printed in logs;
- never be returned by API responses, frontend mock data, browser smoke output, or task files.

Any provider SDK/client must redact these tokens before logging:

- access tokens;
- API keys;
- account secrets;
- request signatures;
- raw auth headers;
- raw provider request or response bodies.

## Capability Mapping

Before coding adapter behavior, map provider capabilities to Billing fields:

- product type;
- plan code or provider SKU;
- CPU/RAM/disk/bandwidth shape if applicable;
- location;
- inventory mode;
- automatic provisioning support;
- manual fallback support;
- cancellation support;
- resize/upgrade support if relevant;
- expected `provider_resource_id` or external resource id shape.

The adapter must fail closed when capability mapping is missing or ambiguous. Do not guess a provider SKU from UI copy.

## Idempotency And Retry

Document provider behavior for:

- duplicate create requests with the same idempotency key;
- duplicate create requests without a key;
- network timeout after create request was accepted;
- provider 429/rate limit response;
- provider 5xx response;
- provider validation error;
- retry after manual review;
- cancel after partial create.

Map each case to one of:

- safe retry;
- manual review required;
- terminal failure;
- do not retry.

If provider idempotency is weak or missing, adapter work must include an external lookup/reconciliation plan before retrying create.

## Local Contract Harness

Before using real sandbox credentials, run the local provider contract harness:

```bash
go test ./internal/modules/provider -run SandboxContract
```

The current adapter interface maps the sandbox contract cases as follows:

- quote/readiness: `CheckStock`
- order/create: `Provision`
- status read: `GetStatus`
- cancel/cleanup: `Terminate`
- idempotency: repeated `Provision` with the same idempotency key

The harness uses the in-repo fake adapter by default and does not require provider network access or credentials. Real provider sandbox tests should reuse the same cases after credentials are stored outside git.

## Error And Timeout Contract

For sandbox evidence, capture redacted examples for:

- auth denied;
- permission denied;
- rate limited;
- invalid plan/SKU;
- out of capacity;
- duplicate request;
- timeout;
- provider internal error;
- resource not found;
- cancel rejected.

Each example must map to an internal provider error code, retry safety, and operator action. Store only redacted samples.

## Inventory And Readiness

Before running provisioning smoke against sandbox:

- confirm provider source status is active;
- confirm plan/source readiness is `ready`;
- confirm the sandbox location has enough quota/capacity;
- confirm product type and plan code map to provider sandbox SKU;
- confirm provider account has permission to create and cancel the target resource;
- confirm no production tenant, order, wallet, invoice, or credential is involved.

If any item is unknown, the provider is not sandbox-ready.

## Audit And Observability

Provider sandbox work must produce:

- audit log for operator recovery actions;
- redacted error message for failed job attempts;
- correlation ID through API, worker, and provider adapter logs;
- job status transition evidence;
- retry decision evidence;
- no raw provider payloads in frontend, API response, logs, or task notes.

## Rollback And Cleanup

Before first sandbox create call, document:

- how to list resources created by the test run;
- how to cancel/delete each resource;
- who owns manual cleanup if API cleanup fails;
- maximum sandbox spend or quota impact;
- stop condition for repeated failures;
- how to disable the provider source after a bad run.

Rollback must not require editing production data or wallet/payment records.

## Ready Decision

Provider sandbox work is ready only when:

- required intake is complete;
- credentials are stored outside git and scoped to sandbox;
- capability mapping is explicit;
- retry/idempotency behavior is known;
- redaction expectations are documented;
- rollback and cleanup path is documented;
- local fake provider smoke still passes;
- a reviewer can explain why sandbox work cannot affect production.

If any item is missing, keep the task `BLOCKED` or `TODO` and ask for the missing provider contract details.
