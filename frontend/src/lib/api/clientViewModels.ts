import { recordLabel } from "./format";
import type { ServiceInstance } from "./types";
import { hiddenReference } from "./viewModels";

export type ClientServiceCategory = "proxies" | "vps" | "bandwidth";

export function clientServiceCategory(service: ServiceInstance): ClientServiceCategory {
  const text = serviceSearchText(service);
  if (text.includes("bandwidth") || text.includes("traffic") || text.includes("gb")) return "bandwidth";
  if (text.includes("vps") || text.includes("vm") || text.includes("server")) return "vps";
  return "proxies";
}

export function clientServiceSourceLabel(service: ServiceInstance): string {
  const region = snapshotValue(service, ["region", "location", "datacenter"]);
  if (region) return region;
  if (service.provider_source_display_id) return recordLabel(service.provider_source_display_id, "SRC-");
  return hiddenReference("Source");
}

export function clientServiceOrderLabel(service: ServiceInstance): string {
  return service.order_display_id ? recordLabel(service.order_display_id, "ORD-") : hiddenReference("Order");
}

export function clientServicePlanLabel(service: ServiceInstance): string {
  return snapshotValue(service, ["name", "plan_name", "plan_code", "product_type"]) || hiddenReference("Plan");
}

function serviceSearchText(service: ServiceInstance): string {
  return [
    service.external_resource_id,
    snapshotText(service.product_snapshot),
    snapshotText(service.plan_snapshot),
    snapshotText(service.price_snapshot),
  ].filter(Boolean).join(" ").toLowerCase();
}

function snapshotValue(service: ServiceInstance, keys: string[]): string {
  for (const snapshot of [service.plan_snapshot, service.product_snapshot, service.price_snapshot]) {
    if (!snapshot || typeof snapshot !== "object" || Array.isArray(snapshot)) continue;
    for (const key of keys) {
      const value = (snapshot as Record<string, unknown>)[key];
      if (typeof value === "string" && value.trim()) return value.trim();
      if (typeof value === "number") return String(value);
    }
  }
  return "";
}

function snapshotText(value: unknown): string {
  if (!value) return "";
  if (typeof value === "string") return value;
  try {
    return JSON.stringify(value);
  } catch {
    return "";
  }
}
