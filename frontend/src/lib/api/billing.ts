import { getApiData } from "./client";
import {
  AdminAuditLogQuery,
  AdminInvoiceQuery,
  AdminOrderQuery,
  AdminServiceQuery,
  AdminTransactionQuery,
  AdminWalletQuery,
  AuditLog,
  CatalogPlan,
  CatalogQuery,
  Invoice,
  LedgerEntry,
  LedgerQuery,
  Order,
  PaymentReconciliation,
  PaymentTransaction,
  ServiceInstance,
  TenantCatalog,
  TenantCatalogQuery,
  TopupRequest,
  TopupRequestQuery,
  Wallet,
} from "./types";

export const billingApi = {
  listClientCatalog: (query: TenantCatalogQuery = {}) =>
    getApiData<TenantCatalog>("/client/catalog", "client", query),
  listClientWallets: () => getApiData<Wallet[]>("/client/wallets", "client"),
  listClientWalletLedger: (walletId: string, query: LedgerQuery = {}) =>
    getApiData<LedgerEntry[]>(`/client/wallets/${walletId}/ledger`, "client", query),
  listClientInvoices: () => getApiData<Invoice[]>("/client/invoices", "client"),
  listClientOrders: () => getApiData<Order[]>("/client/orders", "client"),
  listClientServices: () => getApiData<ServiceInstance[]>("/client/services", "client"),
  listClientTransactions: () => getApiData<PaymentTransaction[]>("/client/transactions", "client"),
  listClientTopupRequests: (query: TopupRequestQuery = {}) =>
    getApiData<TopupRequest[]>("/client/topup-requests", "client", query),

  listResellerCatalog: (query: TenantCatalogQuery = {}) =>
    getApiData<TenantCatalog>("/reseller/catalog", "reseller", query),
  listResellerMasterPlans: (query: CatalogQuery = {}) =>
    getApiData<CatalogPlan[]>("/reseller/catalog/master-plans", "reseller", query),

  listAdminInvoices: (query: AdminInvoiceQuery = {}) =>
    getApiData<Invoice[]>("/admin/invoices", "admin", query),
  listAdminOrders: (query: AdminOrderQuery = {}) =>
    getApiData<Order[]>("/admin/orders", "admin", query),
  listAdminServices: (query: AdminServiceQuery = {}) =>
    getApiData<ServiceInstance[]>("/admin/services", "admin", query),
  listAdminWallets: (query: AdminWalletQuery = {}) =>
    getApiData<Wallet[]>("/admin/wallets", "admin", query),
  listAdminWalletLedger: (walletId: string, query: LedgerQuery = {}) =>
    getApiData<LedgerEntry[]>(`/admin/wallets/${walletId}/ledger`, "admin", query),
  listAdminTopupRequests: (query: TopupRequestQuery = {}) =>
    getApiData<TopupRequest[]>("/admin/topup-requests", "admin", query),
  listAdminTransactions: (query: AdminTransactionQuery = {}) =>
    getApiData<PaymentTransaction[]>("/admin/transactions", "admin", query),
  listAdminReconciliation: (query: AdminTransactionQuery = {}) =>
    getApiData<PaymentReconciliation[]>("/admin/payment-reconciliation", "admin", query),
  listAdminAuditLogs: (query: AdminAuditLogQuery = {}) =>
    getApiData<AuditLog[]>("/admin/audit-logs", "admin", query),
};
