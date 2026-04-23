export type ApiActor = "admin" | "client";

const DEMO_TENANT_ID = "00000000-0000-0000-0000-000000000010";
const DEMO_RESELLER_ID = "00000000-0000-0000-0000-000000000102";
const DEMO_CUSTOMER_ID = "00000000-0000-0000-0000-000000000103";

export function apiBaseUrl(): string {
  return (process.env.NEXT_PUBLIC_BILLING_API_URL ?? "/backend").replace(/\/+$/, "");
}

export function apiEnabled(): boolean {
  return apiBaseUrl().length > 0;
}

export function actorHeaders(actor: ApiActor): HeadersInit {
  const actorId = actor === "admin" ? DEMO_RESELLER_ID : DEMO_CUSTOMER_ID;
  const actorType = actor === "admin" ? "reseller_owner" : "client";
  return {
    "X-Tenant-Id": DEMO_TENANT_ID,
    "X-Actor-Id": actorId,
    "X-Actor-Type": actorType,
    "X-Actor-Tenant-Id": DEMO_TENANT_ID,
  };
}
