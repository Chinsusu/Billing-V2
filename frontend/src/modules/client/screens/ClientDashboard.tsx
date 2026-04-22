import { StatusBadge } from "@/components/ui/StatusBadge";
import { CLIENT_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function ClientDashboard() {
  const suspended = CLIENT_SERVICES.filter((s) => s.status === "suspended");

  return (
    <div className="p-5 flex flex-col gap-4">
      {suspended.map((s) => (
        <div key={s.id} className="bg-amber-50 border border-amber-200 text-amber-700 text-[12px] px-4 py-2.5 rounded flex items-center gap-2">
          <span>⚠</span>
          <span><strong>{s.label}</strong> is suspended. {s.note} Renew to restore access.</span>
          <button className="ml-auto h-6 px-3 text-[11px] font-medium bg-amber-600 text-white rounded-sm border-0 cursor-pointer hover:bg-amber-700">
            Renew now
          </button>
        </div>
      ))}

      {/* Wallet */}
      <div className="bg-white border border-gray-200 rounded p-4 flex items-center justify-between">
        <div>
          <div className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Wallet balance</div>
          <div className="text-2xl font-bold tabular-nums text-gray-900">$128.40</div>
        </div>
        <button className="h-8 px-4 text-[13px] font-medium bg-[#D50C2D] text-white rounded-[3px] border-0 hover:bg-[#B3082A] cursor-pointer">
          + Top up
        </button>
      </div>

      {/* Services */}
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">My services</h3>
          <a href="#" className="text-[12px] text-[#D50C2D]">View all →</a>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Label", "Region", "Bandwidth", "Expires", "Status"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {CLIENT_SERVICES.map((s) => (
              <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="px-3 py-2 font-medium text-gray-900">{s.label}</td>
                <td className="px-3 py-2 font-mono text-[12px] text-gray-400">{s.region}</td>
                <td className="px-3 py-2 text-gray-500">{s.bandwidth}</td>
                <td className="px-3 py-2 text-gray-400">{s.expiry}</td>
                <td className="px-3 py-2"><StatusBadge status={s.status} dot /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
