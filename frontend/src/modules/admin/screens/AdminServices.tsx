"use client";

import { useState } from "react";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { PROXY_SERVICES, VPS_SERVICES, BANDWIDTH_SERVICES } from "@/mocks/billingData";
import { fmtMoney } from "@/mocks/sampleData";

type Tab = "proxies" | "vps" | "bandwidth";

const TABS: { id: Tab; label: string; count: number }[] = [
  { id: "proxies",   label: "Proxies",   count: PROXY_SERVICES.length },
  { id: "vps",       label: "VPS",       count: VPS_SERVICES.length },
  { id: "bandwidth", label: "Bandwidth", count: BANDWIDTH_SERVICES.length },
];

const PROXY_TYPE_LABEL: Record<string, string> = {
  residential: "Residential",
  datacenter:  "Datacenter",
  mobile:      "Mobile",
  isp:         "ISP",
};

export function AdminServices() {
  const [tab, setTab] = useState<Tab>("proxies");

  return (
    <div className="p-4 flex flex-col gap-4">
      {/* Sub-tabs */}
      <div className="flex gap-0 border-b border-gray-200">
        {TABS.map((t) => (
          <button
            key={t.id}
            onClick={() => setTab(t.id)}
            className={`p-4 p-4 text-[13px] font-medium border-b-2 -mb-px transition-colors cursor-pointer bg-transparent
              ${tab === t.id
                ? "border-[#D50C2D] text-[#D50C2D]"
                : "border-transparent text-gray-500 hover:text-gray-800"}`}
          >
            {t.label}
            <span className={`ml-1.5 text-[11px] tabular-nums ${tab === t.id ? "text-[#D50C2D]" : "text-gray-400"}`}>
              {t.count}
            </span>
          </button>
        ))}
      </div>

      {/* Proxies */}
      {tab === "proxies" && (
        <div className="bg-white border border-gray-200 rounded">
          <table className="w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Type", "Label", "Customer", "Tenant", "Region", "IPs", "Protocol", "Usage", "Price/mo", "Status", "Renews"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {PROXY_SERVICES.map((s) => {
                const pct = Math.round((s.usedGB / s.totalGB) * 100);
                return (
                  <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                    <td className="p-4 p-4 text-[11px] text-[#D50C2D]">{s.id}</td>
                    <td className="p-4 p-4">
                      <span className="text-[10px] font-medium px-1.5 py-0.5 rounded bg-indigo-50 text-indigo-700">
                        {PROXY_TYPE_LABEL[s.proxyType]}
                      </span>
                    </td>
                    <td className="p-4 p-4 text-gray-800 max-w-[180px] truncate">{s.label}</td>
                    <td className="p-4 p-4 text-gray-500">{s.customer}</td>
                    <td className="p-4 p-4 text-gray-400 text-[11px]">{s.tenant}</td>
                    <td className="p-4 p-4 text-[11px] text-gray-400">{s.region}</td>
                    <td className="p-4 p-4 text-gray-500 tabular-nums">{s.ipCount > 0 ? s.ipCount : "—"}</td>
                    <td className="p-4 p-4">
                      <span className="text-[10px] px-1 py-px bg-gray-100 text-gray-500 rounded">{s.protocol}</span>
                    </td>
                    <td className="p-4 p-4">
                      <div className="flex items-center gap-4 min-w-[100px]">
                        <div className="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden">
                          <div
                            className={`h-full rounded-full ${pct >= 90 ? "bg-red-500" : pct >= 70 ? "bg-amber-400" : "bg-green-500"}`}
                            style={{ width: `${pct}%` }}
                          />
                        </div>
                        <span className="text-[11px] text-gray-400 tabular-nums w-8 text-right">{pct}%</span>
                      </div>
                    </td>
                    <td className="p-4 p-4 tabular-nums text-right font-medium">{fmtMoney(s.price)}</td>
                    <td className="p-4 p-4"><StatusBadge status={s.status} dot /></td>
                    <td className="p-4 p-4 tabular-nums">
                      <span className={s.renewsIn < 0 ? "text-red-600 font-medium" : s.renewsIn <= 7 ? "text-amber-600" : "text-gray-500"}>
                        {s.renewsIn < 0 ? `${Math.abs(s.renewsIn)}d overdue` : `${s.renewsIn}d`}
                      </span>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {/* VPS */}
      {tab === "vps" && (
        <div className="bg-white border border-gray-200 rounded">
          <table className="w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "OS", "Label", "Customer", "Tenant", "Region", "Spec", "IP", "Provider", "Price/mo", "Status", "Renews"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {VPS_SERVICES.map((s) => (
                <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 p-4 text-[11px] text-[#D50C2D]">{s.id}</td>
                  <td className="p-4 p-4">
                    <span className={`text-[10px] font-medium px-1.5 py-0.5 rounded ${s.os === "linux" ? "bg-orange-50 text-orange-700" : "bg-blue-50 text-blue-700"}`}>
                      {s.os === "linux" ? "Linux" : "Windows"}
                    </span>
                  </td>
                  <td className="p-4 p-4 text-gray-800 max-w-[160px] truncate">{s.label}</td>
                  <td className="p-4 p-4 text-gray-500">{s.customer}</td>
                  <td className="p-4 p-4 text-gray-400 text-[11px]">{s.tenant}</td>
                  <td className="p-4 p-4 text-[11px] text-gray-400">{s.region}</td>
                  <td className="p-4 p-4 text-gray-500 text-[11px] whitespace-nowrap">
                    {s.cpu}C / {s.ram}G / {s.disk}G
                  </td>
                  <td className="p-4 p-4 text-[11px] text-gray-400">{s.ip}</td>
                  <td className="p-4 p-4 text-gray-400 text-[11px]">{s.provider}</td>
                  <td className="p-4 p-4 tabular-nums text-right font-medium">{fmtMoney(s.price)}</td>
                  <td className="p-4 p-4"><StatusBadge status={s.status} dot /></td>
                  <td className="p-4 p-4 tabular-nums">
                    <span className={s.renewsIn < 0 ? "text-red-600 font-medium" : s.renewsIn <= 7 ? "text-amber-600" : "text-gray-500"}>
                      {s.renewsIn < 0 ? `${Math.abs(s.renewsIn)}d overdue` : `${s.renewsIn}d`}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Bandwidth */}
      {tab === "bandwidth" && (
        <div className="bg-white border border-gray-200 rounded">
          <table className="w-full text-[13px] border-collapse">
            <thead>
              <tr className="bg-gray-50">
                {["ID", "Label", "Customer", "Tenant", "Region", "Used", "Total", "Usage %", "Price/mo", "Status", "Renews"].map((h) => (
                  <th key={h} className="text-left text-[11px] font-medium uppercase tracking-wide text-gray-400 p-4 p-4 border-b border-gray-200">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {BANDWIDTH_SERVICES.map((s) => (
                <tr key={s.id} className="hover:bg-gray-50 border-b border-gray-100 last:border-0">
                  <td className="p-4 p-4 text-[11px] text-[#D50C2D]">{s.id}</td>
                  <td className="p-4 p-4 text-gray-800 max-w-[200px] truncate">{s.label}</td>
                  <td className="p-4 p-4 text-gray-500">{s.customer}</td>
                  <td className="p-4 p-4 text-gray-400 text-[11px]">{s.tenant}</td>
                  <td className="p-4 p-4 text-[11px] text-gray-400">{s.region}</td>
                  <td className="p-4 p-4 tabular-nums text-gray-600">{s.usedGB} GB</td>
                  <td className="p-4 p-4 tabular-nums text-gray-400">{s.totalGB} GB</td>
                  <td className="p-4 p-4">
                    <div className="flex items-center gap-4 min-w-[110px]">
                      <div className="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden">
                        <div
                          className={`h-full rounded-full ${s.usedPct >= 90 ? "bg-red-500" : s.usedPct >= 70 ? "bg-amber-400" : "bg-green-500"}`}
                          style={{ width: `${s.usedPct}%` }}
                        />
                      </div>
                      <span className={`text-[11px] tabular-nums w-8 text-right font-medium ${s.usedPct >= 90 ? "text-red-600" : s.usedPct >= 70 ? "text-amber-600" : "text-gray-500"}`}>
                        {s.usedPct}%
                      </span>
                    </div>
                  </td>
                  <td className="p-4 p-4 tabular-nums text-right font-medium">{fmtMoney(s.price)}</td>
                  <td className="p-4 p-4"><StatusBadge status={s.status} dot /></td>
                  <td className="p-4 p-4 tabular-nums">
                    <span className={s.renewsIn < 0 ? "text-red-600 font-medium" : s.renewsIn <= 7 ? "text-amber-600" : "text-gray-500"}>
                      {s.renewsIn < 0 ? `${Math.abs(s.renewsIn)}d overdue` : `${s.renewsIn}d`}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
