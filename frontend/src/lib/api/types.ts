export interface ApiEnvelope<T> {
  data: T;
  request_id: string;
  page?: {
    limit: number;
    next_cursor?: string;
  };
}

export type ApiQueryValue = string | number | null | undefined;
export type ApiQuery = object;

export interface Wallet {
  id: string;
  display_id: number;
  owner_id: string;
  currency: string;
  status: string;
  available_balance_minor: number;
  locked_balance_minor: number;
}

export interface LedgerEntry {
  id: string;
  display_id: number;
  direction: string;
  amount_minor: number;
  currency: string;
  entry_type: string;
  status: string;
  balance_after_minor: number;
  reference_type: string;
  created_at: string;
}

export interface Invoice {
  id: string;
  display_id: number;
  buyer_user_id: string;
  order_id?: string;
  status: string;
  currency: string;
  total_minor: number;
  issued_at?: string;
  due_at?: string;
  paid_at?: string;
}

export interface PaymentTransaction {
  id: string;
  display_id: number;
  account_user_id: string;
  order_id?: string;
  invoice_id?: string;
  type: string;
  status: string;
  currency: string;
  amount_minor: number;
  description?: string;
  created_at: string;
}

export interface PaymentReconciliation {
  transaction: PaymentTransaction;
  provider?: string;
  invoice?: {
    id: string;
    display_id: number;
    status: string;
    total_minor: number;
  };
  ledger?: {
    id: string;
    display_id: number;
    wallet_display_id?: number;
    direction: string;
    entry_type: string;
    status: string;
  };
}

export interface Order {
  id: string;
  display_id: number;
  tenant_plan_id: string;
  quantity: number;
  currency: string;
  total_minor: number;
  order_status: string;
  billing_status: string;
  created_at: string;
}

export interface ServiceInstance {
  id: string;
  display_id: number;
  order_id: string;
  external_resource_id: string;
  status: string;
  billing_status: string;
  term_end: string;
}

export interface AuditLog {
  id: string;
  display_id: number;
  actor_id?: string;
  actor_type: string;
  action: string;
  target_type: string;
  target_id: string;
  correlation_id: string;
  created_at: string;
}

export interface AdminInvoiceQuery {
  display_id?: string;
  buyer_user_id?: string;
  status?: string;
  amount_min?: string;
  amount_max?: string;
}

export interface AdminTransactionQuery {
  display_id?: string;
  account_user_id?: string;
  status?: string;
  amount_min?: string;
  amount_max?: string;
}

export interface AdminAuditLogQuery {
  display_id?: string;
  actor_id?: string;
  action?: string;
  target_type?: string;
}
