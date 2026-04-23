"use client";

import { useState } from "react";
import { PRODUCTS, type ProductCatalog } from "@/mocks/billingData";
import { Modal } from "@/components/ui/Modal";
import { useToast } from "@/lib/toast/ToastContext";
import { fmtMoney, fmtMoneyShort } from "@/mocks/sampleData";

interface NewProductForm {
  sku: string;
  name: string;
  unit: string;
  price: string;
}

const EMPTY: NewProductForm = { sku: "", name: "", unit: "", price: "" };

function FormField({ label, placeholder, value, onChange, error, type = "text" }: {
  label: string; placeholder: string; value: string;
  onChange: (v: string) => void; error?: string; type?: string;
}) {
  return (
    <div className="flex flex-col gap-1">
      <label className="text-[11px] font-medium uppercase tracking-wide text-gray-500">
        {label} <span className="text-red-500">*</span>
      </label>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className="h-9 px-3 text-[13px] border border-gray-200 rounded focus:outline-none focus:ring-1 focus:ring-gray-400"
      />
      {error && <span className="text-[11px] text-red-500">{error}</span>}
    </div>
  );
}

export function AdminProducts() {
  const { toast } = useToast();
  const [products, setProducts] = useState<ProductCatalog[]>(PRODUCTS);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState<NewProductForm>(EMPTY);
  const [errors, setErrors] = useState<Partial<NewProductForm>>({});

  const validate = () => {
    const e: Partial<NewProductForm> = {};
    if (!form.sku.trim()) e.sku = "SKU is required";
    if (!form.name.trim()) e.name = "Name is required";
    if (!form.unit.trim()) e.unit = "Unit is required";
    const p = parseFloat(form.price);
    if (!form.price || isNaN(p) || p < 0) e.price = "Valid price required";
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleClose = () => { setShowForm(false); setForm(EMPTY); setErrors({}); };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    const newProduct: ProductCatalog = {
      sku: form.sku.trim().toUpperCase(),
      name: form.name.trim(),
      unit: form.unit.trim(),
      price: parseFloat(form.price),
      active: 0,
      rev30: 0,
    };
    setProducts((prev) => [...prev, newProduct]);
    toast(`Product "${newProduct.name}" added`, "success");
    handleClose();
  };

  return (
    <div className="p-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Products & Pricing</h3>
          <button
            onClick={() => setShowForm(true)}
            className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-white hover:bg-gray-50 text-gray-700 border border-gray-300 rounded-md cursor-pointer transition-colors shadow-sm"
          >
            + Add product
          </button>
        </div>
        <table className="w-full text-[13px] border-collapse">
          <thead>
            <tr className="bg-gray-50">
              {["SKU", "Name", "Unit", "Price", "Active", "Rev 30d"].map((h) => (
                <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 border-b border-gray-200">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {products.map((p) => (
              <tr key={p.sku} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                <td className="p-4 text-[12px] text-gray-500">{p.sku}</td>
                <td className="p-4 font-medium text-gray-900">{p.name}</td>
                <td className="p-4 text-gray-400 text-[12px]">{p.unit}</td>
                <td className="p-4 tabular-nums font-medium">{fmtMoney(p.price)}</td>
                <td className="p-4 tabular-nums text-right">{p.active.toLocaleString()}</td>
                <td className="p-4 tabular-nums text-right font-medium">{fmtMoneyShort(p.rev30)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal open={showForm} onClose={handleClose} title="Add product" description="Add a new product or pricing plan">
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="grid grid-cols-2 gap-3">
            <FormField label="SKU" placeholder="PRX-RES-STD" value={form.sku} onChange={(v) => setForm((f) => ({ ...f, sku: v }))} error={errors.sku} />
            <FormField label="Price (USD)" placeholder="6.50" value={form.price} onChange={(v) => setForm((f) => ({ ...f, price: v }))} error={errors.price} type="number" />
          </div>
          <FormField label="Name" placeholder="Residential · Standard" value={form.name} onChange={(v) => setForm((f) => ({ ...f, name: v }))} error={errors.name} />
          <FormField label="Unit" placeholder="per GB" value={form.unit} onChange={(v) => setForm((f) => ({ ...f, unit: v }))} error={errors.unit} />

          <div className="flex justify-end gap-2 pt-1 border-t border-gray-100">
            <button type="button" onClick={handleClose} className="h-9 px-4 text-[13px] font-medium bg-white border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 cursor-pointer transition-colors">
              Cancel
            </button>
            <button type="submit" className="h-9 px-4 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors">
              Add product
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
