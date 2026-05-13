"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { ServiceAccessReveal } from "@/components/ui/ServiceAccessReveal";
import { billingApi } from "@/lib/api/billing";
import { clientServiceCategory, clientServiceOrderLabel, clientServicePlanLabel, clientServiceSourceLabel } from "@/lib/api/clientViewModels";
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
        apiId: service.id,
        category: clientServiceCategory(service),
        label: recordLabel(service.display_id, "SVC-"),
        plan: clientServicePlanLabel(service),
        source: clientServiceSourceLabel(service),
        detail: clientServiceOrderLabel(service),
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
                {["Service ID", "Plan", "Region / Source", "Order", "Expires", "Status", "Access"].map((heading) => (
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
                  <td className="p-4 text-[12px] text-gray-500">{service.plan}</td>
                  <td className="p-4 text-gray-500">{service.source}</td>
                  <td className="p-4 text-gray-500">{service.detail}</td>
                  <td className="p-4 text-gray-500">{service.expiry}</td>
                  <td className="p-4"><StatusBadge status={service.status} dot /></td>
                  <td className="p-4"><ServiceAccessReveal scope="client" serviceId={service.apiId} reason="Client portal reveal" /></td>
                </tr>
              ))}
              {services.status === "loading" && <TableMessage colSpan={7} text="Loading services" />}
              {services.status === "error" && <TableMessage colSpan={7} text={services.error ?? "Services unavailable"} tone="error" />}
              {services.status === "success" && rows.length === 0 && <TableMessage colSpan={7} text={config.empty} />}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
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
