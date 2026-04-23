"use client";

import { useState } from "react";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { ConfirmDialog } from "@/components/ui/ConfirmDialog";
import { useToast } from "@/lib/toast/ToastContext";
import { CLIENT_SERVICES, type ClientService } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function ClientDashboard() {
  const { toast } = useToast();
  const [services, setServices] = useState<ClientService[]>(CLIENT_SERVICES);
  const [renewing, setRenewing] = useState<ClientService | null>(null);
  const [balance] = useState(128.40);

  const suspended = services.filter((s) => s.status === "suspended");

  const handleRenew = () => {
    if (!renewing) return;
    setServices((prev) =>
      prev.map((s) => s.id === renewing.id ? { ...s, status: "active", note: undefined } : s),
    );
    toast(`Service "${renewing.label}" renewed successfully`, "success");
    setRenewing(null);
  };

  return (
    <div className="p-4 flex flex-col gap-4">
      {suspended.map((s) => (
        <div key={s.id} className="bg-amber-50 border border-amber-200 text-amber-700 text-[12px] p-3.5 rounded flex items-center gap-3">
          <span>⚠</span>
          <span><strong>{s.label}</strong> is suspended. {s.note} Renew to restore access.</span>
          <button
            onClick={() => setRenewing(s)}
            className="ml-auto inline-flex items-center justify-center gap-2 px-4 h-8 text-[12px] font-medium bg-amber-600 hover:bg-amber-700 text-white rounded border-0 cursor-pointer transition-colors"
          >
            Renew now
          </button>
        </div>
      ))}

      <div className="bg-white border border-gray-200 rounded p-4 flex items-center justify-between">
        <div>
          <div className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Wallet balance</div>
          <div className="text-lg font-medium tabular-nums text-gray-900">{fmtMoney(balance)}</div>
        </div>
        <button className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm">
          + Top up
        </button>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">My services</h3>
          <span className="text-[12px] text-[#D50C2D] cursor-pointer">View all →</span>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Label", "Region", "Bandwidth", "Expires", "Status"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {services.map((s) => (
              <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 font-medium text-gray-900">{s.label}</td>
                <td className="p-4 text-[12px] text-gray-400">{s.region}</td>
                <td className="p-4 text-gray-500">{s.bandwidth}</td>
                <td className="p-4 text-gray-400">{s.expiry}</td>
                <td className="p-4"><StatusBadge status={s.status} dot /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <ConfirmDialog
        open={!!renewing}
        onClose={() => setRenewing(null)}
        onConfirm={handleRenew}
        title="Renew service"
        description={renewing ? `Renew "${renewing.label}"? The renewal cost will be deducted from your wallet.` : ""}
        confirmLabel="Renew"
      />
    </div>
  );
}
