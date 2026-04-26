const http = require("http");
const fs = require("fs");
const path = require("path");
const { spawn } = require("child_process");
const { chromium } = require("playwright");

const args = parseArgs(process.argv.slice(2));
const appRoot = path.resolve(__dirname, "..");
const host = process.env.SMOKE_HOST || "127.0.0.1";
const port = Number(process.env.SMOKE_PORT || "3120");
const serverMode = args.server || process.env.SMOKE_SERVER || "dev";
const baseURL = `http://${host}:${port}`;
const forbiddenText = [
  "payload_json", // sensitive-text-allowlist: browser smoke forbidden text
  "capability_profile", // sensitive-text-allowlist: browser smoke forbidden text
  "provider_account_id", // sensitive-text-allowlist: browser smoke forbidden text
  "secret", // sensitive-text-allowlist: browser smoke forbidden text
  "raw_response", // sensitive-text-allowlist: browser smoke forbidden text
];

async function main() {
  validateServerMode(serverMode);
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
      await expectVisibleText(page, "Under review");
      await assertNoVisibleText(page, ["under_review"], "overview top-up status label");
      await assertNoForbiddenText(page, "overview");

      await openAdminScreen(page, /Provisioning queue/i);
      await expectVisibleText(page, "Live provisioning jobs");
      await expectVisibleText(page, "JOB-3301");
      await expectVisibleText(page, "ORD-42001");
      await expectVisibleText(page, "SRC-10001");
      await page.getByRole("cell", { name: "Manual Review", exact: true }).waitFor({ timeout: 10_000 });
      await page.getByLabel("Status", { exact: true }).selectOption("manual_review");
      const filteredJobs = page.waitForResponse((response) => {
        const url = new URL(response.url());
        return url.pathname === "/backend/admin/jobs"
          && url.searchParams.get("status") === "manual_review";
      });
      await page.getByRole("button", { name: "Apply" }).click();
      await filteredJobs;
      await expectVisibleText(page, "Live provisioning filters applied.");
      await page.getByRole("button", { name: "JOB-3301" }).click();
      await expectVisibleText(page, "SOURCE READINESS");
      await expectVisibleText(page, "PLAN-10000 / CX23 VPS 40GB / SRC-10001");
      await expectVisibleText(page, "Provider timeout");
      await expectVisibleText(page, "Worker A");
      await assertNoVisibleText(page, ["job-uuid-1", "order-uuid-1", "source-ready", "tenant-uuid-1", "vps-cx23-40gb-monthly", "PROVIDER_TIMEOUT", "worker-a"], "provisioning public ID labels");
      await assertNoForbiddenText(page, "provisioning");
      await smokeProvisioningFallback(browser);

      await openAdminScreen(page, /Providers \/ Sources/i);
      await expectVisibleText(page, "Provider readiness");
      await expectVisibleText(page, "PLAN-10000");
      await expectVisibleText(page, "SRC-10001");
      await expectVisibleText(page, "Local Fake Hetzner Ready");
      await page.getByRole("cell", { name: "VPS", exact: true }).waitFor({ timeout: 10_000 });
      await page.getByLabel("Product", { exact: true }).selectOption("vps");
      const filteredReadiness = page.waitForResponse((response) => {
        const url = new URL(response.url());
        return url.pathname === "/backend/admin/catalog/provider-readiness"
          && url.searchParams.get("product_type") === "vps";
      });
      await page.locator("form").filter({ has: page.getByLabel("Product", { exact: true }) }).getByRole("button", { name: "Apply" }).click();
      await filteredReadiness;
      await page.getByLabel("Source type", { exact: true }).selectOption("hetzner");
      const filteredProviderSource = page.waitForResponse((response) => {
        const url = new URL(response.url());
        return url.pathname === "/backend/admin/catalog/provider-sources"
          && url.searchParams.get("source_type") === "hetzner";
      });
      await page.locator("form").filter({ has: page.getByLabel("Source type", { exact: true }) }).getByRole("button", { name: "Apply" }).click();
      await filteredProviderSource;
      await page.getByRole("cell", { name: "Hetzner", exact: true }).first().waitFor({ timeout: 10_000 });
      await page.getByRole("cell", { name: "Provider live", exact: true }).waitFor({ timeout: 10_000 });
      await page.getByRole("cell", { name: "Medium", exact: true }).waitFor({ timeout: 10_000 });
      await expectVisibleText(page, "Live provider source filters applied.");
      await assertNoForbiddenText(page, "providers");

      await openAdminScreen(page, /VPS/i);
      await expectVisibleText(page, "Live VPS inventory");
      await expectVisibleText(page, "SVC-43001");
      await expectVisibleText(page, "ORD-42001");
      await expectVisibleText(page, "ACC-10002");
      await expectVisibleText(page, "SRC-10001");
      await assertNoVisibleText(page, ["service-uuid-1", "order-uuid-1", "buyer-1", "source-ready", "tenant-uuid-1"], "service public ID labels");
      await assertNoForbiddenText(page, "services");

      await openAdminScreen(page, /Top-up verification/i);
      await expectVisibleText(page, "Live top-up queue");
      await expectVisibleText(page, "TUP-51001");
      await expectVisibleText(page, "WAL-60001");
      await expectVisibleText(page, "ACC-10002");
      await page.getByRole("cell", { name: "Bank transfer", exact: true }).waitFor({ timeout: 10_000 });
      await page.getByRole("cell", { name: "Under review", exact: true }).waitFor({ timeout: 10_000 });
      await page.getByLabel("Status", { exact: true }).selectOption("under_review");
      const filteredTopup = page.waitForResponse((response) => {
        const url = new URL(response.url());
        return url.pathname === "/backend/admin/topup-requests"
          && url.searchParams.get("status") === "under_review";
      });
      await page.getByRole("button", { name: "Apply" }).click();
      await filteredTopup;
      await expectVisibleText(page, "Live top-up filters applied.");
      await assertNoVisibleText(page, ["topup-uuid-1", "wallet-1", "buyer-1", "tenant-uuid-1"], "top-up public ID labels");
      await assertNoForbiddenText(page, "topups");

      await openAdminScreen(page, /^Invoices$/i);
      await expectVisibleText(page, "Live invoice data");
      await expectVisibleText(page, "INV-44001");
      await expectVisibleText(page, "ACC-10002");
      await expectVisibleText(page, "ORD-42001");
      await page.getByLabel("Display ID").fill("44001");
      await page.getByLabel("Customer public ID").fill("10002");
      await page.getByLabel("Status", { exact: true }).selectOption("paid");
      const filteredInvoice = page.waitForResponse((response) => {
        const url = new URL(response.url());
        return url.pathname === "/backend/admin/invoices"
          && url.searchParams.get("display_id") === "44001"
          && url.searchParams.get("buyer_display_id") === "10002"
          && url.searchParams.get("status") === "paid";
      });
      await page.getByRole("button", { name: "Apply" }).click();
      await filteredInvoice;
      await expectVisibleText(page, "Live invoice filters applied.");
      await expectVisibleText(page, "INV-44001");
      await assertNoVisibleText(page, ["invoice-uuid-1", "buyer-1", "order-uuid-1", "tenant-uuid-1"], "invoice public ID filter");
      await assertNoForbiddenText(page, "invoice public ID filter");

      await openAdminScreen(page, /^Transactions$/i);
      await expectVisibleText(page, "Live transaction data");
      await expectVisibleText(page, "TX-51001");
      await expectVisibleText(page, "ACC-10002");
      await expectVisibleText(page, "ORD-42001");
      await expectVisibleText(page, "INV-44001");
      await page.getByRole("cell", { name: "Charge", exact: true }).waitFor({ timeout: 10_000 });
      await page.getByLabel("Status", { exact: true }).selectOption("posted");
      const filteredTransaction = page.waitForResponse((response) => {
        const url = new URL(response.url());
        return url.pathname === "/backend/admin/transactions"
          && url.searchParams.get("status") === "posted";
      });
      await page.getByRole("button", { name: "Apply" }).click();
      await filteredTransaction;
      await expectVisibleText(page, "Live transaction filters applied.");
      await assertNoVisibleText(page, ["txn-uuid-1", "buyer-1", "order-uuid-1", "invoice-uuid-1", "tenant-uuid-1"], "transaction public ID labels");
      await assertNoForbiddenText(page, "transactions");

      await openAdminScreen(page, /^Reports$/i);
      await expectVisibleText(page, "Payment reconciliation");
      await page.getByRole("cell", { name: "Wallet", exact: true }).waitFor({ timeout: 10_000 });
      await assertNoVisibleText(page, ["wallet"], "reports payment provider labels");
      await assertNoForbiddenText(page, "reports");

      await openAdminScreen(page, /Audit logs/i);
      await expectVisibleText(page, "Live audit filters applied.");
      await expectVisibleText(page, "AUD-70001");
      await expectVisibleText(page, "ACC-10001");
      await page.getByRole("cell", { name: "Retry job", exact: true }).waitFor({ timeout: 10_000 });
      await expectVisibleText(page, "Job JOB-3301");
      await page.getByLabel("Actor public ID").fill("10001");
      await page.getByLabel("Action", { exact: true }).selectOption("job.retry");
      await page.getByLabel("Target", { exact: true }).selectOption("job");
      await page.getByLabel("Target public ID").fill("3301");
      const filteredAudit = page.waitForResponse((response) => {
        const url = new URL(response.url());
        return url.pathname === "/backend/admin/audit-logs"
          && url.searchParams.get("actor_display_id") === "10001"
          && url.searchParams.get("action") === "job.retry"
          && url.searchParams.get("target_type") === "job"
          && url.searchParams.get("target_display_id") === "3301";
      });
      await page.getByRole("button", { name: "Apply" }).click();
      await filteredAudit;
      await expectVisibleText(page, "Request not shown");
      await assertNoVisibleText(page, ["audit-1", "admin-1", "job-uuid-1", "tenant-uuid-1", "req-smoke", "job.retry"], "audit public ID labels");
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
  if (serverMode === "standalone") {
    return startStandaloneServer();
  }
  const nextBin = path.join(appRoot, "node_modules", "next", "dist", "bin", "next");
  const child = spawn(process.execPath, [nextBin, serverMode, "--hostname", host, "--port", String(port)], {
    cwd: appRoot,
    env: { ...process.env, NEXT_TELEMETRY_DISABLED: "1" },
    stdio: ["ignore", "pipe", "pipe"],
  });
  child.stdout.on("data", (chunk) => process.stdout.write(`[next] ${chunk}`));
  child.stderr.on("data", (chunk) => process.stderr.write(`[next] ${chunk}`));
  child.on("exit", (code, signal) => {
    if (code !== 0 && signal !== "SIGTERM") {
      console.error(`Next ${serverMode} server exited with code ${code ?? signal}.`);
    }
  });
  return child;
}

function startStandaloneServer() {
  ensureStandaloneAssets();
  const serverPath = path.join(appRoot, ".next", "standalone", "server.js");
  const child = spawn(process.execPath, [serverPath], {
    cwd: path.dirname(serverPath),
    env: {
      ...process.env,
      HOSTNAME: host,
      PORT: String(port),
      NEXT_TELEMETRY_DISABLED: "1",
    },
    stdio: ["ignore", "pipe", "pipe"],
  });
  child.stdout.on("data", (chunk) => process.stdout.write(`[standalone] ${chunk}`));
  child.stderr.on("data", (chunk) => process.stderr.write(`[standalone] ${chunk}`));
  child.on("exit", (code, signal) => {
    if (code !== 0 && signal !== "SIGTERM") {
      console.error(`Next standalone server exited with code ${code ?? signal}.`);
    }
  });
  return child;
}

function ensureStandaloneAssets() {
  const serverPath = path.join(appRoot, ".next", "standalone", "server.js");
  if (!fs.existsSync(serverPath)) {
    throw new Error("Standalone build is missing. Run `npm run build` before `npm run smoke:admin:ci`.");
  }
  copyIfExists(path.join(appRoot, ".next", "static"), path.join(appRoot, ".next", "standalone", ".next", "static"));
  copyIfExists(path.join(appRoot, "public"), path.join(appRoot, ".next", "standalone", "public"));
}

function copyIfExists(source, target) {
  if (!fs.existsSync(source)) {
    return;
  }
  fs.cpSync(source, target, { recursive: true });
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
      try {
        return await chromium.launch({ channel: "chrome", headless: true });
      } catch (fallbackError) {
        throw browserLaunchError(fallbackError, error);
      }
    }
    throw browserLaunchError(error);
  }
}

function browserLaunchError(error, primaryError) {
  const details = [
    "Unable to launch Playwright Chromium.",
    "Install the browser runtime with `npx playwright install chromium` locally or `npx playwright install --with-deps chromium` in CI.",
    `Original error: ${error.message}`,
  ];
  if (primaryError && primaryError !== error) {
    details.push(`Primary launch error: ${primaryError.message}`);
  }
  return new Error(details.join("\n"));
}

function validateServerMode(value) {
  if (value !== "dev" && value !== "start" && value !== "standalone") {
    throw new Error(`Unsupported smoke server mode '${value}'. Use 'dev', 'start', or 'standalone'.`);
  }
}

function parseArgs(argv) {
  const parsed = {};
  for (let index = 0; index < argv.length; index += 1) {
    const value = argv[index];
    if (value.startsWith("--server=")) {
      parsed.server = value.slice("--server=".length);
      continue;
    }
    if (value === "--server" && argv[index + 1]) {
      parsed.server = argv[index + 1];
      index += 1;
    }
  }
  return parsed;
}

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

async function openAdminScreen(page, name) {
  await page.getByRole("button", { name }).click();
  await page.waitForTimeout(500);
}

async function smokeProvisioningFallback(browser) {
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
  try {
    await installApiMocks(page, { failPaths: new Set(["/backend/admin/jobs"]) });
    await page.goto(baseURL, { waitUntil: "domcontentloaded" });
    await page.waitForTimeout(800);
    await openAdminScreen(page, /Provisioning queue/i);
    await expectVisibleText(page, "Live job API unavailable. Showing demo queue data");
    await expectVisibleText(page, "Provider Timeout: Resource State Unknown");
    await expectVisibleText(page, "Auth Failed");
    await expectVisibleText(page, "Partial Success: External ID Unknown");
    await assertNoVisibleText(page, ["provider_timeout", "auth_failed", "partial_success", "external_id", "proxy-cheap", "cor_"], "provisioning demo fallback labels");
    await assertNoForbiddenText(page, "provisioning demo fallback");
  } finally {
    await page.close();
  }
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

async function assertNoVisibleText(page, values, screenName) {
  const bodyText = await page.locator("body").innerText();
  for (const value of values) {
    if (bodyText.includes(value)) {
      throw new Error(`Unexpected backend reference '${value}' is visible on ${screenName}.`);
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
