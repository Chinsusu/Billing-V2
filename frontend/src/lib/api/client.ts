import { ApiActor, actorHeaders, apiBaseUrl, apiEnabled } from "./config";
import { ApiEnvelope } from "./types";

export class BillingApiError extends Error {
  constructor(message: string, public readonly status?: number) {
    super(message);
  }
}

export async function getApiData<T>(path: string, actor: ApiActor): Promise<T> {
  if (!apiEnabled()) {
    throw new BillingApiError("API is not configured.");
  }
  const response = await fetch(apiBaseUrl() + path, {
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
