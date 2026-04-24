"use client";

import { useState } from "react";
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

const REVIEWABLE_TOPUP_STATUSES = new Set(["submitted", "under_review"]);

type TopupAction = "approve" | "reject";

interface TopupRow {
  id: string;
  apiId?: string;
  live: boolean;
  tenant: string;
  actor: string;
  amount: string;
  method: string;
  ref: string;
  created: string;
  proof: string;
  status: string;
  note?: string;
}

interface ActionState {
  id: string;
  action: TopupAction;
  status: "running" | "success" | "error";
  message: string;
}

function isPendingTopup(status: string) {
  return PENDING_TOPUP_STATUSES.has(status);
}

function isReviewableTopup(status: string) {
  return REVIEWABLE_TOPUP_STATUSES.has(status);
}

function errorMessage(error: unknown) {
  return error instanceof Error ? error.message : "Top-up review failed.";
}

export function AdminTopups() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [rejectNotes, setRejectNotes] = useState<Record<string, string>>({});
  const [actionState, setActionState] = useState<ActionState | null>(null);
  const topups = useApiResource(
    () => billingApi.listAdminTopupRequests({ limit: 50 }),
    `admin-topups:${refreshKey}`,
  );
  const usingLive = topups.status === "success";
  const rows: TopupRow[] = usingLive
    ? (topups.data ?? []).map((req) => ({
        id: recordLabel(req.display_id, "TUP-"),
        apiId: req.id,
        live: true,
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
        live: false,
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

  function updateRejectNote(id: string, value: string) {
    setRejectNotes((current) => ({ ...current, [id]: value }));
  }

  async function reviewTopup(row: TopupRow, action: TopupAction) {
    if (!row.apiId) return;
    const note = rejectNotes[row.apiId]?.trim() ?? "";
    if (action === "reject" && !note) {
      setActionState({
        id: row.apiId,
        action,
        status: "error",
        message: "Add a rejection note first.",
      });
      return;
    }

    setActionState({
      id: row.apiId,
      action,
      status: "running",
      message: action === "approve" ? "Approving..." : "Rejecting...",
    });

    try {
      if (action === "approve") {
        await billingApi.approveAdminTopupRequest(row.apiId, { review_note: "Approved from admin queue." });
      } else {
        await billingApi.rejectAdminTopupRequest(row.apiId, { review_note: note });
      }
      setActionState({
        id: row.apiId,
        action,
        status: "success",
        message: action === "approve" ? "Approved. Refreshing queue..." : "Rejected. Refreshing queue...",
      });
      if (action === "reject") {
        setRejectNotes((current) => ({ ...current, [row.apiId as string]: "" }));
      }
      setRefreshKey((current) => current + 1);
    } catch (error: unknown) {
      setActionState({
        id: row.apiId,
        action,
        status: "error",
        message: errorMessage(error),
      });
    }
  }

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
          <table className="min-w-[1120px] w-full text-[13px] border-collapse">
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
              {rows.map((req) => {
                const rowState = req.apiId && actionState?.id === req.apiId ? actionState : null;
                const running = rowState?.status === "running";
                const reviewApiId = req.live && req.apiId && isReviewableTopup(req.status) ? req.apiId : "";

                return (
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
                      {reviewApiId ? (
                        <div className="flex min-w-[260px] flex-col gap-2">
                          <input
                            value={rejectNotes[reviewApiId] ?? ""}
                            onChange={(event) => updateRejectNote(reviewApiId, event.target.value)}
                            disabled={running}
                            placeholder="Reject note"
                            className="h-8 rounded-md border border-gray-200 px-2 text-[12px] text-gray-700 outline-none focus:border-[#D50C2D]"
                          />
                          <div className="flex gap-2">
                            <button
                              disabled={running}
                              onClick={() => reviewTopup(req, "approve")}
                              className="inline-flex h-8 items-center justify-center rounded-md border border-emerald-600 bg-emerald-600 px-3 text-[12px] font-medium text-white disabled:cursor-not-allowed disabled:opacity-60"
                            >
                              Approve
                            </button>
                            <button
                              disabled={running}
                              onClick={() => reviewTopup(req, "reject")}
                              className="inline-flex h-8 items-center justify-center rounded-md border border-red-200 bg-white px-3 text-[12px] font-medium text-red-600 disabled:cursor-not-allowed disabled:opacity-60"
                            >
                              Reject
                            </button>
                          </div>
                          {rowState && (
                            <div className={`text-[11px] ${rowState.status === "error" ? "text-red-500" : "text-gray-400"}`}>
                              {rowState.message}
                            </div>
                          )}
                        </div>
                      ) : (
                        <span className="text-[11px] text-gray-400">
                          {req.live ? "Read-only" : "Demo read-only"}
                        </span>
                      )}
                      {req.note && <div className="mt-1 text-[11px] text-red-500">{req.note}</div>}
                    </td>
                  </tr>
                );
              })}
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
