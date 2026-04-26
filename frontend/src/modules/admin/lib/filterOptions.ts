export type AdminFilterOption = {
  value: string;
  label: string;
};

export const ACCOUNT_TYPE_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All types" },
  { value: "admin", label: "Admin" },
  { value: "reseller", label: "Reseller" },
  { value: "client", label: "Client" },
  { value: "user", label: "User" },
  { value: "worker", label: "Worker" },
  { value: "provider_webhook", label: "Provider webhook" },
];

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

export const PRODUCT_TYPE_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All products" },
  { value: "vps", label: "VPS" },
  { value: "proxy", label: "Proxy" },
  { value: "bandwidth", label: "Bandwidth" },
];

export const PROVIDER_SOURCE_TYPE_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All source types" },
  { value: "hetzner", label: "Hetzner" },
  { value: "proxmox", label: "Proxmox" },
  { value: "manual", label: "Manual pool" },
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

export const TENANT_TYPE_OPTIONS: AdminFilterOption[] = [
  { value: "", label: "All types" },
  { value: "admin", label: "Admin" },
  { value: "reseller", label: "Reseller" },
  { value: "client", label: "Client" },
];
