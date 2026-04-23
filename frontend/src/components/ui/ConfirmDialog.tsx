"use client";

import { useState } from "react";
import { Modal } from "./Modal";

interface ConfirmDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: (reason?: string) => void;
  title: string;
  description: string;
  danger?: boolean;
  confirmLabel?: string;
  /** If true, shows a required "reason" textarea before confirming */
  requireReason?: boolean;
  reasonLabel?: string;
}

export function ConfirmDialog({
  open, onClose, onConfirm,
  title, description,
  danger = false,
  confirmLabel = "Confirm",
  requireReason = false,
  reasonLabel = "Reason",
}: ConfirmDialogProps) {
  const [reason, setReason] = useState("");

  const handleConfirm = () => {
    if (requireReason && !reason.trim()) return;
    onConfirm(reason.trim() || undefined);
    setReason("");
    onClose();
  };

  const handleClose = () => {
    setReason("");
    onClose();
  };

  return (
    <Modal open={open} onClose={handleClose} title={title} width="sm">
      <div className="flex flex-col gap-4">
        <p className="text-[13px] text-gray-600 m-0">{description}</p>

        {requireReason && (
          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">
              {reasonLabel} <span className="text-red-500">*</span>
            </label>
            <textarea
              rows={3}
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Enter reason…"
              className="w-full text-[13px] border border-gray-200 rounded px-3 py-2 resize-none focus:outline-none focus:ring-1 focus:ring-gray-400"
            />
          </div>
        )}

        <div className="flex justify-end gap-2 pt-1">
          <button
            onClick={handleClose}
            className="h-9 px-4 text-[13px] font-medium bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 cursor-pointer transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleConfirm}
            disabled={requireReason && !reason.trim()}
            className={`h-9 px-4 text-[13px] font-medium rounded-md text-white cursor-pointer transition-colors disabled:opacity-40 disabled:cursor-not-allowed
              ${danger ? "bg-red-600 hover:bg-red-700" : "bg-[#D50C2D] hover:bg-[#B3082A]"}`}
          >
            {confirmLabel}
          </button>
        </div>
      </div>
    </Modal>
  );
}
