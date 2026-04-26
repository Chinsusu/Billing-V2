"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { walletLedgerEntryTypeLabel, walletLedgerReferenceLabel } from "@/lib/api/walletViewModels";
import { TOPUP_REQUESTS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

const DEMO_LEDGER = [
  { ts: "2026-04-22 14:02", type: "settlement.reseller.debit", amount: -62.00, ref: "ORD-48291 / VPS 4C/8G", balance: 4820.50 },
  { ts: "2026-04-21 10:18", type: "topup.credit.reseller", amount: 2000.00, ref: "TUP-9116 / VietQR", balance: 4882.50 },
  { ts: "2026-04-20 14:08", type: "settlement.reseller.debit", amount: -390.00, ref: "ORD-48280 / Residential batch", balance: 2882.50 },
  { ts: "2026-04-18 11:22", type: "settlement.reseller.debit", amount: -180.00, ref: "ORD-48270 / ISP batch", balance: 3272.50 },
];

export function ResellerWallet() {
  const wallets = useApiResource(
    () => billingApi.listResellerWallets({ limit: 20 }),
    "reseller-wallets",
  );
  const wallet = wallets.data?.[0];
  const ledger = useApiResource(
    () => wallet ? billingApi.listResellerWalletLedger(wallet.id, { limit: 50 }) : Promise.resolve([]),
    `reseller-ledger:${wallet?.id ?? "no-wallet"}`,
  );
  const topups = useApiResource(
    () => billingApi.listResellerTopupRequests({ limit: 20 }),
    "reseller-wallet-topups",
  );
  const usingLive = wallets.status === "success" && ledger.status === "success";
  const ledgerRows = usingLive
    ? (ledger.data ?? []).map((entry) => {
        const signedMinor = entry.direction === "debit" ? -entry.amount_minor : entry.amount_minor;
        return {
          ts: compactDateTime(entry.created_at),
          type: walletLedgerEntryTypeLabel(entry.entry_type),
          amountMinor: signedMinor,
          amount: moneyMinor(signedMinor, entry.currency),
          ref: walletLedgerReferenceLabel(entry),
          balance: moneyMinor(entry.balance_after_minor, entry.currency),
        };
      })
    : DEMO_LEDGER.map((entry) => ({
        ts: entry.ts,
        type: walletLedgerEntryTypeLabel(entry.type),
        amountMinor: Math.round(entry.amount * 100),
        amount: fmtMoney(entry.amount),
        ref: entry.ref,
        balance: fmtMoney(entry.balance),
      }));
  const topupRows = topups.status === "success"
    ? (topups.data ?? []).map((request) => ({
        id: recordLabel(request.display_id, "TUP-"),
        amount: moneyMinor(request.amount_minor, request.currency),
        amountMinor: request.amount_minor,
        method: request.payment_method,
        created: compactDateTime(request.created_at),
        status: request.status,
      }))
    : TOPUP_REQUESTS.filter((request) => request.tenant === "ProxyVN").map((request) => ({
        id: request.id,
        amount: fmtMoney(request.amount),
        amountMinor: Math.round(request.amount * 100),
        method: request.method,
        created: request.created,
        status: request.status,
      }));
  const pendingTopups = topupRows.filter((row) => row.status.includes("pending") || row.status.includes("review"));
  const spentMinor = ledgerRows
    .filter((row) => row.amountMinor < 0)
    .reduce((total, row) => total + Math.abs(row.amountMinor), 0);
  const source = wallets.status === "error" || ledger.status === "error"
    ? "Live API unavailable. Showing demo wallet data."
    : wallets.status === "loading" || ledger.status === "loading"
      ? "Refreshing live wallet..."
      : usingLive
        ? "Live reseller wallet"
        : "Demo wallet data";

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="bg-white border border-gray-200 rounded p-4">
        <div className="flex items-start justify-between gap-4">
          <div>
            <div className="text-[11px] font-medium uppercase tracking-wide text-gray-400 mb-1">Wallet balance</div>
            <div className="text-3xl font-medium tabular-nums text-gray-900">
              {wallet ? moneyMinor(wallet.available_balance_minor, wallet.currency) : "$4,820.50"}
            </div>
            <div className="text-[12px] text-gray-400 mt-1">
              {wallet ? `${recordLabel(wallet.display_id, "WAL-")} / ${wallet.status}` : "ProxyVN / demo wallet"}
            </div>
          </div>
          <span className="text-[11px] text-gray-400 text-right">{source}</span>
        </div>
        <div className="mt-4 pt-4 border-t border-gray-100 grid grid-cols-3 gap-4">
          <SummaryItem
            label="Pending top-ups"
            value={moneyMinor(pendingTopups.reduce((total, row) => total + row.amountMinor, 0), wallet?.currency ?? "USD")}
            sub={`${pendingTopups.length} awaiting admin`}
          />
          <SummaryItem
            label="Spent this month"
            value={moneyMinor(spentMinor, wallet?.currency ?? "USD")}
            sub="settlement debits"
          />
          <SummaryItem
            label="Low balance alert"
            value={wallet && wallet.available_balance_minor < 20000 ? "Active" : "Normal"}
            sub="threshold $200"
          />
        </div>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Top-up requests</h3>
          <span className="text-[11px] text-gray-400">
            {topups.status === "success" ? "Live top-ups" : "Demo top-ups"}
          </span>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-[620px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Amount", "Method", "Created", "Status"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {topupRows.map((row) => (
                <tr key={row.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D]">{row.id}</td>
                  <td className="p-4 tabular-nums font-medium">{row.amount}</td>
                  <td className="p-4 text-gray-500">{row.method}</td>
                  <td className="p-4 text-gray-400">{row.created}</td>
                  <td className="p-4"><StatusBadge status={row.status} dot /></td>
                </tr>
              ))}
              {topups.status === "success" && topupRows.length === 0 && (
                <tr><td colSpan={5} className="p-4 text-center text-[12px] text-gray-400">No top-up requests</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Ledger history</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-[760px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["Timestamp", "Type", "Amount", "Reference", "Balance after"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {ledgerRows.map((entry, index) => (
                <tr key={`${entry.ts}:${entry.ref}:${index}`} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 tabular-nums text-gray-400">{entry.ts}</td>
                  <td className="p-4 text-[12px] text-gray-500">{entry.type}</td>
                  <td className={`p-4 tabular-nums text-right font-medium ${entry.amountMinor < 0 ? "text-red-600" : "text-green-700"}`}>
                    {entry.amount}
                  </td>
                  <td className="p-4 text-gray-500">{entry.ref}</td>
                  <td className="p-4 tabular-nums text-right font-medium">{entry.balance}</td>
                </tr>
              ))}
              {usingLive && ledgerRows.length === 0 && (
                <tr><td colSpan={5} className="p-4 text-center text-[12px] text-gray-400">No ledger entries</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function SummaryItem({ label, value, sub }: { label: string; value: string; sub: string }) {
  return (
    <div>
      <div className="text-[11px] text-gray-400 mb-0.5">{label}</div>
      <div className="text-[14px] font-medium tabular-nums">{value}</div>
      <div className="text-[11px] text-gray-400">{sub}</div>
    </div>
  );
}
