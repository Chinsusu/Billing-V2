"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, recordLabel, shortID } from "@/lib/api/format";
import type { ProvisioningJobAttempt } from "@/lib/api/jobTypes";
import type { ProvisioningJob } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { AdminJobRecoveryAuditPanel } from "./AdminJobRecoveryAuditPanel";

interface AdminJobTimelinePanelProps {
  job: ProvisioningJob | null;
  orderLabel?: string;
  serviceLabel?: string;
  tenantLabel?: string;
  providerLabel?: string;
}

export function AdminJobTimelinePanel({
  job,
  orderLabel = "-",
  serviceLabel = "-",
  tenantLabel = "-",
  providerLabel = "-",
}: AdminJobTimelinePanelProps) {
  const attempts = useApiResource(
    () => job ? billingApi.listAdminJobAttempts(job.id, { limit: 20 }) : Promise.resolve([]),
    `admin-job-attempts:${job?.id ?? "none"}`,
  );

  if (!job) {
    return (
      <aside className="bg-white border border-gray-200 rounded p-4 text-[12px] text-gray-400">
        No live job selected.
      </aside>
    );
  }

  const issue = latestJobIssue(job);
  const attemptRows = attempts.data ?? [];
  return (
    <aside className="bg-white border border-gray-200 rounded">
      <div className="border-b border-gray-100 p-4">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <p className="m-0 text-[11px] uppercase tracking-wide text-gray-400">Job detail</p>
            <h3 className="m-0 mt-1 text-[16px] font-semibold text-[#D50C2D]">
              {recordLabel(job.display_id, "JOB-")}
            </h3>
          </div>
          <StatusBadge status={job.status} dot />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-x-4 gap-y-3 p-4 text-[12px]">
        <Detail label="Attempt" value={`${job.attempt_count}/${job.max_attempts}`} />
        <Detail label="Next" value={compactDateTime(job.next_attempt_at)} />
        <Detail label="Order" value={orderLabel} />
        <Detail label="Service" value={serviceLabel} />
        <Detail label="Tenant" value={tenantLabel} />
        <Detail label="Provider" value={providerLabel} />
        <Detail label="Created" value={compactDateTime(job.created_at)} />
        <Detail label="Updated" value={compactDateTime(job.updated_at)} />
      </div>

      {(issue !== "-" || job.manual_review_reason) && (
        <div className="border-t border-gray-100 p-4 text-[12px]">
          <p className="m-0 text-[11px] uppercase tracking-wide text-gray-400">Recovery context</p>
          {issue !== "-" && <p className="m-0 mt-2 text-red-600">{issue}</p>}
          {job.manual_review_reason && (
            <p className="m-0 mt-2 text-amber-700">{job.manual_review_reason}</p>
          )}
        </div>
      )}

      <div className="border-t border-gray-100 p-4">
        <div className="flex items-center justify-between gap-3">
          <h4 className="m-0 text-[13px] font-medium text-gray-900">Attempt timeline</h4>
          <span className="text-[11px] text-gray-400">{attemptRows.length} item(s)</span>
        </div>
        <AttemptTimeline attempts={attemptRows} loading={attempts.status === "loading"} error={attempts.error} />
      </div>

      <AdminJobRecoveryAuditPanel job={job} />

      <div className="border-t border-gray-100 p-4 text-[11px] text-gray-400">
        UUID {shortID(job.id)} / correlation {shortID(job.correlation_id)}
      </div>
    </aside>
  );
}

function Detail({ label, value }: { label: string; value: string }) {
  return (
    <div className="min-w-0">
      <p className="m-0 text-[11px] uppercase tracking-wide text-gray-400">{label}</p>
      <p className="m-0 mt-1 truncate text-gray-700" title={value}>{value}</p>
    </div>
  );
}

function AttemptTimeline({
  attempts,
  loading,
  error,
}: {
  attempts: ProvisioningJobAttempt[];
  loading: boolean;
  error: string | null;
}) {
  if (loading) {
    return <p className="mt-3 mb-0 text-[12px] text-gray-400">Loading attempts...</p>;
  }
  if (error) {
    return <p className="mt-3 mb-0 text-[12px] text-red-600">{error}</p>;
  }
  if (attempts.length === 0) {
    return <p className="mt-3 mb-0 text-[12px] text-gray-400">No attempts recorded.</p>;
  }
  return (
    <ol className="mt-4 flex flex-col gap-3 border-l border-gray-200 pl-4">
      {attempts.map((attempt) => (
        <li key={attempt.id} className="relative">
          <span className="absolute -left-[21px] top-1 h-2 w-2 rounded-full bg-[#D50C2D]" />
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-[12px] font-medium text-gray-900">
              {recordLabel(attempt.display_id, "ATT-")}
            </span>
            <StatusBadge status={attempt.result} />
            <span className="text-[11px] text-gray-400">#{attempt.attempt_number}</span>
          </div>
          <div className="mt-1 grid grid-cols-2 gap-x-3 gap-y-1 text-[11px] text-gray-500">
            <span>Worker {attempt.worker_id || "-"}</span>
            <span>{durationLabel(attempt.duration_ms)}</span>
            <span>{compactDateTime(attempt.started_at)}</span>
            <span>{compactDateTime(attempt.finished_at)}</span>
          </div>
          {(attempt.error_code || attempt.error_message_redacted) && (
            <p className="m-0 mt-2 text-[11px] text-red-600">
              {[attempt.error_code, attempt.error_message_redacted].filter(Boolean).join(" / ")}
            </p>
          )}
        </li>
      ))}
    </ol>
  );
}

function latestJobIssue(job: ProvisioningJob): string {
  return job.last_error_message_redacted || job.last_error_code || "-";
}

function durationLabel(value?: number): string {
  if (value === undefined || value === null) return "-";
  if (value < 1000) return `${value}ms`;
  return `${(value / 1000).toFixed(1)}s`;
}
