// Services (Proxies + VPS) with tabs

const TypeIcon = ({ type }) => {
  const map = {
    'residential': { icon: 'globe', color: 'var(--accent)', label: 'Residential' },
    'datacenter': { icon: 'server', color: 'var(--ink-1)', label: 'Datacenter' },
    'mobile': { icon: 'globe', color: 'var(--info)', label: 'Mobile' },
    'isp': { icon: 'globe', color: 'var(--warn)', label: 'ISP' },
    'vps-linux': { icon: 'cpu', color: 'var(--ok)', label: 'VPS Linux' },
    'vps-win': { icon: 'cpu', color: '#7c3aed', label: 'VPS Windows' },
  };
  const m = map[type] || map.residential;
  return (
    <div className="hstack" style={{ gap: 6 }}>
      <div style={{ width: 22, height: 22, borderRadius: 3, background: 'var(--bg-alt)', color: m.color, display: 'grid', placeItems: 'center' }}>
        <Icon name={m.icon} size={13}/>
      </div>
      <span style={{ fontSize: 11, color: 'var(--ink-3)' }}>{m.label}</span>
    </div>
  );
};

const ServicesScreen = () => {
  const [tab, setTab] = React.useState('all');
  const filtered = tab === 'all' ? SAMPLE.services
    : tab === 'proxies' ? SAMPLE.services.filter(s => !s.type.startsWith('vps'))
    : SAMPLE.services.filter(s => s.type.startsWith('vps'));

  return (
    <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 12 }}>
      {/* Summary tiles */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(5,1fr)', gap: 12 }}>
        <Kpi label="Total services" value="15,604" sub="active" spark={<Sparkline data={[15400,15420,15510,15490,15580,15620,15604]} w={56} h={20}/>}/>
        <Kpi label="Proxies" value="12,480" sub="4 types" delta={2.1}/>
        <Kpi label="VPS" value="3,124" sub="Linux + Windows" delta={0.8}/>
        <Kpi label="Provisioning" value="14" sub="in progress" spark={<span style={{fontSize:10,color:'var(--info)'}}>●</span>}/>
        <Kpi label="Overdue renewals" value="23" sub="action needed" deltaType="positive" delta="+6"/>
      </div>

      {/* Tabs + toolbar */}
      <div className="card" style={{ padding: '8px 12px', display: 'flex', alignItems: 'center', gap: 8 }}>
        <div style={{ display: 'flex', background: 'var(--bg-alt)', borderRadius: 3, padding: 2 }}>
          {['all','proxies','vps'].map(t => (
            <button key={t} onClick={() => setTab(t)} style={{
              padding: '4px 12px',
              background: tab === t ? 'var(--surface)' : 'transparent',
              boxShadow: tab === t ? 'var(--shadow-1)' : 'none',
              border: 'none', borderRadius: 2,
              fontSize: 12, fontWeight: tab === t ? 500 : 400,
              color: tab === t ? 'var(--ink-0)' : 'var(--ink-2)',
              cursor: 'pointer',
              fontFamily: 'inherit',
            }}>
              {t === 'all' ? 'All services' : t === 'proxies' ? 'Proxies' : 'VPS & Servers'}
            </button>
          ))}
        </div>
        <span style={{ width: 1, height: 20, background: 'var(--line)' }}/>
        <FilterChip label="Status" value="Active"/>
        <FilterChip label="Region" value="Any"/>
        <FilterChip label="Customer" value="Any"/>
        <div style={{ marginLeft: 'auto', display: 'flex', gap: 6 }}>
          <button className="btn btn-sm"><Icon name="download" size={12}/> Export</button>
          <button className="btn btn-primary btn-sm"><Icon name="plus" size={12}/> Provision service</button>
        </div>
      </div>

      <div className="card" style={{ overflow: 'hidden' }}>
        <table className="tbl">
          <thead>
            <tr>
              <th>Service</th>
              <th>Type</th>
              <th>Customer</th>
              <th>Region</th>
              <th className="num">Bandwidth</th>
              <th className="num">MRR</th>
              <th>Renews</th>
              <th>Status</th>
              <th style={{ width: 96 }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map(s => (
              <tr key={s.id}>
                <td>
                  <div style={{ fontWeight: 500 }}>{s.label}</div>
                  <div className="tiny mono muted">{s.id}</div>
                </td>
                <td><TypeIcon type={s.type}/></td>
                <td>{s.customer}</td>
                <td><span className="mono tiny" style={{ padding: '1px 5px', background: 'var(--bg-alt)', borderRadius: 2 }}>{s.region}</span></td>
                <td className="num">{s.bandwidth}</td>
                <td className="num" style={{ fontWeight: 500 }}>${s.price}/mo</td>
                <td>
                  {s.renewsIn > 0
                    ? <span style={{ color: s.renewsIn < 7 ? 'var(--warn)' : 'var(--ink-2)' }}>in {s.renewsIn}d</span>
                    : <span style={{ color: 'var(--danger)' }}>overdue {Math.abs(s.renewsIn)}d</span>}
                </td>
                <td><span className={`badge dot ${STATUS_BADGE[s.status]}`}>{STATUS_LABEL[s.status]}</span></td>
                <td>
                  <div className="hstack" style={{ gap: 2 }}>
                    <button className="btn btn-ghost btn-sm" style={{ padding: 4 }} title="Restart"><Icon name="restart" size={13}/></button>
                    <button className="btn btn-ghost btn-sm" style={{ padding: 4 }} title="Power"><Icon name="power" size={13}/></button>
                    <button className="btn btn-ghost btn-sm" style={{ padding: 4 }} title="Open"><Icon name="external" size={13}/></button>
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
};

window.ServicesScreen = ServicesScreen;
window.TypeIcon = TypeIcon;
