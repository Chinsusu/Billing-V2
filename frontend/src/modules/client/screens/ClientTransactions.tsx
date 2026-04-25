"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { hiddenReference } from "@/lib/api/viewModels";

export function ClientTransactions() {
  const transactions = useApiResource(billingApi.listClientTransactions);
  const rows = transactions.data ?? [];
  const postedTotal = rows
    .filter((txn) => txn.status === "posted" || txn.status === "paid")
    .reduce((total, txn) => total + txn.amount_minor, 0);
  const failedCount = rows.filter((txn) => txn.status === "failed").length;
  const currency = rows[0]?.currency ?? "USD";

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <SummaryTile label="Transactions" value={String(rows.length)} />
        <SummaryTile label="Posted amount" value={moneyMinor(postedTotal, currency)} />
        <SummaryTile label="Failed" value={String(failedCount)} tone={failedCount > 0 ? "warn" : "neutral"} />
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-3">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Transactions</h3>
          <span className="text-[11px] text-gray-400">{transactions.status === "loading" ? "Loading" : `${rows.length} records`}</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[760px]">
            <thead>
              <tr className="bg-gray-50">
                {["Transaction", "Created", "Type", "Reference", "Description", "Amount", "Status"].map((heading) => (
                  <th key={heading} className="text-left text-[11px] font-medium uppercase text-gray-400 p-4 border-b border-gray-200">
                    {heading}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((txn) => (
                <tr key={txn.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D] font-medium">{recordLabel(txn.display_id, "TXN-")}</td>
                  <td className="p-4 text-gray-500">{compactDateTime(txn.created_at)}</td>
                  <td className="p-4 text-gray-500">{txn.type}</td>
                  <td className="p-4 text-[12px] text-gray-400">{transactionReference(txn)}</td>
                  <td className="p-4 text-gray-500 max-w-[260px] truncate">{txn.description ?? "-"}</td>
                  <td className="p-4 text-right font-medium tabular-nums">{moneyMinor(txn.amount_minor, txn.currency)}</td>
                  <td className="p-4"><StatusBadge status={txn.status} dot /></td>
                </tr>
              ))}
              {transactions.status === "loading" && <TableMessage colSpan={7} text="Loading transactions" />}
              {transactions.status === "error" && <TableMessage colSpan={7} text={transactions.error ?? "Transactions unavailable"} tone="error" />}
              {transactions.status === "success" && rows.length === 0 && <TableMessage colSpan={7} text="No transactions" />}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function transactionReference(
  transaction: { invoice_display_id?: number; order_display_id?: number },
): string {
  if (transaction.invoice_display_id) return recordLabel(transaction.invoice_display_id, "INV-");
  if (transaction.order_display_id) return recordLabel(transaction.order_display_id, "ORD-");
  return hiddenReference("Reference");
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
