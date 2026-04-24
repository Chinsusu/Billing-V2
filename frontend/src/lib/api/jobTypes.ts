import type { PageQuery } from "./types";

export interface ProvisioningJobAttempt {
  id: string;
  display_id: number;
  job_id: string;
  worker_id: string;
  attempt_number: number;
  started_at: string;
  finished_at?: string;
  result: string;
  error_code?: string;
  error_message_redacted?: string;
  duration_ms?: number;
  correlation_id?: string;
}

export type JobAttemptQuery = PageQuery;
