"use client";

import { useState } from "react";

const BILLING_CONTACTS = [
  { label: "Primary email", value: "linh.tran@gmail.com" },
  { label: "Invoice currency", value: "USD" },
  { label: "Tenant", value: "ProxyVN" },
];

const NOTIFICATION_OPTIONS = [
  { label: "Invoice issued", checked: true },
  { label: "Payment posted", checked: true },
  { label: "Service renewal", checked: true },
  { label: "Service suspended", checked: true },
];

export function ClientSettings() {
  const [options, setOptions] = useState(NOTIFICATION_OPTIONS);
  const [saved, setSaved] = useState(false);

  function toggleOption(label: string) {
    setSaved(false);
    setOptions((current) =>
      current.map((item) => item.label === label ? { ...item, checked: !item.checked } : item),
    );
  }

  return (
    <div className="p-4 grid grid-cols-1 xl:grid-cols-[minmax(0,1fr)_360px] gap-4">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Billing profile</h3>
        </div>
        <div className="divide-y divide-gray-100">
          {BILLING_CONTACTS.map((item) => (
            <div key={item.label} className="p-4 flex items-center justify-between gap-4">
              <span className="text-[12px] text-gray-500">{item.label}</span>
              <span className="text-[13px] font-medium text-gray-900 text-right">{item.value}</span>
            </div>
          ))}
        </div>
      </div>

      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Notifications</h3>
        </div>
        <div className="p-4 flex flex-col gap-3">
          {options.map((item) => (
            <label key={item.label} className="flex items-center justify-between gap-3 text-[13px] text-gray-700">
              <span>{item.label}</span>
              <input
                type="checkbox"
                checked={item.checked}
                onChange={() => toggleOption(item.label)}
                className="h-4 w-4 accent-[#D50C2D]"
              />
            </label>
          ))}
          <button
            onClick={() => setSaved(true)}
            className="mt-2 inline-flex items-center justify-center px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors"
          >
            {saved ? "Saved" : "Save settings"}
          </button>
        </div>
      </div>
    </div>
  );
}
