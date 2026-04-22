import { PROVIDERS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";

export function AdminProviders() {
  return (
    <div className="p-5">
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Providers / Sources</h3>
          <StatusBadge status="manual_review" dot />
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Name", "Type", "Health", "Capacity %", "Fail Rate %", "Last Sync"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {PROVIDERS.map((p) => (
              <tr key={p.id} className={`hover:bg-gray-50 border-b border-gray-100 last:border-0 ${p.health === "degraded" ? "bg-amber-50/40" : ""}`}>
                <td className="px-3 py-2 font-mono text-[12px] text-gray-400">{p.id}</td>
                <td className="px-3 py-2 font-medium text-gray-900">{p.name}</td>
                <td className="px-3 py-2">
                  <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{p.type}</span>
                </td>
                <td className="px-3 py-2">
                  <StatusBadge status={p.health === "ok" ? "active" : p.health === "degraded" ? "pending" : "failed"} dot />
                </td>
                <td className="px-3 py-2">
                  <div className="flex items-center gap-2">
                    <div className="flex-1 h-1.5 bg-gray-100 rounded">
                      <div
                        className={`h-full rounded ${p.capacity > 70 ? "bg-green-500" : p.capacity > 40 ? "bg-amber-400" : "bg-red-500"}`}
                        style={{ width: `${p.capacity}%` }}
                      />
                    </div>
                    <span className="text-[12px] tabular-nums w-8 text-right">{p.capacity}%</span>
                  </div>
                </td>
                <td className="px-3 py-2 tabular-nums">
                  <span className={p.failRate > 1 ? "text-red-600 font-medium" : "text-gray-700"}>{p.failRate}%</span>
                </td>
                <td className="px-3 py-2 text-gray-400">{p.lastSync}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
