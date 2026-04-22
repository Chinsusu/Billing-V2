// Customer-side dashboard (end-user view)

const CustomerOverviewScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div style={{
      background: 'linear-gradient(180deg, #FFFFFF 0%, #FAFBFC 100%)',
      border: '1px solid var(--line)',
      borderRadius: 4,
      padding: '18px 20px',
      display: 'grid', gridTemplateColumns: '1fr auto', gap: 20, alignItems: 'center',
    }}>
      <div>
        <div className="tiny muted">Welcome back</div>
        <div style={{ fontSize: 20, fontWeight: 600, marginTop: 2 }}>Linh Trần</div>
        <div className="tiny muted" style={{ marginTop: 4 }}>Pro plan · customer since Sep 2024 · ID C-40217</div>
      </div>
      <div style={{ display: 'flex', gap: 18, alignItems: 'center' }}>
        <div style={{ textAlign: 'right' }}>
          <div className="tiny muted">Wallet balance</div>
          <div style={{ fontSize: 22, fontWeight: 600, fontVariantNumeric: 'tabular-nums' }}>$248.50</div>
        </div>
        <button className="btn btn-primary"><Icon name="plus" size={12}/> Top up</button>
      </div>
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4,1fr)', gap: 12 }}>
      <Kpi label="Active services" value="8" sub="2 proxies · 3 VPS · 3 bandwidth"/>
      <Kpi label="This month" value="$189.00" sub="due May 14" delta={0}/>
      <Kpi label="Bandwidth used" value="128" unit="GB" sub="of 500 GB · 25.6%"/>
      <Kpi label="Next renewal" value="14 days" sub="EU Residential · Premium"/>
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: '1.5fr 1fr', gap: 16 }}>
      <div className="card">
        <div className="card-header">
          <h3>My services</h3>
          <button className="btn btn-primary btn-sm"><Icon name="plus" size={12}/> Order new</button>
        </div>
        <table className="tbl">
          <thead>
            <tr><th>Service</th><th>Type</th><th>Region</th><th className="num">Price</th><th>Renews</th><th>Status</th></tr>
          </thead>
          <tbody>
            {SAMPLE.services.slice(0,6).map(s => (
              <tr key={s.id}>
                <td style={{ fontWeight: 500 }}>{s.label.split(' · ')[0]}</td>
                <td><TypeIcon type={s.type}/></td>
                <td><span className="mono tiny" style={{ padding: '1px 5px', background: 'var(--bg-alt)', borderRadius: 2 }}>{s.region}</span></td>
                <td className="num">${s.price}/mo</td>
                <td style={{ color: 'var(--ink-3)' }}>in {s.renewsIn > 0 ? s.renewsIn : 0}d</td>
                <td><span className={`badge dot ${STATUS_BADGE[s.status]}`}>{STATUS_LABEL[s.status]}</span></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="card">
        <div className="card-header"><h3>Bandwidth usage</h3></div>
        <div className="card-body">
          <div style={{ fontSize: 11, color: 'var(--ink-3)', marginBottom: 4 }}>Residential Pool · Premium</div>
          <div className="row" style={{ alignItems: 'baseline', marginBottom: 8 }}>
            <div><span style={{ fontSize: 24, fontWeight: 600 }}>128</span><span style={{ fontSize: 13, color: 'var(--ink-3)' }}> / 500 GB</span></div>
            <span className="badge ok dot">25.6%</span>
          </div>
          <div style={{ height: 6, background: 'var(--line-2)', borderRadius: 3, marginBottom: 18 }}>
            <div style={{ width: '25.6%', height: '100%', background: 'var(--accent)', borderRadius: 3 }}/>
          </div>
          <Sparkline data={SAMPLE.bandwidthDaily.slice(-14).map(v=>v*0.04)} w={280} h={56} color="var(--accent)" fill/>
          <div className="row tiny muted" style={{ marginTop: 4 }}>
            <span>Apr 9</span><span>Apr 15</span><span>Apr 22</span>
          </div>
        </div>
      </div>
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
      <div className="card">
        <div className="card-header"><h3>Recent invoices</h3><a href="#" style={{ fontSize: 12 }}>All invoices →</a></div>
        <table className="tbl">
          <thead><tr><th>Invoice</th><th>Date</th><th className="num">Amount</th><th>Status</th></tr></thead>
          <tbody>
            {SAMPLE.invoices.slice(0,5).map(inv => (
              <tr key={inv.id}>
                <td className="mono" style={{ color: 'var(--accent)' }}>{inv.id}</td>
                <td style={{ color: 'var(--ink-3)' }}>{inv.issued}</td>
                <td className="num" style={{ fontWeight: 500 }}>{fmtMoney(inv.amount)}</td>
                <td><span className={`badge dot ${STATUS_BADGE[inv.status]}`}>{STATUS_LABEL[inv.status]}</span></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="card">
        <div className="card-header"><h3>Quick actions</h3></div>
        <div style={{ padding: 10, display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
          {[
            { icon: 'plus', label: 'Order proxy', desc: 'Residential · DC · ISP · Mobile' },
            { icon: 'server', label: 'Deploy VPS', desc: 'Linux or Windows · 50+ locations' },
            { icon: 'wallet', label: 'Top up wallet', desc: 'Add funds via card or crypto' },
            { icon: 'download', label: 'Download invoices', desc: 'Bulk export for accounting' },
            { icon: 'shield', label: 'API credentials', desc: 'Manage tokens and IP whitelist' },
            { icon: 'ticket', label: 'Contact support', desc: 'Avg response: 8 minutes' },
          ].map((a,i) => (
            <button key={i} style={{
              display: 'flex', gap: 10, padding: 10,
              border: '1px solid var(--line)', borderRadius: 3,
              background: 'var(--surface)', cursor: 'pointer',
              textAlign: 'left', fontFamily: 'inherit',
            }}>
              <div style={{ width: 28, height: 28, borderRadius: 3, background: 'var(--accent-soft)', color: 'var(--accent)', display: 'grid', placeItems: 'center', flexShrink: 0 }}>
                <Icon name={a.icon} size={14}/>
              </div>
              <div style={{ minWidth: 0 }}>
                <div style={{ fontSize: 12.5, fontWeight: 500, color: 'var(--ink-0)' }}>{a.label}</div>
                <div className="tiny muted" style={{ marginTop: 1 }}>{a.desc}</div>
              </div>
            </button>
          ))}
        </div>
      </div>
    </div>
  </div>
);

window.CustomerOverviewScreen = CustomerOverviewScreen;
