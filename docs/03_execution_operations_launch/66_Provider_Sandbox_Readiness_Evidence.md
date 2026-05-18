# Provider Sandbox Readiness Evidence

**Tasks:** T199, T208, T211, T212, T213, T214, T215, T216, T217, T218, T219, T220, T221, T226, T227, T228, T229, T230, T231, T232, T233, T235, T236, T237, T238, T239, T240, T241, T245, T247, T248, T249, T250, T252, T253, T254, T255, T256, T257, T258
**Date:** 2026-05-18
**Decision:** real provider sandbox is not launch-ready yet. Cloudmini V3 intake, authenticated read-only reachability, guarded dev mapping, source-match fail-closed behavior, one controlled dev create/delete pilot, fail-closed non-usable provider status handling, lifecycle-worker provider cleanup code, target test-server deploy/build evidence, non-mutating Cloudmini worker registry activation, bounded usable-status wait behavior, one target lifecycle-worker cleanup activation, target top-up review E2E evidence, target API session/RBAC denial evidence, target credential reveal audit evidence, target finance reconciliation read evidence, target dev/test duplicate-create plus timeout-after-send evidence, safe redacted provider error examples, usable-status owner sign-off, cleanup owner/procedure evidence, and selected-host self-managed secret-store proof are proven for the selected pilot semantics. T239 resolved the dev/test wallet projection mismatch found by T238 and the target finance reconciliation smoke now returns `balanced`. T240 removed target cloudflared token flag exposure from process argv and verified restricted dev/test secret file modes. T241 records `Admin` as the single accountable launch owner for all launch-day roles. T245 records Admin Finance/Security sign-off for completed target evidence. T247 reverified target dev/test secret file metadata and cloudflared token-file handling without reading or printing secret contents. T249 captured live Cloudmini duplicate-create and timeout-after-send evidence on the approved target dev/test provider account with same-session cleanup. T250 captured auth, not-found, and malformed-validation error metadata without raw response bodies. T252 records Admin sign-off for the fail-closed usable-status policy and bounded wait/read semantics. T253 records Admin as cleanup owner with the required cleanup/rollback procedure. T254 records owner-confirmed key rotation and selected-host self-managed secret-store proof without printing secret contents. T255 records source-inspection evidence and a safe execution plan for the remaining provider-controlled error cases. T256 captured the permission-denied runtime case with one temporary low-scope key, same-run revoke, and active key count restoration. T257 captured the out-of-capacity runtime case with one exhausted-group reservation probe and no reservation created. T258 adds Billing support for a one-request rate-limit fixture but confirms the inspected provider source has no safe V3 fixture to run yet. Rate-limit, 5xx, cancel/delete-rejected runtime examples and broader owner approval are still missing before broader pilot or production-like provisioning.

## Scope

This record separates local fake-provider evidence from real sandbox-provider readiness. Local fake evidence is useful for CI and developer smoke tests, but it is not approval to provision real sandbox or production resources.

## Current Readiness

| Target | Current status | Evidence | Pilot decision |
|---|---|---|---|
| VPS local fake path | `ready` for local validation only | Fresh seed maps `vps-cx23-40gb-monthly` to `Local Fake Hetzner Ready`; provider contract tests cover fake Hetzner create, status, terminate, idempotency, and timeout mappings. | OK for local/CI smoke only. |
| Proxy/manual local path | documented non-ready state | Fresh seed maps `proxy-static-10gb-monthly` first to an unsupported VPS-style source and has a manual fallback path. This proves the readiness API surfaces the gap instead of silently treating proxy as ready. | Not ready for real proxy sandbox provisioning. |
| Real provider sandbox | `blocked for broader pilot` | Cloudmini V3 non-production base URL and API version are known. A 2026-05-16 read-only rerun through `https://cz.resvn.net/` reached the V3 app: unauthenticated capabilities returned HTTP `401`, and authenticated capabilities plus inventory returned HTTP `200` V3 success envelopes using bearer, `X-API-Key`, and `X-ACCESS-CODE`. T218 selected a redacted `ipv4_dc` pilot mapping candidate from sellable read-only inventory and defines quota/cleanup approval requirements. T220 applied the guarded Cloudmini pilot mapping on the approved Billing dev runtime env and T221 evidence passed with plan display `10002`, plan-source display `10024`, source display `10012`, source type `cloudmini_v3`, readiness `ready`, priority `1`, and first-pilot guardrails of one create, one active resource, and one worker concurrency. T228 ran one controlled dev Billing checkout/payment/provisioning worker pilot with one provider resource, encrypted credential storage, provider delete cleanup, and Billing service termination. T229 makes Cloudmini provisioning fail closed when the provider resource status is not usable and makes lifecycle-worker termination call provider delete before marking a service terminated. T230 deployed that code to the approved test server and verified focused tests, builds, services, and health checks without mutating provider routes. T231 added and ran a non-mutating worker registry activation check with `PROVIDER_DEFAULT_MODE=cloudmini_v3`, confirming a real Cloudmini adapter can boot from the protected dev credential path without claiming jobs or calling provider APIs. T232 first hit manual review on status `creating` and cleaned up by fallback. T233 added bounded status polling, was deployed to the target server, and the owner-approved rerun created service display `10002`, stored one encrypted active credential, ran `lifecycle-once` for one terminate job, and verified provider final `GET` returned HTTP `404`. T235 separately proved target top-up review E2E without provider side effects. T236 proved target API session/RBAC denial behavior without provider side effects. T237 proved target credential reveal audit/redaction behavior without provider or money mutation routes. T245 records Admin sign-off for completed target finance/security evidence. T247 reverified the target dev/test credential and app secret file metadata without reading contents. T249 proved duplicate-create and timeout-after-send behavior on the approved target dev/test provider account with raw cleanup refs kept outside git and cleanup success for both scenarios. T250 proved redacted error metadata for auth missing, auth invalid, proxy not found, and malformed create validation. T252 records Admin sign-off for the usable-status semantics in doc 74. T253 records Admin cleanup ownership and rollback procedure in doc 75. T254 records selected-host self-managed secret-store proof in doc 76. T255 records provider source-inspection and a safe evidence plan in doc 77. T256 proves the provider-controlled permission-denied runtime case with redacted `403` metadata, same-run temporary key revoke, and active key count restoration. T257 proves the out-of-capacity runtime case with redacted `409 CAPACITY_EXHAUSTED` metadata, one exhausted-group reservation probe, and no reservation created. | Do not broaden pilot provisioning until the remaining provider-controlled runtime examples and broader provider approval are complete. |

## Proxy Cloudmini API V3 Candidate

T211 inspected the local `/opt/proxy-cloudmini` source code and added a Billing adapter for its API V3 contract using local `httptest` coverage only. T212 added disabled-by-default worker registry wiring behind explicit environment config. T213 recorded partial non-production intake for the Cloudmini V3 API. T214 attempted authenticated read-only provider checks and found an edge/gateway access blocker. T216 reran read-only checks with the local dev credential source and reached successful V3 app envelopes through the public hostname. This is still not real sandbox pilot evidence.

Known Cloudmini V3 intake as of 2026-05-18:

- Provider/API candidate: Cloudmini V3.
- Non-production API base URL: `https://cz.resvn.net/`.
- API version: V3.
- Auth boundary check: unauthenticated `GET /api/v3/capabilities` returned HTTP `401` with app JSON in `711ms`, confirming the public hostname routes to the Cloudmini manager and keeps auth enforced.
- Authenticated read-only rerun: T216 used the local dev credential source at `/opt/cred` without printing the raw key. The run used the Billing Go-client-style user-agent plus `X-Request-ID`; the checker was limited to read-only endpoints.
- Header forwarding result: bearer `Authorization`, `X-API-Key`, and `X-ACCESS-CODE` each returned HTTP `200` V3 success envelopes for `GET /api/v3/capabilities`, `GET /api/v3/inventory/groups?kind=ipv4_dc`, and `GET /api/v3/inventory/groups?kind=residential`.
- Capability summary: feature keys returned were `inventory_webhooks`, `prefer_wait`, `reservations`, and `tombstones`.
- Inventory summary: `ipv4_dc` returned `2` groups, with `1` sellable and `1` exhausted, totaling `200` allocatable units. `residential` returned `4` groups, all exhausted, totaling `0` allocatable units.
- Controlled pilot mapping candidate: T218 selected the sellable `ipv4_dc` group as `redacted:c6a7189f0a` with `200` allocatable units and `socks5` protocol. The raw group id is stored only in `/opt/cred-cloudmini-dev.env`.
- Edge note: provider-side evidence reported Cloudflare still blocks the generic `Python-urllib/3.12` user-agent with HTTP `403` code `1010`. The Billing Go-client-style path passed, so launch evidence should use the Billing adapter or the provider checker user-agent override, not a generic scripting user-agent.
- Credential status: credential material must stay outside git, task notes, PR text, logs, and raw command output. The local dev provider credential was split into `/opt/cred-cloudmini-dev.env` with mode `0600`; T247 reverified this target dev/test path as a protected local secret file. T254 promotes the rotated Cloudmini credential to the selected-host canonical path `/etc/billing/secrets/cloudmini.env` with mode `0600`, owner `root:root`, required Cloudmini keys present, and no `DB_DSN` present.
- Multi-endpoint status: T217 adds runtime support for multiple Cloudmini V3 endpoint/API-key mappings through `CLOUDMINI_V3_MAPPINGS_JSON`, keyed by provider source and optionally provider account. T220 covers only the single dev pilot source mapping; no multi-account target mapping or real provider pilot was run.
- Catalog mapping status: T219 adds `migrations/0025_add_cloudmini_provider_type.sql` and `scripts/cloudmini_pilot_mapping.sh` so an approved non-production DB can create the pilot `cloudmini_v3` provider source and plan-source mapping. T220 ran it on the approved Billing dev runtime env with `APP_ENV=dev`; migration plan showed `0` pending migrations and migration apply reported `0` applied migrations.
- Mapping evidence collector status: T221 adds `scripts/cloudmini_mapping_evidence.sh` so an operator can verify the applied mapping on an approved non-production Billing DB without sharing DSNs, tokens, raw group IDs, or raw provider payloads in repo evidence. T220 ran it read-only and recorded `result=PASS`, readiness `ready`, source type `cloudmini_v3`, and redacted guardrails only.
- Mutating pilot preflight status: T226 adds `scripts/cloudmini_mutating_pilot_preflight.sh` to fail closed unless non-production env, mapping evidence, owner fields, cleanup fields, private credential path, and exact one-resource guardrails are present. The script does not call provider mutating routes.
- Runtime source-match status: T227 makes Cloudmini runtime selection fail closed when an operation carries a Billing provider source ID that is not explicitly configured, so an account-level endpoint cannot bypass a source mismatch.
- Target DB access status: T220 used the approved test server Billing runtime env at `/opt/Billing/.env.dev`; the run confirmed `APP_ENV=dev` and `DB_DSN` presence without printing the DSN or provider secrets.
- Pilot status: T228 ran one controlled dev Billing-path create/delete pilot. The one-off worker created one Cloudmini `ipv4_dc` resource through `PROVIDER_DEFAULT_MODE=cloudmini_v3`, stored one encrypted credential with a masked hint, cleaned up the provider resource through V3 `DELETE`, and marked the Billing service terminated. Evidence uses display IDs and redacted hashes only.
- Hardening status: T229 changes Cloudmini provisioning so operation success is not enough by itself. Billing only treats Cloudmini proxy statuses `running`, `active`, `ready`, or `available` as usable for service activation; `creating` and other non-usable statuses go to manual review instead of active service creation. T229 also wires lifecycle-worker termination through provider `Terminate` when a provider registry is configured; provider cleanup timeout/unknown blocks the lifecycle transition and moves the job to manual review.
- Target deploy status: T230 synced the T229 code to the approved Billing dev test server at `/opt/Billing`, preserving local env and credential files. The target server confirmed the T229 source markers, passed `go test ./internal/modules/provider ./internal/modules/order ./cmd/worker`, `go run ./cmd/taskguard`, `go build -o bin/api ./cmd/api`, `go build -o bin/worker ./cmd/worker`, and `npm --prefix frontend run build`. `billing-api`, `billing-worker`, `billing-frontend`, and `cloudflared` were active after restart; `/healthz`, `/readyz`, and local frontend returned HTTP `200`; ports `8080` and `3000` were listening. No Cloudmini create/delete/action route, Billing checkout, payment, provisioning worker mutation, raw DSN, raw token, raw group id, raw provider payload, or proxy credential was printed or recorded. The always-on worker remained on `PROVIDER_DEFAULT_MODE=fake`; the protected dev Cloudmini credential file remained outside git at `/opt/cred-cloudmini-dev.env` with mode `0600`.
- Runtime activation preflight status: T231 added `cmd/worker provider-registry-check`, which builds the worker provider registry from env without opening the DB, claiming jobs, or calling provider APIs. The approved test server ran it with `APP_ENV=dev`, `PROVIDER_DEFAULT_MODE=cloudmini_v3`, `.env.dev`, and `/opt/cred-cloudmini-dev.env`; output was `result=PASS`, `cloudmini_v3_adapter=real`, one source mapping, zero account mappings, `provider_api_called=no`, `mutating_routes_called=no`, `jobs_claimed=0`, and `secrets_printed=no`. The credential file remained mode `0600`, owner `root:root`; no raw DSN, token, group id, source id, provider payload, or proxy credential was printed.
- Lifecycle activation attempt status: T232 ran the approved dev mutating preflight, stopped the always-on fake worker, created one Billing dev order/invoice/payment, and ran `cmd/worker provision-once` with `PROVIDER_DEFAULT_MODE=cloudmini_v3` and batch size `1`. The worker claimed exactly one job and returned manual review with `PROVIDER_PARTIAL_SUCCESS` because the Cloudmini resource status was `creating`. No active Billing service was created, so lifecycle-worker cleanup could not be run without bypassing T229. The resource was found by provider `external_ref`, direct V3 cleanup reached `succeeded`, and final provider `GET /api/v3/proxies/:id` returned HTTP `404`. The always-on worker was restarted in fake mode; no raw DSN, token, provider id, provider payload, or proxy credential was printed.
- Bounded status wait status: T233 changes Cloudmini provisioning so after create operation success, Billing polls `GET /api/v3/proxies/:id` until status is `running`, `active`, `ready`, or `available`, or until the configured timeout expires. A resource that stays `creating` still returns manual review and does not create an active service. A resource that becomes usable but lacks credential fields still fails closed with credential-missing handling.
- Usable-status sign-off status: T252 records Admin sign-off for the selected pilot semantics in `docs/03_execution_operations_launch/74_Cloudmini_Usable_Status_Signoff.md`. The approved policy keeps only `running`, `active`, `ready`, and `available` as service-activation statuses; empty, unknown, pending, `creating`, unrecognized, timeout, or credential-missing cases remain fail-closed/manual-review outcomes.
- Cleanup owner status: T253 records Admin as cleanup owner for selected Cloudmini pilot runs in `docs/03_execution_operations_launch/75_Cloudmini_Cleanup_Owner_And_Rollback.md`. The approved cleanup hierarchy prefers lifecycle-worker provider-backed cleanup, allows direct Cloudmini V3 delete only as an owner-approved dev/test fallback, and requires source disable plus incident/follow-up before any further create attempt if cleanup cannot be confirmed.
- Target lifecycle activation status: T233 was deployed to the approved test server and rerun in a one-resource owner-approved activation window. The run used existing dev wallet balance, stopped the always-on fake worker, created order display `10003`, invoice display `10004`, payment display `10003`, and provider job display `10003`, then ran `provision-once` with `PROVIDER_DEFAULT_MODE=cloudmini_v3`. Provisioning succeeded, service display `10002` became active/paid, one active encrypted credential with masked hint existed, one lifecycle terminate job display `10004` succeeded through `lifecycle-once`, the service ended `terminated`, and provider final `GET` returned HTTP `404`. The always-on worker was restored active. No raw DSN, token, provider id, provider payload, or proxy credential was printed.
- Target top-up review status: T235 was deployed to the approved test server and ran `dev-topup-review` against the local API and DB. The run used a temporary dev/test wallet, approved top-up display `10003` with ledger display `10005` and audit display `10015`, rejected top-up display `10004` with no ledger and audit display `10016`, recorded wallet delta `111`, and observed no order, provider job, service, or provider-resource side effects.
- Target auth/RBAC status: T236 was deployed to the approved test server and ran `dev-target-auth-rbac` against the local API. The run proved cookie-only client session access without `X-Actor-*`, admin 2FA enforcement with `auth.2fa_required`, invalid session denial with `auth.session_invalid`, missing actor denial with `auth.actor_required`, cross-tenant denial with `tenant.context_mismatch`, and three low-permission RBAC denials with `auth.permission_denied`. No raw session token, cookie, password, DSN, provider payload, or credential was printed or recorded.
- Target credential reveal status: T237 was deployed to the approved test server and ran `dev-target-credential-reveal` against the local API and DB. The run created/refreshed one encrypted dev/test credential fixture for service display `43001`, revealed it through client session cookie-only auth, verified `no-store` response headers, verified `last_revealed_by`, reveal rate-limit state, and audit display `10017`, and did not call provider or money mutation routes. No plaintext credential, encrypted payload, raw credential ID, session token, cookie, DSN, provider payload, or provider credential was printed or recorded.
- Target finance reconciliation status: T238 was deployed to the approved test server and ran `dev-target-finance-reconciliation` against the local API and DB. The read-only run verified payment reconciliation list/detail for transaction display `51001`, invoice display `44001`, wallet display `41001`, and ledger display `50002`, but daily reconciliation for `2026-04-23` returned `mismatched` with one wallet mismatch. T239 traced the root cause to dev/test wallet projection drift from an inconsistent seed baseline and later smoke runs, fixed the seed baseline to `3600`, repaired the approved target dev/test projection from posted ledger source-of-truth with audit display `10018`, and reran the target smoke. The rerun returned `balanced` with wallets checked `2`, wallet mismatches `0`, invoices checked `1`, invoice mismatches `0`, payments checked `1`, and duplicate payment references `0`. No ledger rows were inserted or updated by the projection repair, and no money or provider mutation routes were called by the smoke. T245 records Admin Finance sign-off for this evidence.
- Target secret/key handling status: T240 verified target dev/test secret file modes without printing contents: `/opt/Billing/.env.dev` mode `640` owner `root:billing-svc`, `/opt/cred-cloudmini-dev.env` mode `600` owner `root:root`, and `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`. Cloudflared was changed from token flag usage to `--token-file /etc/cloudflared/tunnel.token`, restarted active, and verified with no token present in process arguments. `http://localhost:3000`, `https://billing.resvn.net`, `https://reseller.resvn.net`, and `https://client.resvn.net` returned HTTP `200`. T245 records Admin Security sign-off for this target dev/test evidence; T254 later records the selected-host self-managed secret-store proof.
- Target dev/test secret-store recheck: T247 rechecked metadata on the approved target server without reading file contents. `/opt/Billing/.env.dev` remained mode `640` owner `root:billing-svc`, `/opt/cred-cloudmini-dev.env` remained mode `600` owner `root:root`, and `/etc/cloudflared/tunnel.token` remained mode `600` owner `root:root`. `cloudflared` was active, used `--token-file`, and had no token flag in process arguments. `billing-api` and `billing-worker` were active. No provider create/delete/action route, raw DSN, raw token, API key, file contents, provider payload, cookie, or proxy credential was printed or recorded. This confirms the target dev/test local secret-file boundary; T254 later records selected-host self-managed secret-store proof after owner-confirmed key rotation.
- Self-managed secret-store status: T254 records owner-confirmed API key rotation and selected-host local-only secret-store proof in `docs/03_execution_operations_launch/76_Self_Managed_Secret_Store_Rotation_Evidence.md`. The canonical Cloudmini path is `/etc/billing/secrets/cloudmini.env` with directory mode `700` owner `root:root`, file mode `600` owner `root:root`, required Cloudmini keys present, and `DB_DSN` absent. Cloudflared uses `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`; the running process has `--token-file` and no exact `--token` arg. No secret value or file content was printed.
- Idempotency evidence status: T248 adds `go run ./cmd/smoke cloudmini-idempotency-evidence` for owner-approved non-production `duplicate-create` and `timeout-after-send` evidence. T249 ran both scenarios on the approved target dev/test provider account. Duplicate-create used two create attempts with the same idempotency key and returned one distinct redacted resource, `duplicate_same_resource=true`, and cleanup success. Timeout-after-send used one create attempt with a forced short poll timeout and returned `PROVIDER_TIMEOUT_REQUEST_KNOWN`, `manual_review_required`, and cleanup success. Raw cleanup references remained outside git with mode `0600`; stdout omitted raw DSNs, tokens, provider IDs, provider payloads, proxy credentials, and cookies.
- Error evidence status: T250 adds `go run ./cmd/smoke cloudmini-error-evidence` for owner-approved non-production Cloudmini error evidence. The approved target dev/test run returned `PASS` with four redacted examples: missing auth `401` mapped to `PROVIDER_AUTH_FAILED`; invalid auth `401` mapped to `PROVIDER_AUTH_FAILED`; a valid-format missing proxy UUID returned `404 PROXY_NOT_FOUND` mapped to `PROVIDER_STATE_DRIFT` with `manual_review_required`; malformed create JSON returned `400 INVALID_INPUT` mapped to `PROVIDER_CONFIG_INVALID`. The malformed create case was explicitly approved, called one mutating route with invalid JSON, and did not print or record any raw body, token, provider ID, provider payload, proxy credential, cookie, or file contents. T255 records source-inspection evidence and a safe execution plan for the remaining provider-controlled cases in `docs/03_execution_operations_launch/77_Cloudmini_Provider_Controlled_Error_Evidence.md`. T256 extends the smoke with guarded permission-denied collection and the approved dev/test run returned `403`, `PROVIDER_PERMISSION_DENIED`, `do_not_retry`, temporary API key revoked, active key count restored, and `mutating_routes_called=true` for the API-key create/revoke support routes. T257 extends the smoke with guarded out-of-capacity collection and the approved dev/test run returned `409 CAPACITY_EXHAUSTED`, `PROVIDER_OUT_OF_STOCK`, `do_not_retry`, one exhausted-group reservation probe, and `reservation_created=false`. T258 extends the smoke with guarded rate-limit fixture support and local coverage for `429 RATE_LIMITED` to `PROVIDER_RATE_LIMITED` with `safe_retry`; no live provider rate-limit run was executed because the inspected provider source has no safe V3 fixture. Remaining provider-controlled runtime examples are rate limited, provider 5xx, and cancel/delete rejected.

T250 target dev/test error evidence redacted stdout:

```text
cloudmini_error_evidence result=PASS
pilot_environment=dev
approval_fields_present=yes
owner_fields_present=yes
example_count=4
mutating_routes_called=true
example_1_name=auth_missing_capabilities
example_1_http_status=401
example_1_provider_error_code=none
example_1_normalized_error_code=PROVIDER_AUTH_FAILED
example_1_retry_safety=do_not_retry
example_1_error_envelope_present=true
example_1_error_message_field_present=true
example_1_error_details_field_present=false
example_2_name=auth_invalid_capabilities
example_2_http_status=401
example_2_provider_error_code=none
example_2_normalized_error_code=PROVIDER_AUTH_FAILED
example_2_retry_safety=do_not_retry
example_2_error_envelope_present=true
example_2_error_message_field_present=true
example_2_error_details_field_present=false
example_3_name=not_found_proxy
example_3_http_status=404
example_3_provider_error_code=PROXY_NOT_FOUND
example_3_normalized_error_code=PROVIDER_STATE_DRIFT
example_3_retry_safety=manual_review_required
example_3_error_envelope_present=true
example_3_error_message_field_present=true
example_3_error_details_field_present=false
example_4_name=validation_malformed_create
example_4_http_status=400
example_4_provider_error_code=INVALID_INPUT
example_4_normalized_error_code=PROVIDER_CONFIG_INVALID
example_4_retry_safety=do_not_retry
example_4_error_envelope_present=true
example_4_error_message_field_present=true
example_4_error_details_field_present=false
raw_response_body_printed=no
sensitive_values_printed=no
raw_provider_ids_printed=no
provider_payloads_printed=no
remaining_provider_controlled_examples=permission_denied,rate_limited,out_of_capacity,provider_5xx,cancel_rejected
```

T256 target dev/test permission-denied evidence redacted stdout excerpt:

```text
cloudmini_error_evidence result=PASS
pilot_environment=dev
approval_fields_present=yes
owner_fields_present=yes
example_count=4
mutating_routes_called=true
example_4_name=permission_denied_proxy_list
example_4_http_status=403
example_4_provider_error_code=none
example_4_normalized_error_code=PROVIDER_PERMISSION_DENIED
example_4_retry_safety=do_not_retry
example_4_error_envelope_present=true
example_4_error_message_field_present=true
example_4_error_details_field_present=false
example_4_side_effect_created=cleaned_up
example_4_temporary_api_key_created=true
example_4_temporary_api_key_revoked=true
example_4_active_key_count_restored=true
raw_response_body_printed=no
sensitive_values_printed=no
raw_provider_ids_printed=no
provider_payloads_printed=no
remaining_provider_controlled_examples=rate_limited,out_of_capacity,provider_5xx,cancel_rejected
```

T257 target dev/test out-of-capacity evidence redacted stdout excerpt:

```text
cloudmini_error_evidence result=PASS
pilot_environment=dev
approval_fields_present=yes
owner_fields_present=yes
example_count=4
mutating_routes_called=true
example_4_name=out_of_capacity_reservation
example_4_http_status=409
example_4_provider_error_code=CAPACITY_EXHAUSTED
example_4_normalized_error_code=PROVIDER_OUT_OF_STOCK
example_4_retry_safety=do_not_retry
example_4_error_envelope_present=true
example_4_error_message_field_present=true
example_4_error_details_field_present=true
example_4_side_effect_created=no
example_4_reservation_probe_attempted=true
example_4_exhausted_group_selected=true
example_4_reservation_created=false
example_4_reservation_cleaned_up=false
example_4_reservation_max_attempts=1
example_4_reservation_ttl_seconds=60
raw_response_body_printed=no
sensitive_values_printed=no
raw_provider_ids_printed=no
provider_payloads_printed=no
```

## Cloudmini Edge/Gateway Unblock Runbook

T215 documents the provider-owner handoff needed for authenticated Cloudmini read-only checks. The unblock is outside Billing runtime code because T214 reached the public base URL but initially received provider edge/gateway HTTP `403` responses before a successful V3 app envelope. T216 shows the required read-only route/header path now works for the Billing Go-client-style path.

Provider owner must confirm these items before Billing reruns authenticated checks:

- Public hostname `https://cz.resvn.net/` routes `/api/v3/*` to the Cloudmini manager origin through the approved tunnel or gateway path.
- Edge/WAF/Access policy allows non-browser server-to-server clients for `/api/v3/capabilities` and `/api/v3/inventory/groups`.
- Edge policy allows the headers required by the code-read contract: `Authorization`, `X-API-Key`, `X-ACCESS-CODE`, `X-Request-ID`, and `Idempotency-Key`.
- If IP allowlisting is required, the provider owner records a redacted allowlist reference for the Billing runner egress IP or approved staging egress, not the raw credential.
- The provider API key is active, scoped to sandbox/non-production, and has read permission for capabilities and inventory plus later explicit `proxy_crud` permission only when create/delete pilot is approved.
- The credential shared through chat is rotated or explicitly accepted by Security/Provider owner as a temporary sandbox-only credential before reuse.
- Query-string credentials such as `?token=` or `?access_code=` are avoided for Billing evidence because they can leak through URL logs. Use headers unless a Security Owner signs a temporary exception.

Safe read-only rerun after unblock:

1. Store the rotated credential outside git in an approved secret path or local-only `.env` file.
2. Run only `GET /api/v3/capabilities` and `GET /api/v3/inventory/groups?kind=ipv4_dc|residential`.
3. Capture only status codes, envelope success, feature keys, inventory counts, sell-state counts, and redacted group references.
4. Do not capture raw provider response bodies, raw auth headers, provider-private IDs, proxy credentials, or URL query credentials.
5. Keep pilot readiness blocked unless both read-only checks return a successful V3 app envelope and owner/quota/mapping/cleanup evidence is recorded.

Although read-only evidence and the first dev pilot now pass, do not run these outside the controlled dev pilot boundary until the remaining provider readiness items are complete:

- `POST /api/v3/proxies`
- `DELETE /api/v3/proxies/:id`
- `POST /api/v3/proxies/:id/actions/:action`
- Billing checkout/provisioning worker pilot with `PROVIDER_DEFAULT_MODE=cloudmini_v3`

Code-read contract summary:

- Auth supports `Authorization: Bearer <token>` and API-key fallback headers in the provider service. Billing adapter uses bearer auth.
- Readiness/inventory endpoints: `GET /api/v3/capabilities`, `GET /api/v3/inventory/groups?kind=<kind>`.
- Supported proxy kinds from code: `ipv4_dc` and `residential`.
- Mutating V3 endpoints require `Idempotency-Key`.
- Create path: `POST /api/v3/proxies` returns `202 Accepted` with an async operation id and resource id.
- Status path: `GET /api/v3/operations/:id` polls `accepted/running/succeeded/failed/timed_out/cancelled`.
- Resource paths: `GET /api/v3/proxies/:id`, `DELETE /api/v3/proxies/:id`.
- Action path: `POST /api/v3/proxies/:id/actions/:action` supports `stop`, `start`, and residential-only `change-ip`.
- Credential-bearing proxy response fields are encrypted by Billing adapter tests before returning a provider `CredentialEnvelope`.
- Worker runtime stays on `PROVIDER_DEFAULT_MODE=fake` by default. `PROVIDER_DEFAULT_MODE=cloudmini_v3` requires Cloudmini base URL, API token, Billing source id, kind, group id, protocol, and `ENCRYPTION_KEY` before startup.

Still missing for broader real sandbox readiness:

- sandbox account owner and support contact for broader provider account proof;
- provider edge/gateway approval record for the read-only route/header policy;
- owner-approved source-to-group/SKU mapping beyond the dev pilot evidence;
- active Cloudmini V3 provider source and plan source readiness evidence for any additional target provider source;
- multi-endpoint/account config if more than one Cloudmini V3 URL or API key is needed;
- owner-approved timeout/spend guardrail sign-off beyond the dev defaults;
- provider-controlled redacted runtime examples for rate limited, provider 5xx, and cancel/delete rejected;
- broader provider-owner approval for production-like provisioning beyond the selected pilot semantics.

## Evidence Packet Status

One approved dev mutating pilot evidence packet is stored in this repository as of 2026-05-17. The packet below is still required before changing broader real provider sandbox readiness from `blocked for broader pilot`.

| Evidence area | Required proof | Current repo status |
|---|---|---|
| Provider intake | Provider name, sandbox account owner, support contact, docs/API version, and sandbox base URL. | Partial: Cloudmini V3, API V3, and `https://cz.resvn.net/` are recorded. Sandbox account owner and support contact are missing. |
| Credential safety | Approved secret store or local-only `.env` path, least-privilege scopes, rotation/revocation owner, and confirmation that no secret is committed. | Pass for selected host/pilot scope: no credential is committed in repo evidence. T216 used `/opt/cred` as a local dev-only source without printing the raw key. T240/T247 verify target dev/test secret-file metadata and cloudflared token-file handling. T254 records owner-confirmed key rotation, canonical local-only Cloudmini path `/etc/billing/secrets/cloudmini.env` mode `0600` owner `root:root`, required Cloudmini keys present, `DB_DSN` absent, cloudflared token-file mode `0600`, and no exact `--token` arg in the running process. Repeat this evidence for any new host or secret path. |
| Quota and cost guardrail | Sandbox quota, rate/concurrency limits, maximum spend or credit exposure, and stop condition. | Partial: T220 recorded first-pilot dev guardrails of one create, one active resource, one worker concurrency, no parallel mutating calls, and single-dev-resource exposure. T226 adds a preflight guard that requires these fields, but owner sign-off, sandbox quota, and stop-condition approval are still missing. |
| Capability mapping | Product type, Billing plan code, provider SKU, location, inventory mode, auto/manual provisioning support, cancellation support, and credential retrieval behavior. | Partial: authenticated read-only inventory succeeds and shows `ipv4_dc` has sellable capacity while `residential` is exhausted. T220 applied dev pilot mapping for `proxy-static-10gb-monthly` to `cloudmini_v3` with readiness `ready`; owner-approved SKU/source mapping and cleanup/cancellation behavior are still missing. |
| Retry/idempotency | Duplicate create behavior, timeout-after-send behavior, request/status lookup support, and mapping to retry safety or manual review. | Pass for approved target dev/test: repo tests cover request-known timeout/manual review and non-usable status/manual review. T249 ran duplicate-create against Cloudmini with two create attempts, one distinct redacted resource, duplicate protection, and cleanup success. T249 ran timeout-after-send with `PROVIDER_TIMEOUT_REQUEST_KNOWN`, `manual_review_required`, one redacted resource, and cleanup success. Broader provider approval remains a separate blocker. |
| Error examples | Redacted auth, permission, rate limit, validation, out-of-capacity, timeout, duplicate, 5xx, not-found, and cancel-rejected examples. | Partial: redacted provider edge/gateway HTTP `403` shape from the earlier blocked run, app-level unauthenticated HTTP `401`, T249 duplicate/timeout behavior, T250 auth missing, auth invalid, not-found, and malformed-validation examples, T256 provider-controlled permission-denied runtime output, and T257 provider-controlled out-of-capacity runtime output are captured. T255 records the provider source-inspection and safe execution plan for the remaining cases. T258 adds Billing runner support for a future rate-limit fixture but records no live provider rate-limit evidence. Provider-controlled rate limit, 5xx, and cancel/delete rejected runtime outputs are still missing. |
| Cleanup and rollback | How to list test resources, cancel/delete them, disable the provider source, and assign manual cleanup owner. | Pass for selected pilot scope: T228 cleaned up the single dev resource through provider V3 `DELETE` and marked the Billing service terminated. T229 adds provider-backed lifecycle-worker termination before service `terminated`; T230 proves the hardened code is deployed and builds on the target test server. T231 proves non-mutating worker registry activation with the real Cloudmini adapter. T232 proved same-session fallback cleanup after a non-usable `creating` status. T233 target rerun proved lifecycle-worker provider cleanup on one Cloudmini resource. T249 proved cleanup success for duplicate-create and timeout-after-send scenarios. T253 assigns Admin as cleanup owner and records the rollback procedure. Each future approved run still must record run-specific cleanup result evidence. |
| Pilot run | Redacted evidence for one approved sandbox order through checkout, reservation, provider request, service activation, credential storage/reveal audit, and cleanup. | Partial/pass for controlled dev: T228 created one order/job/service through Billing path, stored encrypted credential metadata, and completed same-session provider cleanup. T233 target rerun created one additional order/job/service, stored one encrypted active credential, and completed lifecycle-worker provider cleanup. T235 proved target top-up review E2E using a temporary dev/test wallet: approve posted one ledger credit and audit row, reject posted no ledger and one audit row, and provider side effects were `none`. T237 proved target credential reveal audit/redaction on a dev/test fixture without provider or money mutation routes. T238 proved target finance reconciliation read paths and T239 resolved the dev/test projection mismatch with audit display `10018`; the target finance smoke now returns `balanced`. T249 proved duplicate-create and timeout-after-send provider behavior with same-session cleanup. Broader owner sign-off remains missing. |

## Evidence Packet Template

Use this template in a task, runbook appendix, or external launch packet. Store only redacted values in git.

```text
Provider:
Sandbox account owner:
Provider support contact:
Sandbox API base URL:
Provider docs/API version:
Approved credential storage path:
Credential scope:
Credential rotation/revocation owner:
Quota/rate/concurrency limits:
Maximum sandbox spend or quota exposure:
Stop condition:
Billing plan code:
Provider SKU:
Sandbox location:
Inventory mode:
Auto provisioning supported:
Manual fallback supported:
Credential retrieval behavior:
Cancellation/cleanup supported:
Provider idempotency level:
Duplicate create behavior:
Timeout-after-send behavior:
Safe status lookup path:
Retry safety map:
Redacted error examples captured:
Pilot resource cleanup owner:
Pilot run date:
Pilot run reviewer:
```

## Required Pilot Run Evidence

Before approving real sandbox provisioning, capture redacted evidence for all of these:

- Source readiness is `ready` for the exact plan/source being tested.
- Checkout debits wallet and creates a single reservation and provisioning job.
- Provider create uses the job idempotency key or a documented equivalent.
- Success creates one active service and one provider resource mapping.
- Duplicate retry does not create a second provider resource.
- Timeout-after-send moves to manual review or safe status lookup, not blind retry.
- Provider credentials are stored encrypted and not printed in logs, provider request records, audit, task notes, or PR text.
- Credential reveal remains a separate audited action.
- Cleanup/cancel removes or disables the sandbox resource and records the cleanup owner.
- Any failed case maps to an internal provider error code and retry safety.

## Redaction Rules

Do not commit or paste:

- API keys, bearer tokens, signatures, cookies, private keys, passwords, root credentials, proxy credentials, or raw auth headers.
- Production provider account IDs, production DSNs, production customer data, or real customer order data.
- Raw provider request or response bodies if they contain secrets, customer data, or provider-private identifiers.

Use stable display IDs or redacted placeholders when a human-readable reference is needed.

## Verification Evidence

Provider contract expectations are covered locally by:

```bash
go test ./internal/modules/provider -run SandboxContract
```

Provisioning worker timeout safety is covered by:

```bash
go test ./internal/modules/order -run ProviderProvisioningHandler
```

The required local billing smoke before pilot remains:

```bash
go run ./cmd/smoke -dsn "$DB_DSN" -base-url "$API_BASE_URL" dev-billing
```

That smoke requires a running API and a non-production database. Do not treat unit tests as a substitute for that smoke before pilot.

## No-Go Until Fixed

Before changing this decision to ready, record the required provider intake from `docs/05_development_standards/60_Provider_Sandbox_Contract_Checklist.md`:

- sandbox provider name, account owner, and support contact;
- sandbox API base URL and docs version;
- sandbox-only credential storage outside git;
- supported VPS/proxy product types, locations, quota, and rate limits;
- idempotency behavior for duplicate create and timeout-after-create;
- redacted examples for auth, rate limit, validation, timeout, and provider 5xx errors;
- cleanup and rollback plan for resources created by sandbox tests.

If any item is missing, keep real provider sandbox readiness blocked.
