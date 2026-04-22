import { fmtMoney } from "@/mocks/sampleData";

const LEDGER = [
  { ts: "2026-04-22 14:02", type: "settlement.reseller.debit", amount: -62.00, ref: "ORD-48291 · VPS 4C/8G", balance: 4820.50 },
  { ts: "2026-04-21 10:18", type: "topup.credit.reseller", amount: 2000.00, ref: "TUP-9116 · VietQR", balance: 4882.50 },
  { ts: "2026-04-20 14:08", type: "settlement.reseller.debit", amount: -390.00, ref: "ORD-48280 · Residential batch", balance: 2882.50 },
  { ts: "2026-04-18 11:22", type: "settlement.reseller.debit", amount: -180.00, ref: "ORD-48270 · ISP batch", balance: 3272.50 },
];

export function ResellerWallet() {
  return (
    <div className="p-5 flex flex-col gap-4">
      {/* Balance card */}
      <div className="bg-white border border-gray-200 rounded p-5">
        <div className="flex items-start justify-between">
          <div>
            <div className="text-[11px] font-medium uppercase tracking-wide text-gray-400 mb-1">Wallet balance</div>
            <div className="text-3xl font-bold tabular-nums text-gray-900">$4,820.50</div>
            <div className="text-[12px] text-gray-400 mt-1">ProxyVN · T-0042</div>
          </div>
          <button className="h-8 px-4 text-[13px] font-medium bg-[#D50C2D] text-white rounded-[3px] border-0 hover:bg-[#B3082A] cursor-pointer">
            + Request top-up
          </button>
        </div>
        <div className="mt-4 pt-4 border-t border-gray-100 grid grid-cols-3 gap-4">
          {[
            { label: "Pending top-ups", value: "$2,000.00", sub: "TUP-9120 · awaiting admin" },
            { label: "Spent this month", value: "$8,240.00", sub: "settlement debits" },
            { label: "Low balance alert", value: "< $200", sub: "notify when below" },
          ].map(({ label, value, sub }) => (
            <div key={label}>
              <div className="text-[11px] text-gray-400 mb-0.5">{label}</div>
              <div className="text-[14px] font-semibold tabular-nums">{value}</div>
              <div className="text-[11px] text-gray-400">{sub}</div>
            </div>
          ))}
        </div>
      </div>

      {/* Ledger */}
      <div className="bg-white border border-gray-200 rounded">
        <div className="px-4 py-3 border-b border-gray-100">
          <h3 className="text-[13px] font-semibold text-gray-900 m-0">Ledger history</h3>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Timestamp", "Type", "Amount", "Reference", "Balance after"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 px-3 py-2 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {LEDGER.map((e, i) => (
              <tr key={i} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="px-3 py-2 tabular-nums text-gray-400">{e.ts}</td>
                <td className="px-3 py-2 font-mono text-[12px] text-gray-500">{e.type}</td>
                <td className={`px-3 py-2 tabular-nums text-right font-medium ${e.amount < 0 ? "text-red-600" : "text-green-700"}`}>
                  {fmtMoney(e.amount)}
                </td>
                <td className="px-3 py-2 text-gray-500">{e.ref}</td>
                <td className="px-3 py-2 tabular-nums text-right font-medium">{fmtMoney(e.balance)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
