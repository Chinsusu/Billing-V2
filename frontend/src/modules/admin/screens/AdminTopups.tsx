"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel, shortID } from "@/lib/api/format";
import { useApiResource } from "@/lib/api/useApiResource";
import { TOPUP_REQUESTS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

const PENDING_TOPUP_STATUSES = new Set([
  "pending",
  "pending_verification",
  "submitted",
  "under_review",
  "manual_review",
  "queued",
]);

function isPendingTopup(status: string) {
  return PENDING_TOPUP_STATUSES.has(status);
}

export function AdminTopups() {
  const topups = useApiResource(
    () => billingApi.listAdminTopupRequests({ limit: 50 }),
    "admin-topups:50",
  );
  const usingLive = topups.status === "success";
  const rows = usingLive
    ? (topups.data ?? []).map((req) => ({
        id: recordLabel(req.display_id, "TUP-"),
        tenant: shortID(req.tenant_id),
        actor: shortID(req.requested_by),
        amount: moneyMinor(req.amount_minor, req.currency),
        method: req.payment_method,
        ref: req.payment_reference ?? "-",
        created: compactDateTime(req.created_at),
        proof: req.payment_reference ? "Ref provided" : "No ref",
        status: req.status,
        note: req.review_note,
      }))
    : TOPUP_REQUESTS.map((req) => ({
        id: req.id,
        tenant: req.tenant,
        actor: req.actor,
        amount: fmtMoney(req.amount),
        method: req.method,
        ref: req.ref,
        created: req.created,
        proof: req.proof ? "Ref provided" : "No ref",
        status: req.status,
        note: req.reason,
      }));
  const pendingCount = rows.filter((req) => isPendingTopup(req.status)).length;
  const statusText = topups.status === "error"
    ? "Live API unavailable. Showing demo top-up data until the backend responds."
    : topups.status === "loading"
      ? "Refreshing live top-up queue..."
      : usingLive
        ? "Live top-up queue"
        : "Demo top-up queue";

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Top-up verification queue</h3>
          <div className="flex flex-wrap items-center justify-end gap-3">
            <span className="text-[11px] text-gray-400">{statusText}</span>
            <span className="text-[11px] text-gray-400">{pendingCount} pending</span>
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-[940px] w-full text-[13px] border-collapse">
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
              {rows.map((req) => (
                <tr key={req.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D]">{req.id}</td>
                  <td className="p-4">
                    <div className="text-gray-900">{req.tenant}</div>
                    <div className="text-[11px] text-gray-400">{req.actor}</div>
                  </td>
                  <td className="p-4 tabular-nums font-medium">{req.amount}</td>
                  <td className="p-4 text-gray-500">{req.method}</td>
                  <td className="p-4 text-[12px] text-gray-400">{req.ref}</td>
                  <td className="p-4 text-gray-400">{req.created}</td>
                  <td className="p-4 text-[12px] text-gray-500">{req.proof}</td>
                  <td className="p-4"><StatusBadge status={req.status} dot /></td>
                  <td className="p-4">
                    {isPendingTopup(req.status) ? (
                      <button
                        disabled
                        className="inline-flex h-8 cursor-not-allowed items-center justify-center rounded-md border border-gray-200 bg-gray-50 px-3 text-[12px] font-medium text-gray-400"
                        title="Read-only until a review action is wired to this screen."
                      >
                        Review pending
                      </button>
                    ) : (
                      <span className="text-[11px] text-gray-400">Read-only</span>
                    )}
                    {req.note && <div className="mt-1 text-[11px] text-red-500">{req.note}</div>}
                  </td>
                </tr>
              ))}
              {usingLive && rows.length === 0 && (
                <tr><td colSpan={9} className="p-4 text-center text-[12px] text-gray-400">No top-up requests</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
