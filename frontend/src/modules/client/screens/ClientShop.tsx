"use client";

import { PRODUCTS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";

export function ClientShop() {
  const invoices = useApiResource(billingApi.listClientInvoices);
  const liveInvoices = invoices.status === "success" ? invoices.data ?? [] : null;

  if (liveInvoices) {
    return (
      <div className="p-4">
        <div className="bg-white border border-gray-200 rounded">
          <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Invoices</h3>
            <span className="text-[11px] text-gray-400">{liveInvoices.length} records</span>
          </div>
          <table className="w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["Invoice", "Issued", "Paid", "Amount", "Status"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {liveInvoices.map((inv) => (
                <tr key={inv.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 p-4 text-[12px] text-[#D50C2D]">{recordLabel(inv.display_id, "INV-")}</td>
                  <td className="p-4 p-4 text-gray-400">{compactDateTime(inv.issued_at)}</td>
                  <td className="p-4 p-4 text-gray-400">{compactDateTime(inv.paid_at)}</td>
                  <td className="p-4 p-4 text-right font-medium tabular-nums">{moneyMinor(inv.total_minor, inv.currency)}</td>
                  <td className="p-4 p-4"><StatusBadge status={inv.status} dot /></td>
                </tr>
              ))}
              {liveInvoices.length === 0 && (
                <tr><td colSpan={5} className="p-4 text-center text-[12px] text-gray-400">No invoices</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    );
  }

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-3 gap-4">
        {PRODUCTS.map((p) => (
          <div key={p.sku} className="bg-white border border-gray-200 rounded p-4 flex flex-col gap-4 hover:border-gray-300 transition-colors">
            <div>
              <div className="text-[13px] font-medium text-gray-900">{p.name}</div>
              <div className="text-[11px] text-gray-400 mt-0.5">{p.sku}</div>
            </div>
            <div className="flex items-baseline gap-1">
              <span className="text-lg font-medium tabular-nums text-gray-900">{fmtMoney(p.price)}</span>
              <span className="text-[12px] text-gray-400">{p.unit}</span>
            </div>
            <div className="text-[11px] text-gray-400">
              {p.active.toLocaleString()} active subscriptions
            </div>
            <button className="mt-auto w-full inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm">
              Order now
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
