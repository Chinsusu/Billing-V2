import { recordLabel } from "./format";
import type { AdminAccount } from "./types";

export type ResellerAccountIdentity = Pick<AdminAccount, "display_id" | "full_name" | "email">;

export function resellerAccountLabel(displayID?: number, account?: ResellerAccountIdentity): string {
  if (!displayID) return "-";
  const publicID = recordLabel(displayID, "ACC-");
  if (!account) return publicID;
  return `${account.full_name || account.email} (${publicID})`;
}
