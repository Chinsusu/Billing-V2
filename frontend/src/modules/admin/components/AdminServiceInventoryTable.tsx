"use client";

import { FormEvent, useState } from "react";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { ServiceAccessReveal } from "@/components/ui/ServiceAccessReveal";
import { TablePagination } from "@/components/ui/TablePagination";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, recordLabel } from "@/lib/api/format";
import type { AdminServiceQuery, ServiceInstance } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { hiddenReference } from "@/lib/api/viewModels";
import { AdminFilterBar, AdminFilterInput, AdminFilterSelect } from "./AdminFilterBar";
import { SERVICE_STATUS_OPTIONS } from "../lib/filterOptions";
import { equalsFilter, hasActiveFilters, includesFilter, trimStringFilters } from "../lib/filterUtils";

export type ServiceFamily = "proxy" | "vps" | "bandwidth";

export interface AdminServiceDemoRow {
  id: string;
  apiId?: string;
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

type ServiceFilterFields = Required<Pick<
  AdminServiceQuery,
  "display_id" | "order_display_id" | "provider_source_display_id" | "status"
>>;

const EMPTY_FILTERS: ServiceFilterFields = {
  display_id: "",
  order_display_id: "",
  provider_source_display_id: "",
  status: "",
};

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
    apiId: service.id,
    service: planName || productName || recordLabel(service.display_id, "Service "),
    owner: service.order_display_id ? recordLabel(service.order_display_id, "ORD-") : hiddenReference("Order"),
    tenant: service.buyer_display_id ? recordLabel(service.buyer_display_id, "ACC-") : hiddenReference("Account"),
    resource: hiddenReference("Resource"),
    plan: planName || hiddenReference("Plan"),
    region: region || "-",
    status: service.status,
    billingStatus: service.billing_status,
    created: compactDateTime(service.created_at),
    expires: compactDateTime(service.term_end),
    provider: service.provider_source_display_id ? recordLabel(service.provider_source_display_id, "SRC-") : hiddenReference("Source"),
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

function filterDemoServices(rows: AdminServiceDemoRow[], filters: ServiceFilterFields): AdminServiceDemoRow[] {
  return rows.filter((row) => (
    includesFilter(row.id, filters.display_id)
    && includesFilter(row.owner, filters.order_display_id)
    && includesFilter(row.provider, filters.provider_source_display_id)
    && equalsFilter(row.status, filters.status)
  ));
}

export function AdminServiceInventoryTable({ family, title, demoRows }: AdminServiceInventoryTableProps) {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [draftFilters, setDraftFilters] = useState(EMPTY_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(EMPTY_FILTERS);
  const services = useApiResource(
    () => billingApi.listAdminServices({ ...appliedFilters, limit: 100 }),
    `admin-services:${family}:${JSON.stringify(appliedFilters)}`,
  );
  const usingLive = services.status === "success";
  const liveRows = (services.data ?? [])
    .map(liveRow)
    .filter((row) => row.family === family);
  const rows: AdminServiceDemoRow[] = usingLive ? liveRows : filterDemoServices(demoRows, appliedFilters);
  const displayed = limit === -1 ? rows : rows.slice((page - 1) * limit, page * limit);
  const statusText = sourceText(services.status, usingLive, family);
  const activeFilters = hasActiveFilters(appliedFilters);
  const statusTone = services.status === "error"
    ? "error"
    : services.status === "loading"
      ? "loading"
      : usingLive
        ? "success"
        : "default";
  const filterStatusText = usingLive
    ? activeFilters
      ? `Live ${FAMILY_LABELS[family]} filters applied.`
      : statusText
    : activeFilters
      ? `Filters are applied to demo ${FAMILY_LABELS[family]} data.`
      : statusText;

  function updateFilter(field: keyof ServiceFilterFields, value: string) {
    setDraftFilters((current) => ({ ...current, [field]: value }));
  }

  function applyFilters(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setAppliedFilters(trimStringFilters(draftFilters));
    setPage(1);
  }

  function resetFilters() {
    setDraftFilters(EMPTY_FILTERS);
    setAppliedFilters(EMPTY_FILTERS);
    setPage(1);
  }

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
        <AdminFilterBar onSubmit={applyFilters} onReset={resetFilters} statusText={filterStatusText} statusTone={statusTone}>
          <AdminFilterInput
            label="Service public ID"
            value={draftFilters.display_id}
            onChange={(event) => updateFilter("display_id", event.target.value)}
            placeholder="50002"
            inputMode="numeric"
          />
          <AdminFilterInput
            label="Order public ID"
            value={draftFilters.order_display_id}
            onChange={(event) => updateFilter("order_display_id", event.target.value)}
            placeholder="30005"
            inputMode="numeric"
          />
          <AdminFilterInput
            label="Source public ID"
            value={draftFilters.provider_source_display_id}
            onChange={(event) => updateFilter("provider_source_display_id", event.target.value)}
            placeholder="23001"
            inputMode="numeric"
          />
          <AdminFilterSelect
            label="Status"
            value={draftFilters.status}
            onChange={(event) => updateFilter("status", event.target.value)}
            options={SERVICE_STATUS_OPTIONS}
          />
        </AdminFilterBar>
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
                    {service.apiId ? (
                      <ServiceAccessReveal scope="admin" serviceId={service.apiId} reason="Admin service support reveal" />
                    ) : (
                      <span className="inline-flex h-8 items-center justify-center rounded-md border border-gray-200 bg-gray-50 px-3 text-[12px] font-medium text-gray-400">
                        Read-only
                      </span>
                    )}
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
