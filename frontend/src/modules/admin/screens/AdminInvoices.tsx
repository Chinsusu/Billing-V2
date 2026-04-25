"use client";

import { FormEvent, useState } from "react";
import { INVOICES } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";
import { billingApi } from "@/lib/api/billing";
import { AdminInvoiceQuery } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { mapAdminInvoiceView } from "@/lib/api/viewModels";
import { AdminFilterBar, AdminFilterInput } from "../components/AdminFilterBar";
import { equalsFilter, hasActiveFilters, includesFilter, matchesAmountRange, trimStringFilters } from "../lib/filterUtils";

const EMPTY_FILTERS: Required<AdminInvoiceQuery> = {
  display_id: "",
  buyer_user_id: "",
  status: "",
  amount_min: "",
  amount_max: "",
};

function filterMockInvoices(filters: Required<AdminInvoiceQuery>) {
  return INVOICES.filter((invoice) => (
    includesFilter(invoice.id, filters.display_id)
    && includesFilter(invoice.customer, filters.buyer_user_id)
    && equalsFilter(invoice.status, filters.status)
    && matchesAmountRange(invoice.amount, filters.amount_min, filters.amount_max)
  ));
}

export function AdminInvoices() {
  const [draftFilters, setDraftFilters] = useState(EMPTY_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(EMPTY_FILTERS);
  const invoices = useApiResource(
    () => billingApi.listAdminInvoices(appliedFilters),
    JSON.stringify(appliedFilters),
  );
  const liveInvoices = invoices.data ?? [];
  const usingLive = invoices.status === "success";
  const rows = usingLive
    ? liveInvoices.map(mapAdminInvoiceView)
    : filterMockInvoices(appliedFilters).map((inv) => ({
        ...inv,
        amount: fmtMoney(inv.amount),
      }));
  const activeFilters = hasActiveFilters(appliedFilters);
  const statusTone = invoices.status === "error"
    ? "error"
    : invoices.status === "loading"
      ? "loading"
      : usingLive
        ? "success"
        : "default";
  const statusText = invoices.status === "error"
    ? "Live API unavailable. Showing demo data for the current filters."
    : invoices.status === "loading"
      ? "Refreshing live invoice data..."
      : usingLive
        ? "Live invoice filters applied."
        : activeFilters
          ? "Filters are applied to demo data."
          : "Demo data is active until the live API responds.";

  function updateFilter(field: keyof AdminInvoiceQuery, value: string) {
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
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Invoices</h3>
          <span className="text-[11px] text-gray-400">{rows.length} records</span>
        </div>
        <AdminFilterBar onSubmit={applyFilters} onReset={resetFilters} statusText={statusText} statusTone={statusTone}>
          <AdminFilterInput
            label="Display ID"
            value={draftFilters.display_id}
            onChange={(event) => updateFilter("display_id", event.target.value)}
            placeholder="44001"
            inputMode="numeric"
          />
          <AdminFilterInput
            label="Customer / account"
            value={draftFilters.buyer_user_id}
            onChange={(event) => updateFilter("buyer_user_id", event.target.value)}
            placeholder="account reference"
          />
          <AdminFilterInput
            label="Status"
            value={draftFilters.status}
            onChange={(event) => updateFilter("status", event.target.value)}
            placeholder="open, paid, overdue"
          />
          <AdminFilterInput
            label="Amount Min"
            type="number"
            min="0"
            step="0.01"
            value={draftFilters.amount_min}
            onChange={(event) => updateFilter("amount_min", event.target.value)}
            placeholder="100.00"
          />
          <AdminFilterInput
            label="Amount Max"
            type="number"
            min="0"
            step="0.01"
            value={draftFilters.amount_max}
            onChange={(event) => updateFilter("amount_max", event.target.value)}
            placeholder="5000.00"
          />
        </AdminFilterBar>
        <div className="overflow-x-auto">
          <table className="min-w-[760px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["Invoice", "Customer", "Issued", "Due", "Amount", "Status"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((inv) => (
                <tr key={inv.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 p-4 text-[12px] text-[#D50C2D]">{inv.id}</td>
                  <td className="p-4 p-4 text-gray-700">{inv.customer}</td>
                  <td className="p-4 p-4 text-gray-400">{inv.issued}</td>
                  <td className="p-4 p-4 text-gray-400">{inv.due}</td>
                  <td className="p-4 p-4 text-right font-medium tabular-nums">{inv.amount}</td>
                  <td className="p-4 p-4"><StatusBadge status={inv.status} dot /></td>
                </tr>
              ))}
              {usingLive && rows.length === 0 && (
                <tr><td colSpan={6} className="p-4 text-center text-[12px] text-gray-400">No invoices</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
