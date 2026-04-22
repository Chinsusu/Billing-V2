import { RESELLER_CLIENTS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";

export function ResellerClients() {
  return (
    <div className="p-5">
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Clients</h3>
          <button className="h-7 px-3 text-[12px] font-medium bg-[#D50C2D] text-white rounded-[3px] border-0 hover:bg-[#B3082A] cursor-pointer">
            + Add client
          </button>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Name", "Email", "Wallet", "Services", "Orders", "Status", "Last login"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {RESELLER_CLIENTS.map((c) => (
              <tr key={c.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="px-3 py-2 font-mono text-[12px] text-[#D50C2D]">{c.id}</td>
                <td className="px-3 py-2 font-medium text-gray-900">{c.name}</td>
                <td className="px-3 py-2 text-gray-400 text-[12px]">{c.email}</td>
                <td className="px-3 py-2 tabular-nums">
                  <span className={c.wallet < 20 ? "text-red-600 font-medium" : ""}>{fmtMoney(c.wallet)}</span>
                </td>
                <td className="px-3 py-2 tabular-nums text-right">{c.services}</td>
                <td className="px-3 py-2 tabular-nums text-right">{c.orders}</td>
                <td className="px-3 py-2"><StatusBadge status={c.status} dot /></td>
                <td className="px-3 py-2 text-gray-400">{c.lastLogin}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
