// Variation 3: Data-dense Stripe-like — icon-only nav, huge numbers, minimal chrome
const OverviewV3 = () => (
  <div style={{ padding: '16px 24px', display: 'flex', flexDirection: 'column', gap: 16, background: '#fff' }}>
    <div className="row">
      <div>
        <h1 style={{ margin: 0, fontSize: 18, fontWeight: 600 }}>Today</h1>
        <div className="tiny muted" style={{ marginTop: 2 }}>Wed, Apr 22 2026 · all amounts in USD</div>
      </div>
      <div className="hstack">
        <div style={{ display: 'flex', border: '1px solid var(--line)', borderRadius: 3, overflow: 'hidden' }}>
          {['Today','7D','30D','QTD','YTD','All'].map((p,i) => (
            <button key={p} style={{
              padding: '4px 12px', border: 'none',
              background: i === 2 ? 'var(--ink-0)' : '#fff',
              color: i === 2 ? '#fff' : 'var(--ink-2)',
              fontSize: 11, fontWeight: 500, cursor: 'pointer',
              borderRight: i < 5 ? '1px solid var(--line)' : 'none',
              fontFamily: 'inherit',
            }}>{p}</button>
          ))}
        </div>
      </div>
    </div>

    {/* Big numbers row, Stripe-style */}
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4,1fr)', gap: 0, borderTop: '1px solid var(--line)', borderBottom: '1px solid var(--line)' }}>
      {[
        { l: 'Gross volume', v: '$29,204.82', d: '+12.1%', dprev: '$26,048.00 previous', sp: SAMPLE.revenue30d, color: 'var(--accent)' },
        { l: 'Net volume', v: '$27,124.48', d: '+11.8%', dprev: '$24,248.92 previous', sp: SAMPLE.revenue30d.map(v=>v*0.93), color: 'var(--ok)' },
        { l: 'New customers', v: '34', d: '+21.4%', dprev: '28 previous', sp: SAMPLE.customers30d.map(v=>v-2800), color: 'var(--info)' },
        { l: 'Failed payments', v: '3', d: '+1', dprev: '2 previous', sp: [1,2,1,0,2,3,2,3,1,3], color: 'var(--danger)' },
      ].map((k,i) => (
        <div key={i} style={{
          padding: '20px 20px 24px',
          borderRight: i < 3 ? '1px solid var(--line)' : 'none',
          minWidth: 0,
        }}>
          <div style={{ fontSize: 12, color: 'var(--ink-3)' }}>{k.l}</div>
          <div style={{ fontSize: 26, fontWeight: 600, letterSpacing: -0.6, fontVariantNumeric: 'tabular-nums', marginTop: 4 }}>{k.v}</div>
          <div className="hstack" style={{ gap: 6, marginTop: 4 }}>
            <span style={{ fontSize: 12, fontWeight: 500, color: k.d.startsWith('-') || k.l === 'Failed payments' && k.d.startsWith('+') ? 'var(--danger)' : 'var(--ok)' }}>{k.d}</span>
            <span className="tiny muted">{k.dprev}</span>
          </div>
          <div style={{ marginTop: 10 }}>
            <Sparkline data={k.sp} color={k.color} w={240} h={30} fill/>
          </div>
        </div>
      ))}
    </div>

    {/* Deep chart */}
    <div style={{ padding: '12px 0' }}>
      <div className="row" style={{ marginBottom: 14 }}>
        <div className="hstack" style={{ gap: 20 }}>
          {['Volume','Payments','Customers','Disputes'].map((t,i) => (
            <button key={t} style={{
              border: 'none', background: 'transparent',
              padding: '4px 0', fontSize: 13,
              color: i === 0 ? 'var(--ink-0)' : 'var(--ink-3)',
              fontWeight: i === 0 ? 600 : 400,
              borderBottom: i === 0 ? '2px solid var(--ink-0)' : 'none',
              cursor: 'pointer', fontFamily: 'inherit',
            }}>{t}</button>
          ))}
        </div>
      </div>
      <LineArea
        series={[SAMPLE.revenue30d, SAMPLE.revenue30d.map(v => v * 0.82 + 2000)]}
        w={940} h={240} color="var(--accent)"
        labels={Array.from({length:30},(_,i)=>i%5===0?`Apr ${i+1}`:'')}
        yFmt={v => '$' + (v/1000).toFixed(0) + 'k'}
      />
    </div>

    {/* Dense tables row */}
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 32, borderTop: '1px solid var(--line)', paddingTop: 16 }}>
      <div>
        <div className="row" style={{ marginBottom: 10 }}>
          <h3 style={{ margin: 0, fontSize: 13, fontWeight: 600 }}>Top customers by revenue</h3>
          <a href="#" style={{ fontSize: 11 }}>View all</a>
        </div>
        {[
          ['DataMine Inc.','$8,420.00','64 services'],
          ['CloudHarvest','$6,240.00','48 services'],
          ['Acme Proxy Co.','$4,280.00','24 services'],
          ['Scrapers Ltd','$1,840.00','16 services'],
          ['Sofia Bergström','$1,480.00','21 services'],
          ['Marie Dubois','$1,120.00','12 services'],
        ].map((r,i) => (
          <div key={i} style={{ display: 'grid', gridTemplateColumns: '1fr auto auto', gap: 16, padding: '6px 0', fontSize: 12.5, borderBottom: i < 5 ? '1px solid var(--line-2)' : 'none' }}>
            <span style={{ fontWeight: 500 }}>{r[0]}</span>
            <span style={{ color: 'var(--ink-3)' }}>{r[2]}</span>
            <span style={{ fontWeight: 500, fontVariantNumeric: 'tabular-nums' }}>{r[1]}</span>
          </div>
        ))}
      </div>
      <div>
        <div className="row" style={{ marginBottom: 10 }}>
          <h3 style={{ margin: 0, fontSize: 13, fontWeight: 600 }}>Recent payments</h3>
          <a href="#" style={{ fontSize: 11 }}>View all</a>
        </div>
        {SAMPLE.transactions.slice(0,6).map((tx,i) => (
          <div key={tx.id} style={{ display: 'grid', gridTemplateColumns: 'auto 1fr auto auto', gap: 10, padding: '6px 0', fontSize: 12.5, borderBottom: i < 5 ? '1px solid var(--line-2)' : 'none', alignItems: 'center' }}>
            <StatusDot status={tx.status}/>
            <span>{tx.customer}</span>
            <span className="tiny muted mono">{tx.id.slice(-6)}</span>
            <span style={{ fontWeight: 500, fontVariantNumeric: 'tabular-nums' }}>{fmtMoney(tx.amount)}</span>
          </div>
        ))}
      </div>
    </div>
  </div>
);

window.OverviewV3 = OverviewV3;
