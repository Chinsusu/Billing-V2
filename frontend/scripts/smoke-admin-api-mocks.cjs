async function installApiMocks(page, options = {}) {
  await page.route("**/backend/**", (route) => {
    const request = route.request();
    const url = new URL(request.url());
    const pathname = url.pathname;
    const query = url.searchParams;
    if (options.failPaths?.has(pathname)) {
      return route.fulfill({
        status: 503,
        contentType: "application/json",
        body: JSON.stringify({ error: { message: "Smoke fallback failure" } }),
      });
    }
    if (request.method() === "POST" && pathname.endsWith("/credentials/credential-1/reveal")) { // sensitive-text-allowlist
      return json(route, {
        id: "credential-1", // sensitive-text-allowlist
        credential_type: "vps_root", // sensitive-text-allowlist
        masked_hint: "root / ****",
        status: "active",
        payload: { username: "root", password: "fixture-access" },
        revealed_at: now(),
        reveal_expires_message: "Shown once. Store it securely before leaving this screen.",
      });
    }
    if (request.method() !== "GET") {
      return json(route, []);
    }

    switch (pathname) {
      case "/backend/admin/wallets":
        return json(route, [
          { id: "wallet-1", display_id: 60001, tenant_id: "tenant-1", owner_type: "tenant", owner_id: "tenant-1", currency: "USD", status: "active", available_balance_minor: 11820000, locked_balance_minor: 0, created_at: now(), updated_at: now() },
        ]);
      case "/backend/admin/orders":
        return json(route, [
          { id: "order-uuid-1", display_id: 42001, tenant_id: "tenant-uuid-1", buyer_user_id: "buyer-1", tenant_plan_id: "tenant-plan-1", quantity: 1, currency: "USD", total_minor: 1400, order_status: "paid", billing_status: "paid", plan_snapshot: { plan_code: "vps-cx23-40gb-monthly" }, created_at: "2026-04-24T08:00:00Z" },
        ]);
      case "/backend/admin/services":
        return json(route, filterRows([
          { id: "service-uuid-1", display_id: 43001, tenant_id: "tenant-uuid-1", order_id: "order-uuid-1", order_display_id: 42001, buyer_display_id: 10002, tenant_plan_id: "tenant-plan-1", provider_source_id: "source-ready", provider_source_display_id: 10001, external_resource_id: "local-vps-405910", status: "active", billing_status: "paid", term_end: "2026-05-24T08:00:00Z", product_snapshot: { product_type: "vps", name: "VPS" }, plan_snapshot: { plan_code: "vps-cx23-40gb-monthly", name: "CX23 VPS 40GB", region: "eu-central" } },
        ], query, [
          ["display_id", (row) => row.display_id],
          ["order_display_id", (row) => row.order_display_id],
          ["provider_source_display_id", (row) => row.provider_source_display_id],
          ["status", (row) => row.status],
        ]));
      case "/backend/admin/services/service-uuid-1":
      case "/backend/reseller/services/service-uuid-1":
      case "/backend/client/services/service-uuid-1":
        return json(route, serviceDetail());
      case "/backend/admin/topup-requests":
        return json(route, filterRows([
          { id: "topup-uuid-1", display_id: 51001, tenant_id: "tenant-uuid-1", wallet_id: "wallet-1", wallet_display_id: 60001, requested_by: "buyer-1", requested_by_display_id: 10002, amount_minor: 50000, currency: "USD", payment_method: "bank_transfer", payment_reference: "LOCAL-REF-51001", status: "under_review", review_note: "", created_at: "2026-04-24T08:05:00Z" },
        ], query, [
          ["display_id", (row) => row.display_id],
          ["wallet_display_id", (row) => row.wallet_display_id],
          ["requested_by_display_id", (row) => row.requested_by_display_id],
          ["status", (row) => row.status],
        ]));
      case "/backend/admin/invoices":
        return json(route, filterRows([
          { id: "invoice-uuid-1", display_id: 44001, tenant_id: "tenant-uuid-1", buyer_user_id: "buyer-1", buyer_display_id: 10002, order_id: "order-uuid-1", order_display_id: 42001, status: "paid", currency: "USD", subtotal_minor: 1400, tax_minor: 0, discount_minor: 0, total_minor: 1400, issued_at: "2026-04-24T08:10:00Z", due_at: "2026-05-24T08:10:00Z", paid_at: "2026-04-24T08:12:00Z", created_at: "2026-04-24T08:10:00Z", updated_at: "2026-04-24T08:12:00Z" },
        ], query, [
          ["display_id", (row) => row.display_id],
          ["buyer_display_id", (row) => row.buyer_display_id],
          ["order_display_id", (row) => row.order_display_id],
          ["status", (row) => row.status],
        ]));
      case "/backend/admin/transactions":
        return json(route, filterRows([
          { id: "txn-uuid-1", display_id: 51001, tenant_id: "tenant-uuid-1", account_user_id: "buyer-1", account_display_id: 10002, order_id: "order-uuid-1", order_display_id: 42001, invoice_id: "invoice-uuid-1", invoice_display_id: 44001, type: "charge", status: "posted", currency: "USD", amount_minor: 1400, created_at: "2026-04-24T08:12:00Z" },
        ], query, [
          ["display_id", (row) => row.display_id],
          ["account_display_id", (row) => row.account_display_id],
          ["order_display_id", (row) => row.order_display_id],
          ["invoice_display_id", (row) => row.invoice_display_id],
          ["status", (row) => row.status],
        ]));
      case "/backend/admin/payment-reconciliation":
        return json(route, [
          { transaction: { id: "txn-uuid-1", display_id: 51001, tenant_id: "tenant-uuid-1", account_user_id: "buyer-1", account_display_id: 10002, order_id: "order-uuid-1", order_display_id: 42001, invoice_id: "invoice-uuid-1", invoice_display_id: 44001, type: "charge", status: "posted", currency: "USD", amount_minor: 1400, created_at: "2026-04-24T08:12:00Z" }, provider: "wallet", invoice: { id: "invoice-uuid-1", display_id: 44001, status: "paid", total_minor: 1400 }, ledger: { id: "ledger-uuid-1", display_id: 50002, wallet_display_id: 60001, direction: "debit", entry_type: "purchase", status: "posted" } },
        ]);
      case "/backend/admin/jobs":
        return json(route, filterRows([
          { id: "job-uuid-1", display_id: 3301, tenant_id: "tenant-uuid-1", job_type: "provider.provision", reference_type: "order", reference_id: "order-uuid-1", reference_display_id: 42001, source_id: "source-ready", source_display_id: 10001, status: "manual_review", priority: 5, attempt_count: 2, max_attempts: 5, next_attempt_at: "2026-04-24T09:00:00Z", last_error_code: "PROVIDER_TIMEOUT", last_error_message_redacted: "Provider timed out", manual_review_reason: "Needs provider check", correlation_id: "req-smoke", created_at: "2026-04-24T08:00:00Z", updated_at: "2026-04-24T08:35:00Z" },
        ], query, [
          ["display_id", (row) => row.display_id],
          ["source_display_id", (row) => row.source_display_id],
          ["status", (row) => row.status],
          ["job_type", (row) => row.job_type],
        ]));
      case "/backend/admin/jobs/summary":
        return json(route, {
          job_type: "provider.provision",
          total: 1,
          attention_count: 1,
          counts: { queued: 0, claimed: 0, running: 0, succeeded: 0, failed_retryable: 0, failed_terminal: 0, manual_review: 1, cancelled: 0 },
          latest_failure: { id: "job-uuid-1", display_id: 3301, status: "manual_review", last_error_code: "PROVIDER_TIMEOUT", last_error_message_redacted: "Provider timed out", manual_review_reason: "Needs provider check", created_at: "2026-04-24T08:00:00Z", updated_at: "2026-04-24T08:35:00Z" },
          generated_at: "2026-04-24T09:00:00Z",
        });
      case "/backend/admin/jobs/job-uuid-1/attempts":
        return json(route, [
          { id: "attempt-1", display_id: 91001, job_id: "job-uuid-1", worker_id: "worker-a", attempt_number: 2, started_at: "2026-04-24T08:30:00Z", finished_at: "2026-04-24T08:30:04Z", result: "manual_review", error_code: "PROVIDER_TIMEOUT", error_message_redacted: "Provider timed out", duration_ms: 4200, correlation_id: "req-smoke" },
        ]);
      case "/backend/admin/catalog/provider-sources":
        return json(route, filterRows([
          { id: "source-ready", display_id: 10001, source_type: "hetzner", name: "Local Fake Hetzner Ready", location: "local-fsn1", status: "active", inventory_mode: "provider_live", risk_level: "medium", created_at: now(), updated_at: now() },
        ], query, [
          ["display_id", (row) => row.display_id],
          ["source_type", (row) => row.source_type],
          ["status", (row) => row.status],
        ]));
      case "/backend/admin/catalog/provider-readiness":
        return json(route, filterRows([
          { plan_display_id: 10000, plan_code: "vps-cx23-40gb-monthly", plan_name: "CX23 VPS 40GB", product_type: "vps", plan_status: "active", plan_source_display_id: 10000, plan_source_status: "active", source_display_id: 10001, source_name: "Local Fake Hetzner Ready", source_type: "hetzner", source_status: "active", inventory_mode: "provider_live", state: "ready", reason: "Source is active and supports automatic provisioning." },
        ], query, [
          ["plan_display_id", (row) => row.plan_display_id],
          ["source_display_id", (row) => row.source_display_id],
          ["product_type", (row) => row.product_type],
          ["status", (row) => row.plan_status],
        ]));
      case "/backend/admin/audit-logs":
        return json(route, filterRows([
          { id: "audit-1", display_id: 70001, actor_id: "admin-1", actor_display_id: 10001, actor_type: "user", action: "job.retry", target_type: "job", target_id: "job-uuid-1", target_display_id: 3301, correlation_id: "req-smoke", created_at: "2026-04-24T08:40:00Z" },
        ], query, [
          ["display_id", (row) => row.display_id],
          ["actor_display_id", (row) => row.actor_display_id],
          ["action", (row) => row.action],
          ["target_type", (row) => row.target_type],
          ["target_display_id", (row) => row.target_display_id],
        ]));
      default:
        return json(route, []);
    }
  });
}

function now() {
  return "2026-04-24T08:00:00Z";
}

function serviceDetail() {
  return {
    id: "service-uuid-1",
    display_id: 43001,
    tenant_id: "tenant-uuid-1",
    order_id: "order-uuid-1",
    order_display_id: 42001,
    buyer_display_id: 10002,
    tenant_plan_id: "tenant-plan-1",
    provider_source_id: "source-ready",
    provider_source_display_id: 10001,
    external_resource_id: "local-vps-405910",
    status: "active",
    billing_status: "paid",
    term_end: "2026-05-24T08:00:00Z",
    credentials: [ // sensitive-text-allowlist
      { id: "credential-1", credential_type: "vps_root", masked_hint: "root / ****", status: "active" }, // sensitive-text-allowlist
    ],
    product_snapshot: { product_type: "vps", name: "VPS" },
    plan_snapshot: { plan_code: "vps-cx23-40gb-monthly", name: "CX23 VPS 40GB", region: "eu-central" },
  };
}

function filterRows(rows, query, filters) {
  return rows.filter((row) => filters.every(([key, value]) => {
    const actual = query.get(key);
    if (!actual) return true;
    return String(value(row)).toLowerCase() === actual.toLowerCase();
  }));
}

function json(route, data) {
  return route.fulfill({
    status: 200,
    contentType: "application/json",
    body: JSON.stringify({ data, request_id: "req-smoke" }),
  });
}

module.exports = { installApiMocks };
