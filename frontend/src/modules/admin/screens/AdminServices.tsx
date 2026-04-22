import { SERVICES } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";

export function AdminServices() {
  return (
    <div className="p-5">
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Services</h3>
          <span className="text-[11px] text-gray-400">{SERVICES.length} shown</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Type", "Label", "Customer", "Region", "Bandwidth", "Price/mo", "Status", "Renews in"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {SERVICES.map((s) => (
              <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="px-3 py-2 font-mono text-[12px] text-[#D50C2D]">{s.id}</td>
                <td className="px-3 py-2">
                  <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{s.type}</span>
                </td>
                <td className="px-3 py-2 text-gray-800">{s.label}</td>
                <td className="px-3 py-2 text-gray-500">{s.customer}</td>
                <td className="px-3 py-2 font-mono text-[12px] text-gray-400">{s.region}</td>
                <td className="px-3 py-2 text-gray-400">{s.bandwidth}</td>
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
