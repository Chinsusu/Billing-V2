import { getApiData } from "./client";
import {
  AdminAuditLogQuery,
  AdminInvoiceQuery,
  AdminTransactionQuery,
  AuditLog,
  Invoice,
  LedgerEntry,
  Order,
  PaymentReconciliation,
  PaymentTransaction,
  ServiceInstance,
  Wallet,
} from "./types";

export const billingApi = {
  listClientWallets: () => getApiData<Wallet[]>("/client/wallets", "client"),
  listClientWalletLedger: (walletId: string) =>
    getApiData<LedgerEntry[]>(`/client/wallets/${walletId}/ledger`, "client"),
  listClientInvoices: () => getApiData<Invoice[]>("/client/invoices", "client"),
  listClientOrders: () => getApiData<Order[]>("/client/orders", "client"),
  listClientServices: () => getApiData<ServiceInstance[]>("/client/services", "client"),
  listClientTransactions: () => getApiData<PaymentTransaction[]>("/client/transactions", "client"),

  listAdminInvoices: (query: AdminInvoiceQuery = {}) =>
    getApiData<Invoice[]>("/admin/invoices", "admin", query),
  listAdminTransactions: (query: AdminTransactionQuery = {}) =>
    getApiData<PaymentTransaction[]>("/admin/transactions", "admin", query),
  listAdminReconciliation: (query: AdminTransactionQuery = {}) =>
    getApiData<PaymentReconciliation[]>("/admin/payment-reconciliation", "admin", query),
  listAdminAuditLogs: (query: AdminAuditLogQuery = {}) =>
    getApiData<AuditLog[]>("/admin/audit-logs", "admin", query),
};
