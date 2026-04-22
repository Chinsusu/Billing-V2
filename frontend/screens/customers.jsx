// Customers list

const CustomersScreen = () => {
  const [selected, setSelected] = React.useState(new Set());
  const toggle = id => setSelected(s => {
    const n = new Set(s);
    n.has(id) ? n.delete(id) : n.add(id);
    return n;
  });

  return (
    <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 12 }}>
      {/* Filter bar */}
      <div className="card" style={{ padding: 12, display: 'flex', alignItems: 'center', gap: 8, flexWrap: 'wrap' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 6, background: 'var(--bg-alt)', padding: '4px 10px', borderRadius: 3, width: 260 }}>
          <Icon name="search" size={13} style={{ color: 'var(--ink-3)' }}/>
          <input placeholder="Search by name, email, ID…" style={{ border: 'none', background: 'transparent', outline: 'none', flex: 1, fontSize: 12, fontFamily: 'inherit' }}/>
        </div>
        <span style={{ width: 1, height: 20, background: 'var(--line)' }}/>
        <FilterChip label="Status" value="Any"/>
        <FilterChip label="Plan" value="Any"/>
        <FilterChip label="Country" value="Any"/>
        <FilterChip label="MRR" value="Any"/>
        <button className="btn btn-ghost btn-sm"><Icon name="plus" size={12}/> Add filter</button>
        <div style={{ marginLeft: 'auto', display: 'flex', gap: 6 }}>
          <button className="btn btn-sm"><Icon name="download" size={12}/> Export CSV</button>
          <button className="btn btn-primary btn-sm"><Icon name="plus" size={12}/> New customer</button>
        </div>
      </div>

      {/* Summary bar */}
      <div style={{ display: 'flex', alignItems: 'center', gap: 24, padding: '8px 12px', fontSize: 12 }}>
        <span style={{ color: 'var(--ink-3)' }}>Showing <strong style={{ color: 'var(--ink-0)' }}>{SAMPLE.customers.length}</strong> of <strong style={{ color: 'var(--ink-0)' }}>2,847</strong> customers</span>
        <span style={{ color: 'var(--ink-3)' }}>Total MRR: <strong style={{ color: 'var(--ink-0)' }}>$24,716</strong></span>
        <span style={{ color: 'var(--ink-3)' }}>Overdue: <strong style={{ color: 'var(--danger)' }}>1</strong></span>
        <span style={{ color: 'var(--ink-3)' }}>Suspended: <strong style={{ color: 'var(--ink-2)' }}>1</strong></span>
      </div>

      {/* Table */}
      <div className="card" style={{ overflow: 'hidden' }}>
        <table className="tbl">
          <thead>
            <tr>
              <th style={{ width: 28, paddingRight: 0 }}><input type="checkbox" style={{ margin: 0 }}/></th>
              <th>Customer</th>
              <th>ID</th>
              <th>Plan</th>
              <th>Country</th>
              <th className="num">Services</th>
              <th className="num">MRR</th>
              <th>Customer since</th>
              <th>Status</th>
              <th style={{ width: 40 }}></th>
            </tr>
          </thead>
          <tbody>
            {SAMPLE.customers.map(c => (
              <tr key={c.id}>
                <td style={{ paddingRight: 0 }}><input type="checkbox" checked={selected.has(c.id)} onChange={() => toggle(c.id)} style={{ margin: 0 }}/></td>
                <td>
                  <div className="hstack" style={{ gap: 8 }}>
                    <div style={{
                      width: 24, height: 24, borderRadius: 12,
                      background: 'var(--bg-alt)', color: 'var(--ink-2)',
                      display: 'grid', placeItems: 'center',
                      fontSize: 10, fontWeight: 600,
                    }}>{c.name.split(' ').map(w => w[0]).slice(0,2).join('')}</div>
                    <div style={{ minWidth: 0 }}>
                      <div style={{ fontWeight: 500, color: 'var(--ink-0)' }}>{c.name}</div>
                      <div className="tiny muted">{c.email}</div>
                    </div>
                  </div>
                </td>
                <td className="mono" style={{ color: 'var(--ink-3)' }}>{c.id}</td>
                <td>
                  <span style={{
                    fontSize: 11, padding: '1px 6px',
                    background: c.plan === 'Enterprise' ? 'var(--ink-0)' : c.plan === 'Business' ? 'var(--ink-1)' : c.plan === 'Pro' ? 'var(--accent-soft)' : 'var(--muted-bg)',
                    color: c.plan === 'Enterprise' || c.plan === 'Business' ? '#fff' : c.plan === 'Pro' ? 'var(--accent)' : 'var(--ink-2)',
                    borderRadius: 2, fontWeight: 500,
                  }}>{c.plan}</span>
                </td>
                <td>
                  <span className="hstack" style={{ gap: 5, color: 'var(--ink-2)' }}>
                    <span style={{ fontSize: 10, fontWeight: 600, padding: '1px 4px', background: 'var(--bg-alt)', borderRadius: 2, fontFamily: 'var(--font-mono)' }}>{c.country}</span>
                  </span>
                </td>
                <td className="num" style={{ fontVariantNumeric: 'tabular-nums' }}>{c.services}</td>
                <td className="num" style={{ fontVariantNumeric: 'tabular-nums', fontWeight: 500 }}>${c.mrr.toLocaleString()}</td>
                <td style={{ color: 'var(--ink-3)' }}>{c.since}</td>
                <td><span className={`badge dot ${STATUS_BADGE[c.status]}`}>{STATUS_LABEL[c.status]}</span></td>
                <td><button className="btn btn-ghost btn-sm" style={{ padding: 4 }}><Icon name="more" size={14}/></button></td>
              </tr>
            ))}
          </tbody>
        </table>
        {/* pagination */}
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '10px 14px', borderTop: '1px solid var(--line-2)' }}>
          <div className="tiny muted">Rows 1–12 of 2,847</div>
          <div className="hstack" style={{ gap: 4 }}>
            <button className="btn btn-sm" disabled><Icon name="chevronLeft" size={12}/></button>
            {[1,2,3,'…',237].map((p,i) => (
              <button key={i} className="btn btn-sm" style={{
                background: p === 1 ? 'var(--accent-soft)' : 'var(--surface)',
                color: p === 1 ? 'var(--accent)' : 'var(--ink-1)',
                borderColor: p === 1 ? 'var(--accent-border)' : 'var(--line-3)',
                minWidth: 26, padding: 0,
              }}>{p}</button>
            ))}
            <button className="btn btn-sm"><Icon name="chevronRight" size={12}/></button>
          </div>
        </div>
      </div>
    </div>
  );
};

const FilterChip = ({ label, value }) => (
  <button className="btn btn-sm" style={{ gap: 4, fontWeight: 400 }}>
    <span style={{ color: 'var(--ink-3)' }}>{label}:</span>
    <span>{value}</span>
    <Icon name="chevronDown" size={10} style={{ opacity: .5 }}/>
  </button>
);

window.CustomersScreen = CustomersScreen;
window.FilterChip = FilterChip;
