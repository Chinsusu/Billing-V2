// Invoices + Transactions

const InvoicesScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 12 }}>
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4,1fr)', gap: 12 }}>
      <Kpi label="Outstanding" value="$6,460" sub="3 invoices open" deltaType="positive" delta="-12%"/>
      <Kpi label="Paid · MTD" value="$22,208" delta={14.3} sub="7 invoices"/>
      <Kpi label="Overdue" value="$420" delta={-68} sub="1 invoice · 7d"/>
      <Kpi label="Avg days to pay" value="3.4" sub="target 5 days" unit="d" delta={-8.1}/>
    </div>

    <div className="card" style={{ padding: 12, display: 'flex', alignItems: 'center', gap: 8, flexWrap: 'wrap' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 6, background: 'var(--bg-alt)', padding: '4px 10px', borderRadius: 3, width: 260 }}>
        <Icon name="search" size={13} style={{ color: 'var(--ink-3)' }}/>
        <input placeholder="Invoice # or customer…" style={{ border: 'none', background: 'transparent', outline: 'none', flex: 1, fontSize: 12, fontFamily: 'inherit' }}/>
      </div>
      <FilterChip label="Status" value="Any"/>
      <FilterChip label="Period" value="Apr 2026"/>
      <FilterChip label="Amount" value="Any"/>
      <div style={{ marginLeft: 'auto', display: 'flex', gap: 6 }}>
        <button className="btn btn-sm"><Icon name="download" size={12}/> Export PDF</button>
        <button className="btn btn-primary btn-sm"><Icon name="plus" size={12}/> Create invoice</button>
      </div>
    </div>

    <div className="card" style={{ overflow: 'hidden' }}>
      <table className="tbl">
        <thead>
          <tr>
            <th>Invoice #</th>
            <th>Customer</th>
            <th>Issued</th>
            <th>Due</th>
            <th>Items</th>
            <th className="num">Subtotal</th>
            <th className="num">Tax</th>
            <th className="num">Total</th>
            <th>Status</th>
            <th style={{ width: 80 }}></th>
          </tr>
        </thead>
        <tbody>
          {SAMPLE.invoices.map(inv => (
            <tr key={inv.id}>
              <td className="mono" style={{ color: 'var(--accent)', fontWeight: 500 }}>{inv.id}</td>
              <td>{inv.customer}</td>
              <td style={{ color: 'var(--ink-3)' }}>{inv.issued}</td>
              <td style={{ color: inv.status === 'overdue' ? 'var(--danger)' : 'var(--ink-3)' }}>{inv.due}</td>
              <td>{Math.floor(Math.random() * 5) + 1} line{Math.random() > 0.5 ? 's' : ''}</td>
              <td className="num">{fmtMoney(inv.amount * 0.9)}</td>
              <td className="num" style={{ color: 'var(--ink-3)' }}>{fmtMoney(inv.amount * 0.1)}</td>
              <td className="num" style={{ fontWeight: 600 }}>{fmtMoney(inv.amount)}</td>
              <td><span className={`badge dot ${STATUS_BADGE[inv.status]}`}>{STATUS_LABEL[inv.status]}</span></td>
              <td>
                <div className="hstack" style={{ gap: 2 }}>
                  <button className="btn btn-ghost btn-sm" style={{ padding: 4 }}><Icon name="eye" size={13}/></button>
                  <button className="btn btn-ghost btn-sm" style={{ padding: 4 }}><Icon name="download" size={13}/></button>
                  <button className="btn btn-ghost btn-sm" style={{ padding: 4 }}><Icon name="more" size={13}/></button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

const TransactionsScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 12 }}>
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4,1fr)', gap: 12 }}>
      <Kpi label="Processed · today" value="$22,028" delta={18.2} sub="9 transactions"/>
      <Kpi label="Success rate" value="94.2" unit="%" delta={-1.4} sub="last 24h"/>
      <Kpi label="Failed charges" value="3" delta="+1" deltaType="negative" sub="$840 at risk"/>
      <Kpi label="Refunded · MTD" value="$129" delta={-42} sub="2 refunds"/>
    </div>

    <div className="card" style={{ padding: 12, display: 'flex', alignItems: 'center', gap: 8 }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 6, background: 'var(--bg-alt)', padding: '4px 10px', borderRadius: 3, width: 260 }}>
        <Icon name="search" size={13} style={{ color: 'var(--ink-3)' }}/>
        <input placeholder="Transaction ID, customer…" style={{ border: 'none', background: 'transparent', outline: 'none', flex: 1, fontSize: 12, fontFamily: 'inherit' }}/>
      </div>
      <FilterChip label="Type" value="All"/>
      <FilterChip label="Method" value="All"/>
      <FilterChip label="Status" value="All"/>
      <FilterChip label="Date" value="Last 7d"/>
    </div>

    <div className="card" style={{ overflow: 'hidden' }}>
      <table className="tbl">
        <thead>
          <tr>
            <th>Transaction ID</th>
            <th>Time</th>
            <th>Customer</th>
            <th>Type</th>
            <th>Method</th>
            <th className="num">Amount</th>
            <th>Status</th>
            <th style={{ width: 40 }}></th>
          </tr>
        </thead>
        <tbody>
          {SAMPLE.transactions.map(tx => (
            <tr key={tx.id}>
              <td className="mono" style={{ color: 'var(--ink-1)' }}>{tx.id}</td>
              <td className="mono tiny" style={{ color: 'var(--ink-3)' }}>{tx.time}</td>
              <td>{tx.customer}</td>
              <td>
                <span style={{
                  fontSize: 11, padding: '1px 6px', borderRadius: 2,
                  background: tx.type === 'charge' ? 'var(--muted-bg)' : tx.type === 'topup' ? 'var(--info-bg)' : 'var(--warn-bg)',
                  color: tx.type === 'charge' ? 'var(--ink-1)' : tx.type === 'topup' ? 'var(--info)' : 'var(--warn)',
                }}>{tx.type}</span>
              </td>
              <td>{tx.method}</td>
              <td className="num" style={{ fontWeight: 500, color: tx.amount < 0 ? 'var(--danger)' : 'var(--ink-0)' }}>{fmtMoney(tx.amount)}</td>
              <td><span className={`badge dot ${STATUS_BADGE[tx.status]}`}>{STATUS_LABEL[tx.status]}</span></td>
              <td><button className="btn btn-ghost btn-sm" style={{ padding: 4 }}><Icon name="more" size={13}/></button></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

window.InvoicesScreen = InvoicesScreen;
window.TransactionsScreen = TransactionsScreen;
