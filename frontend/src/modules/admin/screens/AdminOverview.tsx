"use client";

import type { ReactNode } from "react";
import { CreditCard, Headphones, Server, User, Wallet, XCircle } from "lucide-react";
import { KpiCard } from "@/components/ui/KpiCard";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { statusLabel } from "@/lib/api/displayLabels";
import { compactDateTime, moneyMinor, fmtMoney } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { adminDisplayLabel, mapAdminInvoiceView } from "@/lib/api/viewModels";
import { ACTIVITY_FEED, ActivityEvent, INVOICES, TOPUP_REQUESTS } from "@/mocks/billingData";

const ACTIVITY_ICONS: Record<ActivityEvent["icon"], ReactNode> = {
  payment: <CreditCard size={11} />,
  user: <User size={11} />,
  server: <Server size={11} />,
  ticket: <Headphones size={11} />,
  error: <XCircle size={11} />,
  wallet: <Wallet size={11} />,
};

const PENDING_TOPUP_STATUSES = new Set([
  "pending",
  "pending_verification",
  "submitted",
  "under_review",
  "manual_review",
  "queued",
]);

type ActivityRow = {
  t: string;
  icon: ActivityEvent["icon"];
  text: string;
  type: ActivityEvent["type"];
};

function isPendingTopup(status: string) {
  return PENDING_TOPUP_STATUSES.has(status);
}

function topupActivityTone(status: string): ActivityEvent["type"] {
  if (status === "approved" || status === "paid" || status === "posted") return "ok";
  if (status === "rejected" || status === "failed") return "danger";
  if (isPendingTopup(status)) return "warn";
  return "info";
}

function countByStatus(items: { status: string }[], statuses: string[]) {
  const statusSet = new Set(statuses);
  return items.filter((item) => statusSet.has(item.status)).length;
}

function formatWalletTotal(wallets: { available_balance_minor: number; currency: string }[]) {
  if (wallets.length === 0) return "$0.00";

  const totals = new Map<string, number>();
  for (const wallet of wallets) {
    totals.set(wallet.currency, (totals.get(wallet.currency) ?? 0) + wallet.available_balance_minor);
  }

  const [currency, total] = [...totals.entries()].sort((a, b) => Math.abs(b[1]) - Math.abs(a[1]))[0];
  const suffix = totals.size > 1 ? ` +${totals.size - 1}` : "";
  return `${moneyMinor(total, currency)}${suffix}`;
}

export function AdminOverview() {
  const wallets = useApiResource(
    () => billingApi.listAdminWallets({ limit: 50 }),
    "admin-overview:wallets",
  );
  const orders = useApiResource(
    () => billingApi.listAdminOrders({ limit: 50 }),
    "admin-overview:orders",
  );
  const services = useApiResource(
    () => billingApi.listAdminServices({ limit: 50 }),
    "admin-overview:services",
  );
  const topups = useApiResource(
    () => billingApi.listAdminTopupRequests({ limit: 50 }),
    "admin-overview:topups",
  );
  const invoices = useApiResource(
    () => billingApi.listAdminInvoices(),
    "admin-overview:invoices",
  );

  const liveWallets = wallets.status === "success" ? (wallets.data ?? []) : [];
  const liveOrders = orders.status === "success" ? (orders.data ?? []) : [];
  const liveServices = services.status === "success" ? (services.data ?? []) : [];
  const liveTopups = topups.status === "success" ? (topups.data ?? []) : [];
  const usingLiveInvoices = invoices.status === "success";

  const walletValue = wallets.status === "success" ? formatWalletTotal(liveWallets) : "$118.2k";
  const orderValue = orders.status === "success" ? liveOrders.length.toLocaleString() : "2,847";
  const activeServiceCount = services.status === "success"
    ? countByStatus(liveServices, ["active", "running"])
    : 15604;
  const pendingTopupCount = topups.status === "success"
    ? liveTopups.filter((req) => isPendingTopup(req.status)).length
    : TOPUP_REQUESTS.filter((req) => isPendingTopup(req.status)).length;

  const invoiceRows = usingLiveInvoices
    ? (invoices.data ?? []).slice(0, 7).map(mapAdminInvoiceView)
    : INVOICES.slice(0, 7).map((inv) => ({
        id: inv.id,
        customer: inv.customer,
        issued: inv.issued,
        due: inv.due,
        amount: fmtMoney(inv.amount),
        status: inv.status,
      }));

  const activityRows: ActivityRow[] = topups.status === "success"
    ? liveTopups.slice(0, 7).map((req) => ({
        t: compactDateTime(req.created_at),
        icon: "wallet",
        text: `Top-up ${adminDisplayLabel(req.display_id, "TUP-")} ${statusLabel(req.status)} (${moneyMinor(req.amount_minor, req.currency)})`,
        type: topupActivityTone(req.status),
      }))
    : ACTIVITY_FEED;

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
        <KpiCard
          label="Wallet balance"
          value={walletValue}
          sub={wallets.status === "success" ? `${liveWallets.length} live wallets` : "demo fallback"}
        />
        <KpiCard
          label="Orders"
          value={orderValue}
          sub={orders.status === "success" ? "live order read model" : "demo fallback"}
        />
        <KpiCard
          label="Active services"
          value={activeServiceCount.toLocaleString()}
          sub={services.status === "success" ? "live service read model" : "proxies + VPS"}
        />
        <KpiCard
          label="Pending top-ups"
          value={pendingTopupCount.toLocaleString()}
          sub={topups.status === "success" ? "live review queue" : "demo fallback"}
        />
      </div>

      <div className="grid grid-cols-1 gap-4 xl:grid-cols-[1.7fr_1fr]">
        <div className="bg-white border border-gray-200 rounded">
          <div className="p-4 border-b border-gray-100 flex items-center justify-between">
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Recent invoices</h3>
            <span className="text-[11px] text-gray-400">
              {usingLiveInvoices ? "Live API" : "Demo fallback"}
            </span>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-[760px] w-full text-[13px] border-collapse">
              <thead>
                <tr className="bg-gray-50">
                  {["Invoice", "Customer", "Issued", "Due", "Amount", "Status"].map((h) => (
                    <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {invoiceRows.map((inv) => (
                  <tr key={inv.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                    <td className="p-4 text-[12px] text-[#D50C2D]">{inv.id}</td>
                    <td className="p-4 text-gray-700">{inv.customer}</td>
                    <td className="p-4 text-gray-400">{inv.issued}</td>
                    <td className="p-4 text-gray-400">{inv.due}</td>
                    <td className="p-4 text-right font-medium tabular-nums">{inv.amount}</td>
                    <td className="p-4"><StatusBadge status={inv.status} dot /></td>
                  </tr>
                ))}
                {usingLiveInvoices && invoiceRows.length === 0 && (
                  <tr><td colSpan={6} className="p-4 text-center text-[12px] text-gray-400">No invoices</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>

        <div className="bg-white border border-gray-200 rounded">
          <div className="p-4 border-b border-gray-100 flex items-center justify-between">
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Activity feed</h3>
            <span className="text-[11px] text-gray-400">
              {topups.status === "success" ? "Live top-ups" : "Demo fallback"}
            </span>
          </div>
          <div>
            {activityRows.map((a, i) => (
              <div key={`${a.t}-${i}`} className="flex gap-4 p-4 border-b border-gray-100 last:border-0 items-start">
                <div
                  className={`w-[22px] h-[22px] rounded-full grid place-items-center text-[11px] shrink-0 ${
                    a.type === "ok" ? "bg-green-50 text-green-700"
                      : a.type === "warn" ? "bg-amber-50 text-amber-700"
                        : a.type === "danger" ? "bg-red-50 text-red-700"
                          : "bg-gray-100 text-gray-500"
                  }`}
                >
                  {ACTIVITY_ICONS[a.icon]}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="text-[12px] text-gray-700 leading-snug">{a.text}</div>
                  <div className="text-[11px] text-gray-400 mt-0.5">{a.t}</div>
                </div>
              </div>
            ))}
            {topups.status === "success" && activityRows.length === 0 && (
              <div className="p-4 text-center text-[12px] text-gray-400">No recent top-up activity</div>
            )}
          </div>
        </div>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Infrastructure health</h3>
          <StatusBadge status="active" dot />
        </div>
        <div>
          {[
            { label: "Proxy network uptime", value: "99.98%", bar: 0.9998 },
            { label: "VPS fleet uptime", value: "99.94%", bar: 0.9994 },
            { label: "Payment gateway", value: "100%", bar: 1 },
            { label: "API p95 latency", value: "142ms", bar: 0.82 },
            { label: "Support first response", value: "8m avg", bar: 0.72 },
          ].map((r) => (
            <div key={r.label} className="p-4 border-b border-gray-100 last:border-0">
              <div className="flex items-center justify-between mb-1.5">
                <span className="text-[12px] text-gray-700">{r.label}</span>
                <span className="text-[12px] font-medium tabular-nums">{r.value}</span>
              </div>
              <div className="h-1 bg-gray-100 rounded">
                <div
                  className={`h-full rounded ${r.bar > 0.99 ? "bg-green-500" : r.bar > 0.8 ? "bg-blue-500" : "bg-amber-400"}`}
                  style={{ width: `${r.bar * 100}%` }}
                />
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
