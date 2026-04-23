"use client";

import { useState } from "react";
import { PROVISIONING_JOBS, type ProvisioningJob } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { ConfirmDialog } from "@/components/ui/ConfirmDialog";
import { EmptyState } from "@/components/ui/EmptyState";
import { useToast } from "@/lib/toast/ToastContext";

export function AdminProvisioning() {
  const { toast } = useToast();
  const [jobs, setJobs] = useState<ProvisioningJob[]>(PROVISIONING_JOBS);
  const [retrying, setRetrying] = useState<ProvisioningJob | null>(null);
  const [cancelling, setCancelling] = useState<ProvisioningJob | null>(null);

  const manualReview = jobs.filter((j) => j.status === "manual_review");

  const handleRetry = () => {
    if (!retrying) return;
    setJobs((prev) =>
      prev.map((j) => j.id === retrying.id ? { ...j, status: "queued", attempt: j.attempt + 1, error: "" } : j),
    );
    toast(`Job ${retrying.id} re-queued for retry`, "info");
    setRetrying(null);
  };

  const handleCancel = (reason?: string) => {
    if (!cancelling) return;
    setJobs((prev) => prev.filter((j) => j.id !== cancelling.id));
    toast(`Job ${cancelling.id} cancelled${reason ? ` — ${reason}` : ""}`, "warning");
    setCancelling(null);
  };

  return (
    <div className="p-4 flex flex-col gap-4">
      {manualReview.length > 0 && (
        <div className="bg-amber-50 border border-amber-200 text-amber-700 text-[12px] p-3.5 rounded flex items-center gap-3">
          <span>⚠</span>
          <span>{manualReview.length} job(s) require manual review — provider state unknown. Verify before re-triggering.</span>
        </div>
      )}

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Provisioning queue</h3>
          <span className="text-[11px] text-gray-400">{jobs.length} jobs</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Job ID", "Order", "Service", "Tenant", "Provider", "Status", "Attempt", "Age", "Error", "Actions"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {jobs.length === 0 && <EmptyState title="No provisioning jobs" />}
            {jobs.map((job) => (
              <tr key={job.id} className={`hover:bg-gray-50 border-b border-gray-100 last:border-0 ${job.status === "manual_review" ? "bg-amber-50/40" : ""}`}>
                <td className="p-4 text-[12px] text-gray-500">{job.id}</td>
                <td className="p-4 text-[12px] text-[#D50C2D]">{job.order}</td>
                <td className="p-4 text-gray-700">{job.service}</td>
                <td className="p-4 text-gray-500">{job.tenant}</td>
                <td className="p-4 text-gray-500">{job.provider}</td>
                <td className="p-4"><StatusBadge status={job.status} dot /></td>
                <td className="p-4 text-center tabular-nums">{job.attempt}</td>
                <td className="p-4 text-gray-400 tabular-nums">{job.age}</td>
                <td className="p-4 text-[11px] text-red-600 max-w-[160px] truncate">{job.error || "—"}</td>
                <td className="p-4">
                  {(job.status === "failed" || job.status === "manual_review") && (
                    <div className="flex gap-1.5">
                      <button
                        onClick={() => setRetrying(job)}
                        className="inline-flex items-center px-3 h-8 text-[12px] font-medium bg-white hover:bg-blue-50 text-blue-600 border border-blue-200 rounded cursor-pointer transition-colors"
                      >
                        Re-trigger
                      </button>
                      <button
                        onClick={() => setCancelling(job)}
                        className="inline-flex items-center px-3 h-8 text-[12px] font-medium bg-white hover:bg-red-50 text-red-600 border border-red-200 rounded cursor-pointer transition-colors"
                      >
                        Cancel
                      </button>
                    </div>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <ConfirmDialog
        open={!!retrying}
        onClose={() => setRetrying(null)}
        onConfirm={handleRetry}
        title="Re-trigger provisioning job"
        description={retrying ? `Re-trigger job ${retrying.id} for order ${retrying.order}? Verify the provider state is clean before proceeding.` : ""}
        confirmLabel="Re-trigger"
      />

      <ConfirmDialog
        open={!!cancelling}
        onClose={() => setCancelling(null)}
        onConfirm={handleCancel}
        title="Cancel provisioning job"
        description={cancelling ? `Cancel job ${cancelling.id}? This will remove it from the queue.` : ""}
        danger
        confirmLabel="Cancel job"
        requireReason
        reasonLabel="Cancellation reason"
      />
    </div>
  );
}
