import { PRODUCTS } from "@/mocks/billingData";
import { fmtMoney, fmtMoneyShort } from "@/mocks/sampleData";

export function AdminProducts() {
  return (
    <div className="p-5">
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Products & Pricing</h3>
          <button className="h-7 px-3 text-[12px] font-medium border border-gray-300 rounded-[3px] bg-white hover:bg-gray-50 cursor-pointer">
            + Add product
          </button>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["SKU", "Name", "Unit", "Price", "Active", "Rev 30d"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {PRODUCTS.map((p) => (
              <tr key={p.sku} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="px-3 py-2 font-mono text-[12px] text-gray-500">{p.sku}</td>
                <td className="px-3 py-2 font-medium text-gray-900">{p.name}</td>
                <td className="px-3 py-2 text-gray-400 text-[12px]">{p.unit}</td>
                <td className="px-3 py-2 tabular-nums font-medium">{fmtMoney(p.price)}</td>
                <td className="px-3 py-2 tabular-nums text-right">{p.active.toLocaleString()}</td>
                <td className="px-3 py-2 tabular-nums text-right font-medium">{fmtMoneyShort(p.rev30)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
