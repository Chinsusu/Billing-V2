import type { AdminServiceDemoRow } from "../components/AdminServiceInventoryTable";
import { AdminServiceInventoryTable } from "../components/AdminServiceInventoryTable";
import { BANDWIDTH_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

function renewsInText(days: number) {
  return days < 0 ? `${Math.abs(days)}d overdue` : `${days}d left`;
}

const DEMO_ROWS: AdminServiceDemoRow[] = BANDWIDTH_SERVICES.map((service) => ({
  id: service.id,
  service: service.label,
  owner: service.customer,
  tenant: service.tenant,
  resource: `${service.usedGB} / ${service.totalGB} GB`,
  plan: "bandwidth",
  region: service.region,
  status: service.status,
  billingStatus: service.renewsIn < 0 ? "overdue" : "paid",
  created: "Demo data",
  expires: renewsInText(service.renewsIn),
  provider: "demo",
  note: `${service.usedPct}% used / ${fmtMoney(service.price)} monthly`,
}));

export function AdminServicesBandwidth() {
  return (
    <AdminServiceInventoryTable
      family="bandwidth"
      title="Bandwidth"
      demoRows={DEMO_ROWS}
    />
  );
}
