import { getApiData, newIdempotencyKey, postApiData } from "./client";
import {
  AdminAccount,
  AdminAccountQuery,
  AdminAuditLogQuery,
  AdminInvoiceQuery,
  AdminOrderQuery,
  AdminProviderSourceQuery,
  AdminServiceQuery,
  AdminTenant,
  AdminTenantQuery,
  AdminTransactionQuery,
  AdminWalletQuery,
  AuditLog,
  CatalogPlan,
  CatalogProviderSource,
  CatalogQuery,
  CatalogProduct,
  CloneTenantPlanBody,
  CloneTenantProductBody,
  CreateClientOrderBody,
  CreateTopupRequestBody,
  Invoice,
  InvoiceWalletPayment,
  InvoiceWalletPaymentBody,
  LedgerEntry,
  LedgerQuery,
  Order,
  PaymentReconciliation,
  PaymentTransaction,
  ServiceInstance,
  TenantCatalog,
  TenantCatalogQuery,
  TenantCatalogPlan,
  TenantCatalogProduct,
  TopupRequest,
  TopupRequestQuery,
  TopupReviewBody,
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
  createClientOrder: (body: CreateClientOrderBody, idempotencyKey = newIdempotencyKey("client-order")) =>
    postApiData<Order>("/client/orders", "client", body, { idempotencyKey }),
  createClientTopupRequest: (
    body: CreateTopupRequestBody,
    idempotencyKey = newIdempotencyKey("client-topup"),
  ) => postApiData<TopupRequest>("/client/topup-requests", "client", body, { idempotencyKey }),
  payClientInvoiceFromWallet: (
    body: InvoiceWalletPaymentBody,
    idempotencyKey = newIdempotencyKey("client-invoice-wallet-payment"),
  ) => postApiData<InvoiceWalletPayment>("/client/invoice-wallet-payments", "client", body, { idempotencyKey }),

  listResellerCatalog: (query: TenantCatalogQuery = {}) =>
    getApiData<TenantCatalog>("/reseller/catalog", "reseller", query),
  listResellerMasterPlans: (query: CatalogQuery = {}) =>
    getApiData<CatalogPlan[]>("/reseller/catalog/master-plans", "reseller", query),
  listResellerCustomers: (query: AdminAccountQuery = {}) =>
    getApiData<AdminAccount[]>("/reseller/customers", "reseller", query),
  listResellerOrders: (query: AdminOrderQuery = {}) =>
    getApiData<Order[]>("/reseller/orders", "reseller", query),
  listResellerServices: (query: AdminServiceQuery = {}) =>
    getApiData<ServiceInstance[]>("/reseller/services", "reseller", query),
  listResellerInvoices: (query: AdminInvoiceQuery = {}) =>
    getApiData<Invoice[]>("/reseller/invoices", "reseller", query),
  listResellerTransactions: (query: AdminTransactionQuery = {}) =>
    getApiData<PaymentTransaction[]>("/reseller/transactions", "reseller", query),
  listResellerWallets: (query: AdminWalletQuery = {}) =>
    getApiData<Wallet[]>("/reseller/wallets", "reseller", query),
  listResellerWalletLedger: (walletId: string, query: LedgerQuery = {}) =>
    getApiData<LedgerEntry[]>(`/reseller/wallets/${walletId}/ledger`, "reseller", query),
  listResellerTopupRequests: (query: TopupRequestQuery = {}) =>
    getApiData<TopupRequest[]>("/reseller/topup-requests", "reseller", query),
  cloneResellerCatalogProduct: (body: CloneTenantProductBody) =>
    postApiData<TenantCatalogProduct>("/reseller/catalog/products/clone", "reseller", body),
  cloneResellerCatalogPlan: (body: CloneTenantPlanBody) =>
    postApiData<TenantCatalogPlan>("/reseller/catalog/plans/clone", "reseller", body),

  listAdminInvoices: (query: AdminInvoiceQuery = {}) =>
    getApiData<Invoice[]>("/admin/invoices", "admin", query),
  listAdminTenants: (query: AdminTenantQuery = {}) =>
    getApiData<AdminTenant[]>("/admin/tenants", "admin", query),
  listAdminAccounts: (query: AdminAccountQuery = {}) =>
    getApiData<AdminAccount[]>("/admin/accounts", "admin", query),
  listAdminCustomers: (query: AdminAccountQuery = {}) =>
    getApiData<AdminAccount[]>("/admin/customers", "admin", query),
  listAdminCatalogProducts: (query: CatalogQuery = {}) =>
    getApiData<CatalogProduct[]>("/admin/catalog/products", "admin", query),
  listAdminCatalogPlans: (query: CatalogQuery = {}) =>
    getApiData<CatalogPlan[]>("/admin/catalog/plans", "admin", query),
  listAdminProviderSources: (query: AdminProviderSourceQuery = {}) =>
    getApiData<CatalogProviderSource[]>("/admin/catalog/provider-sources", "admin", query),
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
  approveAdminTopupRequest: (id: string, body: TopupReviewBody = {}) =>
    postApiData<TopupRequest>(`/admin/topup-requests/${encodeURIComponent(id)}/approve`, "admin", body),
  rejectAdminTopupRequest: (id: string, body: TopupReviewBody = {}) =>
    postApiData<TopupRequest>(`/admin/topup-requests/${encodeURIComponent(id)}/reject`, "admin", body),
  listAdminTransactions: (query: AdminTransactionQuery = {}) =>
    getApiData<PaymentTransaction[]>("/admin/transactions", "admin", query),
  listAdminReconciliation: (query: AdminTransactionQuery = {}) =>
    getApiData<PaymentReconciliation[]>("/admin/payment-reconciliation", "admin", query),
  listAdminAuditLogs: (query: AdminAuditLogQuery = {}) =>
    getApiData<AuditLog[]>("/admin/audit-logs", "admin", query),
};
