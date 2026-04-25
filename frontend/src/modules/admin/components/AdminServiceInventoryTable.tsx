"use client";

import { useState } from "react";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { TablePagination } from "@/components/ui/TablePagination";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, recordLabel } from "@/lib/api/format";
import type { ServiceInstance } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { hiddenReference } from "@/lib/api/viewModels";

export type ServiceFamily = "proxy" | "vps" | "bandwidth";

export interface AdminServiceDemoRow {
  id: string;
  service: string;
  owner: string;
  tenant: string;
  resource: string;
  plan: string;
  region: string;
  status: string;
  billingStatus: string;
  created: string;
  expires: string;
  provider: string;
  note: string;
}

interface ServiceInventoryRow extends AdminServiceDemoRow {
  family: ServiceFamily | "unknown";
}

type ServiceInstanceWithSnapshots = ServiceInstance & {
  product_snapshot?: unknown;
  plan_snapshot?: unknown;
  price_snapshot?: unknown;
};

interface AdminServiceInventoryTableProps {
  family: ServiceFamily;
  title: string;
  demoRows: AdminServiceDemoRow[];
}

const FAMILY_LABELS: Record<ServiceFamily, string> = {
  proxy: "proxy",
  vps: "VPS",
  bandwidth: "bandwidth",
};

function isRecord(value: unknown): value is Record<string, unknown> {
  return Boolean(value) && typeof value === "object" && !Array.isArray(value);
}

function recordValue(value: unknown, keys: string[]) {
  if (!isRecord(value)) return "";

  for (const key of keys) {
    const item = value[key];
    if (typeof item === "string" && item.trim()) return item.trim();
    if (typeof item === "number") return String(item);
  }

  return "";
}

function snapshotText(service: ServiceInstance) {
  const record = service as ServiceInstanceWithSnapshots;
  const values = [
    record.product_snapshot,
    record.plan_snapshot,
    record.price_snapshot,
    service.external_resource_id,
    service.provider_source_id,
    service.tenant_plan_id,
  ];

  return values
    .map((value) => {
      if (!value) return "";
      if (typeof value === "string") return value;
      try {
        return JSON.stringify(value);
      } catch {
        return "";
      }
    })
    .join(" ")
    .toLowerCase();
}

function inferFamily(service: ServiceInstance): ServiceFamily | "unknown" {
  const text = snapshotText(service);
  if (text.includes("bandwidth") || text.includes("traffic") || text.includes("quota")) return "bandwidth";
  if (text.includes("vps") || text.includes("virtual") || text.includes("server") || text.includes("rdp")) return "vps";
  if (
    text.includes("proxy")
    || text.includes("residential")
    || text.includes("socks")
    || text.includes("mobile")
    || text.includes("isp")
  ) {
    return "proxy";
  }
  return "unknown";
}

function liveRow(service: ServiceInstance): ServiceInventoryRow {
  const record = service as ServiceInstanceWithSnapshots;
  const productName = recordValue(record.product_snapshot, ["name", "product_name", "product_type"]);
  const planName = recordValue(record.plan_snapshot, ["name", "plan_name", "plan_code"]);
  const region = recordValue(record.plan_snapshot, ["region", "location", "datacenter"])
    || recordValue(record.product_snapshot, ["region", "location", "datacenter"]);

  return {
    id: recordLabel(service.display_id, "SVC-"),
    service: planName || productName || recordLabel(service.display_id, "Service "),
    owner: hiddenReference("Order"),
    tenant: hiddenReference("Tenant"),
    resource: hiddenReference("Resource"),
    plan: planName || hiddenReference("Plan"),
    region: region || "-",
    status: service.status,
    billingStatus: service.billing_status,
    created: compactDateTime(service.created_at),
    expires: compactDateTime(service.term_end),
    provider: hiddenReference("Source"),
    note: service.suspension_reason || "Live read-only",
    family: inferFamily(service),
  };
}

function sourceText(status: string, usingLive: boolean, family: ServiceFamily) {
  if (status === "error") {
    return `Live API unavailable. Showing demo ${FAMILY_LABELS[family]} data.`;
  }
  if (status === "loading") {
    return `Refreshing live ${FAMILY_LABELS[family]} inventory...`;
  }
  if (usingLive) {
    return `Live ${FAMILY_LABELS[family]} inventory`;
  }
  return `Demo ${FAMILY_LABELS[family]} inventory`;
}

export function AdminServiceInventoryTable({ family, title, demoRows }: AdminServiceInventoryTableProps) {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const services = useApiResource(
    () => billingApi.listAdminServices({ limit: 100 }),
    `admin-services:${family}`,
  );
  const usingLive = services.status === "success";
  const liveRows = (services.data ?? [])
    .map(liveRow)
    .filter((row) => row.family === family);
  const rows: AdminServiceDemoRow[] = usingLive ? liveRows : demoRows;
  const displayed = limit === -1 ? rows : rows.slice((page - 1) * limit, page * limit);
  const statusText = sourceText(services.status, usingLive, family);

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded shadow-sm text-[12px]">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">{title}</h3>
          <div className="flex flex-wrap items-center justify-end gap-3">
            <span className="text-[11px] text-gray-400">{statusText}</span>
            <span className="text-[11px] text-gray-400">{rows.length} services</span>
          </div>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="min-w-[980px] w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-50 text-gray-500">
                {["ID", "Service", "Owner", "Resource", "Plan", "Region", "Status", "Billing", "Created", "Expire", "Provider", "Action"].map((label) => (
                  <th key={label} className="font-medium p-3 border-b border-gray-200 text-[11px] tracking-wider">
                    {label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="bg-white text-gray-600">
              {displayed.map((service) => (
                <tr key={service.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="p-3 text-[#D50C2D] font-medium">{service.id}</td>
                  <td className="p-3">
                    <div className="font-medium text-gray-800">{service.service}</div>
                    <div className="text-[11px] text-gray-400">{service.note}</div>
                  </td>
                  <td className="p-3 text-[11px]">
                    <div className="font-medium text-gray-800">{service.owner}</div>
                    <div className="text-gray-400">{service.tenant}</div>
                  </td>
                  <td className="p-3 text-[11px] text-gray-500">{service.resource}</td>
                  <td className="p-3">
                    <span className="rounded bg-indigo-50 px-2 py-1 text-[10px] font-medium text-indigo-700">
                      {service.plan}
                    </span>
                  </td>
                  <td className="p-3 text-[11px] text-gray-500">{service.region}</td>
                  <td className="p-3"><StatusBadge status={service.status} dot /></td>
                  <td className="p-3"><StatusBadge status={service.billingStatus} dot /></td>
                  <td className="p-3 text-[11px] text-gray-400">{service.created}</td>
                  <td className="p-3 text-[11px] text-gray-400">{service.expires}</td>
                  <td className="p-3 text-[11px] text-gray-500">{service.provider}</td>
                  <td className="p-3">
                    <span className="inline-flex h-8 items-center justify-center rounded-md border border-gray-200 bg-gray-50 px-3 text-[12px] font-medium text-gray-400">
                      Read-only
                    </span>
                  </td>
                </tr>
              ))}
              {usingLive && rows.length === 0 && (
                <tr>
                  <td colSpan={12} className="p-4 text-center text-[12px] text-gray-400">
                    No live {FAMILY_LABELS[family]} services found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
        <TablePagination page={page} setPage={setPage} limit={limit} setLimit={setLimit} total={rows.length} />
      </div>
    </div>
  );
}
