package provider

type Type string

const (
	TypeManual             Type = "manual"
	TypeProxmox            Type = "proxmox"
	TypeOVH                Type = "ovh"
	TypeHetzner            Type = "hetzner"
	TypeProxyUpstream      Type = "proxy_upstream"
	TypePreloadedProxyPool Type = "preloaded_proxy_pool"
	TypeCustomAPI          Type = "custom_api"
)

type CapabilityProfile struct {
	SupportsHealthCheck        bool
	SupportsLiveStockCheck     bool
	SupportsAutoProvision      bool
	SupportsManualProvision    bool
	SupportsStatusSync         bool
	SupportsSuspend            bool
	SupportsUnsuspend          bool
	SupportsTerminate          bool
	SupportsRenew              bool
	SupportsResetPassword      bool
	SupportsReinstall          bool
	SupportsChangeIP           bool
	SupportsBandwidthUsage     bool
	SupportsConsole            bool
	SupportsReverseDNS         bool
	SupportsSnapshot           bool
	SupportsBackup             bool
	SupportsCredentialFetch    bool
	SupportsCredentialRotation bool
	VPS                        VPSCapabilities
	Proxy                      ProxyCapabilities
}

type VPSCapabilities struct {
	SupportsOSTemplateSelection bool
	SupportsCustomHostname      bool
	SupportsIPv6                bool
	SupportsPrivateNetwork      bool
	SupportsResize              bool
	SupportsRescueMode          bool
	SupportsVNCConsole          bool
	SupportsSSHKeyInjection     bool
}

type ProxyCapabilities struct {
	SupportsHTTPProtocol   bool
	SupportsSOCKS5Protocol bool
	SupportsRotatingProxy  bool
	SupportsStaticProxy    bool
	SupportsGeoSelection   bool
	SupportsIPWhitelist    bool
	SupportsUserPassAuth   bool
	SupportsBandwidthQuota bool
	SupportsThreadLimit    bool
	SupportsChangeExitIP   bool
}
