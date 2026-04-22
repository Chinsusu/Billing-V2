// Billing-V2 data — aligned with docs spec (tenant model, wallets, ledger, provisioning)

const BILLING = {
  // Tenants: 1 admin + 3 reseller white-label tenants
  tenants: [
    { id: 'T-0001', name: 'HANetwork (Admin)', type: 'admin', domain: 'billing.hanetwork.vn', clients: 1284, services: 8420, wallet: 0, status: 'active' },
    { id: 'T-0042', name: 'ProxyVN Reseller', type: 'reseller', domain: 'proxyvn.io', clients: 312, services: 1840, wallet: 4820.50, walletLow: false, status: 'active', since: '2024-03-12' },
    { id: 'T-0051', name: 'CloudBase.asia', type: 'reseller', domain: 'cloudbase.asia', clients: 184, services: 924, wallet: 184.20, walletLow: true, status: 'active', since: '2024-07-28' },
    { id: 'T-0063', name: 'Saigon Proxy', type: 'reseller', domain: 'proxy.saigon.vn', clients: 94, services: 421, wallet: 2140.80, walletLow: false, status: 'active', since: '2025-01-05' },
    { id: 'T-0071', name: 'DanangHost', type: 'reseller', domain: 'danang.host', clients: 42, services: 128, wallet: 0, walletLow: true, status: 'suspended', since: '2025-03-20' },
  ],

  providers: [
    { id: 'prv-pmx-01', name: 'Proxmox · VN-HCM', type: 'self-host', health: 'ok', capacity: 82, failRate: 0.2, lastSync: '2m ago' },
    { id: 'prv-pmx-02', name: 'Proxmox · VN-HAN', type: 'self-host', health: 'ok', capacity: 64, failRate: 0.1, lastSync: '1m ago' },
    { id: 'prv-pmg', name: 'proxy-manager (self)', type: 'self-host', health: 'ok', capacity: 91, failRate: 0.3, lastSync: '4m ago' },
    { id: 'prv-ovh', name: 'OVH', type: 'upstream', health: 'degraded', capacity: 54, failRate: 2.8, lastSync: '12m ago' },
    { id: 'prv-hzn', name: 'Hetzner', type: 'upstream', health: 'ok', capacity: 78, failRate: 0.4, lastSync: '3m ago' },
    { id: 'prv-smh', name: 'Smarthost', type: 'upstream', health: 'ok', capacity: 42, failRate: 0.9, lastSync: '5m ago' },
    { id: 'prv-pch', name: 'proxy-cheap', type: 'upstream', health: 'ok', capacity: 88, failRate: 0.5, lastSync: '2m ago' },
  ],

  // Provisioning queue per spec: queued / provisioning / failed / manual_review
  provJobs: [
    { id: 'job-8a21', order: 'ORD-48291', service: 'VPS 4C/8G · HCM', tenant: 'ProxyVN', provider: 'Proxmox · VN-HCM', status: 'manual_review', attempt: 2, error: 'provider_timeout — resource state unknown', correlation: 'cor_9a21f8', age: '18m' },
    { id: 'job-8a20', order: 'ORD-48290', service: 'Residential · 5GB', tenant: 'CloudBase', provider: 'proxy-cheap', status: 'failed', attempt: 3, error: 'auth_failed', correlation: 'cor_9a20c4', age: '22m' },
    { id: 'job-8a19', order: 'ORD-48289', service: 'VPS 2C/4G · SG', tenant: 'HANetwork', provider: 'OVH', status: 'provisioning', attempt: 1, error: '', correlation: 'cor_9a19a1', age: '3m' },
    { id: 'job-8a18', order: 'ORD-48288', service: 'DC Proxy · 100 IPs', tenant: 'Saigon Proxy', provider: 'proxy-cheap', status: 'provisioning', attempt: 1, error: '', correlation: 'cor_9a18e2', age: '1m' },
    { id: 'job-8a17', order: 'ORD-48287', service: 'VPS 8C/16G · HEL', tenant: 'HANetwork', provider: 'Hetzner', status: 'queued', attempt: 0, error: '', correlation: 'cor_9a17b8', age: '0m' },
    { id: 'job-8a16', order: 'ORD-48286', service: 'ISP Static · 50 IPs', tenant: 'ProxyVN', provider: 'Smarthost', status: 'manual_review', attempt: 1, error: 'partial_success — external_id unknown', correlation: 'cor_9a16d3', age: '1h 04m' },
  ],

  // Top-up review queue
  topups: [
    { id: 'TUP-9120', tenant: 'ProxyVN', actor: 'reseller_wallet', amount: 2000, currency: 'USD', method: 'VietQR', ref: 'FT26042200832', status: 'pending_verification', created: '2026-04-22 14:02', proof: true },
    { id: 'TUP-9119', tenant: 'CloudBase', actor: 'reseller_wallet', amount: 500, currency: 'USD', method: 'USDT', ref: '0x8a…d4e1', status: 'pending_verification', created: '2026-04-22 13:40', proof: true },
    { id: 'TUP-9118', tenant: 'ProxyVN > linh.tran', actor: 'client_wallet', amount: 100, currency: 'USD', method: 'VietQR', ref: 'FT26042200781', status: 'pending_verification', created: '2026-04-22 13:18', proof: true },
    { id: 'TUP-9117', tenant: 'HANetwork > kenji.w', actor: 'client_wallet', amount: 50, currency: 'USD', method: 'USDT', ref: '0x1c…9a82', status: 'approved', created: '2026-04-22 11:50', proof: true },
    { id: 'TUP-9116', tenant: 'Saigon Proxy', actor: 'reseller_wallet', amount: 1000, currency: 'USD', method: 'VietQR', ref: 'FT26042200412', status: 'approved', created: '2026-04-22 10:12', proof: true },
    { id: 'TUP-9115', tenant: 'CloudBase > huy.nguyen', actor: 'client_wallet', amount: 200, currency: 'USD', method: 'VietQR', ref: 'FT26042200388', status: 'rejected', created: '2026-04-22 09:40', proof: false, reason: 'bank reference not found' },
  ],

  // Reseller client list (for ProxyVN tenant)
  resellerClients: [
    { id: 'RC-2021', name: 'Linh Tran', email: 'linh@scrape.dev', wallet: 128.40, services: 7, orders: 24, status: 'active', lastLogin: '2h ago' },
    { id: 'RC-2020', name: 'Hùng Phạm', email: 'hung@adsvn.co', wallet: 420.00, services: 12, orders: 48, status: 'active', lastLogin: '18m ago' },
    { id: 'RC-2019', name: 'AdBot Studio', email: 'ops@adbot.studio', wallet: 2840.80, services: 42, orders: 182, status: 'active', lastLogin: '1h ago' },
    { id: 'RC-2018', name: 'Mai Ngô', email: 'mai@social.buzz', wallet: 18.20, services: 3, orders: 12, status: 'active', lastLogin: '4h ago' },
    { id: 'RC-2017', name: 'ScrapeHub VN', email: 'team@scrapehub.vn', wallet: 840.00, services: 28, orders: 94, status: 'active', lastLogin: '12m ago' },
    { id: 'RC-2016', name: 'Quang Le', email: 'quang@seoking.vn', wallet: 0, services: 2, orders: 8, status: 'overdue', lastLogin: '3d ago' },
    { id: 'RC-2015', name: 'Tien Đỗ', email: 'tien@datalab.io', wallet: 62.00, services: 5, orders: 18, status: 'active', lastLogin: '6h ago' },
  ],

  // Reseller catalog clones (tenant plans with margin)
  resellerCatalog: [
    { plan: 'Residential · Standard', unit: 'per GB', cost: 4.80, selling: 6.50, margin: 35, stock: 'ok', version: 'v3', status: 'active' },
    { plan: 'Residential · Premium', unit: 'per GB', cost: 7.20, selling: 9.80, margin: 36, stock: 'ok', version: 'v3', status: 'active' },
    { plan: 'Datacenter · Shared', unit: 'per IP/mo', cost: 0.50, selling: 0.80, margin: 60, stock: 'ok', version: 'v2', status: 'active' },
    { plan: 'Datacenter · Dedicated', unit: 'per IP/mo', cost: 1.80, selling: 2.20, margin: 22, stock: 'low', version: 'v2', status: 'active' },
    { plan: 'ISP Static', unit: 'per IP/mo', cost: 2.80, selling: 3.50, margin: 25, stock: 'ok', version: 'v1', status: 'active' },
    { plan: 'Mobile 4G', unit: 'per port/mo', cost: 38.00, selling: 48.00, margin: 26, stock: 'ok', version: 'v1', status: 'active' },
    { plan: 'VPS Linux · Small (HCM)', unit: 'per mo', cost: 14.00, selling: 19.00, margin: 36, stock: 'ok', version: 'v4', status: 'active' },
    { plan: 'VPS Linux · Medium (HCM)', unit: 'per mo', cost: 36.00, selling: 34.00, margin: -6, stock: 'ok', version: 'v4', status: 'warn' },
    { plan: 'VPS Linux · Large (HEL)', unit: 'per mo', cost: 98.00, selling: 129.00, margin: 32, stock: 'out', version: 'v3', status: 'disabled' },
    { plan: 'VPS Windows · Medium', unit: 'per mo', cost: 58.00, selling: 78.00, margin: 34, stock: 'ok', version: 'v2', status: 'active' },
  ],

  // Client-facing services (Linh Tran's view — reseller client)
  clientServices: [
    { id: 'svc-r-9281', type: 'residential', label: 'Residential EU · Premium', identifier: 'res-eu-prm-9281', region: 'EU-MULTI', bandwidth: '4.2 / 10 GB', expiry: '2026-05-14', status: 'active', cycle: 'month_30d' },
    { id: 'svc-v-4421', type: 'vps-linux', label: 'vps-scrape-01 · 2C/4G/60G', identifier: '103.28.44.21', region: 'VN-HCM', bandwidth: '—', expiry: '2026-05-08', status: 'active', cycle: 'calendar_month' },
    { id: 'svc-d-8102', type: 'datacenter', label: 'DC Shared · 10 IPs', identifier: 'dc-us-8102', region: 'US-EAST', bandwidth: '—', expiry: '2026-04-28', status: 'active', cycle: 'month_30d' },
    { id: 'svc-v-4422', type: 'vps-linux', label: 'vps-test · 1C/2G/20G', identifier: '103.28.44.22', region: 'VN-HAN', bandwidth: '—', expiry: '2026-04-24', status: 'suspended', cycle: 'calendar_month', note: 'Grace: 2 days left' },
    { id: 'svc-m-2109', type: 'mobile', label: 'Mobile 4G · 2 ports', identifier: 'mob-vn-2109', region: 'VN', bandwidth: '188 GB', expiry: '2026-06-02', status: 'active', cycle: 'month_30d' },
    { id: 'svc-v-4423', type: 'vps-linux', label: 'vps-pending · 4C/8G', identifier: 'provisioning…', region: 'VN-HCM', bandwidth: '—', expiry: '—', status: 'provisioning', cycle: 'calendar_month' },
    { id: 'svc-r-9282', type: 'residential', label: 'Residential APAC · Std', identifier: 'res-ap-std-9282', region: 'APAC', bandwidth: '0.8 / 5 GB', expiry: '2026-07-10', status: 'active', cycle: 'month_30d' },
  ],

  // Client ledger entries
  clientLedger: [
    { ts: '2026-04-22 14:02', type: 'purchase.client_wallet.debit', amount: -62.00, ref: 'ORD-48290 · Residential EU', balance: 128.40 },
    { ts: '2026-04-21 09:18', type: 'topup.credit.client', amount: 100.00, ref: 'TUP-9110 · VietQR', balance: 190.40 },
    { ts: '2026-04-20 14:08', type: 'renewal.client_wallet.debit', amount: -19.00, ref: 'svc-v-4421 · VPS Small', balance: 90.40 },
    { ts: '2026-04-18 11:22', type: 'purchase.client_wallet.debit', amount: -48.00, ref: 'ORD-48280 · Mobile 4G', balance: 109.40 },
    { ts: '2026-04-14 10:02', type: 'topup.credit.client', amount: 150.00, ref: 'TUP-9088 · VietQR', balance: 157.40 },
    { ts: '2026-04-12 16:44', type: 'refund.client_wallet.credit', amount: 15.20, ref: 'ORD-48252 · partial refund', balance: 7.40 },
  ],

  // Revenue splits
  revenueSplit: {
    retailAdmin: 42800,
    resellerFee: 18200,
    total: 61000,
  },

  // Alerts
  alerts: [
    { icon: 'alert', text: '3 provisioning jobs in manual_review > 1h', type: 'danger', action: 'Review queue' },
    { icon: 'card', text: '2 reseller tenants below wallet threshold', type: 'warn', action: 'View' },
    { icon: 'globe', text: 'OVH API degraded — 2.8% fail rate', type: 'warn', action: 'Provider status' },
    { icon: 'users', text: '6 top-up requests waiting > 30min', type: 'warn', action: 'Verification queue' },
  ],
};

const fmtMargin = (m) => (m >= 0 ? '+' : '') + m + '%';

window.BILLING = BILLING;
window.fmtMargin = fmtMargin;
