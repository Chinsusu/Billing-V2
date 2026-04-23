"use client";

import { useState } from "react";
import { TENANTS, type Tenant } from "@/mocks/billingData";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { Modal } from "@/components/ui/Modal";
import { useToast } from "@/lib/toast/ToastContext";
import { fmtMoney } from "@/mocks/sampleData";

interface NewTenantForm {
  name: string;
  type: "reseller" | "admin";
  domain: string;
}

const EMPTY_FORM: NewTenantForm = { name: "", type: "reseller", domain: "" };

export function AdminTenants() {
  const { toast } = useToast();
  const [tenants, setTenants] = useState<Tenant[]>(TENANTS);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState<NewTenantForm>(EMPTY_FORM);
  const [errors, setErrors] = useState<Partial<NewTenantForm>>({});

  const validate = () => {
    const e: Partial<NewTenantForm> = {};
    if (!form.name.trim()) e.name = "Name is required";
    if (!form.domain.trim()) e.domain = "Domain is required";
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    const newTenant: Tenant = {
      id: `T-${String(tenants.length + 1).padStart(4, "0")}`,
      name: form.name.trim(),
      type: form.type,
      domain: form.domain.trim(),
      clients: 0,
      services: 0,
      wallet: 0,
      walletLow: false,
      status: "active",
      since: new Date().toISOString().slice(0, 10),
    };
    setTenants((prev) => [...prev, newTenant]);
    toast(`Account "${newTenant.name}" created successfully`, "success");
    setShowForm(false);
    setForm(EMPTY_FORM);
    setErrors({});
  };

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Accounts</h3>
          <button
            onClick={() => setShowForm(true)}
            className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-white hover:bg-gray-50 text-gray-700 border border-gray-300 rounded-md cursor-pointer transition-colors shadow-sm"
          >
            + New account
          </button>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {[
                { label: "ID", align: "left" },
                { label: "Name", align: "left" },
                { label: "Type", align: "left" },
                { label: "Domain", align: "left" },
                { label: "Clients", align: "right" },
                { label: "Services", align: "right" },
                { label: "Wallet", align: "right" },
                { label: "Status", align: "left" },
                { label: "Since", align: "left" },
              ].map((h) => (
                <th key={h.label} className={`${h.align === "right" ? "text-right" : "text-left"} text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200`}>
                  {h.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {tenants.map((t) => (
              <tr key={t.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 text-[12px] text-[#D50C2D]">{t.id}</td>
                <td className="p-4 font-medium text-gray-900">{t.name}</td>
                <td className="p-4">
                  <span className="text-[11px] px-1.5 py-px bg-gray-100 text-gray-500 rounded-sm">{t.type}</span>
                </td>
                <td className="p-4 text-[12px] text-gray-500">{t.domain}</td>
                <td className="p-4 tabular-nums text-right">{t.clients.toLocaleString()}</td>
                <td className="p-4 tabular-nums text-right">{t.services.toLocaleString()}</td>
                <td className="p-4 tabular-nums text-right">
                  {t.type === "admin" ? "—" : (
                    <span className={t.walletLow ? "text-red-600 font-medium" : ""}>{fmtMoney(t.wallet)}</span>
                  )}
                </td>
                <td className="p-4"><StatusBadge status={t.status} dot /></td>
                <td className="p-4 text-gray-400">{t.since ?? "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal open={showForm} onClose={() => { setShowForm(false); setForm(EMPTY_FORM); setErrors({}); }} title="New account" description="Create a new reseller or admin account">
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Name <span className="text-red-500">*</span></label>
            <input
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              placeholder="ProxyVN"
              className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400"
            />
            {errors.name && <span className="text-[11px] text-red-500">{errors.name}</span>}
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Type</label>
            <select
              value={form.type}
              onChange={(e) => setForm((f) => ({ ...f, type: e.target.value as "reseller" | "admin" }))}
              className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400 bg-white"
            >
              <option value="reseller">Reseller</option>
              <option value="admin">Admin</option>
            </select>
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">Domain <span className="text-red-500">*</span></label>
            <input
              value={form.domain}
              onChange={(e) => setForm((f) => ({ ...f, domain: e.target.value }))}
              placeholder="proxyvn.io"
              className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400"
            />
            {errors.domain && <span className="text-[11px] text-red-500">{errors.domain}</span>}
          </div>

          <div className="flex justify-end gap-2 pt-1 border-t border-gray-100">
            <button type="button" onClick={() => { setShowForm(false); setForm(EMPTY_FORM); setErrors({}); }} className="h-9 px-4 text-[13px] font-medium bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 cursor-pointer transition-colors">
              Cancel
            </button>
            <button type="submit" className="h-9 px-4 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors">
              Create account
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
