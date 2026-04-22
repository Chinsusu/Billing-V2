import { AUDIT_LOGS, AuditLog } from "@/mocks/billingData";

const LEVEL_STYLE: Record<AuditLog["level"], string> = {
  info:  "bg-blue-50 text-blue-700",
  warn:  "bg-amber-50 text-amber-700",
  error: "bg-red-50 text-red-700",
};

const ACTOR_STYLE: Record<AuditLog["actor"], string> = {
  system:   "bg-gray-100 text-gray-500",
  admin:    "bg-purple-50 text-purple-700",
  reseller: "bg-indigo-50 text-indigo-700",
  client:   "bg-teal-50 text-teal-700",
};

export function AdminLogs() {
  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Audit logs</h3>
          <span className="text-[11px] text-gray-400">{AUDIT_LOGS.length} entries (mock)</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Time", "Level", "Actor", "Action", "Target", "Detail", "req_id"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {AUDIT_LOGS.map((log) => (
              <tr key={log.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 text-gray-400 text-[11px] tabular-nums whitespace-nowrap">{log.ts}</td>
                <td className="p-4 p-4">
                  <span className={`text-[10px] font-medium uppercase px-1.5 py-0.5 rounded ${LEVEL_STYLE[log.level]}`}>
                    {log.level}
                  </span>
                </td>
                <td className="p-4 p-4 whitespace-nowrap">
                  <span className={`text-[10px] font-medium px-1.5 py-0.5 rounded mr-1 ${ACTOR_STYLE[log.actor]}`}>
                    {log.actor}
                  </span>
                  <span className="text-gray-600">{log.actorName}</span>
                </td>
                <td className="p-4 p-4 text-[11px] text-gray-700">{log.action}</td>
                <td className="p-4 p-4 text-[11px] text-[#D50C2D]">{log.target}</td>
                <td className="p-4 p-4 text-gray-500 max-w-[260px] truncate">{log.detail}</td>
                <td className="p-4 p-4 text-[11px] text-gray-300">{log.requestId}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
