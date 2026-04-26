// Sample time-series data for charts
export const SAMPLE_SERIES = {
  revenue30d: [12400,11800,13100,12900,14200,15100,14800,16300,15900,17200,16800,18100,17900,19300,18800,20100,21400,20900,22300,22100,23800,23400,24900,24600,26200,25800,27400,27100,28600,29200],
  mrr30d: [86000,87200,87800,89100,90400,91200,92500,93100,94800,95400,96700,97300,98200,99500,100200,101400,102800,103500,104900,105600,107200,108100,109400,110200,111700,112800,114100,115300,116800,118200],
  customers30d: [2640,2658,2671,2682,2694,2706,2718,2731,2744,2752,2763,2775,2781,2792,2803,2811,2819,2828,2836,2840,2843,2847,2849,2851,2853,2854,2846,2847,2847,2847],
  bandwidthDaily: [184,192,201,198,215,223,219,234,229,245,251,248,262,268,271,284,278,291,295,302,298,311,315,309,322,318,331,327,339,342],
};

export const STATUS_LABEL: Record<string, string> = {
  active: "Active", running: "Running", paid: "Paid", open: "Open",
  pending: "Pending", overdue: "Overdue", failed: "Failed",
  suspended: "Suspended", stopped: "Stopped", provisioning: "Provisioning",
  manual_review: "Manual Review", queued: "Queued", claimed: "Claimed",
  succeeded: "Succeeded", failed_retryable: "Retryable", failed_terminal: "Terminal Failed",
  cancelled: "Cancelled", unknown: "Unknown",
  pending_verification: "Pending verification", submitted: "Submitted", under_review: "Under review",
  approved: "Approved", rejected: "Rejected", posted: "Posted",
};

export const STATUS_VARIANT: Record<string, "ok" | "warn" | "danger" | "info" | "muted"> = {
  active: "ok", running: "info", paid: "ok", approved: "ok", succeeded: "ok", posted: "ok",
  open: "info", provisioning: "info", claimed: "info",
  pending: "warn", manual_review: "warn", pending_verification: "warn", submitted: "warn", under_review: "warn",
  queued: "warn", failed_retryable: "warn",
  overdue: "danger", failed: "danger", rejected: "danger", failed_terminal: "danger",
  suspended: "muted", stopped: "muted", cancelled: "muted", unknown: "muted",
};

export function fmtMoney(v: number): string {
  const sign = v < 0 ? "-" : "";
  return sign + "$" + Math.abs(v).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

export function fmtMoneyShort(v: number): string {
  if (Math.abs(v) >= 1000) return "$" + (v / 1000).toFixed(1) + "k";
  return "$" + v.toFixed(0);
}
