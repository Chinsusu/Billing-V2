"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import type { Order, ServiceInstance } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { BANDWIDTH_SERVICES, PROXY_SERVICES, VPS_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

type ResellerServiceCategory = "proxies" | "vps" | "bandwidth";

interface ResellerServicesProps {
  category: ResellerServiceCategory;
}

const CONFIG: Record<ResellerServiceCategory, { title: string; empty: string }> = {
  proxies: { title: "Proxy services", empty: "No proxy services" },
  vps: { title: "VPS services", empty: "No VPS services" },
  bandwidth: { title: "Bandwidth usage", empty: "No bandwidth records" },
};

export function ResellerServices({ category }: ResellerServicesProps) {
  const services = useApiResource(
    () => billingApi.listResellerServices({ limit: 100 }),
    "reseller-services",
  );
  const orders = useApiResource(
    () => billingApi.listResellerOrders({ limit: 100 }),
    "reseller-service-orders",
  );
  const customers = useApiResource(
    () => billingApi.listResellerCustomers({ limit: 100 }),
    "reseller-service-customers",
  );
  const usingLive = services.status === "success";
  const rows = usingLive
    ? liveServiceRows(category, services.data ?? [], orders.data ?? [], customers.data ?? [])
    : demoServiceRows(category);
  const attention = rows.filter((row) => row.status === "suspended" || row.status === "overdue").length;
  const revenue = rows.reduce((total, row) => total + row.priceMinor, 0);
  const config = CONFIG[category];
  const source = services.status === "error"
    ? "Live API unavailable. Showing demo service data."
    : services.status === "loading"
      ? "Refreshing live services..."
      : usingLive
        ? "Live reseller services"
        : "Demo service data";

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <SummaryTile label="Records" value={String(rows.length)} />
        <SummaryTile label="Monthly value" value={usingLive ? moneyMinor(revenue) : fmtMoney(revenue / 100)} />
        <SummaryTile label="Attention" value={String(attention)} tone={attention > 0 ? "warn" : "neutral"} />
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-3">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">{config.title}</h3>
          <span className="text-[11px] text-gray-400">{source}</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[760px]">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Service", "Client", "Region", "Usage", "Renewal", "Price", "Status"].map((heading) => (
                  <th key={heading} className="text-left text-[11px] font-medium uppercase text-gray-400 p-4 border-b border-gray-200">
                    {heading}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D] font-medium">{row.id}</td>
                  <td className="p-4 font-medium text-gray-900">{row.label}</td>
                  <td className="p-4 text-gray-500">{row.customer}</td>
                  <td className="p-4 text-gray-500">{row.region}</td>
                  <td className="p-4 text-gray-500">{row.usage}</td>
                  <td className="p-4 text-gray-500">{row.renewal}</td>
                  <td className="p-4 text-right font-medium tabular-nums">{row.price}</td>
                  <td className="p-4"><StatusBadge status={row.status} dot /></td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr><td colSpan={8} className="p-4 text-center text-[12px] text-gray-400">{config.empty}</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

interface ServiceRow {
  id: string;
  label: string;
  customer: string;
  region: string;
  usage: string;
  renewal: string;
  price: string;
  priceMinor: number;
  status: string;
}

function demoServiceRows(category: ResellerServiceCategory): ServiceRow[] {
  if (category === "vps") {
    return VPS_SERVICES.filter((item) => item.tenant === "ProxyVN").map((item) => ({
      id: item.id,
      label: item.label,
      customer: item.customer,
      region: item.region,
      usage: `${item.cpu}C / ${item.ram}GB / ${item.disk}GB`,
      renewal: item.renewsIn < 0 ? `${Math.abs(item.renewsIn)}d overdue` : `${item.renewsIn}d`,
      price: fmtMoney(item.price),
      priceMinor: Math.round(item.price * 100),
      status: item.status,
    }));
  }
  if (category === "bandwidth") {
    return BANDWIDTH_SERVICES.filter((item) => item.tenant === "ProxyVN").map((item) => ({
      id: item.id,
      label: item.label,
      customer: item.customer,
      region: item.region,
      usage: `${item.usedGB} / ${item.totalGB} GB`,
      renewal: item.renewsIn < 0 ? `${Math.abs(item.renewsIn)}d overdue` : `${item.renewsIn}d`,
      price: fmtMoney(item.price),
      priceMinor: Math.round(item.price * 100),
      status: item.status,
    }));
  }
  return PROXY_SERVICES.filter((item) => item.tenant === "ProxyVN").map((item) => ({
    id: item.id,
    label: item.label,
    customer: item.customer,
    region: item.region,
    usage: item.ipCount > 0 ? `${item.ipCount} IPs` : `${item.usedGB} / ${item.totalGB} GB`,
    renewal: item.renewsIn < 0 ? `${Math.abs(item.renewsIn)}d overdue` : `${item.renewsIn}d`,
    price: fmtMoney(item.price),
    priceMinor: Math.round(item.price * 100),
    status: item.status,
  }));
}

function liveServiceRows(
  category: ResellerServiceCategory,
  services: ServiceInstance[],
  orders: Order[],
  customers: { id: string; display_id: number; full_name: string; email: string }[],
): ServiceRow[] {
  const ordersByID = new Map(orders.map((order) => [order.id, order]));
  const customerByID = new Map(customers.map((customer) => [customer.id, customer]));
  return services
    .map((service) => {
      const order = ordersByID.get(service.order_id);
      const text = `${snapshotText(order?.product_snapshot)} ${snapshotText(order?.plan_snapshot)} ${service.external_resource_id}`.toLowerCase();
      const detected = text.includes("vps")
        ? "vps"
        : text.includes("bandwidth")
          ? "bandwidth"
          : "proxies";
      const customer = order?.buyer_user_id ? customerByID.get(order.buyer_user_id) : undefined;
      return {
        id: recordLabel(service.display_id, "SVC-"),
        label: service.external_resource_id || snapshotText(order?.plan_snapshot) || recordLabel(service.display_id, "SVC-"),
        customer: customer
          ? `${customer.full_name || customer.email} (${recordLabel(customer.display_id, "ACC-")})`
          : "-",
        region: snapshotText(order?.product_snapshot, ["location", "region"]) || "-",
        usage: order ? `${recordLabel(order.display_id, "ORD-")} x ${order.quantity}` : "-",
        renewal: compactDateTime(service.term_end),
        price: order ? moneyMinor(order.total_minor, order.currency) : "-",
        priceMinor: order?.total_minor ?? 0,
        status: service.status || service.billing_status,
        category: detected,
      };
    })
    .filter((row) => row.category === category || (category === "proxies" && row.category !== "vps" && row.category !== "bandwidth"));
}

function snapshotText(value: unknown, keys = ["name", "plan_code", "product_type", "description"]): string {
  if (!value || typeof value !== "object" || Array.isArray(value)) return "";
  const record = value as Record<string, unknown>;
  for (const key of keys) {
    const data = record[key];
    if (typeof data === "string" && data.trim()) return data;
  }
  return "";
}

function SummaryTile({ label, value, tone = "neutral" }: { label: string; value: string; tone?: "neutral" | "warn" }) {
  return (
    <div className={`bg-white border rounded p-4 ${tone === "warn" ? "border-amber-200" : "border-gray-200"}`}>
      <div className="text-[11px] text-gray-400 uppercase mb-1">{label}</div>
      <div className={`text-lg font-medium tabular-nums ${tone === "warn" ? "text-amber-700" : "text-gray-900"}`}>{value}</div>
    </div>
  );
}
