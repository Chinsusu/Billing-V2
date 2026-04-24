"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { RESELLER_CLIENTS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

interface ClientRow {
  id: string;
  name: string;
  email: string;
  wallet: string;
  services: string;
  orders: string;
  status: string;
  lastLogin: string;
  walletLow: boolean;
}

function sourceText(status: string, usingLive: boolean, extraError?: string | null) {
  if (status === "error") return "Live API unavailable. Showing demo client data.";
  if (status === "loading") return "Refreshing live clients...";
  if (usingLive && extraError) return "Live clients loaded. Some counters are unavailable.";
  return usingLive ? "Live reseller clients" : "Demo client data";
}

export function ResellerClients() {
  const customers = useApiResource(
    () => billingApi.listResellerCustomers({ limit: 100 }),
    "reseller-customers",
  );
  const wallets = useApiResource(
    () => billingApi.listResellerWallets({ limit: 100 }),
    "reseller-client-wallets",
  );
  const orders = useApiResource(
    () => billingApi.listResellerOrders({ limit: 100 }),
    "reseller-client-orders",
  );
  const services = useApiResource(
    () => billingApi.listResellerServices({ limit: 100 }),
    "reseller-client-services",
  );
  const usingLive = customers.status === "success";
  const orderCounts = new Map<string, number>();
  const serviceCounts = new Map<string, number>();
  const orderBuyers = new Map<string, string>();
  const walletByOwner = new Map(
    (wallets.data ?? []).map((wallet) => [wallet.owner_id, wallet]),
  );

  if (orders.status === "success") {
    for (const order of orders.data ?? []) {
      if (!order.buyer_user_id) continue;
      orderBuyers.set(order.id, order.buyer_user_id);
      orderCounts.set(order.buyer_user_id, (orderCounts.get(order.buyer_user_id) ?? 0) + 1);
    }
  }
  if (orders.status === "success" && services.status === "success") {
    for (const service of services.data ?? []) {
      const buyerID = orderBuyers.get(service.order_id);
      if (buyerID) serviceCounts.set(buyerID, (serviceCounts.get(buyerID) ?? 0) + 1);
    }
  }

  const rows: ClientRow[] = usingLive
    ? (customers.data ?? []).map((client) => {
        const wallet = walletByOwner.get(client.id);
        return {
          id: recordLabel(client.display_id, "ACC-"),
          name: client.full_name || client.email,
          email: client.email,
          wallet: wallet ? moneyMinor(wallet.available_balance_minor, wallet.currency) : "-",
          services: services.status === "success" && orders.status === "success"
            ? String(serviceCounts.get(client.id) ?? 0)
            : "-",
          orders: orders.status === "success" ? String(orderCounts.get(client.id) ?? 0) : "-",
          status: client.status,
          lastLogin: compactDateTime(client.last_login_at),
          walletLow: wallet ? wallet.available_balance_minor < 2000 : false,
        };
      })
    : RESELLER_CLIENTS.map((client) => ({
        id: client.id,
        name: client.name,
        email: client.email,
        wallet: fmtMoney(client.wallet),
        services: String(client.services),
        orders: String(client.orders),
        status: client.status,
        lastLogin: client.lastLogin,
        walletLow: client.wallet < 20,
      }));
  const extraError = wallets.error ?? orders.error ?? services.error;

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Clients</h3>
          <div className="flex flex-wrap items-center justify-end gap-3">
            <span className="text-[11px] text-gray-400">{sourceText(customers.status, usingLive, extraError)}</span>
            <span className="text-[11px] text-gray-400">{rows.length} total</span>
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-[860px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Name", "Email", "Wallet", "Services", "Orders", "Status", "Last login"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((client) => (
                <tr key={client.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D]">{client.id}</td>
                  <td className="p-4 font-medium text-gray-900">{client.name}</td>
                  <td className="p-4 text-gray-400 text-[12px]">{client.email}</td>
                  <td className="p-4 tabular-nums">
                    <span className={client.walletLow ? "text-red-600 font-medium" : ""}>{client.wallet}</span>
                  </td>
                  <td className="p-4 tabular-nums text-right">{client.services}</td>
                  <td className="p-4 tabular-nums text-right">{client.orders}</td>
                  <td className="p-4"><StatusBadge status={client.status} dot /></td>
                  <td className="p-4 text-gray-400">{client.lastLogin}</td>
                </tr>
              ))}
              {usingLive && rows.length === 0 && (
                <tr><td colSpan={8} className="p-4 text-center text-[12px] text-gray-400">No reseller clients</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
