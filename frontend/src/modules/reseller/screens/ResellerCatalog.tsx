"use client";

import { billingApi } from "@/lib/api/billing";
import { moneyMinor } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { RESELLER_CATALOG } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function ResellerCatalog() {
  const catalog = useApiResource(billingApi.listResellerCatalog);
  const masterPlans = useApiResource(billingApi.listResellerMasterPlans);
  const loading = catalog.status === "loading" || masterPlans.status === "loading";
  const liveRows = catalog.status === "success" && catalog.data?.plans.length
    ? catalog.data.plans.map((plan) => ({
        key: plan.id,
        plan: snapshotText(plan.plan_snapshot) ?? plan.master_plan_id,
        unit: plan.visibility,
        cost: plan.reseller_cost_minor != null ? moneyMinor(plan.reseller_cost_minor, plan.currency) : "-",
        selling: moneyMinor(plan.selling_price_minor, plan.currency),
        margin: plan.reseller_cost_minor ? Math.round(((plan.selling_price_minor - plan.reseller_cost_minor) / plan.reseller_cost_minor) * 100) : null,
        stock: "catalog",
        status: plan.status,
      }))
    : masterPlans.status === "success"
      ? (masterPlans.data ?? []).map((plan) => ({
          key: plan.id,
          plan: plan.name,
          unit: `${plan.billing_cycle.value} ${plan.billing_cycle.type}`,
          cost: moneyMinor(plan.base_cost_minor, plan.currency),
          selling: moneyMinor(plan.suggested_price_minor, plan.currency),
          margin: plan.base_cost_minor ? Math.round(((plan.suggested_price_minor - plan.base_cost_minor) / plan.base_cost_minor) * 100) : null,
          stock: "master",
          status: plan.status,
        }))
      : null;
  const demoRows = RESELLER_CATALOG.map((item) => ({
    key: item.plan,
    plan: item.plan,
    unit: item.unit,
    cost: fmtMoney(item.cost),
    selling: fmtMoney(item.selling),
    margin: item.margin,
    stock: item.stock,
    status: item.status,
  }));
  const rows = loading ? [] : liveRows ?? demoRows;
  const warningCount = rows.filter((item) => item.margin != null && item.margin < 20).length;

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Catalog / Pricing</h3>
          <span className="text-[11px] text-amber-600 font-medium">{warningCount} margin warning(s)</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[760px]">
            <thead>
              <tr className="bg-gray-50">
                {["Plan", "Unit", "Cost", "Your price", "Margin", "Source", "Status"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => (
                <tr key={item.key} className={`hover:bg-gray-50 border-b border-gray-100 last:border-0 ${item.margin != null && item.margin < 0 ? "bg-amber-50/60" : ""}`}>
                  <td className="p-4 p-4 font-medium text-gray-900">{item.plan}</td>
                  <td className="p-4 p-4 text-gray-400 text-[12px]">{item.unit}</td>
                  <td className="p-4 p-4 tabular-nums text-gray-500">{item.cost}</td>
                  <td className="p-4 p-4 tabular-nums font-medium">{item.selling}</td>
                  <td className="p-4 p-4 tabular-nums">
                    {item.margin == null ? (
                      <span className="text-gray-400">-</span>
                    ) : (
                      <span className={item.margin < 0 ? "text-red-600 font-medium" : item.margin < 20 ? "text-amber-600" : "text-green-700"}>
                        {item.margin < 0 ? "" : "+"}{item.margin}%
                      </span>
                    )}
                  </td>
                  <td className="p-4 p-4 text-gray-500">{item.stock}</td>
                  <td className="p-4 p-4">
                    <span className="text-[11px] px-1.5 py-px rounded-sm bg-green-50 text-green-700 border border-green-200">{item.status}</span>
                  </td>
                </tr>
              ))}
              {loading && (
                <tr><td colSpan={7} className="p-4 text-center text-[12px] text-gray-400">Loading catalog</td></tr>
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
