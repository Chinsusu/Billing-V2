# Bộ tài liệu dự án nền tảng VPS/Proxy hybrid multi-tenant - v1.10 Local Development Runbook

## Mô tả
Gói này là bản mở rộng technical handoff cho dự án web thuê/bán VPS/Proxy theo mô hình:
- Hybrid multi-provider
- Multi-tenant
- White-label reseller
- Wallet-first
- Admin có bán lẻ trực tiếp
- Reseller có storefront/client/wallet/pricing riêng

Bản v1.1 đã vá các điểm P0 về reseller settlement, tenant isolation, inventory locking, provisioning idempotency, credential security, audit traceability và abuse control phase 1.

Bản v1.2 bổ sung bộ tài liệu `14–23` để chuyển blueprint thành tài liệu bàn giao dev/backend/frontend/QA/DevOps.

Bản v1.3 bổ sung lớp execution/operations/launch `24–36`.

Bản v1.4 bổ sung lớp architecture deep dive `37–45` để khóa hướng Go backend, PostgreSQL transaction, async worker/outbox, provider runtime, security, secrets, observability, scaling và ADR.

Bản v1.5 bổ sung chuẩn coding, cấu trúc module/component, rule tách shared component, giới hạn dòng mỗi file và chuẩn đặt tên dễ hiểu.

Bản v1.6 bổ sung workflow Git chuẩn cho tạo branch, dev, build, test, commit, pull request, review, merge, release và hotfix.

Bản v1.7 bổ sung bộ guardrail trước dev: Definition of Ready/Done, testing strategy, API/error/logging standard, environment/config/secrets guide và database migration workflow.

Bản v1.8 bổ sung chuẩn frontend app shell: frontend phải là app chạy được với package scripts, navigation thật, screen registry, mock data layer và build validation; chỉ làm HTML tĩnh không được xem là hoàn thành.

Bản v1.9 bổ sung workflow task board cho nhiều agent: `TASKS.md` chỉ là index, mỗi task active có file riêng trong `tasks/active/` để giảm conflict khi claim, review, done hoặc block task.

Bản v1.10 bổ sung local development runbook để dev mới có thể setup env, chạy API, kiểm tra health/readiness, dùng migration runner và chạy quality gate local theo script trong repo.

## Nguyên tắc chung
- Chưa code.
- Không ràng buộc framework/ngôn ngữ triển khai.
- Tập trung vào data contract, API behavior, permission, worker, UI wireflow, QA, deployment và notification.
- Mọi rule liên quan tiền, tenant, credential, provider và provisioning được xem là P0.
- Dev không nên tự đoán ở các flow tiền, tenant, provisioning, stock, credential.

## Thành phần gói

### Tài liệu nền đã vá từ bản v1.1
- `01_product_foundation/01_Product_Scope_Business_Model.md`
- `01_product_foundation/02_Tenant_Model_Role_Architecture.md`
- `01_product_foundation/03_Product_Catalog_Pricing_Rules.md`
- `01_product_foundation/04_Billing_Wallet_Ledger_Spec.md`
- `01_product_foundation/05_Provisioning_Provider_Adapter_Spec.md`
- `01_product_foundation/06_Order_Service_Lifecycle_State_Machine.md`
- `01_product_foundation/07_Portal_Functional_Spec.md`
- `01_product_foundation/08_Audit_Reports_Operational_Control.md`
- `01_product_foundation/09_Reseller_Settlement_Ledger_Model.md`
- `01_product_foundation/10_Tenant_Security_Access_Control_Spec.md`
- `01_product_foundation/11_Provisioning_Idempotency_And_Inventory_Locking.md`
- `01_product_foundation/12_API_Data_Model_Acceptance_Criteria.md`
- `01_product_foundation/13_Abuse_Fraud_Operational_Policy_Phase1.md`

### Tài liệu technical handoff thêm ở bản v1.2
- `02_technical_handoff/14_System_Architecture_Blueprint.md`
- `02_technical_handoff/15_Database_Schema_And_ERD.md`
- `02_technical_handoff/16_API_Contract_And_Permission_Spec.md`
- `02_technical_handoff/17_RBAC_Permission_Matrix.md`
- `02_technical_handoff/18_Provider_Adapter_Technical_Spec.md`
- `02_technical_handoff/19_Worker_Queue_And_Cron_Jobs_Spec.md`
- `02_technical_handoff/20_UI_Wireflow_And_Screen_Spec.md`
- `02_technical_handoff/21_QA_Test_Cases_And_Acceptance_Plan.md`
- `02_technical_handoff/22_Deployment_DevOps_And_Environment_Runbook.md`
- `02_technical_handoff/23_Notification_Email_Telegram_Template_Spec.md`

### Tài liệu execution/operations/launch thêm ở bản v1.3
- `03_execution_operations_launch/24_Project_Roadmap_Milestones_And_Sprint_Plan.md`
- `03_execution_operations_launch/25_Backlog_Epics_User_Stories_And_Task_Breakdown.md`
- `03_execution_operations_launch/26_MVP_Scope_Lock_And_Non_Goals.md`
- `03_execution_operations_launch/27_Developer_Onboarding_Guide.md`
- `03_execution_operations_launch/28_Finance_Reconciliation_SOP.md`
- `03_execution_operations_launch/29_Customer_Support_SOP_And_Macro_Templates.md`
- `03_execution_operations_launch/30_Provider_Onboarding_And_Scoring_Checklist.md`
- `03_execution_operations_launch/31_Incident_Response_And_Disaster_Recovery_Playbook.md`
- `03_execution_operations_launch/32_Abuse_Compliance_Takedown_SOP.md`
- `03_execution_operations_launch/33_Launch_Checklist_And_Go_No_Go_Criteria.md`
- `03_execution_operations_launch/34_Beta_Pilot_Program_And_Feedback_Loop.md`
- `03_execution_operations_launch/35_Reseller_Acquisition_And_Enablement_Playbook.md`
- `03_execution_operations_launch/36_Go_To_Market_Positioning_And_Offer_Strategy.md`
- `03_execution_operations_launch/65_MVP_Launch_Gap_Audit.md`

### Tài liệu architecture deep dive thêm ở bản v1.4
- `04_architecture_deep_dive/37_Go_Backend_Architecture_And_Module_Boundaries.md`
- `04_architecture_deep_dive/38_PostgreSQL_Data_Consistency_Transaction_Design.md`
- `04_architecture_deep_dive/39_Async_Worker_Outbox_And_Job_Architecture.md`
- `04_architecture_deep_dive/40_Provider_Adapter_Runtime_And_Error_Taxonomy.md`
- `04_architecture_deep_dive/41_Tenant_Isolation_RBAC_And_Security_Architecture.md`
- `04_architecture_deep_dive/42_Secrets_Credential_Encryption_And_Audit_Architecture.md`
- `04_architecture_deep_dive/43_Observability_Logging_Metrics_Tracing_Spec.md`
- `04_architecture_deep_dive/44_Scaling_Performance_And_Failure_Mode_Architecture.md`
- `04_architecture_deep_dive/45_Architecture_Decision_Records_ADR.md`

### Tài liệu development standards thêm ở bản v1.5
- `05_development_standards/46_Coding_Standards_Module_Component_Guide.md`
- `05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md`
- `05_development_standards/48_Definition_Of_Ready_Done_And_Task_Workflow.md`
- `05_development_standards/49_Testing_Strategy_And_Quality_Gates.md`
- `05_development_standards/50_API_Response_Error_Logging_Standard.md`
- `05_development_standards/51_Environment_Config_Secrets_Guide.md`
- `05_development_standards/52_Database_Migration_Seed_Data_Workflow.md`
- `05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md`
- `05_development_standards/54_Multi_Agent_Task_Board_Conflict_Workflow.md`
- `05_development_standards/55_Local_Development_Runbook.md`
- `05_development_standards/56_Billing_API_Operational_Reference.md`
- `05_development_standards/57_Billing_Operations_Runbook.md`
- `05_development_standards/58_Provisioning_Ops_Readiness_Checklist.md`
- `05_development_standards/59_API_Contract_Drift_Guard.md`
- `05_development_standards/60_Provider_Sandbox_Contract_Checklist.md`
- `05_development_standards/61_Task_Board_Consistency_Guard.md`
- `05_development_standards/62_API_Error_Code_Drift_Guard.md`
- `05_development_standards/63_Validation_Command_Matrix.md`
- `05_development_standards/64_Public_Display_ID_And_Backend_Reference_Policy.md`

### Tài liệu tổng hợp và ghi chú
- `VPS_Proxy_Project_Master_Document.md`
- `changelogs/CHANGELOG_FIXES.md`
- `changelogs/CHANGELOG_TECHNICAL_HANDOFF_v1_2.md`
- `changelogs/CHANGELOG_EXECUTION_OPERATIONS_LAUNCH_v1_3.md`
- `changelogs/CHANGELOG_ARCHITECTURE_DEEP_DIVE_v1_4.md`
- `MANIFEST.txt`

## Cách đọc đề xuất

### Cho founder/product owner
1. `01_product_foundation/01_Product_Scope_Business_Model.md`
2. `01_product_foundation/09_Reseller_Settlement_Ledger_Model.md`
3. `02_technical_handoff/14_System_Architecture_Blueprint.md`
4. `02_technical_handoff/20_UI_Wireflow_And_Screen_Spec.md`
5. `02_technical_handoff/21_QA_Test_Cases_And_Acceptance_Plan.md`

### Cho backend/dev lead
1. `02_technical_handoff/14_System_Architecture_Blueprint.md`
2. `04_architecture_deep_dive/37_Go_Backend_Architecture_And_Module_Boundaries.md`
3. `04_architecture_deep_dive/38_PostgreSQL_Data_Consistency_Transaction_Design.md`
4. `04_architecture_deep_dive/39_Async_Worker_Outbox_And_Job_Architecture.md`
5. `04_architecture_deep_dive/40_Provider_Adapter_Runtime_And_Error_Taxonomy.md`
6. `04_architecture_deep_dive/41_Tenant_Isolation_RBAC_And_Security_Architecture.md`
7. `04_architecture_deep_dive/42_Secrets_Credential_Encryption_And_Audit_Architecture.md`
8. `02_technical_handoff/15_Database_Schema_And_ERD.md`
9. `02_technical_handoff/16_API_Contract_And_Permission_Spec.md`
10. `02_technical_handoff/17_RBAC_Permission_Matrix.md`
11. `05_development_standards/46_Coding_Standards_Module_Component_Guide.md`
12. `05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md`
13. `05_development_standards/48_Definition_Of_Ready_Done_And_Task_Workflow.md`
14. `05_development_standards/49_Testing_Strategy_And_Quality_Gates.md`
15. `05_development_standards/50_API_Response_Error_Logging_Standard.md`
16. `05_development_standards/51_Environment_Config_Secrets_Guide.md`
17. `05_development_standards/52_Database_Migration_Seed_Data_Workflow.md`
18. `05_development_standards/54_Multi_Agent_Task_Board_Conflict_Workflow.md`
19. `05_development_standards/55_Local_Development_Runbook.md`
20. `05_development_standards/56_Billing_API_Operational_Reference.md`
21. `05_development_standards/64_Public_Display_ID_And_Backend_Reference_Policy.md`
22. `05_development_standards/57_Billing_Operations_Runbook.md`
23. `05_development_standards/58_Provisioning_Ops_Readiness_Checklist.md`

### Cho frontend
1. `02_technical_handoff/16_API_Contract_And_Permission_Spec.md`
2. `02_technical_handoff/17_RBAC_Permission_Matrix.md`
3. `04_architecture_deep_dive/41_Tenant_Isolation_RBAC_And_Security_Architecture.md`
4. `04_architecture_deep_dive/42_Secrets_Credential_Encryption_And_Audit_Architecture.md`
5. `02_technical_handoff/20_UI_Wireflow_And_Screen_Spec.md`
6. `02_technical_handoff/23_Notification_Email_Telegram_Template_Spec.md`
7. `05_development_standards/46_Coding_Standards_Module_Component_Guide.md`
8. `05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md`
9. `05_development_standards/50_API_Response_Error_Logging_Standard.md`
10. `05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md`
11. `05_development_standards/54_Multi_Agent_Task_Board_Conflict_Workflow.md`
12. `05_development_standards/55_Local_Development_Runbook.md`
13. `05_development_standards/56_Billing_API_Operational_Reference.md`
14. `05_development_standards/64_Public_Display_ID_And_Backend_Reference_Policy.md`

### Cho QA
1. `01_product_foundation/12_API_Data_Model_Acceptance_Criteria.md`
2. `02_technical_handoff/15_Database_Schema_And_ERD.md`
3. `02_technical_handoff/16_API_Contract_And_Permission_Spec.md`
4. `04_architecture_deep_dive/38_PostgreSQL_Data_Consistency_Transaction_Design.md`
5. `04_architecture_deep_dive/41_Tenant_Isolation_RBAC_And_Security_Architecture.md`
6. `04_architecture_deep_dive/42_Secrets_Credential_Encryption_And_Audit_Architecture.md`
7. `02_technical_handoff/21_QA_Test_Cases_And_Acceptance_Plan.md`
8. `05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md`
9. `05_development_standards/48_Definition_Of_Ready_Done_And_Task_Workflow.md`
10. `05_development_standards/49_Testing_Strategy_And_Quality_Gates.md`
11. `05_development_standards/54_Multi_Agent_Task_Board_Conflict_Workflow.md`
12. `05_development_standards/55_Local_Development_Runbook.md`
13. `05_development_standards/57_Billing_Operations_Runbook.md`
14. `05_development_standards/58_Provisioning_Ops_Readiness_Checklist.md`

### Cho DevOps/Ops
1. `02_technical_handoff/14_System_Architecture_Blueprint.md`
2. `02_technical_handoff/18_Provider_Adapter_Technical_Spec.md`
3. `04_architecture_deep_dive/39_Async_Worker_Outbox_And_Job_Architecture.md`
4. `04_architecture_deep_dive/40_Provider_Adapter_Runtime_And_Error_Taxonomy.md`
5. `04_architecture_deep_dive/43_Observability_Logging_Metrics_Tracing_Spec.md`
6. `04_architecture_deep_dive/44_Scaling_Performance_And_Failure_Mode_Architecture.md`
7. `02_technical_handoff/22_Deployment_DevOps_And_Environment_Runbook.md`
8. `02_technical_handoff/23_Notification_Email_Telegram_Template_Spec.md`
9. `05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md`
10. `05_development_standards/51_Environment_Config_Secrets_Guide.md`
11. `05_development_standards/52_Database_Migration_Seed_Data_Workflow.md`
12. `05_development_standards/54_Multi_Agent_Task_Board_Conflict_Workflow.md`
13. `05_development_standards/55_Local_Development_Runbook.md`
14. `05_development_standards/56_Billing_API_Operational_Reference.md`
15. `05_development_standards/57_Billing_Operations_Runbook.md`
16. `05_development_standards/58_Provisioning_Ops_Readiness_Checklist.md`

## 10 luật nền phải giữ
1. Không provision nếu tiền chưa được debit/lock hợp lệ.
2. Không debit ví nếu không tạo được ledger entry.
3. Không có ledger thì giao dịch không tồn tại.
4. Không có tenant context thì không cho đọc/ghi dữ liệu tenant.
5. Không retry provisioning nếu không biết provider đã tạo tài nguyên hay chưa.
6. Không hiển thị action nếu capability snapshot không cho phép.
7. Không sửa transaction cũ; chỉ tạo adjustment/reversal.
8. Không lưu credential plaintext trong log/audit.
9. Không cho client reseller provision nếu reseller wallet không đủ reseller cost.
10. Không dùng giá/policy hiện tại để xử lý tranh chấp order cũ; dùng snapshot lúc mua.

## Mục tiêu sau bản v1.10
Sau khi đọc xong gói này, team dev phải trả lời được:
- Cần tạo bảng nào và bảng nào bắt buộc có tenant_id.
- API nào cần build và role nào được gọi.
- Checkout, wallet, reseller settlement và provisioning chạy theo flow nào.
- Provider adapter phải trả về gì khi success/fail/timeout/partial success.
- Worker/cron nào phải chạy nền.
- UI cần màn nào và action nào phải ẩn/chặn.
- QA phải test case nào trước khi nghiệm thu.
- Production cần backup, monitoring, secret và rollback như thế nào.
- Backend Go nên chia module/process như thế nào.
- PostgreSQL transaction boundary cho checkout/top-up/refund/provisioning nằm ở đâu.
- Outbox/job/worker/scheduler chạy, retry, lock và recover thế nào.
- Provider error nào được retry, error nào phải manual review.
- Tenant/RBAC/credential/observability/scaling được bảo vệ ở lớp nào.
- File code nên được tách thế nào để không vượt 500 dòng.
- Khi nào logic dùng chung phải tách thành module/component riêng.
- Tên package, file, function và component nên đặt thế nào để ít gây hiểu lầm.
- Branch mới phải tạo từ đâu và đặt tên thế nào.
- Trước khi commit, mở PR và merge cần chạy build/test gì.
- Khi nào được merge, revert, hotfix hoặc tag release.
- Khi nào task đủ ready để bắt đầu dev và khi nào được coi là done.
- Test nào bắt buộc cho từng loại thay đổi.
- API success/error/validation/pagination/logging phải theo format nào.
- Config và secret phải được đặt, validate và rotate thế nào.
- Migration, seed, backfill và rollback database phải đi theo quy trình nào.
- Frontend app shell tối thiểu phải có package scripts, entrypoint app, navigation, screen registry, mock data layer và build validation nào.
- Khi nào một task frontend chỉ tạo HTML tĩnh sẽ bị từ chối.
- Nhiều agent claim/review/done task thế nào mà không cùng sửa một bảng task trung tâm.
- Dev mới setup env local, chạy API, kiểm tra health/readiness, dùng migration runner và chạy quality gate nào trước PR.
- Khi conflict `TASKS.md` xảy ra thì giữ task row và task-file status thế nào.


---

## v1.3 Update — Execution, Operations & Launch Package

**Date:** 2026-04-22

This full package includes the previous `00–23` product/technical docs and adds the execution/operations/launch layer `24–36`:

```text
03_execution_operations_launch/24_Project_Roadmap_Milestones_And_Sprint_Plan.md
03_execution_operations_launch/25_Backlog_Epics_User_Stories_And_Task_Breakdown.md
03_execution_operations_launch/26_MVP_Scope_Lock_And_Non_Goals.md
03_execution_operations_launch/27_Developer_Onboarding_Guide.md
03_execution_operations_launch/28_Finance_Reconciliation_SOP.md
03_execution_operations_launch/29_Customer_Support_SOP_And_Macro_Templates.md
03_execution_operations_launch/30_Provider_Onboarding_And_Scoring_Checklist.md
03_execution_operations_launch/31_Incident_Response_And_Disaster_Recovery_Playbook.md
03_execution_operations_launch/32_Abuse_Compliance_Takedown_SOP.md
03_execution_operations_launch/33_Launch_Checklist_And_Go_No_Go_Criteria.md
03_execution_operations_launch/34_Beta_Pilot_Program_And_Feedback_Loop.md
03_execution_operations_launch/35_Reseller_Acquisition_And_Enablement_Playbook.md
03_execution_operations_launch/36_Go_To_Market_Positioning_And_Offer_Strategy.md
```

The v1.3 layer is designed to help the team execute, test, pilot, support, reconcile finance, respond to incidents, onboard providers/resellers, and prepare GTM without losing control of scope.


---

## v1.4 Update — Architecture Deep Dive Package

**Date:** 2026-04-22

This package adds the backend architecture deep-dive layer `37–45`:

```text
04_architecture_deep_dive/37_Go_Backend_Architecture_And_Module_Boundaries.md
04_architecture_deep_dive/38_PostgreSQL_Data_Consistency_Transaction_Design.md
04_architecture_deep_dive/39_Async_Worker_Outbox_And_Job_Architecture.md
04_architecture_deep_dive/40_Provider_Adapter_Runtime_And_Error_Taxonomy.md
04_architecture_deep_dive/41_Tenant_Isolation_RBAC_And_Security_Architecture.md
04_architecture_deep_dive/42_Secrets_Credential_Encryption_And_Audit_Architecture.md
04_architecture_deep_dive/43_Observability_Logging_Metrics_Tracing_Spec.md
04_architecture_deep_dive/44_Scaling_Performance_And_Failure_Mode_Architecture.md
04_architecture_deep_dive/45_Architecture_Decision_Records_ADR.md
```

The v1.4 layer is intended to be the technical decision baseline before implementation starts. It locks the modular monolith direction, PostgreSQL consistency model, async/outbox architecture, provider runtime error taxonomy, tenant/RBAC/security, secrets/credential handling, observability, scaling/failure modes, and ADR record.


---

## v1.5 Update — Coding Standards & Module/Component Guide

**Date:** 2026-04-22

This package adds the implementation standards layer `46`:

```text
05_development_standards/46_Coding_Standards_Module_Component_Guide.md
```

The v1.5 layer defines coding standards, module boundaries, component reuse rules, file length limits, naming rules, shared ownership, and review checklist before the first production code is added.


---

## v1.6 Update — Git Workflow, Build, Test, PR & Merge

**Date:** 2026-04-22

This package adds the team delivery workflow layer `47`:

```text
05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md
```

The v1.6 layer defines branch naming, local development flow, build/test gates, commit style, pull request checklist, review rules, merge policy, revert, hotfix, release tagging, CI expectations, and secret handling in Git.


---

## v1.7 Update — Pre-Dev Guardrails

**Date:** 2026-04-22

This package adds the pre-development guardrail layer `48–52`:

```text
05_development_standards/48_Definition_Of_Ready_Done_And_Task_Workflow.md
05_development_standards/49_Testing_Strategy_And_Quality_Gates.md
05_development_standards/50_API_Response_Error_Logging_Standard.md
05_development_standards/51_Environment_Config_Secrets_Guide.md
05_development_standards/52_Database_Migration_Seed_Data_Workflow.md
```

The v1.7 layer defines task readiness, done criteria, testing gates, API/error/logging format, config/secret handling, database migration rules, seed data safety, and rollback expectations before production code starts.


---

## v1.8 Update — Frontend App Shell Standard

**Date:** 2026-04-22

This package adds the frontend delivery standard layer `53`:

```text
05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md
```

The v1.8 layer defines the minimum frontend app-shell deliverable before backend route wiring: runnable package scripts, app entrypoint, navigation, screen registry, mock data layer, shared layout/component structure, build validation, and PR checklist. Static HTML alone is not accepted for frontend app-shell tasks.


---

## v1.9 Update — Multi-Agent Task Board Conflict Workflow

**Date:** 2026-04-22

This package adds the multi-agent coordination layer `54`:

```text
05_development_standards/54_Multi_Agent_Task_Board_Conflict_Workflow.md
```

The v1.9 layer defines `TASKS.md` as a stable task index and moves mutable active task status into one file per task under `tasks/active/`. Agents claim, review, block, and mark done by editing only their task file, reducing merge conflicts between unrelated coding work.


---

## v1.10 Update — Local Development Runbook

**Date:** 2026-04-22

This package adds the local development runbook layer `55`:

```text
05_development_standards/55_Local_Development_Runbook.md
```

The v1.10 layer defines local backend setup, safe `.env` handling, API run commands, health/readiness checks, migration runner usage, PostgreSQL local notes, and PR validation gates.

## v1.11 Billing operations and provisioning readiness addendum

**Date:** 2026-04-24

This package adds the billing operations reference layer `56-58`:

```text
05_development_standards/56_Billing_API_Operational_Reference.md
05_development_standards/57_Billing_Operations_Runbook.md
05_development_standards/58_Provisioning_Ops_Readiness_Checklist.md
```

The v1.11 layer defines billing API routes, paid-order fulfillment checks, provisioning worker run modes, job recovery actions, smoke verification, and local/sandbox no-go rules for money and provider-state safety.

## v1.12 Validation command matrix addendum

**Date:** 2026-04-25

This package adds one source of truth for validation and smoke commands:

```text
05_development_standards/63_Validation_Command_Matrix.md
```

The v1.12 layer maps docs-only, backend, frontend, DB, provider, full-stack, CI, and task-board changes to the local commands agents should run before PR review.
