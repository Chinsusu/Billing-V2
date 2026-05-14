package provider

type Type string

const (
	TypeManual             Type = "manual"
	TypeProxmox            Type = "proxmox"
	TypeOVH                Type = "ovh"
	TypeHetzner            Type = "hetzner"
	TypeProxyUpstream      Type = "proxy_upstream"
	TypeCloudminiV3        Type = "cloudmini_v3"
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

func DefaultCapabilityProfile(providerType Type) CapabilityProfile {
	switch providerType {
	case TypeManual:
		return CapabilityProfile{
			SupportsHealthCheck:     true,
			SupportsManualProvision: true,
			SupportsStatusSync:      true,
		}
	case TypeProxyUpstream, TypeCloudminiV3, TypePreloadedProxyPool:
		return proxyCapabilityProfile()
	default:
		return vpsCapabilityProfile()
	}
}

func vpsCapabilityProfile() CapabilityProfile {
	return CapabilityProfile{
		SupportsHealthCheck:        true,
		SupportsLiveStockCheck:     true,
		SupportsAutoProvision:      true,
		SupportsStatusSync:         true,
		SupportsSuspend:            true,
		SupportsUnsuspend:          true,
		SupportsTerminate:          true,
		SupportsRenew:              true,
		SupportsResetPassword:      true,
		SupportsReinstall:          true,
		SupportsChangeIP:           true,
		SupportsBandwidthUsage:     true,
		SupportsConsole:            true,
		SupportsReverseDNS:         true,
		SupportsSnapshot:           true,
		SupportsBackup:             true,
		SupportsCredentialFetch:    true,
		SupportsCredentialRotation: true,
		VPS: VPSCapabilities{
			SupportsOSTemplateSelection: true,
			SupportsCustomHostname:      true,
			SupportsIPv6:                true,
			SupportsPrivateNetwork:      true,
			SupportsResize:              true,
			SupportsRescueMode:          true,
			SupportsVNCConsole:          true,
			SupportsSSHKeyInjection:     true,
		},
	}
}

func proxyCapabilityProfile() CapabilityProfile {
	return CapabilityProfile{
		SupportsHealthCheck:        true,
		SupportsLiveStockCheck:     true,
		SupportsAutoProvision:      true,
		SupportsStatusSync:         true,
		SupportsSuspend:            true,
		SupportsUnsuspend:          true,
		SupportsTerminate:          true,
		SupportsRenew:              true,
		SupportsChangeIP:           true,
		SupportsCredentialFetch:    true,
		SupportsCredentialRotation: true,
		Proxy: ProxyCapabilities{
			SupportsHTTPProtocol:   true,
			SupportsSOCKS5Protocol: true,
			SupportsRotatingProxy:  true,
			SupportsStaticProxy:    true,
			SupportsGeoSelection:   true,
			SupportsIPWhitelist:    true,
			SupportsUserPassAuth:   true,
			SupportsBandwidthQuota: true,
			SupportsThreadLimit:    true,
			SupportsChangeExitIP:   true,
		},
	}
}
