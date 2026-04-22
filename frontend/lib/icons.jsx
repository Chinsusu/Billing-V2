// Icon set — outline, 1.5px stroke, 16px default. Hetzner-like.
const Icon = ({ name, size = 16, stroke = 1.5, style, className }) => {
  const p = {
    fill: 'none',
    stroke: 'currentColor',
    strokeWidth: stroke,
    strokeLinecap: 'round',
    strokeLinejoin: 'round',
  };
  const paths = {
    // nav
    dashboard: <g {...p}><rect x="3" y="3" width="7" height="9"/><rect x="14" y="3" width="7" height="5"/><rect x="14" y="12" width="7" height="9"/><rect x="3" y="16" width="7" height="5"/></g>,
    users: <g {...p}><circle cx="9" cy="8" r="3.5"/><path d="M2 20c0-3.5 3-6 7-6s7 2.5 7 6"/><circle cx="17" cy="6" r="2.5"/><path d="M15.5 14c3 .3 5.5 2 5.5 5"/></g>,
    server: <g {...p}><rect x="3" y="4" width="18" height="7" rx="1"/><rect x="3" y="13" width="18" height="7" rx="1"/><circle cx="7" cy="7.5" r=".3" fill="currentColor"/><circle cx="7" cy="16.5" r=".3" fill="currentColor"/><path d="M11 7.5h6M11 16.5h6"/></g>,
    globe: <g {...p}><circle cx="12" cy="12" r="9"/><path d="M3 12h18M12 3c2.5 3 2.5 15 0 18M12 3c-2.5 3-2.5 15 0 18"/></g>,
    file: <g {...p}><path d="M14 3H6a1 1 0 0 0-1 1v16a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V8l-5-5z"/><path d="M14 3v5h5"/></g>,
    card: <g {...p}><rect x="2" y="5" width="20" height="14" rx="1.5"/><path d="M2 10h20M6 15h3"/></g>,
    tag: <g {...p}><path d="M3 12V3h9l9 9-9 9-9-9z"/><circle cx="7.5" cy="7.5" r="1.2"/></g>,
    chart: <g {...p}><path d="M3 3v18h18"/><path d="M7 15l4-5 3 3 5-7"/></g>,
    ticket: <g {...p}><path d="M3 8a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v2a2 2 0 0 0 0 4v2a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-2a2 2 0 0 0 0-4V8z"/><path d="M10 6v12"/></g>,
    settings: <g {...p}><circle cx="12" cy="12" r="3"/><path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41"/></g>,
    // actions
    search: <g {...p}><circle cx="11" cy="11" r="7"/><path d="m20 20-3.5-3.5"/></g>,
    plus: <g {...p}><path d="M12 5v14M5 12h14"/></g>,
    filter: <g {...p}><path d="M4 5h16l-6 8v6l-4-2v-4z"/></g>,
    download: <g {...p}><path d="M12 4v12m0 0l-4-4m4 4l4-4M4 20h16"/></g>,
    chevronDown: <g {...p}><path d="M6 9l6 6 6-6"/></g>,
    chevronRight: <g {...p}><path d="M9 6l6 6-6 6"/></g>,
    chevronLeft: <g {...p}><path d="M15 6l-6 6 6 6"/></g>,
    more: <g {...p}><circle cx="5" cy="12" r="1.2" fill="currentColor"/><circle cx="12" cy="12" r="1.2" fill="currentColor"/><circle cx="19" cy="12" r="1.2" fill="currentColor"/></g>,
    check: <g {...p}><path d="M5 12l5 5L20 7"/></g>,
    x: <g {...p}><path d="M6 6l12 12M18 6L6 18"/></g>,
    bell: <g {...p}><path d="M6 8a6 6 0 0 1 12 0c0 5 3 7 3 7H3s3-2 3-7"/><path d="M10 20a2 2 0 0 0 4 0"/></g>,
    mail: <g {...p}><rect x="3" y="5" width="18" height="14" rx="1.5"/><path d="M3 7l9 6 9-6"/></g>,
    cpu: <g {...p}><rect x="6" y="6" width="12" height="12" rx="1"/><rect x="9" y="9" width="6" height="6"/><path d="M6 10H3M6 14H3M21 10h-3M21 14h-3M10 6V3M14 6V3M10 21v-3M14 21v-3"/></g>,
    hdd: <g {...p}><rect x="3" y="4" width="18" height="16" rx="1.5"/><circle cx="12" cy="13" r="4"/><circle cx="12" cy="13" r="1" fill="currentColor"/></g>,
    shield: <g {...p}><path d="M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6l8-3z"/></g>,
    clock: <g {...p}><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></g>,
    wallet: <g {...p}><rect x="3" y="6" width="18" height="14" rx="1.5"/><path d="M3 10h18"/><circle cx="16" cy="15" r="1.2" fill="currentColor"/></g>,
    power: <g {...p}><path d="M12 3v9"/><path d="M7 6a7 7 0 1 0 10 0"/></g>,
    restart: <g {...p}><path d="M3 12a9 9 0 1 0 3-6.7"/><path d="M3 4v5h5"/></g>,
    copy: <g {...p}><rect x="8" y="8" width="12" height="12" rx="1.5"/><path d="M16 8V5a1 1 0 0 0-1-1H5a1 1 0 0 0-1 1v10a1 1 0 0 0 1 1h3"/></g>,
    external: <g {...p}><path d="M14 4h6v6"/><path d="M20 4l-9 9"/><path d="M19 13v6a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h6"/></g>,
    menu: <g {...p}><path d="M4 7h16M4 12h16M4 17h16"/></g>,
    eye: <g {...p}><path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12z"/><circle cx="12" cy="12" r="3"/></g>,
    arrowUp: <g {...p}><path d="M12 19V5M5 12l7-7 7 7"/></g>,
    arrowDown: <g {...p}><path d="M12 5v14M5 12l7 7 7-7"/></g>,
    arrowUpRight: <g {...p}><path d="M7 17L17 7M8 7h9v9"/></g>,
    calendar: <g {...p}><rect x="3" y="5" width="18" height="16" rx="1.5"/><path d="M3 10h18M8 3v4M16 3v4"/></g>,
    logout: <g {...p}><path d="M15 4h4a1 1 0 0 1 1 1v14a1 1 0 0 1-1 1h-4"/><path d="M10 17l-5-5 5-5M5 12h11"/></g>,
    refresh: <g {...p}><path d="M20 12a8 8 0 1 1-2.3-5.6"/><path d="M20 4v5h-5"/></g>,
    alert: <g {...p}><path d="M12 9v5"/><circle cx="12" cy="17" r=".8" fill="currentColor"/><path d="M12 3l10 18H2z"/></g>,
    db: <g {...p}><ellipse cx="12" cy="6" rx="8" ry="3"/><path d="M4 6v6c0 1.7 3.6 3 8 3s8-1.3 8-3V6"/><path d="M4 12v6c0 1.7 3.6 3 8 3s8-1.3 8-3v-6"/></g>,
    phone: <g {...p}><rect x="7" y="2" width="10" height="20" rx="2"/><path d="M10 18h4"/></g>,
    box: <g {...p}><path d="M3 7l9-4 9 4v10l-9 4-9-4z"/><path d="M3 7l9 4 9-4M12 11v10"/></g>,
  };
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" style={style} className={className} aria-hidden="true">
      {paths[name] || null}
    </svg>
  );
};

window.Icon = Icon;
