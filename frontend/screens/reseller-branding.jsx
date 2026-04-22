// Reseller branding settings — domain, theme, storefront, email
// The reseller configures their own white-labeled storefront.

const ResellerBrandingScreen = () => {
  const [tab, setTab] = React.useState('brand');

  return (
    <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
      <div>
        <h2 style={{ fontSize: 20, fontWeight: 600, margin: 0 }}>Branding &amp; Settings</h2>
        <div className="tiny muted" style={{ marginTop: 4 }}>
          Configure how your storefront at <span className="mono">proxyvn.io</span> looks to your clients.
        </div>
      </div>

      <div style={{ display: 'flex', gap: 4, borderBottom: '1px solid var(--line)' }}>
        {[
          ['brand', 'Brand identity'],
          ['domain', 'Domain &amp; DNS'],
          ['theme', 'Storefront theme'],
          ['email', 'Email &amp; sender'],
          ['legal', 'Legal / Terms'],
          ['team', 'Team &amp; roles'],
        ].map(([k, l]) => (
          <button key={k} onClick={() => setTab(k)} style={{
            padding: '10px 14px', border: 'none', background: 'transparent',
            fontSize: 13, fontWeight: tab === k ? 600 : 400,
            color: tab === k ? 'var(--accent)' : 'var(--ink-2)',
            borderBottom: tab === k ? '2px solid var(--accent)' : '2px solid transparent',
            cursor: 'pointer', marginBottom: -1, fontFamily: 'inherit',
          }} dangerouslySetInnerHTML={{ __html: l }}/>
        ))}
      </div>

      {tab === 'brand' && <BrandIdentityPanel/>}
      {tab === 'domain' && <DomainDnsPanel/>}
      {tab === 'theme' && <StorefrontThemePanel/>}
      {tab === 'email' && <EmailSenderPanel/>}
      {tab === 'legal' && <LegalPanel/>}
      {tab === 'team' && <TeamPanel/>}
    </div>
  );
};

// ─── Brand identity ──────────────────────────────────────────
const BrandIdentityPanel = () => (
  <div style={{ display: 'grid', gridTemplateColumns: '1.2fr 1fr', gap: 20, alignItems: 'start' }}>
    <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
      <div className="card">
        <div className="card-header"><h3>Company identity</h3></div>
        <div style={{ padding: 16, display: 'flex', flexDirection: 'column', gap: 14 }}>
          <FormRow label="Company name" hint="Shown on storefront, invoices, emails.">
            <input className="input" defaultValue="ProxyVN"/>
          </FormRow>
          <FormRow label="Legal entity" hint="For invoices. Leave blank to use company name.">
            <input className="input" defaultValue="CÔNG TY TNHH PROXYVN"/>
          </FormRow>
          <FormRow label="Tax ID (MST)">
            <input className="input" defaultValue="0315xxx9xx"/>
          </FormRow>
          <FormRow label="Support email">
            <input className="input" defaultValue="support@proxyvn.io"/>
          </FormRow>
          <FormRow label="Tagline" hint="One line, displayed under the logo.">
            <input className="input" defaultValue="High-quality proxies &amp; VPS for Vietnam &amp; APAC"/>
          </FormRow>
        </div>
      </div>

      <div className="card">
        <div className="card-header"><h3>Logo assets</h3></div>
        <div style={{ padding: 16, display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
          <LogoUploader label="Primary logo · light bg" note="SVG or PNG · 240×60"/>
          <LogoUploader label="Logo mark · compact" note="Square · 64×64" square/>
          <LogoUploader label="Logo · dark bg" note="For footer &amp; emails"/>
          <LogoUploader label="Favicon" note="32×32 PNG or ICO" square small/>
        </div>
      </div>
    </div>

    <div style={{ position: 'sticky', top: 80, display: 'flex', flexDirection: 'column', gap: 12 }}>
      <div className="tiny muted" style={{ textTransform: 'uppercase', letterSpacing: 0.5, fontWeight: 600 }}>Preview</div>
      <div className="card" style={{ padding: 20, background: 'var(--surface)' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 16 }}>
          <div style={{
            width: 34, height: 34, borderRadius: 4, background: 'var(--accent)',
            color: '#fff', display: 'grid', placeItems: 'center',
            fontWeight: 700, fontSize: 15, letterSpacing: -0.5,
          }}>P</div>
          <div>
            <div style={{ fontSize: 15, fontWeight: 700, letterSpacing: -0.2 }}>ProxyVN</div>
            <div className="tiny muted">High-quality proxies &amp; VPS for Vietnam &amp; APAC</div>
          </div>
        </div>
        <div style={{ padding: 14, background: 'var(--bg)', borderRadius: 3, fontSize: 12 }}>
          <div style={{ fontWeight: 600 }}>Invoice #INV-0421</div>
          <div className="tiny muted" style={{ marginTop: 2 }}>CÔNG TY TNHH PROXYVN · MST 0315xxx9xx</div>
          <div className="tiny muted">support@proxyvn.io</div>
        </div>
      </div>
    </div>

    <div style={{ gridColumn: '1 / -1', display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
      <button className="btn btn-sm">Discard</button>
      <button className="btn btn-sm btn-primary">Save changes</button>
    </div>
  </div>
);

// ─── Domain & DNS ────────────────────────────────────────────
const DomainDnsPanel = () => (
  <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div className="card">
      <div className="card-header">
        <h3>Custom domain</h3>
        <span className="badge ok dot">verified · TLS active</span>
      </div>
      <div style={{ padding: 16 }}>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr auto', gap: 10, alignItems: 'center' }}>
          <input className="input" defaultValue="proxyvn.io" style={{ fontSize: 14, height: 34, fontFamily: 'var(--font-mono)' }}/>
          <button className="btn">Replace domain</button>
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12, marginTop: 16 }}>
          <DomainCheck label="Apex A record" value="76.76.21.21" ok/>
          <DomainCheck label="www CNAME" value="tenants.hanetwork.vn" ok/>
          <DomainCheck label="TLS certificate" value="Let's Encrypt · renews May 14" ok/>
          <DomainCheck label="HSTS" value="max-age=31536000" ok/>
        </div>
      </div>
    </div>

    <div className="card">
      <div className="card-header"><h3>Required DNS records</h3></div>
      <table className="tbl">
        <thead>
          <tr><th>Type</th><th>Host</th><th>Value</th><th>TTL</th><th>Status</th></tr>
        </thead>
        <tbody>
          <tr>
            <td><span className="badge">A</span></td>
            <td className="mono">@</td>
            <td className="mono">76.76.21.21</td>
            <td>3600</td>
            <td><span className="badge ok dot">propagated</span></td>
          </tr>
          <tr>
            <td><span className="badge">CNAME</span></td>
            <td className="mono">www</td>
            <td className="mono">tenants.hanetwork.vn</td>
            <td>3600</td>
            <td><span className="badge ok dot">propagated</span></td>
          </tr>
          <tr>
            <td><span className="badge">CNAME</span></td>
            <td className="mono">gw.pr</td>
            <td className="mono">proxy-gw.hanetwork.vn</td>
            <td>300</td>
            <td><span className="badge ok dot">propagated</span></td>
          </tr>
          <tr>
            <td><span className="badge">TXT</span></td>
            <td className="mono">_hanet-verify</td>
            <td className="mono tiny">hanet-verify=9a21f8...</td>
            <td>3600</td>
            <td><span className="badge ok dot">verified</span></td>
          </tr>
          <tr>
            <td><span className="badge">TXT</span></td>
            <td className="mono">@</td>
            <td className="mono tiny">v=spf1 include:_spf.hanet.email -all</td>
            <td>3600</td>
            <td><span className="badge warn dot">pending</span></td>
          </tr>
        </tbody>
      </table>
    </div>

    <div className="card">
      <div className="card-header">
        <h3>Alternative domains</h3>
        <button className="btn btn-sm"><Icon name="plus" size={11}/> Add domain</button>
      </div>
      <table className="tbl">
        <thead><tr><th>Domain</th><th>Role</th><th>Status</th><th></th></tr></thead>
        <tbody>
          <tr>
            <td className="mono">proxy.proxyvn.io</td><td>Storefront</td>
            <td><span className="badge ok dot">active</span></td>
            <td style={{ textAlign: 'right' }}><button className="btn btn-ghost btn-sm">Remove</button></td>
          </tr>
          <tr>
            <td className="mono">app.proxyvn.io</td><td>Client portal</td>
            <td><span className="badge ok dot">active</span></td>
            <td style={{ textAlign: 'right' }}><button className="btn btn-ghost btn-sm">Remove</button></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
);

// ─── Storefront theme ────────────────────────────────────────
const StorefrontThemePanel = () => {
  const [primary, setPrimary] = React.useState('#1E4FA3');
  const [radius, setRadius] = React.useState(3);
  const [font, setFont] = React.useState('Inter');
  const [mode, setMode] = React.useState('light');

  const SWATCHES = ['#D50C2D', '#1E4FA3', '#0B7A3B', '#7C3AED', '#EA580C', '#0E1116'];

  return (
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1.3fr', gap: 20, alignItems: 'start' }}>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        <div className="card">
          <div className="card-header"><h3>Colors</h3></div>
          <div style={{ padding: 16, display: 'flex', flexDirection: 'column', gap: 14 }}>
            <FormRow label="Primary color" hint="Used for buttons, links, accents.">
              <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                <div style={{ display: 'flex', gap: 4 }}>
                  {SWATCHES.map(c => (
                    <button key={c} onClick={() => setPrimary(c)} style={{
                      width: 28, height: 28, borderRadius: 3, background: c,
                      border: primary === c ? '2px solid var(--ink-0)' : '1px solid var(--line-3)',
                      cursor: 'pointer', padding: 0,
                    }}/>
                  ))}
                </div>
                <input className="input mono" value={primary} onChange={e => setPrimary(e.target.value)} style={{ width: 100, height: 28 }}/>
              </div>
            </FormRow>
            <FormRow label="Background mode">
              <div style={{ display: 'flex', gap: 4, background: 'var(--bg-alt)', padding: 2, borderRadius: 3, width: 200 }}>
                {['light', 'dark'].map(m => (
                  <button key={m} onClick={() => setMode(m)} style={{
                    flex: 1, padding: '5px 8px', fontSize: 12,
                    background: mode === m ? '#fff' : 'transparent', border: 'none', borderRadius: 2,
                    fontWeight: mode === m ? 500 : 400, cursor: 'pointer', textTransform: 'capitalize', fontFamily: 'inherit',
                  }}>{m}</button>
                ))}
              </div>
            </FormRow>
          </div>
        </div>

        <div className="card">
          <div className="card-header"><h3>Typography &amp; shape</h3></div>
          <div style={{ padding: 16, display: 'flex', flexDirection: 'column', gap: 14 }}>
            <FormRow label="Font family">
              <select className="select" value={font} onChange={e => setFont(e.target.value)}>
                <option>Inter</option><option>IBM Plex Sans</option><option>Geist</option><option>Manrope</option><option>System</option>
              </select>
            </FormRow>
            <FormRow label="Border radius" hint={`${radius}px · ${radius < 4 ? 'sharp' : radius < 10 ? 'soft' : 'rounded'}`}>
              <input type="range" min="0" max="16" value={radius} onChange={e => setRadius(+e.target.value)} style={{ width: '100%' }}/>
            </FormRow>
            <FormRow label="Hero style">
              <select className="select"><option>Pricing-first (Hetzner-like)</option><option>Story-led</option><option>Minimal form</option></select>
            </FormRow>
          </div>
        </div>
      </div>

      <div style={{ position: 'sticky', top: 80 }}>
        <div className="tiny muted" style={{ textTransform: 'uppercase', letterSpacing: 0.5, fontWeight: 600, marginBottom: 10 }}>
          Live preview · proxyvn.io
        </div>
        <div style={{
          border: '1px solid var(--line)', borderRadius: 4, overflow: 'hidden',
          background: mode === 'dark' ? '#0E1116' : '#fff',
          color: mode === 'dark' ? '#E4E6EA' : 'var(--ink-0)',
          fontFamily: font === 'System' ? '-apple-system, system-ui, sans-serif' : `"${font}", sans-serif`,
        }}>
          {/* Fake nav */}
          <div style={{
            padding: '12px 18px', borderBottom: '1px solid ' + (mode === 'dark' ? '#1F2937' : 'var(--line)'),
            display: 'flex', alignItems: 'center', gap: 14, fontSize: 13,
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginRight: 'auto' }}>
              <div style={{ width: 22, height: 22, borderRadius: radius + 'px', background: primary, color: '#fff', display: 'grid', placeItems: 'center', fontWeight: 700, fontSize: 12 }}>P</div>
              <span style={{ fontWeight: 700 }}>ProxyVN</span>
            </div>
            <span style={{ opacity: .7 }}>Proxies</span>
            <span style={{ opacity: .7 }}>VPS</span>
            <span style={{ opacity: .7 }}>Pricing</span>
            <button style={{
              border: '1px solid ' + (mode === 'dark' ? '#2A2F36' : 'var(--line-3)'), background: 'transparent',
              color: 'inherit', borderRadius: radius + 'px', padding: '5px 10px', fontSize: 12, cursor: 'pointer',
            }}>Log in</button>
            <button style={{
              border: 'none', background: primary, color: '#fff',
              borderRadius: radius + 'px', padding: '5px 12px', fontSize: 12, fontWeight: 500, cursor: 'pointer',
            }}>Sign up</button>
          </div>
          <div style={{ padding: 28 }}>
            <div style={{ fontSize: 26, fontWeight: 700, letterSpacing: -0.4, marginBottom: 6 }}>
              Proxies &amp; VPS, built for APAC
            </div>
            <div style={{ fontSize: 13, opacity: .7, marginBottom: 20, maxWidth: 440 }}>
              Low-latency infrastructure with wallet-based billing. Start small, scale when you need to.
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 10 }}>
              {[
                { name: 'Residential', price: '6.50', unit: '/GB' },
                { name: 'VPS Small', price: '19', unit: '/mo', popular: true },
                { name: 'Datacenter', price: '8', unit: '/mo' },
              ].map(p => (
                <div key={p.name} style={{
                  padding: 14, borderRadius: radius + 'px',
                  border: '1px solid ' + (p.popular ? primary : (mode === 'dark' ? '#1F2937' : 'var(--line)')),
                  background: p.popular ? (mode === 'dark' ? '#1a1f2e' : '#fafbff') : 'transparent',
                }}>
                  <div style={{ fontSize: 12, fontWeight: 600 }}>{p.name}</div>
                  <div style={{ marginTop: 8, fontSize: 20, fontWeight: 700 }}>${p.price}<span style={{ fontSize: 11, opacity: .6, fontWeight: 400 }}>{p.unit}</span></div>
                  <button style={{
                    marginTop: 10, width: '100%', border: 'none',
                    background: p.popular ? primary : 'transparent',
                    color: p.popular ? '#fff' : primary,
                    borderRadius: radius + 'px', padding: '6px 10px', fontSize: 11, fontWeight: 500, cursor: 'pointer',
                    border: p.popular ? 'none' : `1px solid ${primary}`,
                  }}>Configure</button>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div style={{ gridColumn: '1 / -1', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div className="tiny muted">Changes apply to storefront &amp; client portal within ~30s after save.</div>
        <div style={{ display: 'flex', gap: 8 }}>
          <button className="btn btn-sm">Preview in new tab</button>
          <button className="btn btn-sm btn-primary">Publish theme</button>
        </div>
      </div>
    </div>
  );
};

// ─── Email & sender ──────────────────────────────────────────
const EmailSenderPanel = () => (
  <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
    <div className="card">
      <div className="card-header">
        <h3>Sending domain</h3>
        <span className="badge warn dot">1 DNS record pending</span>
      </div>
      <div style={{ padding: 16 }}>
        <FormRow label="From address" hint="Used for order confirmations, invoices, provisioning notices.">
          <input className="input" defaultValue="billing@proxyvn.io"/>
        </FormRow>
        <div style={{ marginTop: 14 }}>
          <div className="tiny muted" style={{ marginBottom: 8 }}>Authentication</div>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 10 }}>
            <DomainCheck label="SPF" value="include:_spf.hanet.email" ok/>
            <DomainCheck label="DKIM" value="hanet-dkim._domainkey" ok/>
            <DomainCheck label="DMARC" value="v=DMARC1; p=none;" warn/>
          </div>
        </div>
      </div>
    </div>

    <div className="card">
      <div className="card-header"><h3>Transactional templates</h3></div>
      <table className="tbl">
        <thead><tr><th>Template</th><th>Trigger</th><th>Language</th><th>Last edited</th><th></th></tr></thead>
        <tbody>
          {[
            { t: 'Order confirmation', e: 'order.paid', l: 'EN · VI', d: '2 days ago' },
            { t: 'Provisioning complete', e: 'service.ready', l: 'EN · VI', d: '2 days ago' },
            { t: 'Provisioning failed', e: 'service.manual_review', l: 'EN · VI', d: '1 week ago' },
            { t: 'Renewal reminder (7d)', e: 'service.renew_soon', l: 'EN · VI', d: '1 week ago' },
            { t: 'Suspension notice', e: 'service.suspended', l: 'EN · VI', d: '3 weeks ago' },
            { t: 'Top-up verified', e: 'wallet.topup_verified', l: 'EN · VI', d: '1 week ago' },
          ].map(r => (
            <tr key={r.t}>
              <td style={{ fontWeight: 500 }}>{r.t}</td>
              <td className="mono tiny">{r.e}</td>
              <td className="tiny">{r.l}</td>
              <td className="tiny muted">{r.d}</td>
              <td style={{ textAlign: 'right' }}><button className="btn btn-ghost btn-sm">Edit</button></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

// ─── Legal ────────────────────────────────────────────────────
const LegalPanel = () => (
  <div className="card">
    <div className="card-header"><h3>Legal documents</h3></div>
    <table className="tbl">
      <thead><tr><th>Document</th><th>Required on</th><th>Version</th><th>Status</th><th></th></tr></thead>
      <tbody>
        {[
          { d: 'Terms of Service', req: 'Sign-up', v: 'v3 · 2026-02-14', s: 'published' },
          { d: 'Privacy Policy', req: 'Sign-up', v: 'v2 · 2026-01-08', s: 'published' },
          { d: 'Acceptable Use Policy', req: 'Each order', v: 'v4 · 2026-03-20', s: 'published' },
          { d: 'Refund Policy', req: 'Footer', v: 'v1 · 2026-01-08', s: 'published' },
          { d: 'DPA (EU clients)', req: 'On request', v: 'Draft', s: 'draft' },
        ].map(r => (
          <tr key={r.d}>
            <td style={{ fontWeight: 500 }}>{r.d}</td>
            <td className="tiny">{r.req}</td>
            <td className="tiny muted">{r.v}</td>
            <td><span className={`badge dot ${r.s === 'published' ? 'ok' : 'warn'}`}>{r.s}</span></td>
            <td style={{ textAlign: 'right' }}><button className="btn btn-ghost btn-sm">Edit</button></td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

// ─── Team ─────────────────────────────────────────────────────
const TeamPanel = () => (
  <div className="card">
    <div className="card-header">
      <h3>Team members</h3>
      <button className="btn btn-sm"><Icon name="plus" size={11}/> Invite member</button>
    </div>
    <table className="tbl">
      <thead><tr><th>Member</th><th>Role</th><th>2FA</th><th>Last active</th><th></th></tr></thead>
      <tbody>
        {[
          { n: 'Phong Tran', e: 'phong@proxyvn.io', r: 'Owner', f: true, a: '3m ago' },
          { n: 'Ngoc Pham', e: 'ngoc@proxyvn.io', r: 'Billing admin', f: true, a: '18m ago' },
          { n: 'Duy Le', e: 'duy@proxyvn.io', r: 'Support agent', f: false, a: '2h ago' },
          { n: 'An Vo', e: 'an@proxyvn.io', r: 'Support agent', f: true, a: '4d ago' },
        ].map(m => (
          <tr key={m.e}>
            <td>
              <div style={{ fontWeight: 500 }}>{m.n}</div>
              <div className="tiny muted">{m.e}</div>
            </td>
            <td><span className="badge">{m.r}</span></td>
            <td>{m.f ? <span className="badge ok dot">on</span> : <span className="badge warn dot">off</span>}</td>
            <td className="tiny muted">{m.a}</td>
            <td style={{ textAlign: 'right' }}><button className="btn btn-ghost btn-sm">Manage</button></td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

// ─── Shared helpers ──────────────────────────────────────────
const FormRow = ({ label, hint, children }) => (
  <div>
    <div style={{ display: 'flex', alignItems: 'baseline', justifyContent: 'space-between', marginBottom: 6 }}>
      <label style={{ fontSize: 12, fontWeight: 500 }}>{label}</label>
      {hint && <span className="tiny muted">{hint}</span>}
    </div>
    {children}
  </div>
);

const LogoUploader = ({ label, note, square, small }) => (
  <div>
    <div className="tiny" style={{ fontWeight: 500, marginBottom: 4 }}>{label}</div>
    <div style={{
      height: small ? 64 : 90,
      border: '1px dashed var(--line-3)',
      borderRadius: 3,
      background: `repeating-linear-gradient(45deg, var(--bg-alt), var(--bg-alt) 6px, var(--surface-2) 6px, var(--surface-2) 12px)`,
      display: 'grid', placeItems: 'center',
      fontFamily: 'var(--font-mono)', fontSize: 11, color: 'var(--ink-3)',
    }}>
      {square ? 'drop mark' : 'drop logo'}
    </div>
    <div className="tiny muted" style={{ marginTop: 4 }}>{note}</div>
  </div>
);

const DomainCheck = ({ label, value, ok, warn }) => (
  <div style={{
    padding: '8px 10px',
    border: '1px solid ' + (ok ? 'var(--ok-border)' : warn ? 'var(--warn-border)' : 'var(--line)'),
    background: ok ? 'var(--ok-bg)' : warn ? 'var(--warn-bg)' : 'var(--surface-2)',
    borderRadius: 3,
  }}>
    <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
      <Icon name={ok ? 'check' : 'alert'} size={11} style={{ color: ok ? 'var(--ok)' : 'var(--warn)' }}/>
      <span className="tiny" style={{ fontWeight: 600, color: ok ? 'var(--ok)' : warn ? 'var(--warn)' : 'var(--ink-1)' }}>{label}</span>
    </div>
    <div className="mono tiny muted" style={{ marginTop: 2, wordBreak: 'break-all' }}>{value}</div>
  </div>
);

window.ResellerBrandingScreen = ResellerBrandingScreen;
