// Reseller portal screens — "ProxyVN" tenant POV
// Per spec: reseller_wallet (owed to Admin), client-facing catalog (clones w/ margin),
// their clients, their services.

const ResellerDashboard = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    {/* Wallet hero */}
    <div className="card" style={{ padding: 20 }}>
      <div style={{ display: 'grid', gridTemplateColumns: '1.2fr 1fr 1fr 1fr', gap: 24, alignItems: 'stretch' }}>
        <div style={{ borderRight: '1px solid var(--line)', paddingRight: 24 }}>
          <div className="tiny muted" style={{ letterSpacing: 0.3, textTransform: 'uppercase' }}>Reseller wallet</div>
          <div style={{ fontSize: 32, fontWeight: 600, marginTop: 4, fontFeatureSettings: '"tnum"' }}>
            $4,820<span style={{ color: 'var(--ink-3)', fontSize: 22 }}>.50</span>
          </div>
          <div className="tiny muted" style={{ marginTop: 2 }}>Owed to Admin · settles purchases &amp; renewals</div>
          <div style={{ display: 'flex', gap: 6, marginTop: 12 }}>
            <button className="btn btn-primary btn-sm"><Icon name="plus" size={12}/> Top up</button>
            <button className="btn btn-sm">Ledger</button>
          </div>
        </div>
        <div>
          <div className="tiny muted" style={{ letterSpacing: 0.3, textTransform: 'uppercase' }}>Active services (MRR)</div>
          <div style={{ fontSize: 24, fontWeight: 600, marginTop: 4, fontFeatureSettings: '"tnum"' }}>$8,240<span className="muted" style={{fontSize:14}}>/mo</span></div>
          <div className="tiny" style={{ color: 'var(--ok)', marginTop: 2 }}>+6.2% vs last month</div>
        </div>
        <div>
          <div className="tiny muted" style={{ letterSpacing: 0.3, textTransform: 'uppercase' }}>Clients</div>
          <div style={{ fontSize: 24, fontWeight: 600, marginTop: 4 }}>312</div>
          <div className="tiny muted" style={{ marginTop: 2 }}>14 new this week</div>
        </div>
        <div>
          <div className="tiny muted" style={{ letterSpacing: 0.3, textTransform: 'uppercase' }}>Margin (30d)</div>
          <div style={{ fontSize: 24, fontWeight: 600, marginTop: 4, fontFeatureSettings: '"tnum"' }}>33.4%</div>
          <div className="tiny muted" style={{ marginTop: 2 }}>$2,180 net</div>
        </div>
      </div>
    </div>

    <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: 16 }}>
      <div className="card">
        <div className="card-header">
          <h3>Recent client activity</h3>
          <button className="btn btn-ghost btn-sm">All clients →</button>
        </div>
        <table className="tbl">
          <thead>
            <tr>
              <th>Client</th><th>Services</th><th className="num">Wallet</th>
              <th className="num">Orders</th><th>Last</th><th>Status</th>
            </tr>
          </thead>
          <tbody>
            {BILLING.resellerClients.slice(0, 6).map(c => (
              <tr key={c.id}>
                <td>
                  <div style={{ fontSize: 13, fontWeight: 500 }}>{c.name}</div>
                  <div className="tiny muted">{c.email}</div>
                </td>
                <td className="num">{c.services}</td>
                <td className="num" style={{ color: c.wallet === 0 ? 'var(--danger)' : 'var(--ink-0)' }}>
                  ${c.wallet.toFixed(2)}
                </td>
                <td className="num">{c.orders}</td>
                <td className="tiny muted">{c.lastLogin}</td>
                <td><span className={`badge dot ${c.status === 'active' ? 'ok' : 'danger'}`}>{c.status}</span></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        <div className="card" style={{ padding: 16 }}>
          <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 10 }}>Revenue split — this month</div>
          <div style={{ display: 'flex', height: 6, borderRadius: 1, overflow: 'hidden', background: 'var(--line-2)' }}>
            <div style={{ width: '67%', background: 'var(--accent)' }}/>
            <div style={{ width: '33%', background: 'var(--ink-2)' }}/>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 10, marginTop: 14 }}>
            <div>
              <div className="tiny muted">Gross revenue</div>
              <div style={{ fontSize: 18, fontWeight: 600, marginTop: 2 }}>$6,520</div>
            </div>
            <div>
              <div className="tiny muted">Net margin</div>
              <div style={{ fontSize: 18, fontWeight: 600, marginTop: 2, color: 'var(--accent)' }}>$2,180</div>
            </div>
            <div>
              <div className="tiny muted">Admin cost</div>
              <div style={{ fontSize: 14, marginTop: 2 }}>$4,340</div>
            </div>
            <div>
              <div className="tiny muted">Refunds</div>
              <div style={{ fontSize: 14, marginTop: 2 }}>$48.20</div>
            </div>
          </div>
        </div>

        <div className="card" style={{ padding: 16 }}>
          <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 10 }}>Upstream status</div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            {BILLING.providers.slice(0, 4).map(p => (
              <div key={p.id} style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 12 }}>
                <span style={{
                  width: 6, height: 6, borderRadius: '50%',
                  background: p.health === 'ok' ? 'var(--ok)' : 'var(--warn)',
                }}/>
                <span style={{ flex: 1 }}>{p.name}</span>
                <span className="tiny muted">{p.failRate}% fail</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  </div>
);

const ResellerCatalogScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12 }}>
      <Kpi label="Plans published" value="9" sub="of 10 cloned"/>
      <Kpi label="Avg margin" value="31%" delta={2.4} sub="target: 30%"/>
      <Kpi label="Out of stock" value="1" sub="Large VPS HEL"/>
      <Kpi label="Pricing warnings" value="1" sub="margin negative"/>
    </div>

    <div className="card">
      <div className="card-header">
        <h3>Catalog · your pricing</h3>
        <div className="hstack">
          <span className="tiny muted">Sync with admin catalog</span>
          <button className="btn btn-ghost btn-sm">Check updates</button>
          <button className="btn btn-sm">+ Clone plan</button>
        </div>
      </div>
      <table className="tbl">
        <thead>
          <tr>
            <th>Plan</th>
            <th>Unit</th>
            <th className="num">Admin cost</th>
            <th className="num">Your price</th>
            <th className="num">Margin</th>
            <th>Stock</th>
            <th>Version</th>
            <th>Status</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {BILLING.resellerCatalog.map(p => (
            <tr key={p.plan}>
              <td style={{ fontWeight: 500 }}>{p.plan}</td>
              <td className="tiny muted">{p.unit}</td>
              <td className="num">${p.cost.toFixed(2)}</td>
              <td className="num" style={{ fontWeight: 500 }}>${p.selling.toFixed(2)}</td>
              <td className="num" style={{ color: p.margin < 0 ? 'var(--danger)' : p.margin < 20 ? 'var(--warn)' : 'var(--ok)', fontWeight: 500 }}>
                {fmtMargin(p.margin)}
              </td>
              <td>
                <span className={`badge dot ${p.stock === 'ok' ? 'ok' : p.stock === 'low' ? 'warn' : 'danger'}`}>{p.stock}</span>
              </td>
              <td className="mono tiny muted">{p.version}</td>
              <td>
                <span className={`badge ${p.status === 'disabled' ? 'danger' : p.status === 'warn' ? 'warn' : 'ok'} dot`}>{p.status}</span>
              </td>
              <td style={{ textAlign: 'right' }}>
                <button className="btn btn-ghost btn-sm">Edit price</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

const ResellerClientsScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div className="card">
      <div className="card-header">
        <h3>Your clients</h3>
        <div className="hstack">
          <input className="input" style={{ width: 220, height: 28 }} placeholder="Search email, name…"/>
          <button className="btn btn-sm"><Icon name="plus" size={12}/> Invite client</button>
        </div>
      </div>
      <table className="tbl">
        <thead>
          <tr>
            <th>Client</th>
            <th>Email</th>
            <th className="num">Wallet</th>
            <th className="num">Services</th>
            <th className="num">Orders</th>
            <th>Status</th>
            <th>Last login</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {BILLING.resellerClients.map(c => (
            <tr key={c.id}>
              <td>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <div style={{
                    width: 24, height: 24, borderRadius: '50%',
                    background: 'var(--ink-4)', color: 'var(--ink-1)',
                    fontSize: 10, fontWeight: 600,
                    display: 'grid', placeItems: 'center',
                  }}>{c.name.split(' ').map(s => s[0]).slice(0, 2).join('')}</div>
                  <div>
                    <div style={{ fontSize: 13, fontWeight: 500 }}>{c.name}</div>
                    <div className="mono tiny muted">{c.id}</div>
                  </div>
                </div>
              </td>
              <td className="mono" style={{ fontSize: 12 }}>{c.email}</td>
              <td className="num" style={{ color: c.wallet === 0 ? 'var(--danger)' : 'var(--ink-0)', fontWeight: 500 }}>
                ${c.wallet.toFixed(2)}
              </td>
              <td className="num">{c.services}</td>
              <td className="num">{c.orders}</td>
              <td><span className={`badge dot ${c.status === 'active' ? 'ok' : 'danger'}`}>{c.status}</span></td>
              <td className="tiny muted">{c.lastLogin}</td>
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

const ResellerWalletScreen = () => (
  <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div style={{ display: 'grid', gridTemplateColumns: '1.4fr 1fr', gap: 16, alignItems: 'start' }}>
      <div className="card" style={{ padding: 20 }}>
        <div className="tiny muted" style={{ letterSpacing: 0.3, textTransform: 'uppercase' }}>Reseller wallet balance</div>
        <div style={{ fontSize: 40, fontWeight: 600, marginTop: 6, fontFeatureSettings: '"tnum"' }}>
          $4,820<span style={{ color: 'var(--ink-3)', fontSize: 28 }}>.50</span>
        </div>
        <div style={{ display: 'flex', gap: 24, marginTop: 16, paddingTop: 16, borderTop: '1px solid var(--line)' }}>
          <div>
            <div className="tiny muted">This month's spend</div>
            <div style={{ fontSize: 16, fontWeight: 500, marginTop: 2 }}>$4,340.20</div>
          </div>
          <div>
            <div className="tiny muted">Avg daily burn</div>
            <div style={{ fontSize: 16, fontWeight: 500, marginTop: 2 }}>$147.50</div>
          </div>
          <div>
            <div className="tiny muted">Estimated runway</div>
            <div style={{ fontSize: 16, fontWeight: 500, marginTop: 2, color: 'var(--ok)' }}>~32 days</div>
          </div>
        </div>
      </div>

      <div className="card" style={{ padding: 20 }}>
        <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>Add funds</div>
        <div className="tiny muted" style={{ marginBottom: 12 }}>
          Top-ups are held pending verification until payment is matched. Only matched amounts credit your wallet.
        </div>
        <div style={{ display: 'flex', gap: 6, marginBottom: 12 }}>
          {['500', '1,000', '2,500', '5,000'].map(a => (
            <button key={a} className="btn btn-sm" style={{ flex: 1 }}>${a}</button>
          ))}
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 6, marginBottom: 12 }}>
          <button className="btn" style={{ justifyContent: 'flex-start', textAlign: 'left', padding: 10 }}>
            <div>
              <div style={{ fontSize: 12, fontWeight: 600 }}>VietQR</div>
              <div className="tiny muted">Verified in ~5 min</div>
            </div>
          </button>
          <button className="btn" style={{ justifyContent: 'flex-start', textAlign: 'left', padding: 10 }}>
            <div>
              <div style={{ fontSize: 12, fontWeight: 600 }}>USDT (TRC20)</div>
              <div className="tiny muted">Verified on 1 conf</div>
            </div>
          </button>
        </div>
        <button className="btn btn-primary" style={{ width: '100%' }}>Generate payment →</button>
      </div>
    </div>

    <div className="card">
      <div className="card-header">
        <h3>Wallet ledger</h3>
        <div className="hstack">
          <select className="select" style={{ width: 140, height: 28 }}>
            <option>All entries</option>
            <option>Debits</option>
            <option>Credits</option>
          </select>
          <button className="btn btn-ghost btn-sm">Export CSV</button>
        </div>
      </div>
      <table className="tbl">
        <thead>
          <tr>
            <th>Timestamp</th>
            <th>Type</th>
            <th>Reference</th>
            <th className="num">Amount</th>
            <th className="num">Balance</th>
          </tr>
        </thead>
        <tbody>
          {[
            { ts: '2026-04-22 14:02', type: 'purchase.reseller_wallet.debit', amount: -62.00, ref: 'ORD-48290 · Admin settles', balance: 4820.50 },
            { ts: '2026-04-22 09:14', type: 'renewal.reseller_wallet.debit', amount: -19.00, ref: 'svc-v-4421 · VPS Small', balance: 4882.50 },
            { ts: '2026-04-21 22:40', type: 'topup.credit.reseller', amount: 2000.00, ref: 'TUP-9114 · VietQR verified', balance: 4901.50 },
            { ts: '2026-04-21 11:08', type: 'renewal.reseller_wallet.debit', amount: -128.40, ref: 'batch: 14 residential renewals', balance: 2901.50 },
            { ts: '2026-04-20 16:30', type: 'purchase.reseller_wallet.debit', amount: -48.00, ref: 'ORD-48276 · Mobile 4G', balance: 3029.90 },
            { ts: '2026-04-20 09:02', type: 'refund.reseller_wallet.credit', amount: 62.00, ref: 'ORD-48274 · provisioning fail', balance: 3077.90 },
          ].map((r, i) => (
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

window.ResellerDashboard = ResellerDashboard;
window.ResellerCatalogScreen = ResellerCatalogScreen;
window.ResellerClientsScreen = ResellerClientsScreen;
window.ResellerWalletScreen = ResellerWalletScreen;
