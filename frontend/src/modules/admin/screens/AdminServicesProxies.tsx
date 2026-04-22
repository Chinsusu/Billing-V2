import { StatusBadge } from "@/components/ui/StatusBadge";
import { PROXY_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

const TYPE_LABEL: Record<string, string> = {
  residential: "Residential",
  datacenter:  "Datacenter",
  mobile:      "Mobile",
  isp:         "ISP",
};

export function AdminServicesProxies() {
  return (
    <div className="p-5">
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Proxies</h3>
          <span className="text-[11px] text-gray-400">{PROXY_SERVICES.length} services</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Type", "Label", "Customer", "Tenant", "Region", "IPs", "Protocol", "Usage", "Price/mo", "Status", "Renews"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {PROXY_SERVICES.map((s) => {
              const pct = Math.round((s.usedGB / s.totalGB) * 100);
              return (
                <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="px-3 py-2 font-mono text-[11px] text-[#D50C2D]">{s.id}</td>
                  <td className="px-3 py-2">
                    <span className="text-[10px] font-medium px-1.5 py-0.5 rounded bg-indigo-50 text-indigo-700">
                      {TYPE_LABEL[s.proxyType]}
                    </span>
                  </td>
                  <td className="px-3 py-2 text-gray-800 max-w-[180px] truncate">{s.label}</td>
                  <td className="px-3 py-2 text-gray-500">{s.customer}</td>
                  <td className="px-3 py-2 text-gray-400 text-[11px]">{s.tenant}</td>
                  <td className="px-3 py-2 font-mono text-[11px] text-gray-400">{s.region}</td>
                  <td className="px-3 py-2 text-gray-500 tabular-nums">{s.ipCount > 0 ? s.ipCount : "—"}</td>
                  <td className="px-3 py-2">
                    <span className="text-[10px] font-mono px-1 py-px bg-gray-100 text-gray-500 rounded">{s.protocol}</span>
                  </td>
                  <td className="px-3 py-2">
                    <div className="flex items-center gap-2 min-w-[100px]">
                      <div className="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden">
                        <div
                          className={`h-full rounded-full ${pct >= 90 ? "bg-red-500" : pct >= 70 ? "bg-amber-400" : "bg-green-500"}`}
                          style={{ width: `${pct}%` }}
                        />
                      </div>
                      <span className="text-[11px] text-gray-400 tabular-nums w-8 text-right">{pct}%</span>
                    </div>
                  </td>
                  <td className="px-3 py-2 tabular-nums text-right font-medium">{fmtMoney(s.price)}</td>
                  <td className="px-3 py-2"><StatusBadge status={s.status} dot /></td>
                  <td className="px-3 py-2 tabular-nums">
                    <span className={s.renewsIn < 0 ? "text-red-600 font-medium" : s.renewsIn <= 7 ? "text-amber-600" : "text-gray-500"}>
                      {s.renewsIn < 0 ? `${Math.abs(s.renewsIn)}d overdue` : `${s.renewsIn}d`}
                    </span>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
