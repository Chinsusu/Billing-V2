"use client";

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
          {collapsed ? "→" : "←"}
        </button>
      </div>
    </aside>
  );
}

function NavIcon({ id }: { id: string }) {
  const icons: Record<string, string> = {
    "admin-overview": "⊞", "reseller-overview": "⊞", "client-overview": "⊞",
    "admin-tenants": "⊙", "admin-provisioning": "▣", "admin-topups": "◎", "admin-providers": "⊕",
    "admin-customers": "⊚", "reseller-clients": "⊚", "admin-tickets": "◈",
    "admin-services": "▦", "reseller-services": "▦", "client-services": "▦",
    "admin-invoices": "◻", "reseller-orders": "◻",
    "admin-transactions": "▤", "reseller-wallet": "◑", "client-wallet": "◑",
    "admin-products": "◇", "reseller-catalog": "◇", "client-shop": "◇",
    "admin-settings": "⊛", "reseller-settings": "⊛", "client-settings": "⊛",
    "reseller-reports": "△", "client-usage": "△",
    "client-support": "◫",
  };
  return <span className="text-[13px] shrink-0 leading-none">{icons[id] ?? "·"}</span>;
}
