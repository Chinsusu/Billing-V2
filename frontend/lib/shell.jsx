// Shared chrome + data helpers for HANetwork admin

// ─── Logo ───────────────────────────────────────────────────────
const HALogo = ({ collapsed, portal }) => {
  const brand = portal === 'reseller'
    ? { initial: 'P', name: 'ProxyVN', sub: 'Portal', bg: 'var(--accent)' }
    : portal === 'customer'
    ? { initial: 'P', name: 'ProxyVN', sub: 'Client', bg: 'var(--accent)' }
    : { initial: 'H', name: 'HANetwork', sub: 'Console', bg: 'var(--accent)' };
  return (
    <div style={{
      display: 'flex', alignItems: 'center', gap: 8,
      height: 48, padding: collapsed ? '0' : '0 14px',
      justifyContent: collapsed ? 'center' : 'flex-start',
      borderBottom: '1px solid var(--line)',
    }}>
      <div style={{
        width: 22, height: 22,
        display: 'grid', placeItems: 'center',
        background: brand.bg, color: '#fff',
        fontSize: 12, fontWeight: 700, letterSpacing: -0.5,
        borderRadius: 2,
      }}>{brand.initial}</div>
      {!collapsed && (
        <div style={{ fontSize: 13, fontWeight: 600, letterSpacing: -0.1, color: 'var(--ink-0)' }}>
          {brand.name}<span style={{ color: 'var(--ink-3)', fontWeight: 400, marginLeft: 4 }}>{brand.sub}</span>
        </div>
      )}
    </div>
  );
};

// ─── Sidebar ────────────────────────────────────────────────────
const NAV_SECTIONS = [
  {
    label: null,
    items: [
      { id: 'overview', label: 'Overview', icon: 'dashboard' },
    ],
  },
  {
    label: 'Customers',
    items: [
      { id: 'customers', label: 'Customers', icon: 'users', count: 2847 },
      { id: 'tickets', label: 'Support tickets', icon: 'ticket', count: 23, badge: 'danger' },
    ],
  },
  {
    label: 'Services',
    items: [
      { id: 'proxies', label: 'Proxies', icon: 'globe', count: 12480 },
      { id: 'vps', label: 'VPS & Servers', icon: 'server', count: 3124 },
      { id: 'bandwidth', label: 'Bandwidth', icon: 'chart' },
    ],
  },
  {
    label: 'Billing',
    items: [
      { id: 'invoices', label: 'Invoices', icon: 'file' },
      { id: 'transactions', label: 'Transactions', icon: 'card' },
      { id: 'products', label: 'Products & Pricing', icon: 'tag' },
    ],
  },
  {
    label: 'System',
    items: [
      { id: 'reports', label: 'Reports', icon: 'chart' },
      { id: 'settings', label: 'Settings', icon: 'settings' },
    ],
  },
];

const getNav = (portal) => {
  if (portal === 'customer') return CUSTOMER_NAV;
  if (portal === 'reseller') return RESELLER_NAV;
  // admin: append v2 extras to existing NAV_SECTIONS
  return [NAV_SECTIONS[0], ...ADMIN_V2_EXTRA, ...NAV_SECTIONS.slice(1)];
};

const Sidebar = ({ active, onSelect, collapsed, onToggle, customer, portal }) => (
  <aside style={{
    width: collapsed ? 'var(--sidebar-w-collapsed)' : 'var(--sidebar-w)',
    flexShrink: 0,
    background: 'var(--surface)',
    borderRight: '1px solid var(--line)',
    display: 'flex',
    flexDirection: 'column',
    transition: 'width .15s ease',
    overflow: 'hidden',
  }}>
    <HALogo collapsed={collapsed} portal={portal} />

    <nav style={{ flex: 1, overflowY: 'auto', padding: '6px 8px' }} className="no-scrollbar">
      {(portal ? getNav(portal) : customer ? CUSTOMER_NAV : NAV_SECTIONS).map((sec, i) => (
        <div key={i} style={{ marginBottom: 10 }}>
          {sec.label && !collapsed && (
            <div style={{
              fontSize: 10, fontWeight: 600, letterSpacing: 0.6,
              color: 'var(--ink-4)', textTransform: 'uppercase',
              padding: '10px 8px 4px',
            }}>{sec.label}</div>
          )}
          {sec.label && collapsed && <div style={{ height: 6 }} />}
          {sec.items.map(it => (
            <button
              key={it.id}
              onClick={() => onSelect(it.id)}
              title={collapsed ? it.label : undefined}
              style={{
                display: 'flex', alignItems: 'center',
                gap: 10,
                width: '100%',
                padding: collapsed ? '0' : '0 8px',
                height: 28,
                border: 'none',
                background: active === it.id ? 'var(--accent-soft)' : 'transparent',
                color: active === it.id ? 'var(--accent)' : 'var(--ink-1)',
                fontSize: 13,
                fontWeight: active === it.id ? 500 : 400,
                borderRadius: 3,
                cursor: 'pointer',
                justifyContent: collapsed ? 'center' : 'flex-start',
                textAlign: 'left',
                marginBottom: 1,
              }}
              onMouseEnter={e => { if (active !== it.id) e.currentTarget.style.background = 'var(--surface-hover)'; }}
              onMouseLeave={e => { if (active !== it.id) e.currentTarget.style.background = 'transparent'; }}
            >
              <Icon name={it.icon} size={15} />
              {!collapsed && (
                <>
                  <span style={{ flex: 1 }}>{it.label}</span>
                  {it.count != null && (
                    <span style={{
                      fontSize: 10, fontWeight: 500, color: 'var(--ink-3)',
                      fontVariantNumeric: 'tabular-nums',
                    }}>{it.count.toLocaleString()}</span>
                  )}
                  {it.badge === 'danger' && it.count != null && (
                    <span style={{
                      fontSize: 10, fontWeight: 600, color: '#fff',
                      background: 'var(--accent)', padding: '0 5px',
                      borderRadius: 8, lineHeight: '14px', height: 14,
                    }}>{it.count}</span>
                  )}
                </>
              )}
            </button>
          ))}
        </div>
      ))}
    </nav>

    <div style={{
      borderTop: '1px solid var(--line)',
      padding: collapsed ? '6px' : '10px 12px',
      display: 'flex', alignItems: 'center', gap: 10,
      justifyContent: collapsed ? 'center' : 'space-between',
    }}>
      {!collapsed && (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, minWidth: 0 }}>
          <div style={{
            width: 26, height: 26, borderRadius: 13,
            background: portal === 'reseller' ? 'var(--accent)' : portal === 'customer' ? '#4B5563' : '#1F2937',
            color: '#fff',
            display: 'grid', placeItems: 'center',
            fontSize: 11, fontWeight: 600,
          }}>{portal === 'reseller' ? 'PV' : portal === 'customer' ? 'LT' : 'MN'}</div>
          <div style={{ minWidth: 0 }}>
            <div style={{ fontSize: 12, fontWeight: 500, color: 'var(--ink-0)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
              {portal === 'reseller' ? 'ProxyVN' : portal === 'customer' ? 'Linh Tran' : 'Minh Nguyen'}
            </div>
            <div style={{ fontSize: 11, color: 'var(--ink-3)' }}>
              {portal === 'reseller' ? 'Reseller · T-0042' : portal === 'customer' ? 'Client · via ProxyVN' : 'Administrator'}
            </div>
          </div>
        </div>
      )}
      <button className="btn btn-ghost btn-sm" onClick={onToggle} style={{ height: 24, padding: 4 }} title="Toggle sidebar">
        <Icon name={collapsed ? 'chevronRight' : 'chevronLeft'} size={14} />
      </button>
    </div>
  </aside>
);

const CUSTOMER_NAV = [
  { label: null, items: [
    { id: 'c-overview', label: 'Dashboard', icon: 'dashboard' },
    { id: 'c-shop', label: 'Shop', icon: 'tag' },
  ]},
  { label: 'My services', items: [
    { id: 'c-services', label: 'All services', icon: 'server', count: 5 },
    { id: 'c-checkout', label: 'Checkout · preview', icon: 'card' },
  ]},
  { label: 'Billing', items: [
    { id: 'c-wallet', label: 'Wallet', icon: 'wallet' },
    { id: 'c-usage', label: 'Usage', icon: 'chart' },
  ]},
  { label: 'Account', items: [
    { id: 'c-settings', label: 'Settings', icon: 'settings' },
    { id: 'c-support', label: 'Support', icon: 'ticket' },
  ]},
];

// Admin nav adds Tenants / Provisioning / Top-ups / Providers from v2 spec
const ADMIN_V2_EXTRA = [
  { label: 'Platform', items: [
    { id: 'tenants', label: 'Tenants (Resellers)', icon: 'users', count: 5 },
    { id: 'provisioning', label: 'Provisioning queue', icon: 'server', count: 6, badge: 'danger' },
    { id: 'topups', label: 'Top-up verification', icon: 'wallet', count: 3, badge: 'danger' },
    { id: 'providers', label: 'Providers / Sources', icon: 'globe' },
  ]},
];

const RESELLER_NAV = [
  { label: null, items: [
    { id: 'r-overview', label: 'Dashboard', icon: 'dashboard' },
  ]},
  { label: 'My business', items: [
    { id: 'r-clients', label: 'Clients', icon: 'users', count: 312 },
    { id: 'r-catalog', label: 'Catalog / Pricing', icon: 'tag' },
    { id: 'r-services', label: 'Services', icon: 'server' },
    { id: 'r-orders', label: 'Orders', icon: 'file' },
  ]},
  { label: 'Finance', items: [
    { id: 'r-wallet', label: 'Wallet & Top-up', icon: 'wallet' },
    { id: 'r-reports', label: 'Reports', icon: 'chart' },
  ]},
  { label: 'Account', items: [
    { id: 'r-settings', label: 'Branding & Settings', icon: 'settings' },
  ]},
];

// ─── Topbar ─────────────────────────────────────────────────────
const Topbar = ({ title, breadcrumbs, actions, meta }) => (
  <header style={{
    height: 'var(--topbar-h)',
    padding: '0 20px',
    background: 'var(--surface)',
    borderBottom: '1px solid var(--line)',
    display: 'flex', alignItems: 'center', gap: 16,
    flexShrink: 0,
  }}>
    <div style={{ flex: 1, minWidth: 0 }}>
      {breadcrumbs && (
        <div style={{ display: 'flex', alignItems: 'center', gap: 4, fontSize: 11, color: 'var(--ink-3)', marginBottom: 1 }}>
          {breadcrumbs.map((c, i) => (
            <React.Fragment key={i}>
              {i > 0 && <Icon name="chevronRight" size={10} style={{ opacity: .5 }}/>}
              <span>{c}</span>
            </React.Fragment>
          ))}
        </div>
      )}
      <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
        <h1 style={{ margin: 0, fontSize: 15, fontWeight: 600, letterSpacing: -0.1, color: 'var(--ink-0)' }}>{title}</h1>
        {meta}
      </div>
    </div>

    <div style={{
      display: 'flex', alignItems: 'center', gap: 6,
      background: 'var(--bg-alt)', padding: '4px 8px 4px 8px',
      borderRadius: 3, width: 260,
      border: '1px solid transparent',
    }}>
      <Icon name="search" size={13} style={{ color: 'var(--ink-3)' }}/>
      <input placeholder="Search customers, services, invoices…"
        style={{ border: 'none', background: 'transparent', outline: 'none', flex: 1, fontSize: 12, color: 'var(--ink-1)', fontFamily: 'inherit' }}
      />
      <span className="kbd">⌘K</span>
    </div>

    <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
      <button className="btn btn-ghost btn-sm" title="Refresh"><Icon name="refresh" size={14}/></button>
      <button className="btn btn-ghost btn-sm" style={{ position: 'relative' }} title="Notifications">
        <Icon name="bell" size={14}/>
        <span style={{ position: 'absolute', top: 4, right: 4, width: 6, height: 6, borderRadius: 3, background: 'var(--accent)' }}/>
      </button>
      {actions}
    </div>
  </header>
);

// ─── KPI Tile ───────────────────────────────────────────────────
const Kpi = ({ label, value, unit, delta, spark, deltaType = 'auto', sub }) => {
  const num = typeof delta === 'string' ? parseFloat(delta) : delta;
  const positive = deltaType === 'auto' ? num >= 0 : deltaType === 'positive';
  return (
    <div style={{
      background: 'var(--surface)',
      border: '1px solid var(--line)',
      borderRadius: 'var(--r-md)',
      padding: '14px 16px',
      display: 'flex', flexDirection: 'column', gap: 8,
      minWidth: 0,
    }}>
      <div style={{ fontSize: 11, fontWeight: 500, color: 'var(--ink-3)', textTransform: 'uppercase', letterSpacing: 0.5 }}>{label}</div>
      <div style={{ display: 'flex', alignItems: 'baseline', gap: 4 }}>
        <div style={{ fontSize: 22, fontWeight: 600, letterSpacing: -0.4, color: 'var(--ink-0)', fontVariantNumeric: 'tabular-nums' }}>{value}</div>
        {unit && <div style={{ fontSize: 12, color: 'var(--ink-3)' }}>{unit}</div>}
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 6, minHeight: 18 }}>
        {delta != null && (
          <span style={{
            display: 'inline-flex', alignItems: 'center', gap: 2,
            fontSize: 11, fontWeight: 500,
            color: positive ? 'var(--ok)' : 'var(--danger)',
          }}>
            <Icon name={positive ? 'arrowUp' : 'arrowDown'} size={10} stroke={2}/>
            {typeof delta === 'number' ? `${Math.abs(delta)}%` : delta}
          </span>
        )}
        {sub && <span style={{ fontSize: 11, color: 'var(--ink-3)' }}>{sub}</span>}
        {spark && <div style={{ marginLeft: 'auto' }}>{spark}</div>}
      </div>
    </div>
  );
};

// ─── Sparkline ──────────────────────────────────────────────────
const Sparkline = ({ data, w = 80, h = 22, color = 'var(--ink-3)', fill = false }) => {
  const max = Math.max(...data);
  const min = Math.min(...data);
  const rng = max - min || 1;
  const step = w / (data.length - 1);
  const pts = data.map((v, i) => `${(i*step).toFixed(1)},${(h - ((v-min)/rng)*h).toFixed(1)}`).join(' ');
  return (
    <svg width={w} height={h} style={{ display: 'block' }}>
      {fill && <polygon points={`0,${h} ${pts} ${w},${h}`} fill={color} fillOpacity={.12}/>}
      <polyline points={pts} fill="none" stroke={color} strokeWidth="1.3" strokeLinejoin="round" strokeLinecap="round"/>
    </svg>
  );
};

// ─── Status dot ─────────────────────────────────────────────────
const StatusDot = ({ status }) => {
  const map = {
    active: { c: 'var(--ok)', pulse: true },
    running: { c: 'var(--ok)', pulse: true },
    paid: { c: 'var(--ok)' },
    open: { c: 'var(--info)' },
    pending: { c: 'var(--warn)' },
    overdue: { c: 'var(--danger)' },
    failed: { c: 'var(--danger)' },
    suspended: { c: 'var(--ink-4)' },
    stopped: { c: 'var(--ink-4)' },
    provisioning: { c: 'var(--info)', pulse: true },
  };
  const s = map[status] || map.stopped;
  return (
    <span style={{
      display: 'inline-block', width: 7, height: 7, borderRadius: 4,
      background: s.c, flexShrink: 0,
      boxShadow: s.pulse ? `0 0 0 0 ${s.c}` : 'none',
      animation: s.pulse ? 'haPulse 2s infinite' : 'none',
    }}/>
  );
};

// ─── Bar chart (simple columns) ─────────────────────────────────
const BarChart = ({ data, h = 120, color = 'var(--accent)', labels, valueFmt = v => v }) => {
  const max = Math.max(...data);
  return (
    <div style={{ display: 'flex', alignItems: 'flex-end', gap: 4, height: h, position: 'relative' }}>
      {data.map((v, i) => (
        <div key={i} style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 4 }}>
          <div style={{ fontSize: 10, color: 'var(--ink-3)', fontVariantNumeric: 'tabular-nums' }}>
            {valueFmt(v)}
          </div>
          <div style={{
            width: '100%',
            height: `${(v/max) * (h - 36)}px`,
            background: color,
            borderRadius: '2px 2px 0 0',
            transition: 'height .3s',
            minHeight: 2,
          }}/>
          {labels && <div style={{ fontSize: 10, color: 'var(--ink-3)' }}>{labels[i]}</div>}
        </div>
      ))}
    </div>
  );
};

// ─── Line chart (area) ─────────────────────────────────────────
const LineArea = ({ series, w = 600, h = 180, color = 'var(--accent)', labels, yFmt = v => v }) => {
  const max = Math.max(...series.flat());
  const n = series[0].length;
  const step = (w - 40) / (n - 1);
  const mapY = v => h - 30 - (v/max) * (h - 50);
  const toPts = arr => arr.map((v,i) => `${(30 + i*step).toFixed(1)},${mapY(v).toFixed(1)}`).join(' ');
  const gridLines = 4;
  return (
    <svg width={w} height={h} style={{ display: 'block' }}>
      {Array.from({length: gridLines + 1}).map((_,i) => {
        const y = 20 + ((h - 50)/gridLines) * i;
        const val = max * (1 - i/gridLines);
        return (
          <g key={i}>
            <line x1={30} y1={y} x2={w-10} y2={y} stroke="var(--line-2)" strokeWidth="1"/>
            <text x={4} y={y+3} fontSize="10" fill="var(--ink-3)">{yFmt(val)}</text>
          </g>
        );
      })}
      {labels && labels.map((l,i) => (
        <text key={i} x={30 + i*step} y={h-8} fontSize="10" fill="var(--ink-3)" textAnchor="middle">{l}</text>
      ))}
      {series.map((arr,si) => {
        const c = si === 0 ? color : 'var(--ink-4)';
        return (
          <g key={si}>
            {si === 0 && (
              <polygon points={`30,${h-30} ${toPts(arr)} ${30+(n-1)*step},${h-30}`} fill={c} fillOpacity={.08}/>
            )}
            <polyline points={toPts(arr)} fill="none" stroke={c} strokeWidth="1.6" strokeLinejoin="round"/>
          </g>
        );
      })}
    </svg>
  );
};

// Keyframes
if (!document.getElementById('ha-keyframes')) {
  const s = document.createElement('style');
  s.id = 'ha-keyframes';
  s.textContent = `
  @keyframes haPulse {
    0% { box-shadow: 0 0 0 0 currentColor; }
    70% { box-shadow: 0 0 0 5px rgba(0,0,0,0); }
    100% { box-shadow: 0 0 0 0 rgba(0,0,0,0); }
  }`;
  document.head.appendChild(s);
}

Object.assign(window, { Sidebar, Topbar, Kpi, Sparkline, StatusDot, BarChart, LineArea, HALogo });
