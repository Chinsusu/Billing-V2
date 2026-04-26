// Mock service/account data — aligned with Billing-V2 spec
// No real provider auth material, IPs, or customer data.

export type ServiceStatus = "active" | "suspended" | "provisioning" | "stopped" | "overdue";

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

export type ProxyType = "residential" | "datacenter" | "mobile" | "isp";
export type VpsOS = "linux" | "windows";

export interface ProxyService {
  id: string;
  proxyType: ProxyType;
  label: string;
  customer: string;
  tenant: string;
  region: string;
  ipCount: number;
  protocol: "http" | "socks5" | "both";
  usedGB: number;
  totalGB: number;
  price: number;
  status: ServiceStatus;
  renewsIn: number;
}

export interface VpsService {
  id: string;
  os: VpsOS;
  label: string;
  customer: string;
  tenant: string;
  region: string;
  cpu: number;
  ram: number;
  disk: number;
  ip: string;
  provider: string;
  price: number;
  status: ServiceStatus;
  renewsIn: number;
}

export interface BandwidthService {
  id: string;
  label: string;
  customer: string;
  tenant: string;
  region: string;
  usedGB: number;
  totalGB: number;
  usedPct: number;
  price: number;
  status: ServiceStatus;
  renewsIn: number;
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

export const CUSTOMERS: Customer[] = [
  { id: "ACC-40218", name: "Acme Proxy Co.", email: "ops@acmeproxy.io", plan: "Enterprise", services: 48, mrr: 4280, status: "active", since: "2023-04-12", country: "US" },
  { id: "ACC-40217", name: "Linh Tran", email: "linh.tran@gmail.com", plan: "Pro", services: 7, mrr: 189, status: "active", since: "2024-09-03", country: "VN" },
  { id: "ACC-40216", name: "Scrapers Ltd", email: "billing@scrapers.ltd", plan: "Business", services: 24, mrr: 1840, status: "active", since: "2023-11-28", country: "GB" },
  { id: "ACC-40215", name: "Kenji Watanabe", email: "kenji@tokyonet.jp", plan: "Pro", services: 12, mrr: 420, status: "overdue", since: "2024-02-14", country: "JP" },
  { id: "ACC-40214", name: "Marie Dubois", email: "marie@duboisco.fr", plan: "Business", services: 16, mrr: 1120, status: "active", since: "2023-06-21", country: "FR" },
  { id: "ACC-40213", name: "DataMine Inc.", email: "accounts@datamine.io", plan: "Enterprise", services: 82, mrr: 8420, status: "active", since: "2022-12-01", country: "US" },
  { id: "ACC-40212", name: "Hans Müller", email: "h.mueller@web.de", plan: "Starter", services: 3, mrr: 49, status: "suspended", since: "2025-01-18", country: "DE" },
  { id: "ACC-40211", name: "Proxy Garden", email: "hi@proxygarden.co", plan: "Pro", services: 11, mrr: 340, status: "active", since: "2024-07-09", country: "CA" },
  { id: "ACC-40210", name: "Alex Rodriguez", email: "alex@rodriguez.mx", plan: "Pro", services: 9, mrr: 278, status: "active", since: "2024-11-22", country: "MX" },
  { id: "ACC-40209", name: "Sofia Bergström", email: "sofia@bergnordic.se", plan: "Business", services: 21, mrr: 1480, status: "active", since: "2023-08-15", country: "SE" },
  { id: "ACC-40208", name: "CloudHarvest", email: "billing@cloudharvest.ai", plan: "Enterprise", services: 64, mrr: 6240, status: "active", since: "2023-02-28", country: "US" },
  { id: "ACC-40207", name: "Nguyễn Tuấn", email: "tuan.nguyen@startup.vn", plan: "Starter", services: 2, mrr: 29, status: "active", since: "2025-03-14", country: "VN" },
];

export const SERVICES: Service[] = [
  { id: "SVC-61001", type: "residential", label: "US Residential Pool", customer: "Acme Proxy Co.", region: "US-EAST", bandwidth: "2.4 TB", price: 380, status: "active", renewsIn: 12 },
  { id: "SVC-61002", type: "vps-linux", label: "Production Linux VPS · 8 vCPU · 32GB", customer: "DataMine Inc.", region: "EU-HEL", bandwidth: "—", price: 89, status: "active", renewsIn: 8 },
  { id: "SVC-61003", type: "datacenter", label: "DC Pool Alpha · 500 IPs", customer: "Scrapers Ltd", region: "US-WEST", bandwidth: "1.1 TB", price: 240, status: "active", renewsIn: 22 },
  { id: "SVC-61004", type: "vps-win", label: "Windows RDP Workspace · 4 vCPU · 16GB", customer: "Kenji Watanabe", region: "APAC-TYO", bandwidth: "—", price: 64, status: "overdue", renewsIn: -3 },
  { id: "SVC-61005", type: "mobile", label: "Mobile 4G · 20 ports", customer: "CloudHarvest", region: "US-EAST", bandwidth: "842 GB", price: 620, status: "active", renewsIn: 5 },
  { id: "SVC-61006", type: "isp", label: "ISP Static · 100 IPs", customer: "Marie Dubois", region: "EU-FRA", bandwidth: "680 GB", price: 180, status: "active", renewsIn: 18 },
  { id: "SVC-61007", type: "vps-linux", label: "Proxy Automation VPS · 2 vCPU · 4GB", customer: "Proxy Garden", region: "EU-HEL", bandwidth: "—", price: 19, status: "provisioning", renewsIn: 30 },
  { id: "SVC-61008", type: "residential", label: "EU Residential · Premium", customer: "Linh Tran", region: "EU-MULTI", bandwidth: "128 GB", price: 62, status: "active", renewsIn: 14 },
  { id: "SVC-61009", type: "datacenter", label: "DC Pool Beta · 1000 IPs", customer: "Acme Proxy Co.", region: "GLOBAL", bandwidth: "3.2 TB", price: 440, status: "active", renewsIn: 9 },
  { id: "SVC-61010", type: "residential", label: "APAC Residential", customer: "Alex Rodriguez", region: "APAC-SIN", bandwidth: "184 GB", price: 118, status: "suspended", renewsIn: -7 },
];

export const PROXY_SERVICES: ProxyService[] = [
  { id: "SVC-61001", proxyType: "residential", label: "US Residential Pool",       customer: "Acme Proxy Co.",  tenant: "ProxyVN",      region: "US-EAST",   ipCount: 0,    protocol: "http",   usedGB: 2400, totalGB: 5000, price: 380, status: "active",       renewsIn: 12 },
  { id: "SVC-61003", proxyType: "datacenter",  label: "DC Pool Alpha · 500 IPs",   customer: "Scrapers Ltd",    tenant: "HANetwork",    region: "US-WEST",   ipCount: 500,  protocol: "both",   usedGB: 1100, totalGB: 3000, price: 240, status: "active",       renewsIn: 22 },
  { id: "SVC-61005", proxyType: "mobile",      label: "Mobile 4G · 20 ports",      customer: "CloudHarvest",    tenant: "ProxyVN",      region: "US-EAST",   ipCount: 20,   protocol: "socks5", usedGB: 842,  totalGB: 2000, price: 620, status: "active",       renewsIn: 5  },
  { id: "SVC-61006", proxyType: "isp",         label: "ISP Static · 100 IPs",      customer: "Marie Dubois",    tenant: "ProxyVN",      region: "EU-FRA",    ipCount: 100,  protocol: "http",   usedGB: 680,  totalGB: 2000, price: 180, status: "active",       renewsIn: 18 },
  { id: "SVC-61008", proxyType: "residential", label: "EU Residential · Premium",  customer: "Linh Tran",       tenant: "ProxyVN",      region: "EU-MULTI",  ipCount: 0,    protocol: "http",   usedGB: 128,  totalGB: 500,  price: 62,  status: "active",       renewsIn: 14 },
  { id: "SVC-61009", proxyType: "datacenter",  label: "DC Pool Beta · 1000 IPs",   customer: "Acme Proxy Co.",  tenant: "HANetwork",    region: "GLOBAL",    ipCount: 1000, protocol: "both",   usedGB: 3200, totalGB: 8000, price: 440, status: "active",       renewsIn: 9  },
  { id: "SVC-61010", proxyType: "residential", label: "APAC Residential",          customer: "Alex Rodriguez",  tenant: "ProxyVN",      region: "APAC-SIN",  ipCount: 0,    protocol: "http",   usedGB: 184,  totalGB: 500,  price: 118, status: "suspended",    renewsIn: -7 },
  { id: "SVC-61011", proxyType: "isp",         label: "ISP UK · 50 IPs",           customer: "Scrapers Ltd",    tenant: "HANetwork",    region: "EU-LON",    ipCount: 50,   protocol: "both",   usedGB: 310,  totalGB: 1000, price: 95,  status: "active",       renewsIn: 3  },
  { id: "SVC-61012", proxyType: "mobile",      label: "Mobile VN · 10 ports",      customer: "Nguyễn Tuấn",     tenant: "ProxyVN",      region: "VN-HCM",    ipCount: 10,   protocol: "socks5", usedGB: 55,   totalGB: 200,  price: 48,  status: "active",       renewsIn: 21 },
  { id: "SVC-61013", proxyType: "datacenter",  label: "DC APAC · 200 IPs",         customer: "DataMine Inc.",   tenant: "HANetwork",    region: "APAC-SIN",  ipCount: 200,  protocol: "http",   usedGB: 890,  totalGB: 2000, price: 160, status: "provisioning", renewsIn: 30 },
];

export const VPS_SERVICES: VpsService[] = [
  { id: "SVC-61002", os: "linux",   label: "Production Linux VPS",      customer: "DataMine Inc.",   tenant: "HANetwork",  region: "EU-HEL",   cpu: 8,  ram: 32,  disk: 320,  ip: "95.216.x.x",   provider: "Hetzner",  price: 89,  status: "active",       renewsIn: 8  },
  { id: "SVC-61004", os: "windows", label: "Windows RDP Workspace",     customer: "Kenji Watanabe",  tenant: "ProxyVN",    region: "APAC-TYO", cpu: 4,  ram: 16,  disk: 160,  ip: "45.77.x.x",    provider: "OVH",      price: 64,  status: "overdue",      renewsIn: -3 },
  { id: "SVC-61007", os: "linux",   label: "Proxy Automation VPS",      customer: "Proxy Garden",    tenant: "ProxyVN",    region: "EU-HEL",   cpu: 2,  ram: 4,   disk: 40,   ip: "—",            provider: "Proxmox",  price: 19,  status: "provisioning", renewsIn: 30 },
  { id: "SVC-61014", os: "linux",   label: "API Gateway VPS",           customer: "Acme Proxy Co.",  tenant: "HANetwork",  region: "US-EAST",  cpu: 4,  ram: 8,   disk: 80,   ip: "104.21.x.x",   provider: "Hetzner",  price: 34,  status: "active",       renewsIn: 15 },
  { id: "SVC-61015", os: "linux",   label: "Database Replica VPS",      customer: "CloudHarvest",    tenant: "HANetwork",  region: "US-WEST",  cpu: 8,  ram: 32,  disk: 640,  ip: "34.102.x.x",   provider: "OVH",      price: 102, status: "active",       renewsIn: 6  },
  { id: "SVC-61016", os: "windows", label: "Windows Development VPS",   customer: "Marie Dubois",    tenant: "ProxyVN",    region: "EU-FRA",   cpu: 2,  ram: 8,   disk: 80,   ip: "51.75.x.x",    provider: "OVH",      price: 44,  status: "suspended",    renewsIn: -1 },
  { id: "SVC-61017", os: "linux",   label: "Batch Worker VPS",          customer: "Scrapers Ltd",    tenant: "HANetwork",  region: "EU-HEL",   cpu: 16, ram: 64,  disk: 960,  ip: "65.108.x.x",   provider: "Hetzner",  price: 188, status: "active",       renewsIn: 19 },
];

export const BANDWIDTH_SERVICES: BandwidthService[] = [
  { id: "SVC-62001", label: "Residential US 5TB",    customer: "Acme Proxy Co.",  tenant: "ProxyVN",    region: "US-EAST",  usedGB: 2400, totalGB: 5120,  usedPct: 47, price: 380, status: "active",  renewsIn: 12 },
  { id: "SVC-62002", label: "DC Pool Global 8TB",    customer: "DataMine Inc.",   tenant: "HANetwork",  region: "GLOBAL",   usedGB: 5800, totalGB: 8192,  usedPct: 71, price: 620, status: "active",  renewsIn: 9  },
  { id: "SVC-62003", label: "Mobile US 2TB",         customer: "CloudHarvest",    tenant: "ProxyVN",    region: "US-EAST",  usedGB: 842,  totalGB: 2048,  usedPct: 41, price: 210, status: "active",  renewsIn: 5  },
  { id: "SVC-62004", label: "ISP EU 2TB",            customer: "Marie Dubois",    tenant: "ProxyVN",    region: "EU-MULTI", usedGB: 680,  totalGB: 2048,  usedPct: 33, price: 180, status: "active",  renewsIn: 18 },
  { id: "SVC-62005", label: "EU Residential 500GB",  customer: "Linh Tran",       tenant: "ProxyVN",    region: "EU-MULTI", usedGB: 128,  totalGB: 512,   usedPct: 25, price: 62,  status: "active",  renewsIn: 14 },
  { id: "SVC-62006", label: "APAC Residential 500GB",customer: "Alex Rodriguez",  tenant: "ProxyVN",    region: "APAC-SIN", usedGB: 501,  totalGB: 512,   usedPct: 98, price: 118, status: "overdue", renewsIn: -7 },
  { id: "SVC-62007", label: "ISP UK 1TB",            customer: "Scrapers Ltd",    tenant: "HANetwork",  region: "EU-LON",   usedGB: 310,  totalGB: 1024,  usedPct: 30, price: 95,  status: "active",  renewsIn: 3  },
];

export const RESELLER_CLIENTS: ResellerClient[] = [
  { id: "ACC-42021", name: "Linh Tran", email: "linh@scrape.dev", wallet: 128.40, services: 7, orders: 24, status: "active", lastLogin: "2h ago" },
  { id: "ACC-42020", name: "Hùng Phạm", email: "hung@adsvn.co", wallet: 420.00, services: 12, orders: 48, status: "active", lastLogin: "18m ago" },
  { id: "ACC-42019", name: "AdBot Studio", email: "ops@adbot.studio", wallet: 2840.80, services: 42, orders: 182, status: "active", lastLogin: "1h ago" },
  { id: "ACC-42018", name: "Mai Ngô", email: "mai@social.buzz", wallet: 18.20, services: 3, orders: 12, status: "active", lastLogin: "4h ago" },
  { id: "ACC-42017", name: "ScrapeHub VN", email: "team@scrapehub.vn", wallet: 840.00, services: 28, orders: 94, status: "active", lastLogin: "12m ago" },
  { id: "ACC-42016", name: "Quang Le", email: "quang@seoking.vn", wallet: 0, services: 2, orders: 8, status: "overdue", lastLogin: "3d ago" },
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
  { id: "SVC-69281", type: "residential", label: "Residential EU · Premium", identifier: "res-eu-prm-9281", region: "EU-MULTI", bandwidth: "4.2 / 10 GB", expiry: "2026-05-14", status: "active", cycle: "month_30d" },
  { id: "SVC-64421", type: "vps-linux", label: "Proxy Automation VPS · 2C/4G/60G", identifier: "103.28.44.21", region: "VN-HCM", bandwidth: "—", expiry: "2026-05-08", status: "active", cycle: "calendar_month" },
  { id: "SVC-68102", type: "datacenter", label: "DC Shared · 10 IPs", identifier: "dc-us-8102", region: "US-EAST", bandwidth: "—", expiry: "2026-04-28", status: "active", cycle: "month_30d" },
  { id: "SVC-64422", type: "vps-linux", label: "Small Linux VPS · 1C/2G/20G", identifier: "103.28.44.22", region: "VN-HAN", bandwidth: "—", expiry: "2026-04-24", status: "suspended", cycle: "calendar_month", note: "Grace: 2 days left" },
  { id: "SVC-62109", type: "mobile", label: "Mobile 4G · 2 ports", identifier: "mob-vn-2109", region: "VN", bandwidth: "188 GB", expiry: "2026-06-02", status: "active", cycle: "month_30d" },
];
