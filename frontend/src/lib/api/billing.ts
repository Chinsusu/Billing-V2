import { getApiData } from "./client";
import {
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

  listAdminInvoices: () => getApiData<Invoice[]>("/admin/invoices", "admin"),
  listAdminTransactions: () => getApiData<PaymentTransaction[]>("/admin/transactions", "admin"),
  listAdminReconciliation: () =>
    getApiData<PaymentReconciliation[]>("/admin/payment-reconciliation", "admin"),
  listAdminAuditLogs: () => getApiData<AuditLog[]>("/admin/audit-logs", "admin"),
};
