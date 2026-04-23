"use client";

import { useState } from "react";
import { PRODUCTS, type ProductCatalog } from "@/mocks/billingData";
import { ConfirmDialog } from "@/components/ui/ConfirmDialog";
import { useToast } from "@/lib/toast/ToastContext";
import { fmtMoney } from "@/mocks/sampleData";

export function ClientShop() {
  const { toast } = useToast();
  const [ordering, setOrdering] = useState<ProductCatalog | null>(null);

  const handleOrder = () => {
    if (!ordering) return;
    toast(`Order placed for "${ordering.name}" — check My Services`, "success");
    setOrdering(null);
  };

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
            <button
              onClick={() => setOrdering(p)}
              className="mt-auto w-full inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm"
            >
              Order now
            </button>
          </div>
        ))}
      </div>

      <ConfirmDialog
        open={!!ordering}
        onClose={() => setOrdering(null)}
        onConfirm={handleOrder}
        title="Confirm order"
        description={ordering ? `Order "${ordering.name}" at ${fmtMoney(ordering.price)} ${ordering.unit}? The cost will be deducted from your wallet.` : ""}
        confirmLabel="Place order"
      />
    </div>
  );
}
