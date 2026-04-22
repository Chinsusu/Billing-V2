// Variation 2: Modern minimal — more whitespace, larger type, softer borders
const OverviewV2 = () => (
  <div style={{ padding: 28, display: 'flex', flexDirection: 'column', gap: 24, background: '#FAFAF9' }}>
    <div className="row">
      <div>
        <div style={{ fontSize: 11, color: 'var(--ink-3)', letterSpacing: 0.5, textTransform: 'uppercase' }}>Overview</div>
        <h1 style={{ margin: '4px 0 0', fontSize: 22, fontWeight: 600, letterSpacing: -0.4 }}>Good afternoon, Minh.</h1>
        <div className="tiny muted" style={{ marginTop: 4 }}>Here's what's happening across HANetwork today · Wed, Apr 22</div>
      </div>
      <div className="hstack">
        <select className="select" style={{ width: 130 }}><option>Last 30 days</option></select>
        <button className="btn"><Icon name="download" size={13}/> Export</button>
      </div>
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4,1fr)', gap: 16 }}>
      {[
        { l: 'MRR', v: '$118.2k', d: '+8.4%', sub: 'Monthly recurring', sp: SAMPLE.mrr30d, color: 'var(--accent)' },
        { l: 'Revenue · MTD', v: '$29.2k', d: '+12.1%', sub: 'of $32k target', sp: SAMPLE.revenue30d, color: 'var(--ok)' },
        { l: 'Customers', v: '2,847', d: '+34 net', sub: 'this month', sp: SAMPLE.customers30d, color: 'var(--info)' },
        { l: 'Churn rate', v: '1.8%', d: '−0.4pp', sub: 'trailing 30d', sp: [2.4,2.2,2.3,2.1,2.0,1.9,1.8,1.8], color: 'var(--ok)' },
      ].map((k,i) => (
        <div key={i} style={{ background: '#fff', borderRadius: 8, padding: 20, border: '1px solid #F0EEE9', position: 'relative', overflow: 'hidden' }}>
          <div style={{ fontSize: 12, color: 'var(--ink-3)', marginBottom: 10 }}>{k.l}</div>
          <div style={{ fontSize: 28, fontWeight: 600, letterSpacing: -0.6, fontVariantNumeric: 'tabular-nums' }}>{k.v}</div>
          <div className="hstack" style={{ marginTop: 4, justifyContent: 'space-between' }}>
            <div className="hstack" style={{ gap: 6 }}>
              <span style={{ color: k.d.startsWith('-') || k.d.startsWith('−') ? 'var(--ok)' : 'var(--ok)', fontSize: 12, fontWeight: 500 }}>{k.d}</span>
              <span className="tiny muted">{k.sub}</span>
            </div>
          </div>
          <div style={{ marginTop: 12, marginLeft: -6, marginRight: -6 }}>
            <Sparkline data={k.sp} color={k.color} fill w={240} h={36}/>
          </div>
        </div>
      ))}
    </div>

    <div style={{ background: '#fff', borderRadius: 8, border: '1px solid #F0EEE9', padding: 20 }}>
      <div className="row" style={{ marginBottom: 20 }}>
        <div>
          <h3 style={{ margin: 0, fontSize: 15, fontWeight: 600 }}>Revenue trend</h3>
          <div className="tiny muted" style={{ marginTop: 2 }}>Daily revenue with previous period comparison</div>
        </div>
        <div className="hstack" style={{ fontSize: 11 }}>
          <span className="hstack" style={{ gap: 5 }}><span style={{ width: 8, height: 2, background: 'var(--accent)' }}/>Current</span>
          <span className="hstack" style={{ gap: 5, color: 'var(--ink-3)' }}><span style={{ width: 8, height: 2, background: 'var(--ink-4)' }}/>Previous</span>
        </div>
      </div>
      <LineArea
        series={[SAMPLE.revenue30d, SAMPLE.revenue30d.map(v => v * 0.82 + 2000)]}
        w={920} h={220} color="var(--accent)"
        labels={Array.from({length:30},(_,i)=>i%5===0?`Apr ${i+1}`:'')}
        yFmt={v => '$' + (v/1000).toFixed(0) + 'k'}
      />
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
      <div style={{ background: '#fff', borderRadius: 8, border: '1px solid #F0EEE9', padding: 20 }}>
        <h3 style={{ margin: '0 0 14px', fontSize: 14, fontWeight: 600 }}>Revenue by product</h3>
        {SAMPLE.products.map((p,i) => {
          const total = SAMPLE.products.reduce((s,x)=>s+x.rev,0);
          return (
            <div key={i} style={{ padding: '10px 0', borderBottom: i < SAMPLE.products.length-1 ? '1px solid #F5F3EF' : 'none' }}>
              <div className="row" style={{ marginBottom: 6 }}>
                <div className="hstack" style={{ gap: 10 }}>
                  <span style={{ width: 6, height: 6, background: p.color, borderRadius: 3 }}/>
                  <span style={{ fontSize: 13 }}>{p.name}</span>
                </div>
                <span style={{ fontSize: 13, fontWeight: 500, fontVariantNumeric: 'tabular-nums' }}>${(p.rev/1000).toFixed(1)}k</span>
              </div>
              <div style={{ height: 2, background: '#F0EEE9', borderRadius: 1 }}>
                <div style={{ width: `${(p.rev/total)*100}%`, height: '100%', background: p.color, borderRadius: 1 }}/>
              </div>
            </div>
          );
        })}
      </div>

      <div style={{ background: '#fff', borderRadius: 8, border: '1px solid #F0EEE9', padding: 20 }}>
        <h3 style={{ margin: '0 0 14px', fontSize: 14, fontWeight: 600 }}>Attention required</h3>
        {[
          { icon: 'x', type: 'danger', title: '3 failed charges in last 24h', sub: '$840.00 at risk · retry scheduled', cta: 'Review' },
          { icon: 'ticket', type: 'warn', title: '4 high-priority tickets open', sub: '2 have breached SLA · avg wait 42m', cta: 'Triage' },
          { icon: 'clock', type: 'warn', title: '23 services overdue on renewal', sub: 'Total $4,120 MRR at risk', cta: 'View' },
          { icon: 'users', type: 'info', title: '8 customers signed up today', sub: '2 converted from trial to Pro', cta: 'Welcome' },
        ].map((a,i) => (
          <div key={i} style={{ display: 'flex', gap: 12, padding: '12px 0', borderBottom: i < 3 ? '1px solid #F5F3EF' : 'none', alignItems: 'center' }}>
            <div style={{
              width: 32, height: 32, borderRadius: 16,
              background: a.type === 'danger' ? 'var(--danger-bg)' : a.type === 'warn' ? 'var(--warn-bg)' : 'var(--info-bg)',
              color: a.type === 'danger' ? 'var(--danger)' : a.type === 'warn' ? 'var(--warn)' : 'var(--info)',
              display: 'grid', placeItems: 'center', flexShrink: 0,
            }}>
              <Icon name={a.icon} size={14} stroke={2}/>
            </div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <div style={{ fontSize: 13, fontWeight: 500 }}>{a.title}</div>
              <div className="tiny muted" style={{ marginTop: 1 }}>{a.sub}</div>
            </div>
            <button className="btn btn-sm">{a.cta}</button>
          </div>
        ))}
      </div>
    </div>
  </div>
);

window.OverviewV2 = OverviewV2;
