# API Response, Error Code, and Logging Standard

**Version:** v1.7  
**Date:** 2026-04-22  
**Scope:** HTTP API response format, error codes, validation errors, pagination, request tracing, logging, and redaction.

## Mục tiêu

Tài liệu này khóa format API và log trước khi backend dev. Mục tiêu là mọi module trả response, error và log cùng một kiểu, để frontend, QA và support không phải đoán.

## Luật bắt buộc

1. Không module nào tự phát minh response envelope riêng.
2. Error trả ra user phải có `code`, `message`, `request_id`.
3. Log nội bộ phải có `request_id` nếu request đi qua HTTP.
4. Dữ liệu tenant phải log kèm `tenant_id` khi có.
5. Không log secret, password, token, API key, private key hoặc credential provider.
6. Validation error phải chỉ rõ field lỗi.
7. Error code phải ổn định để frontend và support dùng được.

## Success response

Response thành công dùng format:

```json
{
  "data": {
    "id": "ord_123",
    "status": "paid"
  },
  "request_id": "req_abc"
}
```

Nếu response là list:

```json
{
  "data": [
    {
      "id": "ord_123",
      "status": "paid"
    }
  ],
  "page": {
    "limit": 20,
    "next_cursor": "cursor_value"
  },
  "request_id": "req_abc"
}
```

Không trả nhiều format khác nhau cho cùng một API family.

## Error response

Response lỗi dùng format:

```json
{
  "error": {
    "code": "wallet.insufficient_balance",
    "message": "Wallet balance is not enough.",
    "details": {
      "required_amount": "100000",
      "currency": "VND"
    }
  },
  "request_id": "req_abc"
}
```

Rule:

- `code` dùng cho máy đọc.
- `message` dùng cho người đọc.
- `details` chỉ chứa dữ liệu an toàn.
- Không đưa stack trace ra response.
- Không đưa provider raw error có secret ra response.

## Validation error

Validation error dùng format:

```json
{
  "error": {
    "code": "validation.failed",
    "message": "Request validation failed.",
    "fields": [
      {
        "field": "email",
        "code": "email.invalid",
        "message": "Email is invalid."
      },
      {
        "field": "amount",
        "code": "amount.must_be_positive",
        "message": "Amount must be greater than zero."
      }
    ]
  },
  "request_id": "req_abc"
}
```

Không trả validation error dạng plain text.

## Pagination

List API mặc định dùng cursor pagination.

Request:

```text
GET /orders?limit=20&cursor=cursor_value
```

Response:

```json
{
  "data": [],
  "page": {
    "limit": 20,
    "next_cursor": null
  },
  "request_id": "req_abc"
}
```

Rule:

- `limit` có default và max.
- Không cho `limit` quá lớn.
- Sort mặc định phải ổn định.
- List theo tenant phải có tenant filter.

## Error code naming

Format:

```text
<module>.<reason>
```

Ví dụ:

```text
validation.failed
auth.unauthorized
rbac.permission_denied
tenant.context_missing
wallet.insufficient_balance
ledger.entry_not_created
order.not_found
provider.timeout
provisioning.manual_review_required
credential.redacted
```

Rule:

- Dùng chữ thường.
- Dùng dấu `_` trong reason nếu cần.
- Không dùng code chung chung như `error`, `failed`, `bad_request`.
- Không đổi code đã public nếu không có migration plan cho frontend/support.

## HTTP status mapping

```text
400  validation.failed, request malformed
401  auth.unauthorized
403  rbac.permission_denied, tenant.access_denied
404  resource not found trong phạm vi tenant được phép
409  conflict, idempotency conflict, state conflict
422  business rule rejected
429  rate limit
500  lỗi hệ thống không mong muốn
502  provider bad response
503  provider/service unavailable
504  provider timeout
```

Không dùng `500` cho lỗi business rule đã biết.

## Request id

Mỗi request phải có `request_id`.

Rule:

- Nếu client gửi `X-Request-ID`, server có thể dùng nếu hợp lệ.
- Nếu không có, server tự tạo.
- Response luôn trả lại `request_id`.
- Log trong cùng request dùng cùng `request_id`.

Header:

```text
X-Request-ID: req_abc
```

## Logging fields

Log nên là structured log.

Field chuẩn:

```text
level
message
request_id
tenant_id
user_id
role
module
operation
order_id
wallet_id
ledger_entry_id
provider_id
job_id
error_code
duration_ms
```

Không phải log nào cũng cần đủ field. Nhưng nếu field có liên quan, dùng tên chuẩn.

## Log levels

```text
debug  thông tin dev/debug, không bật nhiều ở production
info   flow bình thường đáng ghi nhận
warn   lỗi recover được hoặc cần chú ý
error  lỗi request/job thất bại
fatal  app không thể tiếp tục chạy
```

Ví dụ:

- Checkout thành công: `info`.
- Provider timeout retryable: `warn`.
- Ledger entry tạo fail: `error`.
- Config thiếu khi app boot: `fatal`.

## Redaction

Các field phải redact:

```text
password
token
api_key
secret
private_key
credential
proxy_password
authorization
cookie
provider_raw_response nếu có secret
```

Giá trị thay thế:

```text
"[REDACTED]"
```

Không log request/response body toàn bộ nếu có thể chứa secret.

## User message và internal message

User-facing message phải đơn giản:

```text
Wallet balance is not enough.
```

Internal log có thể chi tiết hơn:

```text
wallet debit rejected because available balance is lower than required amount
```

Không trả message nội bộ cho user nếu có thông tin nhạy cảm.

## Audit khác log

Log dùng để vận hành và debug.

Audit dùng để truy vết hành động quan trọng:

- ai làm
- làm gì
- tenant nào
- resource nào
- trước/sau nếu an toàn
- thời điểm nào

Flow tiền, tenant, quyền, credential, provider và provisioning cần audit event nếu hành động đổi trạng thái quan trọng.

## Provider error

Provider error phải được map về error nội bộ:

```text
provider.timeout
provider.unavailable
provider.validation_failed
provider.auth_failed
provider.unknown_state
provisioning.manual_review_required
```

Không trả nguyên văn lỗi provider cho client nếu lỗi đó chứa credential, IP nội bộ hoặc chi tiết vận hành.

## Idempotency error

Khi idempotency key bị dùng lại:

- Nếu request giống nhau: trả kết quả cũ hoặc trạng thái hiện tại.
- Nếu request khác nhau: trả `409 idempotency.conflict`.
- Không tạo ledger/order/provisioning job trùng.

## Checklist

Trước khi merge API mới:

- Success response đúng envelope chưa?
- Error response đúng envelope chưa?
- Validation error có field-level chưa?
- Error code rõ và ổn định chưa?
- Request id có trong response/log chưa?
- Log có tenant id nếu cần chưa?
- Có redact secret chưa?
- Docs/API contract đã cập nhật chưa?
