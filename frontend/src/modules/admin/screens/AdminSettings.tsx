"use client";

import { useState } from "react";
import { useToast } from "@/lib/toast/ToastContext";

interface PlatformSettings {
  platformName: string;
  domain: string;
  supportEmail: string;
  currency: string;
  timezone: string;
}

const DEFAULTS: PlatformSettings = {
  platformName: "HANetwork Billing",
  domain: "billing.hanetwork.vn",
  supportEmail: "support@hanetwork.vn",
  currency: "USD",
  timezone: "Asia/Ho_Chi_Minh",
};

export function AdminSettings() {
  const { toast } = useToast();
  const [settings, setSettings] = useState<PlatformSettings>(DEFAULTS);
  const [saved, setSaved] = useState<PlatformSettings>(DEFAULTS);

  const isDirty = JSON.stringify(settings) !== JSON.stringify(saved);

  const handleSave = () => {
    setSaved(settings);
    toast("Platform settings saved", "success");
  };

  const fields: { key: keyof PlatformSettings; label: string }[] = [
    { key: "platformName", label: "Platform name" },
    { key: "domain", label: "Primary domain" },
    { key: "supportEmail", label: "Support email" },
    { key: "currency", label: "Default currency" },
    { key: "timezone", label: "Default timezone" },
  ];

  return (
    <div className="p-4 flex flex-col gap-4 max-w-2xl">
      <div className="bg-white border border-gray-200 rounded">
        <div className="p-4 border-b border-gray-100">
          <h3 className="text-[13px] font-medium text-gray-900 m-0">Platform settings</h3>
        </div>
        <div className="p-4 space-y-4">
          {fields.map(({ key, label }) => (
            <div key={key} className="flex items-center gap-4">
              <label className="text-[12px] text-gray-500 w-40 shrink-0">{label}</label>
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
            Save changes
          </button>
        </div>
      </div>

      <div className="bg-amber-50 border border-amber-200 rounded p-4 text-[12px] text-amber-700">
        Settings changes are audited. All modifications are logged with actor, timestamp, and previous value.
      </div>
    </div>
  );
}
