import { compactDateTime, moneyMinor, recordLabel } from "./format";
import {
  accountTypeLabel,
  inventoryModeLabel,
  providerSourceTypeLabel,
  riskLevelLabel,
  securityStatusLabel,
} from "./displayLabels";
import { paymentMethodLabel, paymentTransactionTypeLabel } from "./paymentViewModels";
import type {
  AdminAccount,
  AuditLog,
  CatalogProviderSource,
  Invoice,
  PaymentReconciliation,
  PaymentTransaction,
  TopupRequest,
} from "./types";

export function adminDisplayLabel(displayID: number | string, prefix: string): string {
  return recordLabel(displayID, prefix);
}

export function hiddenReference(label = "Reference"): string {
  return label ? `${label} not shown` : "Not shown";
}

function publicIDLabel(displayID: number | undefined, prefix: string, fallback: string): string {
  return displayID ? adminDisplayLabel(displayID, prefix) : hiddenReference(fallback);
}

function auditTargetLabel(log: AuditLog): string {
  const targetLabels: Record<string, string> = {
    invoice: "Invoice",
    job: "Job",
    order: "Order",
    provider: "Provider",
    provider_source: "Provider",
    service: "Service",
    service_instance: "Service",
    topup_request: "Top-up",
  };
  const targetLabel = targetLabels[log.target_type] ?? "Target";
  if (!log.target_display_id) return `${targetLabel} not shown`;
  const prefixes: Record<string, string> = {
    invoice: "INV-",
    job: "JOB-",
    order: "ORD-",
    provider: "SRC-",
    provider_source: "SRC-",
    service: "SVC-",
    service_instance: "SVC-",
    topup_request: "TUP-",
  };
  return `${targetLabel} ${adminDisplayLabel(log.target_display_id, prefixes[log.target_type] ?? "#")}`;
}

export function requestLabel(): string {
  return hiddenReference("Request");
}

export interface AdminAccountView {
  id: string;
  name: string;
  email: string;
  type: string;
  tenant: string;
  security: string;
  status: string;
  created: string;
  lastLogin: string;
}

export function mapAdminAccountView(account: AdminAccount): AdminAccountView {
  return {
    id: adminDisplayLabel(account.display_id, "ACC-"),
    name: account.full_name || account.email,
    email: account.email,
    type: accountTypeLabel(account.user_type),
    tenant: account.tenant_name || account.tenant_slug,
    security: securityStatusLabel(account.two_factor_status),
    status: account.status,
    created: compactDateTime(account.created_at),
    lastLogin: compactDateTime(account.last_login_at),
  };
}

export interface AdminProviderSourceView {
  id: string;
  name: string;
  type: string;
  status: string;
  location: string;
  inventory: string;
  risk: string;
  account: string;
  updated: string;
}

export function providerSourceLabel(provider: CatalogProviderSource): string {
  return `${provider.name} (${adminDisplayLabel(provider.display_id, "SRC-")})`;
}

export function mapAdminProviderSourceView(provider: CatalogProviderSource): AdminProviderSourceView {
  return {
    id: adminDisplayLabel(provider.display_id, "SRC-"),
    name: provider.name,
    type: providerSourceTypeLabel(provider.source_type),
    status: provider.status,
    location: provider.location || "-",
    inventory: inventoryModeLabel(provider.inventory_mode),
    risk: riskLevelLabel(provider.risk_level),
    account: hiddenReference("Account"),
    updated: compactDateTime(provider.updated_at),
  };
}

export interface AdminInvoiceView {
  id: string;
  customer: string;
  order: string;
  issued: string;
  due: string;
  amount: string;
  status: string;
}

export function mapAdminInvoiceView(invoice: Invoice): AdminInvoiceView {
  return {
    id: adminDisplayLabel(invoice.display_id, "INV-"),
    customer: publicIDLabel(invoice.buyer_display_id, "ACC-", "Account"),
    order: publicIDLabel(invoice.order_display_id, "ORD-", "Order"),
    issued: compactDateTime(invoice.issued_at),
    due: compactDateTime(invoice.due_at),
    amount: moneyMinor(invoice.total_minor, invoice.currency),
    status: invoice.status,
  };
}

export interface AdminTransactionView {
  id: string;
  time: string;
  customer: string;
  order: string;
  invoice: string;
  method: string;
  type: string;
  amount: string;
  status: string;
}

export function mapAdminTransactionView(
  transaction: PaymentTransaction,
  reconciliation?: PaymentReconciliation,
): AdminTransactionView {
  return {
    id: adminDisplayLabel(transaction.display_id, "TX-"),
    time: compactDateTime(transaction.created_at),
    customer: publicIDLabel(transaction.account_display_id, "ACC-", "Account"),
    order: publicIDLabel(transaction.order_display_id, "ORD-", "Order"),
    invoice: publicIDLabel(transaction.invoice_display_id, "INV-", "Invoice"),
    method: paymentMethodLabel(reconciliation?.provider ?? "wallet"),
    type: paymentTransactionTypeLabel(transaction.type),
    amount: moneyMinor(transaction.amount_minor, transaction.currency),
    status: transaction.status,
  };
}

export type AdminAuditActorBadge =
  | "system"
  | "admin"
  | "reseller"
  | "client"
  | "user"
  | "worker"
  | "provider_webhook";

export interface AdminAuditLogView {
  id: string;
  ts: string;
  level: "info";
  actor: AdminAuditActorBadge;
  actorName: string;
  action: string;
  target: string;
  detail: string;
  requestId: string;
}

export function mapAdminActorBadge(actorType: string): AdminAuditActorBadge {
  switch (actorType) {
    case "admin":
    case "reseller":
    case "client":
    case "user":
    case "worker":
    case "provider_webhook":
      return actorType;
    default:
      return "system";
  }
}

export function mapAdminAuditLogView(log: AuditLog): AdminAuditLogView {
  return {
    id: adminDisplayLabel(log.display_id, "AUD-"),
    ts: compactDateTime(log.created_at),
    level: "info",
    actor: mapAdminActorBadge(log.actor_type),
    actorName: publicIDLabel(log.actor_display_id, "ACC-", "Actor"),
    action: log.action,
    target: auditTargetLabel(log),
    detail: log.target_display_id ? "Public target linked" : hiddenReference("Target"),
    requestId: requestLabel(),
  };
}

export interface AdminTopupView {
  id: string;
  apiId: string;
  live: true;
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

export function mapAdminTopupView(request: TopupRequest): AdminTopupView {
  return {
    id: adminDisplayLabel(request.display_id, "TUP-"),
    apiId: request.id,
    live: true,
    tenant: publicIDLabel(request.wallet_display_id, "WAL-", "Wallet"),
    actor: publicIDLabel(request.requested_by_display_id, "ACC-", "Requester"),
    amount: moneyMinor(request.amount_minor, request.currency),
    method: paymentMethodLabel(request.payment_method),
    ref: request.payment_reference ?? "-",
    created: compactDateTime(request.created_at),
    proof: request.payment_reference ? "Ref provided" : "No ref",
    status: request.status,
    note: request.review_note,
  };
}
