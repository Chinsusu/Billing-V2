"use client";

import { useState } from "react";
import { CLIENT_LEDGER, type LedgerEntry } from "@/mocks/billingData";
import { Modal } from "@/components/ui/Modal";
import { useToast } from "@/lib/toast/ToastContext";
import { fmtMoney } from "@/mocks/sampleData";

const METHODS = ["VietQR", "Bank Transfer", "USDT", "PayPal"];

export function ClientWallet() {
  const { toast } = useToast();
  const [balance, setBalance] = useState(128.40);
  const [ledger, setLedger] = useState<LedgerEntry[]>(CLIENT_LEDGER);
  const [showTopup, setShowTopup] = useState(false);
  const [form, setForm] = useState({ amount: "", method: "VietQR" });
  const [errors, setErrors] = useState<{ amount?: string }>({});

  const validate = () => {
    const a = parseFloat(form.amount);
    if (!form.amount || isNaN(a) || a <= 0) {
      setErrors({ amount: "Enter a valid amount" });
      return false;
    }
    setErrors({});
    return true;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    const amount = parseFloat(form.amount);
    const newBalance = balance + amount;
    const entry: LedgerEntry = {
      ts: new Date().toISOString().slice(0, 16).replace("T", " "),
      type: "topup.credit.client",
      amount,
      ref: `${form.method} · pending`,
      balance: newBalance,
    };
    setBalance(newBalance);
    setLedger((prev) => [entry, ...prev]);
    toast(`Top-up request submitted — ${fmtMoney(amount)} via ${form.method}`, "success");
    setShowTopup(false);
    setForm({ amount: "", method: "VietQR" });
  };

  return (
    <div className="p-4 flex flex-col gap-4">
      <div className="bg-white border border-gray-200 rounded p-4 flex items-start justify-between">
        <div>
          <div className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Available balance</div>
          <div className="text-3xl font-medium tabular-nums text-gray-900">{fmtMoney(balance)}</div>
          <div className="text-[12px] text-gray-400 mt-1">Linh Tran · via ProxyVN</div>
        </div>
        <button
          onClick={() => setShowTopup(true)}
          className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm"
        >
          + Top up
        </button>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Transaction history</h3>
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

      <Modal open={showTopup} onClose={() => setShowTopup(false)} title="Top up wallet" description="Submit a top-up request to your reseller" width="sm">
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Amount (USD) <span className="text-red-500">*</span></label>
            <input type="number" min="1" step="0.01" value={form.amount}
              onChange={(e) => setForm((f) => ({ ...f, amount: e.target.value }))}
              placeholder="50.00"
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
          <div className="flex justify-end gap-2 pt-1 border-t border-gray-100">
            <button type="button" onClick={() => setShowTopup(false)} className="h-9 px-4 text-[13px] bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 cursor-pointer">Cancel</button>
            <button type="submit" className="h-9 px-4 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer">Submit</button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
