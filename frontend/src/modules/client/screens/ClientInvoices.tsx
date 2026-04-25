"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { hiddenReference } from "@/lib/api/viewModels";

export function ClientInvoices() {
  const invoices = useApiResource(billingApi.listClientInvoices);
  const rows = invoices.data ?? [];
  const openCount = rows.filter((invoice) => invoice.status === "issued" || invoice.status === "overdue").length;
  const paidTotal = rows
    .filter((invoice) => invoice.status === "paid")
    .reduce((total, invoice) => total + invoice.total_minor, 0);
  const currency = rows[0]?.currency ?? "USD";

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <SummaryTile label="Total invoices" value={String(rows.length)} />
        <SummaryTile label="Open invoices" value={String(openCount)} tone={openCount > 0 ? "warn" : "neutral"} />
        <SummaryTile label="Paid total" value={moneyMinor(paidTotal, currency)} />
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-3">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Invoices</h3>
          <span className="text-[11px] text-gray-400">{invoices.status === "loading" ? "Loading" : `${rows.length} records`}</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[720px]">
            <thead>
              <tr className="bg-gray-50">
                {["Invoice", "Order", "Issued", "Due", "Paid", "Amount", "Status"].map((heading) => (
                  <th key={heading} className="text-left text-[11px] font-medium uppercase text-gray-400 p-4 border-b border-gray-200">
                    {heading}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((invoice) => (
                <tr key={invoice.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D] font-medium">{recordLabel(invoice.display_id, "INV-")}</td>
                  <td className="p-4 text-[12px] text-gray-400">
                    {invoice.order_display_id ? recordLabel(invoice.order_display_id, "ORD-") : hiddenReference("Order")}
                  </td>
                  <td className="p-4 text-gray-500">{compactDateTime(invoice.issued_at)}</td>
                  <td className="p-4 text-gray-500">{compactDateTime(invoice.due_at)}</td>
                  <td className="p-4 text-gray-500">{compactDateTime(invoice.paid_at)}</td>
                  <td className="p-4 text-right font-medium tabular-nums">{moneyMinor(invoice.total_minor, invoice.currency)}</td>
                  <td className="p-4"><StatusBadge status={invoice.status} dot /></td>
                </tr>
              ))}
              {invoices.status === "loading" && <TableMessage colSpan={7} text="Loading invoices" />}
              {invoices.status === "error" && <TableMessage colSpan={7} text={invoices.error ?? "Invoices unavailable"} tone="error" />}
              {invoices.status === "success" && rows.length === 0 && <TableMessage colSpan={7} text="No invoices" />}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function SummaryTile({ label, value, tone = "neutral" }: { label: string; value: string; tone?: "neutral" | "warn" }) {
  return (
    <div className={`bg-white border rounded p-4 ${tone === "warn" ? "border-amber-200" : "border-gray-200"}`}>
      <div className="text-[11px] text-gray-400 uppercase mb-1">{label}</div>
      <div className={`text-lg font-medium tabular-nums ${tone === "warn" ? "text-amber-700" : "text-gray-900"}`}>{value}</div>
    </div>
  );
}

function TableMessage({ colSpan, text, tone = "neutral" }: { colSpan: number; text: string; tone?: "neutral" | "error" }) {
  return (
    <tr>
      <td colSpan={colSpan} className={`p-4 text-center text-[12px] ${tone === "error" ? "text-red-600" : "text-gray-400"}`}>
        {text}
      </td>
    </tr>
  );
}
