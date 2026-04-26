"use client";

import { FormEvent, useState } from "react";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import type { AdminProviderSourceQuery } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { mapAdminProviderSourceView } from "@/lib/api/viewModels";
import { PROVIDERS } from "@/mocks/billingData";
import { AdminFilterBar, AdminFilterInput, AdminFilterSelect } from "../components/AdminFilterBar";
import { AdminProviderReadinessPanel } from "../components/AdminProviderReadinessPanel";
import { PROVIDER_SOURCE_TYPE_OPTIONS, PROVIDER_STATUS_OPTIONS } from "../lib/filterOptions";
import { equalsFilter, hasActiveFilters, includesFilter, trimStringFilters } from "../lib/filterUtils";

interface ProviderRow {
  id: string;
  name: string;
  type: string;
  status: string;
  location: string;
  inventory: string;
  risk: string;
  account: string;
  updated: string;
}

type ProviderFilterFields = Required<Pick<AdminProviderSourceQuery, "display_id" | "source_type" | "status">>;

const EMPTY_FILTERS: ProviderFilterFields = {
  display_id: "",
  source_type: "",
  status: "",
};

function sourceText(status: string, usingLive: boolean) {
  if (status === "error") return "Live API unavailable. Showing demo provider data.";
  if (status === "loading") return "Refreshing live provider sources...";
  return usingLive ? "Live provider sources" : "Demo provider data";
}

function filterDemoProviders(filters: ProviderFilterFields): ProviderRow[] {
  return PROVIDERS.map((provider) => ({
    id: provider.id,
    name: provider.name,
    type: provider.type,
    status: provider.health === "ok" ? "active" : provider.health === "degraded" ? "pending" : "failed",
    location: "-",
    inventory: `${provider.capacity}% capacity`,
    risk: `${provider.failRate}% fail rate`,
    account: "demo",
    updated: provider.lastSync,
  })).filter((provider) => (
    includesFilter(provider.id, filters.display_id)
    && includesFilter(provider.type, filters.source_type)
    && equalsFilter(provider.status, filters.status)
  ));
}

export function AdminProviders() {
  const [draftFilters, setDraftFilters] = useState(EMPTY_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(EMPTY_FILTERS);
  const providers = useApiResource(
    () => billingApi.listAdminProviderSources({ ...appliedFilters, limit: 100 }),
    `admin-provider-sources:${JSON.stringify(appliedFilters)}`,
  );
  const usingLive = providers.status === "success";
  const rows: ProviderRow[] = usingLive
    ? (providers.data ?? []).map(mapAdminProviderSourceView)
    : filterDemoProviders(appliedFilters);
  const activeFilters = hasActiveFilters(appliedFilters);
  const statusTone = providers.status === "error"
    ? "error"
    : providers.status === "loading"
      ? "loading"
      : usingLive
        ? "success"
        : "default";
  const statusText = sourceText(providers.status, usingLive);
  const filterStatusText = usingLive
    ? activeFilters
      ? "Live provider source filters applied."
      : statusText
    : activeFilters
      ? "Filters are applied to demo provider data."
      : statusText;

  function updateFilter(field: keyof ProviderFilterFields, value: string) {
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
    <div className="flex flex-col gap-4 p-4">
      <AdminProviderReadinessPanel />
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Providers / Sources</h3>
          <div className="flex flex-wrap items-center justify-end gap-3">
            <span className="text-[11px] text-gray-400">{statusText}</span>
            <span className="text-[11px] text-gray-400">{rows.length} sources</span>
          </div>
        </div>
        <AdminFilterBar onSubmit={applyFilters} onReset={resetFilters} statusText={filterStatusText} statusTone={statusTone}>
          <AdminFilterInput
            label="Source public ID"
            value={draftFilters.display_id}
            onChange={(event) => updateFilter("display_id", event.target.value)}
            placeholder="23001"
            inputMode="numeric"
          />
          <AdminFilterSelect
            label="Source type"
            value={draftFilters.source_type}
            onChange={(event) => updateFilter("source_type", event.target.value)}
            options={PROVIDER_SOURCE_TYPE_OPTIONS}
          />
          <AdminFilterSelect
            label="Status"
            value={draftFilters.status}
            onChange={(event) => updateFilter("status", event.target.value)}
            options={PROVIDER_STATUS_OPTIONS}
          />
        </AdminFilterBar>
        <div className="overflow-x-auto">
          <table className="min-w-[860px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Name", "Type", "Status", "Location", "Inventory", "Risk", "Account", "Updated"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((provider) => (
                <tr key={provider.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D]">{provider.id}</td>
                  <td className="p-4 font-medium text-gray-900">{provider.name}</td>
                  <td className="p-4">
                    <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{provider.type}</span>
                  </td>
                  <td className="p-4"><StatusBadge status={provider.status} dot /></td>
                  <td className="p-4 text-gray-500">{provider.location}</td>
                  <td className="p-4 text-gray-500">{provider.inventory}</td>
                  <td className="p-4 text-gray-500">{provider.risk}</td>
                  <td className="p-4 text-[12px] text-gray-400">{provider.account}</td>
                  <td className="p-4 text-gray-400">{provider.updated}</td>
                </tr>
              ))}
              {usingLive && rows.length === 0 && (
                <tr><td colSpan={9} className="p-4 text-center text-[12px] text-gray-400">No provider sources</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
