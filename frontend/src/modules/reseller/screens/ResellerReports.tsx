import { RESELLER_CATALOG, RESELLER_CLIENTS } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

export function ResellerReports() {
  const walletTotal = RESELLER_CLIENTS.reduce((total, client) => total + client.wallet, 0);
  const serviceTotal = RESELLER_CLIENTS.reduce((total, client) => total + client.services, 0);
  const negativeMargins = RESELLER_CATALOG.filter((item) => item.margin < 0).length;

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <SummaryTile label="Client wallet total" value={fmtMoney(walletTotal)} />
        <SummaryTile label="Active services" value={serviceTotal.toLocaleString()} />
        <SummaryTile label="Margin warnings" value={String(negativeMargins)} tone={negativeMargins > 0 ? "warn" : "neutral"} />
      </div>
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Profit report</h3>
          <span className="text-[11px] text-gray-400">Backend route pending</span>
        </div>
        <div className="p-4 grid grid-cols-1 lg:grid-cols-2 gap-4">
          {RESELLER_CATALOG.map((item) => (
            <div key={item.plan} className="border border-gray-200 rounded p-4">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <div className="text-[13px] font-medium text-gray-900">{item.plan}</div>
                  <div className="text-[11px] text-gray-400 mt-0.5">{item.unit}</div>
                </div>
                <div className={`text-[13px] font-medium tabular-nums ${item.margin < 0 ? "text-red-600" : "text-green-700"}`}>
                  {item.margin < 0 ? "" : "+"}{item.margin}%
                </div>
              </div>
              <div className="mt-3 grid grid-cols-2 gap-3 text-[12px]">
                <div>
                  <div className="text-gray-400">Cost</div>
                  <div className="font-medium tabular-nums">{fmtMoney(item.cost)}</div>
                </div>
                <div>
                  <div className="text-gray-400">Selling</div>
                  <div className="font-medium tabular-nums">{fmtMoney(item.selling)}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function SummaryTile({ label, value, tone = "neutral" }: { label: string; value: string; tone?: "neutral" | "warn" }) {
  return (
    <div className={`bg-white border rounded p-4 ${tone === "warn" ? "border-amber-200" : "border-gray-200"}`}>
      <div className="text-[11px] text-gray-400 uppercase mb-1">{label}</div>
      <div className={`text-lg font-medium tabular-nums ${tone === "warn" ? "text-amber-700" : "text-gray-900"}`}>{value}</div>
    </div>
  );
}
