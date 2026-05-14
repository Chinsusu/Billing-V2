package provider

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	cloudminiV3KindIPv4DC      = "ipv4_dc"
	cloudminiV3KindResidential = "residential"
	cloudminiV3ProtocolHTTP    = "http"
	cloudminiV3ProtocolSOCKS5  = "socks5"
)

type CloudminiV3SourceConfig struct {
	Kind             string
	GroupID          string
	NodeID           string
	Protocol         string
	BandwidthLimitMB int
	SpeedLimitMBps   int
}

type CloudminiV3Config struct {
	BaseURL          string
	APIToken         string
	CredentialCipher CredentialCipher
	KeyVersion       string
	DefaultSource    CloudminiV3SourceConfig
	SourceConfigs    map[SourceID]CloudminiV3SourceConfig
	PollInterval     time.Duration
	PollTimeout      time.Duration
	HTTPClient       cloudminiV3HTTPDoer
	Now              func() time.Time
}

type CloudminiV3Adapter struct {
	client           *cloudminiV3Client
	credentialCipher CredentialCipher
	keyVersion       string
	defaultSource    CloudminiV3SourceConfig
	sourceConfigs    map[SourceID]CloudminiV3SourceConfig
	pollInterval     time.Duration
	pollTimeout      time.Duration
	now              func() time.Time
}

func NewCloudminiV3Adapter(config CloudminiV3Config) (*CloudminiV3Adapter, error) {
	client, err := newCloudminiV3Client(config.BaseURL, config.APIToken, config.HTTPClient)
	if err != nil {
		return nil, err
	}
	if config.CredentialCipher == nil {
		return nil, ErrCredentialCipherMissing
	}
	pollInterval := config.PollInterval
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	pollTimeout := config.PollTimeout
	if pollTimeout <= 0 {
		pollTimeout = 30 * time.Second
	}
	sourceConfigs := make(map[SourceID]CloudminiV3SourceConfig, len(config.SourceConfigs))
	for sourceID, sourceConfig := range config.SourceConfigs {
		sourceConfigs[sourceID] = normalizeCloudminiV3SourceConfig(sourceConfig)
	}
	return &CloudminiV3Adapter{
		client:           client,
		credentialCipher: config.CredentialCipher,
		keyVersion:       strings.TrimSpace(config.KeyVersion),
		defaultSource:    normalizeCloudminiV3SourceConfig(config.DefaultSource),
		sourceConfigs:    sourceConfigs,
		pollInterval:     pollInterval,
		pollTimeout:      pollTimeout,
		now:              config.Now,
	}, nil
}

func (adapter *CloudminiV3Adapter) ProviderType() Type {
	return TypeCloudminiV3
}

func (adapter *CloudminiV3Adapter) CapabilityProfile() CapabilityProfile {
	return DefaultCapabilityProfile(TypeCloudminiV3)
}

func (adapter *CloudminiV3Adapter) CheckHealth(ctx context.Context, operation OperationContext, request HealthRequest) (HealthResult, error) {
	observedAt := adapter.observedAt()
	if _, err := adapter.client.getCapabilities(ctx, operation); err != nil {
		adapterErr := normalizeAdapterError(err, ErrorTemporary, "cloudmini v3 health check failed")
		return HealthResult{
			HealthStatus: healthStatusForCloudminiError(adapterErr),
			Result:       ResultFromError(adapterErr, observedAt),
		}, err
	}
	return HealthResult{
		HealthStatus: HealthStatusHealthy,
		Result:       SuccessResult(observedAt),
	}, nil
}

func (adapter *CloudminiV3Adapter) CheckStock(ctx context.Context, operation OperationContext, request StockRequest) (StockResult, error) {
	observedAt := adapter.observedAt()
	source, err := adapter.sourceConfig(operation)
	if err != nil {
		adapterErr := normalizeAdapterError(err, ErrorConfigInvalid, "cloudmini v3 source mapping is missing")
		return StockResult{StockStatus: StockStatusUnknown, Result: ResultFromError(adapterErr, observedAt)}, adapterErr
	}
	inventory, err := adapter.client.listGroupInventory(ctx, operation, source.Kind)
	if err != nil {
		adapterErr := normalizeAdapterError(err, ErrorTemporary, "cloudmini v3 inventory check failed")
		return StockResult{StockStatus: StockStatusUnknown, Result: ResultFromError(adapterErr, observedAt)}, err
	}
	for _, group := range inventory {
		if group.ID != source.GroupID {
			continue
		}
		if group.AllocatableUnits <= 0 {
			adapterErr := NewError(ErrorOutOfStock, "cloudmini v3 group has no allocatable units")
			return StockResult{
				StockStatus:   StockStatusOutOfStock,
				CapacityCount: 0,
				Result:        ResultFromError(adapterErr, observedAt),
			}, adapterErr
		}
		return StockResult{
			StockStatus:   StockStatusAvailable,
			CapacityCount: group.AllocatableUnits,
			Result:        SuccessResult(observedAt),
		}, nil
	}
	adapterErr := NewError(ErrorRegionUnavailable, "cloudmini v3 group mapping was not found")
	return StockResult{StockStatus: StockStatusUnknown, Result: ResultFromError(adapterErr, observedAt)}, adapterErr
}

func (adapter *CloudminiV3Adapter) Provision(ctx context.Context, operation OperationContext, request ProvisionRequest) (OperationResult, error) {
	if err := operation.Validate(); err != nil {
		return adapter.resultForError(NewError(ErrorConfigInvalid, "provider operation context is invalid"), "", "")
	}
	source, err := adapter.sourceConfig(operation)
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorConfigInvalid, "cloudmini v3 source mapping is missing"), "", "")
	}
	createResponse, err := adapter.client.createProxy(ctx, operation, cloudminiV3CreateProxyRequest{
		Kind:             source.Kind,
		GroupID:          source.GroupID,
		NodeID:           optionalString(source.NodeID),
		Protocol:         source.Protocol,
		BandwidthLimitMB: source.BandwidthLimitMB,
		SpeedLimitMBps:   source.SpeedLimitMBps,
		ExternalRef:      optionalString(string(operation.OperationID)),
	})
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorTimeoutUnknown, "cloudmini v3 create result is unknown"), "", "")
	}
	externalRequestID := ExternalRequestID(createResponse.Operation.ID)
	externalResourceID := ExternalResourceID(createResponse.Resource.ID)
	providerOperation, err := adapter.waitForOperation(ctx, operation, createResponse.Operation.ID)
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorTimeoutRequestKnown, "cloudmini v3 operation did not finish"), externalRequestID, externalResourceID)
	}
	if providerOperation.State != cloudminiV3OperationSucceeded {
		return adapter.resultForError(NewError(ErrorPartialSuccess, "cloudmini v3 operation did not complete successfully"), externalRequestID, externalResourceID)
	}
	proxy, err := adapter.proxyFromOperationOrRead(ctx, operation, providerOperation, string(externalResourceID))
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorCredentialMissing, "cloudmini v3 credential payload is missing"), externalRequestID, externalResourceID)
	}
	credential, err := adapter.credentialEnvelope(proxy)
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorCredentialMissing, "cloudmini v3 credential payload is missing"), externalRequestID, externalResourceID)
	}
	result := SuccessResult(adapter.observedAt())
	result.ExternalRequestID = externalRequestID
	result.ExternalResourceID = ExternalResourceID(proxy.ID)
	result.ServiceIdentifier = ServiceIdentifier(proxy.serviceIdentifier())
	result.Credential = credential
	result.ProviderStatus = proxy.Status
	return result, nil
}

func (adapter *CloudminiV3Adapter) GetStatus(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	resourceID := strings.TrimSpace(string(request.ExternalResourceID))
	if resourceID == "" {
		return adapter.resultForError(NewError(ErrorConfigInvalid, "cloudmini v3 resource id is missing"), "", "")
	}
	proxy, err := adapter.client.getProxy(ctx, operation, resourceID)
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorTemporary, "cloudmini v3 status read failed"), "", ExternalResourceID(resourceID))
	}
	result := SuccessResult(adapter.observedAt())
	result.ExternalResourceID = ExternalResourceID(proxy.ID)
	result.ServiceIdentifier = ServiceIdentifier(proxy.serviceIdentifier())
	result.ProviderStatus = proxy.Status
	return result, nil
}

func (adapter *CloudminiV3Adapter) Suspend(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.proxyAction(ctx, operation, request, "stop")
}

func (adapter *CloudminiV3Adapter) Unsuspend(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.proxyAction(ctx, operation, request, "start")
}

func (adapter *CloudminiV3Adapter) Terminate(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	if err := operation.Validate(); err != nil {
		return adapter.resultForError(NewError(ErrorConfigInvalid, "provider operation context is invalid"), "", "")
	}
	resourceID := strings.TrimSpace(string(request.ExternalResourceID))
	if resourceID == "" {
		return adapter.resultForError(NewError(ErrorConfigInvalid, "cloudmini v3 resource id is missing"), "", "")
	}
	response, err := adapter.client.deleteProxy(ctx, operation, resourceID)
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorTimeoutUnknown, "cloudmini v3 delete result is unknown"), "", ExternalResourceID(resourceID))
	}
	return adapter.resultFromMutatingOperation(ctx, operation, response.Operation.ID, response.Resource.ID)
}

func (adapter *CloudminiV3Adapter) Renew(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.unsupported("cloudmini v3 renew is not supported")
}

func (adapter *CloudminiV3Adapter) ResetPassword(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	return adapter.unsupported("cloudmini v3 password reset is not supported")
}

func (adapter *CloudminiV3Adapter) ChangeIP(ctx context.Context, operation OperationContext, request ResourceRequest) (OperationResult, error) {
	source, err := adapter.sourceConfig(operation)
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorConfigInvalid, "cloudmini v3 source mapping is missing"), "", request.ExternalResourceID)
	}
	if source.Kind != cloudminiV3KindResidential {
		return adapter.unsupported("cloudmini v3 change ip is residential only")
	}
	return adapter.proxyAction(ctx, operation, request, "change-ip")
}

func (adapter *CloudminiV3Adapter) proxyAction(ctx context.Context, operation OperationContext, request ResourceRequest, action string) (OperationResult, error) {
	if err := operation.Validate(); err != nil {
		return adapter.resultForError(NewError(ErrorConfigInvalid, "provider operation context is invalid"), "", "")
	}
	resourceID := strings.TrimSpace(string(request.ExternalResourceID))
	if resourceID == "" {
		return adapter.resultForError(NewError(ErrorConfigInvalid, "cloudmini v3 resource id is missing"), "", "")
	}
	response, err := adapter.client.proxyAction(ctx, operation, resourceID, action)
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorTimeoutUnknown, "cloudmini v3 action result is unknown"), "", ExternalResourceID(resourceID))
	}
	return adapter.resultFromMutatingOperation(ctx, operation, response.Operation.ID, response.Resource.ID)
}

func (adapter *CloudminiV3Adapter) resultFromMutatingOperation(ctx context.Context, operation OperationContext, operationID string, resourceID string) (OperationResult, error) {
	providerOperation, err := adapter.waitForOperation(ctx, operation, operationID)
	externalRequestID := ExternalRequestID(operationID)
	externalResourceID := ExternalResourceID(resourceID)
	if err != nil {
		return adapter.resultForError(normalizeAdapterError(err, ErrorTimeoutRequestKnown, "cloudmini v3 operation did not finish"), externalRequestID, externalResourceID)
	}
	if providerOperation.State != cloudminiV3OperationSucceeded {
		return adapter.resultForError(NewError(ErrorPartialSuccess, "cloudmini v3 operation did not complete successfully"), externalRequestID, externalResourceID)
	}
	result := SuccessResult(adapter.observedAt())
	result.ExternalRequestID = externalRequestID
	result.ExternalResourceID = externalResourceID
	result.ProviderStatus = string(providerOperation.State)
	return result, nil
}

func (adapter *CloudminiV3Adapter) waitForOperation(ctx context.Context, operation OperationContext, operationID string) (cloudminiV3Operation, error) {
	pollCtx, cancel := context.WithTimeout(ctx, adapter.pollTimeout)
	defer cancel()
	ticker := time.NewTicker(adapter.pollInterval)
	defer ticker.Stop()
	var lastOperation cloudminiV3Operation
	for {
		providerOperation, err := adapter.client.getOperation(pollCtx, operation, operationID)
		if err != nil {
			if pollCtx.Err() != nil {
				return lastOperation, AdapterError{
					Code:            ErrorTimeoutRequestKnown,
					MessageRedacted: "cloudmini v3 operation status is still pending",
					Safety:          RetrySafetyManualReviewRequired,
					Cause:           pollCtx.Err(),
				}
			}
			return cloudminiV3Operation{}, err
		}
		lastOperation = providerOperation
		switch providerOperation.State {
		case cloudminiV3OperationSucceeded, cloudminiV3OperationFailed, cloudminiV3OperationTimedOut, cloudminiV3OperationCancelled:
			return providerOperation, nil
		}
		select {
		case <-pollCtx.Done():
			return providerOperation, AdapterError{
				Code:            ErrorTimeoutRequestKnown,
				MessageRedacted: "cloudmini v3 operation status is still pending",
				Safety:          RetrySafetyManualReviewRequired,
				Cause:           pollCtx.Err(),
			}
		case <-ticker.C:
		}
	}
}

func (adapter *CloudminiV3Adapter) proxyFromOperationOrRead(ctx context.Context, operation OperationContext, providerOperation cloudminiV3Operation, resourceID string) (cloudminiV3Proxy, error) {
	if providerOperation.ResourceSnapshot != nil {
		proxy, err := providerOperation.proxySnapshot()
		if err == nil && proxy.ID != "" {
			return proxy, nil
		}
	}
	return adapter.client.getProxy(ctx, operation, resourceID)
}

func (adapter *CloudminiV3Adapter) credentialEnvelope(proxy cloudminiV3Proxy) (CredentialEnvelope, error) {
	if proxy.Host == "" || proxy.Username == "" || proxy.Password == "" {
		return CredentialEnvelope{}, NewError(ErrorCredentialMissing, "cloudmini v3 proxy credential fields are missing")
	}
	payload := map[string]string{
		"kind":        proxy.Kind,
		"host":        proxy.Host,
		"outbound_ip": proxy.OutboundIP,
		"username":    proxy.Username,
		"password":    proxy.Password,
	}
	if proxy.PortSocks > 0 {
		payload["port_socks"] = strconv.Itoa(proxy.PortSocks)
	}
	if proxy.PortHTTP > 0 {
		payload["port_http"] = strconv.Itoa(proxy.PortHTTP)
	}
	if proxy.ConnectionURI != "" {
		payload["connection_uri"] = proxy.ConnectionURI
	}
	return NewEncryptedCredentialEnvelope(CredentialTypeProxyAuth, payload, proxy.maskedHint(), adapter.keyVersion, adapter.credentialCipher)
}

func (adapter *CloudminiV3Adapter) sourceConfig(operation OperationContext) (CloudminiV3SourceConfig, error) {
	if source, ok := adapter.sourceConfigs[operation.SourceID]; ok {
		return validateCloudminiV3SourceConfig(source)
	}
	return validateCloudminiV3SourceConfig(adapter.defaultSource)
}

func (adapter *CloudminiV3Adapter) resultForError(err AdapterError, externalRequestID ExternalRequestID, externalResourceID ExternalResourceID) (OperationResult, error) {
	result := ResultFromError(err, adapter.observedAt())
	result.ExternalRequestID = externalRequestID
	result.ExternalResourceID = externalResourceID
	return result, err
}

func (adapter *CloudminiV3Adapter) unsupported(message string) (OperationResult, error) {
	return adapter.resultForError(NewError(ErrorCapabilityNotSupported, message), "", "")
}

func (adapter *CloudminiV3Adapter) observedAt() time.Time {
	if adapter.now != nil {
		return adapter.now().UTC()
	}
	return time.Now().UTC()
}

func healthStatusForCloudminiError(err AdapterError) HealthStatus {
	switch err.Code {
	case ErrorAuthFailed, ErrorPermissionDenied, ErrorAccountSuspended:
		return HealthStatusDown
	case ErrorRateLimited, ErrorTemporary, ErrorMaintenance:
		return HealthStatusDegraded
	default:
		return HealthStatusUnknown
	}
}

func normalizeCloudminiV3SourceConfig(config CloudminiV3SourceConfig) CloudminiV3SourceConfig {
	config.Kind = strings.ToLower(strings.TrimSpace(config.Kind))
	config.GroupID = strings.TrimSpace(config.GroupID)
	config.NodeID = strings.TrimSpace(config.NodeID)
	config.Protocol = strings.ToLower(strings.TrimSpace(config.Protocol))
	return config
}

func validateCloudminiV3SourceConfig(config CloudminiV3SourceConfig) (CloudminiV3SourceConfig, error) {
	config = normalizeCloudminiV3SourceConfig(config)
	if config.Kind != cloudminiV3KindIPv4DC && config.Kind != cloudminiV3KindResidential {
		return CloudminiV3SourceConfig{}, NewError(ErrorConfigInvalid, "cloudmini v3 kind must be ipv4_dc or residential")
	}
	if config.GroupID == "" {
		return CloudminiV3SourceConfig{}, NewError(ErrorConfigInvalid, "cloudmini v3 group id is missing")
	}
	if config.Protocol != cloudminiV3ProtocolHTTP && config.Protocol != cloudminiV3ProtocolSOCKS5 {
		return CloudminiV3SourceConfig{}, NewError(ErrorConfigInvalid, "cloudmini v3 protocol must be http or socks5")
	}
	if config.BandwidthLimitMB < 0 || config.SpeedLimitMBps < 0 {
		return CloudminiV3SourceConfig{}, NewError(ErrorConfigInvalid, "cloudmini v3 limits must not be negative")
	}
	return config, nil
}

func normalizeAdapterError(err error, fallbackCode ErrorCode, fallbackMessage string) AdapterError {
	var adapterErr AdapterError
	if errors.As(err, &adapterErr) {
		return adapterErr
	}
	return AdapterError{
		Code:            fallbackCode,
		MessageRedacted: fallbackMessage,
		Cause:           err,
	}
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
