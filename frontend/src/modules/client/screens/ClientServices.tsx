"use client";

import { useState } from "react";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { ServiceAccessReveal } from "@/components/ui/ServiceAccessReveal";
import { billingApi } from "@/lib/api/billing";
import { clientServiceCategory, clientServiceOrderLabel, clientServicePlanLabel, clientServiceSourceLabel } from "@/lib/api/clientViewModels";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";

type ServiceCategory = "proxies" | "vps" | "bandwidth";

interface ClientServicesProps {
  category: ServiceCategory;
}

type Notice = { type: "success" | "error"; text: string };

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
  const [refreshKey, setRefreshKey] = useState(0);
  const [busyServiceID, setBusyServiceID] = useState<string | null>(null);
  const [notice, setNotice] = useState<Notice | null>(null);
  const services = useApiResource(billingApi.listClientServices, `client-services:${refreshKey}`);
  const wallets = useApiResource(billingApi.listClientWallets, `client-service-wallets:${refreshKey}`);
  const wallet = wallets.data?.find((item) => item.status === "active") ?? null;
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
        suspensionReason: service.suspension_reason,
      })) ?? []
    : null;
  const sourceRows = liveRows ?? [];
  const rows = sourceRows.filter((service) => service.category === category);
  const suspended = rows.filter((service) => service.status === "suspended" || service.status === "overdue").length;
  const config = CATEGORY_CONFIG[category];

  async function handleRenew(service: (typeof rows)[number]) {
    if (!wallet) {
      setNotice({ type: "error", text: "No active wallet is available for renewal." });
      return;
    }
    setBusyServiceID(service.apiId);
    setNotice(null);
    try {
      const result = await billingApi.renewClientService(service.apiId, {
        wallet_id: wallet.id,
        from_status: service.status,
        reason: "Client portal renewal",
      });
      setNotice({
        type: "success",
        text: `Renewed ${service.label}. Charged ${moneyMinor(result.amount_minor, result.currency)} on invoice ${recordLabel(result.invoice.display_id, "INV-")}.`,
      });
      setRefreshKey((current) => current + 1);
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : "Renewal failed.";
      setNotice({ type: "error", text: message });
    } finally {
      setBusyServiceID(null);
    }
  }

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
          <div className="flex items-center gap-3">
            {notice && (
              <span className={`text-[11px] font-medium ${notice.type === "error" ? "text-red-600" : "text-green-700"}`}>
                {notice.text}
              </span>
            )}
            <span className="text-[11px] text-gray-400">{services.status === "loading" ? "Loading" : `${rows.length} records`}</span>
          </div>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[860px]">
            <thead>
              <tr className="bg-gray-50">
                {["Service ID", "Plan", "Region / Source", "Order", "Expires", "Status", "Access", "Renew"].map((heading) => (
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
                  <td className="p-4">
                    <RenewAction
                      busy={busyServiceID === service.apiId}
                      service={service}
                      walletReady={Boolean(wallet)}
                      onRenew={() => handleRenew(service)}
                    />
                  </td>
                </tr>
              ))}
              {services.status === "loading" && <TableMessage colSpan={8} text="Loading services" />}
              {services.status === "error" && <TableMessage colSpan={8} text={services.error ?? "Services unavailable"} tone="error" />}
              {services.status === "success" && rows.length === 0 && <TableMessage colSpan={8} text={config.empty} />}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function RenewAction({
  busy,
  service,
  walletReady,
  onRenew,
}: {
  busy: boolean;
  service: { status: string; suspensionReason?: string };
  walletReady: boolean;
  onRenew: () => void;
}) {
  const supportedStatus = service.status === "active" ||
    service.status === "expired" ||
    (service.status === "suspended" && service.suspensionReason === "expiry");
  if (!supportedStatus) {
    return <span className="text-[11px] text-gray-400">Support</span>;
  }
  return (
    <button
      type="button"
      className="inline-flex h-8 items-center justify-center rounded-md border border-[#D50C2D] px-3 text-[12px] font-medium text-[#D50C2D] transition-colors hover:bg-[#D50C2D] hover:text-white disabled:cursor-not-allowed disabled:border-gray-200 disabled:text-gray-400 disabled:hover:bg-white"
      disabled={busy || !walletReady}
      onClick={onRenew}
    >
      {busy ? "Renewing" : "Renew"}
    </button>
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
