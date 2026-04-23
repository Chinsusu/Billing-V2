import { ApiActor, actorHeaders, apiBaseUrl, apiEnabled } from "./config";
import { ApiEnvelope, ApiQuery, ApiQueryValue } from "./types";

export class BillingApiError extends Error {
  constructor(message: string, public readonly status?: number) {
    super(message);
  }
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
