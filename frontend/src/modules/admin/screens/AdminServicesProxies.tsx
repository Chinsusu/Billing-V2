import { StatusBadge } from "@/components/ui/StatusBadge";
import { Settings, Edit2 } from "lucide-react";
import { PROXY_SERVICES } from "@/mocks/billingData";

export function AdminServicesProxies() {
  const cols = [
    { label: "ID", align: "left" },
    { label: "Host:Port", align: "left" },
    { label: "User/Pass", align: "left" },
    { label: "Port(http/Socks)", align: "left" },
    { label: "Status", align: "center" },
    { label: "Region", align: "left" },
    { label: "Member", align: "left" },
    { label: "Plan", align: "center" },
    { label: "Date", align: "left" },
    { label: "Expire", align: "center" },
    { label: "Auto Renew", align: "center" },
    { label: "Protection", align: "center" },
    { label: "Details", align: "center" },
    { label: "Action", align: "center" },
  ];

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded shadow-sm text-[12px]">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Proxies</h3>
          <span className="text-[11px] text-gray-400">{PROXY_SERVICES.length} services</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-50 text-gray-500">
                {cols.map((c) => (
                  <th key={c.label} className={`font-medium p-3 text-${c.align} border-b border-gray-200 text-[11px] tracking-wider`}>
                    {c.label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="bg-white text-gray-600">
              {PROXY_SERVICES.map((s, i) => {
                const now = new Date();
                const exp = new Date(now.getTime() + s.renewsIn * 24 * 3600 * 1000);
                const ord = new Date(exp.getTime() - 30 * 24 * 3600 * 1000);
                const pad = (n: number) => n.toString().padStart(2, '0');
                const f = (d: Date) => `${pad(d.getDate())}-${pad(d.getMonth() + 1)}-${d.getFullYear()} ${pad(d.getHours())}:${pad(d.getMinutes())}`;

                return (
                  <tr key={s.id} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="p-3 text-[#D50C2D] font-medium">{s.id}</td>
                    <td className="p-3">
                      <div className="bg-gray-100/80 rounded px-2 py-1 text-[11px] inline-flex items-center gap-1 border border-gray-200/50">
                        103.160.2.{77 + i}:N/A <Edit2 size={10} className="text-gray-400 cursor-pointer" />
                      </div>
                    </td>
                    <td className="p-3">
                      <div className="flex flex-col gap-1.5 align-start items-start">
                        <div className="bg-gray-100/80 rounded px-2 py-1 text-[11px] inline-flex items-center gap-1 border border-gray-200/50">
                          i5Syd8tg{i} <Edit2 size={10} className="text-gray-400 cursor-pointer" />
                        </div>
                        <div className="bg-gray-100/80 rounded px-2 py-1 text-[11px] inline-flex items-center gap-1 border border-gray-200/50">
                          yluheKv{i} <Edit2 size={10} className="text-gray-400 cursor-pointer" />
                        </div>
                      </div>
                    </td>
                    <td className="p-3 text-[11px]">
                      <div className="flex flex-col text-gray-500 gap-1.5 align-start items-start">
                        <div className="flex items-center gap-1">HTTPS: <span className="bg-gray-100 rounded px-1 text-gray-400 text-[10px]">N/A</span> <Edit2 size={10} className="cursor-pointer text-gray-400" /></div>
                        <div className="flex items-center gap-1">SOCKS5: <span className="bg-gray-100 rounded px-1 text-gray-500 text-[10px]">{41750 + i}</span> <Edit2 size={10} className="cursor-pointer text-gray-400" /></div>
                      </div>
                    </td>
                    <td className="p-3 text-center"><StatusBadge status={s.status} dot /></td>
                    <td className="p-3">{s.region}</td>
                    <td className="p-3 text-[11px]">
                      <div className="font-medium text-gray-800">{s.customer}</div>
                      <div className="text-gray-400">{s.tenant.toLowerCase().replace(/ /g, '')}@gmail.com</div>
                    </td>
                    <td className="p-3 text-center">
                      <div className="flex flex-col items-center gap-1.5">
                        <div className="bg-indigo-500 text-white rounded-full px-3 py-1 text-[10px] font-medium leading-none">proxy-{s.proxyType}</div>
                        <div className="bg-indigo-500 text-white rounded-full px-3 py-1 text-[10px] font-medium leading-none">{s.price.toLocaleString()}</div>
                      </div>
                    </td>
                    <td className="p-3 text-[10px] text-gray-400 leading-tight">
                      <div className="mb-0.5">{f(ord)}</div>
                      <div>{f(exp)}</div>
                    </td>
                    <td className="p-3 text-center">
                      <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-medium text-white ${s.renewsIn < 0 ? "bg-[#D50C2D]" : "bg-[#D50C2D]"}`}>
                        {s.renewsIn < 0 ? `quá hạn ${Math.abs(s.renewsIn)} ngày` : `còn ${s.renewsIn} ngày`}
                        <Edit2 size={10} />
                      </span>
                    </td>
                    <td className="p-3 text-center">
                      <div className="w-8 h-4 bg-gray-200 rounded-full inline-block relative cursor-pointer border border-gray-300">
                        <div className="w-3 h-3 bg-white shadow rounded-full absolute left-0.5 top-0.5"></div>
                      </div>
                    </td>
                    <td className="p-3 text-center">
                      <div className="w-8 h-4 bg-emerald-500 rounded-full inline-block relative cursor-pointer border border-emerald-600">
                        <div className="w-3 h-3 bg-white shadow rounded-full absolute right-0.5 top-0.5"></div>
                      </div>
                    </td>
                    <td className="p-3 text-center">
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
    </div>
  );
}
