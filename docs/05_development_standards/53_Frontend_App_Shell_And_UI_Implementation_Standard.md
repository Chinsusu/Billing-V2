# Frontend App Shell and UI Implementation Standard

**Version:** v1.9
**Date:** 2026-04-25
**Scope:** Frontend app shell, runnable UI deliverables, mock data, navigation, screen structure, and build validation.

## Mục tiêu

Tài liệu này khóa yêu cầu tối thiểu cho frontend implementation. Mục tiêu là agent không chỉ nộp một file HTML tĩnh, mà phải dựng một app frontend chạy được, có structure rõ, có navigation, có build script và có thể phát triển tiếp.

Phase này chưa cần wire backend route/API thật. Nhưng UI phải là app shell hoạt động được với mock data.

## Frontend architecture decision

UI mặc định dùng:

```text
Next.js App Router
React
TypeScript
Tailwind CSS
```

Node.js chỉ được dùng làm runtime/toolchain cho frontend: install dependency, chạy dev server, build, preview/start. Node.js không phải backend nghiệp vụ của hệ thống Billing.

Backend nghiệp vụ vẫn là Go API trong repo này. Logic về tiền, tenant, role, order, provisioning, provider, audit, credential và database không được đặt trong Node.js backend riêng, Next API routes hoặc Next Server Actions.

Nếu task muốn đổi framework frontend khác, PR phải ghi rõ lý do và được chấp nhận trước khi merge. Mặc định agent frontend phải dùng Next.js.

Stack phụ trợ được khuyến nghị:

```text
shadcn/ui hoặc component system tương đương
TanStack Table cho bảng phức tạp
Zod cho schema/validation phía frontend
```

## Luật bắt buộc

1. Không chấp nhận chỉ một file HTML tĩnh làm deliverable frontend.
2. Frontend mặc định phải dùng Next.js App Router, React, TypeScript và Tailwind CSS.
3. Frontend phải có `frontend/package.json`.
4. Phải có scripts tối thiểu: `dev`, `build`, `preview`.
5. `preview` có thể alias tới `next start`; nếu dùng thêm `start` thì vẫn phải giữ `preview` cho workflow chung.
6. Phải có app entrypoint thật, không nhét toàn bộ UI vào HTML.
7. Phải có navigation hoạt động giữa các screen bằng Next router hoặc client state rõ ràng.
8. Phải có mock data layer tách riêng khỏi component.
9. Phải có layout shell dùng chung.
10. Mỗi file không vượt 500 dòng.
11. Không hardcode toàn bộ dữ liệu trong component chính.
12. PR phải chạy `npm install` hoặc package-manager tương đương và `npm run build`.

## Deliverable tối thiểu

Frontend app shell tối thiểu gồm:

```text
frontend/
  package.json
  next.config.ts hoặc next.config.mjs
  tsconfig.json
  src/
    app/
      layout.tsx
      page.tsx
      globals.css
    modules/
      admin/
      reseller/
      client/
    components/
      layout/
      ui/
      shared/
      data-display/
    lib/
      api/
      config/
      navigation/
    mocks/
    styles/
```

Với Next.js, cấu trúc gợi ý:

```text
frontend/
  package.json
  next.config.ts
  src/
    app/
      layout.tsx
      page.tsx
      globals.css
    modules/
      admin/
        AdminOverview.tsx
      reseller/
        ResellerDashboard.tsx
      client/
        ClientServices.tsx
    lib/
      navigation/
        screens.ts
      api/
        mockApi.ts
      config/
        env.ts
    mocks/
      billingData.ts
    components/
      layout/
      ui/
      shared/
      data-display/
```

Không tạo cấu trúc SPA/HTML tĩnh nếu task không cho phép rõ. Output phải chạy được bằng script trong `package.json`.

## Scripts bắt buộc

`frontend/package.json` phải có:

```json
{
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "preview": "next start",
    "start": "next start"
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
- Node.js backend riêng.
- Next API routes hoặc Server Actions cho nghiệp vụ Billing.
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

Không tạo Node.js backend riêng cho frontend app shell. Không dùng Express, Nest, Fastify hoặc Next API routes để xử lý nghiệp vụ tiền, tenant, order, provisioning, provider, credential hoặc audit.

Nếu cần giả lập API, tạo adapter/mock layer:

```text
frontend/src/lib/api/mockApi.ts
```

Không gọi production endpoint.

Không hardcode future API response format trái với:

```text
docs/05_development_standards/50_API_Response_Error_Logging_Standard.md
docs/02_technical_handoff/16_API_Contract_And_Permission_Spec.md
docs/05_development_standards/64_Public_Display_ID_And_Backend_Reference_Policy.md
```

Visible resource labels must use public numeric IDs from `display_id` or related `*_display_id` fields. Do not render backend UUID references as row labels, card titles, search suggestions, or mock examples. If the API does not provide the needed public ID yet, show `not shown` and create or link a backend follow-up.

## Build validation

Use `docs/05_development_standards/63_Validation_Command_Matrix.md` for the exact frontend validation set by change type.

PR frontend app shell phải ghi rõ lệnh đã chạy:

```bash
cd frontend
npm install
npm run build
```

Truoc khi mo PR frontend, chay static guard de chan field backend nhay cam bi dua vao UI copy hoac mock data:

```bash
cd frontend
npm run check:sensitive-text
```

Guard nay chi cho phep ngoai le hep trong API type definitions va explicit redaction tests. Khong render hoac hardcode cac field nhu `payload_json`, `capability_profile`, `provider_account_id`, `secret`, `raw_response`, credential/token variants trong component hoac mock data.

`npm run preview` dùng để smoke test thủ công sau build. Không dùng lệnh preview như CI gate nếu command giữ server chạy liên tục.

Khi task thay đổi admin navigation, admin screen, API adapter, mock data, hoặc response mapping, chạy thêm browser smoke:

```bash
cd frontend
npm run smoke:admin
```

CI smoke runs after the production build and uses the standalone server artifact:

```bash
cd frontend
npx playwright install --with-deps chromium
npm run build
npm run smoke:admin:ci
```

Do not run `npm run smoke:admin`, `npm run smoke:admin:ci`, or `npm run build` in parallel because they read or write `.next`.

Browser smoke dùng dữ liệu mock/intercept an toàn, không cần backend thật hoặc provider credential thật. Nếu command fail vì thiếu browser runtime trên máy mới, cài browser Playwright local bằng `npx playwright install chromium`, rồi chạy lại smoke.

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
- Frontend có dùng Next.js App Router, React, TypeScript và Tailwind CSS chưa?
- Node.js có chỉ được dùng cho frontend toolchain, không làm backend nghiệp vụ chưa?
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
- `frontend/` có Next.js app shell chạy được.
- Node.js chỉ dùng cho frontend toolchain.
- Build pass.
- Navigation hoạt động.
- Mock data tách riêng.
- Component/screen structure rõ.
- Không file nào vượt 500 dòng.
- Không có secret hoặc dữ liệu khách hàng thật.
- `TASKS.md` được cập nhật.
