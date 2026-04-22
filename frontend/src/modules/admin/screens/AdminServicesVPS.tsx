import { StatusBadge } from "@/components/ui/StatusBadge";
import { Settings } from "lucide-react";
import { VPS_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function AdminServicesVPS() {
  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">VPS</h3>
          <span className="text-[11px] text-gray-400">{VPS_SERVICES.length} services</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "OS", "Label", "Customer", "Tenant", "Region", "Spec", "IP", "Provider", "Price/mo", "Status", "Date", "Expire", "Action"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {VPS_SERVICES.map((s) => (
              <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 text-[11px] text-[#D50C2D]">{s.id}</td>
                <td className="p-4 p-4">
                  <span className={`text-[10px] font-medium px-1.5 py-0.5 rounded ${s.os === "linux" ? "bg-orange-50 text-orange-700" : "bg-blue-50 text-blue-700"}`}>
                    {s.os === "linux" ? "Linux" : "Windows"}
                  </span>
                </td>
                <td className="p-4 p-4 text-gray-800 max-w-[160px] truncate">{s.label}</td>
                <td className="p-4 p-4 text-gray-500">{s.customer}</td>
                <td className="p-4 p-4 text-gray-400 text-[11px]">{s.tenant}</td>
                <td className="p-4 p-4 text-[11px] text-gray-400">{s.region}</td>
                <td className="p-4 p-4 text-gray-500 text-[11px] whitespace-nowrap">
                  {s.cpu}C / {s.ram}G / {s.disk}G
                </td>
                <td className="p-4 p-4 text-[11px] text-gray-400">{s.ip}</td>
                <td className="p-4 p-4 text-gray-400 text-[11px]">{s.provider}</td>
                <td className="p-4 p-4 tabular-nums text-right font-medium">{fmtMoney(s.price)}</td>
                <td className="p-4 p-4 w-[110px]"><StatusBadge status={s.status} dot /></td>
                <td className="p-4 p-4 tabular-nums text-[11px] text-gray-500 whitespace-nowrap leading-relaxed">
                  {(() => {
                    const now = new Date();
                    const exp = new Date(now.getTime() + s.renewsIn * 24 * 3600 * 1000);
                    const ord = new Date(exp.getTime() - 30 * 24 * 3600 * 1000);
                    const pad = (n: number) => n.toString().padStart(2, '0');
                    const f = (d: Date) => `${pad(d.getDate())}-${pad(d.getMonth() + 1)}-${d.getFullYear()} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
                    return <><div className="mb-0.5 text-gray-800">{f(ord)}</div><div>{f(exp)}</div></>;
                  })()}
                </td>
                <td className="p-4 p-4 tabular-nums whitespace-nowrap">
                  <span className={`inline-block px-2.5 py-0.5 rounded-full text-[10px] font-medium 
                    ${s.renewsIn < 0 ? "bg-[#D50C2D] text-white" : "bg-[#D50C2D] text-white"}`}>
                    {s.renewsIn < 0 ? `quá hạn ${Math.abs(s.renewsIn)} ngày` : `còn ${s.renewsIn} ngày`}
                  </span>
                </td>
                <td className="p-4 p-4 w-[60px] text-center">
                  <button className="text-gray-400 hover:text-gray-600 p-1 rounded transition-colors cursor-pointer bg-transparent border-0 inline-flex items-center justify-center">
                    <Settings size={14} />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
