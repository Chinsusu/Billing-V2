import { ApiActor, actorHeaders, apiBaseUrl, apiEnabled, devActorHeadersEnabled } from "./config";
import { ApiEnvelope, ApiJson, ApiQuery, ApiQueryValue } from "./types";

export class BillingApiError extends Error {
  constructor(message: string, public readonly status?: number) {
    super(message);
  }
}

interface ApiErrorBody {
  error?: {
    code?: string;
    message?: string;
    fields?: Array<{ message?: string }>;
  };
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
  const headers = devActorHeadersEnabled() ? actorHeaders(actor) : undefined;
  const response = await fetch(apiBaseUrl() + buildApiPath(path, query), {
    method: "GET",
    headers,
    cache: "no-store",
    credentials: "include", // sensitive-text-allowlist: required for HttpOnly session cookies.
  });
  const body = await response.text();
  if (!response.ok) {
    throw new BillingApiError(apiErrorMessage(body, response.statusText), response.status);
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
  const headers = new Headers(devActorHeadersEnabled() ? actorHeaders(actor) : undefined);
  headers.set("Content-Type", "application/json");
  if (options.idempotencyKey) {
    headers.set("Idempotency-Key", options.idempotencyKey);
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

export function apiErrorMessage(body: string, fallback: string): string {
  if (!body) {
    return fallback;
  }
  try {
    const parsed = JSON.parse(body) as ApiErrorBody;
    const fieldMessage = parsed.error?.fields?.find((field) => field.message)?.message;
    return fieldMessage ?? parsed.error?.message ?? fallback;
  } catch {
    return body;
  }
}
