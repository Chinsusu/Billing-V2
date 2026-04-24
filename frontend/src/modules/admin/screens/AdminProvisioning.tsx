"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { jobStatusLabel } from "@/lib/api/fulfillment";
import { compactDateTime, recordLabel, shortID } from "@/lib/api/format";
import type { CatalogProviderSource, Order, ProvisioningJob, ServiceInstance } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { PROVISIONING_JOBS } from "@/mocks/billingData";

interface ProvisioningRow {
  id: string;
  order: string;
  service: string;
  tenant: string;
  provider: string;
  status: string;
  attempt: string;
  created: string;
  error: string;
}

export function AdminProvisioning() {
  const jobs = useApiResource(
    () => billingApi.listAdminJobs({ job_type: "provider.provision", limit: 100 }),
    "admin-provisioning-jobs",
  );
  const orders = useApiResource(
    () => billingApi.listAdminOrders({ limit: 100 }),
    "admin-provisioning-orders",
  );
  const services = useApiResource(
    () => billingApi.listAdminServices({ limit: 100 }),
    "admin-provisioning-services",
  );
  const providers = useApiResource(
    () => billingApi.listAdminProviderSources({ limit: 100 }),
    "admin-provisioning-providers",
  );
  const usingLive = jobs.status === "success";
  const rows = usingLive
    ? liveProvisioningRows(jobs.data ?? [], orders.data ?? [], services.data ?? [], providers.data ?? [])
    : demoProvisioningRows();
  const manualReview = rows.filter((row) => row.status === "manual_review");
  const failed = rows.filter((row) => row.status === "failed_retryable" || row.status === "failed_terminal" || row.status === "failed").length;
  const extraError = orders.error ?? services.error ?? providers.error;
  const source = jobs.status === "error"
    ? "Live job API unavailable. Showing demo queue data."
    : jobs.status === "loading"
      ? "Refreshing live provisioning jobs..."
      : usingLive
        ? extraError ? "Live jobs loaded. Order or service links may be incomplete." : "Live provisioning jobs"
        : "Demo provisioning jobs";

  return (
    <div className="p-4 flex flex-col gap-4">
      {(manualReview.length > 0 || failed > 0) && (
        <div className="bg-amber-50 border border-amber-200 text-amber-700 text-[12px] p-4 rounded flex items-center gap-3">
          <span className="font-medium tabular-nums">{manualReview.length + failed}</span>
          <span>job(s) need operator attention. Verify provider state before retrying.</span>
        </div>
      )}

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-3">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Provisioning queue</h3>
          <span className="text-[11px] text-gray-400">{source}</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[1060px]">
            <thead>
              <tr className="bg-gray-50">
                {["Job ID", "Order", "Service", "Tenant", "Provider", "Status", "Attempt", "Created", "Error"].map((heading) => (
                  <th key={heading} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {heading}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.id} className={`hover:bg-gray-50 border-b border-gray-100 last:border-0 ${row.status === "manual_review" ? "bg-amber-50/40" : ""}`}>
                  <td className="p-4 text-[12px] text-[#D50C2D] font-medium">{row.id}</td>
                  <td className="p-4 text-[12px] text-gray-500">{row.order}</td>
                  <td className="p-4 text-gray-700">{row.service}</td>
                  <td className="p-4 text-gray-500">{row.tenant}</td>
                  <td className="p-4 text-gray-500">{row.provider}</td>
                  <td className="p-4"><StatusBadge status={row.status} dot /></td>
                  <td className="p-4 text-center tabular-nums">{row.attempt}</td>
                  <td className="p-4 text-gray-400 tabular-nums">{row.created}</td>
                  <td className="p-4 text-[11px] text-red-600 max-w-[220px] truncate">{row.error}</td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr><td colSpan={9} className="p-4 text-center text-[12px] text-gray-400">No provisioning jobs</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function liveProvisioningRows(
  jobs: ProvisioningJob[],
  orders: Order[],
  services: ServiceInstance[],
  providers: CatalogProviderSource[],
): ProvisioningRow[] {
  const ordersByID = new Map(orders.map((order) => [order.id, order]));
  const servicesByOrderID = new Map(services.map((service) => [service.order_id, service]));
  const providersByID = new Map(providers.map((provider) => [provider.id, provider]));
  return jobs.map((job) => {
    const order = ordersByID.get(job.reference_id);
    const service = servicesByOrderID.get(job.reference_id);
    const provider = job.source_id ? providersByID.get(job.source_id) : undefined;
    const error = job.manual_review_reason || job.last_error_message_redacted || job.last_error_code || "-";
    return {
      id: recordLabel(job.display_id, "JOB-"),
      order: order ? recordLabel(order.display_id, "ORD-") : shortID(job.reference_id),
      service: service ? recordLabel(service.display_id, "SVC-") : jobStatusLabel(job.status),
      tenant: order?.tenant_id ? shortID(order.tenant_id) : shortID(job.tenant_id),
      provider: provider ? `${provider.name} (${recordLabel(provider.display_id, "SRC-")})` : shortID(job.source_id),
      status: job.status,
      attempt: `${job.attempt_count}/${job.max_attempts}`,
      created: compactDateTime(job.created_at),
      error,
    };
  });
}

function demoProvisioningRows(): ProvisioningRow[] {
  return PROVISIONING_JOBS.map((job) => ({
    id: job.id,
    order: job.order,
    service: job.service,
    tenant: job.tenant,
    provider: job.provider,
    status: job.status === "provisioning" ? "running" : job.status,
    attempt: String(job.attempt),
    created: job.age,
    error: job.error || "-",
  }));
}
