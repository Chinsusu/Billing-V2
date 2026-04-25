"use client";

import { FormEvent, useState } from "react";
import { TRANSACTIONS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";
import { billingApi } from "@/lib/api/billing";
import { AdminTransactionQuery } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { mapAdminTransactionView } from "@/lib/api/viewModels";
import { AdminFilterBar, AdminFilterInput } from "../components/AdminFilterBar";
import { equalsFilter, hasActiveFilters, includesFilter, matchesAmountRange, trimStringFilters } from "../lib/filterUtils";

const EMPTY_FILTERS: Required<AdminTransactionQuery> = {
  display_id: "",
  account_user_id: "",
  status: "",
  amount_min: "",
  amount_max: "",
};

function filterMockTransactions(filters: Required<AdminTransactionQuery>) {
  return TRANSACTIONS.filter((transaction) => (
    includesFilter(transaction.id, filters.display_id)
    && includesFilter(transaction.customer, filters.account_user_id)
    && equalsFilter(transaction.status, filters.status)
    && matchesAmountRange(transaction.amount, filters.amount_min, filters.amount_max)
  ));
}

export function AdminTransactions() {
  const [draftFilters, setDraftFilters] = useState(EMPTY_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(EMPTY_FILTERS);
  const transactions = useApiResource(
    () => billingApi.listAdminTransactions(appliedFilters),
    `transactions:${JSON.stringify(appliedFilters)}`,
  );
  const reconciliation = useApiResource(
    () => billingApi.listAdminReconciliation(appliedFilters),
    `reconciliation:${JSON.stringify(appliedFilters)}`,
  );
  const usingLive = transactions.status === "success";
  const reconciliationByTransactionID = new Map(
    (reconciliation.data ?? []).map((item) => [item.transaction.id, item]),
  );
  const rows = usingLive
    ? (transactions.data ?? []).map((tx) => mapAdminTransactionView(tx, reconciliationByTransactionID.get(tx.id)))
    : filterMockTransactions(appliedFilters).map((tx) => ({
        ...tx,
        amount: fmtMoney(tx.amount),
      }));
  const activeFilters = hasActiveFilters(appliedFilters);
  const statusTone = transactions.status === "error"
    ? "error"
    : transactions.status === "loading"
      ? "loading"
      : usingLive
        ? "success"
        : "default";
  const statusText = transactions.status === "error"
    ? "Live API unavailable. Showing demo transaction data for the current filters."
    : transactions.status === "loading"
      ? "Refreshing live transaction data..."
      : usingLive
        ? "Live transaction filters applied."
        : activeFilters
          ? "Filters are applied to demo transaction data."
          : "Demo transaction data is active until the live API responds.";
  const reconciliationCount = usingLive ? (reconciliation.data?.length ?? 0) : rows.length;

  function updateFilter(field: keyof AdminTransactionQuery, value: string) {
    setDraftFilters((current) => ({ ...current, [field]: value }));
  }

  function applyFilters(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setAppliedFilters(trimStringFilters(draftFilters));
  }

  function resetFilters() {
    setDraftFilters(EMPTY_FILTERS);
    setAppliedFilters(EMPTY_FILTERS);
  }

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Transactions / Ledger</h3>
          <span className="text-[11px] text-gray-400">{reconciliationCount} reconciled</span>
        </div>
        <AdminFilterBar onSubmit={applyFilters} onReset={resetFilters} statusText={statusText} statusTone={statusTone}>
          <AdminFilterInput
            label="Display ID"
            value={draftFilters.display_id}
            onChange={(event) => updateFilter("display_id", event.target.value)}
            placeholder="51001"
            inputMode="numeric"
          />
          <AdminFilterInput
            label="Customer / account"
            value={draftFilters.account_user_id}
            onChange={(event) => updateFilter("account_user_id", event.target.value)}
            placeholder="account reference"
          />
          <AdminFilterInput
            label="Status"
            value={draftFilters.status}
            onChange={(event) => updateFilter("status", event.target.value)}
            placeholder="posted, paid, failed"
          />
          <AdminFilterInput
            label="Amount Min"
            type="number"
            step="0.01"
            value={draftFilters.amount_min}
            onChange={(event) => updateFilter("amount_min", event.target.value)}
            placeholder="100.00"
          />
          <AdminFilterInput
            label="Amount Max"
            type="number"
            step="0.01"
            value={draftFilters.amount_max}
            onChange={(event) => updateFilter("amount_max", event.target.value)}
            placeholder="5000.00"
          />
        </AdminFilterBar>
        <div className="overflow-x-auto">
          <table className="min-w-[820px] w-full text-[13px] border-collapse">
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
    </div>
  );
}
