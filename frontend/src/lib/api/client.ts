import { ApiActor, actorHeaders, apiBaseUrl, apiEnabled } from "./config";
import { ApiEnvelope, ApiJson, ApiQuery, ApiQueryValue } from "./types";

export class BillingApiError extends Error {
  constructor(message: string, public readonly status?: number) {
    super(message);
  }
}

export interface PostApiOptions {
  idempotencyKey?: string;
}

function hasQueryValue(value: ApiQueryValue): boolean {
  return value !== undefined && value !== null && String(value).trim() !== "";
}

function buildApiPath(path: string, query?: ApiQuery): string {
  if (!query) {
    return path;
  }
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(query as Record<string, ApiQueryValue>)) {
    if (!hasQueryValue(value)) {
      continue;
    }
    params.set(key, String(value).trim());
  }
  const queryString = params.toString();
  if (!queryString) {
    return path;
  }
  return `${path}?${queryString}`;
}

export function newIdempotencyKey(action: string): string {
  const label = action.trim().toLowerCase().replace(/[^a-z0-9_-]+/g, "-") || "mutation";
  const randomID = globalThis.crypto?.randomUUID?.();
  const fallbackID = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
  return `${label}:${randomID ?? fallbackID}`;
}

export async function getApiData<T>(path: string, actor: ApiActor, query?: ApiQuery): Promise<T> {
  if (!apiEnabled()) {
    throw new BillingApiError("API is not configured.");
  }
  const response = await fetch(apiBaseUrl() + buildApiPath(path, query), {
    method: "GET",
    headers: actorHeaders(actor),
    cache: "no-store",
  });
  const body = await response.text();
  if (!response.ok) {
    throw new BillingApiError(body || response.statusText, response.status);
  }
  const envelope = JSON.parse(body) as ApiEnvelope<T>;
  return envelope.data;
}

export async function postApiData<T>(
  path: string,
  actor: ApiActor,
  body: ApiJson = {},
  options: PostApiOptions = {},
): Promise<T> {
  if (!apiEnabled()) {
    throw new BillingApiError("API is not configured.");
  }
  const headers = new Headers(actorHeaders(actor));
  headers.set("Content-Type", "application/json");
  if (options.idempotencyKey) {
    headers.set("Idempotency-Key", options.idempotencyKey);
  }
  const response = await fetch(apiBaseUrl() + path, {
    method: "POST",
    headers,
    body: JSON.stringify(body ?? {}),
    cache: "no-store",
  });
  const responseBody = await response.text();
  if (!response.ok) {
    throw new BillingApiError(responseBody || response.statusText, response.status);
  }
  const envelope = JSON.parse(responseBody) as ApiEnvelope<T>;
  return envelope.data;
}
