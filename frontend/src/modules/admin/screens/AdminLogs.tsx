"use client";

import { FormEvent, useState } from "react";
import { AUDIT_LOGS, AuditLog } from "@/mocks/billingData";
import { billingApi } from "@/lib/api/billing";
import { AdminAuditLogQuery } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { mapAdminAuditLogView, type AdminAuditActorBadge, type AdminAuditLogView } from "@/lib/api/viewModels";
import { AdminFilterBar, AdminFilterInput } from "../components/AdminFilterBar";
import { hasActiveFilters, includesFilter, trimStringFilters } from "../lib/filterUtils";

const LEVEL_STYLE: Record<AuditLog["level"], string> = {
  info:  "bg-blue-50 text-blue-700",
  warn:  "bg-amber-50 text-amber-700",
  error: "bg-red-50 text-red-700",
};

type ActorBadge = AdminAuditActorBadge;
type AuditTableRow = AuditLog | AdminAuditLogView;

const ACTOR_STYLE: Record<ActorBadge, string> = {
  system:   "bg-gray-100 text-gray-500",
  admin:    "bg-purple-50 text-purple-700",
  reseller: "bg-indigo-50 text-indigo-700",
  client:   "bg-teal-50 text-teal-700",
  user: "bg-slate-100 text-slate-700",
  worker: "bg-orange-50 text-orange-700",
  provider_webhook: "bg-cyan-50 text-cyan-700",
};

type AuditLogFilters = Required<Pick<AdminAuditLogQuery, "display_id" | "actor_display_id" | "action" | "target_type" | "target_display_id">>;

const EMPTY_FILTERS: AuditLogFilters = {
  display_id: "",
  actor_display_id: "",
  action: "",
  target_type: "",
  target_display_id: "",
};

function filterMockLogs(filters: AuditLogFilters) {
  return AUDIT_LOGS.filter((log) => (
    includesFilter(log.id, filters.display_id)
    && includesFilter(log.actorName, filters.actor_display_id)
    && includesFilter(log.action, filters.action)
    && includesFilter(log.target, filters.target_type)
    && includesFilter(log.target, filters.target_display_id)
  ));
}

export function AdminLogs() {
  const [draftFilters, setDraftFilters] = useState(EMPTY_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(EMPTY_FILTERS);
  const logs = useApiResource(
    () => billingApi.listAdminAuditLogs(appliedFilters),
    JSON.stringify(appliedFilters),
  );
  const usingLive = logs.status === "success";
  const rows: AuditTableRow[] = usingLive
    ? (logs.data ?? []).map(mapAdminAuditLogView)
    : filterMockLogs(appliedFilters);
  const activeFilters = hasActiveFilters(appliedFilters);
  const statusTone = logs.status === "error"
    ? "error"
    : logs.status === "loading"
      ? "loading"
      : usingLive
        ? "success"
        : "default";
  const statusText = logs.status === "error"
    ? "Live API unavailable. Showing demo audit data for the current filters."
    : logs.status === "loading"
      ? "Refreshing live audit data..."
      : usingLive
        ? "Live audit filters applied."
        : activeFilters
          ? "Filters are applied to demo audit data."
          : "Demo audit data is active until the live API responds.";

  function updateFilter(field: keyof AuditLogFilters, value: string) {
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
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Audit logs</h3>
          <span className="text-[11px] text-gray-400">{rows.length} entries</span>
        </div>
        <AdminFilterBar onSubmit={applyFilters} onReset={resetFilters} statusText={statusText} statusTone={statusTone}>
          <AdminFilterInput
            label="Display ID"
            value={draftFilters.display_id}
            onChange={(event) => updateFilter("display_id", event.target.value)}
            placeholder="70001"
            inputMode="numeric"
          />
          <AdminFilterInput
            label="Actor public ID"
            value={draftFilters.actor_display_id}
            onChange={(event) => updateFilter("actor_display_id", event.target.value)}
            placeholder="10001"
            inputMode="numeric"
          />
          <AdminFilterInput
            label="Action"
            value={draftFilters.action}
            onChange={(event) => updateFilter("action", event.target.value)}
            placeholder="invoice.paid"
          />
          <AdminFilterInput
            label="Target Type"
            value={draftFilters.target_type}
            onChange={(event) => updateFilter("target_type", event.target.value)}
            placeholder="invoice"
          />
          <AdminFilterInput
            label="Target public ID"
            value={draftFilters.target_display_id}
            onChange={(event) => updateFilter("target_display_id", event.target.value)}
            placeholder="53001"
            inputMode="numeric"
          />
        </AdminFilterBar>
        <div className="overflow-x-auto">
          <table className="min-w-[920px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Time", "Level", "Actor", "Action", "Target", "Detail", "Request"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((log) => (
                <tr key={log.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 p-4 text-[11px] text-[#D50C2D]">{log.id}</td>
                  <td className="p-4 p-4 text-gray-400 text-[11px] tabular-nums whitespace-nowrap">{log.ts}</td>
                  <td className="p-4 p-4">
                    <span className={`text-[10px] font-medium uppercase px-1.5 py-0.5 rounded ${LEVEL_STYLE[log.level]}`}>
                      {log.level}
                    </span>
                  </td>
                  <td className="p-4 p-4 whitespace-nowrap">
                    <span className={`text-[10px] font-medium px-1.5 py-0.5 rounded mr-1 ${ACTOR_STYLE[log.actor]}`}>
                      {log.actor}
                    </span>
                    <span className="text-gray-600">{log.actorName}</span>
                  </td>
                  <td className="p-4 p-4 text-[11px] text-gray-700">{log.action}</td>
                  <td className="p-4 p-4 text-[11px] text-[#D50C2D]">{log.target}</td>
                  <td className="p-4 p-4 text-gray-500 max-w-[260px] truncate">{log.detail}</td>
                  <td className="p-4 p-4 text-[11px] text-gray-300">{log.requestId}</td>
                </tr>
              ))}
              {usingLive && rows.length === 0 && (
                <tr><td colSpan={8} className="p-4 text-center text-[12px] text-gray-400">No audit logs</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
