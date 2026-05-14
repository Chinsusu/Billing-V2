import type { ServiceInstance } from "./types";

export interface ClientServiceRenewalBody {
  wallet_id: string;
  from_status: string;
  reason?: string;
}

export interface ClientServiceRenewal {
  service: ServiceInstance;
  invoice: {
    id: string;
    display_id: number;
    status: string;
    total_minor: number;
    currency: string;
  };
  payment_transaction: {
    id: string;
    display_id: number;
    status: string;
  };
  ledger: {
    id: string;
    display_id: number;
    wallet_id: string;
    entry_type: string;
  };
  amount_minor: number;
  currency: string;
  renewed: boolean;
}
