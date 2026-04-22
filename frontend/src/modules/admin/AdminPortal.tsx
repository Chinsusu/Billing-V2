"use client";

import { useState } from "react";
import { AppShell } from "@/components/layout/AppShell";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { AdminOverview } from "./screens/AdminOverview";
import { AdminTenants } from "./screens/AdminTenants";
import { AdminProvisioning } from "./screens/AdminProvisioning";
import { AdminTopups } from "./screens/AdminTopups";
import { AdminProviders } from "./screens/AdminProviders";
import { AdminCustomers } from "./screens/AdminCustomers";
import { AdminInvoices } from "./screens/AdminInvoices";
import { AdminTransactions } from "./screens/AdminTransactions";
import { AdminProducts } from "./screens/AdminProducts";
import { AdminTickets } from "./screens/AdminTickets";
import { AdminSettings } from "./screens/AdminSettings";
import { AdminAlerts } from "./screens/AdminAlerts";
import { AdminLogs } from "./screens/AdminLogs";
import { AdminServicesProxies } from "./screens/AdminServicesProxies";
import { AdminServicesVPS } from "./screens/AdminServicesVPS";
import { AdminServicesBandwidth } from "./screens/AdminServicesBandwidth";

interface ScreenConfig {
  title: string;
  breadcrumbs: string[];
  meta?: React.ReactNode;
  component: React.ReactNode;
}

const SCREENS: Record<string, ScreenConfig> = {
  "admin-overview": {
    title: "Overview", breadcrumbs: ["HANetwork", "Overview"],
    meta: <StatusBadge status="active" dot />,
    component: <AdminOverview />,
  },
  "admin-tenants": {
    title: "Accounts", breadcrumbs: ["HANetwork", "Platform", "Accounts"],
    meta: <span className="text-[11px] text-gray-400">2,852 accounts · 4 resellers</span>,
    component: <AdminTenants />,
  },
  "admin-provisioning": {
    title: "Provisioning queue", breadcrumbs: ["HANetwork", "Platform", "Provisioning"],
    meta: <StatusBadge status="manual_review" dot />,
    component: <AdminProvisioning />,
  },
  "admin-topups": {
    title: "Top-up verification", breadcrumbs: ["HANetwork", "Platform", "Top-ups"],
    meta: <span className="text-[11px] text-amber-600 font-medium">3 pending</span>,
    component: <AdminTopups />,
  },
  "admin-providers": {
    title: "Providers / Sources", breadcrumbs: ["HANetwork", "Platform", "Providers"],
    meta: <span className="text-[11px] text-amber-600 font-medium">1 degraded</span>,
    component: <AdminProviders />,
  },
  "admin-customers": {
    title: "Customers", breadcrumbs: ["HANetwork", "Customers"],
    component: <AdminCustomers />,
  },
  "admin-tickets": {
    title: "Support tickets", breadcrumbs: ["HANetwork", "Support", "Tickets"],
    component: <AdminTickets />,
  },
  "admin-services-proxies": {
    title: "Proxies", breadcrumbs: ["HANetwork", "Services", "Proxies"],
    component: <AdminServicesProxies />,
  },
  "admin-services-vps": {
    title: "VPS", breadcrumbs: ["HANetwork", "Services", "VPS"],
    component: <AdminServicesVPS />,
  },
  "admin-services-bandwidth": {
    title: "Bandwidth", breadcrumbs: ["HANetwork", "Services", "Bandwidth"],
    component: <AdminServicesBandwidth />,
  },
  "admin-invoices": {
    title: "Invoices", breadcrumbs: ["HANetwork", "Billing", "Invoices"],
    component: <AdminInvoices />,
  },
  "admin-transactions": {
    title: "Transactions", breadcrumbs: ["HANetwork", "Billing", "Transactions"],
    component: <AdminTransactions />,
  },
  "admin-products": {
    title: "Products & Pricing", breadcrumbs: ["HANetwork", "Billing", "Products"],
    component: <AdminProducts />,
  },
  "admin-reports": {
    title: "Reports", breadcrumbs: ["HANetwork", "Billing", "Reports"],
    component: <AdminOverview />,
  },
  "admin-alerts": {
    title: "Alerts", breadcrumbs: ["HANetwork", "System", "Alerts"],
    meta: <span className="text-[11px] text-red-600 font-medium">5 open</span>,
    component: <AdminAlerts />,
  },
  "admin-logs": {
    title: "Audit logs", breadcrumbs: ["HANetwork", "System", "Logs"],
    component: <AdminLogs />,
  },
  "admin-settings": {
    title: "Settings", breadcrumbs: ["HANetwork", "Settings"],
    component: <AdminSettings />,
  },
};

export function AdminPortal() {
  const [screen, setScreen] = useState("admin-overview");
  const cur = SCREENS[screen] ?? SCREENS["admin-overview"];

  return (
    <AppShell
      portal="admin"
      activeScreen={screen}
      onSelectScreen={setScreen}
      title={cur.title}
      breadcrumbs={cur.breadcrumbs}
      meta={cur.meta}
      actions={
        <button className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-white hover:bg-gray-50 text-gray-700 border border-gray-300 rounded-md cursor-pointer transition-colors shadow-sm">
          + New
        </button>
      }
    >
      {cur.component}
    </AppShell>
  );
}
