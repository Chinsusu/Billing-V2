// Admin Overview — classic Hetzner style (dense, light, minimal)

const OverviewScreen = () => {
  return (
    <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
      {/* KPI row */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12 }}>
        <Kpi label="MRR" value="$118.2k" delta={8.4} sub="vs last month"
          spark={<Sparkline data={SAMPLE.mrr30d} color="var(--accent)" fill w={70} h={24}/>}/>
        <Kpi label="Revenue · MTD" value="$29.2k" delta={12.1} sub="Apr 2026"
          spark={<Sparkline data={SAMPLE.revenue30d} color="var(--ok)" fill w={70} h={24}/>}/>
        <Kpi label="Active customers" value="2,847" delta={1.2} sub="net +34"
          spark={<Sparkline data={SAMPLE.customers30d} w={70} h={24}/>}/>
        <Kpi label="Active services" value="15,604" delta={-0.3} sub="proxies + VPS"
          spark={<Sparkline data={[15400,15420,15510,15490,15580,15620,15604]} w={70} h={24}/>}/>
      </div>

      {/* Main grid */}
      <div style={{ display: 'grid', gridTemplateColumns: '1.7fr 1fr', gap: 16 }}>
        {/* Revenue chart */}
        <div className="card">
          <div className="card-header">
            <h3>Revenue · last 30 days</h3>
            <div className="hstack">
              <div className="hstack" style={{ gap: 12, fontSize: 11, color: 'var(--ink-3)' }}>
                <span className="hstack" style={{ gap: 5 }}>
                  <span style={{ width: 8, height: 2, background: 'var(--accent)' }}/>This period
                </span>
                <span className="hstack" style={{ gap: 5 }}>
                  <span style={{ width: 8, height: 2, background: 'var(--ink-4)' }}/>Previous
                </span>
              </div>
              <select className="select" style={{ width: 110, height: 26 }}>
                <option>Last 30 days</option>
                <option>This quarter</option>
                <option>YTD</option>
              </select>
            </div>
          </div>
          <div className="card-body" style={{ padding: '14px 10px 10px' }}>
            <LineArea
              series={[SAMPLE.revenue30d, SAMPLE.revenue30d.map(v => v * 0.82 + 2000)]}
              w={600} h={200} color="var(--accent)"
              labels={['Apr 1','','','','','','','','Apr 10','','','','','','','','Apr 18','','','','','','Apr 24','','','','','','','Apr 30']}
              yFmt={v => '$' + (v/1000).toFixed(0) + 'k'}
            />
          </div>
        </div>

        {/* Revenue by product */}
        <div className="card">
          <div className="card-header">
            <h3>Revenue by product · 30d</h3>
            <button className="btn btn-ghost btn-sm"><Icon name="more" size={14}/></button>
          </div>
          <div className="card-body" style={{ padding: 12 }}>
            {(() => {
              const total = SAMPLE.products.reduce((s,p)=>s+p.rev,0);
              return SAMPLE.products.map((p,i) => (
                <div key={i} style={{ padding: '6px 4px', borderBottom: i < SAMPLE.products.length-1 ? '1px solid var(--line-2)' : 'none' }}>
                  <div className="row" style={{ marginBottom: 4 }}>
                    <div className="hstack" style={{ gap: 8 }}>
                      <span style={{ width: 8, height: 8, background: p.color, borderRadius: 1 }}/>
                      <span style={{ fontSize: 12, color: 'var(--ink-1)' }}>{p.name}</span>
                    </div>
                    <div className="hstack" style={{ gap: 10, fontSize: 12, fontVariantNumeric: 'tabular-nums' }}>
                      <span style={{ color: 'var(--ink-3)' }}>{p.sold.toLocaleString()}</span>
                      <span style={{ fontWeight: 500, minWidth: 56, textAlign: 'right' }}>${(p.rev/1000).toFixed(1)}k</span>
                    </div>
                  </div>
                  <div style={{ height: 3, background: 'var(--line-2)', borderRadius: 1 }}>
                    <div style={{ width: `${(p.rev/total)*100}%`, height: '100%', background: p.color, borderRadius: 1 }}/>
                  </div>
                </div>
              ));
            })()}
          </div>
        </div>
      </div>

      {/* Second row: recent activity + outstanding */}
      <div style={{ display: 'grid', gridTemplateColumns: '1.7fr 1fr', gap: 16 }}>
        <div className="card">
          <div className="card-header">
            <h3>Recent invoices</h3>
            <a href="#" style={{ fontSize: 12 }}>View all →</a>
          </div>
          <table className="tbl">
            <thead>
              <tr>
                <th>Invoice</th>
                <th>Customer</th>
                <th>Issued</th>
                <th>Due</th>
                <th className="num">Amount</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {SAMPLE.invoices.slice(0, 7).map(inv => (
                <tr key={inv.id}>
                  <td className="mono" style={{ color: 'var(--accent)' }}>{inv.id}</td>
                  <td>{inv.customer}</td>
                  <td style={{ color: 'var(--ink-3)' }}>{inv.issued}</td>
                  <td style={{ color: 'var(--ink-3)' }}>{inv.due}</td>
                  <td className="num" style={{ fontWeight: 500 }}>{fmtMoney(inv.amount)}</td>
                  <td>
                    <span className={`badge dot ${STATUS_BADGE[inv.status]}`}>{STATUS_LABEL[inv.status]}</span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="card">
          <div className="card-header">
            <h3>Activity feed</h3>
            <span className="tiny muted">Live</span>
          </div>
          <div style={{ maxHeight: 340, overflow: 'hidden' }}>
            {SAMPLE.activity.map((a,i) => (
              <div key={i} style={{
                display: 'flex', gap: 10, padding: '9px 14px',
                borderBottom: i < SAMPLE.activity.length-1 ? '1px solid var(--line-2)' : 'none',
                alignItems: 'flex-start',
              }}>
                <div style={{
                  width: 22, height: 22, borderRadius: 11,
                  background: a.type === 'ok' ? 'var(--ok-bg)' : a.type === 'warn' ? 'var(--warn-bg)' : a.type === 'danger' ? 'var(--danger-bg)' : 'var(--muted-bg)',
                  color: a.type === 'ok' ? 'var(--ok)' : a.type === 'warn' ? 'var(--warn)' : a.type === 'danger' ? 'var(--danger)' : 'var(--ink-2)',
                  display: 'grid', placeItems: 'center', flexShrink: 0,
                }}>
                  <Icon name={a.icon} size={11} stroke={2}/>
                </div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ fontSize: 12, color: 'var(--ink-1)', lineHeight: 1.35 }}>{a.text}</div>
                  <div style={{ fontSize: 11, color: 'var(--ink-3)', marginTop: 1 }}>{a.t}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Third row: bandwidth + system health */}
      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: 16 }}>
        <div className="card">
          <div className="card-header">
            <h3>Bandwidth consumption · 30 days (TB)</h3>
            <div className="hstack" style={{ fontSize: 11, color: 'var(--ink-3)' }}>
              Peak: 342 GB · Avg: 267 GB
            </div>
          </div>
          <div className="card-body" style={{ padding: '16px 20px' }}>
            <BarChart data={SAMPLE.bandwidthDaily} h={110} color="var(--accent)"
              valueFmt={() => ''} />
            <div className="row tiny muted" style={{ marginTop: 4 }}>
              <span>Apr 1</span>
              <span>Apr 8</span>
              <span>Apr 15</span>
              <span>Apr 22</span>
              <span>Apr 30</span>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="card-header">
            <h3>Infrastructure health</h3>
            <span className="badge ok dot">Operational</span>
          </div>
          <div style={{ padding: '4px 0' }}>
            {[
              { label: 'Proxy network · uptime', value: '99.98%', bar: 0.9998 },
              { label: 'VPS fleet · uptime', value: '99.94%', bar: 0.9994 },
              { label: 'Payment gateway', value: '100%', bar: 1 },
              { label: 'API · p95 latency', value: '142ms', bar: 0.82 },
              { label: 'Support · first response', value: '8m avg', bar: 0.72 },
            ].map((r,i) => (
              <div key={i} style={{ padding: '8px 16px', borderBottom: i < 4 ? '1px solid var(--line-2)' : 'none' }}>
                <div className="row" style={{ marginBottom: 5 }}>
                  <span style={{ fontSize: 12 }}>{r.label}</span>
                  <span style={{ fontSize: 12, fontWeight: 500, fontVariantNumeric: 'tabular-nums' }}>{r.value}</span>
                </div>
                <div style={{ height: 3, background: 'var(--line-2)' }}>
                  <div style={{ width: `${r.bar*100}%`, height: '100%', background: r.bar > 0.99 ? 'var(--ok)' : r.bar > 0.8 ? 'var(--info)' : 'var(--warn)' }}/>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

window.OverviewScreen = OverviewScreen;
