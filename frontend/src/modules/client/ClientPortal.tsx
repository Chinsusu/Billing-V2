"use client";

import { useState } from "react";
import { AppShell } from "@/components/layout/AppShell";
import { ClientDashboard } from "./screens/ClientDashboard";
import { ClientInvoices } from "./screens/ClientInvoices";
import { ClientPlaceholder } from "./screens/ClientPlaceholder";
import { ClientServices } from "./screens/ClientServices";
import { ClientSettings } from "./screens/ClientSettings";
import { ClientTransactions } from "./screens/ClientTransactions";

interface ScreenConfig {
  title: string;
  breadcrumbs: string[];
  meta?: React.ReactNode;
  component: React.ReactNode;
}

const SCREENS: Record<string, ScreenConfig> = {
  "client-overview": {
    title: "Overview", breadcrumbs: ["ProxyVN", "Overview"],
    meta: <span className="text-[11px] text-gray-400">Wallet $128.40</span>,
    component: <ClientDashboard />,
  },
  "client-tickets": {
    title: "Support tickets", breadcrumbs: ["ProxyVN", "Customers", "Tickets"],
    component: <ClientPlaceholder title="Support tickets" />,
  },
  "client-services-proxies": {
    title: "Proxies", breadcrumbs: ["ProxyVN", "Services", "Proxies"],
    component: <ClientServices category="proxies" />,
  },
  "client-services-vps": {
    title: "VPS", breadcrumbs: ["ProxyVN", "Services", "VPS"],
    component: <ClientServices category="vps" />,
  },
  "client-services-bandwidth": {
    title: "Bandwidth", breadcrumbs: ["ProxyVN", "Services", "Bandwidth"],
    component: <ClientServices category="bandwidth" />,
  },
  "client-invoices": {
    title: "Invoices", breadcrumbs: ["ProxyVN", "Billing", "Invoices"],
    component: <ClientInvoices />,
  },
  "client-transactions": {
    title: "Transactions", breadcrumbs: ["ProxyVN", "Billing", "Transactions"],
    meta: <span className="text-[11px] text-gray-400">$128.40</span>,
    component: <ClientTransactions />,
  },
  "client-settings": {
    title: "Settings", breadcrumbs: ["ProxyVN", "Settings"],
    component: <ClientSettings />,
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
