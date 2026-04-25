import { recordLabel } from "./format";
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
