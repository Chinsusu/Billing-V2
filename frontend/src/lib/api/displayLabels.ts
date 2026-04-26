const ACCOUNT_TYPE_LABELS: Record<string, string> = {
  admin: "Admin",
  client: "Client",
  provider_webhook: "Provider webhook",
  reseller: "Reseller",
  reseller_owner: "Reseller owner",
  user: "User",
  worker: "Worker",
};

const INVENTORY_MODE_LABELS: Record<string, string> = {
  manual: "Manual",
  provider_live: "Provider live",
  stock_pool: "Stock pool",
};

const PRODUCT_TYPE_LABELS: Record<string, string> = {
  bandwidth: "Bandwidth",
  datacenter: "Datacenter",
  isp: "ISP",
  mobile: "Mobile",
  proxy: "Proxy",
  residential: "Residential",
  vps: "VPS",
  "vps-linux": "VPS Linux",
  "vps-win": "VPS Windows",
};

const PROVIDER_SOURCE_TYPE_LABELS: Record<string, string> = {
  hetzner: "Hetzner",
  manual: "Manual pool",
  proxmox: "Proxmox",
  "self-host": "Self-hosted",
  upstream: "Upstream",
};

const COMMON_WORD_LABELS: Record<string, string> = {
  api: "API",
  gb: "GB",
  id: "ID",
  ip: "IP",
  isp: "ISP",
  ui: "UI",
  vps: "VPS",
};

export function accountTypeLabel(type: string): string {
  return labelFromKey(type, ACCOUNT_TYPE_LABELS);
}

export function inventoryModeLabel(mode: string): string {
  return labelFromKey(mode, INVENTORY_MODE_LABELS);
}

export function productTypeLabel(type: string): string {
  return labelFromKey(type, PRODUCT_TYPE_LABELS);
}

export function providerSourceTypeLabel(type: string): string {
  return labelFromKey(type, PROVIDER_SOURCE_TYPE_LABELS);
}

export function riskLevelLabel(level: string): string {
  return labelFromKey(level);
}

export function tenantTypeLabel(type: string): string {
  return accountTypeLabel(type);
}

function labelFromKey(value: string, knownLabels: Record<string, string> = {}): string {
  const normalized = value.trim();
  if (!normalized) return "-";
  const knownLabel = knownLabels[normalized.toLowerCase()];
  if (knownLabel) return knownLabel;

  return normalized
    .split(/[._-]/)
    .filter(Boolean)
    .map((part) => COMMON_WORD_LABELS[part.toLowerCase()] ?? titleCase(part))
    .join(" ");
}

function titleCase(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1);
}
