# Frontend App Shell and UI Implementation Standard

**Version:** v1.8  
**Date:** 2026-04-22  
**Scope:** Frontend app shell, runnable UI deliverables, mock data, navigation, screen structure, and build validation.

## Mục tiêu

Tài liệu này khóa yêu cầu tối thiểu cho frontend implementation. Mục tiêu là agent không chỉ nộp một file HTML tĩnh, mà phải dựng một app frontend chạy được, có structure rõ, có navigation, có build script và có thể phát triển tiếp.

Phase này chưa cần wire backend route/API thật. Nhưng UI phải là app shell hoạt động được với mock data.

## Luật bắt buộc

1. Không chấp nhận chỉ một file HTML tĩnh làm deliverable frontend.
2. Frontend phải có `frontend/package.json`.
3. Phải có scripts tối thiểu: `dev`, `build`, `preview`.
4. Phải có app entrypoint thật, không nhét toàn bộ UI vào HTML.
5. Phải có navigation hoạt động giữa các screen bằng client state hoặc router nội bộ.
6. Phải có mock data layer tách riêng khỏi component.
7. Phải có layout shell dùng chung.
8. Mỗi file không vượt 500 dòng.
9. Không hardcode toàn bộ dữ liệu trong component chính.
10. PR phải chạy `npm install` hoặc package-manager tương đương và `npm run build`.

## Deliverable tối thiểu

Frontend app shell tối thiểu gồm:

```text
frontend/
  package.json
  index.html
  src/ hoặc lib/
    app entrypoint
    app shell/layout
    navigation/screen registry
    mock data
    shared components
    screens
    styles/tokens
```

Nếu dùng Vite/React, cấu trúc gợi ý:

```text
frontend/
  package.json
  index.html
  src/
    main.jsx
    app.jsx
    routes.jsx
    data/
      mock-billing-data.js
    components/
      layout/
      ui/
      shared/
    screens/
      admin/
      reseller/
      client/
    styles/
      tokens.css
      app.css
```

Không bắt buộc đúng framework trên, nhưng output phải chạy được bằng script trong `package.json`.

## Scripts bắt buộc

`frontend/package.json` phải có:

```json
{
  "scripts": {
    "dev": "...",
    "build": "...",
    "preview": "..."
  }
}
```

Nếu có lint/check tool, thêm:

```json
{
  "scripts": {
    "check": "...",
    "lint": "..."
  }
}
```

Nếu không có lint ở phase đầu, PR phải ghi rõ chỉ chạy `build`.

## Không cần trong phase này

Phase app shell chưa cần:

- Auth thật.
- Backend API thật.
- Server-side routing.
- Production deploy.
- Payment/provider action thật.
- Persistent state.
- Database.

Nhưng mock flow phải đủ để QA và product nhìn được screen/state chính.

## Navigation

Navigation phải hoạt động trong app.

Tối thiểu cần:

- Sidebar hoặc top navigation.
- Active screen state.
- Click menu đổi screen.
- Màn admin/reseller/client nếu UI scope có nhiều portal.
- Không reload cả page khi đổi screen nếu framework hỗ trợ client state.

Không để toàn bộ screen nằm nối tiếp nhau trong một HTML dài mà không có navigation.

## Screen registry

App nên có screen registry rõ:

```text
id
label
portal
permission hoặc role hint nếu có
component
```

Ví dụ:

```js
export const screens = [
  { id: "admin-overview", label: "Overview", portal: "admin", component: AdminOverview },
  { id: "reseller-billing", label: "Billing", portal: "reseller", component: ResellerBilling },
  { id: "client-services", label: "Services", portal: "client", component: ClientServices }
]
```

Không hardcode navigation ở nhiều file nếu cùng một list screen được dùng nhiều nơi.

## Mock data

Mock data phải nằm ở file riêng như:

```text
frontend/src/data/mock-billing-data.js
```

Mock data nên mô phỏng:

- tenant/reseller/client
- wallet balance
- ledger entries
- orders
- services
- provider status
- notifications hoặc audit summary nếu screen cần

Không đặt dữ liệu mẫu lớn trực tiếp trong component chính.

Không dùng secret thật, IP nhạy cảm, customer data thật hoặc provider credential thật trong mock data.

## Component structure

Component dùng một screen thì để trong folder screen đó.

Component dùng nhiều screen thì đưa vào:

```text
components/shared/
```

Primitive UI component đưa vào:

```text
components/ui/
```

Layout đưa vào:

```text
components/layout/
```

Không tạo component tên mơ hồ như:

```text
Box
Thing
CommonTable
BaseCard
Panel
```

Dùng tên cụ thể hơn:

```text
WalletBalanceCard
OrderStatusBadge
TenantSwitcher
LedgerEntryTable
ProviderCapabilityTable
```

## Styling

Styling phải có structure rõ:

- token hoặc CSS variables cho màu/spacing/font nếu có nhiều screen
- responsive desktop/mobile tối thiểu
- state hover/focus/active cho controls chính
- không để text tràn button/card
- không dùng palette một màu đơn điệu nếu không có chủ ý

Không dùng inline style khắp nơi nếu có thể tách thành CSS/module/theme file.

## States tối thiểu

Screen có data table/list nên có:

- normal state
- empty state
- loading placeholder hoặc loading state
- error state hoặc disabled action hint

Phase mock chưa cần fetch thật, nhưng component nên có đường để truyền state sau này.

## API boundary

Không wire backend route thật ở phase app shell.

Nếu cần giả lập API, tạo adapter/mock layer:

```text
frontend/src/data/mock-api.js
```

Không gọi production endpoint.

Không hardcode future API response format trái với:

```text
docs/05_development_standards/50_API_Response_Error_Logging_Standard.md
docs/02_technical_handoff/16_API_Contract_And_Permission_Spec.md
```

## Build validation

PR frontend app shell phải ghi rõ lệnh đã chạy:

```bash
cd frontend
npm install
npm run build
```

Nếu dùng package manager khác, ghi rõ:

```bash
pnpm install
pnpm build
```

Build fail thì không merge.

Nếu không thể chạy build vì thiếu tool/hạ tầng, task phải chuyển `BLOCKED` và ghi lý do vào `TASKS.md`.

## Cấm

Không được:

- nộp một file HTML tĩnh rồi coi là app frontend
- nhét toàn bộ UI vào một file trên 500 dòng
- dùng dữ liệu thật hoặc secret thật
- bỏ qua package scripts
- bỏ qua build validation
- tạo UI không có navigation
- tạo screen nhưng không có cách truy cập từ app shell
- tạo code frontend không thể phát triển tiếp

## PR checklist

Trước khi mở PR frontend app shell:

- `frontend/package.json` có script `dev`, `build`, `preview` chưa?
- App chạy được chưa?
- Navigation đổi screen được chưa?
- Mock data tách riêng chưa?
- Layout shell dùng chung chưa?
- File nào vượt 500 dòng không?
- Build đã pass chưa?
- `TASKS.md` đã cập nhật trạng thái task chưa?
- PR body có ghi lệnh build không?

## Definition of done

Frontend app shell task done khi:

- PR đã merge vào `main`.
- `frontend/` có app chạy được.
- Build pass.
- Navigation hoạt động.
- Mock data tách riêng.
- Component/screen structure rõ.
- Không file nào vượt 500 dòng.
- Không có secret hoặc dữ liệu khách hàng thật.
- `TASKS.md` được cập nhật.
