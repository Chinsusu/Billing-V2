package provider

// CapabilityProfile declares which operations a provider source supports.
// Unsupported operations must not be enqueued by the worker or exposed in the API.
type CapabilityProfile struct {
	// Common
	HealthCheck      bool
	LiveStockCheck   bool
	AutoProvision    bool
	ManualProvision  bool
	StatusSync       bool
	Suspend          bool
	Unsuspend        bool
	Terminate        bool
	Renew            bool
	ResetPassword    bool
	Reinstall        bool
	ChangeIP         bool
	BandwidthUsage   bool
	Console          bool
	ReverseDNS       bool
	Snapshot         bool
	Backup           bool
	CredentialFetch  bool
	CredentialRotate bool

	// VPS-specific
	OSTemplateSelection bool
	CustomHostname      bool
	IPv6                bool
	PrivateNetwork      bool
	Resize              bool
	RescueMode          bool
	VNCConsole          bool
	SSHKeyInjection     bool

	// Proxy-specific
	ProxyHTTP        bool
	ProxySOCKS5      bool
	RotatingProxy    bool
	StaticProxy      bool
	GeoSelection     bool
	IPWhitelist      bool
	UserPassAuth     bool
	BandwidthQuota   bool
	ThreadLimit      bool
	ChangeExitIP     bool
}
