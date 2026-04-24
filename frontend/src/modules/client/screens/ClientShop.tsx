"use client";

import { useMemo, useState } from "react";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import type { Invoice, InvoiceWalletPayment, Order, TenantCatalogPlan, Wallet } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { PRODUCTS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

type Notice = { type: "success" | "error"; text: string };
type CheckoutSummary = {
  orderId: string;
  orderDisplayID: number;
  orderStatus: string;
  billingStatus: string;
  invoiceId: string;
  invoiceDisplayID: number;
  invoiceStatus: string;
  amountMinor: number;
  currency: string;
  transactionDisplayID?: number;
};

const PAYABLE_INVOICE_STATUSES = new Set(["issued", "overdue"]);

function snapshotValue(value: unknown, keys: string[]): string | null {
  if (!value || typeof value !== "object" || Array.isArray(value)) return null;
  const record = value as Record<string, unknown>;
  for (const key of keys) {
    const data = record[key];
    if (typeof data === "string" && data.trim()) return data;
  }
  return null;
}

function planTitle(plan: TenantCatalogPlan): string {
  return snapshotValue(plan.plan_snapshot, ["name", "plan_code"]) ?? recordLabel(plan.display_id, "PLAN-");
}

function productTitle(plan: TenantCatalogPlan): string {
  return snapshotValue(plan.product_snapshot, ["name", "product_name"]) ?? "Catalog plan";
}

function orderBody(plan: TenantCatalogPlan) {
  return {
    tenant_plan_id: plan.id,
    quantity: 1,
    currency: plan.currency,
    unit_price_minor: plan.selling_price_minor,
    discount_minor: 0,
    total_minor: plan.selling_price_minor,
    product_snapshot: plan.product_snapshot ?? { tenant_product_id: plan.tenant_product_id },
    plan_snapshot: plan.plan_snapshot ?? { tenant_plan_id: plan.id, master_plan_id: plan.master_plan_id },
    price_snapshot: plan.price_snapshot ?? {
      selling_price_minor: plan.selling_price_minor,
      currency: plan.currency,
    },
  };
}

function canPayInvoice(invoice: Invoice): boolean {
  return PAYABLE_INVOICE_STATUSES.has(invoice.status);
}

function orderLabelForID(orderId: string | undefined, ordersByID: Map<string, Order>): string {
  if (!orderId) return "-";
  const order = ordersByID.get(orderId);
  return order ? recordLabel(order.display_id, "ORD-") : "-";
}

function summaryFromCheckout(order: Order, invoice: Invoice): CheckoutSummary {
  return {
    orderId: order.id,
    orderDisplayID: order.display_id,
    orderStatus: order.order_status,
    billingStatus: order.billing_status,
    invoiceId: invoice.id,
    invoiceDisplayID: invoice.display_id,
    invoiceStatus: invoice.status,
    amountMinor: invoice.total_minor,
    currency: invoice.currency,
  };
}

function summaryAfterPayment(current: CheckoutSummary | null, result: InvoiceWalletPayment): CheckoutSummary {
  return {
    orderId: result.order?.id ?? result.transaction.order_id ?? current?.orderId ?? "",
    orderDisplayID: result.order?.display_id ?? current?.orderDisplayID ?? 0,
    orderStatus: result.order?.order_status ?? current?.orderStatus ?? "paid",
    billingStatus: result.order?.billing_status ?? current?.billingStatus ?? "paid",
    invoiceId: result.invoice.id,
    invoiceDisplayID: result.invoice.display_id,
    invoiceStatus: result.invoice.status,
    amountMinor: result.invoice.total_minor,
    currency: result.invoice.currency,
    transactionDisplayID: result.transaction.display_id,
  };
}

function CheckoutStatusPanel({
  summary,
  liveOrder,
  serviceCount,
  transactionCount,
  walletBalance,
}: {
  summary: CheckoutSummary;
  liveOrder?: Order | null;
  serviceCount: number;
  transactionCount: number;
  walletBalance: string;
}) {
  const orderStatus = liveOrder?.order_status ?? summary.orderStatus;
  const billingStatus = liveOrder?.billing_status ?? summary.billingStatus;
  const orderDisplayID = liveOrder?.display_id ?? summary.orderDisplayID;
  const orderLabel = orderDisplayID > 0 ? recordLabel(orderDisplayID, "ORD-") : "-";

  return (
    <div className="bg-white border border-gray-200 rounded">
      <div className="p-4 border-b border-gray-100 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Checkout status</h3>
          <div className="text-[11px] text-gray-400 mt-0.5">{walletBalance} wallet balance</div>
        </div>
        <div className="flex items-center gap-2">
          <StatusBadge status={orderStatus} dot />
          <StatusBadge status={billingStatus} dot />
        </div>
      </div>
      <div className="grid grid-cols-2 lg:grid-cols-4 divide-x divide-gray-100">
        <FlowMetric label="Order" value={orderLabel} meta={orderStatus} />
        <FlowMetric
          label="Invoice"
          value={recordLabel(summary.invoiceDisplayID, "INV-")}
          meta={`${summary.invoiceStatus} / ${moneyMinor(summary.amountMinor, summary.currency)}`}
        />
        <FlowMetric
          label="Payment"
          value={summary.transactionDisplayID ? recordLabel(summary.transactionDisplayID, "TXN-") : "-"}
          meta={`${transactionCount} transactions`}
        />
        <FlowMetric label="Provisioning" value={String(serviceCount)} meta={serviceCount === 1 ? "service" : "services"} />
      </div>
    </div>
  );
}

function FlowMetric({ label, value, meta }: { label: string; value: string; meta: string }) {
  return (
    <div className="p-4 min-w-0">
      <div className="text-[11px] text-gray-400 uppercase mb-1">{label}</div>
      <div className="text-[14px] font-medium text-gray-900 tabular-nums truncate">{value}</div>
      <div className="text-[11px] text-gray-400 mt-1 truncate">{meta}</div>
    </div>
  );
}

export function ClientShop() {
  const [refreshKey, setRefreshKey] = useState(0);
  const [busyAction, setBusyAction] = useState<string | null>(null);
  const [notice, setNotice] = useState<Notice | null>(null);
  const [checkoutSummary, setCheckoutSummary] = useState<CheckoutSummary | null>(null);
  const catalog = useApiResource(
    () => billingApi.listClientCatalog({ limit: 100, status: "active", visibility: "public" }),
    `client-catalog:${refreshKey}`,
  );
  const invoices = useApiResource(billingApi.listClientInvoices, `client-invoices:${refreshKey}`);
  const orders = useApiResource(billingApi.listClientOrders, `client-orders:${refreshKey}`);
  const transactions = useApiResource(billingApi.listClientTransactions, `client-transactions:${refreshKey}`);
  const services = useApiResource(billingApi.listClientServices, `client-services:${refreshKey}`);
  const wallets = useApiResource(billingApi.listClientWallets, `client-wallets:${refreshKey}`);
  const livePlans = catalog.status === "success" ? (catalog.data?.plans ?? []) : null;
  const liveInvoices = invoices.status === "success" ? (invoices.data ?? []) : null;
  const wallet = wallets.data?.[0];
  const ordersByID = useMemo(() => new Map((orders.data ?? []).map((order) => [order.id, order])), [orders.data]);
  const latestSummaryOrder = checkoutSummary?.orderId ? ordersByID.get(checkoutSummary.orderId) : null;
  const serviceCountForSummary = checkoutSummary
    ? (services.data ?? []).filter((service) => service.order_id === checkoutSummary.orderId).length
    : 0;
  const loading = catalog.status === "loading" || invoices.status === "loading" || orders.status === "loading";
  const liveError = catalog.error ?? invoices.error ?? orders.error ?? transactions.error ?? wallets.error ?? services.error;

  async function handleOrder(plan: TenantCatalogPlan) {
    setBusyAction(`order:${plan.id}`);
    setNotice(null);
    try {
      const order = await billingApi.createClientOrder(orderBody(plan));
      const invoice = await billingApi.checkoutClientOrder({ order_id: order.id });
      setCheckoutSummary(summaryFromCheckout(order, invoice));
      setNotice({
        type: "success",
        text: `Order ${recordLabel(order.display_id, "ORD-")} created. Invoice ${recordLabel(invoice.display_id, "INV-")} issued.`,
      });
      setRefreshKey((current) => current + 1);
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : "Checkout failed.";
      setNotice({ type: "error", text: message });
    } finally {
      setBusyAction(null);
    }
  }

  async function handlePay(invoice: Invoice, activeWallet?: Wallet) {
    if (!activeWallet) {
      setNotice({ type: "error", text: "No wallet is available for payment." });
      return;
    }
    setBusyAction(`pay:${invoice.id}`);
    setNotice(null);
    try {
      const result = await billingApi.payClientInvoiceFromWallet({
        invoice_id: invoice.id,
        wallet_id: activeWallet.id,
      });
      setCheckoutSummary((current) => summaryAfterPayment(current, result));
      const orderText = result.order ? ` Order ${recordLabel(result.order.display_id, "ORD-")} is ${result.order.order_status}.` : "";
      setNotice({ type: "success", text: `Invoice ${recordLabel(result.invoice.display_id, "INV-")} paid.${orderText}` });
      setRefreshKey((current) => current + 1);
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : "Payment failed.";
      setNotice({ type: "error", text: message });
    } finally {
      setBusyAction(null);
    }
  }

  return (
    <div className="p-4 flex flex-col gap-4">
      {checkoutSummary && (
        <CheckoutStatusPanel
          summary={checkoutSummary}
          liveOrder={latestSummaryOrder}
          serviceCount={serviceCountForSummary}
          transactionCount={transactions.data?.length ?? 0}
          walletBalance={wallet ? moneyMinor(wallet.available_balance_minor, wallet.currency) : "-"}
        />
      )}

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
          <div>
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Plans</h3>
            <div className="text-[11px] text-gray-400 mt-0.5">
              {livePlans ? `${livePlans.length} available` : "Demo catalog"}
            </div>
          </div>
          <div className="flex flex-wrap items-center justify-end gap-3">
            {notice && (
              <span className={`text-[11px] font-medium ${notice.type === "error" ? "text-red-600" : "text-green-700"}`}>
                {notice.text}
              </span>
            )}
            {liveError && <span className="text-[11px] text-amber-600">Live data unavailable.</span>}
          </div>
        </div>
        {livePlans ? (
          <div className="overflow-x-auto max-w-full">
            <table className="w-full text-[13px] border-collapse min-w-[760px]">
              <thead>
                <tr className="bg-gray-50">
                  {["Plan", "Product", "Price", "Status", "Action"].map((h) => (
                    <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {livePlans.map((plan) => (
                  <tr key={plan.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                    <td className="p-4">
                      <div className="font-medium text-gray-900">{planTitle(plan)}</div>
                      <div className="text-[11px] text-gray-400">{recordLabel(plan.display_id, "PLAN-")}</div>
                    </td>
                    <td className="p-4 text-gray-500">{productTitle(plan)}</td>
                    <td className="p-4 tabular-nums text-right font-medium">{moneyMinor(plan.selling_price_minor, plan.currency)}</td>
                    <td className="p-4"><StatusBadge status={plan.status} dot /></td>
                    <td className="p-4">
                      <button
                        className="inline-flex h-8 items-center justify-center rounded-md border border-[#D50C2D] bg-[#D50C2D] px-3 text-[12px] font-medium text-white transition-colors hover:bg-[#B3082A] disabled:cursor-not-allowed disabled:border-gray-200 disabled:bg-gray-100 disabled:text-gray-400"
                        disabled={busyAction === `order:${plan.id}` || plan.status !== "active" || plan.selling_price_minor <= 0}
                        onClick={() => handleOrder(plan)}
                      >
                        {busyAction === `order:${plan.id}` ? "Ordering" : "Order"}
                      </button>
                    </td>
                  </tr>
                ))}
                {!loading && livePlans.length === 0 && (
                  <tr><td colSpan={5} className="p-4 text-center text-[12px] text-gray-400">No plans</td></tr>
                )}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="grid grid-cols-3 gap-4 p-4">
            {PRODUCTS.map((product) => (
              <div key={product.sku} className="bg-white border border-gray-200 rounded p-4 flex flex-col gap-4">
                <div>
                  <div className="text-[13px] font-medium text-gray-900">{product.name}</div>
                  <div className="text-[11px] text-gray-400 mt-0.5">{product.sku}</div>
                </div>
                <div className="flex items-baseline gap-1">
                  <span className="text-lg font-medium tabular-nums text-gray-900">{fmtMoney(product.price)}</span>
                  <span className="text-[12px] text-gray-400">{product.unit}</span>
                </div>
                <div className="text-[11px] text-gray-400">{product.active.toLocaleString()} active subscriptions</div>
                <button
                  className="mt-auto w-full inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-gray-100 text-gray-400 rounded-md border-0 cursor-not-allowed"
                  disabled
                >
                  Order
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <div>
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Invoices</h3>
            <div className="text-[11px] text-gray-400 mt-0.5">
              {wallet ? `${moneyMinor(wallet.available_balance_minor, wallet.currency)} available` : "No wallet loaded"}
            </div>
          </div>
          <span className="text-[11px] text-gray-400">{liveInvoices?.length ?? 0} records</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[860px]">
            <thead>
              <tr className="bg-gray-50">
                {["Invoice", "Order", "Issued", "Paid", "Amount", "Status", "Action"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {(liveInvoices ?? []).map((invoice) => (
                <tr key={invoice.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D]">{recordLabel(invoice.display_id, "INV-")}</td>
                  <td className="p-4 text-[12px] text-gray-500">{orderLabelForID(invoice.order_id, ordersByID)}</td>
                  <td className="p-4 text-gray-400">{compactDateTime(invoice.issued_at)}</td>
                  <td className="p-4 text-gray-400">{compactDateTime(invoice.paid_at)}</td>
                  <td className="p-4 text-right font-medium tabular-nums">{moneyMinor(invoice.total_minor, invoice.currency)}</td>
                  <td className="p-4"><StatusBadge status={invoice.status} dot /></td>
                  <td className="p-4">
                    {canPayInvoice(invoice) ? (
                      <button
                        className="inline-flex h-8 items-center justify-center rounded-md border border-gray-200 bg-white px-3 text-[12px] font-medium text-gray-700 transition-colors hover:border-gray-300 hover:bg-gray-50 disabled:cursor-not-allowed disabled:bg-gray-100 disabled:text-gray-400"
                        disabled={busyAction === `pay:${invoice.id}` || !wallet}
                        onClick={() => handlePay(invoice, wallet)}
                      >
                        {busyAction === `pay:${invoice.id}` ? "Paying" : "Pay"}
                      </button>
                    ) : (
                      <span className="text-[12px] text-gray-400">-</span>
                    )}
                  </td>
                </tr>
              ))}
              {liveInvoices && liveInvoices.length === 0 && (
                <tr><td colSpan={7} className="p-4 text-center text-[12px] text-gray-400">No invoices</td></tr>
              )}
              {!liveInvoices && (
                <tr><td colSpan={7} className="p-4 text-center text-[12px] text-gray-400">Invoice data unavailable</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
