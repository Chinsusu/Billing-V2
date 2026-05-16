package provider

import "strings"

type CloudminiV3EndpointConfig struct {
	BaseURL  string
	APIToken string
	Source   CloudminiV3SourceConfig
}

type cloudminiV3RuntimeConfig struct {
	client *cloudminiV3Client
	source CloudminiV3SourceConfig
}

type cloudminiV3RuntimeSet struct {
	defaultClient    *cloudminiV3Client
	sourceConfigs    map[SourceID]CloudminiV3SourceConfig
	sourceEndpoints  map[SourceID]cloudminiV3RuntimeConfig
	accountEndpoints map[AccountID]cloudminiV3RuntimeConfig
}

func cloudminiV3RuntimeFromConfig(config CloudminiV3Config) (cloudminiV3RuntimeSet, error) {
	defaultClient, err := cloudminiV3DefaultClientFromConfig(config)
	if err != nil {
		return cloudminiV3RuntimeSet{}, err
	}
	sourceConfigs := make(map[SourceID]CloudminiV3SourceConfig, len(config.SourceConfigs))
	for sourceID, sourceConfig := range config.SourceConfigs {
		sourceConfigs[sourceID] = normalizeCloudminiV3SourceConfig(sourceConfig)
	}
	sourceEndpoints, err := cloudminiV3SourceEndpointsFromConfig(config)
	if err != nil {
		return cloudminiV3RuntimeSet{}, err
	}
	accountEndpoints, err := cloudminiV3AccountEndpointsFromConfig(config)
	if err != nil {
		return cloudminiV3RuntimeSet{}, err
	}
	return cloudminiV3RuntimeSet{
		defaultClient:    defaultClient,
		sourceConfigs:    sourceConfigs,
		sourceEndpoints:  sourceEndpoints,
		accountEndpoints: accountEndpoints,
	}, nil
}

func cloudminiV3DefaultClientFromConfig(config CloudminiV3Config) (*cloudminiV3Client, error) {
	hasDefaultEndpoint := strings.TrimSpace(config.BaseURL) != "" || strings.TrimSpace(config.APIToken) != ""
	hasExplicitEndpoints := len(config.SourceEndpoints) > 0 || len(config.AccountEndpoints) > 0
	if !hasDefaultEndpoint && hasExplicitEndpoints {
		return nil, nil
	}
	return newCloudminiV3Client(config.BaseURL, config.APIToken, config.HTTPClient)
}

func cloudminiV3SourceEndpointsFromConfig(config CloudminiV3Config) (map[SourceID]cloudminiV3RuntimeConfig, error) {
	endpoints := make(map[SourceID]cloudminiV3RuntimeConfig, len(config.SourceEndpoints))
	for sourceID, endpoint := range config.SourceEndpoints {
		runtime, err := cloudminiV3RuntimeEndpoint(endpoint, config.HTTPClient)
		if err != nil {
			return nil, err
		}
		endpoints[sourceID] = runtime
	}
	return endpoints, nil
}

func cloudminiV3AccountEndpointsFromConfig(config CloudminiV3Config) (map[AccountID]cloudminiV3RuntimeConfig, error) {
	endpoints := make(map[AccountID]cloudminiV3RuntimeConfig, len(config.AccountEndpoints))
	for accountID, endpoint := range config.AccountEndpoints {
		runtime, err := cloudminiV3RuntimeEndpoint(endpoint, config.HTTPClient)
		if err != nil {
			return nil, err
		}
		endpoints[accountID] = runtime
	}
	return endpoints, nil
}

func cloudminiV3RuntimeEndpoint(endpoint CloudminiV3EndpointConfig, httpClient cloudminiV3HTTPDoer) (cloudminiV3RuntimeConfig, error) {
	client, err := newCloudminiV3Client(endpoint.BaseURL, endpoint.APIToken, httpClient)
	if err != nil {
		return cloudminiV3RuntimeConfig{}, err
	}
	return cloudminiV3RuntimeConfig{
		client: client,
		source: normalizeCloudminiV3SourceConfig(endpoint.Source),
	}, nil
}

func (adapter *CloudminiV3Adapter) runtimeConfig(operation OperationContext) (cloudminiV3RuntimeConfig, error) {
	if operation.SourceID != "" {
		if runtime, ok := adapter.sourceEndpoints[operation.SourceID]; ok {
			return validateCloudminiV3RuntimeConfig(runtime)
		}
		if source, ok := adapter.sourceConfigs[operation.SourceID]; ok && adapter.defaultClient != nil {
			return validateCloudminiV3RuntimeConfig(cloudminiV3RuntimeConfig{client: adapter.defaultClient, source: source})
		}
		return cloudminiV3RuntimeConfig{}, NewError(ErrorConfigInvalid, "cloudmini v3 endpoint mapping is missing")
	}
	if operation.ProviderAccountID != "" {
		if runtime, ok := adapter.accountEndpoints[operation.ProviderAccountID]; ok {
			return validateCloudminiV3RuntimeConfig(runtime)
		}
	}
	if adapter.defaultClient == nil {
		return cloudminiV3RuntimeConfig{}, NewError(ErrorConfigInvalid, "cloudmini v3 endpoint mapping is missing")
	}
	return validateCloudminiV3RuntimeConfig(cloudminiV3RuntimeConfig{client: adapter.defaultClient, source: adapter.defaultSource})
}

func validateCloudminiV3RuntimeConfig(runtime cloudminiV3RuntimeConfig) (cloudminiV3RuntimeConfig, error) {
	if runtime.client == nil {
		return cloudminiV3RuntimeConfig{}, NewError(ErrorConfigInvalid, "cloudmini v3 endpoint mapping is missing")
	}
	source, err := validateCloudminiV3SourceConfig(runtime.source)
	if err != nil {
		return cloudminiV3RuntimeConfig{}, err
	}
	runtime.source = source
	return runtime, nil
}
