import { billingApi } from "@/lib/api/billing";
import type { ProvisioningJob } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { mapAdminAuditLogView, type AdminAuditLogView } from "@/lib/api/viewModels";

interface AdminJobRecoveryAuditPanelProps {
  job: ProvisioningJob;
}

const RECOVERY_ACTIONS = new Set(["job.retry", "job.manual_review", "job.cancel"]);

const ACTION_LABEL: Record<string, string> = {
  "job.retry": "Retry",
  "job.manual_review": "Manual review",
  "job.cancel": "Cancel",
};

const ACTION_STYLE: Record<string, string> = {
  "job.retry": "bg-emerald-50 text-emerald-700",
  "job.manual_review": "bg-amber-50 text-amber-700",
  "job.cancel": "bg-red-50 text-red-700",
};

const ACTOR_STYLE: Record<string, string> = {
  user: "bg-slate-100 text-slate-700",
  worker: "bg-orange-50 text-orange-700",
  provider_webhook: "bg-cyan-50 text-cyan-700",
  system: "bg-gray-100 text-gray-500",
};

export function AdminJobRecoveryAuditPanel({ job }: AdminJobRecoveryAuditPanelProps) {
  const audits = useApiResource(
    () => billingApi.listAdminAuditLogs({
      target_type: "job",
      target_id: job.id,
      limit: 20,
    }),
    `admin-job-recovery-audit:${job.id}`,
  );
  const rows = (audits.data ?? [])
    .filter((log) => RECOVERY_ACTIONS.has(log.action))
    .map(mapAdminAuditLogView);

  return (
    <div className="border-t border-gray-100 p-4">
      <div className="flex items-center justify-between gap-3">
        <h4 className="m-0 text-[13px] font-medium text-gray-900">Recovery audit</h4>
        <span className="text-[11px] text-gray-400">{rows.length} action(s)</span>
      </div>
      <RecoveryAuditContent rows={rows} loading={audits.status === "loading"} error={audits.error} />
    </div>
  );
}

function RecoveryAuditContent({
  rows,
  loading,
  error,
}: {
  rows: AdminAuditLogView[];
  loading: boolean;
  error: string | null;
}) {
  if (loading) {
    return <p className="mt-3 mb-0 text-[12px] text-gray-400">Loading recovery audit...</p>;
  }
  if (error) {
    return <p className="mt-3 mb-0 text-[12px] text-amber-700">Recovery audit unavailable.</p>;
  }
  if (rows.length === 0) {
    return <p className="mt-3 mb-0 text-[12px] text-gray-400">No recovery actions recorded.</p>;
  }
  return (
    <ol className="mt-4 flex flex-col gap-3">
      {rows.map((log) => (
        <li key={log.id} className="rounded border border-gray-100 bg-gray-50 p-3">
          <div className="flex flex-wrap items-center gap-2">
            <span className={`rounded px-1.5 py-0.5 text-[10px] font-medium ${ACTION_STYLE[log.action] ?? "bg-gray-100 text-gray-500"}`}>
              {ACTION_LABEL[log.action] ?? log.action}
            </span>
            <span className={`rounded px-1.5 py-0.5 text-[10px] font-medium ${actorStyle(log.actor)}`}>
              {log.actor}
            </span>
            <span className="text-[11px] text-gray-500">{log.actorName}</span>
          </div>
          <div className="mt-2 grid grid-cols-2 gap-x-3 gap-y-1 text-[11px] text-gray-500">
            <span>{log.id}</span>
            <span>{log.ts}</span>
            <span>{log.requestId}</span>
            <span>{log.detail}</span>
          </div>
        </li>
      ))}
    </ol>
  );
}

function actorStyle(actorType: string): string {
  return ACTOR_STYLE[actorType] ?? ACTOR_STYLE.system;
}
