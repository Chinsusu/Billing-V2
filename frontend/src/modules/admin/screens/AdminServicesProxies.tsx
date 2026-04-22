import { StatusBadge } from "@/components/ui/StatusBadge";
import { Settings } from "lucide-react";
import { PROXY_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

const TYPE_LABEL: Record<string, string> = {
  residential: "Residential",
  datacenter: "Datacenter",
  mobile: "Mobile",
  isp: "ISP",
};

export function AdminServicesProxies() {
  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Proxies</h3>
          <span className="text-[11px] text-gray-400">{PROXY_SERVICES.length} services</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {
              [
                { label: "ID", align: "left" },
                { label: "Type", align: "left" },
                { label: "Label", align: "left" },
                { label: "Customer", align: "left" },
                { label: "Tenant", align: "left" },
                { label: "Region", align: "left" },
                { label: "IPs", align: "right" },
                { label: "Protocol", align: "left" },
                { label: "Usage", align: "left" },
                { label: "Price/mo", align: "right" },
                { label: "Status", align: "center" },
                { label: "Date", align: "left" },
                { label: "Expire", align: "center" },
                { label: "Action", align: "center" }
              ].map((h) => (
                <th key={h.label} className={`text-${h.align} text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200`}>
                  {h.label}
                </th>
              ))
            }
            </tr>
          </thead>
          <tbody>
            {PROXY_SERVICES.map((s) => {
              const pct = Math.round((s.usedGB / s.totalGB) * 100);
              return (
                <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[11px] text-[#D50C2D]">{s.id}</td>
                  <td className="p-4">
                    <span className="text-[10px] font-medium px-1.5 py-0.5 rounded bg-indigo-50 text-indigo-700">
                      {TYPE_LABEL[s.proxyType]}
                    </span>
                  </td>
                  <td className="p-4 text-gray-800 max-w-[180px] truncate">{s.label}</td>
                  <td className="p-4 text-gray-500">{s.customer}</td>
                  <td className="p-4 text-gray-400 text-[11px]">{s.tenant}</td>
                  <td className="p-4 text-[11px] text-gray-400">{s.region}</td>
                  <td className="p-4 text-gray-500 tabular-nums text-right">{s.ipCount > 0 ? s.ipCount : "—"}</td>
                  <td className="p-4">
                    <span className="text-[10px] px-1 py-px bg-gray-100 text-gray-500 rounded">{s.protocol}</span>
                  </td>
                  <td className="p-4">
                    <div className="flex items-center gap-4 min-w-[100px]">
                      <div className="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden">
                        <div
                          className={`h-full rounded-full ${pct >= 90 ? "bg-red-500" : pct >= 70 ? "bg-amber-400" : "bg-green-500"}`}
                          style={{ width: `${pct}%` }}
                        />
                      </div>
                      <span className="text-[11px] text-gray-400 tabular-nums w-8 text-right">{pct}%</span>
                    </div>
                  </td>
                  <td className="p-4 tabular-nums text-right font-medium">{fmtMoney(s.price)}</td>
                  <td className="p-4 w-[110px] text-center"><StatusBadge status={s.status} dot /></td>
                  <td className="p-4 tabular-nums text-[11px] text-gray-500 whitespace-nowrap leading-relaxed">
                    {(() => {
                      const now = new Date();
                      const exp = new Date(now.getTime() + s.renewsIn * 24 * 3600 * 1000);
                      const ord = new Date(exp.getTime() - 30 * 24 * 3600 * 1000);
                      const pad = (n: number) => n.toString().padStart(2, '0');
                      const f = (d: Date) => `${pad(d.getDate())}-${pad(d.getMonth() + 1)}-${d.getFullYear()} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
                      return <><div className="mb-0.5 text-gray-800">{f(ord)}</div><div>{f(exp)}</div></>;
                    })()}
                  </td>
                  <td className="p-4 tabular-nums whitespace-nowrap text-center">
                    <span className={`inline-block px-2.5 py-0.5 rounded-full text-[10px] font-medium 
                      ${s.renewsIn < 0 ? "bg-[#D50C2D] text-white" : "bg-[#D50C2D] text-white"}`}>
                      {s.renewsIn < 0 ? `quá hạn ${Math.abs(s.renewsIn)} ngày` : `còn ${s.renewsIn} ngày`}
                    </span>
                  </td>
                  <td className="p-4 w-[60px] text-center">
                    <button className="text-gray-400 hover:text-gray-600 p-1 rounded transition-colors cursor-pointer bg-transparent border-0 inline-flex items-center justify-center">
                      <Settings size={14} />
                    </button>
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
