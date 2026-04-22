# Database Migrations

This directory stores PostgreSQL migration files.

No domain tables are added in the DB skeleton task. Future migration files must follow:

```text
0001_create_tenants.sql
0002_create_users.sql
0003_create_wallets_and_ledger.sql
```

Rules:

- Use four-digit increasing versions.
- Use lowercase descriptive names with underscores.
- Do not edit migrations after they have run in a shared environment.
- Add a new migration to fix a previous migration.
- Explain rollback/data impact in the PR.
