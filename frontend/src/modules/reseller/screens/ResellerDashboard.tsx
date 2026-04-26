"use client";

import { KpiCard } from "@/components/ui/KpiCard";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { statusLabel } from "@/lib/api/displayLabels";
import { compactDateTime, moneyMinor, recordLabel, fmtMoney } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { RESELLER_CLIENTS } from "@/mocks/billingData";

interface DashboardClientRow {
  id: string;
  name: string;
  wallet: string;
  walletLow: boolean;
  services: string;
  orders: string;
  status: string;
  lastLogin: string;
}

export function ResellerDashboard() {
  const customers = useApiResource(
    () => billingApi.listResellerCustomers({ limit: 100 }),
    "reseller-dashboard-customers",
  );
  const services = useApiResource(
    () => billingApi.listResellerServices({ limit: 100 }),
    "reseller-dashboard-services",
  );
  const orders = useApiResource(
    () => billingApi.listResellerOrders({ limit: 100 }),
    "reseller-dashboard-orders",
  );
  const invoices = useApiResource(billingApi.listResellerInvoices, "reseller-dashboard-invoices");
  const wallets = useApiResource(
    () => billingApi.listResellerWallets({ limit: 100 }),
    "reseller-dashboard-wallets",
  );
  const wallet = wallets.data?.[0];
  const usingLiveRows = customers.status === "success";
  const orderCounts = new Map<number, number>();
  const serviceCounts = new Map<number, number>();
  const orderBuyerDisplayIDs = new Map<string, number>();
  const serviceOrderIDs = new Set((services.data ?? []).map((service) => service.order_id));
  const walletByOwnerDisplayID = new Map(
    (wallets.data ?? [])
      .filter((item) => item.owner_display_id)
      .map((item) => [item.owner_display_id!, item]),
  );

  if (orders.status === "success") {
    for (const order of orders.data ?? []) {
      if (!order.buyer_display_id) continue;
      orderBuyerDisplayIDs.set(order.id, order.buyer_display_id);
      orderCounts.set(order.buyer_display_id, (orderCounts.get(order.buyer_display_id) ?? 0) + 1);
    }
  }
  if (orders.status === "success" && services.status === "success") {
    for (const service of services.data ?? []) {
      const buyerDisplayID = service.buyer_display_id ?? orderBuyerDisplayIDs.get(service.order_id);
      if (buyerDisplayID) serviceCounts.set(buyerDisplayID, (serviceCounts.get(buyerDisplayID) ?? 0) + 1);
    }
  }

  const rows: DashboardClientRow[] = usingLiveRows
    ? (customers.data ?? []).slice(0, 5).map((client) => {
        const clientWallet = walletByOwnerDisplayID.get(client.display_id);
        return {
          id: recordLabel(client.display_id, "ACC-"),
          name: client.full_name || client.email,
          wallet: clientWallet ? moneyMinor(clientWallet.available_balance_minor, clientWallet.currency) : "-",
          walletLow: clientWallet ? clientWallet.available_balance_minor < 2000 : false,
          services: services.status === "success" && orders.status === "success"
            ? String(serviceCounts.get(client.display_id) ?? 0)
            : "-",
          orders: orders.status === "success" ? String(orderCounts.get(client.display_id) ?? 0) : "-",
          status: client.status,
          lastLogin: compactDateTime(client.last_login_at),
        };
      })
    : RESELLER_CLIENTS.slice(0, 5).map((client) => ({
        id: client.id,
        name: client.name,
        wallet: fmtMoney(client.wallet),
        walletLow: client.wallet < 20,
        services: String(client.services),
        orders: String(client.orders),
        status: client.status,
        lastLogin: client.lastLogin,
      }));
  const pendingFulfillment = (orders.data ?? []).filter((order) => (
    order.order_status === "paid" &&
    order.billing_status === "paid" &&
    !serviceOrderIDs.has(order.id)
  )).length;
  const revenueMinor = (invoices.data ?? []).reduce((total, invoice) => total + invoice.total_minor, 0);
  const source = customers.status === "error"
    ? "Live API unavailable. Showing demo dashboard data."
    : customers.status === "loading"
      ? "Refreshing live reseller dashboard..."
      : usingLiveRows
        ? "Live reseller dashboard"
        : "Demo reseller dashboard";

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="bg-white border border-gray-200 rounded p-4 flex items-center justify-between gap-4">
        <div>
          <div className="text-[11px] font-medium uppercase tracking-wide text-gray-400 mb-1">Reseller wallet</div>
          <div className="text-3xl font-medium tabular-nums text-gray-900">
            {wallet ? moneyMinor(wallet.available_balance_minor, wallet.currency) : "$4,820.50"}
          </div>
          <div className="text-[12px] text-gray-400 mt-1">
            {wallet ? `${recordLabel(wallet.display_id, "WAL-")} / ${statusLabel(wallet.status)}` : "Available balance / demo"}
          </div>
        </div>
        <span className="text-[11px] text-gray-400 text-right">{source}</span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <KpiCard label="Total clients" value={usingLiveRows ? String(customers.data?.length ?? 0) : "312"} sub="reseller scope" />
        <KpiCard
          label="Active services"
          value={services.status === "success" ? String((services.data ?? []).filter((item) => item.status === "active").length) : "1,840"}
          sub="across clients"
        />
        <KpiCard
          label="Revenue MTD"
          value={invoices.status === "success" ? moneyMinor(revenueMinor, invoices.data?.[0]?.currency ?? "USD") : "$12.4k"}
          sub="issued invoices"
        />
        <KpiCard
          label="Pending fulfillment"
          value={orders.status === "success" && services.status === "success" ? String(pendingFulfillment) : "-"}
          sub="paid orders"
        />
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Recent clients</h3>
          <span className="text-[12px] text-[#D50C2D]">View all</span>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-[760px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Name", "Wallet", "Services", "Orders", "Status", "Last login"].map((h) => (
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
                  <td className="p-4 tabular-nums">
                    <span className={client.walletLow ? "text-red-600 font-medium" : ""}>{client.wallet}</span>
                  </td>
                  <td className="p-4 tabular-nums text-right">{client.services}</td>
                  <td className="p-4 tabular-nums text-right">{client.orders}</td>
                  <td className="p-4"><StatusBadge status={client.status} dot /></td>
                  <td className="p-4 text-gray-400">{client.lastLogin}</td>
                </tr>
              ))}
              {usingLiveRows && rows.length === 0 && (
                <tr><td colSpan={7} className="p-4 text-center text-[12px] text-gray-400">No reseller clients</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
