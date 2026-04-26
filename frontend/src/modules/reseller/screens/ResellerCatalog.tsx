"use client";

import { useState } from "react";

import { billingApi } from "@/lib/api/billing";
import { billingCycleLabel } from "@/lib/api/displayLabels";
import { moneyMinor, fmtMoney } from "@/lib/api/format";
import type { CatalogPlan, TenantCatalogPlan, TenantCatalogProduct } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { RESELLER_CATALOG } from "@/mocks/billingData";

interface CatalogRow {
  key: string;
  plan: string;
  unit: string;
  cost: string;
  sellingMinor: number | null;
  currency: string;
  margin: number | null;
  stock: string;
  status: string;
  masterPlan?: CatalogPlan;
  tenantPlan?: TenantCatalogPlan;
}

function draftPrice(plan: CatalogPlan, drafts: Record<string, string>) {
  return drafts[plan.id] ?? (plan.suggested_price_minor / 100).toFixed(2);
}

function draftPriceMinor(plan: CatalogPlan, drafts: Record<string, string>) {
  const parsed = Number.parseFloat(draftPrice(plan, drafts));
  if (!Number.isFinite(parsed) || parsed <= 0) return plan.suggested_price_minor;
  return Math.round(parsed * 100);
}

function marginPercent(sellingMinor: number, costMinor: number) {
  if (costMinor <= 0) return null;
  return Math.round(((sellingMinor - costMinor) / costMinor) * 100);
}

function marginClass(margin: number | null) {
  if (margin == null) return "text-gray-400";
  if (margin < 0) return "text-red-600 font-medium";
  if (margin < 20) return "text-amber-600";
  return "text-green-700";
}

function statusText(catalogStatus: string, masterStatus: string, usingLive: boolean) {
  if (catalogStatus === "error" || masterStatus === "error") return "Live catalog partially unavailable.";
  if (catalogStatus === "loading" || masterStatus === "loading") return "Refreshing live catalog...";
  return usingLive ? "Live catalog" : "Demo catalog data";
}

function clonePlanBody(plan: CatalogPlan, tenantProductID: string, sellingPriceMinor: number) {
  return {
    tenant_product_id: tenantProductID,
    master_plan_id: plan.id,
    selling_price_minor: sellingPriceMinor,
    reseller_cost_minor: plan.base_cost_minor,
    currency: plan.currency,
    margin_policy: {
      min_margin_minor: Math.max(0, plan.reseller_min_price_minor - plan.base_cost_minor),
    },
    visibility: "public",
    status: "active",
    clone_version: plan.version,
    product_snapshot: {
      product_id: plan.product_id,
    },
    plan_snapshot: {
      plan_id: plan.id,
      plan_code: plan.plan_code,
      name: plan.name,
      specs: plan.specs,
      billing_cycle: plan.billing_cycle,
    },
    price_snapshot: {
      selling_price_minor: sellingPriceMinor,
      currency: plan.currency,
    },
    capability_snapshot: {},
  };
}

export function ResellerCatalog() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [priceDrafts, setPriceDrafts] = useState<Record<string, string>>({});
  const [busyPlanID, setBusyPlanID] = useState<string | null>(null);
  const [notice, setNotice] = useState<{ type: "success" | "error"; text: string } | null>(null);
  const catalog = useApiResource(
    () => billingApi.listResellerCatalog({ limit: 100 }),
    `reseller-catalog:${refreshKey}`,
  );
  const masterPlans = useApiResource(
    () => billingApi.listResellerMasterPlans({ limit: 100, status: "active" }),
    "reseller-master-plans",
  );
  const loading = catalog.status === "loading" || masterPlans.status === "loading";
  const tenantProducts = catalog.data?.products ?? [];
  const tenantPlans = catalog.data?.plans ?? [];
  const productByMaster = new Map(tenantProducts.map((product) => [product.master_product_id, product]));
  const planByMaster = new Map(tenantPlans.map((plan) => [plan.master_plan_id, plan]));
  const masterByID = new Map((masterPlans.data ?? []).map((plan) => [plan.id, plan]));
  const usingLive = catalog.status === "success" || masterPlans.status === "success";

  async function handleClone(plan: CatalogPlan) {
    setBusyPlanID(plan.id);
    setNotice(null);
    let createdProduct = false;
    try {
      let tenantProduct: TenantCatalogProduct | undefined = productByMaster.get(plan.product_id);
      if (!tenantProduct) {
        tenantProduct = await billingApi.cloneResellerCatalogProduct({
          master_product_id: plan.product_id,
          status: "active",
          clone_version: plan.version,
        });
        createdProduct = true;
      }
      await billingApi.cloneResellerCatalogPlan(clonePlanBody(plan, tenantProduct.id, draftPriceMinor(plan, priceDrafts)));
      setPriceDrafts((current) => {
        const next = { ...current };
        delete next[plan.id];
        return next;
      });
      setNotice({ type: "success", text: `Added ${plan.name} to catalog.` });
      setRefreshKey((current) => current + 1);
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : "Clone action failed.";
      setNotice({ type: "error", text: message });
      if (createdProduct) setRefreshKey((current) => current + 1);
    } finally {
      setBusyPlanID(null);
    }
  }

  const masterRows: CatalogRow[] = (masterPlans.data ?? []).map((plan) => {
    const tenantPlan = planByMaster.get(plan.id);
    const sellingMinor = tenantPlan?.selling_price_minor ?? draftPriceMinor(plan, priceDrafts);
    return {
      key: tenantPlan?.id ?? plan.id,
      plan: tenantPlan ? snapshotText(tenantPlan.plan_snapshot) ?? plan.name : plan.name,
      unit: tenantPlan?.visibility ?? billingCycleLabel(plan.billing_cycle),
      cost: moneyMinor(tenantPlan?.reseller_cost_minor ?? plan.base_cost_minor, plan.currency),
      sellingMinor,
      currency: tenantPlan?.currency ?? plan.currency,
      margin: marginPercent(sellingMinor, tenantPlan?.reseller_cost_minor ?? plan.base_cost_minor),
      stock: tenantPlan ? "catalog" : "master",
      status: tenantPlan?.status ?? plan.status,
      masterPlan: plan,
      tenantPlan,
    };
  });
  const orphanCatalogRows: CatalogRow[] = tenantPlans
    .filter((plan) => !masterByID.has(plan.master_plan_id))
    .map((plan) => ({
      key: plan.id,
      plan: snapshotText(plan.plan_snapshot) ?? plan.master_plan_id,
      unit: plan.visibility,
      cost: plan.reseller_cost_minor != null ? moneyMinor(plan.reseller_cost_minor, plan.currency) : "-",
      sellingMinor: plan.selling_price_minor,
      currency: plan.currency,
      margin: plan.reseller_cost_minor ? marginPercent(plan.selling_price_minor, plan.reseller_cost_minor) : null,
      stock: "catalog",
      status: plan.status,
      tenantPlan: plan,
    }));
  const liveRows = usingLive ? [...masterRows, ...orphanCatalogRows] : null;
  const demoRows: CatalogRow[] = RESELLER_CATALOG.map((item) => ({
    key: item.plan,
    plan: item.plan,
    unit: item.unit,
    cost: fmtMoney(item.cost),
    sellingMinor: null,
    currency: "USD",
    margin: item.margin,
    stock: item.stock,
    status: item.status,
  }));
  const rows = loading ? [] : liveRows ?? demoRows;
  const warningCount = rows.filter((item) => item.margin != null && item.margin < 20).length;

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Catalog / Pricing</h3>
          <div className="flex flex-wrap items-center justify-end gap-3">
            {notice && (
              <span className={`text-[11px] font-medium ${notice.type === "error" ? "text-red-600" : "text-green-700"}`}>
                {notice.text}
              </span>
            )}
            <span className="text-[11px] text-gray-400">{statusText(catalog.status, masterPlans.status, usingLive)}</span>
            <span className="text-[11px] text-amber-600 font-medium">{warningCount} margin warning(s)</span>
          </div>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[880px]">
            <thead>
              <tr className="bg-gray-50">
                {["Plan", "Unit", "Cost", "Your price", "Margin", "Source", "Status", "Action"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => (
                <tr key={item.key} className={`hover:bg-gray-50 border-b border-gray-100 last:border-0 ${item.margin != null && item.margin < 0 ? "bg-amber-50/60" : ""}`}>
                  <td className="p-4 font-medium text-gray-900">{item.plan}</td>
                  <td className="p-4 text-gray-400 text-[12px]">{item.unit}</td>
                  <td className="p-4 tabular-nums text-gray-500">{item.cost}</td>
                  <td className="p-4 tabular-nums font-medium">
                    {item.masterPlan && !item.tenantPlan ? (
                      <input
                        className="h-8 w-24 rounded border border-gray-200 px-2 text-right text-[12px] tabular-nums outline-none focus:border-gray-400"
                        inputMode="decimal"
                        min="0"
                        type="number"
                        value={draftPrice(item.masterPlan, priceDrafts)}
                        onChange={(event) => setPriceDrafts((current) => ({ ...current, [item.masterPlan!.id]: event.target.value }))}
                      />
                    ) : item.sellingMinor == null ? (
                      fmtMoney(RESELLER_CATALOG.find((row) => row.plan === item.plan)?.selling ?? 0)
                    ) : (
                      moneyMinor(item.sellingMinor, item.currency)
                    )}
                  </td>
                  <td className="p-4 tabular-nums">
                    <span className={marginClass(item.margin)}>
                      {item.margin == null ? "-" : `${item.margin < 0 ? "" : "+"}${item.margin}%`}
                    </span>
                  </td>
                  <td className="p-4 text-gray-500">{item.stock}</td>
                  <td className="p-4">
                    <span className="text-[11px] px-1.5 py-px rounded-sm bg-green-50 text-green-700 border border-green-200">{item.status}</span>
                  </td>
                  <td className="p-4">
                    {item.masterPlan && !item.tenantPlan ? (
                      <button
                        className="inline-flex h-8 items-center justify-center rounded-md border border-[#D50C2D] bg-[#D50C2D] px-3 text-[12px] font-medium text-white transition-colors hover:bg-[#B3082A] disabled:cursor-not-allowed disabled:border-gray-200 disabled:bg-gray-100 disabled:text-gray-400"
                        disabled={busyPlanID === item.masterPlan.id}
                        onClick={() => item.masterPlan && handleClone(item.masterPlan)}
                      >
                        {busyPlanID === item.masterPlan.id ? "Adding" : "Add"}
                      </button>
                    ) : item.tenantPlan ? (
                      <span className="text-[12px] text-gray-400">Added</span>
                    ) : (
                      <span className="text-[12px] text-gray-400">Read-only</span>
                    )}
                  </td>
                </tr>
              ))}
              {loading && (
                <tr><td colSpan={8} className="p-4 text-center text-[12px] text-gray-400">Loading catalog</td></tr>
              )}
              {usingLive && !loading && rows.length === 0 && (
                <tr><td colSpan={8} className="p-4 text-center text-[12px] text-gray-400">No catalog rows</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function snapshotText(value: unknown): string | null {
  if (!value || typeof value !== "object" || Array.isArray(value)) return null;
  const data = value as Record<string, unknown>;
  return typeof data.name === "string" ? data.name : typeof data.plan_code === "string" ? data.plan_code : null;
}
