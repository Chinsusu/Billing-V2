import { recordLabel } from "./format";
import { paymentTransactionTypeLabel } from "./paymentViewModels";
import type { LedgerEntry } from "./types";

type WalletLedgerReference = Pick<LedgerEntry, "display_id" | "reference_display_id" | "reference_type">;

const REFERENCE_PREFIXES: Record<string, string> = {
  invoice: "INV-",
  order: "ORD-",
  payment_transaction: "TX-",
  topup_request: "TUP-",
};

export function walletLedgerReferenceLabel(entry: WalletLedgerReference): string {
  if (entry.reference_display_id) {
    return recordLabel(entry.reference_display_id, REFERENCE_PREFIXES[entry.reference_type] ?? "#");
  }
  return recordLabel(entry.display_id, "LED-");
}

const ENTRY_TYPE_LABELS: Record<string, string> = {
  "purchase.client_wallet.debit": "Purchase",
  "renewal.client_wallet.debit": "Service renewal",
  "settlement.reseller.debit": "Reseller settlement",
  "topup.credit.client": "Top-up credit",
  "topup.credit.reseller": "Top-up credit",
};

export function walletLedgerEntryTypeLabel(entryType: string): string {
  const knownLabel = ENTRY_TYPE_LABELS[entryType];
  if (knownLabel) return knownLabel;

  return paymentTransactionTypeLabel(entryType);
}
