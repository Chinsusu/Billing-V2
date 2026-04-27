import { StatusBadge } from "@/components/ui/StatusBadge";
import type { ProviderReadiness, ProvisioningJob } from "@/lib/api/types";

export interface ProvisioningRow {
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

export type JobAction = "retry" | "manual-review" | "cancel";

export interface ActionState {
  id: string;
  action: JobAction;
  status: "running" | "success" | "error";
  message: string;
}

interface AdminProvisioningQueueTableProps {
  rows: ProvisioningRow[];
  selectedJobID: string | null;
  actionState: ActionState | null;
  manualReasons: Record<string, string>;
  onSelectJob: (jobID: string) => void;
  onReasonChange: (jobID: string, reason: string) => void;
  onAction: (row: ProvisioningRow, action: JobAction) => void;
}

export function AdminProvisioningQueueTable({
  rows,
  selectedJobID,
  actionState,
  manualReasons,
  onSelectJob,
  onReasonChange,
  onAction,
}: AdminProvisioningQueueTableProps) {
  return (
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
                      onClick={() => onSelectJob(row.apiId as string)}
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
                    onReasonChange={onReasonChange}
                    onAction={onAction}
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
  );
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
