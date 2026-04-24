import type { ProviderReadinessState } from "@/lib/api/types";

const STATE_LABEL: Record<ProviderReadinessState, string> = {
  ready: "Ready",
  inactive_source: "Inactive source",
  missing_plan_source: "Missing source",
  unsupported_capability: "Unsupported",
  fake_provider_only: "Fake only",
};

const STATE_CLASS: Record<ProviderReadinessState, string> = {
  ready: "border-emerald-200 bg-emerald-50 text-emerald-700",
  inactive_source: "border-amber-200 bg-amber-50 text-amber-700",
  missing_plan_source: "border-red-200 bg-red-50 text-red-700",
  unsupported_capability: "border-red-200 bg-red-50 text-red-700",
  fake_provider_only: "border-blue-200 bg-blue-50 text-blue-700",
};

export function ProviderReadinessStateBadge({ state }: { state: ProviderReadinessState }) {
  return (
    <span className={`inline-flex items-center rounded-sm border px-1.5 py-px text-[11px] font-medium ${STATE_CLASS[state] ?? "border-gray-200 bg-gray-100 text-gray-500"}`}>
      {STATE_LABEL[state] ?? state}
    </span>
  );
}
