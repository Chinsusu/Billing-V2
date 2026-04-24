"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";

type ServiceCategory = "proxies" | "vps" | "bandwidth";

interface ClientServicesProps {
  category: ServiceCategory;
}

const CATEGORY_CONFIG: Record<ServiceCategory, { title: string; metric: string; empty: string }> = {
  proxies: {
    title: "Proxy services",
    metric: "Proxy records",
    empty: "No proxy services",
  },
  vps: {
    title: "VPS services",
    metric: "VPS records",
    empty: "No VPS services",
  },
  bandwidth: {
    title: "Bandwidth usage",
    metric: "Usage records",
    empty: "No bandwidth records",
  },
};

export function ClientServices({ category }: ClientServicesProps) {
  const services = useApiResource(billingApi.listClientServices);
  const liveRows = services.status === "success"
    ? services.data?.map((service) => ({
        id: service.id,
        category: inferCategory([service.external_resource_id, service.tenant_plan_id, service.provider_source_id]),
        label: recordLabel(service.display_id, "SVC-"),
        identifier: service.external_resource_id || "-",
        region: service.provider_source_id ?? "-",
        detail: service.order_id ? recordLabel(service.order_id.slice(-6), "ORD-") : "-",
        expiry: compactDateTime(service.term_end),
        status: service.status,
      })) ?? []
    : null;
  const sourceRows = liveRows ?? [];
  const rows = sourceRows.filter((service) => service.category === category);
  const suspended = rows.filter((service) => service.status === "suspended" || service.status === "overdue").length;
  const config = CATEGORY_CONFIG[category];

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <SummaryTile label={config.metric} value={String(rows.length)} />
        <SummaryTile label="Attention" value={String(suspended)} tone={suspended > 0 ? "warn" : "neutral"} />
        <SummaryTile label="Data source" value={liveRows ? "Live API" : "API pending"} />
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-3">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">{config.title}</h3>
          <span className="text-[11px] text-gray-400">{services.status === "loading" ? "Loading" : `${rows.length} records`}</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[760px]">
            <thead>
              <tr className="bg-gray-50">
                {["Service", "Identifier", "Region", category === "bandwidth" ? "Usage" : "Cycle / Order", "Expires", "Status"].map((heading) => (
                  <th key={heading} className="text-left text-[11px] font-medium uppercase text-gray-400 p-4 border-b border-gray-200">
                    {heading}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((service) => (
                <tr key={service.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 font-medium text-gray-900">{service.label}</td>
                  <td className="p-4 text-[12px] text-gray-500">{service.identifier}</td>
                  <td className="p-4 text-gray-500">{service.region}</td>
                  <td className="p-4 text-gray-500">{service.detail}</td>
                  <td className="p-4 text-gray-500">{service.expiry}</td>
                  <td className="p-4"><StatusBadge status={service.status} dot /></td>
                </tr>
              ))}
              {services.status === "loading" && <TableMessage colSpan={6} text="Loading services" />}
              {services.status === "error" && <TableMessage colSpan={6} text={services.error ?? "Services unavailable"} tone="error" />}
              {services.status === "success" && rows.length === 0 && <TableMessage colSpan={6} text={config.empty} />}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function inferCategory(values: Array<string | undefined>): ServiceCategory {
  const text = values.filter(Boolean).join(" ").toLowerCase();
  if (text.includes("bandwidth") || text.includes("traffic") || text.includes("gb")) return "bandwidth";
  if (text.includes("vps") || text.includes("vm") || text.includes("server")) return "vps";
  return "proxies";
}

function SummaryTile({ label, value, tone = "neutral" }: { label: string; value: string; tone?: "neutral" | "warn" }) {
  return (
    <div className={`bg-white border rounded p-4 ${tone === "warn" ? "border-amber-200" : "border-gray-200"}`}>
      <div className="text-[11px] text-gray-400 uppercase mb-1">{label}</div>
      <div className={`text-lg font-medium tabular-nums ${tone === "warn" ? "text-amber-700" : "text-gray-900"}`}>{value}</div>
    </div>
  );
}

function TableMessage({ colSpan, text, tone = "neutral" }: { colSpan: number; text: string; tone?: "neutral" | "error" }) {
  return (
    <tr>
      <td colSpan={colSpan} className={`p-4 text-center text-[12px] ${tone === "error" ? "text-red-600" : "text-gray-400"}`}>
        {text}
      </td>
    </tr>
  );
}
