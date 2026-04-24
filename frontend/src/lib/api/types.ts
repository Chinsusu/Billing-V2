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
export type ApiJson = unknown;

export interface PageQuery {
  limit?: string | number;
  cursor?: string;
}

export interface Wallet {
  id: string;
  display_id: number;
  tenant_id?: string;
  owner_type?: string;
  owner_id: string;
  currency: string;
  status: string;
  available_balance_minor: number;
  locked_balance_minor: number;
  metadata?: ApiJson;
  created_at?: string;
  updated_at?: string;
}

export interface LedgerEntry {
  id: string;
  display_id: number;
  wallet_id?: string;
  tenant_id?: string;
  direction: string;
  amount_minor: number;
  currency: string;
  entry_type: string;
  status: string;
  balance_after_minor: number;
  reference_type: string;
  reference_id?: string;
  created_by?: string;
  reason?: string;
  correlation_id?: string;
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
  tenant_id?: string;
  buyer_user_id?: string;
  tenant_plan_id: string;
  quantity: number;
  currency: string;
  unit_price_minor?: number;
  discount_minor?: number;
  total_minor: number;
  order_status: string;
  billing_status: string;
  product_snapshot?: ApiJson;
  plan_snapshot?: ApiJson;
  price_snapshot?: ApiJson;
  created_at: string;
  updated_at?: string;
}

export interface CreateClientOrderBody {
  tenant_plan_id: string;
  quantity: number;
  currency: string;
  unit_price_minor: number;
  discount_minor: number;
  total_minor: number;
  product_snapshot?: ApiJson;
  plan_snapshot?: ApiJson;
  price_snapshot?: ApiJson;
}

export interface CheckoutClientOrderBody {
  order_id: string;
}

export interface ServiceInstance {
  id: string;
  display_id: number;
  tenant_id?: string;
  order_id: string;
  tenant_plan_id?: string;
  provider_source_id?: string;
  external_resource_id: string;
  status: string;
  billing_status: string;
  suspension_reason?: string;
  term_start?: string;
  term_end: string;
  created_at?: string;
  updated_at?: string;
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

export interface AdminTenant {
  id: string;
  display_id: number;
  parent_tenant_id?: string;
  tenant_type: string;
  name: string;
  slug: string;
  status: string;
  default_currency: string;
  timezone: string;
  owner_user_id?: string;
  primary_domain?: string;
  user_count: number;
  created_at: string;
  updated_at: string;
}

export interface AdminAccount {
  id: string;
  display_id: number;
  tenant_id: string;
  tenant_name: string;
  tenant_slug: string;
  email: string;
  email_verified_at?: string;
  full_name: string;
  user_type: string;
  status: string;
  two_factor_status: string;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
}

export interface BillingCycle {
  type: string;
  value: number;
}

export interface CatalogProduct {
  id: string;
  display_id: number;
  product_type: string;
  name: string;
  description: string;
  status: string;
  display_order: number;
  created_by?: string;
  created_at: string;
  updated_at: string;
}

export interface CatalogPlan {
  id: string;
  display_id: number;
  product_id: string;
  plan_code: string;
  name: string;
  specs: ApiJson;
  billing_cycle: BillingCycle;
  base_cost_minor: number;
  suggested_price_minor: number;
  reseller_min_price_minor: number;
  currency: string;
  status: string;
  version: number;
  created_at: string;
  updated_at: string;
}

export interface CatalogProviderSource {
  id: string;
  display_id: number;
  source_type: string;
  name: string;
  provider_account_id: string;
  location: string;
  status: string;
  capability_profile?: ApiJson;
  inventory_mode: string;
  risk_level: string;
  created_at: string;
  updated_at: string;
}

export interface TenantCatalogProduct {
  id: string;
  display_id: number;
  tenant_id?: string;
  master_product_id: string;
  name_override: string;
  description_override: string;
  status: string;
  clone_version: number;
  created_at: string;
  updated_at: string;
}

export interface TenantCatalogPlan {
  id: string;
  display_id: number;
  tenant_id?: string;
  tenant_product_id: string;
  master_plan_id: string;
  selling_price_minor: number;
  reseller_cost_minor?: number;
  currency: string;
  margin_policy?: ApiJson;
  visibility: string;
  status: string;
  clone_version: number;
  product_snapshot?: ApiJson;
  plan_snapshot?: ApiJson;
  price_snapshot?: ApiJson;
  capability_snapshot?: ApiJson;
  created_at: string;
  updated_at: string;
}

export interface TenantCatalog {
  products: TenantCatalogProduct[];
  plans: TenantCatalogPlan[];
}

export interface CloneTenantProductBody {
  master_product_id: string;
  name_override?: string;
  description_override?: string;
  status: string;
  clone_version?: number;
}

export interface CloneTenantPlanBody {
  tenant_product_id: string;
  master_plan_id: string;
  selling_price_minor: number;
  reseller_cost_minor: number;
  currency: string;
  margin_policy?: ApiJson;
  visibility: string;
  status: string;
  clone_version?: number;
  product_snapshot?: ApiJson;
  plan_snapshot?: ApiJson;
  price_snapshot?: ApiJson;
  capability_snapshot?: ApiJson;
}

export interface TopupRequest {
  id: string;
  display_id: number;
  tenant_id: string;
  wallet_id: string;
  requested_by: string;
  amount_minor: number;
  currency: string;
  payment_method: string;
  payment_reference?: string;
  status: string;
  reviewed_by?: string;
  reviewed_at?: string;
  review_note?: string;
  ledger_entry_id?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTopupRequestBody {
  wallet_id: string;
  amount_minor: number;
  currency: string;
  payment_method: string;
  payment_reference?: string;
}

export interface InvoiceWalletPaymentBody {
  invoice_id: string;
  wallet_id: string;
}

export interface InvoiceWalletPayment {
  invoice: {
    id: string;
    display_id: number;
    status: string;
    total_minor: number;
    currency: string;
    paid_at?: string;
  };
  transaction: PaymentTransaction;
  order?: {
    id: string;
    display_id: number;
    order_status: string;
    billing_status: string;
  };
  ledger?: {
    id: string;
    display_id: number;
    wallet_id: string;
    direction: string;
    entry_type: string;
    status: string;
    currency: string;
    amount_minor: number;
    balance_after_minor: number;
  };
}

export interface TopupReviewBody {
  review_note?: string;
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

export interface AdminTenantQuery extends PageQuery {
  display_id?: string;
  parent_tenant_id?: string;
  type?: string;
  status?: string;
}

export interface AdminAccountQuery extends PageQuery {
  display_id?: string;
  type?: string;
  status?: string;
  email?: string;
}

export interface AdminOrderQuery extends PageQuery {
  buyer_user_id?: string;
  display_id?: string;
  status?: string;
  billing_status?: string;
  amount_min?: string;
  amount_max?: string;
}

export interface AdminServiceQuery extends PageQuery {
  buyer_user_id?: string;
  display_id?: string;
  order_id?: string;
  order_display_id?: string;
  status?: string;
}

export interface AdminWalletQuery extends PageQuery {
  display_id?: string;
  owner_type?: string;
  owner_id?: string;
  status?: string;
}

export interface LedgerQuery extends PageQuery {
  display_id?: string;
  direction?: string;
  entry_type?: string;
  status?: string;
  amount_min?: string;
  amount_max?: string;
}

export interface TopupRequestQuery extends PageQuery {
  wallet_id?: string;
  requested_by?: string;
  display_id?: string;
  payment_method?: string;
  status?: string;
  amount_min?: string;
  amount_max?: string;
}

export interface CatalogQuery extends PageQuery {
  product_type?: string;
  status?: string;
}

export interface AdminProviderSourceQuery extends PageQuery {
  source_type?: string;
  status?: string;
}

export interface TenantCatalogQuery extends CatalogQuery {
  visibility?: string;
}
