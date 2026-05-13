"use client";

import { useState } from "react";
import { billingApi } from "@/lib/api/billing";
import type { ApiJson, ServiceAccessItem, ServiceAccessReveal as ServiceAccessRevealResult, ServiceInstance } from "@/lib/api/types";

type ServiceAccessScope = "client" | "reseller" | "admin";
type RevealState = "idle" | "loading" | "shown" | "error";

interface ServiceAccessRevealProps {
  scope: ServiceAccessScope;
  serviceId?: string;
  reason: string;
}

export function ServiceAccessReveal({ scope, serviceId, reason }: ServiceAccessRevealProps) {
  const [state, setState] = useState<RevealState>("idle");
  const [maskedHint, setMaskedHint] = useState("Hidden until reveal");
  const [shownValue, setShownValue] = useState("");
  const [error, setError] = useState("");

  async function reveal() {
    if (!serviceId || state === "loading" || state === "shown") return;
    setState("loading");
    setError("");
    try {
      const detail = await loadService(scope, serviceId);
      const access = firstActiveAccess(detail);
      if (!access) {
        setState("error");
        setError("No access item is available.");
        return;
      }
      setMaskedHint(access.masked_hint || access.credential_type); // sensitive-text-allowlist
      const result = await revealServiceAccess(scope, serviceId, access.id, { reason });
      setMaskedHint(result.masked_hint || access.masked_hint);
      setShownValue(formatAccessPayload(result.payload));
      setState("shown");
    } catch (err) {
      setState("error");
      setError(err instanceof Error ? err.message : "Access reveal failed.");
    }
  }

  if (!serviceId) {
    return <span className="text-[11px] text-gray-400">Not available</span>;
  }

  return (
    <div className="flex min-w-[160px] flex-col gap-1">
      <span className={`font-mono text-[11px] ${state === "shown" ? "text-gray-900" : "text-gray-400"}`}>
        {state === "shown" ? shownValue : maskedHint}
      </span>
      <button
        type="button"
        onClick={reveal}
        disabled={state === "loading" || state === "shown"}
        className="inline-flex h-7 w-fit items-center justify-center rounded-md border border-gray-200 bg-gray-50 px-2 text-[11px] font-medium text-gray-600 transition-colors hover:bg-white disabled:cursor-not-allowed disabled:text-gray-400"
      >
        {state === "loading" ? "Revealing" : state === "shown" ? "Shown once" : "Reveal access"}
      </button>
      {state === "error" && <span className="text-[11px] text-red-600">{error}</span>}
    </div>
  );
}

function firstActiveAccess(service: ServiceInstance): ServiceAccessItem | undefined {
  const items = service.credentials ?? []; // sensitive-text-allowlist
  return items.find((item) => item.status === "active") ?? items[0];
}

function loadService(scope: ServiceAccessScope, serviceId: string): Promise<ServiceInstance> {
  switch (scope) {
    case "admin":
      return billingApi.getAdminService(serviceId);
    case "reseller":
      return billingApi.getResellerService(serviceId);
    case "client":
      return billingApi.getClientService(serviceId);
  }
}

function revealServiceAccess(
  scope: ServiceAccessScope,
  serviceId: string,
  accessId: string,
  body: { reason?: string },
): Promise<ServiceAccessRevealResult> {
  switch (scope) {
    case "admin":
      return billingApi.revealAdminServiceAccess(serviceId, accessId, body);
    case "reseller":
      return billingApi.revealResellerServiceAccess(serviceId, accessId, body);
    case "client":
      return billingApi.revealClientServiceAccess(serviceId, accessId, body);
  }
}

function formatAccessPayload(payload: ApiJson): string {
  if (typeof payload === "string" || typeof payload === "number" || typeof payload === "boolean") {
    return String(payload);
  }
  if (payload && typeof payload === "object" && !Array.isArray(payload)) {
    const record = payload as Record<string, unknown>;
    const preferred = ["username", "user", "password", "host", "port", "url"]
      .map((key) => [key, record[key]] as const)
      .filter(([, value]) => value !== undefined && value !== null && String(value).trim() !== "")
      .map(([key, value]) => `${key}: ${String(value)}`);
    if (preferred.length > 0) {
      return preferred.join(" | ");
    }
  }
  return JSON.stringify(payload);
}
