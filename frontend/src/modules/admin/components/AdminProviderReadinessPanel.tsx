import { billingApi } from "@/lib/api/billing";
import { recordLabel } from "@/lib/api/format";
import type { ProviderReadiness, ProviderReadinessState } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";

const DEMO_READINESS: ProviderReadiness[] = [
  {
    plan_display_id: 21001,
    plan_code: "vps-linux-small",
    plan_name: "VPS Linux Small",
    product_type: "vps",
    plan_status: "active",
    plan_source_display_id: 22001,
    plan_source_status: "active",
    source_display_id: 23001,
    source_name: "Hetzner Falkenstein",
    source_type: "hetzner",
    source_status: "active",
    inventory_mode: "provider_live",
    state: "ready",
    reason: "Source is active and supports automatic provisioning.",
  },
  {
    plan_display_id: 21002,
    plan_code: "proxy-residential",
    plan_name: "Residential Proxy",
    product_type: "proxy",
    plan_status: "active",
    plan_source_display_id: 22002,
    plan_source_status: "active",
    source_display_id: 23002,
    source_name: "Manual pool",
    source_type: "manual",
    source_status: "active",
    inventory_mode: "manual_unlimited",
    state: "fake_provider_only",
    reason: "Manual source only works with the local fake provider path.",
  },
  {
    plan_display_id: 21003,
    plan_code: "proxy-dc-shared",
    plan_name: "Datacenter Shared",
    product_type: "proxy",
    plan_status: "active",
    plan_source_display_id: 22003,
    plan_source_status: "active",
    source_display_id: 23003,
    source_name: "VPS-only source",
    source_type: "proxmox",
    source_status: "active",
    inventory_mode: "provider_live",
    state: "unsupported_capability",
    reason: "Source does not support automatic provisioning for this product type.",
  },
];

const STATE_LABEL: Record<ProviderReadinessState, string> = {
  ready: "Ready",
  inactive_source: "Inactive source",
  missing_plan_source: "Missing source",
  unsupported_capability: "Unsupported",
  fake_provider_only: "Fake only",
};

const STATE_CLASS: Record<ProviderReadinessState, string> = {
  ready: "border-emerald-200 bg-emerald-50 text-emerald-700",
  inactive_source: "border-amber-200 bg-amber-50 text-amber-700",
  missing_plan_source: "border-red-200 bg-red-50 text-red-700",
  unsupported_capability: "border-red-200 bg-red-50 text-red-700",
  fake_provider_only: "border-blue-200 bg-blue-50 text-blue-700",
};

export function AdminProviderReadinessPanel() {
  const readiness = useApiResource(
    () => billingApi.listAdminProviderReadiness({ status: "active", limit: 100 }),
    "admin-provider-readiness",
  );
  const usingLive = readiness.status === "success";
  const rows = usingLive ? readiness.data ?? [] : DEMO_READINESS;
  const attentionCount = rows.filter((row) => row.state !== "ready").length;
  const statusText = readinessStatusText(readiness.status, usingLive);

  return (
    <section className="bg-white border border-gray-200 rounded">
      <div className="flex flex-col gap-3 border-b border-gray-100 p-4 md:flex-row md:items-center md:justify-between">
        <div className="min-w-0">
          <h3 className="m-0 text-[13px] font-medium text-gray-900">Provider readiness</h3>
          <p className="m-0 mt-1 text-[12px] text-gray-400">
            {rows.length} plan(s), {attentionCount} need attention
          </p>
        </div>
        <span className={`inline-flex w-fit items-center rounded border px-2.5 py-1 text-[11px] ${statusTone(readiness.status, usingLive)}`}>
          {statusText}
        </span>
      </div>

      {readiness.status === "loading" && !readiness.data && (
        <div className="border-b border-gray-100 bg-gray-50 px-4 py-3 text-[12px] text-gray-400">
          Loading readiness...
        </div>
      )}
      {readiness.status === "error" && (
        <div className="border-b border-amber-100 bg-amber-50 px-4 py-3 text-[12px] text-amber-700">
          Live readiness API unavailable. Demo rows are shown.
        </div>
      )}

      <div className="overflow-x-auto">
        <table className="w-full min-w-[900px] border-collapse text-[13px]">
          <thead>
            <tr className="bg-gray-50">
              {["Plan", "Source", "Product", "Source type", "State", "Reason"].map((heading) => (
                <th key={heading} className="border-b border-gray-200 p-4 text-left text-[11px] font-medium uppercase tracking-wide text-gray-400">
                  {heading}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={`${row.plan_display_id}-${row.source_display_id ?? "missing"}`} className="border-b border-gray-100 last:border-0 hover:bg-gray-50">
                <td className="p-4">
                  <div className="font-medium text-[#D50C2D]">{recordLabel(row.plan_display_id, "PLAN-")}</div>
                  <div className="mt-1 text-[11px] text-gray-500">{row.plan_code || row.plan_name}</div>
                </td>
                <td className="p-4">
                  <div className="font-medium text-gray-900">{sourceLabel(row)}</div>
                  <div className="mt-1 text-[11px] text-gray-400">{row.source_name || "No source linked"}</div>
                </td>
                <td className="p-4">
                  <SmallBadge>{row.product_type}</SmallBadge>
                </td>
                <td className="p-4">
                  <SmallBadge>{row.source_type || "-"}</SmallBadge>
                </td>
                <td className="p-4">
                  <ReadinessStateBadge state={row.state} />
                </td>
                <td className="p-4 text-[12px] text-gray-500">{row.reason}</td>
              </tr>
            ))}
            {rows.length === 0 && (
              <tr>
                <td colSpan={6} className="p-4 text-center text-[12px] text-gray-400">No provider readiness rows</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </section>
  );
}

function readinessStatusText(status: string, usingLive: boolean): string {
  if (status === "error") return "Demo readiness";
  if (status === "loading") return "Refreshing readiness";
  return usingLive ? "Live readiness" : "Demo readiness";
}

function statusTone(status: string, usingLive: boolean): string {
  if (status === "error") return "border-amber-200 bg-amber-50 text-amber-700";
  if (status === "loading") return "border-blue-200 bg-blue-50 text-blue-700";
  return usingLive ? "border-emerald-200 bg-emerald-50 text-emerald-700" : "border-gray-200 bg-white text-gray-500";
}

function sourceLabel(row: ProviderReadiness): string {
  return row.source_display_id ? recordLabel(row.source_display_id, "SRC-") : "-";
}

function ReadinessStateBadge({ state }: { state: ProviderReadinessState }) {
  return (
    <span className={`inline-flex items-center rounded-sm border px-1.5 py-px text-[11px] font-medium ${STATE_CLASS[state] ?? "border-gray-200 bg-gray-100 text-gray-500"}`}>
      {STATE_LABEL[state] ?? state}
    </span>
  );
}

function SmallBadge({ children }: { children: string }) {
  return (
    <span className="inline-flex items-center rounded-sm bg-gray-100 px-1.5 py-px text-[11px] font-medium text-gray-500">
      {children}
    </span>
  );
}

