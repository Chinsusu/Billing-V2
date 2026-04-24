import { StatusBadge } from "@/components/ui/StatusBadge";
import { BANDWIDTH_SERVICES, PROXY_SERVICES, VPS_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

type ResellerServiceCategory = "proxies" | "vps" | "bandwidth";

interface ResellerServicesProps {
  category: ResellerServiceCategory;
}

const CONFIG: Record<ResellerServiceCategory, { title: string; empty: string }> = {
  proxies: { title: "Proxy services", empty: "No proxy services" },
  vps: { title: "VPS services", empty: "No VPS services" },
  bandwidth: { title: "Bandwidth usage", empty: "No bandwidth records" },
};

export function ResellerServices({ category }: ResellerServicesProps) {
  const rows = serviceRows(category);
  const attention = rows.filter((row) => row.status === "suspended" || row.status === "overdue").length;
  const revenue = rows.reduce((total, row) => total + row.price, 0);
  const config = CONFIG[category];

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <SummaryTile label="Records" value={String(rows.length)} />
        <SummaryTile label="Monthly value" value={fmtMoney(revenue)} />
        <SummaryTile label="Attention" value={String(attention)} tone={attention > 0 ? "warn" : "neutral"} />
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between gap-3">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">{config.title}</h3>
          <span className="text-[11px] text-gray-400">Demo adapter</span>
        </div>
        <div className="overflow-x-auto max-w-full">
          <table className="w-full text-[13px] border-collapse min-w-[760px]">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Service", "Client", "Region", "Usage", "Renewal", "Price", "Status"].map((heading) => (
                  <th key={heading} className="text-left text-[11px] font-medium uppercase text-gray-400 p-4 border-b border-gray-200">
                    {heading}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 text-[12px] text-[#D50C2D] font-medium">{row.id}</td>
                  <td className="p-4 font-medium text-gray-900">{row.label}</td>
                  <td className="p-4 text-gray-500">{row.customer}</td>
                  <td className="p-4 text-gray-500">{row.region}</td>
                  <td className="p-4 text-gray-500">{row.usage}</td>
                  <td className="p-4 text-gray-500">{row.renewal}</td>
                  <td className="p-4 text-right font-medium tabular-nums">{fmtMoney(row.price)}</td>
                  <td className="p-4"><StatusBadge status={row.status} dot /></td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr><td colSpan={8} className="p-4 text-center text-[12px] text-gray-400">{config.empty}</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function serviceRows(category: ResellerServiceCategory) {
  if (category === "vps") {
    return VPS_SERVICES.filter((item) => item.tenant === "ProxyVN").map((item) => ({
      id: item.id,
      label: item.label,
      customer: item.customer,
      region: item.region,
      usage: `${item.cpu}C / ${item.ram}GB / ${item.disk}GB`,
      renewal: item.renewsIn < 0 ? `${Math.abs(item.renewsIn)}d overdue` : `${item.renewsIn}d`,
      price: item.price,
      status: item.status,
    }));
  }
  if (category === "bandwidth") {
    return BANDWIDTH_SERVICES.filter((item) => item.tenant === "ProxyVN").map((item) => ({
      id: item.id,
      label: item.label,
      customer: item.customer,
      region: item.region,
      usage: `${item.usedGB} / ${item.totalGB} GB`,
      renewal: item.renewsIn < 0 ? `${Math.abs(item.renewsIn)}d overdue` : `${item.renewsIn}d`,
      price: item.price,
      status: item.status,
    }));
  }
  return PROXY_SERVICES.filter((item) => item.tenant === "ProxyVN").map((item) => ({
    id: item.id,
    label: item.label,
    customer: item.customer,
    region: item.region,
    usage: item.ipCount > 0 ? `${item.ipCount} IPs` : `${item.usedGB} / ${item.totalGB} GB`,
    renewal: item.renewsIn < 0 ? `${Math.abs(item.renewsIn)}d overdue` : `${item.renewsIn}d`,
    price: item.price,
    status: item.status,
  }));
}

function SummaryTile({ label, value, tone = "neutral" }: { label: string; value: string; tone?: "neutral" | "warn" }) {
  return (
    <div className={`bg-white border rounded p-4 ${tone === "warn" ? "border-amber-200" : "border-gray-200"}`}>
      <div className="text-[11px] text-gray-400 uppercase mb-1">{label}</div>
      <div className={`text-lg font-medium tabular-nums ${tone === "warn" ? "text-amber-700" : "text-gray-900"}`}>{value}</div>
    </div>
  );
}
