import { apiBaseUrl, apiEnabled, loginTenantID } from "./config";
import { ApiEnvelope, ApiJson } from "./types";
import { BillingApiError, apiErrorMessage } from "./client";

export type AuthActorType = "platform_admin" | "platform_staff" | "reseller_owner" | "reseller_staff" | "client";

export interface AuthSession {
  session_id: string;
  user_id: string;
  tenant_id: string;
  actor_type: AuthActorType;
  expires_at: string;
  two_factor_required: boolean;
  two_factor_satisfied: boolean;
  two_factor_setup_required: boolean;
}

export interface TwoFactorSetup {
  method: string;
  secret: string; // sensitive-text-allowlist: backend auth contract field, displayed only after login.
  provision_uri: string;
}

export interface TwoFactorVerification {
  session_id: string;
  user_id: string;
  tenant_id: string;
  two_factor_satisfied: boolean;
}

export async function login(email: string, password: string): Promise<AuthSession> {
  const headers = new Headers({ "Content-Type": "application/json" });
  const tenantID = loginTenantID();
  if (tenantID) {
    headers.set("X-Tenant-Id", tenantID);
  }
  return postAuthData<AuthSession>("/auth/login", { email, password }, headers);
}

export async function logout(): Promise<void> {
  await postAuthData<{ status: string }>("/auth/logout", {});
}

export async function setupTwoFactor(): Promise<TwoFactorSetup> {
  return postAuthData<TwoFactorSetup>("/auth/2fa/setup", {});
}

export async function verifyTwoFactor(code: string): Promise<TwoFactorVerification> {
  return postAuthData<TwoFactorVerification>("/auth/2fa/verify", { code });
}

async function postAuthData<T>(path: string, body: ApiJson, headers = new Headers({ "Content-Type": "application/json" })): Promise<T> {
  if (!apiEnabled()) {
    throw new BillingApiError("API is not configured.");
  }
  const response = await fetch(apiBaseUrl() + path, {
    method: "POST",
    headers,
    body: JSON.stringify(body ?? {}),
    cache: "no-store",
    credentials: "include", // sensitive-text-allowlist: required for HttpOnly session cookies.
  });
  const responseBody = await response.text();
  if (!response.ok) {
    throw new BillingApiError(apiErrorMessage(responseBody, response.statusText), response.status);
  }
  const envelope = JSON.parse(responseBody) as ApiEnvelope<T>;
  return envelope.data;
}
