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

export interface JobSummaryCounts {
  queued: number;
  claimed: number;
  running: number;
  succeeded: number;
  failed_retryable: number;
  failed_terminal: number;
  manual_review: number;
  cancelled: number;
}

export interface JobSummaryFailure {
  id: string;
  display_id: number;
  status: string;
  last_error_code?: string;
  last_error_message_redacted?: string;
  manual_review_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface JobSummary {
  job_type: string;
  total: number;
  attention_count: number;
  counts: JobSummaryCounts;
  oldest_queued_at?: string;
  oldest_queued_age_seconds?: number;
  latest_failure?: JobSummaryFailure;
  generated_at: string;
}

export interface JobSummaryQuery {
  job_type?: string;
}
