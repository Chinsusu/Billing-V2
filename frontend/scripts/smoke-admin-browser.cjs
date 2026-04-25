const http = require("http");
const path = require("path");
const { spawn } = require("child_process");
const { chromium } = require("playwright");

const appRoot = path.resolve(__dirname, "..");
const host = process.env.SMOKE_HOST || "127.0.0.1";
const port = Number(process.env.SMOKE_PORT || "3120");
const baseURL = `http://${host}:${port}`;
const forbiddenText = [
  "payload_json",
  "capability_profile",
  "provider_account_id",
  "secret",
  "raw_response",
];

async function main() {
  const server = startNextServer();
  try {
    await waitForServer(baseURL, 45_000);
    const browser = await launchBrowser();
    try {
      const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
      await installApiMocks(page);
      await page.goto(baseURL, { waitUntil: "domcontentloaded" });
      await page.waitForTimeout(800);

      await expectVisibleText(page, "Overview");
      await expectVisibleText(page, "INV-44001");
      await assertNoForbiddenText(page, "overview");

      await openAdminScreen(page, /Provisioning queue/i);
      await expectVisibleText(page, "Live provisioning jobs");
      await expectVisibleText(page, "JOB-3301");
      await expectVisibleText(page, "Manual Review");
      await page.getByRole("button", { name: "JOB-3301" }).click();
      await expectVisibleText(page, "SOURCE READINESS");
      await expectVisibleText(page, "PLAN-10000 / vps-cx23-40gb-monthly / SRC-10001");
      await assertNoForbiddenText(page, "provisioning");

      await openAdminScreen(page, /Providers \/ Sources/i);
      await expectVisibleText(page, "Provider readiness");
      await expectVisibleText(page, "PLAN-10000");
      await expectVisibleText(page, "SRC-10001");
      await expectVisibleText(page, "Local Fake Hetzner Ready");
      await assertNoForbiddenText(page, "providers");

      await openAdminScreen(page, /Top-up verification/i);
      await expectVisibleText(page, "Live top-up queue");
      await expectVisibleText(page, "TUP-51001");
      await expectVisibleText(page, "under_review");
      await assertNoForbiddenText(page, "topups");

      await openAdminScreen(page, /Audit logs/i);
      await expectVisibleText(page, "Live audit filters applied.");
      await expectVisibleText(page, "job.retry");
      await expectVisibleText(page, "req-smoke");
      await assertNoForbiddenText(page, "audit logs");
    } finally {
      await browser.close();
    }
  } finally {
    await stopServer(server);
  }
  console.log("Admin browser smoke passed.");
}

function startNextServer() {
  const nextBin = path.join(appRoot, "node_modules", "next", "dist", "bin", "next");
  const child = spawn(process.execPath, [nextBin, "dev", "--hostname", host, "--port", String(port)], {
    cwd: appRoot,
    env: { ...process.env, NEXT_TELEMETRY_DISABLED: "1" },
    stdio: ["ignore", "pipe", "pipe"],
  });
  child.stdout.on("data", (chunk) => process.stdout.write(`[next] ${chunk}`));
  child.stderr.on("data", (chunk) => process.stderr.write(`[next] ${chunk}`));
  child.on("exit", (code, signal) => {
    if (code !== 0 && signal !== "SIGTERM") {
      console.error(`Next dev server exited with code ${code ?? signal}.`);
    }
  });
  return child;
}

function waitForServer(url, timeoutMs) {
  const start = Date.now();
  return new Promise((resolve, reject) => {
    const check = () => {
      http.get(url, (res) => {
        res.resume();
        resolve();
      }).on("error", (error) => {
        if (Date.now() - start > timeoutMs) {
          reject(new Error(`Timed out waiting for ${url}: ${error.message}`));
          return;
        }
        setTimeout(check, 500);
      });
    };
    check();
  });
}

async function launchBrowser() {
  try {
    return await chromium.launch({ headless: true });
  } catch (error) {
    if (process.platform === "win32") {
      return chromium.launch({ channel: "chrome", headless: true });
    }
    throw error;
  }
}

async function installApiMocks(page) {
  await page.route("**/backend/**", (route) => {
    const request = route.request();
    const pathname = new URL(request.url()).pathname;
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
        return json(route, [
          { id: "service-uuid-1", display_id: 43001, tenant_id: "tenant-uuid-1", order_id: "order-uuid-1", tenant_plan_id: "tenant-plan-1", provider_source_id: "source-ready", external_resource_id: "srv-local-1", status: "active", billing_status: "paid", term_end: "2026-05-24T08:00:00Z" },
        ]);
      case "/backend/admin/topup-requests":
        return json(route, [
          { id: "topup-uuid-1", display_id: 51001, tenant_id: "tenant-uuid-1", wallet_id: "wallet-1", requested_by: "buyer-1", amount_minor: 50000, currency: "USD", payment_method: "bank_transfer", payment_reference: "LOCAL-REF-51001", status: "under_review", review_note: "", created_at: "2026-04-24T08:05:00Z" },
        ]);
      case "/backend/admin/invoices":
        return json(route, [
          { id: "invoice-uuid-1", display_id: 44001, tenant_id: "tenant-uuid-1", buyer_user_id: "buyer-1", order_id: "order-uuid-1", status: "paid", currency: "USD", subtotal_minor: 1400, tax_minor: 0, discount_minor: 0, total_minor: 1400, issued_at: "2026-04-24T08:10:00Z", due_at: "2026-05-24T08:10:00Z", paid_at: "2026-04-24T08:12:00Z", created_at: "2026-04-24T08:10:00Z", updated_at: "2026-04-24T08:12:00Z" },
        ]);
      case "/backend/admin/jobs":
        return json(route, [
          { id: "job-uuid-1", display_id: 3301, tenant_id: "tenant-uuid-1", job_type: "provider.provision", reference_type: "order", reference_id: "order-uuid-1", source_id: "source-ready", status: "manual_review", priority: 5, attempt_count: 2, max_attempts: 5, next_attempt_at: "2026-04-24T09:00:00Z", last_error_code: "PROVIDER_TIMEOUT", last_error_message_redacted: "Provider timed out", manual_review_reason: "Needs provider check", correlation_id: "req-smoke", created_at: "2026-04-24T08:00:00Z", updated_at: "2026-04-24T08:35:00Z" },
        ]);
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
        return json(route, [
          { id: "source-ready", display_id: 10001, source_type: "hetzner", name: "Local Fake Hetzner Ready", location: "local-fsn1", status: "active", inventory_mode: "provider_live", risk_level: "medium", created_at: now(), updated_at: now() },
        ]);
      case "/backend/admin/catalog/provider-readiness":
        return json(route, [
          { plan_display_id: 10000, plan_code: "vps-cx23-40gb-monthly", plan_name: "CX23 VPS 40GB", product_type: "vps", plan_status: "active", plan_source_display_id: 10000, plan_source_status: "active", source_display_id: 10001, source_name: "Local Fake Hetzner Ready", source_type: "hetzner", source_status: "active", inventory_mode: "provider_live", state: "ready", reason: "Source is active and supports automatic provisioning." },
        ]);
      case "/backend/admin/audit-logs":
        return json(route, [
          { id: "audit-1", display_id: 70001, actor_id: "admin-1", actor_type: "user", action: "job.retry", target_type: "job", target_id: "job-uuid-1", correlation_id: "req-smoke", created_at: "2026-04-24T08:40:00Z" },
        ]);
      default:
        return json(route, []);
    }
  });
}

function now() {
  return "2026-04-24T08:00:00Z";
}

function json(route, data) {
  return route.fulfill({
    status: 200,
    contentType: "application/json",
    body: JSON.stringify({ data, request_id: "req-smoke" }),
  });
}

async function openAdminScreen(page, name) {
  await page.getByRole("button", { name }).click();
  await page.waitForTimeout(500);
}

async function expectVisibleText(page, text) {
  await page.getByText(text, { exact: false }).first().waitFor({ timeout: 10_000 });
}

async function assertNoForbiddenText(page, screenName) {
  const bodyText = await page.locator("body").innerText();
  for (const value of forbiddenText) {
    if (bodyText.includes(value)) {
      throw new Error(`Forbidden text '${value}' is visible on ${screenName}.`);
    }
  }
}

function stopServer(server) {
  return new Promise((resolve) => {
    if (server.killed || server.exitCode !== null) {
      resolve();
      return;
    }
    server.once("exit", resolve);
    server.kill("SIGTERM");
    setTimeout(resolve, 2_000);
  });
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
