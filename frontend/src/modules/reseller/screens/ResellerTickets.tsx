import { StatusBadge } from "@/components/ui/StatusBadge";
import { TICKETS } from "@/mocks/billingData";

export function ResellerTickets() {
  const rows = TICKETS.slice(0, 6);
  const open = rows.filter((ticket) => ticket.status !== "closed").length;

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <SummaryTile label="Tickets" value={String(rows.length)} />
        <SummaryTile label="Open" value={String(open)} tone={open > 0 ? "warn" : "neutral"} />
        <SummaryTile label="Data source" value="Demo adapter" />
      </div>
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Support tickets</h3>
          <span className="text-[11px] text-gray-400">Backend route pending</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[720px]">
            <thead>
              <tr className="bg-gray-50">
                {["Ticket", "Subject", "Client", "Priority", "Status", "Updated"].map((heading) => (
                  <th key={heading} className="text-left text-[11px] font-medium uppercase text-gray-400 p-4 border-b border-gray-200">
                    {heading}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((ticket) => (
                <tr key={ticket.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D] font-medium">{ticket.id}</td>
                  <td className="p-4 font-medium text-gray-900 max-w-[320px] truncate">{ticket.subject}</td>
                  <td className="p-4 text-gray-500">{ticket.customer}</td>
                  <td className="p-4 text-gray-500">{ticket.priority}</td>
                  <td className="p-4"><StatusBadge status={ticket.status} dot /></td>
                  <td className="p-4 text-gray-400">{ticket.updated}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function SummaryTile({ label, value, tone = "neutral" }: { label: string; value: string; tone?: "neutral" | "warn" }) {
  return (
    <div className={`bg-white border rounded p-4 ${tone === "warn" ? "border-amber-200" : "border-gray-200"}`}>
      <div className="text-[11px] text-gray-400 uppercase mb-1">{label}</div>
      <div className={`text-lg font-medium tabular-nums ${tone === "warn" ? "text-amber-700" : "text-gray-900"}`}>{value}</div>
    </div>
  );
}
