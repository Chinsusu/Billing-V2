"use client";

import { useState } from "react";
import { RESELLER_CLIENTS, type ResellerClient } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { Modal } from "@/components/ui/Modal";
import { useToast } from "@/lib/toast/ToastContext";
import { fmtMoney } from "@/mocks/sampleData";

interface NewClientForm { name: string; email: string; wallet: string; }
const EMPTY: NewClientForm = { name: "", email: "", wallet: "0" };

export function ResellerClients() {
  const { toast } = useToast();
  const [clients, setClients] = useState<ResellerClient[]>(RESELLER_CLIENTS);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState<NewClientForm>(EMPTY);
  const [errors, setErrors] = useState<Partial<NewClientForm>>({});

  const validate = () => {
    const e: Partial<NewClientForm> = {};
    if (!form.name.trim()) e.name = "Name is required";
    if (!form.email.trim() || !form.email.includes("@")) e.email = "Valid email required";
    if (isNaN(parseFloat(form.wallet)) || parseFloat(form.wallet) < 0) e.wallet = "Must be 0 or more";
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    const newClient: ResellerClient = {
      id: `RC-${String(clients.length + 3000).slice(-4)}`,
      name: form.name.trim(),
      email: form.email.trim(),
      wallet: parseFloat(form.wallet),
      services: 0,
      orders: 0,
      status: "active",
      lastLogin: "just now",
    };
    setClients((prev) => [...prev, newClient]);
    toast(`Client "${newClient.name}" created`, "success");
    setShowForm(false);
    setForm(EMPTY);
    setErrors({});
  };

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Clients</h3>
          <button
            onClick={() => setShowForm(true)}
            className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm"
          >
            + Add client
          </button>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Name", "Email", "Wallet", "Services", "Orders", "Status", "Last login"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {clients.map((c) => (
              <tr key={c.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 text-[12px] text-[#D50C2D]">{c.id}</td>
                <td className="p-4 font-medium text-gray-900">{c.name}</td>
                <td className="p-4 text-gray-400 text-[12px]">{c.email}</td>
                <td className="p-4 tabular-nums">
                  <span className={c.wallet < 20 ? "text-red-600 font-medium" : ""}>{fmtMoney(c.wallet)}</span>
                </td>
                <td className="p-4 tabular-nums text-right">{c.services}</td>
                <td className="p-4 tabular-nums text-right">{c.orders}</td>
                <td className="p-4"><StatusBadge status={c.status} dot /></td>
                <td className="p-4 text-gray-400">{c.lastLogin}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal open={showForm} onClose={() => { setShowForm(false); setForm(EMPTY); setErrors({}); }} title="Add client" description="Create a new client under your reseller account">
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          {([
            { key: "name" as const, label: "Full name", placeholder: "Linh Tran", type: "text" },
            { key: "email" as const, label: "Email", placeholder: "linh@example.com", type: "email" },
            { key: "wallet" as const, label: "Initial wallet (USD)", placeholder: "0", type: "number" },
          ]).map(({ key, label, placeholder, type }) => (
            <div key={key} className="flex flex-col gap-1">
              <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">
                {label} <span className="text-red-500">*</span>
              </label>
              <input
                type={type}
                min={type === "number" ? "0" : undefined}
                value={form[key]}
                onChange={(e) => setForm((f) => ({ ...f, [key]: e.target.value }))}
                placeholder={placeholder}
                className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400"
              />
              {errors[key] && <span className="text-[11px] text-red-500">{errors[key]}</span>}
            </div>
          ))}
          <div className="flex justify-end gap-2 pt-1 border-t border-gray-100">
            <button type="button" onClick={() => { setShowForm(false); setForm(EMPTY); setErrors({}); }} className="h-9 px-4 text-[13px] bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 cursor-pointer">Cancel</button>
            <button type="submit" className="h-9 px-4 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer">Create client</button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
