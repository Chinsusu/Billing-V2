import type { AdminServiceDemoRow } from "../components/AdminServiceInventoryTable";
import { AdminServiceInventoryTable } from "../components/AdminServiceInventoryTable";
import { VPS_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/lib/api/format";

function renewsInText(days: number) {
  return days < 0 ? `${Math.abs(days)}d overdue` : `${days}d left`;
}

const DEMO_ROWS: AdminServiceDemoRow[] = VPS_SERVICES.map((service) => ({
  id: service.id,
  service: service.label,
  owner: service.customer,
  tenant: service.tenant,
  resource: service.ip,
  plan: `${service.cpu}C / ${service.ram}G / ${service.disk}G`,
  region: service.region,
  status: service.status,
  billingStatus: service.renewsIn < 0 ? "overdue" : "paid",
  created: "Demo data",
  expires: renewsInText(service.renewsIn),
  provider: service.provider,
  note: `${service.os} / ${fmtMoney(service.price)} monthly`,
}));

export function AdminServicesVPS() {
  return (
    <AdminServiceInventoryTable
      family="vps"
      title="VPS"
      demoRows={DEMO_ROWS}
    />
  );
}
