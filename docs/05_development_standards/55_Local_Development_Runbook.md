# Local Development Runbook

**Version:** v1.10  
**Date:** 2026-04-22  
**Scope:** Local backend setup, environment loading, API run commands, migration runner usage, and local validation gates.

## Mục tiêu

Tài liệu này mô tả cách chạy backend Go trên máy local bằng các script hiện có trong repo. Mục tiêu là developer mới có thể boot API, kiểm tra health/readiness, validate migration và chạy quality gate mà không cần secret thật hoặc môi trường production.

## Luật bắt buộc

1. Không commit `.env` thật.
2. Không dùng production database, production provider credential hoặc production webhook khi chạy local.
3. Local chỉ dùng secret giả, provider fake/sandbox và database local.
4. Khi thêm env mới, cập nhật `.env.example` và tài liệu config liên quan trong cùng PR.
5. Trước khi mở PR backend phải chạy `make test` và `make build`.
6. Nếu thay đổi migration, phải chạy thêm `make migrate-validate`.

## Yêu cầu máy local

Cần có:

```text
Go 1.18 hoặc tương thích với go.mod
make
Git
PostgreSQL local nếu muốn chạy migrate plan/up với database thật
```

Không cần PostgreSQL để chạy unit test hiện tại, build binary hoặc validate format migration.

## Chuẩn bị workspace

Luôn bắt đầu từ `main` mới nhất và tạo branch riêng theo task:

```bash
git switch main
git pull --ff-only origin main
git switch -c <type>/<short-task-name>
```

Ví dụ:

```bash
git switch -c feat/identity-tenant-rbac-skeleton
```

Nếu dùng nhiều agent hoặc nhiều task song song, ưu tiên dùng `git worktree` để mỗi task có workspace riêng, tránh sửa cùng file không liên quan.

## Tạo file env local

Tạo `.env` từ file mẫu:

```bash
cp .env.example .env
```

Giữ giá trị local an toàn. Ví dụ mặc định hợp lệ:

```text
APP_ENV=local
APP_NAME=billing-v2
APP_HTTP_ADDR=:8080
LOG_LEVEL=debug
DB_DSN=postgres://billing:billing@localhost:5432/billing?sslmode=disable
JWT_SECRET=change-me-local-only
ENCRYPTION_KEY=change-me-32-byte-local-only
PROVIDER_DEFAULT_MODE=fake
```

Không thay placeholder bằng secret thật trong repo. Nếu cần dùng credential sandbox, lưu trong `.env` local hoặc secret manager được duyệt, không paste vào log, task, issue hoặc PR.

## Load env cho shell

Ứng dụng đọc config từ environment variable. Trong terminal local có thể load `.env` bằng:

```bash
set -a
. ./.env
set +a
```

Sau khi load, kiểm tra nhanh:

```bash
printenv APP_ENV
printenv APP_HTTP_ADDR
```

Không in các biến secret như `JWT_SECRET`, `ENCRYPTION_KEY`, token provider, SMTP password hoặc database password vào log chia sẻ.

## Chạy API local

Chạy API:

```bash
make run-api
```

API mặc định listen theo `APP_HTTP_ADDR`; với `.env.example` là `:8080`.

Kiểm tra health:

```bash
curl -i http://localhost:8080/healthz
```

Kiểm tra readiness:

```bash
curl -i http://localhost:8080/readyz
```

Response success dùng envelope chuẩn của `internal/platform/httpserver`. Nếu port `8080` đã bận, đổi `APP_HTTP_ADDR` trong `.env`, load lại env rồi chạy lại API.

## Migration runner

Migration runner nằm ở `cmd/migrate` và đọc migration từ thư mục `migrations` theo mặc định.

Validate migration file mà không cần database:

```bash
make migrate-validate
```

Lệnh tương đương:

```bash
go run ./cmd/migrate validate
```

Xem plan không cần database:

```bash
go run ./cmd/migrate plan
```

Khi `DB_DSN` rỗng, `plan` chỉ in danh sách migration có trong repo. Ở trạng thái skeleton hiện tại, repo có thể chưa có file `.sql` nên kết quả hợp lệ có thể là `0 migration(s)`.

## Chạy migration với PostgreSQL local

Chỉ dùng database local hoặc sandbox được duyệt. Nếu muốn dùng đúng DSN trong `.env.example`, tạo role và database local tương ứng:

```bash
psql -d postgres -c "CREATE ROLE billing WITH LOGIN PASSWORD 'billing';"
psql -d postgres -c "CREATE DATABASE billing OWNER billing;"
```

Nếu máy local đã có user/database khác, cập nhật `DB_DSN` trong `.env` theo user/database đó.

Set DSN local cho terminal hiện tại:

```bash
export DB_DSN='postgres://billing:billing@localhost:5432/billing?sslmode=disable'
```

Xem pending migration trên database:

```bash
go run ./cmd/migrate plan
```

Apply migration:

```bash
go run ./cmd/migrate up
```

`up` bắt buộc có `DB_DSN` hoặc flag `-dsn`. Nếu cần override thư mục hoặc timeout:

```bash
go run ./cmd/migrate -dir migrations -timeout 30s plan
go run ./cmd/migrate -dsn "$DB_DSN" up
```

Không chạy `up` vào staging/production từ máy local nếu không có quy trình vận hành và approval rõ.

## Seed dữ liệu dev

Sau khi migration đã chạy trên PostgreSQL local, tạo dữ liệu demo:

```bash
go run ./cmd/seed dev
```

Seed local hiện có các actor mẫu:

```text
admin@local.billing      Platform admin
reseller@local.billing   Demo reseller owner
customer@local.billing   Demo customer
```

Billing flow mẫu tạo wallet, top-up đã duyệt, order đã paid, service instance, invoice đã paid, wallet ledger debit và payment transaction liên kết với nhau bằng UUID cố định và `display_id` dạng số. Các API mẫu sau dùng header local, không dùng credential thật:

```bash
curl -H "X-Tenant-Id: 00000000-0000-0000-0000-000000000010" \
  -H "X-Actor-Id: 00000000-0000-0000-0000-000000000102" \
  -H "X-Actor-Tenant-Id: 00000000-0000-0000-0000-000000000010" \
  -H "X-Actor-Type: reseller_owner" \
  http://localhost:8080/admin/payment-reconciliation

curl -H "X-Tenant-Id: 00000000-0000-0000-0000-000000000010" \
  -H "X-Actor-Id: 00000000-0000-0000-0000-000000000103" \
  -H "X-Actor-Tenant-Id: 00000000-0000-0000-0000-000000000010" \
  -H "X-Actor-Type: client" \
  http://localhost:8080/client/wallets
```

## Smoke test database dev

Sau khi có PostgreSQL local/dev và `DB_DSN`, chạy smoke để kiểm tra migration, seed, và billing flow mẫu trên database thật:

```bash
make smoke-dev-db
```

Lệnh này tương đương:

```bash
go run ./cmd/smoke dev-db
```

Smoke command sẽ:

- Từ chối chạy nếu `APP_ENV=production` hoặc `APP_ENV=prod`.
- Apply toàn bộ migration còn thiếu.
- Chạy seed dev hai lần để kiểm tra idempotency.
- Kiểm tra các record mẫu cho tenant, user, permission, catalog, wallet, top-up, order, service, invoice, ledger và payment.

Nếu cần truyền DSN trực tiếp:

```bash
go run ./cmd/smoke -dsn "$DB_DSN" dev-db
```

Chỉ chạy lệnh này trên database local hoặc sandbox được phép. Không dùng production DSN.

## Smoke test API billing

Sau khi `make smoke-dev-db` pass và API đang chạy với cùng `DB_DSN`, kiểm tra các endpoint billing đọc dữ liệu seed:

```bash
make smoke-dev-api
```

Lệnh này tương đương:

```bash
go run ./cmd/smoke dev-api
```

Mặc định smoke gọi `http://localhost:8080`. Nếu API chạy ở địa chỉ khác:

```bash
go run ./cmd/smoke -base-url "http://localhost:8081" dev-api
```

Smoke API dùng actor local:

```text
reseller@local.billing   Admin/read checks
customer@local.billing   Client/read checks
```

Các check bao gồm health, readiness, wallet, ledger, order, service, invoice, payment transaction, payment reconciliation, top-up request và audit list. Lệnh chỉ dùng header local/dev, không dùng token hoặc credential thật.

## Smoke test billing mutation flow

Sau khi `make smoke-dev-db` pass va API dang chay voi cung `DB_DSN`, chay mutation smoke de kiem tra billing flow that:

```bash
make smoke-dev-billing
```

Lenh nay tuong duong:

```bash
go run ./cmd/smoke dev-billing
```

Flow duoc test:

- client tao top-up request, admin approve top-up;
- client tao order voi `Idempotency-Key`;
- client goi `POST /client/checkouts` voi `order_id` de lay invoice `issued`;
- smoke goi lai checkout cung idempotency key de kiem tra duplicate submit khong tao invoice moi;
- client tra invoice bang `POST /client/invoice-wallet-payments`;
- payment finalizes the order and creates or reuses one `provider.provision` job for that order;
- smoke doc lai `/client/orders/{order_id}` de xac nhan order thanh `order_status=paid` va `billing_status=paid`;
- smoke kiem tra bang `jobs` co dung mot `provider.provision` job cho order vua tra tien;
- smoke chay fake-provider provisioning worker trong process de xu ly job vua tao;
- smoke doc lai `/client/services?order_id=...` de xac nhan service `active/paid` duoc tao dung order;
- smoke doc lai invoice va audit log de xac nhan flow co the debug duoc.

Inspect provisioning jobs when an order is paid but fulfillment is stuck:

```sql
SELECT display_id, job_type, reference_type, reference_id, source_id, status, attempt_count, last_error_code, updated_at
FROM jobs
WHERE job_type = 'provider.provision'
ORDER BY created_at DESC
LIMIT 20;
```

If payment returns `order.provisioning_source_not_found`, check that the order's tenant plan points to a master plan with an active `plan_sources` row and an active `provider_sources` row.

For deeper billing operations checks, use `docs/05_development_standards/57_Billing_Operations_Runbook.md`.

Yeu cau:

- `DB_DSN` tro toi database local/dev da seed;
- `API_BASE_URL` mac dinh la `http://localhost:8080`, co the doi bang `-base-url`;
- API dang chay phai dung cung database voi `DB_DSN`, vi smoke vua goi API vua doc bang `jobs` va chay worker tren cung database;
- provider local dung fake registry, khong can provider credential that;
- chi chay tren local hoac sandbox, khong chay voi production DSN.

## Local provisioning worker

Khi can xu ly job `provider.provision` tren database local/dev, dung worker fake-provider trong `cmd/worker`.

Chay mot pass de claim va xu ly batch hien tai:

```bash
go run ./cmd/worker provision-once -dsn "$DB_DSN"
```

Chay loop local/sandbox de tiep tuc polling job moi:

```bash
go run ./cmd/worker provision-loop -dsn "$DB_DSN" -interval 5s -batch-size 10
```

Loop se in summary tung pass theo cac count `claimed`, `succeeded`, `retried`, `manual_review`, `terminal_failed`, va `cancelled`. Khi khong claim duoc job nao, worker cho het `-interval` truoc khi thu lai de tranh busy-spin. Dung `Ctrl+C` de dung loop, hoac them `-timeout 5m` cho run gioi han.

Khong chay worker local voi `APP_ENV=prod` hoac `APP_ENV=production`. Provider local dung fake registry, khong can provider credential that.

## Quality gate trước PR

Chạy format:

```bash
make fmt
```

Chạy test:

```bash
make test
```

Build binary API và migration runner:

```bash
make build
```

Validate migration:

```bash
make migrate-validate
```

Frontend local validation for the runnable Next.js app:

```bash
cd frontend
npm ci
npm audit --omit=dev
npm run lint
npm run build
```

Dung `npm ci` de giong CI hon `npm install`. Neu frontend PR chi doi UI hoac CI, van nen chay du 4 lenh tren truoc khi mo PR.

Với PR chỉ sửa docs, vẫn nên chạy tối thiểu `make test` và `make build` nếu acceptance criteria của task yêu cầu backend còn build được.

## Lỗi thường gặp

`APP_HTTP_ADDR is required` hoặc `APP_HTTP_ADDR is invalid`:

- Kiểm tra `.env`.
- Load lại env bằng `set -a; . ./.env; set +a`.
- Dùng format hợp lệ như `:8080` hoặc `127.0.0.1:8080`.

Port đã bận:

- Đổi `APP_HTTP_ADDR`, ví dụ `:8081`.
- Load lại env rồi chạy lại `make run-api`.

`DB_DSN or -dsn is required for up`:

- Set `DB_DSN` local.
- Hoặc truyền `-dsn` trực tiếp cho `go run ./cmd/migrate`.

Migration validate fail:

- Kiểm tra tên file theo chuẩn `0001_descriptive_name.sql`.
- Không sửa migration đã chạy ở shared environment.
- Thêm migration mới để sửa migration cũ.

Không kết nối được PostgreSQL:

- Kiểm tra PostgreSQL đang chạy.
- Kiểm tra database/user/password trong `DB_DSN`.
- Chỉ dùng database local hoặc sandbox.

## Checklist trước khi mở PR

- Branch tạo từ `main` mới nhất.
- `.env` không bị commit.
- Không có secret thật trong diff, log, task file hoặc PR description.
- `make fmt` đã chạy nếu có đổi Go code.
- `make test` pass.
- `make build` pass.
- `make migrate-validate` pass nếu có đổi migration hoặc runner.
- Task file trong `tasks/active/` có Agent Log và trạng thái đúng theo workflow nhiều agent.
