"use client";

import { FormEvent, useState } from "react";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import type { TopupRequestQuery } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { mapAdminTopupView } from "@/lib/api/viewModels";
import { paymentMethodLabel } from "@/lib/api/walletViewModels";
import { TOPUP_REQUESTS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";
import { AdminFilterBar, AdminFilterInput, AdminFilterSelect } from "../components/AdminFilterBar";
import { TOPUP_STATUS_OPTIONS } from "../lib/filterOptions";
import { equalsFilter, hasActiveFilters, includesFilter, trimStringFilters } from "../lib/filterUtils";

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
type TopupFilterFields = Required<Pick<
  TopupRequestQuery,
  "display_id" | "wallet_display_id" | "requested_by_display_id" | "status"
>>;

const EMPTY_FILTERS: TopupFilterFields = {
  display_id: "",
  wallet_display_id: "",
  requested_by_display_id: "",
  status: "",
};

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

function filterDemoTopups(filters: TopupFilterFields) {
  return TOPUP_REQUESTS.filter((request) => (
    includesFilter(request.id, filters.display_id)
    && includesFilter(request.tenant, filters.wallet_display_id)
    && includesFilter(request.actor, filters.requested_by_display_id)
    && equalsFilter(request.status, filters.status)
  ));
}

export function AdminTopups() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [draftFilters, setDraftFilters] = useState(EMPTY_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(EMPTY_FILTERS);
  const [rejectNotes, setRejectNotes] = useState<Record<string, string>>({});
  const [actionState, setActionState] = useState<ActionState | null>(null);
  const topups = useApiResource(
    () => billingApi.listAdminTopupRequests({ ...appliedFilters, limit: 50 }),
    `admin-topups:${refreshKey}:${JSON.stringify(appliedFilters)}`,
  );
  const usingLive = topups.status === "success";
  const rows: TopupRow[] = usingLive
    ? (topups.data ?? []).map(mapAdminTopupView)
    : filterDemoTopups(appliedFilters).map((req) => ({
        id: req.id,
        live: false,
        tenant: req.tenant,
        actor: req.actor,
        amount: fmtMoney(req.amount),
        method: paymentMethodLabel(req.method),
        ref: req.ref,
        created: req.created,
        proof: req.proof ? "Ref provided" : "No ref",
        status: req.status,
        note: req.reason,
      }));
  const pendingCount = rows.filter((req) => isPendingTopup(req.status)).length;
  const activeFilters = hasActiveFilters(appliedFilters);
  const statusTone = topups.status === "error"
    ? "error"
    : topups.status === "loading"
      ? "loading"
      : usingLive
        ? "success"
        : "default";
  const statusText = topups.status === "error"
    ? "Live API unavailable. Showing demo top-up data for the current filters."
    : topups.status === "loading"
      ? "Refreshing live top-up queue..."
      : usingLive
        ? activeFilters
          ? "Live top-up filters applied."
          : "Live top-up queue"
        : activeFilters
          ? "Filters are applied to demo top-up data."
          : "Demo top-up queue";

  function updateFilter(field: keyof TopupFilterFields, value: string) {
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
        <AdminFilterBar onSubmit={applyFilters} onReset={resetFilters} statusText={statusText} statusTone={statusTone}>
          <AdminFilterInput
            label="Top-up public ID"
            value={draftFilters.display_id}
            onChange={(event) => updateFilter("display_id", event.target.value)}
            placeholder="62001"
            inputMode="numeric"
          />
          <AdminFilterInput
            label="Wallet public ID"
            value={draftFilters.wallet_display_id}
            onChange={(event) => updateFilter("wallet_display_id", event.target.value)}
            placeholder="41001"
            inputMode="numeric"
          />
          <AdminFilterInput
            label="Requester public ID"
            value={draftFilters.requested_by_display_id}
            onChange={(event) => updateFilter("requested_by_display_id", event.target.value)}
            placeholder="10002"
            inputMode="numeric"
          />
          <AdminFilterSelect
            label="Status"
            value={draftFilters.status}
            onChange={(event) => updateFilter("status", event.target.value)}
            options={TOPUP_STATUS_OPTIONS}
          />
        </AdminFilterBar>
        <div className="overflow-x-auto">
          <table className="min-w-[1120px] w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Wallet / Requester", "Amount", "Method", "Ref", "Created", "Proof", "Status", "Actions"].map((h) => (
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
