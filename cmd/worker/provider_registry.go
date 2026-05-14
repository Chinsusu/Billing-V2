package main

import (
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

func cloudminiV3ConfigFromWorkerEnv(env workerProviderEnv) (provider.CloudminiV3Config, error) {
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
