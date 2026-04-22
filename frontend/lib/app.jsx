// Main apps — Admin / Reseller / Client (portal-aware shell)

const AdminApp = ({ initialScreen = 'overview', density = 'compact', collapsed: initCollapsed = false }) => {
  const [screen, setScreen] = React.useState(initialScreen);
  const [collapsed, setCollapsed] = React.useState(initCollapsed);

  const screenMap = {
    overview: { el: <OverviewScreen/>, title: 'Overview', bc: ['Home', 'Overview'], meta: <span className="badge ok dot">All systems operational</span> },
    tenants: { el: <AdminTenantsScreen/>, title: 'Tenants', bc: ['Home', 'Platform', 'Tenants'], meta: <span className="tiny muted">5 tenants · 4 resellers</span> },
    provisioning: { el: <AdminProvisioningScreen/>, title: 'Provisioning queue', bc: ['Home', 'Platform', 'Provisioning'], meta: <span className="badge warn dot">3 manual_review</span> },
    topups: { el: <AdminTopupsScreen/>, title: 'Top-up verification', bc: ['Home', 'Platform', 'Top-ups'], meta: <span className="badge warn dot">3 pending</span> },
    providers: { el: <AdminProvidersScreen/>, title: 'Providers / Sources', bc: ['Home', 'Platform', 'Providers'], meta: <span className="badge warn dot">1 degraded</span> },
    customers: { el: <CustomersScreen/>, title: 'Retail customers', bc: ['Home', 'Customers'] },
    tickets: { el: <TicketsScreen/>, title: 'Support tickets', bc: ['Home', 'Support', 'Tickets'] },
    proxies: { el: <ServicesScreen/>, title: 'Services', bc: ['Home', 'Services'] },
    vps: { el: <ServicesScreen/>, title: 'Services', bc: ['Home', 'Services'] },
    bandwidth: { el: <BandwidthScreen/>, title: 'Bandwidth', bc: ['Home', 'Services', 'Bandwidth'] },
    invoices: { el: <InvoicesScreen/>, title: 'Invoices', bc: ['Home', 'Billing', 'Invoices'] },
    transactions: { el: <TransactionsScreen/>, title: 'Transactions', bc: ['Home', 'Billing', 'Transactions'] },
    products: { el: <ProductsScreen/>, title: 'Products & Pricing', bc: ['Home', 'Billing', 'Products'] },
    reports: { el: <ReportsScreen/>, title: 'Reports', bc: ['Home', 'Reports'] },
    settings: { el: <SettingsScreen/>, title: 'Settings', bc: ['Home', 'Settings'] },
  };
  const cur = screenMap[screen] || screenMap.overview;

  return (
    <div data-density={density} style={{ display: 'flex', height: '100%', background: 'var(--bg)' }}>
      <Sidebar portal="admin" active={screen} onSelect={setScreen} collapsed={collapsed} onToggle={() => setCollapsed(c => !c)}/>
      <main style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0, overflow: 'hidden' }}>
        <Topbar title={cur.title} breadcrumbs={cur.bc} meta={cur.meta}
          actions={<button className="btn btn-sm"><Icon name="plus" size={12}/> New</button>}
        />
        <div style={{ flex: 1, overflow: 'auto' }} className="no-scrollbar">
          {cur.el}
        </div>
      </main>
    </div>
  );
};

// Reseller portal ("ProxyVN" POV)
const ResellerApp = ({ initialScreen = 'r-overview', density = 'compact', collapsed: initCollapsed = false }) => {
  const [screen, setScreen] = React.useState(initialScreen);
  const [collapsed, setCollapsed] = React.useState(initCollapsed);

  const screenMap = {
    'r-overview': { el: <ResellerDashboard/>, title: 'Dashboard', bc: ['ProxyVN', 'Dashboard'], meta: <span className="badge ok dot">Wallet healthy</span> },
    'r-clients': { el: <ResellerClientsScreen/>, title: 'Clients', bc: ['ProxyVN', 'Clients'], meta: <span className="tiny muted">312 clients</span> },
    'r-catalog': { el: <ResellerCatalogScreen/>, title: 'Catalog / Pricing', bc: ['ProxyVN', 'Catalog'], meta: <span className="badge warn dot">1 margin warning</span> },
    'r-services': { el: <ResellerClientsScreen/>, title: 'Services', bc: ['ProxyVN', 'Services'] },
    'r-orders': { el: <ResellerClientsScreen/>, title: 'Orders', bc: ['ProxyVN', 'Orders'] },
    'r-wallet': { el: <ResellerWalletScreen/>, title: 'Wallet & Top-up', bc: ['ProxyVN', 'Finance', 'Wallet'], meta: <span className="badge ok">$4,820.50</span> },
    'r-reports': { el: <ResellerDashboard/>, title: 'Reports', bc: ['ProxyVN', 'Reports'] },
    'r-settings': { el: <ResellerBrandingScreen/>, title: 'Branding & Settings', bc: ['ProxyVN', 'Settings'] },
  };
  const cur = screenMap[screen] || screenMap['r-overview'];

  return (
    <div data-density={density} style={{ display: 'flex', height: '100%', background: 'var(--bg)' }}>
      <Sidebar portal="reseller" active={screen} onSelect={setScreen} collapsed={collapsed} onToggle={() => setCollapsed(c => !c)}/>
      <main style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0, overflow: 'hidden' }}>
        <Topbar title={cur.title} breadcrumbs={cur.bc} meta={cur.meta}
          actions={<>
            <button className="btn btn-ghost btn-sm">proxyvn.io ↗</button>
            <button className="btn btn-sm btn-primary"><Icon name="plus" size={12}/> Top up</button>
          </>}
        />
        <div style={{ flex: 1, overflow: 'auto' }} className="no-scrollbar">
          {cur.el}
        </div>
      </main>
    </div>
  );
};

// Client portal (Linh Tran — reseller client of ProxyVN)
const CustomerApp = ({ initialScreen = 'c-overview', density = 'compact', collapsed: initCollapsed = false }) => {
  const [screen, setScreen] = React.useState(initialScreen);
  const [collapsed, setCollapsed] = React.useState(initCollapsed);

  const screenMap = {
    'c-overview': { el: <ClientDashboard/>, title: 'Dashboard', bc: ['ProxyVN', 'Dashboard'], meta: <span className="tiny muted">Wallet $128.40</span> },
    'c-shop': { el: <ClientShopScreen/>, title: 'Shop', bc: ['ProxyVN', 'Shop'] },
    'c-services': { el: <ClientDashboard/>, title: 'My services', bc: ['ProxyVN', 'Services'] },
    'c-vps-detail': { el: <ClientVPSDetailScreen/>, title: 'vps-scrape-01', bc: ['ProxyVN', 'Services', 'vps-scrape-01'], meta: <span className="badge ok dot">running</span> },
    'c-proxy-detail': { el: <ClientProxyDetailScreen/>, title: 'Residential EU · Premium', bc: ['ProxyVN', 'Services', 'Residential EU'], meta: <span className="badge ok dot">active</span> },
    'c-checkout': { el: <ClientCheckoutScreen/>, title: 'Checkout', bc: ['ProxyVN', 'Shop', 'VPS', 'Checkout'] },
    'c-wallet': { el: <ClientWalletScreen/>, title: 'Wallet', bc: ['ProxyVN', 'Billing', 'Wallet'], meta: <span className="badge ok">$128.40</span> },
    'c-usage': { el: <ClientDashboard/>, title: 'Usage', bc: ['ProxyVN', 'Billing', 'Usage'] },
    'c-settings': { el: <ClientDashboard/>, title: 'Settings', bc: ['ProxyVN', 'Account'] },
    'c-support': { el: <ClientDashboard/>, title: 'Support', bc: ['ProxyVN', 'Support'] },
  };
  const cur = screenMap[screen] || screenMap['c-overview'];

  return (
    <div data-density={density} style={{ display: 'flex', height: '100%', background: 'var(--bg)' }}>
      <Sidebar portal="customer" active={screen} onSelect={setScreen} collapsed={collapsed} onToggle={() => setCollapsed(c => !c)}/>
      <main style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0, overflow: 'hidden' }}>
        <Topbar title={cur.title} breadcrumbs={cur.bc} meta={cur.meta}/>
        <div style={{ flex: 1, overflow: 'auto' }} className="no-scrollbar">
          {cur.el}
        </div>
      </main>
    </div>
  );
};

// Standalone overview (no sidebar switching) — for focused variation artboards
const OverviewInShell = ({ Screen, density = 'compact', collapsed = false }) => (
  <div data-density={density} style={{ display: 'flex', height: '100%', background: 'var(--bg)' }}>
    <Sidebar portal="admin" active="overview" onSelect={() => {}} collapsed={collapsed} onToggle={() => {}}/>
    <main style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0, overflow: 'hidden' }}>
      <Topbar title="Overview" breadcrumbs={['Home', 'Overview']} meta={<span className="badge ok dot">All systems operational</span>}/>
      <div style={{ flex: 1, overflow: 'auto' }} className="no-scrollbar">
        <Screen/>
      </div>
    </main>
  </div>
);

Object.assign(window, { AdminApp, ResellerApp, CustomerApp, OverviewInShell });
