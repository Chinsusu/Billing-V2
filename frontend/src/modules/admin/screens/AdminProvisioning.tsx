"use client";

import { FormEvent, useState } from "react";
import { billingApi } from "@/lib/api/billing";
import { technicalCodeLabel } from "@/lib/api/displayLabels";
import { canCancelJob, canMarkJobManualReview, canRetryJob, jobStatusLabel } from "@/lib/api/fulfillment";
import { compactDateTime, recordLabel } from "@/lib/api/format";
import type { CatalogProviderSource, JobQuery, Order, ProviderReadiness, ProvisioningJob, ServiceInstance } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { hiddenReference, providerSourceLabel } from "@/lib/api/viewModels";
import { PROVISIONING_JOBS } from "@/mocks/billingData";
import {
  AdminProvisioningQueueTable,
  type ActionState,
  type JobAction,
  type ProvisioningRow,
} from "../components/AdminProvisioningQueueTable";
import { AdminFilterBar, AdminFilterInput, AdminFilterSelect } from "../components/AdminFilterBar";
import { AdminJobTimelinePanel } from "../components/AdminJobTimelinePanel";
import { AdminProvisioningSummaryPanel } from "../components/AdminProvisioningSummaryPanel";
import { JOB_STATUS_OPTIONS } from "../lib/filterOptions";
import { equalsFilter, hasActiveFilters, includesFilter, trimStringFilters } from "../lib/filterUtils";

type ProvisioningFilterFields = Required<Pick<JobQuery, "display_id" | "source_display_id" | "status">>;

const EMPTY_FILTERS: ProvisioningFilterFields = {
  display_id: "",
  source_display_id: "",
  status: "",
};

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
            <AdminFilterSelect
              label="Status"
              value={draftFilters.status}
              onChange={(event) => updateFilter("status", event.target.value)}
              options={JOB_STATUS_OPTIONS}
            />
          </AdminFilterBar>
          <AdminProvisioningQueueTable
            rows={rows}
            selectedJobID={selectedJobID}
            actionState={actionState}
            manualReasons={manualReasons}
            onSelectJob={setSelectedJobID}
            onReasonChange={updateManualReason}
            onAction={runJobAction}
          />
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
    const error = job.manual_review_reason || job.last_error_message_redacted || technicalCodeLabel(job.last_error_code) || "-";
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
    provider: demoProviderLabel(job.provider),
    status: job.status === "provisioning" ? "running" : job.status,
    attempt: String(job.attempt),
    created: job.age,
    error: demoErrorLabel(job.error),
    canRetry: false,
    canReview: false,
    canCancel: false,
  }));
}

function demoProviderLabel(provider: string): string {
  return provider.includes("-") && !provider.includes(" ") ? technicalCodeLabel(provider) : provider;
}

function demoErrorLabel(error: string): string {
  const normalized = error.trim();
  if (!normalized) return "-";
  return normalized
    .split(/\s+(?:-|\u2014)\s+/)
    .filter(Boolean)
    .map((part) => technicalCodeLabel(part.replace(/\s+/g, "_")))
    .join(": ");
}

function filterDemoProvisioningRows(rows: ProvisioningRow[], filters: ProvisioningFilterFields): ProvisioningRow[] {
  return rows.filter((row) => (
    includesFilter(row.id, filters.display_id)
    && includesFilter(row.provider, filters.source_display_id)
    && equalsFilter(row.status, filters.status)
  ));
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
