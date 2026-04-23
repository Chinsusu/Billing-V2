import { CreditCard, User, Server, Headphones, XCircle, Wallet } from "lucide-react";
import { KpiCard } from "@/components/ui/KpiCard";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { INVOICES, ACTIVITY_FEED, ActivityEvent } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

const ACTIVITY_ICONS: Record<ActivityEvent["icon"], React.ReactNode> = {
  payment: <CreditCard size={11} />,
  user:    <User size={11} />,
  server:  <Server size={11} />,
  ticket:  <Headphones size={11} />,
  error:   <XCircle size={11} />,
  wallet:  <Wallet size={11} />,
};

export function AdminOverview() {
  return (
    <div className="p-4 flex flex-col gap-4">
      {/* KPIs */}
      <div className="grid grid-cols-4 gap-4">
        <KpiCard label="MRR" value="$118.2k" delta={8.4} sub="vs last month" />
        <KpiCard label="Revenue · MTD" value="$29.2k" delta={12.1} sub="Apr 2026" />
        <KpiCard label="Active customers" value="2,847" delta={1.2} sub="net +34" />
        <KpiCard label="Active services" value="15,604" delta={-0.3} sub="proxies + VPS" />
      </div>

      {/* Main grid */}
      <div className="grid grid-cols-[1.7fr_1fr] gap-4">
        {/* Recent invoices */}
        <div className="bg-white border border-gray-200 rounded">
          <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Recent invoices</h3>
            <a href="#" className="text-[12px] text-[#D50C2D]">View all →</a>
          </div>
          <table className="w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["Invoice", "Customer", "Issued", "Due", "Amount", "Status"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {INVOICES.slice(0, 7).map((inv) => (
                <tr key={inv.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 p-4 text-[12px] text-[#D50C2D]">{inv.id}</td>
                  <td className="p-4 p-4 text-gray-700">{inv.customer}</td>
                  <td className="p-4 p-4 text-gray-400">{inv.issued}</td>
                  <td className="p-4 p-4 text-gray-400">{inv.due}</td>
                  <td className="p-4 p-4 text-right font-medium tabular-nums">{fmtMoney(inv.amount)}</td>
                  <td className="p-4 p-4"><StatusBadge status={inv.status} dot /></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Activity feed */}
        <div className="bg-white border border-gray-200 rounded">
          <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Activity feed</h3>
            <span className="text-[11px] text-gray-400">Live</span>
          </div>
          <div>
            {ACTIVITY_FEED.map((a, i) => (
              <div
                key={i}
                className="flex gap-4.5 p-4 p-4.5 border-b border-gray-100 last:border-0 items-start"
              >
                <div
                  className={`w-[22px] h-[22px] rounded-full grid place-items-center text-[11px] shrink-0
                    ${a.type === "ok" ? "bg-green-50 text-green-700"
                    : a.type === "warn" ? "bg-amber-50 text-amber-700"
                    : a.type === "danger" ? "bg-red-50 text-red-700"
                    : "bg-gray-100 text-gray-500"}`}
                >
                  {ACTIVITY_ICONS[a.icon]}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="text-[12px] text-gray-700 leading-snug">{a.text}</div>
                  <div className="text-[11px] text-gray-400 mt-0.5">{a.t}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Infrastructure health */}
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Infrastructure health</h3>
          <StatusBadge status="active" dot />
        </div>
        <div>
          {[
            { label: "Proxy network · uptime", value: "99.98%", bar: 0.9998 },
            { label: "VPS fleet · uptime", value: "99.94%", bar: 0.9994 },
            { label: "Payment gateway", value: "100%", bar: 1 },
            { label: "API · p95 latency", value: "142ms", bar: 0.82 },
            { label: "Support · first response", value: "8m avg", bar: 0.72 },
          ].map((r, i) => (
            <div key={i} className="p-4 p-4.5 border-b border-gray-100 last:border-0">
              <div className="flex items-center justify-between mb-1.5">
                <span className="text-[12px] text-gray-700">{r.label}</span>
                <span className="text-[12px] font-medium tabular-nums">{r.value}</span>
              </div>
              <div className="h-1 bg-gray-100 rounded">
                <div
                  className={`h-full rounded ${r.bar > 0.99 ? "bg-green-500" : r.bar > 0.8 ? "bg-blue-500" : "bg-amber-400"}`}
                  style={{ width: `${r.bar * 100}%` }}
                />
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
