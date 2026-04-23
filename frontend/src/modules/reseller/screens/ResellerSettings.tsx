"use client";

import { useState } from "react";
import { useToast } from "@/lib/toast/ToastContext";

interface BrandingSettings { storeName: string; domain: string; supportEmail: string; accentColor: string; }

const DEFAULTS: BrandingSettings = {
  storeName: "ProxyVN",
  domain: "proxyvn.io",
  supportEmail: "support@proxyvn.io",
  accentColor: "#D50C2D",
};

export function ResellerSettings() {
  const { toast } = useToast();
  const [settings, setSettings] = useState<BrandingSettings>(DEFAULTS);
  const [saved, setSaved] = useState<BrandingSettings>(DEFAULTS);

  const isDirty = JSON.stringify(settings) !== JSON.stringify(saved);

  const handleSave = () => {
    setSaved(settings);
    toast("Branding settings saved", "success");
  };

  const fields: { key: keyof BrandingSettings; label: string }[] = [
    { key: "storeName", label: "Store name" },
    { key: "domain", label: "Custom domain" },
    { key: "supportEmail", label: "Support email" },
    { key: "accentColor", label: "Accent color" },
  ];

  return (
    <div className="p-4 flex flex-col gap-4 max-w-2xl">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Branding</h3>
        </div>
        <div className="p-4 space-y-4">
          {fields.map(({ key, label }) => (
            <div key={key} className="flex items-center gap-4">
              <label className="text-[12px] text-gray-500 w-36 shrink-0">{label}</label>
              <input
                value={settings[key]}
                onChange={(e) => setSettings((s) => ({ ...s, [key]: e.target.value }))}
                className="flex-1 h-8 px-2.5 border border-gray-300 rounded-[3px] text-[13px] font-[inherit] text-gray-800 bg-white outline-none focus:border-[#D50C2D]"
              />
            </div>
          ))}
        </div>
        <div className="p-4 border-t border-gray-100 flex items-center justify-between">
          {isDirty
            ? <span className="text-[11px] text-amber-600">Unsaved changes</span>
            : <span className="text-[11px] text-gray-400">All changes saved</span>
          }
          <button
            onClick={handleSave}
            disabled={!isDirty}
            className="inline-flex items-center justify-center gap-2 px-4 h-9 text-[13px] font-medium bg-[#D50C2D] hover:bg-[#B3082A] text-white rounded-md border-0 cursor-pointer transition-colors shadow-sm disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Save branding
          </button>
        </div>
      </div>
    </div>
  );
}
