import { AlertTriangle, Info, XCircle, CheckCircle } from "lucide-react";
import { ALERTS, PlatformAlert } from "@/mocks/billingData";

const SEV_STYLE: Record<PlatformAlert["severity"], string> = {
  danger: "bg-red-50 border-red-200 text-red-700",
  warn:   "bg-amber-50 border-amber-200 text-amber-700",
  info:   "bg-blue-50 border-blue-200 text-blue-700",
};

const SEV_ICON: Record<PlatformAlert["severity"], React.ReactNode> = {
  danger: <XCircle size={14} />,
  warn:   <AlertTriangle size={14} />,
  info:   <Info size={14} />,
};

const CAT_LABEL: Record<PlatformAlert["category"], string> = {
  provisioning: "Provisioning",
  provider:     "Provider",
  billing:      "Billing",
  security:     "Security",
  system:       "System",
};

export function AdminAlerts() {
  const open     = ALERTS.filter((a) => !a.resolved);
  const resolved = ALERTS.filter((a) => a.resolved);

  return (
    <div className="p-4 flex flex-col gap-4">
      {/* Summary strip */}
      <div className="flex gap-4">
        {(["danger", "warn", "info"] as const).map((sev) => {
          const count = open.filter((a) => a.severity === sev).length;
          return (
            <div key={sev} className={`flex items-center gap-4 p-4 p-4.5 rounded border text-[12px] font-medium ${SEV_STYLE[sev]}`}>
              {SEV_ICON[sev]}
              <span>{count} {sev === "danger" ? "critical" : sev === "warn" ? "warning" : "info"}</span>
            </div>
          );
        })}
      </div>

      {/* Open alerts */}
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Open alerts</h3>
        </div>
        <div className="divide-y divide-gray-100">
          {open.map((a) => (
            <div key={a.id} className="flex items-start gap-4 p-4 p-4">
              <div className={`mt-0.5 shrink-0 ${a.severity === "danger" ? "text-red-500" : a.severity === "warn" ? "text-amber-500" : "text-blue-500"}`}>
                {SEV_ICON[a.severity]}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-4 flex-wrap">
                  <span className="text-[12px] font-medium text-gray-900">{a.title}</span>
                  <span className="text-[10px] font-medium uppercase tracking-wide px-1.5 py-0.5 rounded bg-gray-100 text-gray-500">
                    {CAT_LABEL[a.category]}
                  </span>
                </div>
                <div className="text-[12px] text-gray-500 mt-0.5">{a.detail}</div>
              </div>
              <div className="shrink-0 text-right">
                <div className="text-[11px] text-gray-400">{a.ts}</div>
                <div className="text-[11px] text-gray-300 mt-0.5">{a.id}</div>
              </div>
            </div>
          ))}
          {open.length === 0 && (
            <div className="p-4 py-6 text-center text-[12px] text-gray-400">No open alerts</div>
          )}
        </div>
      </div>

      {/* Resolved alerts */}
      {resolved.length > 0 && (
        <div className="bg-white border border-gray-200 rounded">
          <div className="p-4 p-4 border-b border-gray-100 flex items-center gap-4">
            <CheckCircle size={13} className="text-green-500" />
            <h3 className="text-[13px] font-medium text-gray-900 m-0">Resolved</h3>
          </div>
          <div className="divide-y divide-gray-100">
            {resolved.map((a) => (
              <div key={a.id} className="flex items-start gap-4 p-4 p-4 opacity-60">
                <div className="mt-0.5 shrink-0 text-green-500"><CheckCircle size={14} /></div>
                <div className="flex-1 min-w-0">
                  <div className="text-[12px] font-medium text-gray-700 line-through">{a.title}</div>
                  <div className="text-[12px] text-gray-400 mt-0.5">{a.detail}</div>
                </div>
                <div className="text-[11px] text-gray-400 shrink-0">{a.ts}</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
