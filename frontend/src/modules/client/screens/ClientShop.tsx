import { PRODUCTS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function ClientShop() {
  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-3 gap-4">
        {PRODUCTS.map((p) => (
          <div key={p.sku} className="bg-white border border-gray-200 rounded p-4 flex flex-col gap-4 hover:border-gray-300 transition-colors">
            <div>
              <div className="text-[13px] font-medium text-gray-900">{p.name}</div>
              <div className="text-[11px] text-gray-400 mt-0.5">{p.sku}</div>
            </div>
            <div className="flex items-baseline gap-1">
              <span className="text-lg font-medium tabular-nums text-gray-900">{fmtMoney(p.price)}</span>
              <span className="text-[12px] text-gray-400">{p.unit}</span>
            </div>
            <div className="text-[11px] text-gray-400">
              {p.active.toLocaleString()} active subscriptions
            </div>
            <button className="mt-auto h-8 w-full text-[13px] font-medium bg-[#D50C2D] text-white rounded-[3px] border-0 hover:bg-[#B3082A] cursor-pointer">
              Order now
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
