import { TOPUP_REQUESTS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";

export function AdminTopups() {
  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Top-up verification queue</h3>
          <span className="text-[11px] text-gray-400">
            {TOPUP_REQUESTS.filter((t) => t.status === "pending_verification").length} pending
          </span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Tenant / Actor", "Amount", "Method", "Ref", "Created", "Proof", "Status", "Actions"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {TOPUP_REQUESTS.map((req) => (
              <tr key={req.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 text-[12px] text-[#D50C2D]">{req.id}</td>
                <td className="p-4 p-4">
                  <div className="text-gray-900">{req.tenant}</div>
                  <div className="text-[11px] text-gray-400">{req.actor}</div>
                </td>
                <td className="p-4 p-4 tabular-nums font-medium">{fmtMoney(req.amount)}</td>
                <td className="p-4 p-4 text-gray-500">{req.method}</td>
                <td className="p-4 p-4 text-[12px] text-gray-400">{req.ref}</td>
                <td className="p-4 p-4 text-gray-400">{req.created}</td>
                <td className="p-4 p-4 text-center">{req.proof ? "✓" : <span className="text-red-500">✕</span>}</td>
                <td className="p-4 p-4"><StatusBadge status={req.status} dot /></td>
                <td className="p-4 p-4">
                  {req.status === "pending_verification" && (
                    <div className="flex gap-1.5">
                      <button className="h-6 p-4 text-[11px] font-medium bg-green-600 text-white rounded-sm hover:bg-green-700 cursor-pointer border-0">
                        Approve
                      </button>
                      <button className="h-6 p-4 text-[11px] font-medium border border-red-200 text-red-600 rounded-sm hover:bg-red-50 cursor-pointer bg-white">
                        Reject
                      </button>
                    </div>
                  )}
                  {req.reason && <span className="text-[11px] text-red-500">{req.reason}</span>}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
