"use client";

import { useState } from "react";
import { TOPUP_REQUESTS, type TopupRequest } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { ConfirmDialog } from "@/components/ui/ConfirmDialog";
import { EmptyState } from "@/components/ui/EmptyState";
import { useToast } from "@/lib/toast/ToastContext";
import { fmtMoney } from "@/mocks/sampleData";

export function AdminTopups() {
  const { toast } = useToast();
  const [requests, setRequests] = useState<TopupRequest[]>(TOPUP_REQUESTS);
  const [approving, setApproving] = useState<TopupRequest | null>(null);
  const [rejecting, setRejecting] = useState<TopupRequest | null>(null);

  const handleApprove = () => {
    if (!approving) return;
    setRequests((prev) =>
      prev.map((r) => r.id === approving.id ? { ...r, status: "approved" } : r),
    );
    toast(`Top-up ${approving.id} approved — ${fmtMoney(approving.amount)} credited to ${approving.tenant}`, "success");
    setApproving(null);
  };

  const handleReject = (reason?: string) => {
    if (!rejecting) return;
    setRequests((prev) =>
      prev.map((r) => r.id === rejecting.id ? { ...r, status: "rejected", reason: reason || "Rejected by admin" } : r),
    );
    toast(`Top-up ${rejecting.id} rejected`, "warning");
    setRejecting(null);
  };

  const pending = requests.filter((r) => r.status === "pending_verification").length;

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Top-up verification queue</h3>
          <span className="text-[11px] text-gray-400">{pending} pending</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Tenant / Actor", "Amount", "Method", "Ref", "Created", "Proof", "Status", "Actions"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {requests.length === 0 && <EmptyState title="No top-up requests" />}
            {requests.map((req) => (
              <tr key={req.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 text-[12px] text-[#D50C2D]">{req.id}</td>
                <td className="p-4">
                  <div className="text-gray-900">{req.tenant}</div>
                  <div className="text-[11px] text-gray-400">{req.actor}</div>
                </td>
                <td className="p-4 tabular-nums font-medium">{fmtMoney(req.amount)}</td>
                <td className="p-4 text-gray-500">{req.method}</td>
                <td className="p-4 text-[12px] text-gray-400">{req.ref}</td>
                <td className="p-4 text-gray-400">{req.created}</td>
                <td className="p-4 text-center">{req.proof ? "✓" : <span className="text-red-500">✕</span>}</td>
                <td className="p-4"><StatusBadge status={req.status} dot /></td>
                <td className="p-4">
                  {req.status === "pending_verification" && (
                    <div className="flex gap-1.5">
                      <button
                        onClick={() => setApproving(req)}
                        className="inline-flex items-center justify-center px-3 h-8 text-[12px] font-medium bg-emerald-600 hover:bg-emerald-700 text-white rounded border-0 cursor-pointer transition-colors"
                      >
                        Approve
                      </button>
                      <button
                        onClick={() => setRejecting(req)}
                        className="inline-flex items-center justify-center px-3 h-8 text-[12px] font-medium bg-white hover:bg-red-50 text-red-600 border border-red-200 rounded cursor-pointer transition-colors"
                      >
                        Reject
                      </button>
                    </div>
                  )}
                  {req.reason && <span className="text-[11px] text-red-500">{req.reason}</span>}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <ConfirmDialog
        open={!!approving}
        onClose={() => setApproving(null)}
        onConfirm={handleApprove}
        title="Approve top-up"
        description={approving ? `Approve ${fmtMoney(approving.amount)} top-up from ${approving.tenant}? This will credit the wallet immediately.` : ""}
        confirmLabel="Approve"
      />

      <ConfirmDialog
        open={!!rejecting}
        onClose={() => setRejecting(null)}
        onConfirm={handleReject}
        title="Reject top-up"
        description={rejecting ? `Reject top-up request ${rejecting.id} from ${rejecting.tenant}?` : ""}
        danger
        confirmLabel="Reject"
        requireReason
        reasonLabel="Rejection reason"
      />
    </div>
  );
}
