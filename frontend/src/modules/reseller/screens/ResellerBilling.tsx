"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { fulfillmentForOrder } from "@/lib/api/fulfillment";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
import { resellerAccountLabel } from "@/lib/api/resellerViewModels";
import { useApiResource } from "@/lib/api/useApiResource";
import { INVOICES, TRANSACTIONS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

type BillingKind = "invoices" | "transactions";

interface ResellerBillingProps {
  kind: BillingKind;
}

export function ResellerBilling({ kind }: ResellerBillingProps) {
  return kind === "invoices" ? <ResellerInvoices /> : <ResellerTransactions />;
}

function ResellerInvoices() {
  const invoices = useApiResource(billingApi.listResellerInvoices, "reseller-invoices");
  const customers = useApiResource(
    () => billingApi.listResellerCustomers({ limit: 100 }),
    "reseller-invoice-customers",
  );
  const orders = useApiResource(
    () => billingApi.listResellerOrders({ limit: 100 }),
    "reseller-invoice-orders",
  );
  const services = useApiResource(
    () => billingApi.listResellerServices({ limit: 100 }),
    "reseller-invoice-services",
  );
  const jobs = useApiResource(
    () => billingApi.listResellerJobs({ job_type: "provider.provision", limit: 100 }),
    "reseller-invoice-jobs",
  );
  const customerByID = new Map((customers.data ?? []).map((customer) => [customer.id, customer]));
  const customerByDisplayID = new Map((customers.data ?? []).map((customer) => [customer.display_id, customer]));
  const orderByID = new Map((orders.data ?? []).map((order) => [order.id, order]));
  const usingLive = invoices.status === "success";
  const rows = usingLive
    ? (invoices.data ?? []).map((invoice) => {
        const customer = customerByDisplayID.get(invoice.buyer_display_id ?? 0) ?? customerByID.get(invoice.buyer_user_id);
        const buyerDisplayID = invoice.buyer_display_id ?? customer?.display_id;
        const fulfillment = fulfillmentForOrder(invoice.order_id ? orderByID.get(invoice.order_id) : undefined, services.data ?? [], {
          jobs: jobs.data ?? [],
          jobsUnavailable: jobs.status === "error",
        });
        return {
          id: recordLabel(invoice.display_id, "INV-"),
          order: fulfillment.orderLabel,
          customer: resellerAccountLabel(buyerDisplayID, customer),
          issued: compactDateTime(invoice.issued_at),
          due: compactDateTime(invoice.due_at),
          amount: moneyMinor(invoice.total_minor, invoice.currency),
          amountMinor: invoice.total_minor,
          status: invoice.status,
          fulfillment,
        };
      })
    : INVOICES.slice(0, 6).map((invoice) => ({
        id: invoice.id,
        order: "-",
        customer: invoice.customer,
        issued: invoice.issued,
        due: invoice.due,
        amount: fmtMoney(invoice.amount),
        amountMinor: Math.round(invoice.amount * 100),
        status: invoice.status,
        fulfillment: { status: invoice.status, label: invoice.status, serviceLabel: "-", jobLabel: "-", orderLabel: "-" },
      }));
  const open = rows.filter((invoice) => invoice.status !== "paid").length;
  const total = rows.reduce((sum, invoice) => sum + invoice.amountMinor, 0);
  const extraError = orders.error ?? services.error ?? customers.error ?? jobs.error;

  return (
    <BillingShell title="Invoices" records={rows.length} open={open} total={moneyMinor(total)} source={sourceText(invoices.status, usingLive, extraError)}>
      <table className="w-full text-[13px] border-collapse min-w-[920px]">
        <thead>
          <tr className="bg-gray-50">
            {["Invoice", "Order", "Client", "Issued", "Due", "Amount", "Invoice", "Fulfillment", "Service / Job"].map((heading) => (
              <th key={heading} className="text-left text-[11px] font-medium uppercase text-gray-400 p-4 border-b border-gray-200">
                {heading}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((invoice) => (
            <tr key={invoice.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
              <td className="p-4 text-[12px] text-[#D50C2D] font-medium">{invoice.id}</td>
              <td className="p-4 text-[12px] text-gray-500">{invoice.order}</td>
              <td className="p-4 font-medium text-gray-900">{invoice.customer}</td>
              <td className="p-4 text-gray-500">{invoice.issued}</td>
              <td className="p-4 text-gray-500">{invoice.due}</td>
              <td className="p-4 text-right font-medium tabular-nums">{invoice.amount}</td>
              <td className="p-4"><StatusBadge status={invoice.status} dot /></td>
              <td className="p-4"><StatusBadge status={invoice.fulfillment.status} dot /></td>
              <td className="p-4 text-[12px] text-gray-500">
                {invoice.fulfillment.serviceLabel !== "-" ? invoice.fulfillment.serviceLabel : invoice.fulfillment.jobLabel}
              </td>
            </tr>
          ))}
          {usingLive && rows.length === 0 && (
            <tr><td colSpan={9} className="p-4 text-center text-[12px] text-gray-400">No invoices</td></tr>
          )}
        </tbody>
      </table>
    </BillingShell>
  );
}

function ResellerTransactions() {
  const transactions = useApiResource(billingApi.listResellerTransactions, "reseller-transactions");
  const customers = useApiResource(
    () => billingApi.listResellerCustomers({ limit: 100 }),
    "reseller-transaction-customers",
  );
  const orders = useApiResource(
    () => billingApi.listResellerOrders({ limit: 100 }),
    "reseller-transaction-orders",
  );
  const services = useApiResource(
    () => billingApi.listResellerServices({ limit: 100 }),
    "reseller-transaction-services",
  );
  const jobs = useApiResource(
    () => billingApi.listResellerJobs({ job_type: "provider.provision", limit: 100 }),
    "reseller-transaction-jobs",
  );
  const customerByID = new Map((customers.data ?? []).map((customer) => [customer.id, customer]));
  const customerByDisplayID = new Map((customers.data ?? []).map((customer) => [customer.display_id, customer]));
  const orderByID = new Map((orders.data ?? []).map((order) => [order.id, order]));
  const usingLive = transactions.status === "success";
  const rows = usingLive
    ? (transactions.data ?? []).map((transaction) => {
        const customer = customerByDisplayID.get(transaction.account_display_id ?? 0) ?? customerByID.get(transaction.account_user_id);
        const accountDisplayID = transaction.account_display_id ?? customer?.display_id;
        const fulfillment = fulfillmentForOrder(transaction.order_id ? orderByID.get(transaction.order_id) : undefined, services.data ?? [], {
          jobs: jobs.data ?? [],
          jobsUnavailable: jobs.status === "error",
        });
        return {
          id: recordLabel(transaction.display_id, "TX-"),
          time: compactDateTime(transaction.created_at),
          customer: resellerAccountLabel(accountDisplayID, customer),
          order: fulfillment.orderLabel,
          method: transaction.description ?? "wallet",
          type: transaction.type,
          amount: moneyMinor(transaction.amount_minor, transaction.currency),
          amountMinor: transaction.amount_minor,
          status: transaction.status,
          fulfillment,
        };
      })
    : TRANSACTIONS.slice(0, 8).map((transaction) => ({
        id: transaction.id,
        time: transaction.time,
        customer: transaction.customer,
        order: "-",
        method: transaction.method,
        type: transaction.type,
        amount: fmtMoney(transaction.amount),
        amountMinor: Math.round(transaction.amount * 100),
        status: transaction.status,
        fulfillment: { status: transaction.status, label: transaction.status, serviceLabel: "-", jobLabel: "-", orderLabel: "-" },
      }));
  const failed = rows.filter((txn) => txn.status === "failed").length;
  const total = rows.reduce((sum, txn) => sum + txn.amountMinor, 0);
  const extraError = orders.error ?? services.error ?? customers.error ?? jobs.error;

  return (
    <BillingShell title="Transactions" records={rows.length} open={failed} total={moneyMinor(total)} openLabel="Failed" source={sourceText(transactions.status, usingLive, extraError)}>
      <table className="w-full text-[13px] border-collapse min-w-[960px]">
        <thead>
          <tr className="bg-gray-50">
            {["Transaction", "Time", "Client", "Order", "Method", "Type", "Amount", "Payment", "Fulfillment"].map((heading) => (
              <th key={heading} className="text-left text-[11px] font-medium uppercase text-gray-400 p-4 border-b border-gray-200">
                {heading}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((txn) => (
            <tr key={txn.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
              <td className="p-4 text-[12px] text-[#D50C2D] font-medium">{txn.id}</td>
              <td className="p-4 text-gray-500">{txn.time}</td>
              <td className="p-4 font-medium text-gray-900">{txn.customer}</td>
              <td className="p-4 text-[12px] text-gray-500">{txn.order}</td>
              <td className="p-4 text-gray-500">{txn.method}</td>
              <td className="p-4 text-gray-500">{txn.type}</td>
              <td className="p-4 text-right font-medium tabular-nums">{txn.amount}</td>
              <td className="p-4"><StatusBadge status={txn.status} dot /></td>
              <td className="p-4"><StatusBadge status={txn.fulfillment.status} dot /></td>
            </tr>
          ))}
          {usingLive && rows.length === 0 && (
            <tr><td colSpan={9} className="p-4 text-center text-[12px] text-gray-400">No transactions</td></tr>
          )}
        </tbody>
      </table>
    </BillingShell>
  );
}

function BillingShell({
  title,
  records,
  open,
  total,
  source,
  openLabel = "Open",
  children,
}: {
  title: string;
  records: number;
  open: number;
  total: string;
  source: string;
  openLabel?: string;
  children: React.ReactNode;
}) {
  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <SummaryTile label="Records" value={String(records)} />
        <SummaryTile label={openLabel} value={String(open)} tone={open > 0 ? "warn" : "neutral"} />
        <SummaryTile label="Total" value={total} />
      </div>
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">{title}</h3>
          <span className="text-[11px] text-gray-400">{source}</span>
        </div>
        <div className="overflow-x-auto max-w-full">{children}</div>
      </div>
    </div>
  );
}

function sourceText(status: string, usingLive: boolean, extraError?: string | null) {
  if (status === "error") return "Live API unavailable. Showing demo billing data.";
  if (status === "loading") return "Refreshing live billing data...";
  if (usingLive && extraError) return "Live billing loaded. Fulfillment details may be incomplete.";
  return usingLive ? "Live reseller billing" : "Demo billing data";
}

function SummaryTile({ label, value, tone = "neutral" }: { label: string; value: string; tone?: "neutral" | "warn" }) {
  return (
    <div className={`bg-white border rounded p-4 ${tone === "warn" ? "border-amber-200" : "border-gray-200"}`}>
      <div className="text-[11px] text-gray-400 uppercase mb-1">{label}</div>
      <div className={`text-lg font-medium tabular-nums ${tone === "warn" ? "text-amber-700" : "text-gray-900"}`}>{value}</div>
    </div>
  );
}
