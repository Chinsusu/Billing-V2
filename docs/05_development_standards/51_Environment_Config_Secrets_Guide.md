# Environment, Config, and Secrets Guide

**Version:** v1.7  
**Date:** 2026-04-22  
**Scope:** Local/dev/staging/production config, environment variables, secrets, config validation, and rotation rules.

## Mục tiêu

Tài liệu này khóa cách cấu hình app trước khi bootstrap code. Mục tiêu là app khởi động rõ ràng, thiếu config thì fail sớm, và secret không lọt vào Git, log hoặc audit.

## Luật bắt buộc

1. Không commit `.env` thật.
2. Chỉ commit `.env.example`.
3. Secret không được ghi vào log, audit hoặc error response.
4. App phải validate config khi khởi động.
5. Config production không được nằm trong repo.
6. Default config chỉ dùng cho local/dev an toàn.
7. Thay đổi config bắt buộc phải cập nhật `.env.example` và tài liệu liên quan.

## File config

Repo chỉ được commit file mẫu:

```text
.env.example
config.example.yaml nếu sau này cần
```

Không commit:

```text
.env
.env.local
.env.production
*.pem
*.key
provider credential
database dump
```

## Environment names

Dùng các environment:

```text
local      máy dev
dev        môi trường dev shared
staging    môi trường kiểm thử gần production
production môi trường thật
```

Không dùng tên mơ hồ như `test2`, `live1`, `real`, `temp-prod`.

## Env variable naming

Tên env dùng uppercase và prefix rõ:

```text
APP_ENV
APP_NAME
APP_HTTP_ADDR
DB_DSN
LOG_LEVEL
JWT_SECRET
ENCRYPTION_KEY
AUTH_SESSION_COOKIE_NAME
AUTH_SESSION_COOKIE_SECURE
AUTH_SESSION_TTL
TELEGRAM_BOT_TOKEN
SMTP_HOST
SMTP_USERNAME
SMTP_PASSWORD
```

Rule:

- Tên nói rõ mục đích.
- Không viết tắt khó hiểu.
- Secret có hậu tố hoặc tên rõ như `SECRET`, `TOKEN`, `PASSWORD`, `KEY`.
- Không dùng một biến cho nhiều ý nghĩa.

## Config required và optional

Config required là config thiếu thì app không được boot.

Ví dụ:

```text
APP_ENV
DB_DSN
LOG_LEVEL
ENCRYPTION_KEY
```

Config optional phải có default an toàn.

Ví dụ:

```text
APP_HTTP_ADDR=:8080
AUTH_SESSION_COOKIE_NAME=billing_session
AUTH_SESSION_COOKIE_SECURE=true
AUTH_SESSION_TTL=12h
METRICS_ENABLED=false
TRACING_ENABLED=false
```

Không dùng default production secret.

## Config validation

App phải validate:

- `APP_ENV` thuộc danh sách hợp lệ.
- `DB_DSN` không rỗng.
- `LOG_LEVEL` hợp lệ.
- Secret required đủ dài.
- URL/host/port đúng format.
- Production không dùng default local.

Nếu config sai, app fail khi boot với error rõ nhưng không lộ secret.

## `.env.example`

`.env.example` phải:

- Có đủ key required.
- Dùng placeholder an toàn.
- Không chứa secret thật.
- Có comment ngắn cho key dễ nhầm.
- Được cập nhật trong cùng PR khi thêm config mới.

Ví dụ secret placeholder:

```text
JWT_SECRET=change-me-local-only
ENCRYPTION_KEY=change-me-32-byte-local-only
AUTH_SESSION_COOKIE_NAME=billing_session
AUTH_SESSION_COOKIE_SECURE=false
AUTH_SESSION_TTL=12h
AUTH_PASSWORD_RESET_TTL=30m
```

## Local development

Local dev được phép dùng:

- Database local.
- Secret giả.
- Provider fake/sandbox.
- Email/Telegram disabled hoặc sandbox.

Local không được gọi production provider hoặc production database.

## Staging

Staging nên giống production về:

- Migration flow.
- Config key.
- Worker/scheduler topology.
- Logging format.
- Provider sandbox nếu có.

Staging không dùng dữ liệu production thật nếu chưa được ẩn danh và cho phép.

## Production

Production config phải đến từ secret manager, CI/CD secret store hoặc hệ thống vận hành được duyệt.

Production không được:

- Load `.env` từ repo.
- Dùng secret local.
- Bật debug log dài hạn.
- Dùng provider sandbox trừ khi maintenance có chủ ý.

## Secret handling

Secret bao gồm:

```text
password
token
api key
private key
JWT secret
encryption key
provider credential
SMTP password
Telegram token
database password
```

Rule:

- Không log secret.
- Không trả secret qua API.
- Không ghi secret plaintext vào audit.
- Không commit secret.
- Không paste secret vào issue/PR.
- Không gửi secret qua chat nếu có kênh secret manager.

## Secret rotation

Cần rotate secret khi:

- Secret lộ trong Git hoặc log.
- Nhân sự không còn quyền truy cập.
- Provider báo leak.
- Theo lịch bảo mật định kỳ.

Rotation plan cần có:

- Secret nào đổi.
- Service nào bị ảnh hưởng.
- Thời điểm đổi.
- Cách verify.
- Cách rollback nếu app không boot.

## Config ownership

Owner theo nhóm:

```text
APP_*          internal/app hoặc platform/config
DB_*           platform/db
LOG_*          platform/logger
JWT_*          identity/auth
AUTH_SESSION_* identity/auth
ENCRYPTION_*   platform/crypto
SMTP_*         platform/email
TELEGRAM_*     platform/telegram
PROVIDER_*     modules/provider
```

Không thêm config mới nếu chưa rõ owner.

## Logging config

Không log raw config.

Khi boot, chỉ log config an toàn:

```text
app_env
app_name
http_addr
log_level
metrics_enabled
tracing_enabled
```

Không log:

```text
DB_DSN đầy đủ có password
JWT_SECRET
ENCRYPTION_KEY
SMTP_PASSWORD
PROVIDER_API_KEY
```

## Config change checklist

Khi thêm hoặc đổi config:

- `.env.example` đã cập nhật chưa?
- Config validation đã cập nhật chưa?
- README/local runbook cần cập nhật không?
- Secret có bị log không?
- Production có cần migration secret không?
- Default có an toàn không?
- CI cần thêm secret không?

## Incident khi lộ secret

Nếu secret lộ:

1. Báo owner ngay.
2. Rotate secret.
3. Revoke secret cũ.
4. Kiểm tra log/Git/CI artifact.
5. Tạo post-incident note nếu ảnh hưởng production.

Không chỉ xóa secret khỏi file rồi coi là xong.
