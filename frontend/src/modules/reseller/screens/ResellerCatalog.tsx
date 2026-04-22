import { RESELLER_CATALOG } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function ResellerCatalog() {
  return (
    <div className="p-5">
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Catalog / Pricing</h3>
          <span className="text-[11px] text-amber-600 font-medium">1 margin warning</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Plan", "Unit", "Cost (from admin)", "Your price", "Margin", "Stock", "Status"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {RESELLER_CATALOG.map((item) => (
              <tr key={item.plan} className={`hover:bg-gray-50 border-b border-gray-100 last:border-0 ${item.status === "warn" ? "bg-amber-50/60" : ""}`}>
                <td className="px-3 py-2 font-medium text-gray-900">{item.plan}</td>
                <td className="px-3 py-2 text-gray-400 text-[12px]">{item.unit}</td>
                <td className="px-3 py-2 tabular-nums text-gray-500">{fmtMoney(item.cost)}</td>
                <td className="px-3 py-2 tabular-nums font-medium">{fmtMoney(item.selling)}</td>
                <td className="px-3 py-2 tabular-nums">
                  <span className={item.margin < 0 ? "text-red-600 font-semibold" : item.margin < 20 ? "text-amber-600" : "text-green-700"}>
                    {item.margin < 0 ? "" : "+"}{item.margin}%
                  </span>
                </td>
                <td className="px-3 py-2">
                  <span className={`text-[11px] px-1.5 py-px rounded-sm ${
                    item.stock === "ok" ? "bg-green-50 text-green-700 border border-green-200"
                    : item.stock === "low" ? "bg-amber-50 text-amber-700 border border-amber-200"
                    : "bg-red-50 text-red-700 border border-red-200"
                  }`}>{item.stock}</span>
                </td>
                <td className="px-3 py-2">
                  <span className={`text-[11px] px-1.5 py-px rounded-sm ${
                    item.status === "active" ? "bg-green-50 text-green-700 border border-green-200"
                    : item.status === "warn" ? "bg-amber-50 text-amber-700 border border-amber-200"
                    : "bg-gray-100 text-gray-400 border border-transparent"
                  }`}>{item.status}</span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
