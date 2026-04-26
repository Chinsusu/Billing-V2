"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { productTypeLabel } from "@/lib/api/displayLabels";
import { compactDateTime, moneyMinor, recordLabel, fmtMoney, fmtMoneyShort } from "@/lib/api/format";
import type { CatalogPlan } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { PRODUCTS } from "@/mocks/billingData";

interface ProductRow {
  id: string;
  name: string;
  type: string;
  plans: string;
  price: string;
  status: string;
  updated: string;
}

function lowestSuggestedPrice(plans: CatalogPlan[]) {
  if (plans.length === 0) return "-";
  const sorted = [...plans].sort((a, b) => a.suggested_price_minor - b.suggested_price_minor);
  return moneyMinor(sorted[0].suggested_price_minor, sorted[0].currency);
}

function sourceText(productStatus: string, planStatus: string, usingLive: boolean) {
  if (productStatus === "error") return "Live API unavailable. Showing demo catalog data.";
  if (productStatus === "loading" || planStatus === "loading") return "Refreshing live catalog data...";
  if (usingLive) return planStatus === "success" ? "Live products and plans" : "Live products, plans unavailable";
  return "Demo catalog data";
}

export function AdminProducts() {
  const products = useApiResource(
    () => billingApi.listAdminCatalogProducts({ limit: 100 }),
    "admin-catalog-products",
  );
  const plans = useApiResource(
    () => billingApi.listAdminCatalogPlans({ limit: 100 }),
    "admin-catalog-plans",
  );
  const usingLive = products.status === "success";
  const plansByProduct = new Map<string, CatalogPlan[]>();
  for (const plan of plans.data ?? []) {
    plansByProduct.set(plan.product_id, [...(plansByProduct.get(plan.product_id) ?? []), plan]);
  }

  const rows: ProductRow[] = usingLive
    ? (products.data ?? []).map((product) => {
        const productPlans = plansByProduct.get(product.id) ?? [];
        const activePlans = productPlans.filter((plan) => plan.status === "active").length;
        return {
          id: recordLabel(product.display_id, "PROD-"),
          name: product.name,
          type: productTypeLabel(product.product_type),
          plans: productPlans.length > 0 ? `${activePlans}/${productPlans.length} active` : "No plans",
          price: lowestSuggestedPrice(productPlans),
          status: product.status,
          updated: compactDateTime(product.updated_at),
        };
      })
    : PRODUCTS.map((product) => ({
        id: product.sku,
        name: product.name,
        type: product.unit,
        plans: `${product.active.toLocaleString()} active`,
        price: fmtMoney(product.price),
        status: "active",
        updated: `Rev 30d ${fmtMoneyShort(product.rev30)}`,
      }));
  const statusText = sourceText(products.status, plans.status, usingLive);

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Products & Pricing</h3>
          <div className="flex flex-wrap items-center justify-end gap-3">
            <span className="text-[11px] text-gray-400">{statusText}</span>
            <span className="text-[11px] text-gray-400">{rows.length} products</span>
            <span className="inline-flex h-8 items-center justify-center rounded-md border border-gray-200 bg-gray-50 px-3 text-[12px] font-medium text-gray-400">
              Read-only
            </span>
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-[760px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Name", "Type", "Plans", "Suggested Price", "Status", "Updated"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((product) => (
                <tr key={product.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D]">{product.id}</td>
                  <td className="p-4 font-medium text-gray-900">{product.name}</td>
                  <td className="p-4 text-gray-500 text-[12px]">{product.type}</td>
                  <td className="p-4 tabular-nums text-gray-500">{product.plans}</td>
                  <td className="p-4 tabular-nums font-medium">{product.price}</td>
                  <td className="p-4"><StatusBadge status={product.status} dot /></td>
                  <td className="p-4 text-gray-400">{product.updated}</td>
                </tr>
              ))}
              {usingLive && rows.length === 0 && (
                <tr><td colSpan={7} className="p-4 text-center text-[12px] text-gray-400">No catalog products</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
