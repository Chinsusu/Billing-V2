// Sample data for the admin dashboard
const SAMPLE = {
  revenue30d: [12400, 11800, 13100, 12900, 14200, 15100, 14800, 16300, 15900, 17200, 16800, 18100, 17900, 19300, 18800, 20100, 21400, 20900, 22300, 22100, 23800, 23400, 24900, 24600, 26200, 25800, 27400, 27100, 28600, 29200],
  mrr30d: [86000, 87200, 87800, 89100, 90400, 91200, 92500, 93100, 94800, 95400, 96700, 97300, 98200, 99500, 100200, 101400, 102800, 103500, 104900, 105600, 107200, 108100, 109400, 110200, 111700, 112800, 114100, 115300, 116800, 118200],
  customers30d: [2640, 2658, 2671, 2682, 2694, 2706, 2718, 2731, 2744, 2752, 2763, 2775, 2781, 2792, 2803, 2811, 2819, 2828, 2836, 2840, 2843, 2847, 2849, 2851, 2853, 2854, 2846, 2847, 2847, 2847],
  bandwidthDaily: [184, 192, 201, 198, 215, 223, 219, 234, 229, 245, 251, 248, 262, 268, 271, 284, 278, 291, 295, 302, 298, 311, 315, 309, 322, 318, 331, 327, 339, 342],
  monthLabels: Array.from({length: 30}, (_, i) => i % 5 === 0 ? `Apr ${i+1}` : ''),

  products: [
    { name: 'Residential Proxy', sold: 8421, rev: 124800, color: '#D50C2D' },
    { name: 'Datacenter Proxy', sold: 3284, rev: 52400, color: '#1F2937' },
    { name: 'VPS Linux', sold: 1892, rev: 84600, color: '#0B7A3B' },
    { name: 'ISP Proxy', sold: 612, rev: 38900, color: '#A15C00' },
    { name: 'Mobile Proxy', sold: 268, rev: 27100, color: '#1E4FA3' },
    { name: 'VPS Windows', sold: 402, rev: 24800, color: '#7c3aed' },
    { name: 'Bandwidth', sold: 1128, rev: 14200, color: '#6B7280' },
  ],

  customers: [
    { id: 'C-40218', name: 'Acme Proxy Co.', email: 'ops@acmeproxy.io', plan: 'Enterprise', services: 48, mrr: 4280, status: 'active', since: '2023-04-12', country: 'US' },
    { id: 'C-40217', name: 'Linh Tran', email: 'linh.tran@gmail.com', plan: 'Pro', services: 7, mrr: 189, status: 'active', since: '2024-09-03', country: 'VN' },
    { id: 'C-40216', name: 'Scrapers Ltd', email: 'billing@scrapers.ltd', plan: 'Business', services: 24, mrr: 1840, status: 'active', since: '2023-11-28', country: 'GB' },
    { id: 'C-40215', name: 'Kenji Watanabe', email: 'kenji@tokyonet.jp', plan: 'Pro', services: 12, mrr: 420, status: 'overdue', since: '2024-02-14', country: 'JP' },
    { id: 'C-40214', name: 'Marie Dubois', email: 'marie@duboisco.fr', plan: 'Business', services: 16, mrr: 1120, status: 'active', since: '2023-06-21', country: 'FR' },
    { id: 'C-40213', name: 'DataMine Inc.', email: 'accounts@datamine.io', plan: 'Enterprise', services: 82, mrr: 8420, status: 'active', since: '2022-12-01', country: 'US' },
    { id: 'C-40212', name: 'Hans Müller', email: 'h.mueller@web.de', plan: 'Starter', services: 3, mrr: 49, status: 'suspended', since: '2025-01-18', country: 'DE' },
    { id: 'C-40211', name: 'Proxy Garden', email: 'hi@proxygarden.co', plan: 'Pro', services: 11, mrr: 340, status: 'active', since: '2024-07-09', country: 'CA' },
    { id: 'C-40210', name: 'Alex Rodriguez', email: 'alex@rodriguez.mx', plan: 'Pro', services: 9, mrr: 278, status: 'active', since: '2024-11-22', country: 'MX' },
    { id: 'C-40209', name: 'Sofia Bergström', email: 'sofia@bergnordic.se', plan: 'Business', services: 21, mrr: 1480, status: 'active', since: '2023-08-15', country: 'SE' },
    { id: 'C-40208', name: 'CloudHarvest', email: 'billing@cloudharvest.ai', plan: 'Enterprise', services: 64, mrr: 6240, status: 'active', since: '2023-02-28', country: 'US' },
    { id: 'C-40207', name: 'Nguyễn Tuấn', email: 'tuan.nguyen@startup.vn', plan: 'Starter', services: 2, mrr: 29, status: 'active', since: '2025-03-14', country: 'VN' },
  ],

  services: [
    { id: 'prx-9a12', type: 'residential', label: 'US Residential Pool', customer: 'Acme Proxy Co.', region: 'US-EAST', bandwidth: '2.4 TB', price: 380, status: 'active', renewsIn: 12 },
    { id: 'vps-7c8f', type: 'vps-linux', label: 'vps-prod-01 · 8 vCPU · 32GB', customer: 'DataMine Inc.', region: 'EU-HEL', bandwidth: '—', price: 89, status: 'active', renewsIn: 8 },
    { id: 'prx-4e21', type: 'datacenter', label: 'DC Pool Alpha · 500 IPs', customer: 'Scrapers Ltd', region: 'US-WEST', bandwidth: '1.1 TB', price: 240, status: 'active', renewsIn: 22 },
    { id: 'vps-3d9a', type: 'vps-win', label: 'win-rdp-gamma · 4 vCPU · 16GB', customer: 'Kenji Watanabe', region: 'APAC-TYO', bandwidth: '—', price: 64, status: 'overdue', renewsIn: -3 },
    { id: 'prx-8b67', type: 'mobile', label: 'Mobile 4G · 20 ports', customer: 'CloudHarvest', region: 'US-EAST', bandwidth: '842 GB', price: 620, status: 'active', renewsIn: 5 },
    { id: 'prx-2f01', type: 'isp', label: 'ISP Static · 100 IPs', customer: 'Marie Dubois', region: 'EU-FRA', bandwidth: '680 GB', price: 180, status: 'active', renewsIn: 18 },
    { id: 'vps-6a11', type: 'vps-linux', label: 'vps-scrape-02 · 2 vCPU · 4GB', customer: 'Proxy Garden', region: 'EU-HEL', bandwidth: '—', price: 19, status: 'provisioning', renewsIn: 30 },
    { id: 'prx-5c23', type: 'residential', label: 'EU Residential · Premium', customer: 'Linh Tran', region: 'EU-MULTI', bandwidth: '128 GB', price: 62, status: 'active', renewsIn: 14 },
    { id: 'vps-9e44', type: 'vps-linux', label: 'vps-db-01 · 16 vCPU · 64GB', customer: 'DataMine Inc.', region: 'US-ASH', bandwidth: '—', price: 189, status: 'active', renewsIn: 27 },
    { id: 'prx-1a88', type: 'datacenter', label: 'DC Pool Beta · 1000 IPs', customer: 'Acme Proxy Co.', region: 'GLOBAL', bandwidth: '3.2 TB', price: 440, status: 'active', renewsIn: 9 },
    { id: 'prx-7b34', type: 'residential', label: 'APAC Residential', customer: 'Alex Rodriguez', region: 'APAC-SIN', bandwidth: '184 GB', price: 118, status: 'suspended', renewsIn: -7 },
    { id: 'vps-4f56', type: 'vps-linux', label: 'vps-api-03 · 4 vCPU · 8GB', customer: 'Sofia Bergström', region: 'EU-HEL', bandwidth: '—', price: 32, status: 'active', renewsIn: 21 },
  ],

  invoices: [
    { id: 'INV-2026-04218', customer: 'Acme Proxy Co.', issued: '2026-04-20', due: '2026-05-04', amount: 4280.00, status: 'open' },
    { id: 'INV-2026-04217', customer: 'DataMine Inc.', issued: '2026-04-20', due: '2026-05-04', amount: 8420.00, status: 'paid' },
    { id: 'INV-2026-04216', customer: 'Kenji Watanabe', issued: '2026-04-15', due: '2026-04-29', amount: 420.00, status: 'overdue' },
    { id: 'INV-2026-04215', customer: 'CloudHarvest', issued: '2026-04-18', due: '2026-05-02', amount: 6240.00, status: 'paid' },
    { id: 'INV-2026-04214', customer: 'Marie Dubois', issued: '2026-04-18', due: '2026-05-02', amount: 1120.00, status: 'paid' },
    { id: 'INV-2026-04213', customer: 'Scrapers Ltd', issued: '2026-04-17', due: '2026-05-01', amount: 1840.00, status: 'open' },
    { id: 'INV-2026-04212', customer: 'Sofia Bergström', issued: '2026-04-15', due: '2026-04-29', amount: 1480.00, status: 'paid' },
    { id: 'INV-2026-04211', customer: 'Linh Tran', issued: '2026-04-14', due: '2026-04-28', amount: 189.00, status: 'paid' },
    { id: 'INV-2026-04210', customer: 'Proxy Garden', issued: '2026-04-12', due: '2026-04-26', amount: 340.00, status: 'paid' },
    { id: 'INV-2026-04209', customer: 'Alex Rodriguez', issued: '2026-04-10', due: '2026-04-24', amount: 278.00, status: 'paid' },
  ],

  transactions: [
    { id: 'txn_9A8f21', time: '2026-04-22 14:22', customer: 'DataMine Inc.', method: 'Visa •• 4242', type: 'charge', amount: 8420.00, status: 'paid' },
    { id: 'txn_9A8e12', time: '2026-04-22 13:48', customer: 'CloudHarvest', method: 'ACH', type: 'charge', amount: 6240.00, status: 'paid' },
    { id: 'txn_9A8d04', time: '2026-04-22 11:17', customer: 'Linh Tran', method: 'Wallet', type: 'topup', amount: 500.00, status: 'paid' },
    { id: 'txn_9A8c92', time: '2026-04-22 10:02', customer: 'Kenji Watanabe', method: 'Visa •• 0914', type: 'charge', amount: 420.00, status: 'failed' },
    { id: 'txn_9A8b77', time: '2026-04-22 09:41', customer: 'Marie Dubois', method: 'Mastercard •• 1821', type: 'charge', amount: 1120.00, status: 'paid' },
    { id: 'txn_9A8a31', time: '2026-04-22 08:12', customer: 'Proxy Garden', method: 'PayPal', type: 'charge', amount: 340.00, status: 'paid' },
    { id: 'txn_9A8920', time: '2026-04-21 22:58', customer: 'Acme Proxy Co.', method: 'Wire', type: 'charge', amount: 4280.00, status: 'pending' },
    { id: 'txn_9A880f', time: '2026-04-21 19:34', customer: 'Alex Rodriguez', method: 'Visa •• 6601', type: 'charge', amount: 278.00, status: 'paid' },
    { id: 'txn_9A87e1', time: '2026-04-21 16:03', customer: 'Sofia Bergström', method: 'Klarna', type: 'charge', amount: 1480.00, status: 'paid' },
    { id: 'txn_9A86b5', time: '2026-04-21 14:21', customer: 'Nguyễn Tuấn', method: 'Wallet', type: 'refund', amount: -29.00, status: 'paid' },
  ],

  tickets: [
    { id: 'T-8124', subject: 'Proxy pool authentication failing intermittently', customer: 'Acme Proxy Co.', priority: 'high', status: 'open', updated: '12m ago', assignee: 'Linh' },
    { id: 'T-8123', subject: 'Invoice INV-2026-04216 — requesting extension', customer: 'Kenji Watanabe', priority: 'medium', status: 'pending', updated: '38m ago', assignee: 'Minh' },
    { id: 'T-8122', subject: 'VPS Windows license key issue', customer: 'DataMine Inc.', priority: 'low', status: 'open', updated: '1h ago', assignee: 'Tùng' },
    { id: 'T-8121', subject: 'Bandwidth overage pricing clarification', customer: 'Scrapers Ltd', priority: 'low', status: 'open', updated: '2h ago', assignee: '—' },
    { id: 'T-8120', subject: 'Request: dedicated IP block allocation', customer: 'CloudHarvest', priority: 'medium', status: 'pending', updated: '3h ago', assignee: 'Linh' },
    { id: 'T-8119', subject: 'Can\'t access control panel — 2FA lockout', customer: 'Linh Tran', priority: 'high', status: 'open', updated: '4h ago', assignee: 'Minh' },
    { id: 'T-8118', subject: 'Refund for accidentally renewed service', customer: 'Hans Müller', priority: 'low', status: 'pending', updated: '6h ago', assignee: 'Tùng' },
    { id: 'T-8117', subject: 'API rate limit increase request', customer: 'DataMine Inc.', priority: 'medium', status: 'open', updated: '8h ago', assignee: 'Linh' },
  ],

  products_catalog: [
    { sku: 'PRX-RES-STD', name: 'Residential · Standard', unit: 'per GB', price: 6.50, active: 2841, rev30: 124800 },
    { sku: 'PRX-RES-PRM', name: 'Residential · Premium', unit: 'per GB', price: 9.80, active: 1204, rev30: 68200 },
    { sku: 'PRX-DC-SHR', name: 'Datacenter · Shared', unit: 'per IP/mo', price: 0.80, active: 8920, rev30: 52400 },
    { sku: 'PRX-DC-DED', name: 'Datacenter · Dedicated', unit: 'per IP/mo', price: 2.20, active: 1840, rev30: 38200 },
    { sku: 'PRX-ISP-STC', name: 'ISP Static', unit: 'per IP/mo', price: 3.50, active: 612, rev30: 38900 },
    { sku: 'PRX-MOB-4G', name: 'Mobile 4G · Port', unit: 'per port/mo', price: 48.00, active: 268, rev30: 27100 },
    { sku: 'VPS-LNX-S', name: 'VPS Linux · Small', unit: 'per mo', price: 19.00, active: 842, rev30: 18400 },
    { sku: 'VPS-LNX-M', name: 'VPS Linux · Medium', unit: 'per mo', price: 48.00, active: 614, rev30: 34800 },
    { sku: 'VPS-LNX-L', name: 'VPS Linux · Large', unit: 'per mo', price: 129.00, active: 312, rev30: 28600 },
    { sku: 'VPS-WIN-M', name: 'VPS Windows · Medium', unit: 'per mo', price: 78.00, active: 402, rev30: 24800 },
    { sku: 'BW-PKG-1TB', name: 'Bandwidth 1TB Package', unit: 'per pkg', price: 12.00, active: 1128, rev30: 14200 },
  ],

  activity: [
    { t: '14:22', icon: 'card', text: 'Payment of $8,420.00 from DataMine Inc.', type: 'ok' },
    { t: '14:11', icon: 'users', text: 'New customer signup: startup-dev-42@proton.me', type: 'info' },
    { t: '14:02', icon: 'server', text: 'VPS vps-scrape-02 provisioned for Proxy Garden', type: 'info' },
    { t: '13:48', icon: 'card', text: 'Payment of $6,240.00 from CloudHarvest', type: 'ok' },
    { t: '13:17', icon: 'ticket', text: 'New ticket T-8124 opened by Acme Proxy Co. (high)', type: 'warn' },
    { t: '12:59', icon: 'globe', text: 'Residential pool usage spike: +34% past hour', type: 'info' },
    { t: '12:32', icon: 'x', text: 'Charge failed: Kenji Watanabe — Visa •• 0914', type: 'danger' },
    { t: '12:21', icon: 'file', text: '9 invoices issued for April billing cycle', type: 'info' },
    { t: '11:58', icon: 'users', text: 'Customer Hans Müller suspended (overdue 21d)', type: 'warn' },
    { t: '11:17', icon: 'wallet', text: 'Wallet top-up of $500.00 from Linh Tran', type: 'ok' },
  ],
};

const STATUS_LABEL = {
  active: 'Active', running: 'Running', paid: 'Paid', open: 'Open',
  pending: 'Pending', overdue: 'Overdue', failed: 'Failed',
  suspended: 'Suspended', stopped: 'Stopped', provisioning: 'Provisioning',
};
const STATUS_BADGE = {
  active: 'ok', running: 'ok', paid: 'ok',
  open: 'info', pending: 'warn', overdue: 'danger', failed: 'danger',
  suspended: '', stopped: '', provisioning: 'info',
};

const fmtMoney = v => {
  const sign = v < 0 ? '-' : '';
  const n = Math.abs(v);
  return sign + '$' + n.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
};
const fmtMoneyShort = v => {
  if (Math.abs(v) >= 1000) return '$' + (v/1000).toFixed(1) + 'k';
  return '$' + v.toFixed(0);
};

window.SAMPLE = SAMPLE;
window.STATUS_LABEL = STATUS_LABEL;
window.STATUS_BADGE = STATUS_BADGE;
window.fmtMoney = fmtMoney;
window.fmtMoneyShort = fmtMoneyShort;
