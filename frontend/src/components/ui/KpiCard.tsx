interface KpiCardProps {
  label: string;
  value: string;
  delta?: number;
  sub?: string;
}

export function KpiCard({ label, value, delta, sub }: KpiCardProps) {
  const positive = delta == null || delta >= 0;
  return (
    <div className="bg-white border border-gray-200 rounded p-4 flex flex-col gap-2 min-w-0">
      <div className="text-[11px] font-medium text-gray-400 uppercase tracking-wide">{label}</div>
      <div className="text-[22px] font-semibold tracking-tight text-gray-900 tabular-nums">{value}</div>
      <div className="flex items-center gap-1.5 min-h-[18px]">
        {delta != null && (
          <span className={`text-[11px] font-medium ${positive ? "text-green-700" : "text-red-600"}`}>
            {positive ? "↑" : "↓"} {Math.abs(delta)}%
          </span>
        )}
        {sub && <span className="text-[11px] text-gray-400">{sub}</span>}
      </div>
    </div>
  );
}
