// Products / Pricing, Tickets, Bandwidth, Reports, Settings

const ProductsScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 12 }}>
    <div className="row">
      <div>
        <h2 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>Products & Pricing</h2>
        <div className="tiny muted" style={{ marginTop: 2 }}>11 active SKUs across 7 product families</div>
      </div>
      <div className="hstack">
        <button className="btn btn-sm"><Icon name="copy" size={12}/> Price history</button>
        <button className="btn btn-primary btn-sm"><Icon name="plus" size={12}/> New product</button>
      </div>
    </div>

    <div className="card" style={{ overflow: 'hidden' }}>
      <table className="tbl">
        <thead>
          <tr>
            <th>SKU</th>
            <th>Product</th>
            <th>Billing unit</th>
            <th className="num">Unit price</th>
            <th className="num">Active subs</th>
            <th className="num">Revenue · 30d</th>
            <th>Status</th>
            <th style={{ width: 40 }}></th>
          </tr>
        </thead>
        <tbody>
          {SAMPLE.products_catalog.map(p => (
            <tr key={p.sku}>
              <td className="mono" style={{ color: 'var(--ink-3)' }}>{p.sku}</td>
              <td style={{ fontWeight: 500 }}>{p.name}</td>
              <td style={{ color: 'var(--ink-2)' }}>{p.unit}</td>
              <td className="num" style={{ fontWeight: 500 }}>${p.price.toFixed(2)}</td>
              <td className="num">{p.active.toLocaleString()}</td>
              <td className="num" style={{ fontWeight: 500 }}>${p.rev30.toLocaleString()}</td>
              <td><span className="badge ok dot">Active</span></td>
              <td><button className="btn btn-ghost btn-sm" style={{ padding: 4 }}><Icon name="more" size={13}/></button></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

const TicketsScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 12 }}>
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4,1fr)', gap: 12 }}>
      <Kpi label="Open tickets" value="23" delta="+3" deltaType="negative" sub="needs response"/>
      <Kpi label="Avg first response" value="8m" delta={-12} sub="SLA 15m"/>
      <Kpi label="Resolved · 24h" value="47" delta={18} sub="91% CSAT"/>
      <Kpi label="High priority" value="4" delta="+1" deltaType="negative" sub="2 overdue"/>
    </div>

    <div className="card" style={{ overflow: 'hidden' }}>
      <div style={{ padding: '10px 14px', borderBottom: '1px solid var(--line)', display: 'flex', gap: 16 }}>
        {['Open · 23','Pending · 12','Resolved','All'].map((t,i) => (
          <button key={i} style={{
            border: 'none', background: 'transparent',
            padding: '4px 0', fontSize: 13,
            color: i === 0 ? 'var(--ink-0)' : 'var(--ink-3)',
            fontWeight: i === 0 ? 600 : 400,
            borderBottom: i === 0 ? '2px solid var(--accent)' : '2px solid transparent',
            cursor: 'pointer',
          }}>{t}</button>
        ))}
      </div>
      <table className="tbl">
        <thead>
          <tr>
            <th>#</th>
            <th>Subject</th>
            <th>Customer</th>
            <th>Priority</th>
            <th>Assignee</th>
            <th>Updated</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {SAMPLE.tickets.map(t => (
            <tr key={t.id}>
              <td className="mono" style={{ color: 'var(--ink-3)' }}>{t.id}</td>
              <td>
                <div style={{ fontWeight: 500, color: 'var(--ink-0)' }}>{t.subject}</div>
              </td>
              <td style={{ color: 'var(--ink-2)' }}>{t.customer}</td>
              <td>
                <span className={`badge dot ${t.priority === 'high' ? 'danger' : t.priority === 'medium' ? 'warn' : ''}`}>
                  {t.priority}
                </span>
              </td>
              <td>{t.assignee}</td>
              <td style={{ color: 'var(--ink-3)' }}>{t.updated}</td>
              <td><span className={`badge dot ${STATUS_BADGE[t.status]}`}>{STATUS_LABEL[t.status]}</span></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

const BandwidthScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 12 }}>
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4,1fr)', gap: 12 }}>
      <Kpi label="Total usage · MTD" value="8.4 TB" delta={12.1} sub="Apr 2026"/>
      <Kpi label="Avg daily" value="267 GB" delta={4.2}/>
      <Kpi label="Peak day" value="342 GB" sub="Apr 29"/>
      <Kpi label="Overage revenue" value="$1,280" delta={-8.4} sub="42 customers"/>
    </div>

    <div className="card">
      <div className="card-header">
        <h3>Bandwidth by region · 30 days</h3>
      </div>
      <div className="card-body">
        <LineArea
          series={[
            SAMPLE.bandwidthDaily,
            SAMPLE.bandwidthDaily.map(v => v * 0.6),
            SAMPLE.bandwidthDaily.map(v => v * 0.3),
          ]}
          w={860} h={220}
          color="var(--accent)"
          labels={['Apr 1','','','','','','','','Apr 10','','','','','','','','Apr 18','','','','','','Apr 24','','','','','','','Apr 30']}
          yFmt={v => v.toFixed(0) + ' GB'}
        />
      </div>
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
      <div className="card">
        <div className="card-header"><h3>Top consumers · 30d</h3></div>
        <table className="tbl">
          <thead>
            <tr><th>Customer</th><th className="num">Usage</th><th className="num">Cost</th></tr>
          </thead>
          <tbody>
            {[
              { c: 'Acme Proxy Co.', u: '2.4 TB', cost: 380 },
              { c: 'DataMine Inc.', u: '1.8 TB', cost: 284 },
              { c: 'CloudHarvest', u: '842 GB', cost: 164 },
              { c: 'Scrapers Ltd', u: '624 GB', cost: 118 },
              { c: 'Marie Dubois', u: '380 GB', cost: 72 },
            ].map((r,i) => (
              <tr key={i}>
                <td>{r.c}</td>
                <td className="num">{r.u}</td>
                <td className="num" style={{ fontWeight: 500 }}>${r.cost}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="card">
        <div className="card-header"><h3>Pool utilization</h3></div>
        <div style={{ padding: '4px 0' }}>
          {[
            { name: 'US Residential · Pool A', used: 78 },
            { name: 'US Residential · Pool B', used: 42 },
            { name: 'EU Residential · Premium', used: 61 },
            { name: 'APAC Residential', used: 24 },
            { name: 'DC · US-EAST', used: 89 },
            { name: 'DC · US-WEST', used: 54 },
            { name: 'Mobile 4G · US', used: 92 },
          ].map((p,i) => (
            <div key={i} style={{ padding: '8px 16px', borderBottom: i < 6 ? '1px solid var(--line-2)' : 'none' }}>
              <div className="row" style={{ marginBottom: 5 }}>
                <span style={{ fontSize: 12 }}>{p.name}</span>
                <span style={{ fontSize: 12, fontWeight: 500, color: p.used > 85 ? 'var(--danger)' : p.used > 70 ? 'var(--warn)' : 'var(--ink-2)' }}>{p.used}%</span>
              </div>
              <div style={{ height: 3, background: 'var(--line-2)' }}>
                <div style={{ width: `${p.used}%`, height: '100%', background: p.used > 85 ? 'var(--danger)' : p.used > 70 ? 'var(--warn)' : 'var(--ok)' }}/>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  </div>
);

const ReportsScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 12 }}>
    <div className="row">
      <h2 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>Reports</h2>
      <button className="btn btn-primary btn-sm"><Icon name="plus" size={12}/> New report</button>
    </div>
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3,1fr)', gap: 12 }}>
      {[
        { title: 'Monthly Recurring Revenue', desc: 'MRR breakdown by product and plan', updated: '2h ago' },
        { title: 'Customer Cohort Analysis', desc: 'Retention by signup month', updated: 'yesterday' },
        { title: 'Churn & Cancellations', desc: 'Why customers leave — reasons & values', updated: '3d ago' },
        { title: 'Product Mix', desc: 'Revenue distribution across product lines', updated: '1h ago' },
        { title: 'Bandwidth & Overage', desc: 'Usage vs plan limits; overage fees', updated: '4h ago' },
        { title: 'Support Performance', desc: 'SLA, first response, CSAT by agent', updated: 'today' },
      ].map((r,i) => (
        <div key={i} className="card" style={{ padding: 16, cursor: 'pointer', transition: 'border-color .12s' }}>
          <div className="hstack" style={{ justifyContent: 'space-between', marginBottom: 6 }}>
            <Icon name="chart" size={18} style={{ color: 'var(--accent)' }}/>
            <span className="tiny muted">{r.updated}</span>
          </div>
          <div style={{ fontSize: 13, fontWeight: 500, marginBottom: 2 }}>{r.title}</div>
          <div className="tiny muted">{r.desc}</div>
        </div>
      ))}
    </div>
  </div>
);

const SettingsScreen = () => {
  const [sec, setSec] = React.useState('billing');
  const sections = [
    { id: 'profile', label: 'Organization profile', icon: 'settings' },
    { id: 'billing', label: 'Billing & tax', icon: 'card' },
    { id: 'payment', label: 'Payment methods', icon: 'wallet' },
    { id: 'team', label: 'Team members', icon: 'users' },
    { id: 'api', label: 'API keys', icon: 'shield' },
    { id: 'webhooks', label: 'Webhooks', icon: 'external' },
    { id: 'notifications', label: 'Notifications', icon: 'bell' },
    { id: 'security', label: 'Security & 2FA', icon: 'shield' },
  ];
  return (
    <div style={{ padding: 20, display: 'grid', gridTemplateColumns: '220px 1fr', gap: 20 }}>
      <div>
        {sections.map(s => (
          <button key={s.id} onClick={() => setSec(s.id)} style={{
            display: 'flex', alignItems: 'center', gap: 8,
            width: '100%', padding: '6px 10px', border: 'none',
            background: sec === s.id ? 'var(--accent-soft)' : 'transparent',
            color: sec === s.id ? 'var(--accent)' : 'var(--ink-1)',
            fontSize: 12.5, fontWeight: sec === s.id ? 500 : 400,
            borderRadius: 3, cursor: 'pointer', textAlign: 'left', fontFamily: 'inherit',
            marginBottom: 1,
          }}>
            <Icon name={s.icon} size={14}/>{s.label}
          </button>
        ))}
      </div>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
        <div className="card">
          <div className="card-header"><h3>Billing & tax configuration</h3></div>
          <div className="card-body" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
            <SettingsRow label="Billing currency" value="USD · US Dollar" />
            <SettingsRow label="Default tax rate" value="10% VAT (Vietnam)" />
            <SettingsRow label="Invoice numbering" value="INV-{YYYY}-{NNNNN}" />
            <SettingsRow label="Invoice due period" value="Net 14 days" />
            <SettingsRow label="Auto-collect" value="Enabled · 2 retry attempts" />
            <SettingsRow label="Grace period before suspension" value="7 days after due date" />
          </div>
        </div>
        <div className="card">
          <div className="card-header"><h3>Revenue recognition</h3></div>
          <div className="card-body" style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
            <SettingsRow label="Recognition method" value="Accrual — daily" />
            <SettingsRow label="Fiscal year" value="Jan 1 – Dec 31" />
            <SettingsRow label="Accounting export" value="Xero connected" />
          </div>
        </div>
      </div>
    </div>
  );
};

const SettingsRow = ({ label, value }) => (
  <div className="row" style={{ alignItems: 'flex-start' }}>
    <div style={{ fontSize: 12, color: 'var(--ink-2)', width: 240 }}>{label}</div>
    <div style={{ flex: 1, fontSize: 13, color: 'var(--ink-0)', fontWeight: 500 }}>{value}</div>
    <button className="btn btn-sm">Edit</button>
  </div>
);

Object.assign(window, { ProductsScreen, TicketsScreen, BandwidthScreen, ReportsScreen, SettingsScreen });
