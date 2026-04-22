# Database Migration, Seed, and Data Safety Workflow

**Version:** v1.7  
**Date:** 2026-04-22  
**Scope:** PostgreSQL migrations, seed data, backfill, rollback, data safety, and database review rules.

## Mục tiêu

Tài liệu này khóa cách thay đổi database an toàn. Mục tiêu là schema thay đổi có kiểm soát, không mất dữ liệu, không phá flow tiền/tenant và có đường rollback rõ.

## Luật bắt buộc

1. Mọi thay đổi schema phải đi qua migration.
2. Không sửa migration đã merge nếu đã chạy ở shared environment.
3. Tạo migration mới để sửa migration cũ.
4. Migration liên quan tiền, ledger, order, tenant hoặc credential phải review kỹ.
5. Không drop dữ liệu production nếu chưa có backup và plan rõ.
6. Seed data production không được nằm chung với seed local.
7. Migration phải có rollback plan.

## Naming

Tên migration dùng số thứ tự và mục đích rõ:

```text
0001_create_tenants.sql
0002_create_wallets_and_ledger.sql
0003_create_orders.sql
0004_create_provider_inventory.sql
0005_add_order_idempotency_key.sql
```

Tránh tên mơ hồ:

```text
update.sql
fix_db.sql
new_tables.sql
changes.sql
```

## Up và down

Mỗi migration nên có hướng lên và hướng rollback.

Nếu tool dùng file riêng:

```text
0001_create_tenants.up.sql
0001_create_tenants.down.sql
```

Nếu tool dùng chung một file, phải có section rõ.

Rollback không phải lúc nào cũng drop dữ liệu. Với migration rủi ro cao, rollback plan có thể là deploy code cũ cộng migration sửa tiếp.

## Không sửa migration đã chạy

Nếu migration đã merge và chạy ở dev/staging/production:

- Không edit file cũ.
- Không đổi ý nghĩa file cũ.
- Không re-order migration cũ.
- Tạo migration mới để sửa.

Chỉ được sửa migration cũ khi chắc chắn migration đó chưa chạy ở bất kỳ shared environment nào.

## Transaction rule

Migration nên chạy trong transaction nếu PostgreSQL và tool cho phép.

Ngoại lệ:

- Một số index concurrent không chạy trong transaction.
- Một số operation lớn cần tách batch.

Nếu migration không transaction-safe, PR phải ghi rõ.

## Tenant rule

Bảng chứa dữ liệu theo tenant phải có `tenant_id`, trừ khi tài liệu schema nói rõ là global.

Khi thêm bảng tenant-scoped:

- Có `tenant_id`.
- Có index phù hợp theo `tenant_id`.
- Query sau này phải filter theo tenant.
- Unique constraint cần tính tenant nếu nghiệp vụ yêu cầu.

Ví dụ:

```text
unique(tenant_id, external_code)
index(tenant_id, created_at)
```

## Money và ledger rule

Bảng tiền/ledger phải ưu tiên tính đúng hơn tiện sửa.

Rule:

- Ledger entry append-only.
- Không update amount của ledger entry cũ.
- Không delete ledger entry ở production.
- Reversal/adjustment là entry mới.
- Amount dùng integer minor unit hoặc numeric có scale rõ.
- Currency phải rõ nếu hệ thống có nhiều tiền tệ.

Migration tiền phải được review với checklist riêng trong PR.

## Index rule

Thêm index khi:

- Query list/filter cần chạy thường xuyên.
- Foreign key được join nhiều.
- Tenant filter xuất hiện ở bảng lớn.
- Job worker claim theo trạng thái/thời gian.

Với bảng lớn, cân nhắc:

- `CREATE INDEX CONCURRENTLY`.
- Tách migration index riêng.
- Chạy ngoài giờ thấp tải.
- Có rollback plan.

Không thêm index chỉ vì "có thể cần".

## Constraint rule

Ưu tiên constraint cho invariant rõ:

- `NOT NULL` cho field bắt buộc.
- `CHECK` cho amount không âm nếu rule cho phép.
- `UNIQUE` cho idempotency key.
- Foreign key nếu không làm flow vận hành quá rủi ro.

Business rule phức tạp vẫn phải nằm ở service và có test, không chỉ dựa vào database.

## Seed data

Seed local dùng để dev nhanh.

Seed local có thể tạo:

- Tenant mẫu.
- User mẫu.
- Role/permission mẫu.
- Product catalog mẫu.
- Provider fake/sandbox.

Seed local không được chứa:

- Secret thật.
- Dữ liệu khách hàng thật.
- Provider production credential.
- Dump production.

Nếu cần seed staging/production, tạo quy trình riêng có approval.

## Backfill

Backfill dùng khi thêm dữ liệu mới cho row cũ.

Backfill plan cần có:

- Bảng nào bị ảnh hưởng.
- Số lượng row ước tính.
- Chạy một lần hay nhiều batch.
- Có lock lâu không.
- Có thể resume không.
- Cách verify.
- Cách rollback hoặc sửa nếu sai.

Backfill lớn không nên chạy chung migration boot app nếu có thể gây timeout deploy.

## Rollback plan

Mỗi migration PR phải ghi:

- Nếu migration fail giữa chừng thì làm gì.
- Nếu deploy app mới fail thì database có còn tương thích app cũ không.
- Có cần restore backup không.
- Có migration down không.
- Có data loss không.

Rollback với data loss phải được owner approve.

## Local workflow

Khi có migration tool, workflow nên là:

```bash
make db-create
make migrate-up
make test
make migrate-down
make migrate-up
```

Nếu chưa có `Makefile`, lệnh migration phải được ghi trong README hoặc runbook.

## Review checklist

Reviewer kiểm:

- Tên migration rõ không?
- Migration có chạy được trên database sạch không?
- Có sửa migration cũ đã merge không?
- Bảng tenant-scoped có `tenant_id` không?
- Bảng tiền/ledger có append-only rule không?
- Index/constraint có hợp lý không?
- Có ảnh hưởng dữ liệu cũ không?
- Backfill có batch/resume không?
- Rollback plan rõ không?
- App code và migration có tương thích thứ tự deploy không?

## Deploy order

Ưu tiên migration backward-compatible:

```text
1. Add column/table nullable hoặc không phá app cũ
2. Deploy app ghi dữ liệu mới
3. Backfill nếu cần
4. Enforce constraint sau khi dữ liệu sạch
5. Xóa field cũ ở release sau
```

Không đổi schema theo kiểu app cũ chết ngay nếu deploy app mới fail.

## Data deletion

Không hard delete dữ liệu quan trọng nếu chưa có policy.

Dữ liệu cần giữ:

- ledger
- wallet movements
- order
- audit event
- provider provisioning history
- credential access audit

Nếu cần xóa, dùng soft delete hoặc archival theo policy được duyệt.

## Definition of done

Migration task done khi:

- Migration file đúng naming.
- Chạy được trên database sạch.
- Rollback/down hoặc rollback plan rõ.
- Test liên quan pass.
- Docs/schema cập nhật nếu cần.
- PR mô tả data impact.
- Không có secret hoặc production data trong seed.
