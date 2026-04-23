"use client";

import { CLIENT_LEDGER } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";

export function ClientWallet() {
  const wallets = useApiResource(billingApi.listClientWallets);
  const transactions = useApiResource(billingApi.listClientTransactions);
  const wallet = wallets.data?.[0];
  const ledger = useApiResource(
    () => wallet ? billingApi.listClientWalletLedger(wallet.id) : Promise.resolve([]),
    wallet?.id ?? "no-wallet",
  );
  const usingLive = ledger.status === "success";
  const rows = usingLive
    ? (ledger.data ?? []).map((entry) => ({
        ts: compactDateTime(entry.created_at),
        type: entry.entry_type,
        amount: entry.direction === "debit" ? -entry.amount_minor / 100 : entry.amount_minor / 100,
        ref: `${entry.reference_type} ${recordLabel(entry.display_id)}`,
        balance: entry.balance_after_minor / 100,
      }))
    : CLIENT_LEDGER;

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="bg-white border border-gray-200 rounded p-4 flex items-start justify-between">
        <div>
          <div className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Available balance</div>
          <div className="text-3xl font-medium tabular-nums text-gray-900">
            {wallet ? moneyMinor(wallet.available_balance_minor, wallet.currency) : "$128.40"}
          </div>
          <div className="text-[12px] text-gray-400 mt-1">Linh Tran · via ProxyVN</div>
        </div>
        <button className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm">
          + Top up
        </button>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Transaction history</h3>
          <div className="text-[11px] text-gray-400 mt-0.5">{transactions.data?.length ?? 0} payment records</div>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Timestamp", "Type", "Amount", "Reference", "Balance after"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((e, i) => (
              <tr key={i} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 tabular-nums text-gray-400">{e.ts}</td>
                <td className="p-4 p-4 text-[12px] text-gray-500">{e.type}</td>
                <td className={`p-4 p-4 tabular-nums text-right font-medium ${e.amount < 0 ? "text-red-600" : "text-green-700"}`}>
                  {fmtMoney(e.amount)}
                </td>
                <td className="p-4 p-4 text-gray-500">{e.ref}</td>
                <td className="p-4 p-4 tabular-nums text-right font-medium">{fmtMoney(e.balance)}</td>
              </tr>
            ))}
            {usingLive && rows.length === 0 && (
              <tr><td colSpan={5} className="p-4 text-center text-[12px] text-gray-400">No ledger entries</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
