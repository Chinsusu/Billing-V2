export type ApiActor = "admin" | "reseller" | "client";

const DEMO_TENANT_ID = "00000000-0000-0000-0000-000000000010";
const DEMO_RESELLER_ID = "00000000-0000-0000-0000-000000000102";
const DEMO_CUSTOMER_ID = "00000000-0000-0000-0000-000000000103";

const ACTOR_PROFILES: Record<ApiActor, { id: string; type: string }> = {
  admin: { id: DEMO_RESELLER_ID, type: "reseller_owner" },
  reseller: { id: DEMO_RESELLER_ID, type: "reseller_owner" },
  client: { id: DEMO_CUSTOMER_ID, type: "client" },
};

export function apiBaseUrl(): string {
  return (process.env.NEXT_PUBLIC_BILLING_API_URL ?? "/backend").replace(/\/+$/, "");
}

export function apiEnabled(): boolean {
  return apiBaseUrl().length > 0;
}

export function actorHeaders(actor: ApiActor): HeadersInit {
  const profile = ACTOR_PROFILES[actor];
  return {
    "X-Tenant-Id": DEMO_TENANT_ID,
    "X-Actor-Id": profile.id,
    "X-Actor-Type": profile.type,
    "X-Actor-Tenant-Id": DEMO_TENANT_ID,
  };
}
