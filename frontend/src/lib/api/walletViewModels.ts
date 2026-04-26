import { recordLabel } from "./format";
import type { LedgerEntry } from "./types";

type WalletLedgerReference = Pick<LedgerEntry, "display_id" | "reference_display_id" | "reference_type">;

const REFERENCE_PREFIXES: Record<string, string> = {
  invoice: "INV-",
  order: "ORD-",
  payment_transaction: "TX-",
  topup_request: "TUP-",
};

const PAYMENT_METHOD_LABELS: Record<string, string> = {
  bank: "Bank transfer",
  bank_transfer: "Bank transfer",
  card: "Card",
  crypto: "Crypto",
  manual: "Manual",
  other: "Other",
  paypal: "PayPal",
  usdt: "USDT",
  vietqr: "VietQR",
  wallet: "Wallet",
  wire: "Wire transfer",
  wire_transfer: "Wire transfer",
};

export function walletLedgerReferenceLabel(entry: WalletLedgerReference): string {
  if (entry.reference_display_id) {
    return recordLabel(entry.reference_display_id, REFERENCE_PREFIXES[entry.reference_type] ?? "#");
  }
  return recordLabel(entry.display_id, "LED-");
}

export function paymentMethodLabel(method: string): string {
  const normalized = method.trim();
  if (!normalized) return "-";
  const knownLabel = PAYMENT_METHOD_LABELS[normalized.toLowerCase()];
  if (knownLabel) return knownLabel;

  return normalized
    .split(/[_-]/)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

const ENTRY_TYPE_LABELS: Record<string, string> = {
  "purchase.client_wallet.debit": "Purchase",
  "renewal.client_wallet.debit": "Service renewal",
  "settlement.reseller.debit": "Reseller settlement",
  "topup.credit.client": "Top-up credit",
  "topup.credit.reseller": "Top-up credit",
};

const ENTRY_WORD_LABELS: Record<string, string> = {
  credit: "Credit",
  debit: "Debit",
  renewal: "Renewal",
  reseller: "Reseller",
  settlement: "Settlement",
  topup: "Top-up",
};

export function walletLedgerEntryTypeLabel(entryType: string): string {
  const knownLabel = ENTRY_TYPE_LABELS[entryType];
  if (knownLabel) return knownLabel;

  return entryType
    .split(/[._-]/)
    .filter(Boolean)
    .filter((part) => part !== "client" && part !== "wallet")
    .map((part) => ENTRY_WORD_LABELS[part] ?? part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}
