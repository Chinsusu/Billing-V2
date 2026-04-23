"use client";

import { TRANSACTIONS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel, shortID } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";

export function AdminTransactions() {
  const transactions = useApiResource(billingApi.listAdminTransactions);
  const reconciliation = useApiResource(billingApi.listAdminReconciliation);
  const usingLive = transactions.status === "success";
  const rows = usingLive
    ? (transactions.data ?? []).map((tx) => ({
        id: recordLabel(tx.display_id, "TX-"),
        time: compactDateTime(tx.created_at),
        customer: shortID(tx.account_user_id),
        method: reconciliation.data?.find((item) => item.transaction.id === tx.id)?.provider ?? "wallet",
        type: tx.type,
        amount: moneyMinor(tx.amount_minor, tx.currency),
        status: tx.status,
      }))
    : TRANSACTIONS.map((tx) => ({
        ...tx,
        amount: fmtMoney(tx.amount),
      }));

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Transactions / Ledger</h3>
          <span className="text-[11px] text-gray-400">{reconciliation.data?.length ?? 0} reconciled</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Time", "Customer", "Method", "Type", "Amount", "Status"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((tx) => (
              <tr key={tx.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 text-[12px] text-gray-400">{tx.id}</td>
                <td className="p-4 p-4 text-gray-400 tabular-nums">{tx.time}</td>
                <td className="p-4 p-4 text-gray-700">{tx.customer}</td>
                <td className="p-4 p-4 text-gray-500">{tx.method}</td>
                <td className="p-4 p-4">
                  <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{tx.type}</span>
                </td>
                <td className="p-4 p-4 text-right tabular-nums font-medium">{tx.amount}</td>
                <td className="p-4 p-4"><StatusBadge status={tx.status} dot /></td>
              </tr>
            ))}
            {usingLive && rows.length === 0 && (
              <tr><td colSpan={7} className="p-4 text-center text-[12px] text-gray-400">No transactions</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
