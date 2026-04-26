import { StatusBadge } from "@/components/ui/StatusBadge";
import { technicalCodeLabel } from "@/lib/api/displayLabels";
import { compactDateTime, recordLabel } from "@/lib/api/format";
import type { JobSummary, JobSummaryFailure } from "@/lib/api/jobTypes";

interface AdminProvisioningSummaryPanelProps {
  summary: JobSummary | null;
  loading: boolean;
  error: string | null;
}

interface SummaryTile {
  label: string;
  value: number;
  status: string;
  sub?: string;
  attention?: boolean;
}

export function AdminProvisioningSummaryPanel({ summary, loading, error }: AdminProvisioningSummaryPanelProps) {
  const tiles = provisioningSummaryTiles(summary);
  return (
    <section className="bg-white border border-gray-200 rounded">
      <div className="flex flex-col gap-3 border-b border-gray-100 p-4 md:flex-row md:items-center md:justify-between">
        <div className="min-w-0">
          <h3 className="m-0 text-[13px] font-medium text-gray-900">Provisioning health</h3>
          <p className="m-0 mt-1 text-[12px] text-gray-400">
            {summary ? `${summary.total} job(s), ${summary.attention_count} need attention` : "Live summary"}
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2 text-[11px] text-gray-400">
          {loading && <span>Loading summary...</span>}
          {summary?.generated_at && <span>Updated {compactDateTime(summary.generated_at)}</span>}
          {summary?.latest_failure && (
            <span className="inline-flex items-center gap-2 rounded-sm border border-red-100 bg-red-50 px-2 py-1 text-red-700">
              Latest {recordLabel(summary.latest_failure.display_id, "JOB-")}
              <StatusBadge status={summary.latest_failure.status} />
            </span>
          )}
        </div>
      </div>
      {error && (
        <div className="border-b border-amber-100 bg-amber-50 px-4 py-3 text-[12px] text-amber-700">
          Summary API unavailable. Queue table remains visible.
        </div>
      )}
      <div className="grid gap-3 p-4 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-6">
        {loading && !summary
          ? Array.from({ length: 6 }).map((_, index) => <SummarySkeleton key={index} />)
          : tiles.map((tile) => <SummaryStatusTile key={tile.label} tile={tile} />)}
      </div>
      {summary?.latest_failure && (
        <div className="border-t border-gray-100 px-4 py-3 text-[12px] text-gray-500">
          <span className="font-medium text-gray-700">Latest issue: </span>
          {latestFailureIssue(summary.latest_failure)}
        </div>
      )}
    </section>
  );
}

function provisioningSummaryTiles(summary: JobSummary | null): SummaryTile[] {
  const counts = summary?.counts;
  return [
    {
      label: "Queued",
      value: counts?.queued ?? 0,
      status: "queued",
      sub: summary?.oldest_queued_age_seconds ? `oldest ${ageLabel(summary.oldest_queued_age_seconds)}` : undefined,
    },
    { label: "Running", value: counts?.running ?? 0, status: "running", sub: claimedLabel(counts?.claimed ?? 0) },
    { label: "Retryable", value: counts?.failed_retryable ?? 0, status: "failed_retryable", attention: true },
    { label: "Manual review", value: counts?.manual_review ?? 0, status: "manual_review", attention: true },
    { label: "Terminal", value: counts?.failed_terminal ?? 0, status: "failed_terminal", attention: true },
    { label: "Cancelled", value: counts?.cancelled ?? 0, status: "cancelled" },
  ];
}

function SummaryStatusTile({ tile }: { tile: SummaryTile }) {
  const attentionClass = tile.attention && tile.value > 0 ? "border-amber-200 bg-amber-50/50" : "border-gray-200 bg-white";
  return (
    <div className={`flex min-h-[94px] min-w-0 flex-col justify-between rounded border p-3 ${attentionClass}`}>
      <div className="flex items-start justify-between gap-2">
        <span className="text-[11px] font-medium uppercase tracking-wide text-gray-400">{tile.label}</span>
        <StatusBadge status={tile.status} dot />
      </div>
      <div className="text-[24px] font-medium tracking-tight text-gray-900 tabular-nums">{tile.value}</div>
      <div className="min-h-[16px] truncate text-[11px] text-gray-400">{tile.sub ?? " "}</div>
    </div>
  );
}

function SummarySkeleton() {
  return (
    <div className="flex min-h-[94px] flex-col justify-between rounded border border-gray-200 bg-gray-50 p-3">
      <div className="h-3 w-20 rounded bg-gray-200" />
      <div className="h-7 w-12 rounded bg-gray-200" />
      <div className="h-3 w-24 rounded bg-gray-200" />
    </div>
  );
}

function claimedLabel(value: number): string | undefined {
  return value > 0 ? `${value} claimed` : undefined;
}

function ageLabel(seconds: number): string {
  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m`;
  const hours = Math.floor(minutes / 60);
  const restMinutes = minutes % 60;
  return restMinutes > 0 ? `${hours}h ${restMinutes}m` : `${hours}h`;
}

function latestFailureIssue(failure: JobSummaryFailure): string {
  const values = [
    technicalCodeLabel(failure.last_error_code),
    failure.last_error_message_redacted,
    failure.manual_review_reason,
  ].filter(Boolean);
  return values.join(" / ") || "No redacted failure detail.";
}
