"use client";

import { FormEvent, useState } from "react";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import type { AdminAccountQuery } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { mapAdminAccountView } from "@/lib/api/viewModels";
import { CUSTOMERS } from "@/mocks/billingData";
import { fmtMoneyShort } from "@/mocks/sampleData";
import { AdminFilterBar, AdminFilterInput } from "../components/AdminFilterBar";
import { equalsFilter, hasActiveFilters, includesFilter, trimStringFilters } from "../lib/filterUtils";

interface CustomerRow {
  id: string;
  name: string;
  email: string;
  type: string;
  tenant: string;
  security: string;
  status: string;
  created: string;
  lastLogin: string;
}

type CustomerFilterFields = Required<Pick<AdminAccountQuery, "display_id" | "email" | "type" | "status">>;

const EMPTY_FILTERS: CustomerFilterFields = {
  display_id: "",
  email: "",
  type: "",
  status: "",
};

function sourceText(status: string, usingLive: boolean) {
  if (status === "error") return "Live API unavailable. Showing demo customer data.";
  if (status === "loading") return "Refreshing live customers...";
  return usingLive ? "Live customer accounts" : "Demo customer data";
}

function filterDemoCustomers(filters: CustomerFilterFields): CustomerRow[] {
  return CUSTOMERS.map((customer) => ({
    id: customer.id,
    name: customer.name,
    email: customer.email,
    type: customer.plan,
    tenant: customer.country,
    security: `${customer.services} services`,
    status: customer.status,
    created: customer.since,
    lastLogin: fmtMoneyShort(customer.mrr),
  })).filter((customer) => (
    includesFilter(customer.id, filters.display_id)
    && includesFilter(customer.email, filters.email)
    && includesFilter(customer.type, filters.type)
    && equalsFilter(customer.status, filters.status)
  ));
}

export function AdminCustomers() {
  const [draftFilters, setDraftFilters] = useState(EMPTY_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(EMPTY_FILTERS);
  const customers = useApiResource(
    () => billingApi.listAdminCustomers({ ...appliedFilters, limit: 100 }),
    `admin-customers:${JSON.stringify(appliedFilters)}`,
  );
  const usingLive = customers.status === "success";
  const rows: CustomerRow[] = usingLive
    ? (customers.data ?? []).map(mapAdminAccountView)
    : filterDemoCustomers(appliedFilters);
  const activeFilters = hasActiveFilters(appliedFilters);
  const statusTone = customers.status === "error"
    ? "error"
    : customers.status === "loading"
      ? "loading"
      : usingLive
        ? "success"
        : "default";
  const statusText = sourceText(customers.status, usingLive);
  const filterStatusText = usingLive
    ? activeFilters
      ? "Live customer filters applied."
      : statusText
    : activeFilters
      ? "Filters are applied to demo customer data."
      : statusText;

  function updateFilter(field: keyof CustomerFilterFields, value: string) {
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
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Customers</h3>
          <div className="flex flex-wrap items-center justify-end gap-3">
            <span className="text-[11px] text-gray-400">{sourceText(customers.status, usingLive)}</span>
            <span className="text-[11px] text-gray-400">{rows.length.toLocaleString()} total</span>
          </div>
        </div>
        <AdminFilterBar onSubmit={applyFilters} onReset={resetFilters} statusText={filterStatusText} statusTone={statusTone}>
          <AdminFilterInput
            label="Account public ID"
            value={draftFilters.display_id}
            onChange={(event) => updateFilter("display_id", event.target.value)}
            placeholder="20002"
            inputMode="numeric"
          />
          <AdminFilterInput
            label="Email"
            value={draftFilters.email}
            onChange={(event) => updateFilter("email", event.target.value)}
            placeholder="buyer@example.com"
          />
          <AdminFilterInput
            label="Type"
            value={draftFilters.type}
            onChange={(event) => updateFilter("type", event.target.value)}
            placeholder="client, reseller"
          />
          <AdminFilterInput
            label="Status"
            value={draftFilters.status}
            onChange={(event) => updateFilter("status", event.target.value)}
            placeholder="active, suspended"
          />
        </AdminFilterBar>
        <div className="overflow-x-auto">
          <table className="min-w-[860px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Name", "Email", "Type", "Tenant", "2FA", "Status", "Created", "Last login"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((customer) => (
                <tr key={customer.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D]">{customer.id}</td>
                  <td className="p-4 font-medium text-gray-900">{customer.name}</td>
                  <td className="p-4 text-gray-400 text-[12px]">{customer.email}</td>
                  <td className="p-4">
                    <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{customer.type}</span>
                  </td>
                  <td className="p-4 text-gray-500 text-[12px]">{customer.tenant}</td>
                  <td className="p-4 text-gray-500 text-[12px]">{customer.security}</td>
                  <td className="p-4"><StatusBadge status={customer.status} dot /></td>
                  <td className="p-4 text-gray-400">{customer.created}</td>
                  <td className="p-4 text-gray-400">{customer.lastLogin}</td>
                </tr>
              ))}
              {usingLive && rows.length === 0 && (
                <tr><td colSpan={9} className="p-4 text-center text-[12px] text-gray-400">No customer accounts</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
