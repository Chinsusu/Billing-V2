import { PROVISIONING_JOBS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";

export function AdminProvisioning() {
  const manualReview = PROVISIONING_JOBS.filter((j) => j.status === "manual_review");

  return (
    <div className="p-5 flex flex-col gap-4">
      {manualReview.length > 0 && (
        <div className="bg-amber-50 border border-amber-200 text-amber-700 text-[12px] px-4 py-2.5 rounded flex items-center gap-2">
          <span>⚠</span>
          <span>{manualReview.length} job(s) require manual review — provider state unknown. Do not re-trigger without verifying.</span>
        </div>
      )}

      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Provisioning queue</h3>
          <span className="text-[11px] text-gray-400">{PROVISIONING_JOBS.length} jobs</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Job ID", "Order", "Service", "Tenant", "Provider", "Status", "Attempt", "Age", "Error"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {PROVISIONING_JOBS.map((job) => (
              <tr key={job.id} className={`hover:bg-gray-50 border-b border-gray-100 last:border-0 ${job.status === "manual_review" ? "bg-amber-50/40" : ""}`}>
                <td className="px-3 py-2 font-mono text-[12px] text-gray-500">{job.id}</td>
                <td className="px-3 py-2 font-mono text-[12px] text-[#D50C2D]">{job.order}</td>
                <td className="px-3 py-2 text-gray-700">{job.service}</td>
                <td className="px-3 py-2 text-gray-500">{job.tenant}</td>
                <td className="px-3 py-2 text-gray-500">{job.provider}</td>
                <td className="px-3 py-2"><StatusBadge status={job.status} dot /></td>
                <td className="px-3 py-2 text-center tabular-nums">{job.attempt}</td>
                <td className="px-3 py-2 text-gray-400 tabular-nums">{job.age}</td>
                <td className="px-3 py-2 text-[11px] text-red-600 font-mono max-w-[200px] truncate">{job.error || "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
