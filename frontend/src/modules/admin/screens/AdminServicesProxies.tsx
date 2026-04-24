import type { AdminServiceDemoRow } from "../components/AdminServiceInventoryTable";
import { AdminServiceInventoryTable } from "../components/AdminServiceInventoryTable";
import { PROXY_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

function renewsInText(days: number) {
  return days < 0 ? `${Math.abs(days)}d overdue` : `${days}d left`;
}

const DEMO_ROWS: AdminServiceDemoRow[] = PROXY_SERVICES.map((service) => ({
  id: service.id,
  service: service.label,
  owner: service.customer,
  tenant: service.tenant,
  resource: service.ipCount > 0 ? `${service.ipCount} IPs` : service.protocol,
  plan: `proxy-${service.proxyType}`,
  region: service.region,
  status: service.status,
  billingStatus: service.renewsIn < 0 ? "overdue" : "paid",
  created: "Demo data",
  expires: renewsInText(service.renewsIn),
  provider: "demo",
  note: `${service.protocol} / ${fmtMoney(service.price)} monthly`,
}));

export function AdminServicesProxies() {
  return (
    <AdminServiceInventoryTable
      family="proxy"
      title="Proxies"
      demoRows={DEMO_ROWS}
    />
  );
}
