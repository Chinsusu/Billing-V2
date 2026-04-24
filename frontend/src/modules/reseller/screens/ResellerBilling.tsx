"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { billingApi } from "@/lib/api/billing";
import { compactDateTime, moneyMinor, recordLabel } from "@/lib/api/format";
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
  const customerByID = new Map((customers.data ?? []).map((customer) => [customer.id, customer]));
  const usingLive = invoices.status === "success";
  const rows = usingLive
    ? (invoices.data ?? []).map((invoice) => {
        const customer = customerByID.get(invoice.buyer_user_id);
        return {
          id: recordLabel(invoice.display_id, "INV-"),
          customer: customer ? `${customer.full_name || customer.email} (${recordLabel(customer.display_id, "ACC-")})` : "-",
          issued: compactDateTime(invoice.issued_at),
          due: compactDateTime(invoice.due_at),
          amount: moneyMinor(invoice.total_minor, invoice.currency),
          amountMinor: invoice.total_minor,
          status: invoice.status,
        };
      })
    : INVOICES.slice(0, 6).map((invoice) => ({
        id: invoice.id,
        customer: invoice.customer,
        issued: invoice.issued,
        due: invoice.due,
        amount: fmtMoney(invoice.amount),
        amountMinor: Math.round(invoice.amount * 100),
        status: invoice.status,
      }));
  const open = rows.filter((invoice) => invoice.status !== "paid").length;
  const total = rows.reduce((sum, invoice) => sum + invoice.amountMinor, 0);

  return (
    <BillingShell title="Invoices" records={rows.length} open={open} total={moneyMinor(total)} source={sourceText(invoices.status, usingLive)}>
      <table className="w-full text-[13px] border-collapse min-w-[720px]">
        <thead>
          <tr className="bg-gray-50">
            {["Invoice", "Client", "Issued", "Due", "Amount", "Status"].map((heading) => (
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
              <td className="p-4 font-medium text-gray-900">{invoice.customer}</td>
              <td className="p-4 text-gray-500">{invoice.issued}</td>
              <td className="p-4 text-gray-500">{invoice.due}</td>
              <td className="p-4 text-right font-medium tabular-nums">{invoice.amount}</td>
              <td className="p-4"><StatusBadge status={invoice.status} dot /></td>
            </tr>
          ))}
          {usingLive && rows.length === 0 && (
            <tr><td colSpan={6} className="p-4 text-center text-[12px] text-gray-400">No invoices</td></tr>
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
  const customerByID = new Map((customers.data ?? []).map((customer) => [customer.id, customer]));
  const usingLive = transactions.status === "success";
  const rows = usingLive
    ? (transactions.data ?? []).map((transaction) => {
        const customer = customerByID.get(transaction.account_user_id);
        return {
          id: recordLabel(transaction.display_id, "TX-"),
          time: compactDateTime(transaction.created_at),
          customer: customer ? `${customer.full_name || customer.email} (${recordLabel(customer.display_id, "ACC-")})` : "-",
          method: transaction.description ?? "wallet",
          type: transaction.type,
          amount: moneyMinor(transaction.amount_minor, transaction.currency),
          amountMinor: transaction.amount_minor,
          status: transaction.status,
        };
      })
    : TRANSACTIONS.slice(0, 8).map((transaction) => ({
        id: transaction.id,
        time: transaction.time,
        customer: transaction.customer,
        method: transaction.method,
        type: transaction.type,
        amount: fmtMoney(transaction.amount),
        amountMinor: Math.round(transaction.amount * 100),
        status: transaction.status,
      }));
  const failed = rows.filter((txn) => txn.status === "failed").length;
  const total = rows.reduce((sum, txn) => sum + txn.amountMinor, 0);

  return (
    <BillingShell title="Transactions" records={rows.length} open={failed} total={moneyMinor(total)} openLabel="Failed" source={sourceText(transactions.status, usingLive)}>
      <table className="w-full text-[13px] border-collapse min-w-[760px]">
        <thead>
          <tr className="bg-gray-50">
            {["Transaction", "Time", "Client", "Method", "Type", "Amount", "Status"].map((heading) => (
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
              <td className="p-4 text-gray-500">{txn.method}</td>
              <td className="p-4 text-gray-500">{txn.type}</td>
              <td className="p-4 text-right font-medium tabular-nums">{txn.amount}</td>
              <td className="p-4"><StatusBadge status={txn.status} dot /></td>
            </tr>
          ))}
          {usingLive && rows.length === 0 && (
            <tr><td colSpan={7} className="p-4 text-center text-[12px] text-gray-400">No transactions</td></tr>
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

function sourceText(status: string, usingLive: boolean) {
  if (status === "error") return "Live API unavailable. Showing demo billing data.";
  if (status === "loading") return "Refreshing live billing data...";
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
