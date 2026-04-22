// Admin screens aligned with Billing-V2 spec:
// - Tenants (resellers + admin retail)
// - Provisioning queue (queued / provisioning / failed / manual_review)
// - Top-up verification
// - Provider health

const AdminTenantsScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12 }}>
      <Kpi label="Tenants" value="5" sub="1 admin · 4 resellers"/>
      <Kpi label="Retail clients (Admin)" value="1,284" delta={4.2} sub="direct-buy"/>
      <Kpi label="Reseller clients" value="632" delta={8.1} sub="under 4 resellers"/>
      <Kpi label="Low-balance resellers" value="2" delta={null} sub="below threshold"
        spark={<span className="badge danger dot">alert</span>}/>
    </div>

    <div className="card">
      <div className="card-header">
        <h3>Tenants</h3>
        <div className="hstack">
          <input className="input" style={{ width: 220, height: 28 }} placeholder="Search tenant, domain…"/>
          <button className="btn btn-sm"><Icon name="plus" size={12}/> New reseller</button>
        </div>
      </div>
      <table className="tbl">
        <thead>
          <tr>
            <th>Tenant</th>
            <th>Type</th>
            <th>Domain</th>
            <th className="num">Clients</th>
            <th className="num">Services</th>
            <th className="num">Wallet (USD)</th>
            <th>Status</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {BILLING.tenants.map(t => (
            <tr key={t.id}>
              <td>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <div style={{
                    width: 24, height: 24, borderRadius: 3,
                    background: t.type === 'admin' ? 'var(--accent)' : '#1F2937',
                    color: '#fff', fontSize: 10, fontWeight: 700,
                    display: 'grid', placeItems: 'center',
                  }}>{t.name[0]}</div>
                  <div>
                    <div style={{ fontSize: 13, fontWeight: 500 }}>{t.name}</div>
                    <div className="mono tiny muted">{t.id}</div>
                  </div>
                </div>
              </td>
              <td>
                <span className="badge" style={{ textTransform: 'capitalize' }}>{t.type}</span>
              </td>
              <td className="mono" style={{ fontSize: 12 }}>{t.domain}</td>
              <td className="num">{t.clients.toLocaleString()}</td>
              <td className="num">{t.services.toLocaleString()}</td>
              <td className="num">
                <span style={{ fontWeight: 500, color: t.walletLow ? 'var(--danger)' : 'var(--ink-0)' }}>
                  {t.wallet > 0 ? '$' + t.wallet.toLocaleString(undefined, {minimumFractionDigits:2, maximumFractionDigits:2}) : '—'}
                </span>
                {t.walletLow && <span className="badge danger" style={{ marginLeft: 6 }}>low</span>}
              </td>
              <td>
                <span className={`badge dot ${t.status === 'active' ? 'ok' : ''}`}>{t.status}</span>
              </td>
              <td style={{ textAlign: 'right' }}>
                <button className="btn btn-ghost btn-sm">Open</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

const AdminProvisioningScreen = () => {
  const [filter, setFilter] = React.useState('all');
  const counts = {
    all: BILLING.provJobs.length,
    queued: BILLING.provJobs.filter(j => j.status === 'queued').length,
    provisioning: BILLING.provJobs.filter(j => j.status === 'provisioning').length,
    failed: BILLING.provJobs.filter(j => j.status === 'failed').length,
    manual_review: BILLING.provJobs.filter(j => j.status === 'manual_review').length,
  };
  const filtered = filter === 'all' ? BILLING.provJobs : BILLING.provJobs.filter(j => j.status === filter);

  return (
    <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
      <div style={{
        display: 'flex', alignItems: 'flex-start', gap: 12,
        padding: '12px 14px', background: 'var(--warn-bg)',
        border: '1px solid var(--warn-border)', borderRadius: 4,
      }}>
        <Icon name="alert" size={16} style={{ color: 'var(--warn)', marginTop: 2 }}/>
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 13, fontWeight: 600, color: 'var(--warn)' }}>Do not retry blindly.</div>
          <div style={{ fontSize: 12, color: 'var(--ink-2)', marginTop: 2 }}>
            Provider timeout may mean partial success. Jobs in <b>manual_review</b> need human verification before retry — otherwise duplicate resources will be created.
          </div>
        </div>
      </div>

      <div style={{ display: 'flex', gap: 4, borderBottom: '1px solid var(--line)' }}>
        {[
          ['all', 'All'], ['queued', 'Queued'], ['provisioning', 'Running'],
          ['failed', 'Failed'], ['manual_review', 'Manual review'],
        ].map(([k, l]) => (
          <button key={k} onClick={() => setFilter(k)} style={{
            padding: '8px 14px', border: 'none', background: 'transparent',
            fontSize: 13, fontWeight: filter === k ? 600 : 400,
            color: filter === k ? 'var(--accent)' : 'var(--ink-2)',
            borderBottom: filter === k ? '2px solid var(--accent)' : '2px solid transparent',
            cursor: 'pointer', marginBottom: -1, fontFamily: 'inherit',
          }}>
            {l} <span className="muted" style={{ marginLeft: 4 }}>{counts[k]}</span>
          </button>
        ))}
      </div>

      <div className="card">
        <table className="tbl">
          <thead>
            <tr>
              <th>Job</th><th>Order</th><th>Service</th><th>Tenant</th>
              <th>Provider</th><th>Attempt</th><th>Error</th>
              <th>Status</th><th>Age</th><th></th>
            </tr>
          </thead>
          <tbody>
            {filtered.map(j => (
              <tr key={j.id}>
                <td className="mono" style={{ color: 'var(--accent)' }}>{j.id}</td>
                <td className="mono" style={{ fontSize: 12 }}>{j.order}</td>
                <td>{j.service}</td>
                <td>{j.tenant}</td>
                <td className="tiny muted">{j.provider}</td>
                <td className="num">{j.attempt}</td>
                <td className="tiny" style={{ color: j.error ? 'var(--danger)' : 'var(--ink-3)' }}>
                  {j.error || '—'}
                </td>
                <td>
                  <span className={`badge dot ${
                    j.status === 'queued' ? 'info' :
                    j.status === 'provisioning' ? 'info' :
                    j.status === 'failed' ? 'danger' :
                    j.status === 'manual_review' ? 'warn' : ''
                  }`}>{j.status.replace('_', ' ')}</span>
                </td>
                <td className="tiny muted">{j.age}</td>
                <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                  {j.status === 'manual_review' && (
                    <>
                      <button className="btn btn-sm" style={{ marginRight: 4 }}>Investigate</button>
                      <button className="btn btn-ghost btn-sm">Refund</button>
                    </>
                  )}
                  {j.status === 'failed' && (
                    <>
                      <button className="btn btn-sm" style={{ marginRight: 4 }}>Retry safe</button>
                      <button className="btn btn-ghost btn-sm">Refund</button>
                    </>
                  )}
                  {(j.status === 'queued' || j.status === 'provisioning') && (
                    <button className="btn btn-ghost btn-sm">View log</button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

const AdminTopupsScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12 }}>
      <Kpi label="Pending verification" value="3" sub="> 30m: 2"/>
      <Kpi label="Approved today" value="14" delta={12.0} sub="$4,820.50"/>
      <Kpi label="Rejected today" value="2" sub="not found: 2"/>
      <Kpi label="Reseller top-ups" value="$3,500" sub="2 requests"/>
    </div>

    <div className="card">
      <div className="card-header">
        <h3>Top-up queue · wallet settlement</h3>
        <div className="hstack">
          <select className="select" style={{ width: 140, height: 28 }}>
            <option>All methods</option>
            <option>VietQR</option>
            <option>USDT</option>
          </select>
          <select className="select" style={{ width: 140, height: 28 }}>
            <option>All actors</option>
            <option>Reseller wallet</option>
            <option>Client wallet</option>
          </select>
        </div>
      </div>
      <table className="tbl">
        <thead>
          <tr>
            <th>Request</th>
            <th>Tenant / actor</th>
            <th>Wallet</th>
            <th>Method</th>
            <th>Reference</th>
            <th className="num">Amount</th>
            <th>Status</th>
            <th>Created</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {BILLING.topups.map(t => (
            <tr key={t.id}>
              <td className="mono" style={{ color: 'var(--accent)' }}>{t.id}</td>
              <td>{t.tenant}</td>
              <td>
                <span className="badge" style={{ background: t.actor === 'reseller_wallet' ? 'var(--info-bg)' : 'var(--muted-bg)', color: t.actor === 'reseller_wallet' ? 'var(--info)' : 'var(--ink-2)', borderColor: t.actor === 'reseller_wallet' ? 'var(--info-border)' : 'transparent' }}>
                  {t.actor === 'reseller_wallet' ? 'Reseller' : 'Client'}
                </span>
              </td>
              <td>{t.method}</td>
              <td className="mono tiny">{t.ref}</td>
              <td className="num" style={{ fontWeight: 500 }}>${t.amount.toLocaleString()}</td>
              <td>
                <span className={`badge dot ${
                  t.status === 'pending_verification' ? 'warn' :
                  t.status === 'approved' ? 'ok' : 'danger'
                }`}>{t.status.replace('_', ' ')}</span>
                {t.reason && <div className="tiny muted" style={{ marginTop: 2 }}>{t.reason}</div>}
              </td>
              <td className="tiny muted">{t.created}</td>
              <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                {t.status === 'pending_verification' && (
                  <>
                    <button className="btn btn-sm btn-primary" style={{ marginRight: 4 }}>Approve</button>
                    <button className="btn btn-ghost btn-sm">Reject</button>
                  </>
                )}
                {t.status !== 'pending_verification' && <button className="btn btn-ghost btn-sm">View</button>}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

const AdminProvidersScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div className="card">
      <div className="card-header">
        <h3>Provider / source health</h3>
        <span className="tiny muted">Updated every 60s · no auto-failover</span>
      </div>
      <table className="tbl">
        <thead>
          <tr>
            <th>Provider</th><th>Type</th><th>Health</th>
            <th>Capacity</th><th>Fail rate (7d)</th><th>Last sync</th><th></th>
          </tr>
        </thead>
        <tbody>
          {BILLING.providers.map(p => (
            <tr key={p.id}>
              <td>
                <div style={{ fontSize: 13, fontWeight: 500 }}>{p.name}</div>
                <div className="mono tiny muted">{p.id}</div>
              </td>
              <td>
                <span className="badge" style={{ textTransform: 'capitalize' }}>{p.type.replace('-', ' ')}</span>
              </td>
              <td>
                <span className={`badge dot ${p.health === 'ok' ? 'ok' : 'warn'}`}>{p.health}</span>
              </td>
              <td style={{ width: 160 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <div style={{ flex: 1, height: 4, background: 'var(--line-2)' }}>
                    <div style={{ width: p.capacity + '%', height: '100%', background: p.capacity > 80 ? 'var(--danger)' : p.capacity > 60 ? 'var(--warn)' : 'var(--ok)' }}/>
                  </div>
                  <span className="tiny" style={{ width: 30, textAlign: 'right' }}>{p.capacity}%</span>
                </div>
              </td>
              <td className="num" style={{ color: p.failRate > 1 ? 'var(--danger)' : 'var(--ink-1)' }}>{p.failRate}%</td>
              <td className="tiny muted">{p.lastSync}</td>
              <td style={{ textAlign: 'right' }}>
                <button className="btn btn-ghost btn-sm">Test</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

window.AdminTenantsScreen = AdminTenantsScreen;
window.AdminProvisioningScreen = AdminProvisioningScreen;
window.AdminTopupsScreen = AdminTopupsScreen;
window.AdminProvidersScreen = AdminProvidersScreen;
