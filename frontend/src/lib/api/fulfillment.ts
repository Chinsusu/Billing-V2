import { recordLabel } from "./format";
import type { Order, ProvisioningJob, ServiceInstance } from "./types";

export interface FulfillmentState {
  status: string;
  label: string;
  orderLabel: string;
  serviceLabel: string;
  jobLabel: string;
  job?: ProvisioningJob;
}

export interface FulfillmentOptions {
  jobs?: ProvisioningJob[];
  jobsUnavailable?: boolean;
}

export function fulfillmentForOrder(
  order: Order | undefined,
  services: ServiceInstance[],
  options: FulfillmentOptions = {},
): FulfillmentState {
  if (!order) {
    return {
      status: "pending",
      label: "Order not loaded",
      orderLabel: "-",
      serviceLabel: "-",
      jobLabel: "-",
    };
  }

  const job = jobForOrder(order, options.jobs ?? []);
  const jobLabelText = job ? recordLabel(job.display_id, "JOB-") : options.jobsUnavailable ? "Job unavailable" : "-";
  const linkedServices = services.filter((service) => service.order_id === order.id);
  const activeService = linkedServices.find((service) => service.status === "active");
  if (activeService) {
    return {
      status: "active",
      label: "Active service",
      orderLabel: recordLabel(order.display_id, "ORD-"),
      serviceLabel: recordLabel(activeService.display_id, "SVC-"),
      jobLabel: jobLabelText,
      job,
    };
  }

  const visibleService = linkedServices[0];
  if (visibleService) {
    return {
      status: visibleService.status,
      label: visibleService.status,
      orderLabel: recordLabel(order.display_id, "ORD-"),
      serviceLabel: recordLabel(visibleService.display_id, "SVC-"),
      jobLabel: jobLabelText,
      job,
    };
  }

  if (order.order_status === "failed") {
    return baseOrderState(order, "failed", "Failed", jobLabelText, job);
  }
  if (order.order_status === "cancelled" || order.order_status === "refunded") {
    return baseOrderState(order, order.order_status, order.order_status, jobLabelText, job);
  }
  if (job) {
    return baseOrderState(order, job.status, jobStatusLabel(job.status), jobLabelText, job);
  }
  if (order.order_status === "paid" && order.billing_status === "paid") {
    return baseOrderState(
      order,
      options.jobsUnavailable ? "unknown" : "queued",
      options.jobsUnavailable ? "Job unavailable" : "Paid / no job found",
      jobLabelText,
      job,
    );
  }
  return baseOrderState(order, "pending", "Pending payment", jobLabelText, job);
}

export function fulfillmentForService(
  service: ServiceInstance,
  order: Order | undefined,
  options: FulfillmentOptions = {},
): FulfillmentState {
  const state = order ? fulfillmentForOrder(order, [service], options) : fulfillmentForOrder(undefined, [], options);
  return {
    ...state,
    status: service.status || state.status,
    label: service.status || state.label,
    serviceLabel: recordLabel(service.display_id, "SVC-"),
  };
}

export function jobForOrder(order: Order, jobs: ProvisioningJob[]): ProvisioningJob | undefined {
  return jobs.find((job) =>
    job.job_type === "provider.provision" &&
    job.reference_type === "order" &&
    job.reference_id === order.id
  );
}

export function jobStatusLabel(status: string): string {
  switch (status) {
    case "queued":
      return "Queued";
    case "claimed":
    case "running":
      return "Provisioning";
    case "succeeded":
      return "Provisioned";
    case "failed_retryable":
      return "Retryable failure";
    case "failed_terminal":
      return "Terminal failure";
    case "manual_review":
      return "Manual review";
    case "cancelled":
      return "Cancelled";
    default:
      return status || "-";
  }
}

export function canRetryJob(status: string): boolean {
  return status === "failed_retryable" || status === "manual_review";
}

export function canMarkJobManualReview(status: string): boolean {
  return status === "queued" || status === "failed_retryable" || status === "failed_terminal" || status === "manual_review";
}

export function canCancelJob(status: string): boolean {
  return status === "queued" || status === "failed_retryable" || status === "failed_terminal" || status === "manual_review";
}

function baseOrderState(
  order: Order,
  status: string,
  label: string,
  jobLabel: string,
  job?: ProvisioningJob,
): FulfillmentState {
  return {
    status,
    label,
    orderLabel: recordLabel(order.display_id, "ORD-"),
    serviceLabel: "-",
    jobLabel,
    job,
  };
}
