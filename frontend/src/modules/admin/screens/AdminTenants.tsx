import { TENANTS } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { fmtMoney } from "@/mocks/sampleData";

export function AdminTenants() {
  return (
    <div className="p-5">
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Tenants</h3>
          <button className="h-7 px-3 text-[12px] font-medium border border-gray-300 rounded-[3px] bg-white hover:bg-gray-50 cursor-pointer">
            + New tenant
          </button>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Name", "Type", "Domain", "Clients", "Services", "Wallet", "Status", "Since"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {TENANTS.map((t) => (
              <tr key={t.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="px-3 py-2 font-mono text-[12px] text-[#D50C2D]">{t.id}</td>
                <td className="px-3 py-2 font-medium text-gray-900">{t.name}</td>
                <td className="px-3 py-2">
                  <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{t.type}</span>
                </td>
                <td className="px-3 py-2 font-mono text-[12px] text-gray-500">{t.domain}</td>
                <td className="px-3 py-2 tabular-nums text-right">{t.clients.toLocaleString()}</td>
                <td className="px-3 py-2 tabular-nums text-right">{t.services.toLocaleString()}</td>
                <td className="px-3 py-2 tabular-nums text-right">
                  {t.type === "admin" ? "—" : (
                    <span className={t.walletLow ? "text-red-600 font-medium" : ""}>{fmtMoney(t.wallet)}</span>
                  )}
                </td>
                <td className="px-3 py-2"><StatusBadge status={t.status} dot /></td>
                <td className="px-3 py-2 text-gray-400">{t.since ?? "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
