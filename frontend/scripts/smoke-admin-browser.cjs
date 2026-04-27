const http = require("http");
const fs = require("fs");
const path = require("path");
const { spawn } = require("child_process");
const { chromium } = require("playwright");
const { installApiMocks } = require("./smoke-admin-api-mocks.cjs");
const { createFallbackSmokeFlows } = require("./smoke-admin-fallback-flows.cjs");

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
  "proxy-manager",
  "proxy-cheap",
  "vps-prod-01",
  "vps-scrape-01",
  "vps-scrape-02",
  "vps-test",
  "vps-api-gateway",
  "vps-db-replica",
  "vps-worker-03",
  "win-rdp-gamma",
  "win-dev-01",
  "vps-linux-small",
  "proxy-residential",
  "proxy-dc-shared",
];
const fallbackSmoke = createFallbackSmokeFlows({
  baseURL,
  installApiMocks,
  openAdminScreen,
  expectVisibleText,
  assertNoVisibleText,
  assertNoForbiddenText,
});

async function main() {
  validateServerMode(serverMode);
  const server = startNextServer();
  try {
    await waitForServer(baseURL, 45_000);
    const browser = await launchBrowser();
    try {
      const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
      await installApiMocks(page, { failPaths: new Set(["/backend/client/services", "/backend/reseller/services"]) });
      await page.goto(baseURL, { waitUntil: "domcontentloaded" });
      await page.waitForTimeout(800);

      await expectVisibleText(page, "Overview");
      await expectVisibleText(page, "INV-44001");
      await expectVisibleText(page, "Under review");
      await assertNoVisibleText(page, ["under_review", "vps-scrape-02"], "overview top-up and activity labels");
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
      await fallbackSmoke.provisioning(browser);

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
      await fallbackSmoke.providerSources(browser);
      await fallbackSmoke.providerReadiness(browser);

      await openAdminScreen(page, /VPS/i);
      await expectVisibleText(page, "Live VPS inventory");
      await expectVisibleText(page, "SVC-43001");
      await expectVisibleText(page, "ORD-42001");
      await expectVisibleText(page, "ACC-10002");
      await expectVisibleText(page, "SRC-10001");
      await assertNoVisibleText(page, ["service-uuid-1", "order-uuid-1", "buyer-1", "source-ready", "tenant-uuid-1"], "service public ID labels");
      await assertNoForbiddenText(page, "services");
      await fallbackSmoke.adminService(browser);

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
      await fallbackSmoke.topup(browser);

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
      await fallbackSmoke.audit(browser);

      await openAdminScreen(page, /Alerts/i);
      await expectVisibleText(page, "Open alerts");
      await expectVisibleText(page, "3 jobs stuck in manual review > 1h");
      await expectVisibleText(page, "Provider timeout on OVH.");
      await assertNoVisibleText(page, ["manual_review"], "admin alert labels");
      await assertNoForbiddenText(page, "admin alerts");

      await page.getByRole("button", { name: "Reseller · ProxyVN" }).click();
      await page.waitForTimeout(500);
      await openAdminScreen(page, /^VPS$/i);
      await expectVisibleText(page, "Windows RDP Workspace");
      await expectVisibleText(page, "Proxy Automation VPS");
      await assertNoForbiddenText(page, "reseller service demo labels");

      await page.getByRole("button", { name: "Client · Linh Tran" }).click();
      await page.waitForTimeout(500);
      await expectVisibleText(page, "My services");
      await expectVisibleText(page, "Proxy Automation VPS");
      await expectVisibleText(page, "Small Linux VPS");
      await assertNoForbiddenText(page, "client service demo labels");
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
