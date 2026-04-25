"use client";

import { FormEvent, useState } from "react";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { canCancelJob, canMarkJobManualReview, canRetryJob, jobStatusLabel } from "@/lib/api/fulfillment";
import { compactDateTime, recordLabel } from "@/lib/api/format";
import type { CatalogProviderSource, JobQuery, Order, ProviderReadiness, ProvisioningJob, ServiceInstance } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { hiddenReference, providerSourceLabel } from "@/lib/api/viewModels";
import { PROVISIONING_JOBS } from "@/mocks/billingData";
import { AdminFilterBar, AdminFilterInput } from "../components/AdminFilterBar";
import { AdminJobTimelinePanel } from "../components/AdminJobTimelinePanel";
import { AdminProvisioningSummaryPanel } from "../components/AdminProvisioningSummaryPanel";
import { equalsFilter, hasActiveFilters, includesFilter, trimStringFilters } from "../lib/filterUtils";

interface ProvisioningRow {
  id: string;
  apiId?: string;
  live: boolean;
  order: string;
  service: string;
  tenant: string;
  provider: string;
  status: string;
  attempt: string;
  created: string;
  error: string;
  canRetry: boolean;
  canReview: boolean;
  canCancel: boolean;
  readiness?: ProviderReadiness;
  job?: ProvisioningJob;
}

type JobAction = "retry" | "manual-review" | "cancel";
type ProvisioningFilterFields = Required<Pick<JobQuery, "display_id" | "source_display_id" | "status">>;

const EMPTY_FILTERS: ProvisioningFilterFields = {
  display_id: "",
  source_display_id: "",
  status: "",
};

interface ActionState {
  id: string;
  action: JobAction;
  status: "running" | "success" | "error";
  message: string;
}

export function AdminProvisioning() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [draftFilters, setDraftFilters] = useState(EMPTY_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(EMPTY_FILTERS);
  const [manualReasons, setManualReasons] = useState<Record<string, string>>({});
  const [actionState, setActionState] = useState<ActionState | null>(null);
  const [selectedJobID, setSelectedJobID] = useState<string | null>(null);
  const jobs = useApiResource(
    () => billingApi.listAdminJobs({ job_type: "provider.provision", ...appliedFilters, limit: 100 }),
    `admin-provisioning-jobs:${refreshKey}:${JSON.stringify(appliedFilters)}`,
  );
  const summary = useApiResource(
    () => billingApi.getAdminJobSummary({ job_type: "provider.provision" }),
    `admin-provisioning-summary:${refreshKey}`,
  );
  const orders = useApiResource(
    () => billingApi.listAdminOrders({ limit: 100 }),
    `admin-provisioning-orders:${refreshKey}`,
  );
  const services = useApiResource(
    () => billingApi.listAdminServices({ limit: 100 }),
    `admin-provisioning-services:${refreshKey}`,
  );
  const providers = useApiResource(
    () => billingApi.listAdminProviderSources({ limit: 100 }),
    `admin-provisioning-providers:${refreshKey}`,
  );
  const readiness = useApiResource(
    () => billingApi.listAdminProviderReadiness({ status: "active", limit: 100 }),
    `admin-provisioning-readiness:${refreshKey}`,
  );
  const usingLive = jobs.status === "success";
  const rows = usingLive
    ? liveProvisioningRows(jobs.data ?? [], orders.data ?? [], services.data ?? [], providers.data ?? [], readiness.data ?? [])
    : filterDemoProvisioningRows(demoProvisioningRows(), appliedFilters);
  const selectedRow = rows.find((row) => row.apiId === selectedJobID) ?? null;
  const manualReview = rows.filter((row) => row.status === "manual_review");
  const failed = rows.filter((row) => row.status === "failed_retryable" || row.status === "failed_terminal" || row.status === "failed").length;
  const extraError = orders.error ?? services.error ?? providers.error;
  const activeFilters = hasActiveFilters(appliedFilters);
  let source = "Demo provisioning jobs";
  if (jobs.status === "error") {
    source = "Live job API unavailable. Showing demo queue data for the current filters.";
  } else if (jobs.status === "loading") {
    source = "Refreshing live provisioning jobs...";
  } else if (usingLive) {
    source = extraError
      ? "Live jobs loaded. Order or service links may be incomplete."
      : activeFilters
        ? "Live provisioning filters applied."
        : "Live provisioning jobs";
  } else if (activeFilters) {
    source = "Filters are applied to demo provisioning data.";
  }
  const statusTone = jobs.status === "error"
    ? "error"
    : jobs.status === "loading"
      ? "loading"
      : usingLive
        ? "success"
        : "default";

  function updateFilter(field: keyof ProvisioningFilterFields, value: string) {
    setDraftFilters((current) => ({ ...current, [field]: value }));
  }

  function applyFilters(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setAppliedFilters(trimStringFilters(draftFilters));
  }

  function resetFilters() {
    setDraftFilters(EMPTY_FILTERS);
    setAppliedFilters(EMPTY_FILTERS);
  }

  function updateManualReason(jobID: string, reason: string) {
    setManualReasons((current) => ({ ...current, [jobID]: reason }));
  }

  async function runJobAction(row: ProvisioningRow, action: JobAction) {
    if (!row.apiId) return;
    const reason = (manualReasons[row.apiId] ?? "").trim();
    if (action === "manual-review" && reason === "") {
      setActionState({
        id: row.apiId,
        action,
        status: "error",
        message: "Manual review reason is required.",
      });
      return;
    }
    if (action !== "manual-review" && !window.confirm(jobActionConfirmText(row, action))) {
      return;
    }

    setActionState({ id: row.apiId, action, status: "running", message: "Submitting..." });
    try {
      if (action === "retry") {
        await billingApi.retryAdminJob(row.apiId);
      } else if (action === "manual-review") {
        await billingApi.markAdminJobManualReview(row.apiId, { reason });
        setManualReasons((current) => ({ ...current, [row.apiId as string]: "" }));
      } else {
        await billingApi.cancelAdminJob(row.apiId, { reason: reason || undefined });
      }
      setActionState({ id: row.apiId, action, status: "success", message: "Updated" });
      setRefreshKey((current) => current + 1);
    } catch (error: unknown) {
      setActionState({
        id: row.apiId,
        action,
        status: "error",
        message: jobActionErrorMessage(error),
      });
    }
  }

  return (
    <div className="p-4 flex flex-col gap-4">
      <AdminProvisioningSummaryPanel
        summary={summary.data}
        loading={summary.status === "loading"}
        error={summary.error}
      />

      {(manualReview.length > 0 || failed > 0) && (
        <div className="bg-amber-50 border border-amber-200 text-amber-700 text-[12px] p-4 rounded flex items-center gap-3">
          <span className="font-medium tabular-nums">{manualReview.length + failed}</span>
          <span>job(s) need operator attention. Verify provider state before retrying.</span>
        </div>
      )}

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <div className="min-w-0 bg-white border border-gray-200 rounded">
          <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-3">
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Provisioning queue</h3>
            <span className="text-[11px] text-gray-400">{source}</span>
          </div>
          <AdminFilterBar onSubmit={applyFilters} onReset={resetFilters} statusText={source} statusTone={statusTone}>
            <AdminFilterInput
              label="Job public ID"
              value={draftFilters.display_id}
              onChange={(event) => updateFilter("display_id", event.target.value)}
              placeholder="71001"
              inputMode="numeric"
            />
            <AdminFilterInput
              label="Source public ID"
              value={draftFilters.source_display_id}
              onChange={(event) => updateFilter("source_display_id", event.target.value)}
              placeholder="23001"
              inputMode="numeric"
            />
            <AdminFilterInput
              label="Status"
              value={draftFilters.status}
              onChange={(event) => updateFilter("status", event.target.value)}
              placeholder="queued, failed, manual_review"
            />
          </AdminFilterBar>
          <div className="overflow-x-auto max-w-full">
            <table className="w-full text-[13px] border-collapse min-w-[1060px]">
              <thead>
                <tr className="bg-gray-50">
                  {["Job ID", "Order", "Service", "Tenant", "Provider", "Status", "Attempt", "Created", "Error", "Actions"].map((heading) => (
                    <th key={heading} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                      {heading}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => {
                  const rowState = row.apiId && actionState?.id === row.apiId ? actionState : null;
                  const running = rowState?.status === "running";
                  const selected = row.apiId === selectedJobID;
                  return (
                    <tr key={row.id} className={`hover:bg-gray-50 border-b border-gray-100 last:border-0 ${row.status === "manual_review" ? "bg-amber-50/40" : ""} ${selected ? "bg-red-50/50" : ""}`}>
                      <td className="p-4 text-[12px] text-[#D50C2D] font-medium">
                        {row.live && row.apiId ? (
                          <button
                            type="button"
                            onClick={() => setSelectedJobID(row.apiId as string)}
                            className="text-left font-medium text-[#D50C2D] underline-offset-2 hover:underline"
                          >
                            {row.id}
                          </button>
                        ) : row.id}
                      </td>
                      <td className="p-4 text-[12px] text-gray-500">{row.order}</td>
                      <td className="p-4 text-gray-700">{row.service}</td>
                      <td className="p-4 text-gray-500">{row.tenant}</td>
                      <td className="p-4 text-gray-500">{row.provider}</td>
                      <td className="p-4"><StatusBadge status={row.status} dot /></td>
                      <td className="p-4 text-center tabular-nums">{row.attempt}</td>
                      <td className="p-4 text-gray-400 tabular-nums">{row.created}</td>
                      <td className="p-4 text-[11px] text-red-600 max-w-[220px] truncate">{row.error}</td>
                      <td className="p-4">
                        <JobRecoveryControls
                          row={row}
                          reason={row.apiId ? manualReasons[row.apiId] ?? "" : ""}
                          running={running}
                          state={rowState}
                          onReasonChange={updateManualReason}
                          onAction={runJobAction}
                        />
                      </td>
                    </tr>
                  );
                })}
                {rows.length === 0 && (
                  <tr><td colSpan={10} className="p-4 text-center text-[12px] text-gray-400">No provisioning jobs</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
        <AdminJobTimelinePanel
          job={selectedRow?.job ?? null}
          orderLabel={selectedRow?.order}
          serviceLabel={selectedRow?.service}
          tenantLabel={selectedRow?.tenant}
          providerLabel={selectedRow?.provider}
          readiness={selectedRow?.readiness}
          readinessStatus={readiness.status}
        />
      </div>
    </div>
  );
}

function liveProvisioningRows(
  jobs: ProvisioningJob[],
  orders: Order[],
  services: ServiceInstance[],
  providers: CatalogProviderSource[],
  readinessRows: ProviderReadiness[],
): ProvisioningRow[] {
  const ordersByID = new Map(orders.map((order) => [order.id, order]));
  const servicesByOrderID = new Map(services.map((service) => [service.order_id, service]));
  const providersByID = new Map(providers.map((provider) => [provider.id, provider]));
  return jobs.map((job) => {
    const order = ordersByID.get(job.reference_id);
    const service = servicesByOrderID.get(job.reference_id);
    const provider = job.source_id ? providersByID.get(job.source_id) : undefined;
    const sourceDisplayID = provider?.display_id ?? job.source_display_id;
    const providerReadiness = readinessForJobSource(readinessRows, sourceDisplayID, order);
    const error = job.manual_review_reason || job.last_error_message_redacted || job.last_error_code || "-";
    return {
      id: recordLabel(job.display_id, "JOB-"),
      apiId: job.id,
      live: true,
      order: order
        ? recordLabel(order.display_id, "ORD-")
        : job.reference_type === "order" && job.reference_display_id
          ? recordLabel(job.reference_display_id, "ORD-")
          : hiddenReference("Order"),
      service: service ? recordLabel(service.display_id, "SVC-") : jobStatusLabel(job.status),
      tenant: hiddenReference("Tenant"),
      provider: provider
        ? providerSourceLabel(provider)
        : job.source_display_id
          ? recordLabel(job.source_display_id, "SRC-")
          : hiddenReference("Source"),
      status: job.status,
      attempt: `${job.attempt_count}/${job.max_attempts}`,
      created: compactDateTime(job.created_at),
      error,
      canRetry: canRetryJob(job.status),
      canReview: canMarkJobManualReview(job.status),
      canCancel: canCancelJob(job.status),
      readiness: providerReadiness,
      job,
    };
  });
}

function readinessForJobSource(
  rows: ProviderReadiness[],
  sourceDisplayID: number | undefined,
  order: Order | undefined,
): ProviderReadiness | undefined {
  if (!sourceDisplayID) return undefined;
  const planCode = planCodeFromSnapshot(order?.plan_snapshot);
  if (planCode) {
    const exact = rows.find((row) => row.source_display_id === sourceDisplayID && row.plan_code === planCode);
    if (exact) return exact;
  }
  return rows.find((row) => row.source_display_id === sourceDisplayID);
}

function planCodeFromSnapshot(snapshot: unknown): string | undefined {
  if (!snapshot || typeof snapshot !== "object" || Array.isArray(snapshot)) return undefined;
  const value = (snapshot as { plan_code?: unknown }).plan_code;
  return typeof value === "string" && value.trim() !== "" ? value : undefined;
}

function demoProvisioningRows(): ProvisioningRow[] {
  return PROVISIONING_JOBS.map((job) => ({
    id: job.id,
    live: false,
    order: job.order,
    service: job.service,
    tenant: job.tenant,
    provider: job.provider,
    status: job.status === "provisioning" ? "running" : job.status,
    attempt: String(job.attempt),
    created: job.age,
    error: job.error || "-",
    canRetry: false,
    canReview: false,
    canCancel: false,
  }));
}

function filterDemoProvisioningRows(rows: ProvisioningRow[], filters: ProvisioningFilterFields): ProvisioningRow[] {
  return rows.filter((row) => (
    includesFilter(row.id, filters.display_id)
    && includesFilter(row.provider, filters.source_display_id)
    && equalsFilter(row.status, filters.status)
  ));
}

interface JobRecoveryControlsProps {
  row: ProvisioningRow;
  reason: string;
  running: boolean;
  state: ActionState | null;
  onReasonChange: (jobID: string, reason: string) => void;
  onAction: (row: ProvisioningRow, action: JobAction) => void;
}

function JobRecoveryControls({ row, reason, running, state, onReasonChange, onAction }: JobRecoveryControlsProps) {
  if (!row.live || !row.apiId) {
    return <span className="text-[11px] text-gray-400">Demo read-only</span>;
  }
  const hasAction = row.canRetry || row.canReview || row.canCancel;
  if (!hasAction) {
    return <span className="text-[11px] text-gray-400">No action</span>;
  }
  return (
    <div className="flex min-w-[290px] flex-col gap-2">
      {row.canReview && (
        <input
          value={reason}
          onChange={(event) => onReasonChange(row.apiId as string, event.target.value)}
          disabled={running}
          placeholder="Review reason"
          className="h-8 rounded-md border border-gray-200 px-2 text-[12px] text-gray-700 outline-none focus:border-[#D50C2D]"
        />
      )}
      <div className="flex flex-wrap gap-2">
        {row.canRetry && (
          <button
            disabled={running}
            onClick={() => onAction(row, "retry")}
            className="inline-flex h-8 items-center justify-center rounded-md border border-emerald-600 bg-emerald-600 px-3 text-[12px] font-medium text-white disabled:cursor-not-allowed disabled:opacity-60"
          >
            Retry
          </button>
        )}
        {row.canReview && (
          <button
            disabled={running}
            onClick={() => onAction(row, "manual-review")}
            className="inline-flex h-8 items-center justify-center rounded-md border border-amber-200 bg-white px-3 text-[12px] font-medium text-amber-700 disabled:cursor-not-allowed disabled:opacity-60"
          >
            Review
          </button>
        )}
        {row.canCancel && (
          <button
            disabled={running}
            onClick={() => onAction(row, "cancel")}
            className="inline-flex h-8 items-center justify-center rounded-md border border-red-200 bg-white px-3 text-[12px] font-medium text-red-600 disabled:cursor-not-allowed disabled:opacity-60"
          >
            Cancel
          </button>
        )}
      </div>
      {state && (
        <div className={`text-[11px] ${state.status === "error" ? "text-red-500" : "text-gray-400"}`}>
          {state.message}
        </div>
      )}
    </div>
  );
}

function jobActionConfirmText(row: ProvisioningRow, action: JobAction): string {
  if (action === "retry") {
    return `Retry ${row.id} now?`;
  }
  return `Cancel ${row.id}?`;
}

function jobActionErrorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Job action failed.";
}
