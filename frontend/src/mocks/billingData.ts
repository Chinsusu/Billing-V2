// Mock billing platform data — aligned with Billing-V2 spec
// No real provider auth material, IPs, or customer data.

export type {
  BandwidthService,
  ClientService,
  Customer,
  ProxyService,
  ProxyType,
  ResellerCatalogItem,
  ResellerClient,
  Service,
  ServiceStatus,
  VpsOS,
  VpsService,
} from "./serviceData";
export {
  BANDWIDTH_SERVICES,
  CLIENT_SERVICES,
  CUSTOMERS,
  PROXY_SERVICES,
  RESELLER_CATALOG,
  RESELLER_CLIENTS,
  SERVICES,
  VPS_SERVICES,
} from "./serviceData";

export type TenantType = "admin" | "reseller";
export type ProviderHealth = "ok" | "degraded" | "down";
export type JobStatus = "queued" | "provisioning" | "failed" | "manual_review";
export type TopupStatus = "pending_verification" | "approved" | "rejected";
export type InvoiceStatus = "open" | "paid" | "overdue";
export type TransactionType = "charge" | "topup" | "refund";

export interface Tenant {
  id: string;
  name: string;
  type: TenantType;
  domain: string;
  clients: number;
  services: number;
  wallet: number;
  walletLow: boolean;
  status: "active" | "suspended";
  since?: string;
}

export interface Provider {
  id: string;
  name: string;
  type: "self-host" | "upstream";
  health: ProviderHealth;
  capacity: number;
  failRate: number;
  lastSync: string;
}

export interface ProvisioningJob {
  id: string;
  order: string;
  service: string;
  tenant: string;
  provider: string;
  status: JobStatus;
  attempt: number;
  error: string;
  correlation: string;
  age: string;
}

export interface TopupRequest {
  id: string;
  tenant: string;
  actor: string;
  amount: number;
  currency: string;
  method: string;
  ref: string;
  status: TopupStatus;
  created: string;
  proof: boolean;
  reason?: string;
}

export interface Invoice {
  id: string;
  customer: string;
  issued: string;
  due: string;
  amount: number;
  status: InvoiceStatus;
}

export interface Transaction {
  id: string;
  time: string;
  customer: string;
  method: string;
  type: TransactionType;
  amount: number;
  status: "paid" | "failed" | "pending";
}

export interface Ticket {
  id: string;
  subject: string;
  customer: string;
  priority: "high" | "medium" | "low";
  status: "open" | "pending" | "closed";
  updated: string;
  assignee: string;
}

export interface ProductCatalog {
  sku: string;
  name: string;
  unit: string;
  price: number;
  active: number;
  rev30: number;
}

export interface ActivityEvent {
  t: string;
  icon: "payment" | "user" | "server" | "ticket" | "error" | "wallet";
  text: string;
  type: "ok" | "warn" | "danger" | "info";
}

export interface LedgerEntry {
  ts: string;
  type: string;
  amount: number;
  ref: string;
  balance: number;
}

// ─── Data ──────────────────────────────────────────────────────

export const TENANTS: Tenant[] = [
  { id: "T-0001", name: "HANetwork (Admin)", type: "admin", domain: "billing.hanetwork.vn", clients: 1284, services: 8420, wallet: 0, walletLow: false, status: "active" },
  { id: "T-0042", name: "ProxyVN Reseller", type: "reseller", domain: "proxyvn.io", clients: 312, services: 1840, wallet: 4820.50, walletLow: false, status: "active", since: "2024-03-12" },
  { id: "T-0051", name: "CloudBase.asia", type: "reseller", domain: "cloudbase.asia", clients: 184, services: 924, wallet: 184.20, walletLow: true, status: "active", since: "2024-07-28" },
  { id: "T-0063", name: "Saigon Proxy", type: "reseller", domain: "proxy.saigon.vn", clients: 94, services: 421, wallet: 2140.80, walletLow: false, status: "active", since: "2025-01-05" },
  { id: "T-0071", name: "DanangHost", type: "reseller", domain: "danang.host", clients: 42, services: 128, wallet: 0, walletLow: true, status: "suspended", since: "2025-03-20" },
];

export const PROVIDERS: Provider[] = [
  { id: "SRC-23001", name: "Proxmox · VN-HCM", type: "self-host", health: "ok", capacity: 82, failRate: 0.2, lastSync: "2m ago" },
  { id: "SRC-23002", name: "Proxmox · VN-HAN", type: "self-host", health: "ok", capacity: 64, failRate: 0.1, lastSync: "1m ago" },
  { id: "SRC-23003", name: "Self-hosted Proxy Manager", type: "self-host", health: "ok", capacity: 91, failRate: 0.3, lastSync: "4m ago" },
  { id: "SRC-23004", name: "OVH", type: "upstream", health: "degraded", capacity: 54, failRate: 2.8, lastSync: "12m ago" },
  { id: "SRC-23005", name: "Hetzner", type: "upstream", health: "ok", capacity: 78, failRate: 0.4, lastSync: "3m ago" },
  { id: "SRC-23006", name: "Smarthost", type: "upstream", health: "ok", capacity: 42, failRate: 0.9, lastSync: "5m ago" },
  { id: "SRC-23007", name: "Budget Proxy Upstream", type: "upstream", health: "ok", capacity: 88, failRate: 0.5, lastSync: "2m ago" },
];

export const PROVISIONING_JOBS: ProvisioningJob[] = [
  { id: "JOB-3301", order: "ORD-48291", service: "VPS 4C/8G · HCM", tenant: "ProxyVN", provider: "Proxmox · VN-HCM", status: "manual_review", attempt: 2, error: "provider_timeout — resource state unknown", correlation: "cor_9a21f8", age: "18m" },
  { id: "JOB-3298", order: "ORD-48290", service: "Residential · 5GB", tenant: "CloudBase", provider: "Budget Proxy Upstream", status: "failed", attempt: 3, error: "auth_failed", correlation: "cor_9a20c4", age: "22m" },
  { id: "JOB-3291", order: "ORD-48289", service: "VPS 2C/4G · SG", tenant: "HANetwork", provider: "OVH", status: "provisioning", attempt: 1, error: "", correlation: "cor_9a19a1", age: "3m" },
  { id: "JOB-3288", order: "ORD-48288", service: "DC Proxy · 100 IPs", tenant: "Saigon Proxy", provider: "Budget Proxy Upstream", status: "provisioning", attempt: 1, error: "", correlation: "cor_9a18e2", age: "1m" },
  { id: "JOB-3287", order: "ORD-48287", service: "VPS 8C/16G · HEL", tenant: "HANetwork", provider: "Hetzner", status: "queued", attempt: 0, error: "", correlation: "cor_9a17b8", age: "0m" },
  { id: "JOB-3286", order: "ORD-48286", service: "ISP Static · 50 IPs", tenant: "ProxyVN", provider: "Smarthost", status: "manual_review", attempt: 1, error: "partial_success — external_id unknown", correlation: "cor_9a16d3", age: "1h 04m" },
];

export const TOPUP_REQUESTS: TopupRequest[] = [
  { id: "TUP-9120", tenant: "ProxyVN", actor: "reseller_wallet", amount: 2000, currency: "USD", method: "VietQR", ref: "FT26042200832", status: "pending_verification", created: "2026-04-22 14:02", proof: true },
  { id: "TUP-9119", tenant: "CloudBase", actor: "reseller_wallet", amount: 500, currency: "USD", method: "USDT", ref: "0x8a…d4e1", status: "pending_verification", created: "2026-04-22 13:40", proof: true },
  { id: "TUP-9118", tenant: "ProxyVN > linh.tran", actor: "client_wallet", amount: 100, currency: "USD", method: "VietQR", ref: "FT26042200781", status: "pending_verification", created: "2026-04-22 13:18", proof: true },
  { id: "TUP-9117", tenant: "HANetwork > kenji.w", actor: "client_wallet", amount: 50, currency: "USD", method: "USDT", ref: "0x1c…9a82", status: "approved", created: "2026-04-22 11:50", proof: true },
  { id: "TUP-9116", tenant: "Saigon Proxy", actor: "reseller_wallet", amount: 1000, currency: "USD", method: "VietQR", ref: "FT26042200412", status: "approved", created: "2026-04-22 10:12", proof: true },
  { id: "TUP-9115", tenant: "CloudBase > huy.nguyen", actor: "client_wallet", amount: 200, currency: "USD", method: "VietQR", ref: "FT26042200388", status: "rejected", created: "2026-04-22 09:40", proof: false, reason: "bank reference not found" },
];

export const INVOICES: Invoice[] = [
  { id: "INV-2026-04218", customer: "Acme Proxy Co.", issued: "2026-04-20", due: "2026-05-04", amount: 4280.00, status: "open" },
  { id: "INV-2026-04217", customer: "DataMine Inc.", issued: "2026-04-20", due: "2026-05-04", amount: 8420.00, status: "paid" },
  { id: "INV-2026-04216", customer: "Kenji Watanabe", issued: "2026-04-15", due: "2026-04-29", amount: 420.00, status: "overdue" },
  { id: "INV-2026-04215", customer: "CloudHarvest", issued: "2026-04-18", due: "2026-05-02", amount: 6240.00, status: "paid" },
  { id: "INV-2026-04214", customer: "Marie Dubois", issued: "2026-04-18", due: "2026-05-02", amount: 1120.00, status: "paid" },
  { id: "INV-2026-04213", customer: "Scrapers Ltd", issued: "2026-04-17", due: "2026-05-01", amount: 1840.00, status: "open" },
  { id: "INV-2026-04212", customer: "Sofia Bergström", issued: "2026-04-15", due: "2026-04-29", amount: 1480.00, status: "paid" },
  { id: "INV-2026-04211", customer: "Linh Tran", issued: "2026-04-14", due: "2026-04-28", amount: 189.00, status: "paid" },
];

export const TRANSACTIONS: Transaction[] = [
  { id: "TX-51001", time: "2026-04-22 14:22", customer: "DataMine Inc.", method: "Visa •• 4242", type: "charge", amount: 8420.00, status: "paid" },
  { id: "TX-51002", time: "2026-04-22 13:48", customer: "CloudHarvest", method: "ACH", type: "charge", amount: 6240.00, status: "paid" },
  { id: "TX-51003", time: "2026-04-22 11:17", customer: "Linh Tran", method: "Wallet", type: "topup", amount: 500.00, status: "paid" },
  { id: "TX-51004", time: "2026-04-22 10:02", customer: "Kenji Watanabe", method: "Visa •• 0914", type: "charge", amount: 420.00, status: "failed" },
  { id: "TX-51005", time: "2026-04-22 09:41", customer: "Marie Dubois", method: "Mastercard •• 1821", type: "charge", amount: 1120.00, status: "paid" },
  { id: "TX-51006", time: "2026-04-22 08:12", customer: "Proxy Garden", method: "PayPal", type: "charge", amount: 340.00, status: "paid" },
  { id: "TX-51007", time: "2026-04-21 22:58", customer: "Acme Proxy Co.", method: "Wire", type: "charge", amount: 4280.00, status: "pending" },
  { id: "TX-51008", time: "2026-04-21 14:21", customer: "Nguyễn Tuấn", method: "Wallet", type: "refund", amount: -29.00, status: "paid" },
];

export const TICKETS: Ticket[] = [
  { id: "T-8124", subject: "Proxy pool authentication failing intermittently", customer: "Acme Proxy Co.", priority: "high", status: "open", updated: "12m ago", assignee: "Linh" },
  { id: "T-8123", subject: "Invoice INV-2026-04216 — requesting extension", customer: "Kenji Watanabe", priority: "medium", status: "pending", updated: "38m ago", assignee: "Minh" },
  { id: "T-8122", subject: "VPS Windows license key issue", customer: "DataMine Inc.", priority: "low", status: "open", updated: "1h ago", assignee: "Tùng" },
  { id: "T-8121", subject: "Bandwidth overage pricing clarification", customer: "Scrapers Ltd", priority: "low", status: "open", updated: "2h ago", assignee: "—" },
  { id: "T-8120", subject: "Request: dedicated IP block allocation", customer: "CloudHarvest", priority: "medium", status: "pending", updated: "3h ago", assignee: "Linh" },
  { id: "T-8119", subject: "Can't access control panel — 2FA lockout", customer: "Linh Tran", priority: "high", status: "open", updated: "4h ago", assignee: "Minh" },
];

export const PRODUCTS: ProductCatalog[] = [
  { sku: "PRX-RES-STD", name: "Residential · Standard", unit: "per GB", price: 6.50, active: 2841, rev30: 124800 },
  { sku: "PRX-RES-PRM", name: "Residential · Premium", unit: "per GB", price: 9.80, active: 1204, rev30: 68200 },
  { sku: "PRX-DC-SHR", name: "Datacenter · Shared", unit: "per IP/mo", price: 0.80, active: 8920, rev30: 52400 },
  { sku: "PRX-DC-DED", name: "Datacenter · Dedicated", unit: "per IP/mo", price: 2.20, active: 1840, rev30: 38200 },
  { sku: "PRX-ISP-STC", name: "ISP Static", unit: "per IP/mo", price: 3.50, active: 612, rev30: 38900 },
  { sku: "PRX-MOB-4G", name: "Mobile 4G · Port", unit: "per port/mo", price: 48.00, active: 268, rev30: 27100 },
  { sku: "VPS-LNX-S", name: "VPS Linux · Small", unit: "per mo", price: 19.00, active: 842, rev30: 18400 },
  { sku: "VPS-LNX-M", name: "VPS Linux · Medium", unit: "per mo", price: 48.00, active: 614, rev30: 34800 },
  { sku: "VPS-LNX-L", name: "VPS Linux · Large", unit: "per mo", price: 129.00, active: 312, rev30: 28600 },
  { sku: "VPS-WIN-M", name: "VPS Windows · Medium", unit: "per mo", price: 78.00, active: 402, rev30: 24800 },
];

export const ACTIVITY_FEED: ActivityEvent[] = [
  { t: "14:22", icon: "payment", text: "Payment of $8,420.00 from DataMine Inc.", type: "ok" },
  { t: "14:11", icon: "user",    text: "New customer signup: startup-dev-42@proton.me", type: "info" },
  { t: "14:02", icon: "server",  text: "VPS service provisioned for Proxy Garden", type: "info" },
  { t: "13:48", icon: "payment", text: "Payment of $6,240.00 from CloudHarvest", type: "ok" },
  { t: "13:17", icon: "ticket",  text: "New ticket T-8124 opened by Acme Proxy Co. (high)", type: "warn" },
  { t: "12:32", icon: "error",   text: "Charge failed: Kenji Watanabe — Visa •• 0914", type: "danger" },
  { t: "11:17", icon: "wallet",  text: "Wallet top-up of $500.00 from Linh Tran", type: "ok" },
];

export const CLIENT_LEDGER: LedgerEntry[] = [
  { ts: "2026-04-22 14:02", type: "purchase.client_wallet.debit", amount: -62.00, ref: "ORD-48290 · Residential EU", balance: 128.40 },
  { ts: "2026-04-21 09:18", type: "topup.credit.client", amount: 100.00, ref: "TUP-9110 · VietQR", balance: 190.40 },
  { ts: "2026-04-20 14:08", type: "renewal.client_wallet.debit", amount: -19.00, ref: "SVC-64421 · VPS Small", balance: 90.40 },
  { ts: "2026-04-18 11:22", type: "purchase.client_wallet.debit", amount: -48.00, ref: "ORD-48280 · Mobile 4G", balance: 109.40 },
  { ts: "2026-04-14 10:02", type: "topup.credit.client", amount: 150.00, ref: "TUP-9088 · VietQR", balance: 157.40 },
];

export const PLATFORM_ALERTS = [
  { text: "3 provisioning jobs in manual review > 1h", type: "danger" as const, screen: "admin-provisioning" },
  { text: "2 reseller tenants below wallet threshold", type: "warn" as const, screen: "admin-topups" },
  { text: "OVH API degraded — 2.8% fail rate", type: "warn" as const, screen: "admin-providers" },
];

export type AlertSeverity = "danger" | "warn" | "info";
export type AlertCategory = "provisioning" | "provider" | "billing" | "security" | "system";

export interface PlatformAlert {
  id: string;
  severity: AlertSeverity;
  category: AlertCategory;
  title: string;
  detail: string;
  screen: string;
  ts: string;
  resolved: boolean;
}

export const ALERTS: PlatformAlert[] = [
  { id: "ALT-001", severity: "danger", category: "provisioning", title: "3 jobs stuck in manual review > 1h", detail: "Jobs JOB-3301, JOB-3298, JOB-3291 have not progressed. Provider timeout on OVH.", screen: "admin-provisioning", ts: "2026-04-22 14:05", resolved: false },
  { id: "ALT-002", severity: "danger", category: "billing",      title: "Charge failed: Kenji Watanabe", detail: "Visa •• 0914 declined. Invoice INV-8821 overdue $82.00. Auto-suspend in 24h.", screen: "admin-invoices", ts: "2026-04-22 12:32", resolved: false },
  { id: "ALT-003", severity: "warn",   category: "provider",     title: "OVH API degraded — 2.8% fail rate", detail: "Error rate above threshold over last 30 min. Provisioning continues but monitored.", screen: "admin-providers", ts: "2026-04-22 11:50", resolved: false },
  { id: "ALT-004", severity: "warn",   category: "billing",      title: "2 reseller tenants below wallet threshold", detail: "DataMine Inc. ($41.20) and Proxy Garden ($18.80) are below $50 floor.", screen: "admin-tenants", ts: "2026-04-22 10:15", resolved: false },
  { id: "ALT-005", severity: "warn",   category: "provisioning", title: "Proxmox source SRC-23001 at 91% memory", detail: "High memory utilisation on source. New VPS provisioning on this source may fail.", screen: "admin-providers", ts: "2026-04-22 09:44", resolved: false },
  { id: "ALT-006", severity: "info",   category: "system",       title: "DB migration 0003 applied successfully", detail: "Migration ran in 1.2s. No rollback needed.", screen: "admin-settings", ts: "2026-04-21 22:01", resolved: true },
  { id: "ALT-007", severity: "info",   category: "security",     title: "New admin login from new IP", detail: "User Minh Nguyen logged in from 103.21.x.x — Vietnam. Session flagged for review.", screen: "admin-settings", ts: "2026-04-21 18:30", resolved: true },
  { id: "ALT-008", severity: "danger", category: "security",     title: "API key rotation overdue — Hetzner", detail: "Provider source SRC-23005 key has not been rotated in 90 days.", screen: "admin-providers", ts: "2026-04-20 08:00", resolved: false },
];

export type AuditLogLevel = "info" | "warn" | "error";
export type AuditActor = "system" | "admin" | "reseller" | "client";

export interface AuditLog {
  id: string;
  ts: string;
  level: AuditLogLevel;
  actor: AuditActor;
  actorName: string;
  action: string;
  target: string;
  detail: string;
  requestId: string;
  tenantId: string;
}

export const AUDIT_LOGS: AuditLog[] = [
  { id: "AUD-70091", ts: "2026-04-22 14:22", level: "info",  actor: "client",   actorName: "Linh Tran",       action: "wallet.topup.submitted",       target: "TUP-9115",    detail: "Amount $200 via VietQR",                    requestId: "Request not shown", tenantId: "T-0042" },
  { id: "AUD-70090", ts: "2026-04-22 14:11", level: "info",  actor: "system",   actorName: "billing-worker",  action: "invoice.auto_charged",         target: "INV-8820",    detail: "Charged $8,420.00 from DataMine Inc.",      requestId: "Request not shown", tenantId: "T-0018" },
  { id: "AUD-70089", ts: "2026-04-22 14:05", level: "error", actor: "system",   actorName: "prov-worker",     action: "provisioning.job.stuck",       target: "JOB-3301",    detail: "manual_review threshold exceeded",          requestId: "Request not shown", tenantId: "T-0031" },
  { id: "AUD-70088", ts: "2026-04-22 13:58", level: "info",  actor: "admin",    actorName: "Minh Nguyen",     action: "tenant.topup.approved",        target: "TUP-9110",    detail: "Approved $500 for ProxyVN (T-0042)",        requestId: "Request not shown", tenantId: "T-0001" },
  { id: "AUD-70087", ts: "2026-04-22 13:44", level: "info",  actor: "reseller", actorName: "ProxyVN",         action: "service.renewed",              target: "SVC-68821",  detail: "Proxy bundle renewed 30d",                  requestId: "Request not shown", tenantId: "T-0042" },
  { id: "AUD-70086", ts: "2026-04-22 13:17", level: "warn",  actor: "client",   actorName: "Acme Proxy Co.", action: "ticket.opened",                target: "T-8124",      detail: "Priority: high. Subject: IP blocked",       requestId: "Request not shown", tenantId: "T-0031" },
  { id: "AUD-70085", ts: "2026-04-22 12:32", level: "error", actor: "system",   actorName: "billing-worker",  action: "invoice.charge.failed",        target: "INV-8821",    detail: "Visa •• 0914 declined for Kenji Watanabe",  requestId: "Request not shown", tenantId: "T-0042" },
  { id: "AUD-70084", ts: "2026-04-22 12:01", level: "info",  actor: "admin",    actorName: "Minh Nguyen",     action: "product.price.updated",        target: "VPS-SMALL",   detail: "$12→$14/mo, effective next renewal",        requestId: "Request not shown", tenantId: "T-0001" },
  { id: "AUD-70083", ts: "2026-04-22 11:50", level: "warn",  actor: "system",   actorName: "health-worker",   action: "provider.health.degraded",     target: "SRC-23004", detail: "Error rate 2.8% over 30min window",         requestId: "Request not shown", tenantId: "T-0001" },
  { id: "AUD-70082", ts: "2026-04-22 11:17", level: "info",  actor: "client",   actorName: "Linh Tran",       action: "wallet.topup.approved",        target: "TUP-9088",    detail: "Credited $500 to client wallet",             requestId: "Request not shown", tenantId: "T-0042" },
  { id: "AUD-70081", ts: "2026-04-22 10:44", level: "info",  actor: "system",   actorName: "prov-worker",     action: "service.provisioned",          target: "SVC-65512",  detail: "Proxy automation VPS active on OVH",        requestId: "Request not shown", tenantId: "T-0031" },
  { id: "AUD-70080", ts: "2026-04-22 10:15", level: "warn",  actor: "system",   actorName: "billing-worker",  action: "tenant.wallet.low_balance",    target: "T-0018",      detail: "DataMine Inc. balance $41.20 below floor",  requestId: "Request not shown", tenantId: "T-0001" },
  { id: "AUD-70079", ts: "2026-04-22 09:44", level: "warn",  actor: "system",   actorName: "health-worker",   action: "provider.node.high_memory",    target: "SRC-23001", detail: "91% memory utilisation on Proxmox node",    requestId: "Request not shown", tenantId: "T-0001" },
  { id: "AUD-70078", ts: "2026-04-22 09:01", level: "info",  actor: "reseller", actorName: "ProxyVN",         action: "catalog.price.updated",        target: "RES-PROX-4G", detail: "Markup adjusted from 35%→40%",              requestId: "Request not shown", tenantId: "T-0042" },
  { id: "AUD-70077", ts: "2026-04-21 22:01", level: "info",  actor: "system",   actorName: "migrator",        action: "db.migration.applied",         target: "0003",        detail: "Migration 0003_rbac ran in 1.2s",           requestId: "Request not shown", tenantId: "T-0001" },
  { id: "AUD-70076", ts: "2026-04-21 18:30", level: "warn",  actor: "admin",    actorName: "Minh Nguyen",     action: "auth.login.new_ip",            target: "session-991", detail: "Login from 103.21.x.x (Vietnam, new IP)",   requestId: "Request not shown", tenantId: "T-0001" },
];
