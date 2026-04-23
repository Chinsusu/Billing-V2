"use client";

import { useState } from "react";
import { AdminPortal } from "@/modules/admin/AdminPortal";
import { ResellerPortal } from "@/modules/reseller/ResellerPortal";
import { ClientPortal } from "@/modules/client/ClientPortal";

type Portal = "admin" | "reseller" | "client";

const PORTALS: { id: Portal; label: string }[] = [
  { id: "admin", label: "Admin · HANetwork" },
  { id: "reseller", label: "Reseller · ProxyVN" },
  { id: "client", label: "Client · Linh Tran" },
];

export default function Home() {
  const [portal, setPortal] = useState<Portal>("admin");

  return (
    <div className="flex flex-col h-full">
      {/* Portal switcher */}
      <div className="shrink-0 h-8 bg-[#0E1116] flex items-center gap-0.5 p-4 z-50">
        <span className="text-[11px] text-gray-500 mr-2">Portal</span>
        {PORTALS.map((p) => (
          <button
            key={p.id}
            onClick={() => setPortal(p.id)}
            className={`h-[22px] p-4.5 text-[11px] font-medium rounded-[3px] border-0 cursor-pointer transition-colors
              ${portal === p.id
                ? "bg-[#D50C2D] text-white"
                : "bg-transparent text-gray-400 hover:bg-[#1F2937] hover:text-gray-200"
              }`}
          >
            {p.label}
          </button>
        ))}
        <span className="ml-auto text-[10px] text-gray-600">Demo — switch portal above</span>
      </div>

      {/* Portal content */}
      <div className="flex-1 min-h-0">
        {portal === "admin" && <AdminPortal />}
        {portal === "reseller" && <ResellerPortal />}
        {portal === "client" && <ClientPortal />}
      </div>
    </div>
  );
}
