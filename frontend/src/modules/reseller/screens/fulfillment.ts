import { recordLabel } from "@/lib/api/format";
import type { Order, ServiceInstance } from "@/lib/api/types";

export interface FulfillmentState {
  status: string;
  label: string;
  orderLabel: string;
  serviceLabel: string;
  jobLabel: string;
}

export function fulfillmentForOrder(order: Order | undefined, services: ServiceInstance[]): FulfillmentState {
  if (!order) {
    return {
      status: "pending",
      label: "Order not loaded",
      orderLabel: "-",
      serviceLabel: "-",
      jobLabel: "-",
    };
  }

  const linkedServices = services.filter((service) => service.order_id === order.id);
  const activeService = linkedServices.find((service) => service.status === "active");
  if (activeService) {
    return {
      status: "active",
      label: "Active service",
      orderLabel: recordLabel(order.display_id, "ORD-"),
      serviceLabel: recordLabel(activeService.display_id, "SVC-"),
      jobLabel: "-",
    };
  }

  const visibleService = linkedServices[0];
  if (visibleService) {
    return {
      status: visibleService.status,
      label: visibleService.status,
      orderLabel: recordLabel(order.display_id, "ORD-"),
      serviceLabel: recordLabel(visibleService.display_id, "SVC-"),
      jobLabel: "-",
    };
  }

  if (order.order_status === "failed") {
    return baseOrderState(order, "failed", "Failed");
  }
  if (order.order_status === "cancelled" || order.order_status === "refunded") {
    return baseOrderState(order, order.order_status, order.order_status);
  }
  if (order.order_status === "paid" && order.billing_status === "paid") {
    return {
      status: "queued",
      label: "Paid / pending provisioning",
      orderLabel: recordLabel(order.display_id, "ORD-"),
      serviceLabel: "-",
      jobLabel: "provider.provision",
    };
  }
  return baseOrderState(order, "pending", "Pending payment");
}

export function fulfillmentForService(service: ServiceInstance, order: Order | undefined): FulfillmentState {
  const state = order ? fulfillmentForOrder(order, [service]) : fulfillmentForOrder(undefined, []);
  return {
    ...state,
    status: service.status || state.status,
    label: service.status || state.label,
    serviceLabel: recordLabel(service.display_id, "SVC-"),
  };
}

function baseOrderState(order: Order, status: string, label: string): FulfillmentState {
  return {
    status,
    label,
    orderLabel: recordLabel(order.display_id, "ORD-"),
    serviceLabel: "-",
    jobLabel: "-",
  };
}
