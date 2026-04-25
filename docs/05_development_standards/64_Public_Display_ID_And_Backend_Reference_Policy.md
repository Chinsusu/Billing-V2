# Public Display ID and Backend Reference Policy

**Scope:** API contracts, frontend labels, table filters, support notes, and PR review for IDs.

## Purpose

People need short numeric IDs in the UI for lookup, support, and operations. The backend still needs stable UUID or opaque references for writes, relations, and audit safety.

Use these simple terms:

- `public ID`: numeric ID shown to people, usually stored as `display_id` or returned as a related `*_display_id`.
- `backend reference`: UUID or opaque ID used by code, joins, API paths, and action bodies.
- `external reference`: provider or payment processor ID. It is not the main UI label.

## Base Rules

1. User-facing tables, cards, page titles, exports, support notes, and PR screenshots must use public IDs.
2. Backend references may stay in API responses when the frontend needs them for links or actions.
3. The frontend must not use a UUID as a fallback label. If a public ID is missing, show `not shown` or add the missing public ID in the API.
4. Public IDs are for readability, not authorization. Tenant scope, RBAC, and permission checks still use backend rules.
5. Public IDs are not guaranteed to be gapless. Deleted or failed records can leave gaps.
6. Do not encode business meaning into a public ID. Treat it as an opaque number.
7. External references are secondary support detail only. Do not use provider IDs as the primary label.

## API Response Rule

For a resource that can appear in the UI, the response should include:

```text
id          backend reference for routes/actions
display_id  public ID for labels and table search
```

For related resources, prefer both references when the UI needs to show the relation:

```text
order_id           backend reference for detail links/actions
order_display_id   public ID for the visible relation label
```

If adding both would expose a sensitive relationship, expose the public ID only when the caller already has permission to see that related resource.

## Query Filter Rule

Use these names consistently:

| Filter name | Meaning | Example |
| --- | --- | --- |
| `display_id` | public ID of the resource being listed | `/admin/orders?display_id=10042` |
| `<resource>_display_id` | public ID of a related resource | `/admin/services?order_display_id=10042` |
| `<resource>_id` | backend reference used by actions, detail routes, or internal filters | `/admin/services?order_id=<uuid>` |

Public ID filters must accept only positive integers. Non-positive or non-numeric values return a validation error with a field-level code.

When both public and backend filters exist, frontend search boxes should default to the public filter. Backend reference filters can remain available for internal API clients, but they are not the normal UI path.

## Resource Field Policy

| Resource | Public label fields | Backend reference fields | UI rule |
| --- | --- | --- | --- |
| Tenant | `display_id`, future `tenant_display_id` when related | `tenant_id`, `parent_tenant_id` | Show public ID with tenant name/domain. |
| Account/user/customer | `display_id`, future `account_display_id`, `buyer_display_id`, `requested_by_display_id` when related | `user_id`, `buyer_user_id`, `account_user_id`, `requested_by`, `created_by`, `owner_id` | Show public ID plus email/name when allowed. Do not show raw user UUIDs. |
| Provider/source | `display_id`, `source_display_id`, `plan_source_display_id` | `source_id`, `provider_source_id`, `provider_account_id` | Show public source ID plus provider name/location. Never show provider credential or account references as labels. |
| Product/plan | `display_id`, `product_display_id`, `plan_display_id`, `tenant_plan_display_id` | `product_id`, `plan_id`, `tenant_product_id`, `tenant_plan_id` | Show public ID plus code/name. |
| Order | `display_id`, `order_display_id` | `order_id` | Show public order ID in all order-facing UI. |
| Service | `display_id`, `service_display_id`, `order_display_id`, `provider_source_display_id` | `service_id`, `service_instance_id`, `order_id`, `provider_source_id` | Show public service ID and related public IDs. External resource ID is support detail only. |
| Invoice | `display_id`, `invoice_display_id`, `order_display_id` | `invoice_id`, `order_id` | Show public invoice ID; use related public order ID in tables. |
| Wallet and ledger | `display_id`, `wallet_display_id`, `ledger_display_id` | `wallet_id`, `ledger_entry_id`, `reference_id` | Show public wallet or ledger ID. Do not label rows by backend reference. |
| Transaction | `display_id`, `transaction_display_id`, `invoice_display_id`, `order_display_id` | `transaction_id`, `payment_transaction_id`, `invoice_id`, `order_id` | Show public transaction ID and related public IDs. |
| Top-up | `display_id`, `topup_display_id`, `wallet_display_id` | `topup_request_id`, `wallet_id`, `requested_by` | Show public top-up ID and public wallet ID. |
| Job/outbox | `display_id`, `job_display_id`, `attempt_display_id`, related public IDs when available | `job_id`, `job_attempt_id`, `outbox_event_id`, `reference_id`, `source_id` | Show public job ID. Related backend references should be hidden unless clearly labeled for internal ops. |
| Audit log | `display_id`, `audit_display_id`, related actor/target public IDs when available | `audit_id`, `actor_id`, `target_id`, `correlation_id` | Show public audit ID. Actor/target UUIDs are not row labels. |

## Frontend Mapping Rule

Frontend API mapping should create view models with explicit label fields such as:

```text
publicLabel
orderPublicLabel
invoicePublicLabel
walletPublicLabel
```

These labels must come from public IDs or allowed business text, not from UUIDs. If the API does not return a required public ID, the frontend task should either add a backend/API follow-up or display `not shown`.

## PR Review Checklist

Before merging a route or UI that shows a resource:

- Does the visible label use a public ID?
- Does the response include related `*_display_id` fields when the table needs them?
- Are backend references still available for actions that need them?
- Are public ID filters named `display_id` or `<resource>_display_id`?
- Do tests cover positive integer validation for new public ID filters?
- Did the PR body mention any temporary `not shown` fields that need a follow-up task?

## Related Docs

- API reference: `docs/05_development_standards/56_Billing_API_Operational_Reference.md`
- Frontend standard: `docs/05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md`
- API contract: `docs/02_technical_handoff/16_API_Contract_And_Permission_Spec.md`
- Validation matrix: `docs/05_development_standards/63_Validation_Command_Matrix.md`
