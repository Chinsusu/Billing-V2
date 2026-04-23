import { TENANTS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";

export function AdminTenants() {
  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Accounts</h3>
          <button className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-white hover:bg-gray-50 text-gray-700 border border-gray-300 rounded-md cursor-pointer transition-colors shadow-sm">
            + New account
          </button>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {[
                { label: "ID", align: "left" },
                { label: "Name", align: "left" },
                { label: "Type", align: "left" },
                { label: "Domain", align: "left" },
                { label: "Clients", align: "right" },
                { label: "Services", align: "right" },
                { label: "Wallet", align: "right" },
                { label: "Status", align: "left" },
                { label: "Since", align: "left" }
              ].map((h) => (
                <th key={h.label} className={`${h.align === 'right' ? 'text-right' : 'text-left'} text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200`}>
                  {h.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {TENANTS.map((t) => (
              <tr key={t.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 text-[12px] text-[#D50C2D]">{t.id}</td>
                <td className="p-4 p-4 font-medium text-gray-900">{t.name}</td>
                <td className="p-4 p-4">
                  <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{t.type}</span>
                </td>
                <td className="p-4 p-4 text-[12px] text-gray-500">{t.domain}</td>
                <td className="p-4 p-4 tabular-nums text-right">{t.clients.toLocaleString()}</td>
                <td className="p-4 p-4 tabular-nums text-right">{t.services.toLocaleString()}</td>
                <td className="p-4 p-4 tabular-nums text-right">
                  {t.type === "admin" ? "—" : (
                    <span className={t.walletLow ? "text-red-600 font-medium" : ""}>{fmtMoney(t.wallet)}</span>
                  )}
                </td>
                <td className="p-4 p-4"><StatusBadge status={t.status} dot /></td>
                <td className="p-4 p-4 text-gray-400">{t.since ?? "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
