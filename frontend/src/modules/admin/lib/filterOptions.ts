export type AdminFilterOption = {
  value: string;
  label: string;
};

export const ACCOUNT_STATUS_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All statuses" },
  { value: "active", label: "Active" },
  { value: "overdue", label: "Overdue" },
  { value: "suspended", label: "Suspended" },
];

export const INVOICE_STATUS_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All statuses" },
  { value: "open", label: "Open" },
  { value: "paid", label: "Paid" },
  { value: "overdue", label: "Overdue" },
];

export const JOB_STATUS_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All statuses" },
  { value: "queued", label: "Queued" },
  { value: "claimed", label: "Claimed" },
  { value: "running", label: "Running" },
  { value: "manual_review", label: "Manual Review" },
  { value: "failed_retryable", label: "Retryable" },
  { value: "failed_terminal", label: "Terminal Failed" },
  { value: "cancelled", label: "Cancelled" },
];

export const PROVIDER_STATUS_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All statuses" },
  { value: "active", label: "Active" },
  { value: "pending", label: "Pending" },
  { value: "failed", label: "Failed" },
];

export const SERVICE_STATUS_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All statuses" },
  { value: "active", label: "Active" },
  { value: "overdue", label: "Overdue" },
  { value: "suspended", label: "Suspended" },
  { value: "cancelled", label: "Cancelled" },
];

export const TOPUP_STATUS_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All statuses" },
  { value: "submitted", label: "Submitted" },
  { value: "under_review", label: "Under review" },
  { value: "pending_verification", label: "Pending verification" },
  { value: "approved", label: "Approved" },
  { value: "rejected", label: "Rejected" },
];

export const TRANSACTION_STATUS_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All statuses" },
  { value: "posted", label: "Posted" },
  { value: "paid", label: "Paid" },
  { value: "pending", label: "Pending" },
  { value: "failed", label: "Failed" },
  { value: "refunded", label: "Refunded" },
];
