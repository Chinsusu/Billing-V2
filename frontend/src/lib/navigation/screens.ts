export type Portal = "admin" | "reseller" | "client";

export interface NavItem {
  id: string;
  label: string;
  portal: Portal;
  section?: string;
  badge?: "danger" | "warn";
  count?: number;
}

export const ADMIN_NAV: NavItem[] = [
  { id: "admin-overview", label: "Overview", portal: "admin" },

  { id: "admin-tenants", label: "Accounts", portal: "admin", section: "Platform", count: 2852 },
  { id: "admin-provisioning", label: "Provisioning queue", portal: "admin", section: "Platform", badge: "danger", count: 6 },
  { id: "admin-topups", label: "Top-up verification", portal: "admin", section: "Platform", badge: "danger", count: 3 },
  { id: "admin-providers", label: "Providers / Sources", portal: "admin", section: "Platform" },

  { id: "admin-tickets", label: "Support tickets", portal: "admin", section: "Customers", badge: "danger", count: 23 },

  { id: "admin-services-proxies", label: "Proxies", portal: "admin", section: "Services", count: 10 },
  { id: "admin-services-vps", label: "VPS", portal: "admin", section: "Services", count: 7 },
  { id: "admin-services-bandwidth", label: "Bandwidth", portal: "admin", section: "Services", count: 7 },

  { id: "admin-invoices", label: "Invoices", portal: "admin", section: "Billing" },
  { id: "admin-transactions", label: "Transactions", portal: "admin", section: "Billing" },
  { id: "admin-products", label: "Products & Pricing", portal: "admin", section: "Billing" },

  { id: "admin-alerts", label: "Alerts", portal: "admin", section: "System", badge: "danger", count: 5 },
  { id: "admin-logs", label: "Audit logs", portal: "admin", section: "System" },

  { id: "admin-settings", label: "Settings", portal: "admin", section: "Settings" },
];

export const RESELLER_NAV: NavItem[] = [
  { id: "reseller-overview", label: "Overview", portal: "reseller" },

  { id: "reseller-clients", label: "Clients", portal: "reseller", section: "Customers", count: 312 },

  { id: "reseller-services", label: "Services", portal: "reseller", section: "Services" },

  { id: "reseller-catalog", label: "Catalog & Pricing", portal: "reseller", section: "Billing" },
  { id: "reseller-orders", label: "Orders", portal: "reseller", section: "Billing" },
  { id: "reseller-wallet", label: "Wallet & Top-up", portal: "reseller", section: "Billing" },
  { id: "reseller-reports", label: "Reports", portal: "reseller", section: "Billing" },

  { id: "reseller-settings", label: "Branding & Settings", portal: "reseller", section: "Settings" },
];

export const CLIENT_NAV: NavItem[] = [
  { id: "client-overview", label: "Overview", portal: "client" },

  { id: "client-services", label: "My services", portal: "client", section: "Services", count: 5 },

  { id: "client-shop", label: "Shop", portal: "client", section: "Billing" },
  { id: "client-wallet", label: "Wallet", portal: "client", section: "Billing" },
  { id: "client-usage", label: "Usage", portal: "client", section: "Billing" },

  { id: "client-support", label: "Support tickets", portal: "client", section: "Settings" },
  { id: "client-settings", label: "Settings", portal: "client", section: "Settings" },
];

export const NAV_BY_PORTAL: Record<Portal, NavItem[]> = {
  admin: ADMIN_NAV,
  reseller: RESELLER_NAV,
  client: CLIENT_NAV,
};

export interface PortalMeta {
  label: string;
  initial: string;
  user: string;
  role: string;
  domain: string;
}

export const PORTAL_META: Record<Portal, PortalMeta> = {
  admin: { label: "HANetwork", initial: "H", user: "Minh Nguyen", role: "Administrator", domain: "billing.hanetwork.vn" },
  reseller: { label: "ProxyVN", initial: "P", user: "ProxyVN", role: "Reseller · T-0042", domain: "proxyvn.io" },
  client: { label: "ProxyVN", initial: "P", user: "Linh Tran", role: "Client · via ProxyVN", domain: "proxyvn.io" },
};
