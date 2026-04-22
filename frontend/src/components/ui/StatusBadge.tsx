import { STATUS_LABEL, STATUS_VARIANT } from "@/mocks/sampleData";

interface StatusBadgeProps {
  status: string;
  dot?: boolean;
}

const VARIANT_CLASSES: Record<string, string> = {
  ok: "bg-green-50 text-green-700 border-green-200",
  warn: "bg-amber-50 text-amber-700 border-amber-200",
  danger: "bg-red-50 text-red-700 border-red-200",
  info: "bg-blue-50 text-blue-700 border-blue-200",
  muted: "bg-gray-100 text-gray-500 border-transparent",
};

export function StatusBadge({ status, dot = false }: StatusBadgeProps) {
  const variant = STATUS_VARIANT[status] ?? "muted";
  const label = STATUS_LABEL[status] ?? status;
  return (
    <span className={`inline-flex items-center gap-1 px-1.5 py-px text-[11px] font-medium border rounded-sm ${VARIANT_CLASSES[variant]}`}>
      {dot && <span className="w-1.5 h-1.5 rounded-full bg-current" />}
      {label}
    </span>
  );
}
