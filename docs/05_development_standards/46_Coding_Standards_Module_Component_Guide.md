# Coding Standards, Module Structure, and Component Reuse Guide

**Version:** v1.5  
**Date:** 2026-04-22  
**Scope:** Go backend, future frontend, scripts, migrations, and project documentation.

## Mục tiêu

Tài liệu này khóa cách viết code và cách chia module/component để team triển khai dễ đọc, dễ review, dễ test và ít hiểu lầm.

Các rule ở đây ưu tiên sự rõ ràng hơn sự phức tạp. Nếu một cách đặt tên hoặc cách chia file khiến người mới đọc phải đoán, cách đó chưa đạt chuẩn.

## Luật bắt buộc

1. Mỗi file code hoặc tài liệu triển khai không vượt quá 500 dòng.
2. Khi file vượt 350 dòng, người viết phải cân nhắc tách trước khi merge.
3. Không tạo package tên `common`, `utils`, `helpers`, `misc`, `base`, `core` nếu không có nghĩa nghiệp vụ rõ.
4. Code liên quan tiền, ví, ledger, tenant, quyền truy cập, credential, provider, provisioning và audit phải có owner module rõ ràng.
5. Logic dùng ở nhiều nơi không được copy lan rộng. Phải tách thành module/component dùng chung khi có từ 3 nơi sử dụng, hoặc từ 2 nơi nếu logic đó liên quan tiền, bảo mật, tenant, audit hoặc provider.
6. Không để business rule nằm trong `cmd/`, HTTP handler, worker loop, scheduler loop hoặc script.
7. Không log credential, token, API key, password, private key, proxy password hoặc raw provider response có secret.
8. Không dùng thuật ngữ dễ nhầm nếu có từ rõ hơn.
9. Mỗi module phải có test cho rule tiền, tenant, quyền, provisioning và error quan trọng trước khi coi là xong.
10. Nếu phải phá rule trong tài liệu này, PR phải ghi rõ lý do và kế hoạch sửa.

## Cấu trúc workspace

```text
cmd/
  api/         chỉ khởi động HTTP API
  worker/      chỉ khởi động worker chạy job nền
  scheduler/   chỉ khởi động lịch chạy định kỳ
  migrate/     chỉ chạy migration
  cli/         chỉ chứa command cho operator nội bộ

internal/
  app/         nối config, db, logger, module và entrypoint lại với nhau
  platform/    hạ tầng dùng chung, không chứa business rule
  modules/     nghiệp vụ chính của sản phẩm

migrations/    migration PostgreSQL
scripts/       script dev/ops nhỏ, rõ mục đích
docs/          tài liệu dự án
```

## Rule import

`cmd/*` chỉ gọi `internal/app`.

`internal/app` được phép nối các module và platform lại với nhau, nhưng không chứa business rule.

`internal/platform/*` không được import `internal/modules/*`.

`internal/modules/<module>` được import `internal/platform/*` và interface rõ ràng của module khác khi cần.

Một module không được gọi thẳng HTTP handler, SQL store hoặc worker private của module khác.

Nếu module A cần dùng năng lực của module B, module B phải public một service/interface nhỏ, có tên rõ, và che giấu chi tiết bên trong.

## Cấu trúc module backend

Mỗi module trong `internal/modules/<module>` nên bắt đầu đơn giản. Không tách quá sớm, nhưng phải tách khi file hoặc flow bắt đầu dài.

```text
internal/modules/wallet/
  README.md              module làm gì, không làm gì, bảng nào, rule P0 nào
  model.go               kiểu dữ liệu nghiệp vụ
  errors.go              lỗi có nghĩa, dùng lại trong module
  permissions.go         quyền cần kiểm tra trong module
  routes.go              đăng ký route HTTP nếu module có API
  handler_topup.go       nhận request, validate input, gọi service
  service_topup.go       rule nạp tiền
  service_debit.go       rule trừ/lock tiền
  store.go               interface lưu/đọc dữ liệu
  postgres_store.go      SQL PostgreSQL
  events.go              event phát ra cho module khác hoặc outbox
  jobs.go                job nền thuộc module nếu có
  *_test.go              test theo flow và rule quan trọng
```

Không phải module nào cũng cần đủ các file trên. Chỉ tạo file khi có nội dung thật.

## Vai trò từng lớp

HTTP handler chỉ làm các việc sau:

- Đọc request.
- Kiểm tra định dạng input.
- Lấy user, tenant và request id từ context.
- Gọi service của module.
- Map kết quả/lỗi sang response.

Service của module làm các việc sau:

- Kiểm tra business rule.
- Quyết định transaction boundary.
- Gọi store hoặc service module khác qua interface.
- Tạo ledger entry, audit event, outbox event khi flow yêu cầu.

Store làm các việc sau:

- Chứa SQL.
- Map row sang model.
- Không tự quyết định business rule.
- Không tự gọi module khác.

Worker/scheduler làm các việc sau:

- Lấy job cần chạy.
- Gọi service đúng module.
- Ghi kết quả job.
- Không nhét rule nghiệp vụ trực tiếp trong vòng lặp worker.

## Module nào nên là shared backend

Các phần sau là shared rõ ràng và nên nằm trong `internal/platform` hoặc module owner tương ứng:

```text
internal/platform/config       đọc config/env
internal/platform/db           kết nối database, transaction helper
internal/platform/logger       log chuẩn, field chuẩn, redaction
internal/platform/httpserver   HTTP server setup
internal/platform/middleware   request id, auth shell, recover, tenant context shell
internal/platform/clock        thời gian dùng cho test và job
internal/platform/queue        queue/job/outbox helper
internal/platform/crypto       mã hóa/giải mã secret
internal/platform/metrics      metric chung
internal/platform/tracing      trace chung
internal/platform/ratelimit    giới hạn request
internal/platform/email        gửi email
internal/platform/telegram     gửi Telegram
```

Business shared phải có owner module, không đưa vào `platform` chỉ vì nhiều nơi dùng.

Ví dụ:

- `ledger` sở hữu ledger entry, adjustment, reversal và invariant tiền.
- `wallet` sở hữu ví, số dư, lock, debit, credit.
- `tenant` sở hữu tenant, reseller tenant, storefront mapping và tenant context ở cấp nghiệp vụ.
- `rbac` sở hữu role, permission và check quyền.
- `audit` sở hữu audit event và format audit.
- `provider` sở hữu provider adapter contract, capability snapshot và error taxonomy.
- `provisioning` sở hữu flow cấp phát, retry, reconcile và manual review.
- `notification` sở hữu template và quyết định gửi thông báo.

## Khi nào tách shared

Giữ local nếu logic chỉ dùng ở 1 nơi và chưa có dấu hiệu dùng lại.

Cân nhắc tách khi logic dùng ở 2 nơi.

Bắt buộc tách khi logic dùng ở 3 nơi.

Bắt buộc tách ngay khi logic dùng ở 2 nơi và liên quan:

- Tiền hoặc ledger.
- Tenant isolation.
- RBAC/permission.
- Credential/secret.
- Provider/provisioning.
- Audit/compliance.
- Idempotency.
- Rate limit hoặc abuse control.

Không tách bằng cách tạo `utils`. Hãy đặt tên theo mục đích thật:

```text
money_amount.go
tenant_context.go
permission_check.go
credential_redaction.go
provider_error.go
idempotency_key.go
pagination.go
```

## Rule tách file

Tách theo flow, không tách ngẫu nhiên.

Ví dụ tốt:

```text
service_checkout.go
service_refund.go
service_renewal.go
handler_checkout.go
handler_refund.go
store_order.go
store_order_item.go
```

Ví dụ nên tránh:

```text
service.go          quá rộng khi module lớn
helper.go           không nói rõ giúp gì
common.go           không nói rõ thuộc rule nào
misc.go             không review được ý định
manager.go          dễ thành nơi chứa mọi thứ
```

Nếu một function vượt 80 dòng, phải kiểm tra lại. Nếu function đó xử lý nhiều bước nghiệp vụ, tách thành các hàm nhỏ theo bước rõ ràng.

Nếu một test file vượt 500 dòng, tách theo flow:

```text
checkout_test.go
refund_test.go
renewal_test.go
permission_test.go
```

## Rule đặt tên

Tên phải trả lời được 3 câu hỏi:

1. Đây là dữ liệu hay hành động gì?
2. Thuộc module nào?
3. Có thể bị hiểu nhầm với khái niệm khác không?

Ưu tiên tên nghiệp vụ rõ:

```text
CreateLedgerEntry
ReserveWalletBalance
ReleaseWalletHold
CheckTenantAccess
ProvisionService
RecordAuditEvent
SendOrderPaidEmail
```

Tránh viết tắt nếu không phổ biến:

```text
cfg      chỉ dùng local rất ngắn, không dùng làm field public
svc      chỉ dùng local rất ngắn, không dùng làm package hoặc struct public
tx       chỉ dùng cho database transaction trong phạm vi ngắn
amt      không dùng; dùng amount
usr      không dùng; dùng user
prov     không dùng; dùng provider
```

Các từ dễ nhầm phải dùng cụ thể:

```text
transaction        tránh dùng một mình
db_transaction     transaction của database
ledger_entry       dòng tiền trong ledger
wallet_movement    thay đổi số dư ví

account            tránh dùng một mình
user_account       tài khoản đăng nhập
wallet_account     tài khoản ví
provider_account   tài khoản ở nhà cung cấp

service            tránh dùng một mình nếu không rõ
customer_service   dịch vụ khách mua
module_service     service xử lý rule của module
provider_service   dịch vụ phía provider
```

## Ngôn ngữ trong code và tài liệu

Comment và tài liệu phải dùng câu ngắn, rõ, tránh thuật ngữ nặng nếu không cần.

Nếu bắt buộc dùng thuật ngữ như `outbox`, `idempotency`, `capability snapshot`, `tenant isolation`, phải giải thích một lần trong README module hoặc tài liệu liên quan.

Không dùng từ mơ hồ trong lỗi, log hoặc comment:

```text
bad
invalid
failed
error happened
something went wrong
not allowed
```

Hãy nói rõ hơn:

```text
wallet balance is not enough
tenant context is missing
provider response is pending review
permission order.refund is required
ledger entry was not created
```

Thông báo cho user cuối phải đơn giản hơn log nội bộ. Log có thể có `error_code`, `request_id`, `tenant_id`, `order_id`, nhưng không có secret.

## Rule frontend khi thêm web/app

Nếu thêm frontend sau này, dùng cấu trúc theo feature:

```text
web/src/features/orders/
web/src/features/wallet/
web/src/features/catalog/
web/src/features/tenants/
web/src/features/providers/

web/src/components/ui/        button, input, dialog, table, badge
web/src/components/layout/    sidebar, topbar, page shell
web/src/components/shared/    component nghiệp vụ dùng ở nhiều feature
web/src/lib/                  API client, date, money format, validation nhỏ
```

Component chỉ dùng trong một feature thì để trong feature đó.

Component dùng ở 2 feature thì có thể tách nếu ổn định.

Component dùng ở 3 feature thì phải tách vào shared.

Không tạo component tên `Box`, `Panel`, `Thing`, `BaseCard`, `CommonTable` nếu không nói rõ mục đích. Dùng tên cụ thể hơn:

```text
OrderStatusBadge
WalletBalanceCard
TenantSwitcher
ProviderCapabilityTable
LedgerEntryTable
```

Mỗi component không vượt 500 dòng. Nếu vượt, tách phần data loading, table columns, form schema, dialog hoặc subcomponent.

## Migration và SQL

Migration phải nhỏ theo nhóm thay đổi. Không gom mọi bảng vào một file quá dài.

Tên migration phải nói rõ mục đích:

```text
0001_create_tenants.sql
0002_create_wallets_and_ledger.sql
0003_create_orders.sql
0004_create_provider_inventory.sql
```

SQL trong Go phải nằm trong store của module owner. Nếu query dùng ở nhiều module, kiểm tra lại ownership thay vì copy query.

Không để business rule chỉ nằm trong SQL nếu Go service không biết rule đó tồn tại. Rule quan trọng cần có test ở service hoặc integration test.

## Review checklist

Trước khi merge code, kiểm tra:

- File nào vượt 350 dòng chưa?
- File nào vượt 500 dòng không?
- Có package/file tên `utils`, `common`, `helpers`, `misc`, `manager` không?
- Business rule có nằm trong handler, worker loop, scheduler loop hoặc script không?
- Logic dùng ở nhiều nơi đã có owner chưa?
- Có log secret hoặc raw credential không?
- Có tên nào dễ nhầm giữa database transaction và ledger transaction không?
- Có test cho flow tiền, tenant, permission, provider/provisioning không?
- Module mới có README ngắn giải thích scope chưa?
- Lỗi trả ra user có dễ hiểu không?

## Definition of done cho module mới

Một module mới chỉ được coi là sẵn sàng khi có:

- README module mô tả module làm gì và không làm gì.
- Model chính.
- Service chứa rule nghiệp vụ.
- Store hoặc interface lưu trữ nếu cần database.
- Handler/routes nếu có API.
- Test cho rule chính.
- Audit/ledger/outbox/permission hook nếu flow yêu cầu.
- Không có file vượt 500 dòng.
- Tên package, file, type và function rõ nghĩa.
