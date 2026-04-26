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

const PAYMENT_TRANSACTION_TYPE_LABELS: Record<string, string> = {
  adjustment: "Adjustment",
  charge: "Charge",
  refund: "Refund",
  topup: "Top-up",
};

const PAYMENT_WORD_LABELS: Record<string, string> = {
  adjustment: "Adjustment",
  charge: "Charge",
  credit: "Credit",
  debit: "Debit",
  payment: "Payment",
  purchase: "Purchase",
  refund: "Refund",
  renewal: "Service renewal",
  reseller: "Reseller",
  settlement: "Settlement",
  topup: "Top-up",
  transaction: "Transaction",
};

export function paymentMethodLabel(method: string): string {
  return labelValue(method, PAYMENT_METHOD_LABELS);
}

export function paymentTransactionTypeLabel(type: string): string {
  return labelValue(type, PAYMENT_TRANSACTION_TYPE_LABELS);
}

function labelValue(value: string, knownLabels: Record<string, string>): string {
  const normalized = value.trim();
  if (!normalized) return "-";
  const knownLabel = knownLabels[normalized.toLowerCase()];
  if (knownLabel) return knownLabel;

  return normalized
    .split(/[._-]/)
    .filter(Boolean)
    .filter((part) => part !== "client" && part !== "wallet")
    .map((part) => PAYMENT_WORD_LABELS[part.toLowerCase()] ?? titleCase(part))
    .join(" ");
}

function titleCase(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1);
}
