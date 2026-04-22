// Mock billing platform data — aligned with Billing-V2 spec
// No real credentials, IPs, or customer data.

export type TenantType = "admin" | "reseller";
export type ServiceStatus = "active" | "suspended" | "provisioning" | "stopped" | "overdue";
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

export interface Customer {
  id: string;
  name: string;
  email: string;
  plan: string;
  services: number;
  mrr: number;
  status: "active" | "overdue" | "suspended";
  since: string;
  country: string;
}

export interface Service {
  id: string;
  type: string;
  label: string;
  customer: string;
  region: string;
  bandwidth: string;
  price: number;
  status: ServiceStatus;
  renewsIn: number;
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
  icon: string;
  text: string;
  type: "ok" | "warn" | "danger" | "info";
}

export interface ResellerClient {
  id: string;
  name: string;
  email: string;
  wallet: number;
  services: number;
  orders: number;
  status: "active" | "overdue";
  lastLogin: string;
}

export interface ResellerCatalogItem {
  plan: string;
  unit: string;
  cost: number;
  selling: number;
  margin: number;
  stock: "ok" | "low" | "out";
  version: string;
  status: "active" | "disabled" | "warn";
}

export interface ClientService {
  id: string;
  type: string;
  label: string;
  identifier: string;
  region: string;
  bandwidth: string;
  expiry: string;
  status: ServiceStatus;
  cycle: string;
  note?: string;
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
  { id: "prv-pmx-01", name: "Proxmox · VN-HCM", type: "self-host", health: "ok", capacity: 82, failRate: 0.2, lastSync: "2m ago" },
  { id: "prv-pmx-02", name: "Proxmox · VN-HAN", type: "self-host", health: "ok", capacity: 64, failRate: 0.1, lastSync: "1m ago" },
  { id: "prv-pmg", name: "proxy-manager (self)", type: "self-host", health: "ok", capacity: 91, failRate: 0.3, lastSync: "4m ago" },
  { id: "prv-ovh", name: "OVH", type: "upstream", health: "degraded", capacity: 54, failRate: 2.8, lastSync: "12m ago" },
  { id: "prv-hzn", name: "Hetzner", type: "upstream", health: "ok", capacity: 78, failRate: 0.4, lastSync: "3m ago" },
  { id: "prv-smh", name: "Smarthost", type: "upstream", health: "ok", capacity: 42, failRate: 0.9, lastSync: "5m ago" },
  { id: "prv-pch", name: "proxy-cheap", type: "upstream", health: "ok", capacity: 88, failRate: 0.5, lastSync: "2m ago" },
];

export const PROVISIONING_JOBS: ProvisioningJob[] = [
  { id: "job-8a21", order: "ORD-48291", service: "VPS 4C/8G · HCM", tenant: "ProxyVN", provider: "Proxmox · VN-HCM", status: "manual_review", attempt: 2, error: "provider_timeout — resource state unknown", correlation: "cor_9a21f8", age: "18m" },
  { id: "job-8a20", order: "ORD-48290", service: "Residential · 5GB", tenant: "CloudBase", provider: "proxy-cheap", status: "failed", attempt: 3, error: "auth_failed", correlation: "cor_9a20c4", age: "22m" },
  { id: "job-8a19", order: "ORD-48289", service: "VPS 2C/4G · SG", tenant: "HANetwork", provider: "OVH", status: "provisioning", attempt: 1, error: "", correlation: "cor_9a19a1", age: "3m" },
  { id: "job-8a18", order: "ORD-48288", service: "DC Proxy · 100 IPs", tenant: "Saigon Proxy", provider: "proxy-cheap", status: "provisioning", attempt: 1, error: "", correlation: "cor_9a18e2", age: "1m" },
  { id: "job-8a17", order: "ORD-48287", service: "VPS 8C/16G · HEL", tenant: "HANetwork", provider: "Hetzner", status: "queued", attempt: 0, error: "", correlation: "cor_9a17b8", age: "0m" },
  { id: "job-8a16", order: "ORD-48286", service: "ISP Static · 50 IPs", tenant: "ProxyVN", provider: "Smarthost", status: "manual_review", attempt: 1, error: "partial_success — external_id unknown", correlation: "cor_9a16d3", age: "1h 04m" },
];

export const TOPUP_REQUESTS: TopupRequest[] = [
  { id: "TUP-9120", tenant: "ProxyVN", actor: "reseller_wallet", amount: 2000, currency: "USD", method: "VietQR", ref: "FT26042200832", status: "pending_verification", created: "2026-04-22 14:02", proof: true },
  { id: "TUP-9119", tenant: "CloudBase", actor: "reseller_wallet", amount: 500, currency: "USD", method: "USDT", ref: "0x8a…d4e1", status: "pending_verification", created: "2026-04-22 13:40", proof: true },
  { id: "TUP-9118", tenant: "ProxyVN > linh.tran", actor: "client_wallet", amount: 100, currency: "USD", method: "VietQR", ref: "FT26042200781", status: "pending_verification", created: "2026-04-22 13:18", proof: true },
  { id: "TUP-9117", tenant: "HANetwork > kenji.w", actor: "client_wallet", amount: 50, currency: "USD", method: "USDT", ref: "0x1c…9a82", status: "approved", created: "2026-04-22 11:50", proof: true },
  { id: "TUP-9116", tenant: "Saigon Proxy", actor: "reseller_wallet", amount: 1000, currency: "USD", method: "VietQR", ref: "FT26042200412", status: "approved", created: "2026-04-22 10:12", proof: true },
  { id: "TUP-9115", tenant: "CloudBase > huy.nguyen", actor: "client_wallet", amount: 200, currency: "USD", method: "VietQR", ref: "FT26042200388", status: "rejected", created: "2026-04-22 09:40", proof: false, reason: "bank reference not found" },
];

export const CUSTOMERS: Customer[] = [
  { id: "C-40218", name: "Acme Proxy Co.", email: "ops@acmeproxy.io", plan: "Enterprise", services: 48, mrr: 4280, status: "active", since: "2023-04-12", country: "US" },
  { id: "C-40217", name: "Linh Tran", email: "linh.tran@gmail.com", plan: "Pro", services: 7, mrr: 189, status: "active", since: "2024-09-03", country: "VN" },
  { id: "C-40216", name: "Scrapers Ltd", email: "billing@scrapers.ltd", plan: "Business", services: 24, mrr: 1840, status: "active", since: "2023-11-28", country: "GB" },
  { id: "C-40215", name: "Kenji Watanabe", email: "kenji@tokyonet.jp", plan: "Pro", services: 12, mrr: 420, status: "overdue", since: "2024-02-14", country: "JP" },
  { id: "C-40214", name: "Marie Dubois", email: "marie@duboisco.fr", plan: "Business", services: 16, mrr: 1120, status: "active", since: "2023-06-21", country: "FR" },
  { id: "C-40213", name: "DataMine Inc.", email: "accounts@datamine.io", plan: "Enterprise", services: 82, mrr: 8420, status: "active", since: "2022-12-01", country: "US" },
  { id: "C-40212", name: "Hans Müller", email: "h.mueller@web.de", plan: "Starter", services: 3, mrr: 49, status: "suspended", since: "2025-01-18", country: "DE" },
  { id: "C-40211", name: "Proxy Garden", email: "hi@proxygarden.co", plan: "Pro", services: 11, mrr: 340, status: "active", since: "2024-07-09", country: "CA" },
  { id: "C-40210", name: "Alex Rodriguez", email: "alex@rodriguez.mx", plan: "Pro", services: 9, mrr: 278, status: "active", since: "2024-11-22", country: "MX" },
  { id: "C-40209", name: "Sofia Bergström", email: "sofia@bergnordic.se", plan: "Business", services: 21, mrr: 1480, status: "active", since: "2023-08-15", country: "SE" },
  { id: "C-40208", name: "CloudHarvest", email: "billing@cloudharvest.ai", plan: "Enterprise", services: 64, mrr: 6240, status: "active", since: "2023-02-28", country: "US" },
  { id: "C-40207", name: "Nguyễn Tuấn", email: "tuan.nguyen@startup.vn", plan: "Starter", services: 2, mrr: 29, status: "active", since: "2025-03-14", country: "VN" },
];

export const SERVICES: Service[] = [
  { id: "prx-9a12", type: "residential", label: "US Residential Pool", customer: "Acme Proxy Co.", region: "US-EAST", bandwidth: "2.4 TB", price: 380, status: "active", renewsIn: 12 },
  { id: "vps-7c8f", type: "vps-linux", label: "vps-prod-01 · 8 vCPU · 32GB", customer: "DataMine Inc.", region: "EU-HEL", bandwidth: "—", price: 89, status: "active", renewsIn: 8 },
  { id: "prx-4e21", type: "datacenter", label: "DC Pool Alpha · 500 IPs", customer: "Scrapers Ltd", region: "US-WEST", bandwidth: "1.1 TB", price: 240, status: "active", renewsIn: 22 },
  { id: "vps-3d9a", type: "vps-win", label: "win-rdp-gamma · 4 vCPU · 16GB", customer: "Kenji Watanabe", region: "APAC-TYO", bandwidth: "—", price: 64, status: "overdue", renewsIn: -3 },
  { id: "prx-8b67", type: "mobile", label: "Mobile 4G · 20 ports", customer: "CloudHarvest", region: "US-EAST", bandwidth: "842 GB", price: 620, status: "active", renewsIn: 5 },
  { id: "prx-2f01", type: "isp", label: "ISP Static · 100 IPs", customer: "Marie Dubois", region: "EU-FRA", bandwidth: "680 GB", price: 180, status: "active", renewsIn: 18 },
  { id: "vps-6a11", type: "vps-linux", label: "vps-scrape-02 · 2 vCPU · 4GB", customer: "Proxy Garden", region: "EU-HEL", bandwidth: "—", price: 19, status: "provisioning", renewsIn: 30 },
  { id: "prx-5c23", type: "residential", label: "EU Residential · Premium", customer: "Linh Tran", region: "EU-MULTI", bandwidth: "128 GB", price: 62, status: "active", renewsIn: 14 },
  { id: "prx-1a88", type: "datacenter", label: "DC Pool Beta · 1000 IPs", customer: "Acme Proxy Co.", region: "GLOBAL", bandwidth: "3.2 TB", price: 440, status: "active", renewsIn: 9 },
  { id: "prx-7b34", type: "residential", label: "APAC Residential", customer: "Alex Rodriguez", region: "APAC-SIN", bandwidth: "184 GB", price: 118, status: "suspended", renewsIn: -7 },
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
  { id: "txn_9A8f21", time: "2026-04-22 14:22", customer: "DataMine Inc.", method: "Visa •• 4242", type: "charge", amount: 8420.00, status: "paid" },
  { id: "txn_9A8e12", time: "2026-04-22 13:48", customer: "CloudHarvest", method: "ACH", type: "charge", amount: 6240.00, status: "paid" },
  { id: "txn_9A8d04", time: "2026-04-22 11:17", customer: "Linh Tran", method: "Wallet", type: "topup", amount: 500.00, status: "paid" },
  { id: "txn_9A8c92", time: "2026-04-22 10:02", customer: "Kenji Watanabe", method: "Visa •• 0914", type: "charge", amount: 420.00, status: "failed" },
  { id: "txn_9A8b77", time: "2026-04-22 09:41", customer: "Marie Dubois", method: "Mastercard •• 1821", type: "charge", amount: 1120.00, status: "paid" },
  { id: "txn_9A8a31", time: "2026-04-22 08:12", customer: "Proxy Garden", method: "PayPal", type: "charge", amount: 340.00, status: "paid" },
  { id: "txn_9A8920", time: "2026-04-21 22:58", customer: "Acme Proxy Co.", method: "Wire", type: "charge", amount: 4280.00, status: "pending" },
  { id: "txn_9A86b5", time: "2026-04-21 14:21", customer: "Nguyễn Tuấn", method: "Wallet", type: "refund", amount: -29.00, status: "paid" },
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
  { t: "14:22", icon: "💳", text: "Payment of $8,420.00 from DataMine Inc.", type: "ok" },
  { t: "14:11", icon: "👤", text: "New customer signup: startup-dev-42@proton.me", type: "info" },
  { t: "14:02", icon: "🖥", text: "VPS vps-scrape-02 provisioned for Proxy Garden", type: "info" },
  { t: "13:48", icon: "💳", text: "Payment of $6,240.00 from CloudHarvest", type: "ok" },
  { t: "13:17", icon: "🎫", text: "New ticket T-8124 opened by Acme Proxy Co. (high)", type: "warn" },
  { t: "12:32", icon: "✕", text: "Charge failed: Kenji Watanabe — Visa •• 0914", type: "danger" },
  { t: "11:17", icon: "💰", text: "Wallet top-up of $500.00 from Linh Tran", type: "ok" },
];

export const RESELLER_CLIENTS: ResellerClient[] = [
  { id: "RC-2021", name: "Linh Tran", email: "linh@scrape.dev", wallet: 128.40, services: 7, orders: 24, status: "active", lastLogin: "2h ago" },
  { id: "RC-2020", name: "Hùng Phạm", email: "hung@adsvn.co", wallet: 420.00, services: 12, orders: 48, status: "active", lastLogin: "18m ago" },
  { id: "RC-2019", name: "AdBot Studio", email: "ops@adbot.studio", wallet: 2840.80, services: 42, orders: 182, status: "active", lastLogin: "1h ago" },
  { id: "RC-2018", name: "Mai Ngô", email: "mai@social.buzz", wallet: 18.20, services: 3, orders: 12, status: "active", lastLogin: "4h ago" },
  { id: "RC-2017", name: "ScrapeHub VN", email: "team@scrapehub.vn", wallet: 840.00, services: 28, orders: 94, status: "active", lastLogin: "12m ago" },
  { id: "RC-2016", name: "Quang Le", email: "quang@seoking.vn", wallet: 0, services: 2, orders: 8, status: "overdue", lastLogin: "3d ago" },
];

export const RESELLER_CATALOG: ResellerCatalogItem[] = [
  { plan: "Residential · Standard", unit: "per GB", cost: 4.80, selling: 6.50, margin: 35, stock: "ok", version: "v3", status: "active" },
  { plan: "Residential · Premium", unit: "per GB", cost: 7.20, selling: 9.80, margin: 36, stock: "ok", version: "v3", status: "active" },
  { plan: "Datacenter · Shared", unit: "per IP/mo", cost: 0.50, selling: 0.80, margin: 60, stock: "ok", version: "v2", status: "active" },
  { plan: "Datacenter · Dedicated", unit: "per IP/mo", cost: 1.80, selling: 2.20, margin: 22, stock: "low", version: "v2", status: "active" },
  { plan: "VPS Linux · Small (HCM)", unit: "per mo", cost: 14.00, selling: 19.00, margin: 36, stock: "ok", version: "v4", status: "active" },
  { plan: "VPS Linux · Medium (HCM)", unit: "per mo", cost: 36.00, selling: 34.00, margin: -6, stock: "ok", version: "v4", status: "warn" },
];

export const CLIENT_SERVICES: ClientService[] = [
  { id: "svc-r-9281", type: "residential", label: "Residential EU · Premium", identifier: "res-eu-prm-9281", region: "EU-MULTI", bandwidth: "4.2 / 10 GB", expiry: "2026-05-14", status: "active", cycle: "month_30d" },
  { id: "svc-v-4421", type: "vps-linux", label: "vps-scrape-01 · 2C/4G/60G", identifier: "103.28.44.21", region: "VN-HCM", bandwidth: "—", expiry: "2026-05-08", status: "active", cycle: "calendar_month" },
  { id: "svc-d-8102", type: "datacenter", label: "DC Shared · 10 IPs", identifier: "dc-us-8102", region: "US-EAST", bandwidth: "—", expiry: "2026-04-28", status: "active", cycle: "month_30d" },
  { id: "svc-v-4422", type: "vps-linux", label: "vps-test · 1C/2G/20G", identifier: "103.28.44.22", region: "VN-HAN", bandwidth: "—", expiry: "2026-04-24", status: "suspended", cycle: "calendar_month", note: "Grace: 2 days left" },
  { id: "svc-m-2109", type: "mobile", label: "Mobile 4G · 2 ports", identifier: "mob-vn-2109", region: "VN", bandwidth: "188 GB", expiry: "2026-06-02", status: "active", cycle: "month_30d" },
];

export const CLIENT_LEDGER: LedgerEntry[] = [
  { ts: "2026-04-22 14:02", type: "purchase.client_wallet.debit", amount: -62.00, ref: "ORD-48290 · Residential EU", balance: 128.40 },
  { ts: "2026-04-21 09:18", type: "topup.credit.client", amount: 100.00, ref: "TUP-9110 · VietQR", balance: 190.40 },
  { ts: "2026-04-20 14:08", type: "renewal.client_wallet.debit", amount: -19.00, ref: "svc-v-4421 · VPS Small", balance: 90.40 },
  { ts: "2026-04-18 11:22", type: "purchase.client_wallet.debit", amount: -48.00, ref: "ORD-48280 · Mobile 4G", balance: 109.40 },
  { ts: "2026-04-14 10:02", type: "topup.credit.client", amount: 150.00, ref: "TUP-9088 · VietQR", balance: 157.40 },
];

export const PLATFORM_ALERTS = [
  { text: "3 provisioning jobs in manual_review > 1h", type: "danger" as const, screen: "admin-provisioning" },
  { text: "2 reseller tenants below wallet threshold", type: "warn" as const, screen: "admin-topups" },
  { text: "OVH API degraded — 2.8% fail rate", type: "warn" as const, screen: "admin-providers" },
];
