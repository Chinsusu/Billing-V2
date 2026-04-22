"use client";

import { useState } from "react";
import { AppShell } from "@/components/layout/AppShell";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { ResellerDashboard } from "./screens/ResellerDashboard";
import { ResellerClients } from "./screens/ResellerClients";
import { ResellerCatalog } from "./screens/ResellerCatalog";
import { ResellerWallet } from "./screens/ResellerWallet";
import { ResellerSettings } from "./screens/ResellerSettings";

interface ScreenConfig {
  title: string;
  breadcrumbs: string[];
  meta?: React.ReactNode;
  component: React.ReactNode;
}

const SCREENS: Record<string, ScreenConfig> = {
  "reseller-overview": {
    title: "Overview", breadcrumbs: ["ProxyVN", "Overview"],
    meta: <StatusBadge status="active" dot />,
    component: <ResellerDashboard />,
  },
  "reseller-accounts": {
    title: "Accounts", breadcrumbs: ["ProxyVN", "Customers", "Accounts"],
    meta: <span className="text-[11px] text-gray-400">312 accounts</span>,
    component: <ResellerClients />,
  },
  "reseller-tickets": {
    title: "Support tickets", breadcrumbs: ["ProxyVN", "Customers", "Tickets"],
    component: <ResellerClients />,
  },
  "reseller-services-proxies": {
    title: "Proxies", breadcrumbs: ["ProxyVN", "Services", "Proxies"],
    component: <ResellerClients />,
  },
  "reseller-services-vps": {
    title: "VPS", breadcrumbs: ["ProxyVN", "Services", "VPS"],
    component: <ResellerClients />,
  },
  "reseller-services-bandwidth": {
    title: "Bandwidth", breadcrumbs: ["ProxyVN", "Services", "Bandwidth"],
    component: <ResellerClients />,
  },
  "reseller-invoices": {
    title: "Invoices", breadcrumbs: ["ProxyVN", "Billing", "Invoices"],
    component: <ResellerWallet />,
  },
  "reseller-transactions": {
    title: "Transactions", breadcrumbs: ["ProxyVN", "Billing", "Transactions"],
    component: <ResellerWallet />,
  },
  "reseller-products": {
    title: "Products & Pricing", breadcrumbs: ["ProxyVN", "Billing", "Products"],
    component: <ResellerCatalog />,
  },
  "reseller-reports": {
    title: "Reports", breadcrumbs: ["ProxyVN", "Billing", "Reports"],
    component: <ResellerDashboard />,
  },
  "reseller-settings": {
    title: "Settings", breadcrumbs: ["ProxyVN", "Settings"],
    component: <ResellerSettings />,
  },
};

export function ResellerPortal() {
  const [screen, setScreen] = useState("reseller-overview");
  const cur = SCREENS[screen] ?? SCREENS["reseller-overview"];

  return (
    <AppShell
      portal="reseller"
      activeScreen={screen}
      onSelectScreen={setScreen}
      title={cur.title}
      breadcrumbs={cur.breadcrumbs}
      meta={cur.meta}
      actions={
        <button className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm">
          + Top up
        </button>
      }
    >
      {cur.component}
    </AppShell>
  );
}
