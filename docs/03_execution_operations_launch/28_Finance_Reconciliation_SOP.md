# 28 — Finance Reconciliation SOP

**Version:** v1.3  
**Date:** 2026-04-22  
**Owner:** Finance/Admin Operations  
**Scope:** Wallet, ledger, reseller settlement, revenue reconciliation.

## 1. Purpose

This document defines daily, weekly, and monthly reconciliation so the platform never reaches the dangerous state where orders are flowing but nobody knows whether the money is real.

Principles:

```text
Ledger is the source of truth.
Do not edit old transactions.
Every correction is a new adjustment/reversal with reason and approver.
```

## 2. Financial entities

```text
Client wallet: internal wallet used by client to buy services at selling_price.
Reseller wallet: settlement wallet between reseller and platform; debited at reseller_cost.
Platform revenue: amount earned by platform from reseller_cost/direct sales.
Reseller gross profit: selling_price - reseller_cost.
Ledger entry: immutable row recording every balance change.
```

## 3. Ledger entry types

```text
client_topup_credit
client_purchase_debit
client_refund_credit
client_adjustment_credit
client_adjustment_debit
reseller_topup_credit
reseller_settlement_debit
reseller_refund_credit
reseller_adjustment_credit
reseller_adjustment_debit
platform_revenue_recognition
platform_refund_reversal
manual_correction
```

## 4. Core formulas

### Client wallet balance

```text
opening_balance
+ client_topup_credit
+ client_refund_credit
+ client_adjustment_credit
- client_purchase_debit
- client_adjustment_debit
```

### Reseller wallet balance

```text
opening_balance
+ reseller_topup_credit
+ reseller_refund_credit
+ reseller_adjustment_credit
- reseller_settlement_debit
- reseller_adjustment_debit
```

### Platform gross revenue

```text
sum(reseller_settlement_debit)
+ direct_platform_purchase_revenue
- refund_reversals
```

### Reseller gross profit

```text
sum(client_selling_price_snapshot - reseller_cost_snapshot)
for successful orders
minus reseller-side refunds/credits based on policy
```

### Invariant

```text
wallet.current_balance must equal sum(wallet_ledger_entries.amount)
```

If not, open finance incident. Ledger wins.

## 5. Daily reconciliation checklist

```text
1. Export wallet balances.
2. Recompute balances from ledger.
3. Compare materialized wallet balance vs ledger sum.
4. List mismatches.
5. Check orders paid but not provisioned/active/manual_review.
6. Check active services with no paid order.
7. Check expired reservations still counted as reserved.
8. Check negative reseller wallets.
9. Check top-up approved but no ledger credit.
10. Check top-up rejected but balance changed.
11. Check refunds with missing reversal/credit.
12. Check large manual adjustments.
13. Check provisioning failed after debit without refund/manual review.
```

## 6. Exception report

| Exception | Meaning | Action |
|---|---|---|
| Ledger mismatch | wallet balance != sum ledger | P0 finance incident |
| Paid not provisioned | order paid but no active/manual review | Ops review |
| Active without paid | service active without valid order | P0 revenue leak |
| Negative reseller balance | reseller wallet < 0 | Freeze auto-provision |
| Stuck top-up | pending too long | Finance follow-up |
| Large adjustment | exceeds threshold | Manager approval |
| Duplicate debit | same order charged twice | P0 refund/reversal review |

## 7. Weekly reconciliation

```text
- revenue by tenant
- revenue by product
- revenue by provider source
- reseller gross margin
- refunds and adjustments
- failed provisioning financial impact
- provider cost vs reseller cost snapshot
- manual top-up approval audit
- top risk clients/resellers
```

## 8. Monthly close

```text
1. Freeze reporting period.
2. Export ledger entries.
3. Export wallet balances.
4. Export order/service revenue.
5. Export refunds/adjustments.
6. Reconcile platform revenue.
7. Reconcile reseller profit.
8. Review negative balances.
9. Archive close report.
10. Lock old period from normal adjustment.
```

Old-period corrections must be posted in the current period as a new adjustment referencing the old period.

## 9. Top-up approval SOP

Checks before approval:

```text
- payment reference exists
- amount matches request
- currency/FX rule applied
- no duplicate reference
- correct tenant/wallet
- evidence is sufficient
```

Approve flow:

```text
1. Verify payment evidence.
2. Approve top-up request.
3. System creates ledger credit.
4. Wallet balance updates.
5. Audit wallet.topup.approved.
6. Notify user.
```

Reject reasons:

```text
payment_not_received
amount_mismatch
duplicate_reference
wrong_account
fraud_suspected
insufficient_evidence
```

## 10. Refund SOP

Refund decision must define:

```text
- refund_to_client amount
- refund_to_reseller amount if settlement was debited
- reason
- order/service reference
- approver
```

Default rule:

```text
If provider confirmed no resource was created: refund according to policy.
If provider state is uncertain: do not auto-refund; move to manual review.
```

## 11. Adjustment SOP

Allowed reasons:

```text
finance correction
compensation
manual settlement
migration correction
dispute resolution
```

Every adjustment needs:

```text
amount, direction, wallet, reason, reference, maker, approver if above threshold.
```

Suggested threshold:

```text
<= 10 USD: Finance Agent
10–100 USD: Finance Lead
> 100 USD: Super Admin
```

## 12. Revenue recognition

Recommended MVP rule:

```text
Top-up = wallet liability.
Successful service activation/provisioning = revenue recognition.
Refund = revenue reversal.
```

Do not recognize top-up as revenue.

## 13. Red flags

```text
- wallet negative without approved credit policy
- order paid without ledger debit
- ledger debit without order
- duplicate payment reference approved
- service active without reseller settlement
- adjustment without reason
- admin approving own large adjustment
```

## 14. Finance dashboard minimum

```text
- total platform revenue
- revenue by tenant
- revenue by provider
- wallet liability total
- pending top-up amount
- refund amount
- adjustment amount
- failed provisioning exposure
- negative wallet list
```

## 15. Closing principle

```text
Dashboard revenue is a feeling.
Reconciled ledger is the truth.
```
