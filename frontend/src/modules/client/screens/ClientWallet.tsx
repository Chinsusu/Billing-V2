"use client";

import type { FormEvent } from "react";
import { useState } from "react";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { walletLedgerReferenceLabel } from "@/lib/api/walletViewModels";
import { CLIENT_LEDGER } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

type Notice = { type: "success" | "error"; text: string };

const PAYMENT_METHODS = [
  { value: "bank_transfer", label: "Bank" },
  { value: "crypto", label: "Crypto" },
  { value: "manual", label: "Manual" },
  { value: "other", label: "Other" },
];

function amountToMinor(value: string): number | null {
  const parsed = Number.parseFloat(value);
  if (!Number.isFinite(parsed) || parsed <= 0) return null;
  return Math.round(parsed * 100);
}

export function ClientWallet() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [amount, setAmount] = useState("25.00");
  const [paymentMethod, setPaymentMethod] = useState("bank_transfer");
  const [paymentReference, setPaymentReference] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [notice, setNotice] = useState<Notice | null>(null);
  const wallets = useApiResource(billingApi.listClientWallets, `client-wallets:${refreshKey}`);
  const transactions = useApiResource(billingApi.listClientTransactions, `client-transactions:${refreshKey}`);
  const topups = useApiResource(
    () => billingApi.listClientTopupRequests({ limit: 5 }),
    `client-topups:${refreshKey}`,
  );
  const wallet = wallets.data?.[0];
  const ledger = useApiResource(
    () => wallet ? billingApi.listClientWalletLedger(wallet.id) : Promise.resolve([]),
    `client-ledger:${wallet?.id ?? "no-wallet"}:${refreshKey}`,
  );
  const usingLive = wallets.status === "success" && ledger.status === "success";
  const rows = usingLive
    ? (ledger.data ?? []).map((entry) => ({
        ts: compactDateTime(entry.created_at),
        type: entry.entry_type,
        amountMinor: entry.direction === "debit" ? -entry.amount_minor : entry.amount_minor,
        amountText: moneyMinor(entry.direction === "debit" ? -entry.amount_minor : entry.amount_minor, entry.currency),
        ref: walletLedgerReferenceLabel(entry),
        balanceText: moneyMinor(entry.balance_after_minor, entry.currency),
      }))
    : CLIENT_LEDGER.map((entry) => ({
        ts: entry.ts,
        type: entry.type,
        amountMinor: entry.amount,
        amountText: fmtMoney(entry.amount),
        ref: entry.ref,
        balanceText: fmtMoney(entry.balance),
      }));

  async function handleTopup(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!wallet) {
      setNotice({ type: "error", text: "No wallet is available." });
      return;
    }
    const amountMinor = amountToMinor(amount);
    if (!amountMinor) {
      setNotice({ type: "error", text: "Amount must be greater than zero." });
      return;
    }
    setSubmitting(true);
    setNotice(null);
    try {
      const request = await billingApi.createClientTopupRequest({
        wallet_id: wallet.id,
        amount_minor: amountMinor,
        currency: wallet.currency,
        payment_method: paymentMethod,
        payment_reference: paymentReference.trim() || undefined,
      });
      setNotice({ type: "success", text: `Top-up ${recordLabel(request.display_id, "TUP-")} submitted.` });
      setPaymentReference("");
      setRefreshKey((current) => current + 1);
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : "Top-up request failed.";
      setNotice({ type: "error", text: message });
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="bg-white border border-gray-200 rounded p-4 grid grid-cols-[1fr_420px] gap-6 items-start">
        <div>
          <div className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Available balance</div>
          <div className="text-3xl font-medium tabular-nums text-gray-900">
            {wallet ? moneyMinor(wallet.available_balance_minor, wallet.currency) : "$128.40"}
          </div>
          <div className="text-[12px] text-gray-400 mt-1">
            {wallet ? `${recordLabel(wallet.display_id, "WAL-")} - ${wallet.status}` : "Linh Tran - via ProxyVN"}
          </div>
          {notice && (
            <div className={`text-[12px] font-medium mt-4 ${notice.type === "error" ? "text-red-600" : "text-green-700"}`}>
              {notice.text}
            </div>
          )}
        </div>

        <form className="grid grid-cols-2 gap-3" onSubmit={handleTopup}>
          <label className="flex flex-col gap-1 text-[11px] uppercase tracking-wide text-gray-400">
            Amount
            <input
              className="h-9 rounded border border-gray-200 px-3 text-[13px] normal-case tracking-normal text-gray-900 outline-none focus:border-gray-400"
              inputMode="decimal"
              min="0"
              step="0.01"
              type="number"
              value={amount}
              onChange={(event) => setAmount(event.target.value)}
            />
          </label>
          <label className="flex flex-col gap-1 text-[11px] uppercase tracking-wide text-gray-400">
            Method
            <select
              className="h-9 rounded border border-gray-200 px-3 text-[13px] normal-case tracking-normal text-gray-900 outline-none focus:border-gray-400"
              value={paymentMethod}
              onChange={(event) => setPaymentMethod(event.target.value)}
            >
              {PAYMENT_METHODS.map((method) => (
                <option key={method.value} value={method.value}>{method.label}</option>
              ))}
            </select>
          </label>
          <label className="col-span-2 flex flex-col gap-1 text-[11px] uppercase tracking-wide text-gray-400">
            Reference
            <input
              className="h-9 rounded border border-gray-200 px-3 text-[13px] normal-case tracking-normal text-gray-900 outline-none focus:border-gray-400"
              placeholder="Transfer code"
              value={paymentReference}
              onChange={(event) => setPaymentReference(event.target.value)}
            />
          </label>
          <button
            className="col-span-2 inline-flex h-9 items-center justify-center rounded-md border border-[#D50C2D] bg-[#D50C2D] px-4 text-[13px] font-medium text-white transition-colors hover:bg-[#B3082A] disabled:cursor-not-allowed disabled:border-gray-200 disabled:bg-gray-100 disabled:text-gray-400"
            disabled={submitting || !wallet}
            type="submit"
          >
            {submitting ? "Submitting" : "Top up"}
          </button>
        </form>
      </div>

      {topups.status === "success" && (topups.data?.length ?? 0) > 0 && (
        <div className="bg-white border border-gray-200 rounded">
          <div className="p-4 border-b border-gray-100">
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Top-ups</h3>
          </div>
          <table className="w-full text-[13px] border-collapse">
            <tbody>
              {(topups.data ?? []).map((request) => (
                <tr key={request.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[#D50C2D] text-[12px]">{recordLabel(request.display_id, "TUP-")}</td>
                  <td className="p-4 tabular-nums font-medium">{moneyMinor(request.amount_minor, request.currency)}</td>
                  <td className="p-4 text-gray-500">{request.payment_method}</td>
                  <td className="p-4 text-gray-400">{compactDateTime(request.created_at)}</td>
                  <td className="p-4"><StatusBadge status={request.status} dot /></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Transaction history</h3>
          <div className="text-[11px] text-gray-400 mt-0.5">{transactions.data?.length ?? 0} payment records</div>
        </div>
        <table className="w-full text-[13px] border-collapse">
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
            {rows.map((entry, index) => (
              <tr key={`${entry.ts}:${entry.ref}:${index}`} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 tabular-nums text-gray-400">{entry.ts}</td>
                <td className="p-4 text-[12px] text-gray-500">{entry.type}</td>
                <td className={`p-4 tabular-nums text-right font-medium ${entry.amountMinor < 0 ? "text-red-600" : "text-green-700"}`}>
                  {entry.amountText}
                </td>
                <td className="p-4 text-gray-500">{entry.ref}</td>
                <td className="p-4 tabular-nums text-right font-medium">{entry.balanceText}</td>
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
