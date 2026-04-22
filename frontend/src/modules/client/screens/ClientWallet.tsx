import { CLIENT_LEDGER } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function ClientWallet() {
  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="bg-white border border-gray-200 rounded p-4 flex items-start justify-between">
        <div>
          <div className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Available balance</div>
          <div className="text-3xl font-medium tabular-nums text-gray-900">$128.40</div>
          <div className="text-[12px] text-gray-400 mt-1">Linh Tran · via ProxyVN</div>
        </div>
        <button className="h-8 p-4 text-[13px] font-medium bg-[#D50C2D] text-white rounded-[3px] border-0 hover:bg-[#B3082A] cursor-pointer">
          + Top up
        </button>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Transaction history</h3>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Timestamp", "Type", "Amount", "Reference", "Balance after"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {CLIENT_LEDGER.map((e, i) => (
              <tr key={i} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 p-4 tabular-nums text-gray-400">{e.ts}</td>
                <td className="p-4 p-4 text-[12px] text-gray-500">{e.type}</td>
                <td className={`p-4 p-4 tabular-nums text-right font-medium ${e.amount < 0 ? "text-red-600" : "text-green-700"}`}>
                  {fmtMoney(e.amount)}
                </td>
                <td className="p-4 p-4 text-gray-500">{e.ref}</td>
                <td className="p-4 p-4 tabular-nums text-right font-medium">{fmtMoney(e.balance)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
