# T217 - Cloudmini V3 multi-endpoint config

Status: TODO
Owner: -
Branch: feat/cloudmini-v3-multi-endpoint-config
PR: -
Risk: provider/provisioning/credential/config
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Support multiple Cloudmini V3 endpoint/API-key mappings when different provider sources or provider accounts use different V3 base URLs and credentials.

## Scope

- Add runtime config that maps a provider source or provider account to its own Cloudmini V3 base URL, API credential reference, kind, group, node, protocol, and limits.
- Route Cloudmini V3 operations to the matching endpoint/key based on the provisioning operation source/account.
- Keep credential material outside git, logs, task notes, and raw command output.
- Preserve fail-closed behavior when a source/account mapping is missing or invalid.
- Do not model multiple endpoints by inventing provider types such as `cloudmini_v3_a` or `cloudmini_v3_b`.

## Acceptance Criteria

- A worker can provision through at least two Cloudmini V3 source/account mappings with different base URLs or API keys in tests.
- Missing source/account mapping returns a config error and does not call any provider endpoint.
- Secret values are loaded by approved secret reference or local-only environment input and are never logged.
- Existing single-endpoint Cloudmini V3 behavior remains backward compatible or has an explicit migration path.
- Relevant provider and worker tests pass.

## Notes

- Current limitation: `CloudminiV3Config` has one `BaseURL` and one `APIToken` per adapter, and the worker env currently populates one `CLOUDMINI_V3_SOURCE_ID` mapping only.
- `CloudminiV3Adapter` already supports per-source kind/group/node/protocol config, but not per-source endpoint/key.
- `provider.Registry` is keyed by provider type, so registering multiple `cloudmini_v3` adapters is not the right fix.

## Agent Log

- 2026-05-16: Task created as a follow-up note after provider evidence T216.
