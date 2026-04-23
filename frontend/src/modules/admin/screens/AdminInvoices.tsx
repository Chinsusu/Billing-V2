import { INVOICES } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";

export function AdminInvoices() {
  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Invoices</h3>
          <span className="text-[11px] text-gray-400">{INVOICES.length} records</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Invoice", "Customer", "Issued", "Due", "Amount", "Status"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {INVOICES.map((inv) => (
              <tr key={inv.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 text-[12px] text-[#D50C2D]">{inv.id}</td>
                <td className="p-4 p-4 text-gray-700">{inv.customer}</td>
                <td className="p-4 p-4 text-gray-400">{inv.issued}</td>
                <td className="p-4 p-4 text-gray-400">{inv.due}</td>
                <td className="p-4 p-4 text-right font-medium tabular-nums">{fmtMoney(inv.amount)}</td>
                <td className="p-4 p-4"><StatusBadge status={inv.status} dot /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
