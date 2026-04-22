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
    title: "Dashboard", breadcrumbs: ["ProxyVN", "Dashboard"],
    meta: <StatusBadge status="active" dot />,
    component: <ResellerDashboard />,
  },
  "reseller-clients": {
    title: "Clients", breadcrumbs: ["ProxyVN", "Clients"],
    meta: <span className="text-[11px] text-gray-400">312 clients</span>,
    component: <ResellerClients />,
  },
  "reseller-catalog": {
    title: "Catalog / Pricing", breadcrumbs: ["ProxyVN", "Catalog"],
    meta: <span className="text-[11px] text-amber-600 font-medium">1 margin warning</span>,
    component: <ResellerCatalog />,
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
  "reseller-orders": {
    title: "Orders", breadcrumbs: ["ProxyVN", "Orders"],
    component: <ResellerClients />,
  },
  "reseller-wallet": {
    title: "Wallet & Top-up", breadcrumbs: ["ProxyVN", "Finance", "Wallet"],
    meta: <StatusBadge status="active" dot />,
    component: <ResellerWallet />,
  },
  "reseller-reports": {
    title: "Reports", breadcrumbs: ["ProxyVN", "Reports"],
    component: <ResellerDashboard />,
  },
  "reseller-settings": {
    title: "Branding & Settings", breadcrumbs: ["ProxyVN", "Settings"],
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
