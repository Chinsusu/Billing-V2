// Client service detail — VPS console view + Proxy credentials view
// Reached when a client clicks "Manage" on a service row.

const ClientVPSDetailScreen = () => {
  const [tab, setTab] = React.useState('overview');
  const svc = {
    label: 'vps-scrape-01',
    specs: '2 vCPU · 4 GB RAM · 60 GB NVMe',
    region: 'VN-HCM · Proxmox node pmx-01',
    os: 'Ubuntu 24.04 LTS',
    ipv4: '103.28.44.21',
    ipv6: '2a01:4ff:2f0:3a28::1',
    hostname: 'vps-scrape-01.proxyvn.io',
    renewal: '2026-05-08',
    cycle: '$19.00 / month',
    createdAt: '2026-01-08',
  };

  return (
    <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
      {/* Service header */}
      <div style={{ display: 'flex', alignItems: 'flex-start', gap: 16 }}>
        <div style={{
          width: 44, height: 44, borderRadius: 4,
          background: 'var(--ink-6)', border: '1px solid var(--line)',
          display: 'grid', placeItems: 'center',
        }}>
          <Icon name="server" size={22}/>
        </div>
        <div style={{ flex: 1 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>{svc.label}</h2>
            <span className="badge ok dot">running</span>
          </div>
          <div className="tiny muted" style={{ marginTop: 4 }}>
            {svc.specs} · {svc.region} · <span className="mono">svc-v-4421</span>
          </div>
        </div>
        <div className="hstack">
          <button className="btn btn-sm"><Icon name="restart" size={12}/> Reboot</button>
          <button className="btn btn-sm"><Icon name="power" size={12}/> Stop</button>
          <button className="btn btn-sm btn-primary">Renew · $19</button>
        </div>
      </div>

      {/* Quick stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12 }}>
        <Kpi label="CPU" value="18%" sub="2 vCPU · load 0.42"/>
        <Kpi label="RAM" value="2.4 / 4 GB" sub="60% used"/>
        <Kpi label="Disk" value="18.2 / 60 GB" sub="read: 2.1 MB/s"/>
        <Kpi label="Bandwidth" value="128 GB" sub="this month · unmetered"/>
      </div>

      {/* Tabs */}
      <div style={{ display: 'flex', gap: 4, borderBottom: '1px solid var(--line)' }}>
        {[
          ['overview', 'Overview'], ['console', 'Console / VNC'],
          ['networking', 'Networking'], ['snapshots', 'Snapshots'],
          ['backups', 'Backups'], ['activity', 'Activity log'],
        ].map(([k, l]) => (
          <button key={k} onClick={() => setTab(k)} style={{
            padding: '10px 14px', border: 'none', background: 'transparent',
            fontSize: 13, fontWeight: tab === k ? 600 : 400,
            color: tab === k ? 'var(--accent)' : 'var(--ink-2)',
            borderBottom: tab === k ? '2px solid var(--accent)' : '2px solid transparent',
            cursor: 'pointer', marginBottom: -1, fontFamily: 'inherit',
          }}>{l}</button>
        ))}
      </div>

      {tab === 'overview' && (
        <div style={{ display: 'grid', gridTemplateColumns: '1.5fr 1fr', gap: 16 }}>
          <div className="card">
            <div className="card-header"><h3>Server details</h3></div>
            <div style={{ padding: 16 }}>
              <table className="tbl" style={{ border: 'none' }}>
                <tbody>
                  {[
                    ['Hostname', <span className="mono">{svc.hostname}</span>, true],
                    ['IPv4', <span className="mono">{svc.ipv4}</span>, true],
                    ['IPv6', <span className="mono">{svc.ipv6}</span>, true],
                    ['Operating system', svc.os, false],
                    ['Node', svc.region, false],
                    ['Created', svc.createdAt, false],
                    ['Next renewal', <><span>{svc.renewal}</span> <span className="tiny muted">· {svc.cycle}</span></>, false],
                  ].map(([k, v, copy], i) => (
                    <tr key={i}>
                      <td className="tiny muted" style={{ width: 140, border: 'none', padding: '6px 0' }}>{k}</td>
                      <td style={{ border: 'none', padding: '6px 0', fontSize: 13 }}>
                        {v}
                        {copy && <button className="btn btn-ghost btn-sm" style={{ marginLeft: 6, padding: 2 }} title="Copy"><Icon name="copy" size={11}/></button>}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
            <div className="card" style={{ padding: 16 }}>
              <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 10 }}>SSH access</div>
              <div style={{
                background: 'var(--ink-0)', color: '#F4F5F7',
                padding: '10px 12px', borderRadius: 3,
                fontFamily: 'var(--font-mono)', fontSize: 12,
                display: 'flex', alignItems: 'center', justifyContent: 'space-between',
              }}>
                <span>ssh root@{svc.ipv4}</span>
                <button style={{
                  background: 'transparent', border: 'none', color: '#A0A4AB',
                  cursor: 'pointer', padding: 2,
                }}><Icon name="copy" size={12}/></button>
              </div>
              <div className="tiny muted" style={{ marginTop: 8 }}>
                Root password was sent to your email at provisioning. <a href="#">Reset password →</a>
              </div>
            </div>

            <div className="card" style={{ padding: 16 }}>
              <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 10 }}>Danger zone</div>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                <button className="btn btn-sm" style={{ justifyContent: 'space-between' }}>
                  Rebuild from snapshot <Icon name="chevronRight" size={11}/>
                </button>
                <button className="btn btn-sm" style={{ justifyContent: 'space-between' }}>
                  Reinstall OS <Icon name="chevronRight" size={11}/>
                </button>
                <button className="btn btn-sm btn-danger" style={{ justifyContent: 'space-between' }}>
                  Cancel &amp; delete server <Icon name="chevronRight" size={11}/>
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {tab === 'console' && (
        <div className="card" style={{ padding: 0, overflow: 'hidden' }}>
          <div style={{
            display: 'flex', alignItems: 'center', justifyContent: 'space-between',
            padding: '8px 12px', borderBottom: '1px solid var(--line)',
            background: 'var(--surface-2)',
          }}>
            <div className="hstack">
              <span style={{ width: 8, height: 8, borderRadius: '50%', background: 'var(--ok)' }}/>
              <span className="tiny">Connected · noVNC · 1024×768</span>
            </div>
            <div className="hstack">
              <button className="btn btn-ghost btn-sm"><Icon name="refresh" size={11}/> Reconnect</button>
              <button className="btn btn-ghost btn-sm">Send Ctrl+Alt+Del</button>
              <button className="btn btn-ghost btn-sm"><Icon name="external" size={11}/> Fullscreen</button>
            </div>
          </div>
          <div style={{
            background: '#0B0D10',
            padding: 20,
            fontFamily: 'var(--font-mono)', fontSize: 13,
            color: '#D4D4D4', minHeight: 420,
            lineHeight: 1.55,
          }}>
            <div style={{ color: '#808080' }}>Ubuntu 24.04.1 LTS vps-scrape-01 tty1</div>
            <div style={{ color: '#808080' }}>&nbsp;</div>
            <div>vps-scrape-01 login: <span style={{ color: '#4EC9B0' }}>root</span></div>
            <div>Password: <span style={{ color: '#6A6A6A' }}>••••••••••</span></div>
            <div>Last login: Tue Apr 22 13:48:22 UTC 2026 on tty1</div>
            <div>&nbsp;</div>
            <div><span style={{ color: '#4EC9B0' }}>root@vps-scrape-01</span>:<span style={{ color: '#569CD6' }}>~</span># uptime</div>
            <div>&nbsp;14:02:18 up 9 days, 22:14, 1 user, load average: 0.42, 0.38, 0.31</div>
            <div><span style={{ color: '#4EC9B0' }}>root@vps-scrape-01</span>:<span style={{ color: '#569CD6' }}>~</span># systemctl status nginx</div>
            <div>&nbsp;● nginx.service - A high performance web server</div>
            <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Active: <span style={{ color: '#4EC9B0' }}>active (running)</span> since Tue 2026-04-22 08:12:04 UTC</div>
            <div>&nbsp;&nbsp;&nbsp;&nbsp;Process: 412 ExecStart=/usr/sbin/nginx -g daemon on; master_process on;</div>
            <div>&nbsp;&nbsp;&nbsp;Main PID: 418 (nginx)</div>
            <div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Tasks: 3 (limit: 4659)</div>
            <div>&nbsp;&nbsp;&nbsp;&nbsp;Memory: 12.4M</div>
            <div>&nbsp;</div>
            <div><span style={{ color: '#4EC9B0' }}>root@vps-scrape-01</span>:<span style={{ color: '#569CD6' }}>~</span># <span style={{ background: '#D4D4D4', color: '#0B0D10', padding: '0 2px' }}>&nbsp;</span></div>
          </div>
        </div>
      )}

      {tab === 'networking' && (
        <div className="card">
          <div className="card-header"><h3>Network interfaces &amp; firewall</h3></div>
          <table className="tbl">
            <thead>
              <tr><th>Interface</th><th>Address</th><th>Type</th><th>Bandwidth (mo)</th><th>Rule set</th></tr>
            </thead>
            <tbody>
              <tr><td className="mono">eth0</td><td className="mono">103.28.44.21/24</td><td>Public IPv4</td><td>128 GB</td><td><span className="badge ok dot">strict</span></td></tr>
              <tr><td className="mono">eth0</td><td className="mono">2a01:4ff:2f0:3a28::1/64</td><td>Public IPv6</td><td>—</td><td><span className="badge ok dot">strict</span></td></tr>
              <tr><td className="mono">vmbr0</td><td className="mono">10.0.0.21/16</td><td>Private</td><td>unlimited</td><td>—</td></tr>
            </tbody>
          </table>
        </div>
      )}

      {(tab === 'snapshots' || tab === 'backups' || tab === 'activity') && (
        <div className="card" style={{ padding: 40, textAlign: 'center' }}>
          <div style={{ fontSize: 13, fontWeight: 500 }}>No {tab} yet.</div>
          <div className="tiny muted" style={{ marginTop: 4 }}>This tab shows up once you create your first {tab.slice(0, -1)}.</div>
          <button className="btn btn-sm btn-primary" style={{ marginTop: 12 }}>
            <Icon name="plus" size={11}/> Create {tab.slice(0, -1)}
          </button>
        </div>
      )}
    </div>
  );
};

const ClientProxyDetailScreen = () => {
  const [tab, setTab] = React.useState('credentials');
  const [showPassword, setShowPassword] = React.useState(false);

  return (
    <div style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
      <div style={{ display: 'flex', alignItems: 'flex-start', gap: 16 }}>
        <div style={{
          width: 44, height: 44, borderRadius: 4,
          background: 'var(--ink-6)', border: '1px solid var(--line)',
          display: 'grid', placeItems: 'center',
        }}>
          <Icon name="globe" size={22}/>
        </div>
        <div style={{ flex: 1 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Residential EU · Premium</h2>
            <span className="badge ok dot">active</span>
          </div>
          <div className="tiny muted" style={{ marginTop: 4 }}>
            Rotating residential · EU-multi · <span className="mono">svc-r-9281</span>
          </div>
        </div>
        <div className="hstack">
          <button className="btn btn-sm"><Icon name="download" size={12}/> Export list</button>
          <button className="btn btn-sm btn-primary">Top up bandwidth</button>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12 }}>
        <Kpi label="Bandwidth used" value="4.2 GB" sub="of 10 GB · 42%"/>
        <Kpi label="Requests (24h)" value="182K" delta={12.4} sub="2.1 req/s avg"/>
        <Kpi label="Success rate" value="98.6%" sub="last 24h"/>
        <Kpi label="Renews" value="May 14" sub="30-day cycle"/>
      </div>

      <div style={{ display: 'flex', gap: 4, borderBottom: '1px solid var(--line)' }}>
        {[
          ['credentials', 'Credentials'], ['generator', 'List generator'],
          ['targeting', 'Targeting'], ['usage', 'Usage'], ['docs', 'Integration'],
        ].map(([k, l]) => (
          <button key={k} onClick={() => setTab(k)} style={{
            padding: '10px 14px', border: 'none', background: 'transparent',
            fontSize: 13, fontWeight: tab === k ? 600 : 400,
            color: tab === k ? 'var(--accent)' : 'var(--ink-2)',
            borderBottom: tab === k ? '2px solid var(--accent)' : '2px solid transparent',
            cursor: 'pointer', marginBottom: -1, fontFamily: 'inherit',
          }}>{l}</button>
        ))}
      </div>

      {tab === 'credentials' && (
        <div style={{ display: 'grid', gridTemplateColumns: '1.3fr 1fr', gap: 16, alignItems: 'start' }}>
          <div className="card">
            <div className="card-header">
              <h3>Proxy credentials</h3>
              <button className="btn btn-ghost btn-sm"><Icon name="refresh" size={11}/> Rotate password</button>
            </div>
            <div style={{ padding: 16, display: 'flex', flexDirection: 'column', gap: 12 }}>
              {[
                { label: 'Endpoint (rotating)', value: 'gw.pr.proxyvn.io:10000', mono: true, copy: true },
                { label: 'Endpoint (sticky 10min)', value: 'gw-sticky.pr.proxyvn.io:10001', mono: true, copy: true },
                { label: 'Username', value: 'linh_res_9281', mono: true, copy: true },
                { label: 'Password', value: showPassword ? 'Pr0xy!vN.2026-9281' : '••••••••••••••••', mono: true, copy: true, toggle: true },
              ].map((f, i) => (
                <div key={i}>
                  <div className="tiny muted" style={{ marginBottom: 4 }}>{f.label}</div>
                  <div style={{
                    display: 'flex', alignItems: 'center', gap: 6,
                    padding: '6px 10px',
                    border: '1px solid var(--line)', borderRadius: 3,
                    background: 'var(--surface-2)',
                    fontFamily: f.mono ? 'var(--font-mono)' : 'inherit',
                    fontSize: 12,
                  }}>
                    <span style={{ flex: 1 }}>{f.value}</span>
                    {f.toggle && (
                      <button onClick={() => setShowPassword(s => !s)} style={{ background: 'transparent', border: 'none', color: 'var(--ink-3)', cursor: 'pointer', padding: 2 }}>
                        <Icon name="eye" size={13}/>
                      </button>
                    )}
                    {f.copy && (
                      <button style={{ background: 'transparent', border: 'none', color: 'var(--ink-3)', cursor: 'pointer', padding: 2 }}>
                        <Icon name="copy" size={13}/>
                      </button>
                    )}
                  </div>
                </div>
              ))}

              <div style={{ marginTop: 4, padding: '10px 12px', background: 'var(--warn-bg)', border: '1px solid var(--warn-border)', borderRadius: 3 }}>
                <div style={{ fontSize: 12, color: 'var(--warn)', fontWeight: 500, display: 'flex', alignItems: 'center', gap: 6 }}>
                  <Icon name="alert" size={12}/> Keep credentials private
                </div>
                <div className="tiny" style={{ color: 'var(--ink-2)', marginTop: 2 }}>
                  These credentials route through your bandwidth quota. Rotating resets the password everywhere.
                </div>
              </div>
            </div>
          </div>

          <div className="card">
            <div className="card-header"><h3>Quick test</h3></div>
            <div style={{ padding: 16 }}>
              <div className="tiny muted" style={{ marginBottom: 8 }}>cURL example</div>
              <div style={{
                background: 'var(--ink-0)', color: '#F4F5F7',
                padding: 12, borderRadius: 3,
                fontFamily: 'var(--font-mono)', fontSize: 11.5, lineHeight: 1.55,
                whiteSpace: 'pre-wrap',
              }}>
{`curl -x http://gw.pr.proxyvn.io:10000 \\
  -U linh_res_9281:Pr0xy!vN.2026-9281 \\
  -L https://ipinfo.io`}
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 12, paddingTop: 12, borderTop: '1px solid var(--line)' }}>
                <span className="tiny muted">Last test</span>
                <span className="tiny" style={{ color: 'var(--ok)' }}>✓ 200 · 1.2s · from 185.94.x.x (DE)</span>
              </div>
              <button className="btn btn-sm" style={{ width: '100%', marginTop: 10 }}>Run test again</button>
            </div>
          </div>
        </div>
      )}

      {tab === 'generator' && (
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1.3fr', gap: 16, alignItems: 'start' }}>
          <div className="card" style={{ padding: 16 }}>
            <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 12 }}>Build a list</div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
              <div>
                <div className="tiny muted" style={{ marginBottom: 4 }}>Count</div>
                <input className="input" defaultValue="50" style={{ height: 28 }}/>
              </div>
              <div>
                <div className="tiny muted" style={{ marginBottom: 4 }}>Country</div>
                <select className="select" style={{ height: 28 }}>
                  <option>Any (EU pool)</option><option>Germany</option><option>France</option><option>UK</option><option>Netherlands</option>
                </select>
              </div>
              <div>
                <div className="tiny muted" style={{ marginBottom: 4 }}>Session</div>
                <div style={{ display: 'flex', gap: 4, background: 'var(--bg-alt)', padding: 2, borderRadius: 3 }}>
                  {['Rotating', 'Sticky 10m', 'Sticky 30m'].map((s, i) => (
                    <button key={s} style={{
                      flex: 1, padding: '5px 6px',
                      background: i === 0 ? '#fff' : 'transparent', border: 'none', borderRadius: 2,
                      fontSize: 11, fontWeight: i === 0 ? 500 : 400, cursor: 'pointer',
                      fontFamily: 'inherit',
                    }}>{s}</button>
                  ))}
                </div>
              </div>
              <div>
                <div className="tiny muted" style={{ marginBottom: 4 }}>Format</div>
                <select className="select" style={{ height: 28 }}>
                  <option>host:port:user:pass</option><option>user:pass@host:port</option><option>host:port (auth header)</option>
                </select>
              </div>
            </div>
            <button className="btn btn-primary" style={{ width: '100%', marginTop: 14 }}>Generate list</button>
          </div>

          <div className="card">
            <div className="card-header">
              <h3>Generated list · 50 entries</h3>
              <div className="hstack">
                <button className="btn btn-ghost btn-sm"><Icon name="copy" size={11}/> Copy all</button>
                <button className="btn btn-ghost btn-sm"><Icon name="download" size={11}/> Download .txt</button>
              </div>
            </div>
            <div style={{
              padding: 14, background: 'var(--ink-0)', color: '#D4D4D4',
              fontFamily: 'var(--font-mono)', fontSize: 11.5, lineHeight: 1.7,
              maxHeight: 340, overflow: 'auto',
            }}>
              {Array.from({ length: 12 }).map((_, i) => (
                <div key={i}>gw.pr.proxyvn.io:{10000 + i}:linh_res_9281-eu-{String(i).padStart(3,'0')}:Pr0xy!vN.2026-9281</div>
              ))}
              <div style={{ color: '#808080' }}>… 38 more lines</div>
            </div>
          </div>
        </div>
      )}

      {(tab === 'targeting' || tab === 'usage' || tab === 'docs') && (
        <div className="card" style={{ padding: 40, textAlign: 'center' }}>
          <div style={{ fontSize: 13, fontWeight: 500 }}>{tab[0].toUpperCase() + tab.slice(1)} view</div>
          <div className="tiny muted" style={{ marginTop: 4 }}>Placeholder — data tables &amp; charts go here.</div>
        </div>
      )}
    </div>
  );
};

window.ClientVPSDetailScreen = ClientVPSDetailScreen;
window.ClientProxyDetailScreen = ClientProxyDetailScreen;
