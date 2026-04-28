function createFallbackSmokeFlows(context) {
  const {
    baseURL,
    installApiMocks,
    openAdminScreen,
    expectVisibleText,
    assertNoVisibleText,
    assertNoForbiddenText,
  } = context;

  async function withFallbackPage(browser, failPaths, screenName, runAssertions) {
    const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
    try {
      await installApiMocks(page, { failPaths: new Set(failPaths) });
      await page.goto(baseURL, { waitUntil: "domcontentloaded" });
      await page.waitForTimeout(800);
      await openAdminScreen(page, screenName);
      await runAssertions(page);
    } finally {
      await page.close();
    }
  }

  function provisioning(browser) {
    return withFallbackPage(browser, ["/backend/admin/jobs"], /Provisioning queue/i, async (page) => {
      await expectVisibleText(page, "Live job API unavailable. Showing demo queue data");
      await expectVisibleText(page, "Provider Timeout: Resource State Unknown");
      await expectVisibleText(page, "Auth Failed");
      await expectVisibleText(page, "Partial Success: External ID Unknown");
      await expectVisibleText(page, "Budget Proxy Upstream");
      await assertNoVisibleText(page, ["provider_timeout", "auth_failed", "partial_success", "external_id", "proxy-cheap", "cor_"], "provisioning demo fallback labels");
      await assertNoForbiddenText(page, "provisioning demo fallback");
    });
  }

  function topup(browser) {
    return withFallbackPage(browser, ["/backend/admin/topup-requests"], /Top-up verification/i, async (page) => {
      await expectVisibleText(page, "Live API unavailable. Showing demo top-up data");
      await expectVisibleText(page, "Reseller Wallet");
      await expectVisibleText(page, "Client Wallet");
      await expectVisibleText(page, "Ref provided");
      await assertNoVisibleText(page, ["reseller_wallet", "client_wallet", "pending_verification"], "top-up demo fallback labels");
      await assertNoForbiddenText(page, "top-up demo fallback");
    });
  }

  function overview(browser) {
    return withFallbackPage(browser, ["/backend/admin/topup-requests"], /Overview/i, async (page) => {
      await expectVisibleText(page, "New high-priority ticket opened by Acme Proxy Co.");
      await assertNoVisibleText(page, ["under_review", "vps-scrape-02", "T-8124"], "overview demo fallback labels");
      await assertNoForbiddenText(page, "overview demo fallback");
    });
  }

  function audit(browser) {
    return withFallbackPage(browser, ["/backend/admin/audit-logs"], /Audit logs/i, async (page) => {
      await expectVisibleText(page, "Live API unavailable. Showing demo audit data");
      await expectVisibleText(page, "Provisioning Worker");
      await expectVisibleText(page, "manual review threshold exceeded");
      await expectVisibleText(page, "RBAC migration");
      await expectVisibleText(page, "VPS Small plan");
      await expectVisibleText(page, "Session not shown");
      await assertNoVisibleText(page, ["prov-worker", "billing-worker", "health-worker", "manual_review", "0003_rbac", "vps-scrape-02", "VPS-SMALL", "RES-PROX-4G", "session-991"], "audit demo fallback labels");
      await assertNoForbiddenText(page, "audit demo fallback");
    });
  }

  function adminService(browser) {
    return withFallbackPage(browser, ["/backend/admin/services"], /VPS/i, async (page) => {
      await expectVisibleText(page, "Live API unavailable. Showing demo VPS data.");
      await expectVisibleText(page, "Production Linux VPS");
      await expectVisibleText(page, "API Gateway VPS");
      await expectVisibleText(page, "Database Replica VPS");
      await assertNoForbiddenText(page, "admin service demo fallback");
    });
  }

  function providerReadiness(browser) {
    return withFallbackPage(browser, ["/backend/admin/catalog/provider-readiness"], /Providers \/ Sources/i, async (page) => {
      await expectVisibleText(page, "Live readiness API unavailable. Demo rows are shown.");
      await expectVisibleText(page, "VPS Linux Small");
      await expectVisibleText(page, "Residential Proxy");
      await expectVisibleText(page, "Datacenter Shared");
      await assertNoVisibleText(page, ["vps-linux-small", "proxy-residential", "proxy-dc-shared"], "provider readiness demo labels");
      await assertNoForbiddenText(page, "provider readiness demo fallback");
    });
  }

  function providerSources(browser) {
    return withFallbackPage(browser, ["/backend/admin/catalog/provider-sources"], /Providers \/ Sources/i, async (page) => {
      await expectVisibleText(page, "Live API unavailable. Showing demo provider data.");
      await expectVisibleText(page, "Self-hosted Proxy Manager");
      await expectVisibleText(page, "Budget Proxy Upstream");
      await assertNoVisibleText(page, ["proxy-manager", "proxy-cheap"], "provider source demo fallback labels");
      await assertNoForbiddenText(page, "provider source demo fallback");
    });
  }

  return { adminService, audit, overview, providerReadiness, providerSources, provisioning, topup };
}

module.exports = { createFallbackSmokeFlows };
