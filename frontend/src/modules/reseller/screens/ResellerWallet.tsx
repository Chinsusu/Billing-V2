"use client";

import { useState } from "react";
import { Modal } from "@/components/ui/Modal";
import { useToast } from "@/lib/toast/ToastContext";
import { fmtMoney } from "@/mocks/sampleData";

interface LedgerRow { ts: string; type: string; amount: number; ref: string; balance: number; }

const INITIAL_LEDGER: LedgerRow[] = [
  { ts: "2026-04-22 14:02", type: "settlement.reseller.debit", amount: -62.00, ref: "ORD-48291 · VPS 4C/8G", balance: 4820.50 },
  { ts: "2026-04-21 10:18", type: "topup.credit.reseller", amount: 2000.00, ref: "TUP-9116 · VietQR", balance: 4882.50 },
  { ts: "2026-04-20 14:08", type: "settlement.reseller.debit", amount: -390.00, ref: "ORD-48280 · Residential batch", balance: 2882.50 },
  { ts: "2026-04-18 11:22", type: "settlement.reseller.debit", amount: -180.00, ref: "ORD-48270 · ISP batch", balance: 3272.50 },
];

const METHODS = ["VietQR", "Bank Transfer", "USDT", "PayPal"];

export function ResellerWallet() {
  const { toast } = useToast();
  const [balance, setBalance] = useState(4820.50);
  const [ledger, setLedger] = useState<LedgerRow[]>(INITIAL_LEDGER);
  const [showTopup, setShowTopup] = useState(false);
  const [form, setForm] = useState({ amount: "", method: "VietQR", reference: "" });
  const [errors, setErrors] = useState<{ amount?: string; reference?: string }>({});

  const validate = () => {
    const e: typeof errors = {};
    const a = parseFloat(form.amount);
    if (!form.amount || isNaN(a) || a <= 0) e.amount = "Enter a valid amount";
    if (!form.reference.trim()) e.reference = "Reference is required";
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    const amount = parseFloat(form.amount);
    const newBalance = balance + amount;
    const entry: LedgerRow = {
      ts: new Date().toISOString().slice(0, 16).replace("T", " "),
      type: "topup.credit.reseller",
      amount,
      ref: `${form.method} · ${form.reference}`,
      balance: newBalance,
    };
    setBalance(newBalance);
    setLedger((prev) => [entry, ...prev]);
    toast(`Top-up request submitted — ${fmtMoney(amount)} via ${form.method}`, "success");
    setShowTopup(false);
    setForm({ amount: "", method: "VietQR", reference: "" });
    setErrors({});
  };

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="bg-white border border-gray-200 rounded p-4">
        <div className="flex items-start justify-between">
          <div>
            <div className="text-[11px] font-medium uppercase tracking-wide text-gray-400 mb-1">Wallet balance</div>
            <div className="text-3xl font-medium tabular-nums text-gray-900">{fmtMoney(balance)}</div>
            <div className="text-[12px] text-gray-400 mt-1">ProxyVN · T-0042</div>
          </div>
          <button
            onClick={() => setShowTopup(true)}
            className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm"
          >
            + Request top-up
          </button>
        </div>
        <div className="mt-4 pt-4 border-t border-gray-100 grid grid-cols-3 gap-4">
          {[
            { label: "Pending top-ups", value: "$2,000.00", sub: "TUP-9120 · awaiting admin" },
            { label: "Spent this month", value: "$8,240.00", sub: "settlement debits" },
            { label: "Low balance alert", value: "< $200", sub: "notify when below" },
          ].map(({ label, value, sub }) => (
            <div key={label}>
              <div className="text-[11px] text-gray-400 mb-0.5">{label}</div>
              <div className="text-[14px] font-medium tabular-nums">{value}</div>
              <div className="text-[11px] text-gray-400">{sub}</div>
            </div>
          ))}
        </div>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Ledger history</h3>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["Timestamp", "Type", "Amount", "Reference", "Balance after"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {ledger.map((e, i) => (
              <tr key={i} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 tabular-nums text-gray-400">{e.ts}</td>
                <td className="p-4 text-[12px] text-gray-500">{e.type}</td>
                <td className={`p-4 tabular-nums text-right font-medium ${e.amount < 0 ? "text-red-600" : "text-green-700"}`}>{fmtMoney(e.amount)}</td>
                <td className="p-4 text-gray-500">{e.ref}</td>
                <td className="p-4 tabular-nums text-right font-medium">{fmtMoney(e.balance)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal open={showTopup} onClose={() => setShowTopup(false)} title="Request top-up" description="Submit a top-up request to the platform admin">
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Amount (USD) <span className="text-red-500">*</span></label>
            <input type="number" min="1" step="0.01" value={form.amount} onChange={(e) => setForm((f) => ({ ...f, amount: e.target.value }))} placeholder="500.00"
              className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400" />
            {errors.amount && <span className="text-[11px] text-red-500">{errors.amount}</span>}
          </div>
          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Payment method</label>
            <select value={form.method} onChange={(e) => setForm((f) => ({ ...f, method: e.target.value }))}
              className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400 bg-white">
              {METHODS.map((m) => <option key={m}>{m}</option>)}
            </select>
          </div>
          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Bank reference / TxID <span className="text-red-500">*</span></label>
            <input value={form.reference} onChange={(e) => setForm((f) => ({ ...f, reference: e.target.value }))} placeholder="FT26042200832"
              className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400" />
            {errors.reference && <span className="text-[11px] text-red-500">{errors.reference}</span>}
          </div>
          <div className="flex justify-end gap-2 pt-1 border-t border-gray-100">
            <button type="button" onClick={() => setShowTopup(false)} className="h-9 px-4 text-[13px] bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 cursor-pointer">Cancel</button>
            <button type="submit" className="h-9 px-4 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer">Submit request</button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
