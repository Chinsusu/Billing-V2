"use client";

import { useState } from "react";
import { RESELLER_CATALOG, type ResellerCatalogItem } from "@/mocks/billingData";
import { useToast } from "@/lib/toast/ToastContext";
import { fmtMoney } from "@/mocks/sampleData";

export function ResellerCatalog() {
  const { toast } = useToast();
  const [catalog, setCatalog] = useState<ResellerCatalogItem[]>(RESELLER_CATALOG);
  const [editing, setEditing] = useState<string | null>(null);
  const [editValue, setEditValue] = useState("");

  const marginWarnings = catalog.filter((i) => i.status === "warn").length;

  const startEdit = (item: ResellerCatalogItem) => {
    setEditing(item.plan);
    setEditValue(String(item.selling));
  };

  const commitEdit = (item: ResellerCatalogItem) => {
    const newPrice = parseFloat(editValue);
    if (isNaN(newPrice) || newPrice <= 0) { setEditing(null); return; }
    const margin = Math.round(((newPrice - item.cost) / item.cost) * 100);
    setCatalog((prev) =>
      prev.map((c) => c.plan === item.plan
        ? { ...c, selling: newPrice, margin, status: margin < 0 ? "warn" : "active" }
        : c
      ),
    );
    toast(`Price updated for "${item.plan}" → ${fmtMoney(newPrice)}`, "success");
    setEditing(null);
  };

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Catalog / Pricing</h3>
          {marginWarnings > 0 && (
            <span className="text-[11px] text-amber-600 font-medium">{marginWarnings} margin warning{marginWarnings > 1 ? "s" : ""}</span>
          )}
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Plan", "Unit", "Cost (from admin)", "Your price", "Margin", "Stock", "Status"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {catalog.map((item) => (
              <tr key={item.plan} className={`hover:bg-gray-50 border-b border-gray-100 last:border-0 ${item.status === "warn" ? "bg-amber-50/60" : ""}`}>
                <td className="p-4 font-medium text-gray-900">{item.plan}</td>
                <td className="p-4 text-gray-400 text-[12px]">{item.unit}</td>
                <td className="p-4 tabular-nums text-gray-500">{fmtMoney(item.cost)}</td>
                <td className="p-4">
                  {editing === item.plan ? (
                    <div className="flex items-center gap-1">
                      <input
                        type="number"
                        step="0.01"
                        min="0"
                        value={editValue}
                        onChange={(e) => setEditValue(e.target.value)}
                        onBlur={() => commitEdit(item)}
                        onKeyDown={(e) => { if (e.key === "Enter") commitEdit(item); if (e.key === "Escape") setEditing(null); }}
                        autoFocus
                        className="w-24 h-7 px-2 text-[13px] border border-blue-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-400"
                      />
                    </div>
                  ) : (
                    <button
                      onClick={() => startEdit(item)}
                      className="tabular-nums font-medium text-left hover:text-blue-600 cursor-pointer bg-transparent border-0 p-0 underline-offset-2 hover:underline"
                    >
                      {fmtMoney(item.selling)}
                    </button>
                  )}
                </td>
                <td className="p-4 tabular-nums">
                  <span className={item.margin < 0 ? "text-red-600 font-medium" : item.margin < 20 ? "text-amber-600" : "text-green-700"}>
                    {item.margin < 0 ? "" : "+"}{item.margin}%
                  </span>
                </td>
                <td className="p-4">
                  <span className={`text-[11px] px-1.5 py-px rounded-sm ${
                    item.stock === "ok" ? "bg-green-50 text-green-700 border border-green-200"
                    : item.stock === "low" ? "bg-amber-50 text-amber-700 border border-amber-200"
                    : "bg-red-50 text-red-700 border border-red-200"
                  }`}>{item.stock}</span>
                </td>
                <td className="p-4">
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
        <div className="p-3 border-t border-gray-100 text-[11px] text-gray-400">
          Click on a price to edit it inline. Press Enter to save, Esc to cancel.
        </div>
      </div>
    </div>
  );
}
