import { KpiCard } from "@/components/ui/KpiCard";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { RESELLER_CLIENTS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function ResellerDashboard() {
  return (
    <div className="p-5 flex flex-col gap-4">
      {/* Wallet hero */}
      <div className="bg-white border border-gray-200 rounded p-5 flex items-center justify-between">
        <div>
          <div className="text-[11px] font-medium uppercase tracking-wide text-gray-400 mb-1">Reseller wallet · T-0042</div>
          <div className="text-3xl font-bold tabular-nums text-gray-900">$4,820.50</div>
          <div className="text-[12px] text-gray-400 mt-1">Available balance · ProxyVN</div>
        </div>
        <div className="flex gap-2">
          <button className="h-8 px-4 text-[13px] font-medium border border-gray-300 rounded-[3px] bg-white hover:bg-gray-50 cursor-pointer">
            Ledger history
          </button>
          <button className="h-8 px-4 text-[13px] font-medium bg-[#D50C2D] text-white rounded-[3px] hover:bg-[#B3082A] cursor-pointer border-0">
            + Top up
          </button>
        </div>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-4 gap-3">
        <KpiCard label="Total clients" value="312" delta={2.8} sub="this month" />
        <KpiCard label="Active services" value="1,840" delta={1.4} sub="across clients" />
        <KpiCard label="Revenue · MTD" value="$12.4k" delta={6.2} sub="Apr 2026" />
        <KpiCard label="Pending top-ups" value="2" sub="awaiting admin" />
      </div>

      {/* Recent clients */}
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Recent clients</h3>
          <a href="#" className="text-[12px] text-[#D50C2D]">View all →</a>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Name", "Wallet", "Services", "Orders", "Status", "Last login"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {RESELLER_CLIENTS.slice(0, 5).map((c) => (
              <tr key={c.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="px-3 py-2 font-mono text-[12px] text-[#D50C2D]">{c.id}</td>
                <td className="px-3 py-2 font-medium text-gray-900">{c.name}</td>
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
