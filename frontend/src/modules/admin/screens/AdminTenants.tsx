"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { TENANTS } from "@/mocks/billingData";

interface TenantRow {
  id: string;
  name: string;
  type: string;
  domain: string;
  users: string;
  currency: string;
  status: string;
  created: string;
}

function sourceText(status: string, usingLive: boolean) {
  if (status === "error") return "Live API unavailable. Showing demo account data.";
  if (status === "loading") return "Refreshing live accounts...";
  return usingLive ? "Live tenant accounts" : "Demo account data";
}

export function AdminTenants() {
  const tenants = useApiResource(
    () => billingApi.listAdminTenants({ limit: 100 }),
    "admin-tenants",
  );
  const usingLive = tenants.status === "success";
  const rows: TenantRow[] = usingLive
    ? (tenants.data ?? []).map((tenant) => ({
        id: recordLabel(tenant.display_id, "TEN-"),
        name: tenant.name,
        type: tenant.tenant_type,
        domain: tenant.primary_domain || tenant.slug,
        users: tenant.user_count.toLocaleString(),
        currency: tenant.default_currency,
        status: tenant.status,
        created: compactDateTime(tenant.created_at),
      }))
    : TENANTS.map((tenant) => ({
        id: tenant.id,
        name: tenant.name,
        type: tenant.type,
        domain: tenant.domain,
        users: tenant.clients.toLocaleString(),
        currency: "demo",
        status: tenant.status,
        created: tenant.since ?? "-",
      }));

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Accounts</h3>
          <div className="flex flex-wrap items-center justify-end gap-3">
            <span className="text-[11px] text-gray-400">{sourceText(tenants.status, usingLive)}</span>
            <span className="text-[11px] text-gray-400">{rows.length} accounts</span>
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-[760px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Name", "Type", "Domain", "Users", "Currency", "Status", "Created"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((tenant) => (
                <tr key={tenant.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D]">{tenant.id}</td>
                  <td className="p-4 font-medium text-gray-900">{tenant.name}</td>
                  <td className="p-4">
                    <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{tenant.type}</span>
                  </td>
                  <td className="p-4 text-[12px] text-gray-500">{tenant.domain}</td>
                  <td className="p-4 tabular-nums text-right">{tenant.users}</td>
                  <td className="p-4 text-[12px] text-gray-500">{tenant.currency}</td>
                  <td className="p-4"><StatusBadge status={tenant.status} dot /></td>
                  <td className="p-4 text-gray-400">{tenant.created}</td>
                </tr>
              ))}
              {usingLive && rows.length === 0 && (
                <tr><td colSpan={8} className="p-4 text-center text-[12px] text-gray-400">No tenant accounts</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
