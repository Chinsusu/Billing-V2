import { StatusBadge } from "@/components/ui/StatusBadge";
import { VPS_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function AdminServicesVPS() {
  return (
    <div className="p-5">
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">VPS</h3>
          <span className="text-[11px] text-gray-400">{VPS_SERVICES.length} services</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "OS", "Label", "Customer", "Tenant", "Region", "Spec", "IP", "Provider", "Price/mo", "Status", "Renews"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {VPS_SERVICES.map((s) => (
              <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="px-3 py-2 font-mono text-[11px] text-[#D50C2D]">{s.id}</td>
                <td className="px-3 py-2">
                  <span className={`text-[10px] font-medium px-1.5 py-0.5 rounded ${s.os === "linux" ? "bg-orange-50 text-orange-700" : "bg-blue-50 text-blue-700"}`}>
                    {s.os === "linux" ? "Linux" : "Windows"}
                  </span>
                </td>
                <td className="px-3 py-2 text-gray-800 max-w-[160px] truncate">{s.label}</td>
                <td className="px-3 py-2 text-gray-500">{s.customer}</td>
                <td className="px-3 py-2 text-gray-400 text-[11px]">{s.tenant}</td>
                <td className="px-3 py-2 font-mono text-[11px] text-gray-400">{s.region}</td>
                <td className="px-3 py-2 text-gray-500 text-[11px] whitespace-nowrap">
                  {s.cpu}C / {s.ram}G / {s.disk}G
                </td>
                <td className="px-3 py-2 font-mono text-[11px] text-gray-400">{s.ip}</td>
                <td className="px-3 py-2 text-gray-400 text-[11px]">{s.provider}</td>
                <td className="px-3 py-2 tabular-nums text-right font-medium">{fmtMoney(s.price)}</td>
                <td className="px-3 py-2"><StatusBadge status={s.status} dot /></td>
                <td className="px-3 py-2 tabular-nums">
                  <span className={s.renewsIn < 0 ? "text-red-600 font-medium" : s.renewsIn <= 7 ? "text-amber-600" : "text-gray-500"}>
                    {s.renewsIn < 0 ? `${Math.abs(s.renewsIn)}d overdue` : `${s.renewsIn}d`}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
