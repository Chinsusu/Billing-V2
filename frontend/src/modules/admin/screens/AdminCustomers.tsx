"use client";

import { useState } from "react";
import { CUSTOMERS, type Customer } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { Modal } from "@/components/ui/Modal";
import { useToast } from "@/lib/toast/ToastContext";
import { fmtMoneyShort } from "@/mocks/sampleData";

const PLANS = ["Starter", "Pro", "Business", "Enterprise"];

interface NewCustomerForm {
  name: string;
  email: string;
  plan: string;
  country: string;
}

const EMPTY: NewCustomerForm = { name: "", email: "", plan: "Starter", country: "" };

export function AdminCustomers() {
  const { toast } = useToast();
  const [customers, setCustomers] = useState<Customer[]>(CUSTOMERS);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState<NewCustomerForm>(EMPTY);
  const [errors, setErrors] = useState<Partial<NewCustomerForm>>({});

  const validate = () => {
    const e: Partial<NewCustomerForm> = {};
    if (!form.name.trim()) e.name = "Name is required";
    if (!form.email.trim() || !form.email.includes("@")) e.email = "Valid email required";
    if (!form.country.trim()) e.country = "Country is required";
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleClose = () => { setShowForm(false); setForm(EMPTY); setErrors({}); };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    const newCustomer: Customer = {
      id: `C-${String(40218 + customers.length + 1)}`,
      name: form.name.trim(),
      email: form.email.trim(),
      plan: form.plan,
      services: 0,
      mrr: 0,
      status: "active",
      since: new Date().toISOString().slice(0, 10),
      country: form.country.trim().toUpperCase().slice(0, 2),
    };
    setCustomers((prev) => [newCustomer, ...prev]);
    toast(`Customer "${newCustomer.name}" created`, "success");
    handleClose();
  };

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Customers</h3>
          <div className="flex items-center gap-3">
            <span className="text-[11px] text-gray-400">{customers.length.toLocaleString()} total</span>
            <button
              onClick={() => setShowForm(true)}
              className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-white hover:bg-gray-50 text-gray-700 border border-gray-300 rounded-md cursor-pointer transition-colors shadow-sm"
            >
              + Add customer
            </button>
          </div>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["ID", "Name", "Email", "Plan", "Services", "MRR", "Status", "Since", "Country"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {customers.map((c) => (
              <tr key={c.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 text-[12px] text-[#D50C2D]">{c.id}</td>
                <td className="p-4 font-medium text-gray-900">{c.name}</td>
                <td className="p-4 text-gray-400 text-[12px]">{c.email}</td>
                <td className="p-4">
                  <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{c.plan}</span>
                </td>
                <td className="p-4 tabular-nums text-right">{c.services}</td>
                <td className="p-4 tabular-nums text-right font-medium">{fmtMoneyShort(c.mrr)}</td>
                <td className="p-4"><StatusBadge status={c.status} dot /></td>
                <td className="p-4 text-gray-400">{c.since}</td>
                <td className="p-4 text-gray-400">{c.country}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal open={showForm} onClose={handleClose} title="Add customer" description="Create a new customer account" width="sm">
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Full name <span className="text-red-500">*</span></label>
            <input value={form.name} onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))} placeholder="Acme Proxy Co."
              className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400" />
            {errors.name && <span className="text-[11px] text-red-500">{errors.name}</span>}
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Email <span className="text-red-500">*</span></label>
            <input type="email" value={form.email} onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))} placeholder="ops@example.com"
              className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400" />
            {errors.email && <span className="text-[11px] text-red-500">{errors.email}</span>}
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="flex flex-col gap-1">
              <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Plan</label>
              <select value={form.plan} onChange={(e) => setForm((f) => ({ ...f, plan: e.target.value }))}
                className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400 bg-white">
                {PLANS.map((p) => <option key={p}>{p}</option>)}
              </select>
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Country code <span className="text-red-500">*</span></label>
              <input value={form.country} onChange={(e) => setForm((f) => ({ ...f, country: e.target.value }))} placeholder="VN"
                maxLength={2}
                className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400 uppercase" />
              {errors.country && <span className="text-[11px] text-red-500">{errors.country}</span>}
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-1 border-t border-gray-100">
            <button type="button" onClick={handleClose} className="h-9 px-4 text-[13px] bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 cursor-pointer">Cancel</button>
            <button type="submit" className="h-9 px-4 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer">Create customer</button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
