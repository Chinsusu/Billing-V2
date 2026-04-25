# Testing Strategy and Quality Gates

**Version:** v1.7  
**Date:** 2026-04-22  
**Scope:** Unit tests, integration tests, database tests, CI gates, and release quality checks.

Central command matrix: `docs/05_development_standards/63_Validation_Command_Matrix.md`.

## Mục tiêu

Tài liệu này định nghĩa test nào cần có trước khi merge. Mục tiêu không phải là coverage đẹp, mà là bắt được lỗi thật trong flow tiền, tenant, quyền, provider và provisioning.

## Luật bắt buộc

1. Code không có test phù hợp thì chưa done.
2. Flow tiền phải test idempotency và không double debit.
3. Tenant/RBAC phải test đúng quyền và sai quyền.
4. Provider/provisioning phải test timeout, fail và retry/manual review.
5. Migration phải được kiểm tra trên database dev sạch nếu có thể.
6. CI fail thì không merge.
7. Test flaky phải sửa, không ignore lâu dài.

## Các lớp test

```text
unit test          test service/function nhỏ, nhanh, không cần database thật
integration test   test với database, HTTP server hoặc module thật
contract test      test format API/request/response không đổi ngoài ý định
race test          test cạnh tranh dữ liệu ở flow tiền/ledger
smoke test         test app khởi động và endpoint health cơ bản
manual test        chỉ dùng bổ sung, không thay thế test tự động cho rule P0
```

## Test theo loại thay đổi

```text
docs only                  không bắt buộc test code
config                     test app load config nếu có code config
handler/API                test route chính, status code, response body
service business rule      unit test case chính và case lỗi
store/SQL                  integration test query hoặc repository
migration                  chạy migration trên database dev sạch
platform shared            chạy test toàn repo
worker/job                 test retry, lock, success, fail
provider adapter           test success, fail, timeout, partial
```

## Quality gate tối thiểu

Trước khi mở PR ready:

```bash
gofmt -w <files>
go test ./...
```

Khi có nhiều entrypoint:

```bash
go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker
```

Khi có `Makefile`, dùng:

```bash
make fmt
make test
make build
```

## CI gate tối thiểu

CI nên chạy:

```text
format check
lint
unit test
build all cmd entrypoints
migration check
secret scan
```

CI mở rộng khi có hạ tầng:

```text
PostgreSQL integration test
race test for wallet/ledger
API contract test
Docker image build
dependency vulnerability scan
```

## Wallet và ledger tests

Bắt buộc test:

- Credit thành công.
- Debit thành công.
- Debit khi không đủ số dư.
- Hold/lock tiền nếu flow dùng hold.
- Release hold.
- Refund/reversal.
- Không sửa ledger entry cũ.
- Idempotency key bị gọi lại.
- Retry không tạo thêm debit.
- Concurrent debit không làm số dư âm.

Test nên kiểm cả:

- Số dư cuối.
- Ledger entry được tạo đúng loại.
- Audit event được ghi nếu flow yêu cầu.
- Error code đúng.

## Order và checkout tests

Bắt buộc test:

- Checkout thành công.
- Checkout khi wallet không đủ tiền.
- Checkout khi tenant không hợp lệ.
- Checkout khi product không còn bán.
- Checkout khi provider capability không cho phép action.
- Order snapshot giữ giá/policy tại thời điểm mua.
- Retry checkout không tạo order trùng hoặc debit trùng.

## Tenant và RBAC tests

Bắt buộc test:

- User đúng tenant đọc được dữ liệu.
- User sai tenant không đọc được dữ liệu.
- User thiếu permission bị chặn.
- User đủ permission được phép.
- Admin/reseller/client khác nhau đúng behavior.
- Query list không leak dữ liệu tenant khác.

## Provider và provisioning tests

Bắt buộc test:

- Provider success.
- Provider validation fail.
- Provider timeout.
- Provider network fail.
- Provider partial success nếu adapter hỗ trợ.
- Retry với lỗi retryable.
- Manual review với lỗi không chắc trạng thái.
- Không retry khi không biết provider đã tạo resource hay chưa.
- Credential không xuất hiện trong log/error.

## Worker và scheduler tests

Bắt buộc test:

- Job được claim một lần.
- Job success đổi trạng thái đúng.
- Job fail tăng retry count.
- Retry backoff đúng.
- Job quá số lần retry vào manual review hoặc dead state.
- Worker restart không mất job.
- Outbox event không publish trùng ngoài chính sách cho phép.

## Database tests

Migration test nên kiểm:

- Database sạch chạy lên được.
- Bảng/cột/index/constraint đúng.
- Down migration hoặc rollback plan rõ.
- Seed local không tạo dữ liệu production giả.
- Constraint tiền/tenant quan trọng không bị thiếu.

Query test nên kiểm:

- Query có tenant filter.
- Query list có pagination.
- Query update dùng điều kiện đủ chặt.
- Transaction rollback khi một bước fail.

## API contract tests

API test nên kiểm:

- Success envelope.
- Error envelope.
- Validation error format.
- Pagination format.
- Request id xuất hiện trong response/log nếu có.
- Error code không đổi ngoài ý định.

Không để mỗi module tự phát minh response format.

## Test data

Test data phải rõ nghĩa:

```text
tenantAdmin
resellerTenant
clientUser
walletWithEnoughBalance
walletWithLowBalance
retryableProviderTimeout
```

Tránh tên mơ hồ:

```text
test1
userA
data
mockObj
```

## Mock và fake

Dùng fake khi cần kiểm behavior rõ. Dùng mock khi cần kiểm một call cụ thể.

Provider adapter nên có fake provider để test:

- success
- fail
- timeout
- partial success
- slow response

Không dùng mock để che giấu business rule chưa test.

## Race và concurrency

Chạy race test cho flow tiền khi có code thật:

```bash
go test -race ./internal/modules/wallet/...
go test -race ./internal/modules/ledger/...
go test -race ./internal/modules/order/...
```

Concurrency test nên tập trung vào:

- debit cùng ví cùng lúc
- checkout retry song song
- worker claim cùng job
- provider reconcile chạy cùng retry

## Coverage rule

Không đặt coverage làm mục tiêu duy nhất.

Ưu tiên test các rule:

- tiền không sai
- tenant không leak
- quyền không bypass
- credential không lộ
- provider retry an toàn
- ledger không bị sửa lịch sử

Coverage thấp ở code đơn giản có thể chấp nhận. Thiếu test ở rule P0 thì không chấp nhận.

## Flaky test

Nếu test flaky:

- Tìm nguyên nhân.
- Tách dependency thời gian bằng clock fake.
- Tách network bằng fake provider.
- Tránh sleep cố định nếu có thể.
- Không mark skip dài hạn nếu không có issue follow-up.

## Test naming

Tên test nên nói rõ case:

```text
TestDebitWalletRejectsInsufficientBalance
TestCheckoutDoesNotDoubleDebitOnRetry
TestTenantUserCannotReadOtherTenantOrders
TestProviderTimeoutMovesJobToManualReview
```

Không dùng:

```text
TestWallet
TestCheckout2
TestError
```

## Definition of done cho testing

Một thay đổi đạt testing gate khi:

- Test liên quan đã được thêm hoặc cập nhật.
- `go test ./...` pass.
- Entry point bị ảnh hưởng build được.
- Migration được kiểm nếu có.
- Test P0 pass cho flow tiền/tenant/quyền/provider.
- PR ghi rõ lệnh test đã chạy.
