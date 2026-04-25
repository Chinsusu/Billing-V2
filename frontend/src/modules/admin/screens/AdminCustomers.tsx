"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { useApiResource } from "@/lib/api/useApiResource";
import { mapAdminAccountView } from "@/lib/api/viewModels";
import { CUSTOMERS } from "@/mocks/billingData";
import { fmtMoneyShort } from "@/mocks/sampleData";

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

function sourceText(status: string, usingLive: boolean) {
  if (status === "error") return "Live API unavailable. Showing demo customer data.";
  if (status === "loading") return "Refreshing live customers...";
  return usingLive ? "Live customer accounts" : "Demo customer data";
}

export function AdminCustomers() {
  const customers = useApiResource(
    () => billingApi.listAdminCustomers({ limit: 100 }),
    "admin-customers",
  );
  const usingLive = customers.status === "success";
  const rows: CustomerRow[] = usingLive
    ? (customers.data ?? []).map(mapAdminAccountView)
    : CUSTOMERS.map((customer) => ({
        id: customer.id,
        name: customer.name,
        email: customer.email,
        type: customer.plan,
        tenant: customer.country,
        security: `${customer.services} services`,
        status: customer.status,
        created: customer.since,
        lastLogin: fmtMoneyShort(customer.mrr),
      }));

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
