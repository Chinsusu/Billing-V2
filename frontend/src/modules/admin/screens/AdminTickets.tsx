import { TICKETS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { ticketPriorityLabel } from "@/lib/api/displayLabels";

const PRIORITY_CLASS: Record<string, string> = {
  high: "text-red-600 font-medium",
  medium: "text-amber-600",
  low: "text-gray-400",
};

export function AdminTickets() {
  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Support tickets</h3>
          <span className="text-[11px] text-gray-400">{TICKETS.filter((t) => t.status === "open").length} open</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Subject", "Customer", "Priority", "Status", "Updated", "Assignee"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {TICKETS.map((t) => (
              <tr key={t.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 text-[12px] text-[#D50C2D]">{t.id}</td>
                <td className="p-4 p-4 text-gray-800 max-w-[300px] truncate">{t.subject}</td>
                <td className="p-4 p-4 text-gray-500">{t.customer}</td>
                <td className={`p-4 p-4 text-[12px] ${PRIORITY_CLASS[t.priority]}`}>{ticketPriorityLabel(t.priority)}</td>
                <td className="p-4 p-4"><StatusBadge status={t.status} dot /></td>
                <td className="p-4 p-4 text-gray-400">{t.updated}</td>
                <td className="p-4 p-4 text-gray-500">{t.assignee}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
