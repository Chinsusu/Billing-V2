"use client";

import { StatusBadge } from "@/components/ui/StatusBadge";
import { KpiCard } from "@/components/ui/KpiCard";
import { billingApi } from "@/lib/api/billing";
import { moneyMinor, recordLabel } from "@/lib/api/format";
import type { PaymentTransaction } from "@/lib/api/types";
import { useApiResource } from "@/lib/api/useApiResource";
import { mapAdminAuditLogView } from "@/lib/api/viewModels";
import { AUDIT_LOGS, INVOICES, TRANSACTIONS } from "@/mocks/billingData";
import { fmtMoney, fmtMoneyShort } from "@/mocks/sampleData";

interface ReconciliationRow {
  id: string;
  provider: string;
  invoice: string;
  ledger: string;
  amount: string;
  status: string;
}

interface AuditRow {
  id: string;
  time: string;
  actor: string;
  action: string;
  target: string;
  request: string;
}

function liveMoneyTotal(transactions: PaymentTransaction[]) {
  if (transactions.length === 0) return "$0.00";
  const currency = transactions[0]?.currency ?? "USD";
  const total = transactions
    .filter((transaction) => transaction.currency === currency)
    .reduce((sum, transaction) => sum + transaction.amount_minor, 0);
  const suffix = transactions.some((transaction) => transaction.currency !== currency) ? " +" : "";
  return `${moneyMinor(total, currency)}${suffix}`;
}

function sectionSource(status: string, liveLabel: string, fallbackLabel: string) {
  if (status === "error") return `${liveLabel} unavailable. Showing ${fallbackLabel}.`;
  if (status === "loading") return `Refreshing ${liveLabel}...`;
  return status === "success" ? liveLabel : fallbackLabel;
}

export function AdminReports() {
  const transactions = useApiResource(
    () => billingApi.listAdminTransactions(),
    "admin-reports:transactions",
  );
  const reconciliation = useApiResource(
    () => billingApi.listAdminReconciliation(),
    "admin-reports:reconciliation",
  );
  const invoices = useApiResource(
    () => billingApi.listAdminInvoices(),
    "admin-reports:invoices",
  );
  const auditLogs = useApiResource(
    () => billingApi.listAdminAuditLogs(),
    "admin-reports:audit",
  );

  const usingLiveTransactions = transactions.status === "success";
  const usingLiveReconciliation = reconciliation.status === "success";
  const usingLiveInvoices = invoices.status === "success";
  const usingLiveAudit = auditLogs.status === "success";

  const transactionRows: ReconciliationRow[] = usingLiveReconciliation
    ? (reconciliation.data ?? []).slice(0, 8).map((item) => ({
        id: recordLabel(item.transaction.display_id, "TX-"),
        provider: item.provider ?? "wallet",
        invoice: item.invoice ? recordLabel(item.invoice.display_id, "INV-") : "-",
        ledger: item.ledger ? recordLabel(item.ledger.display_id, "LED-") : "-",
        amount: moneyMinor(item.transaction.amount_minor, item.transaction.currency),
        status: item.transaction.status,
      }))
    : TRANSACTIONS.slice(0, 8).map((transaction) => ({
        id: transaction.id,
        provider: transaction.method,
        invoice: "-",
        ledger: "-",
        amount: fmtMoney(transaction.amount),
        status: transaction.status,
      }));

  const auditRows: AuditRow[] = usingLiveAudit
    ? (auditLogs.data ?? []).slice(0, 8).map((log) => {
        const row = mapAdminAuditLogView(log);
        return {
          id: row.id,
          time: row.ts,
          actor: `${row.actor} ${row.actorName}`,
          action: row.action,
          target: row.target,
          request: row.requestId,
        };
      })
    : AUDIT_LOGS.slice(0, 8).map((log) => ({
        id: log.id,
        time: log.ts,
        actor: `${log.actor} ${log.actorName}`,
        action: log.action,
        target: log.target,
        request: log.requestId,
      }));

  const liveTransactions = transactions.data ?? [];
  const liveInvoices = invoices.data ?? [];
  const invoiceTotal = usingLiveInvoices
    ? moneyMinor(liveInvoices.reduce((sum, invoice) => sum + invoice.total_minor, 0), liveInvoices[0]?.currency ?? "USD")
    : fmtMoneyShort(INVOICES.reduce((sum, invoice) => sum + invoice.amount, 0));
  const transactionTotal = usingLiveTransactions
    ? liveMoneyTotal(liveTransactions)
    : fmtMoneyShort(TRANSACTIONS.reduce((sum, transaction) => sum + transaction.amount, 0));

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
        <KpiCard
          label="Transaction volume"
          value={transactionTotal}
          sub={usingLiveTransactions ? `${liveTransactions.length} live records` : "demo fallback"}
        />
        <KpiCard
          label="Invoice value"
          value={invoiceTotal}
          sub={usingLiveInvoices ? `${liveInvoices.length} live invoices` : "demo fallback"}
        />
        <KpiCard
          label="Reconciled rows"
          value={(usingLiveReconciliation ? (reconciliation.data?.length ?? 0) : TRANSACTIONS.length).toLocaleString()}
          sub={usingLiveReconciliation ? "live reconciliation" : "demo fallback"}
        />
        <KpiCard
          label="Audit events"
          value={(usingLiveAudit ? (auditLogs.data?.length ?? 0) : AUDIT_LOGS.length).toLocaleString()}
          sub={usingLiveAudit ? "live audit log" : "demo fallback"}
        />
      </div>

      <div className="grid grid-cols-1 gap-4 xl:grid-cols-[1.3fr_1fr]">
        <div className="bg-white border border-gray-200 rounded">
          <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Payment reconciliation</h3>
            <span className="text-[11px] text-gray-400">
              {sectionSource(reconciliation.status, "Live reconciliation", "demo transaction data")}
            </span>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-[760px] w-full text-[13px] border-collapse">
              <thead>
                <tr className="bg-gray-50">
                  {["Transaction", "Provider", "Invoice", "Ledger", "Amount", "Status"].map((h) => (
                    <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {transactionRows.map((row) => (
                  <tr key={row.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                    <td className="p-4 text-[12px] text-[#D50C2D]">{row.id}</td>
                    <td className="p-4 text-gray-500">{row.provider}</td>
                    <td className="p-4 text-gray-500">{row.invoice}</td>
                    <td className="p-4 text-gray-400">{row.ledger}</td>
                    <td className="p-4 text-right tabular-nums font-medium">{row.amount}</td>
                    <td className="p-4"><StatusBadge status={row.status} dot /></td>
                  </tr>
                ))}
                {usingLiveReconciliation && transactionRows.length === 0 && (
                  <tr><td colSpan={6} className="p-4 text-center text-[12px] text-gray-400">No reconciliation rows</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>

        <div className="bg-white border border-gray-200 rounded">
          <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-4">
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Audit activity</h3>
            <span className="text-[11px] text-gray-400">
              {sectionSource(auditLogs.status, "Live audit logs", "demo audit data")}
            </span>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-[680px] w-full text-[13px] border-collapse">
              <thead>
                <tr className="bg-gray-50">
                  {["ID", "Time", "Actor", "Action", "Target", "Request"].map((h) => (
                    <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {auditRows.map((row) => (
                  <tr key={row.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                    <td className="p-4 text-[12px] text-[#D50C2D]">{row.id}</td>
                    <td className="p-4 text-gray-400">{row.time}</td>
                    <td className="p-4 text-gray-600">{row.actor}</td>
                    <td className="p-4 text-[11px] text-gray-700">{row.action}</td>
                    <td className="p-4 text-[11px] text-gray-500">{row.target}</td>
                    <td className="p-4 text-[11px] text-gray-300">{row.request}</td>
                  </tr>
                ))}
                {usingLiveAudit && auditRows.length === 0 && (
                  <tr><td colSpan={6} className="p-4 text-center text-[12px] text-gray-400">No audit events</td></tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
