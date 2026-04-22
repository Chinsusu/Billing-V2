# VPS/Proxy Billing Platform

This workspace is organized as the Go modular monolith described in the architecture docs.

## Directory Layout

```text
cmd/
  api/         HTTP API entrypoint
  worker/      async worker entrypoint
  scheduler/   recurring job entrypoint
  migrate/     migration runner entrypoint
  cli/         internal operator CLI

internal/
  app/         application wiring
  platform/    shared infrastructure packages
  modules/     business domain modules

migrations/    PostgreSQL migrations
scripts/       local/dev/ops scripts
docs/          project documentation package
tasks/         per-task coordination files for multi-agent work
```

## Documentation

Start with [docs/00_README.md](docs/00_README.md).

The full documentation manifest is in [docs/MANIFEST.txt](docs/MANIFEST.txt).

Task coordination starts with [TASKS.md](TASKS.md) and [tasks/README.md](tasks/README.md).

Architecture implementation baseline:

- [Go backend architecture](docs/04_architecture_deep_dive/37_Go_Backend_Architecture_And_Module_Boundaries.md)
- [PostgreSQL consistency and transactions](docs/04_architecture_deep_dive/38_PostgreSQL_Data_Consistency_Transaction_Design.md)
- [Async worker/outbox/jobs](docs/04_architecture_deep_dive/39_Async_Worker_Outbox_And_Job_Architecture.md)
- [Provider runtime/error taxonomy](docs/04_architecture_deep_dive/40_Provider_Adapter_Runtime_And_Error_Taxonomy.md)
- [Tenant/RBAC/security](docs/04_architecture_deep_dive/41_Tenant_Isolation_RBAC_And_Security_Architecture.md)
- [Secrets/credentials/audit](docs/04_architecture_deep_dive/42_Secrets_Credential_Encryption_And_Audit_Architecture.md)
- [Observability](docs/04_architecture_deep_dive/43_Observability_Logging_Metrics_Tracing_Spec.md)
- [Scaling/failure modes](docs/04_architecture_deep_dive/44_Scaling_Performance_And_Failure_Mode_Architecture.md)
- [Architecture decisions](docs/04_architecture_deep_dive/45_Architecture_Decision_Records_ADR.md)

Development standards:

- [Coding standards, module structure, and component reuse](docs/05_development_standards/46_Coding_Standards_Module_Component_Guide.md)
- [Git workflow, build, test, PR, and merge guide](docs/05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md)
- [Definition of Ready, Definition of Done, and task workflow](docs/05_development_standards/48_Definition_Of_Ready_Done_And_Task_Workflow.md)
- [Testing strategy and quality gates](docs/05_development_standards/49_Testing_Strategy_And_Quality_Gates.md)
- [API response, error code, and logging standard](docs/05_development_standards/50_API_Response_Error_Logging_Standard.md)
- [Environment, config, and secrets guide](docs/05_development_standards/51_Environment_Config_Secrets_Guide.md)
- [Database migration, seed, and data safety workflow](docs/05_development_standards/52_Database_Migration_Seed_Data_Workflow.md)
- [Frontend app shell and UI implementation standard](docs/05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md)
- [Multi-agent task board conflict workflow](docs/05_development_standards/54_Multi_Agent_Task_Board_Conflict_Workflow.md)
