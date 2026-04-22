"use client";

import React from "react";
import {
  LayoutDashboard, Users, Server, ArrowDownCircle, Plug, Headphones,
  Receipt, ArrowLeftRight, Package, Settings, BarChart2, ShoppingBag,
  Wallet, Wrench, UserCheck, ChevronLeft, ChevronRight, Bell, ScrollText,
  Globe, Gauge,
} from "lucide-react";
import { NAV_BY_PORTAL, PORTAL_META, Portal, NavItem } from "@/lib/navigation/screens";

interface SidebarProps {
  portal: Portal;
  activeScreen: string;
  onSelect: (id: string) => void;
  collapsed: boolean;
  onToggle: () => void;
}

export function Sidebar({ portal, activeScreen, onSelect, collapsed, onToggle }: SidebarProps) {
  const nav = NAV_BY_PORTAL[portal];
  const meta = PORTAL_META[portal];

  const sections = nav.reduce<{ label: string | undefined; items: NavItem[] }[]>((acc, item) => {
    const last = acc[acc.length - 1];
    if (!last || last.label !== item.section) {
      acc.push({ label: item.section, items: [item] });
    } else {
      last.items.push(item);
    }
    return acc;
  }, []);

  return (
    <aside
      className="flex flex-col bg-white border-r border-gray-200 shrink-0 transition-all duration-150"
      style={{ width: collapsed ? 56 : 224 }}
    >
      {/* Logo */}
      <div
        className="flex items-center gap-2 border-b border-gray-200 shrink-0"
        style={{ height: 48, padding: collapsed ? "0" : "0 14px", justifyContent: collapsed ? "center" : "flex-start" }}
      >
        <div className="w-[22px] h-[22px] grid place-items-center bg-[#D50C2D] text-white text-xs font-bold rounded-sm shrink-0">
          {meta.initial}
        </div>
        {!collapsed && (
          <span className="text-[13px] font-semibold text-gray-900 truncate">
            {meta.label}
          </span>
        )}
      </div>

      {/* Nav */}
      <nav className="flex-1 overflow-y-auto px-2 py-1.5 space-y-2" style={{ scrollbarWidth: "none" }}>
        {sections.map((sec, i) => (
          <div key={i}>
            {sec.label && !collapsed && (
              <div className="text-[10px] font-semibold uppercase tracking-widest text-gray-400 px-2 pt-2.5 pb-1">
                {sec.label}
              </div>
            )}
            {sec.items.map((item) => {
              const active = activeScreen === item.id;
              return (
                <button
                  key={item.id}
                  onClick={() => onSelect(item.id)}
                  title={collapsed ? item.label : undefined}
                  className={`flex items-center gap-2.5 w-full h-7 rounded-[3px] text-[13px] cursor-pointer border-0 transition-colors
                    ${collapsed ? "justify-center px-0" : "px-2"}
                    ${active
                      ? "bg-red-50 text-[#D50C2D] font-medium"
                      : "text-gray-700 hover:bg-gray-100 bg-transparent font-normal"
                    }`}
                >
                  <NavIcon id={item.id} />
                  {!collapsed && (
                    <>
                      <span className="flex-1 text-left truncate">{item.label}</span>
                      {item.badge === "danger" && item.count != null && (
                        <span className="text-[10px] font-semibold text-white bg-[#D50C2D] px-1 rounded-full leading-[14px]">
                          {item.count}
                        </span>
                      )}
                      {!item.badge && item.count != null && (
                        <span className="text-[10px] text-gray-400 tabular-nums">{item.count.toLocaleString()}</span>
                      )}
                    </>
                  )}
                </button>
              );
            })}
          </div>
        ))}
      </nav>

      {/* User + toggle */}
      <div
        className="border-t border-gray-200 flex items-center gap-2.5"
        style={{ padding: collapsed ? 6 : "10px 12px", justifyContent: collapsed ? "center" : "space-between" }}
      >
        {!collapsed && (
          <div className="flex items-center gap-2 min-w-0">
            <div className="w-[26px] h-[26px] rounded-full bg-gray-800 text-white grid place-items-center text-[11px] font-semibold shrink-0">
              {meta.user.slice(0, 2).toUpperCase()}
            </div>
            <div className="min-w-0">
              <div className="text-[12px] font-medium text-gray-900 truncate">{meta.user}</div>
              <div className="text-[11px] text-gray-400 truncate">{meta.role}</div>
            </div>
          </div>
        )}
        <button
          onClick={onToggle}
          className="text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded p-1 transition-colors border-0 bg-transparent cursor-pointer"
          title="Toggle sidebar"
        >
          {collapsed ? <ChevronRight size={14} /> : <ChevronLeft size={14} />}
        </button>
      </div>
    </aside>
  );
}

const NAV_ICONS: Record<string, React.ReactNode> = {
  "admin-overview":     <LayoutDashboard size={14} />,
  "reseller-overview":  <LayoutDashboard size={14} />,
  "client-overview":    <LayoutDashboard size={14} />,
  "admin-tenants":      <Users size={14} />,
  "admin-provisioning": <Wrench size={14} />,
  "admin-topups":       <ArrowDownCircle size={14} />,
  "admin-providers":    <Plug size={14} />,
  "admin-customers":    <UserCheck size={14} />,
  "reseller-clients":   <UserCheck size={14} />,
  "admin-tickets":      <Headphones size={14} />,
  "client-support":     <Headphones size={14} />,
  "admin-services-proxies":   <Globe size={14} />,
  "admin-services-vps":       <Server size={14} />,
  "admin-services-bandwidth": <Gauge size={14} />,
  "reseller-services":  <Server size={14} />,
  "client-services":    <Server size={14} />,
  "admin-invoices":     <Receipt size={14} />,
  "reseller-orders":    <Receipt size={14} />,
  "admin-transactions": <ArrowLeftRight size={14} />,
  "reseller-wallet":    <Wallet size={14} />,
  "client-wallet":      <Wallet size={14} />,
  "admin-products":     <Package size={14} />,
  "reseller-catalog":   <Package size={14} />,
  "client-shop":        <ShoppingBag size={14} />,
  "admin-alerts":       <Bell size={14} />,
  "admin-logs":         <ScrollText size={14} />,
  "admin-settings":     <Settings size={14} />,
  "reseller-settings":  <Settings size={14} />,
  "client-settings":    <Settings size={14} />,
  "reseller-reports":   <BarChart2 size={14} />,
  "client-usage":       <BarChart2 size={14} />,
};

function NavIcon({ id }: { id: string }) {
  return (
    <span className="shrink-0 flex items-center justify-center w-[14px]">
      {NAV_ICONS[id] ?? <span className="w-1.5 h-1.5 rounded-full bg-current opacity-40" />}
    </span>
  );
}
