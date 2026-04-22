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
  "reseller-services": {
    title: "Services", breadcrumbs: ["ProxyVN", "Services"],
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
        <button className="h-7 p-4 text-[12px] font-medium bg-[#D50C2D] text-white rounded-[3px] border-0 hover:bg-[#B3082A] cursor-pointer">
          + Top up
        </button>
      }
    >
      {cur.component}
    </AppShell>
  );
}
