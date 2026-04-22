# Changelog — Architecture Deep Dive Package v1.4

**Date:** 2026-04-22

## Added

- `04_architecture_deep_dive/37_Go_Backend_Architecture_And_Module_Boundaries.md`
- `04_architecture_deep_dive/38_PostgreSQL_Data_Consistency_Transaction_Design.md`
- `04_architecture_deep_dive/39_Async_Worker_Outbox_And_Job_Architecture.md`
- `04_architecture_deep_dive/40_Provider_Adapter_Runtime_And_Error_Taxonomy.md`
- `04_architecture_deep_dive/41_Tenant_Isolation_RBAC_And_Security_Architecture.md`
- `04_architecture_deep_dive/42_Secrets_Credential_Encryption_And_Audit_Architecture.md`
- `04_architecture_deep_dive/43_Observability_Logging_Metrics_Tracing_Spec.md`
- `04_architecture_deep_dive/44_Scaling_Performance_And_Failure_Mode_Architecture.md`
- `04_architecture_deep_dive/45_Architecture_Decision_Records_ADR.md`

## Updated

- `00_README.md` with the v1.4 architecture deep-dive index and reading order.
- `MANIFEST.txt` regenerated as v1.4 full package.
- `04_architecture_deep_dive/38_PostgreSQL_Data_Consistency_Transaction_Design.md` related-doc references updated to include the full architecture layer.
- `VPS_Proxy_Project_Master_Document.md` appended with the v1.4 architecture layer index.

## Intent

This release moves the project from product/technical handoff plus execution planning into implementation-grade architecture alignment:

```text
Go backend module boundaries
PostgreSQL transaction and consistency model
Outbox/job async architecture
Provider runtime error taxonomy
Tenant isolation and RBAC security
Secret/credential encryption and audit
Observability and incident signals
Scaling and failure mode strategy
Architecture Decision Records
```
