import { TRANSACTIONS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";

export function AdminTransactions() {
  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Transactions / Ledger</h3>
          <span className="text-[11px] text-gray-400">Apr 2026</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Time", "Customer", "Method", "Type", "Amount", "Status"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {TRANSACTIONS.map((tx) => (
              <tr key={tx.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 text-[12px] text-gray-400">{tx.id}</td>
                <td className="p-4 p-4 text-gray-400 tabular-nums">{tx.time}</td>
                <td className="p-4 p-4 text-gray-700">{tx.customer}</td>
                <td className="p-4 p-4 text-gray-500">{tx.method}</td>
                <td className="p-4 p-4">
                  <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{tx.type}</span>
                </td>
                <td className={`p-4 p-4 text-right tabular-nums font-medium ${tx.amount < 0 ? "text-red-600" : ""}`}>
                  {fmtMoney(tx.amount)}
                </td>
                <td className="p-4 p-4"><StatusBadge status={tx.status} dot /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
