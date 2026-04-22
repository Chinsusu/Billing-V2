"use client";

import { useState } from "react";
import { AppShell } from "@/components/layout/AppShell";
import { ClientDashboard } from "./screens/ClientDashboard";
import { ClientShop } from "./screens/ClientShop";
import { ClientWallet } from "./screens/ClientWallet";
import { ClientPlaceholder } from "./screens/ClientPlaceholder";

interface ScreenConfig {
  title: string;
  breadcrumbs: string[];
  meta?: React.ReactNode;
  component: React.ReactNode;
}

const SCREENS: Record<string, ScreenConfig> = {
  "client-overview": {
    title: "Dashboard", breadcrumbs: ["ProxyVN", "Dashboard"],
    meta: <span className="text-[11px] text-gray-400">Wallet $128.40</span>,
    component: <ClientDashboard />,
  },
  "client-shop": {
    title: "Shop", breadcrumbs: ["ProxyVN", "Shop"],
    component: <ClientShop />,
  },
  "client-services-proxies": {
    title: "Proxies", breadcrumbs: ["ProxyVN", "Services", "Proxies"],
    component: <ClientDashboard />,
  },
  "client-services-vps": {
    title: "VPS", breadcrumbs: ["ProxyVN", "Services", "VPS"],
    component: <ClientDashboard />,
  },
  "client-services-bandwidth": {
    title: "Bandwidth", breadcrumbs: ["ProxyVN", "Services", "Bandwidth"],
    component: <ClientDashboard />,
  },
  "client-wallet": {
    title: "Wallet", breadcrumbs: ["ProxyVN", "Billing", "Wallet"],
    meta: <span className="text-[11px] text-gray-400">$128.40</span>,
    component: <ClientWallet />,
  },
  "client-usage": {
    title: "Usage", breadcrumbs: ["ProxyVN", "Billing", "Usage"],
    component: <ClientPlaceholder title="Usage" />,
  },
  "client-settings": {
    title: "Settings", breadcrumbs: ["ProxyVN", "Account", "Settings"],
    component: <ClientPlaceholder title="Settings" />,
  },
  "client-support": {
    title: "Support", breadcrumbs: ["ProxyVN", "Support"],
    component: <ClientPlaceholder title="Support" />,
  },
};

export function ClientPortal() {
  const [screen, setScreen] = useState("client-overview");
  const cur = SCREENS[screen] ?? SCREENS["client-overview"];

  return (
    <AppShell
      portal="client"
      activeScreen={screen}
      onSelectScreen={setScreen}
      title={cur.title}
      breadcrumbs={cur.breadcrumbs}
      meta={cur.meta}
    >
      {cur.component}
    </AppShell>
  );
}
