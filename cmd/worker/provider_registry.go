package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/platform/secrets"
)

const (
	providerModeFake        = "fake"
	providerModeCloudminiV3 = "cloudmini_v3"
)

type workerProviderEnv struct {
	Mode                    string
	EncryptionKey           string
	CloudminiV3BaseURL      string
	CloudminiV3APIToken     string
	CloudminiV3MappingsJSON string
	CloudminiV3SourceID     string
	CloudminiV3Kind         string
	CloudminiV3GroupID      string
	CloudminiV3NodeID       string
	CloudminiV3Protocol     string
	CloudminiV3BandwidthMB  string
	CloudminiV3SpeedMBps    string
	CloudminiV3PollInterval string
	CloudminiV3PollTimeout  string
}

func newWorkerProviderRegistry() (*provider.Registry, error) {
	return buildWorkerProviderRegistry(readWorkerProviderEnv())
}

func readWorkerProviderEnv() workerProviderEnv {
	return workerProviderEnv{
		Mode:                    os.Getenv("PROVIDER_DEFAULT_MODE"),
		EncryptionKey:           os.Getenv("ENCRYPTION_KEY"),
		CloudminiV3BaseURL:      os.Getenv("CLOUDMINI_V3_BASE_URL"),
		CloudminiV3APIToken:     os.Getenv("CLOUDMINI_V3_API_TOKEN"),
		CloudminiV3MappingsJSON: os.Getenv("CLOUDMINI_V3_MAPPINGS_JSON"),
		CloudminiV3SourceID:     os.Getenv("CLOUDMINI_V3_SOURCE_ID"),
		CloudminiV3Kind:         os.Getenv("CLOUDMINI_V3_KIND"),
		CloudminiV3GroupID:      os.Getenv("CLOUDMINI_V3_GROUP_ID"),
		CloudminiV3NodeID:       os.Getenv("CLOUDMINI_V3_NODE_ID"),
		CloudminiV3Protocol:     os.Getenv("CLOUDMINI_V3_PROTOCOL"),
		CloudminiV3BandwidthMB:  os.Getenv("CLOUDMINI_V3_BANDWIDTH_LIMIT_MB"),
		CloudminiV3SpeedMBps:    os.Getenv("CLOUDMINI_V3_SPEED_LIMIT_MBPS"),
		CloudminiV3PollInterval: os.Getenv("CLOUDMINI_V3_POLL_INTERVAL"),
		CloudminiV3PollTimeout:  os.Getenv("CLOUDMINI_V3_POLL_TIMEOUT"),
	}
}

func buildWorkerProviderRegistry(env workerProviderEnv) (*provider.Registry, error) {
	mode := strings.ToLower(strings.TrimSpace(env.Mode))
	if mode == "" {
		mode = providerModeFake
	}
	switch mode {
	case providerModeFake:
		return provider.NewFakeRegistry()
	case providerModeCloudminiV3:
		return newCloudminiV3WorkerRegistry(env)
	default:
		return nil, fmt.Errorf("PROVIDER_DEFAULT_MODE must be fake or cloudmini_v3")
	}
}

func newCloudminiV3WorkerRegistry(env workerProviderEnv) (*provider.Registry, error) {
	config, err := cloudminiV3ConfigFromWorkerEnv(env)
	if err != nil {
		return nil, err
	}
	adapter, err := provider.NewCloudminiV3Adapter(config)
	if err != nil {
		return nil, err
	}
	registry, err := provider.NewFakeRegistry(
		provider.TypeManual,
		provider.TypeProxmox,
		provider.TypeOVH,
		provider.TypeHetzner,
		provider.TypeProxyUpstream,
		provider.TypePreloadedProxyPool,
		provider.TypeCustomAPI,
	)
	if err != nil {
		return nil, err
	}
	if err := registry.Register(adapter); err != nil {
		return nil, err
	}
	return registry, nil
}

type cloudminiV3MappingEnv struct {
	SourceID          string `json:"source_id"`
	ProviderAccountID string `json:"provider_account_id"`
	BaseURL           string `json:"base_url"`
	APIToken          string `json:"api_token"`
	Kind              string `json:"kind"`
	GroupID           string `json:"group_id"`
	NodeID            string `json:"node_id"`
	Protocol          string `json:"protocol"`
	BandwidthLimitMB  int    `json:"bandwidth_limit_mb"`
	SpeedLimitMBps    int    `json:"speed_limit_mbps"`
}

func cloudminiV3ConfigFromWorkerEnv(env workerProviderEnv) (provider.CloudminiV3Config, error) {
	if strings.TrimSpace(env.CloudminiV3MappingsJSON) != "" {
		return cloudminiV3MultiConfigFromWorkerEnv(env)
	}
	baseURL, err := requiredCloudminiV3Env("CLOUDMINI_V3_BASE_URL", env.CloudminiV3BaseURL)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	apiToken, err := requiredCloudminiV3Env("CLOUDMINI_V3_API_TOKEN", env.CloudminiV3APIToken)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	sourceID, err := requiredCloudminiV3Env("CLOUDMINI_V3_SOURCE_ID", env.CloudminiV3SourceID)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	kind, err := requiredCloudminiV3Enum("CLOUDMINI_V3_KIND", env.CloudminiV3Kind, "ipv4_dc", "residential")
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	groupID, err := requiredCloudminiV3Env("CLOUDMINI_V3_GROUP_ID", env.CloudminiV3GroupID)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	protocol, err := requiredCloudminiV3Enum("CLOUDMINI_V3_PROTOCOL", env.CloudminiV3Protocol, "http", "socks5")
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	bandwidthMB, err := optionalNonNegativeIntEnv("CLOUDMINI_V3_BANDWIDTH_LIMIT_MB", env.CloudminiV3BandwidthMB)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	speedMBps, err := optionalNonNegativeIntEnv("CLOUDMINI_V3_SPEED_LIMIT_MBPS", env.CloudminiV3SpeedMBps)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	pollInterval, err := optionalPositiveDurationEnv("CLOUDMINI_V3_POLL_INTERVAL", env.CloudminiV3PollInterval)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	pollTimeout, err := optionalPositiveDurationEnv("CLOUDMINI_V3_POLL_TIMEOUT", env.CloudminiV3PollTimeout)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	cipher, err := secrets.NewAESGCMCipher(env.EncryptionKey)
	if err != nil {
		return provider.CloudminiV3Config{}, fmt.Errorf("ENCRYPTION_KEY is required for cloudmini_v3 provider mode")
	}
	return provider.CloudminiV3Config{
		BaseURL:          baseURL,
		APIToken:         apiToken,
		CredentialCipher: cipher,
		SourceConfigs: map[provider.SourceID]provider.CloudminiV3SourceConfig{
			provider.SourceID(sourceID): {
				Kind:             kind,
				GroupID:          groupID,
				NodeID:           strings.TrimSpace(env.CloudminiV3NodeID),
				Protocol:         protocol,
				BandwidthLimitMB: bandwidthMB,
				SpeedLimitMBps:   speedMBps,
			},
		},
		PollInterval: pollInterval,
		PollTimeout:  pollTimeout,
	}, nil
}

func cloudminiV3MultiConfigFromWorkerEnv(env workerProviderEnv) (provider.CloudminiV3Config, error) {
	mappings, err := parseCloudminiV3MappingsEnv(env.CloudminiV3MappingsJSON)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	sourceEndpoints := make(map[provider.SourceID]provider.CloudminiV3EndpointConfig)
	accountEndpoints := make(map[provider.AccountID]provider.CloudminiV3EndpointConfig)
	for _, mapping := range mappings {
		endpoint, sourceID, accountID, err := cloudminiV3EndpointFromMappingEnv(mapping)
		if err != nil {
			return provider.CloudminiV3Config{}, err
		}
		if sourceID != "" {
			if _, exists := sourceEndpoints[sourceID]; exists {
				return provider.CloudminiV3Config{}, fmt.Errorf("CLOUDMINI_V3_MAPPINGS_JSON contains a duplicate source_id")
			}
			sourceEndpoints[sourceID] = endpoint
		}
		if accountID != "" {
			if _, exists := accountEndpoints[accountID]; exists {
				return provider.CloudminiV3Config{}, fmt.Errorf("CLOUDMINI_V3_MAPPINGS_JSON contains a duplicate provider_account_id")
			}
			accountEndpoints[accountID] = endpoint
		}
	}
	pollInterval, err := optionalPositiveDurationEnv("CLOUDMINI_V3_POLL_INTERVAL", env.CloudminiV3PollInterval)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	pollTimeout, err := optionalPositiveDurationEnv("CLOUDMINI_V3_POLL_TIMEOUT", env.CloudminiV3PollTimeout)
	if err != nil {
		return provider.CloudminiV3Config{}, err
	}
	cipher, err := secrets.NewAESGCMCipher(env.EncryptionKey)
	if err != nil {
		return provider.CloudminiV3Config{}, fmt.Errorf("ENCRYPTION_KEY is required for cloudmini_v3 provider mode")
	}
	return provider.CloudminiV3Config{
		CredentialCipher: cipher,
		SourceEndpoints:  sourceEndpoints,
		AccountEndpoints: accountEndpoints,
		PollInterval:     pollInterval,
		PollTimeout:      pollTimeout,
	}, nil
}

func parseCloudminiV3MappingsEnv(value string) ([]cloudminiV3MappingEnv, error) {
	var mappings []cloudminiV3MappingEnv
	if err := json.Unmarshal([]byte(value), &mappings); err != nil {
		return nil, fmt.Errorf("CLOUDMINI_V3_MAPPINGS_JSON must be valid JSON")
	}
	if len(mappings) == 0 {
		return nil, fmt.Errorf("CLOUDMINI_V3_MAPPINGS_JSON must contain at least one mapping")
	}
	return mappings, nil
}

func cloudminiV3EndpointFromMappingEnv(mapping cloudminiV3MappingEnv) (provider.CloudminiV3EndpointConfig, provider.SourceID, provider.AccountID, error) {
	sourceID := provider.SourceID(strings.TrimSpace(mapping.SourceID))
	accountID := provider.AccountID(strings.TrimSpace(mapping.ProviderAccountID))
	if sourceID == "" && accountID == "" {
		return provider.CloudminiV3EndpointConfig{}, "", "", fmt.Errorf("CLOUDMINI_V3_MAPPINGS_JSON entries require source_id or provider_account_id")
	}
	baseURL, err := requiredCloudminiV3Env("CLOUDMINI_V3_MAPPINGS_JSON.base_url", mapping.BaseURL)
	if err != nil {
		return provider.CloudminiV3EndpointConfig{}, "", "", err
	}
	apiToken, err := requiredCloudminiV3Env("CLOUDMINI_V3_MAPPINGS_JSON.api_token", mapping.APIToken)
	if err != nil {
		return provider.CloudminiV3EndpointConfig{}, "", "", err
	}
	kind, err := requiredCloudminiV3Enum("CLOUDMINI_V3_MAPPINGS_JSON.kind", mapping.Kind, "ipv4_dc", "residential")
	if err != nil {
		return provider.CloudminiV3EndpointConfig{}, "", "", err
	}
	groupID, err := requiredCloudminiV3Env("CLOUDMINI_V3_MAPPINGS_JSON.group_id", mapping.GroupID)
	if err != nil {
		return provider.CloudminiV3EndpointConfig{}, "", "", err
	}
	protocol, err := requiredCloudminiV3Enum("CLOUDMINI_V3_MAPPINGS_JSON.protocol", mapping.Protocol, "http", "socks5")
	if err != nil {
		return provider.CloudminiV3EndpointConfig{}, "", "", err
	}
	if mapping.BandwidthLimitMB < 0 || mapping.SpeedLimitMBps < 0 {
		return provider.CloudminiV3EndpointConfig{}, "", "", fmt.Errorf("CLOUDMINI_V3_MAPPINGS_JSON limits must be non-negative")
	}
	return provider.CloudminiV3EndpointConfig{
		BaseURL:  baseURL,
		APIToken: apiToken,
		Source: provider.CloudminiV3SourceConfig{
			Kind:             kind,
			GroupID:          groupID,
			NodeID:           strings.TrimSpace(mapping.NodeID),
			Protocol:         protocol,
			BandwidthLimitMB: mapping.BandwidthLimitMB,
			SpeedLimitMBps:   mapping.SpeedLimitMBps,
		},
	}, sourceID, accountID, nil
}

func requiredCloudminiV3Env(key string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s is required when PROVIDER_DEFAULT_MODE=cloudmini_v3", key)
	}
	return value, nil
}

func requiredCloudminiV3Enum(key string, value string, allowed ...string) (string, error) {
	value, err := requiredCloudminiV3Env(key, value)
	if err != nil {
		return "", err
	}
	value = strings.ToLower(value)
	for _, allowedValue := range allowed {
		if value == allowedValue {
			return value, nil
		}
	}
	return "", fmt.Errorf("%s must be one of %s", key, strings.Join(allowed, ", "))
}

func optionalNonNegativeIntEnv(key string, value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", key)
	}
	return parsed, nil
}

func optionalPositiveDurationEnv(key string, value string) (time.Duration, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", key)
	}
	return parsed, nil
}
