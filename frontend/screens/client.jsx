// Client portal — "Linh Tran" (reseller client of ProxyVN)
// Per spec: services, client_wallet, checkout, grace-period handling

const ClientDashboard = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    {/* Grace alert */}
    <div style={{
      padding: '12px 14px', background: 'var(--warn-bg)',
      border: '1px solid var(--warn-border)', borderRadius: 4,
      display: 'flex', alignItems: 'center', gap: 12,
    }}>
      <Icon name="alert" size={16} style={{ color: 'var(--warn)' }}/>
      <div style={{ flex: 1, fontSize: 13 }}>
        <b>vps-test</b> is suspended — grace period ends in <b>2 days</b>.
        Top up your wallet or pay the renewal to restore.
      </div>
      <button className="btn btn-sm btn-primary">Renew $19</button>
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: '1.2fr 1fr 1fr 1fr', gap: 12 }}>
      <div className="card" style={{ padding: 16 }}>
        <div className="tiny muted" style={{ letterSpacing: 0.3, textTransform: 'uppercase' }}>Wallet balance</div>
        <div style={{ fontSize: 26, fontWeight: 600, marginTop: 4, fontFeatureSettings: '"tnum"' }}>
          $128<span style={{ color: 'var(--ink-3)', fontSize: 18 }}>.40</span>
        </div>
        <div className="tiny muted" style={{ marginTop: 2 }}>Auto-renews turned on · 4 services</div>
        <div style={{ display: 'flex', gap: 6, marginTop: 10 }}>
          <button className="btn btn-primary btn-sm" style={{ flex: 1 }}>Top up</button>
          <button className="btn btn-sm">Ledger</button>
        </div>
      </div>
      <Kpi label="Active services" value="5" sub="1 provisioning · 1 suspended"/>
      <Kpi label="This month" value="$62.00" sub="4 renewals paid"/>
      <Kpi label="Next 7 days" value="$48" sub="3 renewals due"/>
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: 16 }}>
      <div className="card">
        <div className="card-header">
          <h3>Your services</h3>
          <button className="btn btn-sm btn-primary">+ New order</button>
        </div>
        <table className="tbl">
          <thead>
            <tr>
              <th>Service</th>
              <th>Identifier</th>
              <th>Region</th>
              <th>Usage</th>
              <th>Renews</th>
              <th>Status</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {BILLING.clientServices.map(s => (
              <tr key={s.id}>
                <td>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    <div style={{
                      width: 22, height: 22, borderRadius: 2,
                      background: 'var(--ink-6)', border: '1px solid var(--line)',
                      display: 'grid', placeItems: 'center',
                    }}>
                      <Icon name={
                        s.type === 'vps-linux' ? 'server' :
                        s.type === 'residential' ? 'globe' :
                        s.type === 'datacenter' ? 'db' :
                        s.type === 'mobile' ? 'phone' : 'box'
                      } size={12}/>
                    </div>
                    <div>
                      <div style={{ fontSize: 13, fontWeight: 500 }}>{s.label}</div>
                      <div className="tiny muted" style={{ textTransform: 'capitalize' }}>{s.type.replace('-', ' · ')}</div>
                    </div>
                  </div>
                </td>
                <td className="mono" style={{ fontSize: 11 }}>{s.identifier}</td>
                <td className="tiny">{s.region}</td>
                <td className="tiny muted">{s.bandwidth}</td>
                <td className="tiny muted">
                  {s.expiry}
                  {s.note && <div style={{ color: 'var(--warn)', marginTop: 2 }}>{s.note}</div>}
                </td>
                <td>
                  <span className={`badge dot ${
                    s.status === 'active' ? 'ok' :
                    s.status === 'suspended' ? 'danger' :
                    s.status === 'provisioning' ? 'info' : ''
                  }`}>{s.status}</span>
                </td>
                <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                  <button className="btn btn-ghost btn-sm">Manage</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        <div className="card" style={{ padding: 16 }}>
          <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 10 }}>Residential usage · April</div>
          <div style={{ display: 'flex', alignItems: 'baseline', gap: 6 }}>
            <span style={{ fontSize: 22, fontWeight: 600 }}>5.0</span>
            <span className="muted">/ 15 GB</span>
          </div>
          <div style={{ height: 4, background: 'var(--line-2)', marginTop: 8 }}>
            <div style={{ width: '33%', height: '100%', background: 'var(--accent)' }}/>
          </div>
          <div className="tiny muted" style={{ marginTop: 6 }}>Resets May 14 · 2 pools</div>
        </div>

        <div className="card" style={{ padding: 16 }}>
          <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 10 }}>Recent activity</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
            {BILLING.clientLedger.slice(0, 5).map((l, i) => (
              <div key={i} style={{ display: 'flex', alignItems: 'flex-start', gap: 8, fontSize: 12 }}>
                <div style={{
                  width: 6, height: 6, borderRadius: '50%',
                  background: l.amount > 0 ? 'var(--ok)' : 'var(--ink-3)',
                  marginTop: 6, flexShrink: 0,
                }}/>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ fontSize: 12 }}>{l.ref}</div>
                  <div className="tiny muted">{l.ts.slice(5, 16)}</div>
                </div>
                <div className="mono" style={{ fontSize: 12, fontWeight: 500, color: l.amount > 0 ? 'var(--ok)' : 'var(--ink-1)' }}>
                  {l.amount > 0 ? '+' : ''}${Math.abs(l.amount).toFixed(2)}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  </div>
);

const ClientShopScreen = () => {
  const [tab, setTab] = React.useState('vps');
  const plans = {
    vps: [
      { name: 'Starter', specs: '1 vCPU · 2 GB · 20 GB NVMe', region: 'VN-HCM / HAN', price: 12, popular: false },
      { name: 'Small', specs: '2 vCPU · 4 GB · 60 GB NVMe', region: 'VN-HCM / HAN', price: 19, popular: true },
      { name: 'Medium', specs: '4 vCPU · 8 GB · 120 GB NVMe', region: 'VN-HCM / HAN / SG', price: 34, popular: false },
      { name: 'Large', specs: '8 vCPU · 16 GB · 240 GB NVMe', region: 'VN-HCM / HEL', price: 68, popular: false },
      { name: 'XL', specs: '16 vCPU · 32 GB · 480 GB NVMe', region: 'HEL / US', price: 129, popular: false },
      { name: 'CPU-Optimized', specs: '8 vCPU (ded.) · 16 GB · 200 GB', region: 'HEL', price: 148, popular: false },
    ],
    proxy: [
      { name: 'Residential · Std', specs: 'Global pool · 50M+ IPs · rotating', region: 'Global', price: 6.50, unit: '/GB', popular: false },
      { name: 'Residential · Prm', specs: 'Premium ASN · sticky sessions', region: 'Global', price: 9.80, unit: '/GB', popular: true },
      { name: 'Datacenter · Shared', specs: '10 IPs · unlim BW · rotating', region: 'US / EU / ASIA', price: 8, unit: '/mo', popular: false },
      { name: 'Datacenter · Dedicated', specs: '10 IPs · static · unlim BW', region: 'US / EU', price: 22, unit: '/mo', popular: false },
      { name: 'ISP Static', specs: '10 IPs · real ISP · unlim', region: 'US / UK', price: 35, unit: '/mo', popular: false },
      { name: 'Mobile 4G', specs: '1 port · rotating IMEI', region: 'VN', price: 48, unit: '/mo', popular: false },
    ],
  };
  const activePlans = plans[tab];

  return (
    <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
      <div>
        <h2 style={{ fontSize: 22, fontWeight: 600, marginBottom: 4 }}>Shop</h2>
        <div className="tiny muted">Prices set by your reseller (ProxyVN) · monthly billing · cancel anytime</div>
      </div>

      <div style={{ display: 'flex', gap: 4, borderBottom: '1px solid var(--line)' }}>
        {[['vps', 'VPS'], ['proxy', 'Proxies']].map(([k, l]) => (
          <button key={k} onClick={() => setTab(k)} style={{
            padding: '10px 18px', border: 'none', background: 'transparent',
            fontSize: 14, fontWeight: tab === k ? 600 : 400,
            color: tab === k ? 'var(--accent)' : 'var(--ink-2)',
            borderBottom: tab === k ? '2px solid var(--accent)' : '2px solid transparent',
            cursor: 'pointer', marginBottom: -1, fontFamily: 'inherit',
          }}>{l}</button>
        ))}
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 12 }}>
        {activePlans.map(p => (
          <div key={p.name} className="card" style={{
            padding: 18, position: 'relative',
            borderColor: p.popular ? 'var(--accent)' : 'var(--line)',
          }}>
            {p.popular && (
              <div style={{
                position: 'absolute', top: -8, left: 14,
                padding: '2px 8px', background: 'var(--accent)', color: '#fff',
                fontSize: 10, fontWeight: 600, letterSpacing: 0.4, textTransform: 'uppercase',
              }}>Popular</div>
            )}
            <div style={{ fontSize: 15, fontWeight: 600 }}>{p.name}</div>
            <div className="tiny muted" style={{ marginTop: 4, minHeight: 32 }}>{p.specs}</div>
            <div className="tiny muted" style={{ marginTop: 8 }}>Region: {p.region}</div>
            <div style={{ display: 'flex', alignItems: 'baseline', gap: 4, marginTop: 14, paddingTop: 14, borderTop: '1px solid var(--line)' }}>
              <span style={{ fontSize: 24, fontWeight: 600, fontFeatureSettings: '"tnum"' }}>${p.price}</span>
              <span className="tiny muted">{p.unit || '/mo'}</span>
            </div>
            <button className="btn btn-primary" style={{ width: '100%', marginTop: 12 }}>Configure →</button>
          </div>
        ))}
      </div>
    </div>
  );
};

const ClientCheckoutScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div>
      <div className="tiny muted">Step 2 of 3 · VPS Small</div>
      <h2 style={{ fontSize: 22, fontWeight: 600, marginTop: 4 }}>Configure &amp; checkout</h2>
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: '1.5fr 1fr', gap: 20, alignItems: 'start' }}>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        <div className="card" style={{ padding: 20 }}>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 14 }}>Location</div>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 8 }}>
            {['VN-HCM', 'VN-HAN', 'SG'].map((r, i) => (
              <label key={r} style={{
                padding: 14, border: '1px solid ' + (i === 0 ? 'var(--accent)' : 'var(--line)'),
                borderRadius: 3, cursor: 'pointer', background: i === 0 ? 'var(--accent-bg)' : 'transparent',
              }}>
                <input type="radio" name="region" defaultChecked={i === 0} style={{ marginRight: 6 }}/>
                <span style={{ fontSize: 13, fontWeight: 500 }}>{r}</span>
                <div className="tiny muted" style={{ marginTop: 2 }}>
                  {i === 0 ? 'TPHCM · 2ms' : i === 1 ? 'Hà Nội · 12ms' : 'Singapore · 28ms'}
                </div>
              </label>
            ))}
          </div>
        </div>

        <div className="card" style={{ padding: 20 }}>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 14 }}>Operating system</div>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 8 }}>
            {['Ubuntu 24.04', 'Debian 12', 'AlmaLinux 9', 'Windows Server'].map((os, i) => (
              <label key={os} style={{
                padding: '10px 12px', border: '1px solid ' + (i === 0 ? 'var(--accent)' : 'var(--line)'),
                borderRadius: 3, cursor: 'pointer', fontSize: 12, fontWeight: 500,
                background: i === 0 ? 'var(--accent-bg)' : 'transparent',
              }}>
                <input type="radio" name="os" defaultChecked={i === 0} style={{ marginRight: 6 }}/>
                {os}
              </label>
            ))}
          </div>
        </div>

        <div className="card" style={{ padding: 20 }}>
          <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 14 }}>Billing cycle</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            {[
              { label: 'Monthly — calendar month', price: 19.00, pro: '' },
              { label: 'Monthly — 30-day rolling', price: 19.00, pro: 'Same price, exact 30-day periods' },
              { label: 'Quarterly (save 5%)', price: 54.15, pro: 'Billed $54.15 every 3 months' },
              { label: 'Annual (save 10%)', price: 205.20, pro: 'Billed $205.20/yr' },
            ].map((c, i) => (
              <label key={i} style={{
                padding: 12, border: '1px solid ' + (i === 0 ? 'var(--accent)' : 'var(--line)'),
                borderRadius: 3, cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 10,
                background: i === 0 ? 'var(--accent-bg)' : 'transparent',
              }}>
                <input type="radio" name="cycle" defaultChecked={i === 0}/>
                <div style={{ flex: 1 }}>
                  <div style={{ fontSize: 13, fontWeight: 500 }}>{c.label}</div>
                  {c.pro && <div className="tiny muted" style={{ marginTop: 2 }}>{c.pro}</div>}
                </div>
                <div className="mono" style={{ fontSize: 13, fontWeight: 500 }}>${c.price.toFixed(2)}</div>
              </label>
            ))}
          </div>
        </div>
      </div>

      <div className="card" style={{ padding: 20, position: 'sticky', top: 80 }}>
        <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 14 }}>Order summary</div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8, fontSize: 13 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span>VPS Small · VN-HCM</span>
            <span className="mono">$19.00</span>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between' }} className="tiny muted">
            <span>2 vCPU · 4 GB · 60 GB NVMe</span>
            <span/>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between' }} className="tiny muted">
            <span>Ubuntu 24.04 · calendar month</span>
            <span/>
          </div>
          <div style={{ borderTop: '1px solid var(--line)', margin: '8px 0', paddingTop: 8, display: 'flex', justifyContent: 'space-between' }}>
            <span>Subtotal</span>
            <span className="mono">$19.00</span>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 16, fontWeight: 600 }}>
            <span>Total due today</span>
            <span className="mono">$19.00</span>
          </div>
        </div>

        <div style={{ marginTop: 18, paddingTop: 18, borderTop: '1px solid var(--line)' }}>
          <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 8 }}>Payment</div>
          <div style={{ padding: 10, background: 'var(--ok-bg)', border: '1px solid var(--ok-border)', borderRadius: 3, display: 'flex', alignItems: 'center', gap: 8 }}>
            <Icon name="check" size={14} style={{ color: 'var(--ok)' }}/>
            <div style={{ flex: 1, fontSize: 12 }}>
              <div style={{ fontWeight: 500 }}>Pay from wallet</div>
              <div className="tiny muted">Balance: $128.40 · enough for this order</div>
            </div>
            <input type="radio" defaultChecked/>
          </div>
          <label style={{ display: 'flex', alignItems: 'center', gap: 8, padding: 10, marginTop: 6, border: '1px solid var(--line)', borderRadius: 3, cursor: 'pointer' }}>
            <input type="radio" name="pay"/>
            <div style={{ flex: 1 }}>
              <div style={{ fontSize: 12, fontWeight: 500 }}>Top up &amp; pay</div>
              <div className="tiny muted">VietQR · USDT</div>
            </div>
          </label>
        </div>

        <button className="btn btn-primary" style={{ width: '100%', marginTop: 14, height: 38, fontSize: 14 }}>
          Place order · provision now
        </button>
        <div className="tiny muted" style={{ textAlign: 'center', marginTop: 8 }}>
          By placing this order you agree to the Terms of Service.
        </div>
      </div>
    </div>
  </div>
);

const ClientWalletScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div style={{ display: 'grid', gridTemplateColumns: '1.4fr 1fr', gap: 16, alignItems: 'start' }}>
      <div className="card" style={{ padding: 20 }}>
        <div className="tiny muted" style={{ letterSpacing: 0.3, textTransform: 'uppercase' }}>Wallet balance</div>
        <div style={{ fontSize: 40, fontWeight: 600, marginTop: 6, fontFeatureSettings: '"tnum"' }}>
          $128<span style={{ color: 'var(--ink-3)', fontSize: 28 }}>.40</span>
        </div>
        <div style={{ display: 'flex', gap: 24, marginTop: 16, paddingTop: 16, borderTop: '1px solid var(--line)' }}>
          <div>
            <div className="tiny muted">This month</div>
            <div style={{ fontSize: 16, fontWeight: 500, marginTop: 2 }}>-$62.00</div>
          </div>
          <div>
            <div className="tiny muted">Auto-renew reserve</div>
            <div style={{ fontSize: 16, fontWeight: 500, marginTop: 2 }}>$48.00</div>
          </div>
          <div>
            <div className="tiny muted">Available</div>
            <div style={{ fontSize: 16, fontWeight: 500, marginTop: 2, color: 'var(--ok)' }}>$80.40</div>
          </div>
        </div>
      </div>

      <div className="card" style={{ padding: 20 }}>
        <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>Add funds</div>
        <div style={{ display: 'flex', gap: 6, marginBottom: 12 }}>
          {['50', '100', '200', '500'].map(a => (
            <button key={a} className="btn btn-sm" style={{ flex: 1 }}>${a}</button>
          ))}
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 6, marginBottom: 12 }}>
          <button className="btn" style={{ padding: 10, textAlign: 'left', justifyContent: 'flex-start' }}>
            <div>
              <div style={{ fontSize: 12, fontWeight: 600 }}>VietQR</div>
              <div className="tiny muted">~5 min</div>
            </div>
          </button>
          <button className="btn" style={{ padding: 10, textAlign: 'left', justifyContent: 'flex-start' }}>
            <div>
              <div style={{ fontSize: 12, fontWeight: 600 }}>USDT</div>
              <div className="tiny muted">TRC20</div>
            </div>
          </button>
        </div>
        <button className="btn btn-primary" style={{ width: '100%' }}>Generate payment →</button>
      </div>
    </div>

    <div className="card">
      <div className="card-header">
        <h3>Wallet ledger</h3>
        <button className="btn btn-ghost btn-sm">Export</button>
      </div>
      <table className="tbl">
        <thead>
          <tr>
            <th>Timestamp</th><th>Type</th><th>Reference</th>
            <th className="num">Amount</th><th className="num">Balance</th>
          </tr>
        </thead>
        <tbody>
          {BILLING.clientLedger.map((r, i) => (
            <tr key={i}>
              <td className="mono tiny muted">{r.ts}</td>
              <td className="mono" style={{ fontSize: 11, color: r.amount > 0 ? 'var(--ok)' : 'var(--ink-1)' }}>{r.type}</td>
              <td style={{ fontSize: 12 }}>{r.ref}</td>
              <td className="num" style={{ color: r.amount > 0 ? 'var(--ok)' : 'var(--ink-0)', fontWeight: 500 }}>
                {r.amount > 0 ? '+' : ''}${Math.abs(r.amount).toFixed(2)}
              </td>
              <td className="num mono tiny">{r.balance.toFixed(2)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

window.ClientDashboard = ClientDashboard;
window.ClientShopScreen = ClientShopScreen;
window.ClientCheckoutScreen = ClientCheckoutScreen;
window.ClientWalletScreen = ClientWalletScreen;
