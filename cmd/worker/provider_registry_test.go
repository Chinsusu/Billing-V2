package main

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

func TestBuildWorkerProviderRegistryDefaultsToFake(t *testing.T) {
	registry, err := buildWorkerProviderRegistry(workerProviderEnv{})
	if err != nil {
		t.Fatalf("expected fake registry: %v", err)
	}
	adapter, err := registry.Get(provider.TypeCloudminiV3)
	if err != nil {
		t.Fatalf("expected cloudmini fake adapter: %v", err)
	}
	if _, ok := adapter.(*provider.FakeAdapter); !ok {
		t.Fatalf("expected fake cloudmini adapter, got %T", adapter)
	}
}

func TestBuildWorkerProviderRegistryRejectsUnknownMode(t *testing.T) {
	_, err := buildWorkerProviderRegistry(workerProviderEnv{Mode: "prod"})
	if err == nil || !strings.Contains(err.Error(), "PROVIDER_DEFAULT_MODE") {
		t.Fatalf("expected provider mode error, got %v", err)
	}
}

func TestBuildWorkerProviderRegistryCloudminiRequiresExplicitConfig(t *testing.T) {
	_, err := buildWorkerProviderRegistry(workerProviderEnv{Mode: "cloudmini_v3"})
	if err == nil || !strings.Contains(err.Error(), "CLOUDMINI_V3_BASE_URL") {
		t.Fatalf("expected missing base URL error, got %v", err)
	}
}

func TestBuildWorkerProviderRegistryRegistersCloudminiAdapter(t *testing.T) {
	registry, err := buildWorkerProviderRegistry(validCloudminiWorkerProviderEnv())
	if err != nil {
		t.Fatalf("expected cloudmini registry: %v", err)
	}
	adapter, err := registry.Get(provider.TypeCloudminiV3)
	if err != nil {
		t.Fatalf("expected cloudmini adapter: %v", err)
	}
	cloudminiAdapter, ok := adapter.(*provider.CloudminiV3Adapter)
	if !ok {
		t.Fatalf("expected real cloudmini adapter, got %T", adapter)
	}

	operation := provider.OperationContext{SourceID: "other-source"}
	result, err := cloudminiAdapter.CheckStock(context.Background(), operation, provider.StockRequest{PlanKey: "proxy"})
	var adapterErr provider.AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != provider.ErrorConfigInvalid {
		t.Fatalf("expected fail-closed source mapping error, got result=%+v err=%v", result, err)
	}

	manual, err := registry.Get(provider.TypeManual)
	if err != nil {
		t.Fatalf("expected manual fake adapter: %v", err)
	}
	if _, ok := manual.(*provider.FakeAdapter); !ok {
		t.Fatalf("expected other provider types to remain fake, got %T", manual)
	}
}

func TestBuildWorkerProviderRegistryAcceptsCloudminiMappingsJSON(t *testing.T) {
	env := workerProviderEnv{
		Mode:                    "cloudmini_v3",
		EncryptionKey:           strings.Repeat("1", 32),
		CloudminiV3MappingsJSON: `[{"source_id":"source-a","base_url":"http://cloudmini-a.example","api_token":"token-a","kind":"ipv4_dc","group_id":"group-a","protocol":"socks5"},{"source_id":"source-b","provider_account_id":"account-b","base_url":"http://cloudmini-b.example","api_token":"token-b","kind":"residential","group_id":"group-b","node_id":"node-b","protocol":"http","bandwidth_limit_mb":100,"speed_limit_mbps":10}]`,
		CloudminiV3PollInterval: "10ms",
		CloudminiV3PollTimeout:  "1s",
	}

	config, err := cloudminiV3ConfigFromWorkerEnv(env)
	if err != nil {
		t.Fatalf("expected multi mapping config: %v", err)
	}
	if len(config.SourceEndpoints) != 2 {
		t.Fatalf("expected two source endpoints, got %d", len(config.SourceEndpoints))
	}
	if len(config.AccountEndpoints) != 1 {
		t.Fatalf("expected one account endpoint, got %d", len(config.AccountEndpoints))
	}
	if config.SourceEndpoints["source-a"].BaseURL != "http://cloudmini-a.example" ||
		config.SourceEndpoints["source-a"].Source.GroupID != "group-a" {
		t.Fatalf("unexpected source A endpoint: %+v", config.SourceEndpoints["source-a"])
	}
	if config.SourceEndpoints["source-b"].Source.NodeID != "node-b" ||
		config.SourceEndpoints["source-b"].Source.BandwidthLimitMB != 100 ||
		config.SourceEndpoints["source-b"].Source.SpeedLimitMBps != 10 {
		t.Fatalf("unexpected source B endpoint: %+v", config.SourceEndpoints["source-b"])
	}

	registry, err := buildWorkerProviderRegistry(env)
	if err != nil {
		t.Fatalf("expected cloudmini registry: %v", err)
	}
	adapter, err := registry.Get(provider.TypeCloudminiV3)
	if err != nil {
		t.Fatalf("expected cloudmini adapter: %v", err)
	}
	result, err := adapter.CheckStock(context.Background(), provider.OperationContext{SourceID: "missing-source"}, provider.StockRequest{})
	var adapterErr provider.AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != provider.ErrorConfigInvalid {
		t.Fatalf("expected fail-closed missing mapping, got result=%+v err=%v", result, err)
	}
}

func TestCloudminiMappingsJSONRequiresSourceOrAccount(t *testing.T) {
	env := workerProviderEnv{
		Mode:                    "cloudmini_v3",
		EncryptionKey:           strings.Repeat("1", 32),
		CloudminiV3MappingsJSON: `[{"base_url":"http://cloudmini.example","api_token":"token-a","kind":"ipv4_dc","group_id":"group-a","protocol":"socks5"}]`,
	}

	_, err := cloudminiV3ConfigFromWorkerEnv(env)
	if err == nil || !strings.Contains(err.Error(), "source_id or provider_account_id") {
		t.Fatalf("expected selector error, got %v", err)
	}
}

func TestCloudminiConfigRejectsInvalidSourceShape(t *testing.T) {
	env := validCloudminiWorkerProviderEnv()
	env.CloudminiV3Kind = "vps"

	_, err := buildWorkerProviderRegistry(env)
	if err == nil || !strings.Contains(err.Error(), "CLOUDMINI_V3_KIND") {
		t.Fatalf("expected invalid kind error, got %v", err)
	}
}

func validCloudminiWorkerProviderEnv() workerProviderEnv {
	return workerProviderEnv{
		Mode:                    "cloudmini_v3",
		EncryptionKey:           strings.Repeat("1", 32),
		CloudminiV3BaseURL:      "http://cloudmini.example",
		CloudminiV3APIToken:     "sandbox-token",
		CloudminiV3SourceID:     "source-cloudmini-1",
		CloudminiV3Kind:         "ipv4_dc",
		CloudminiV3GroupID:      "group-cloudmini-1",
		CloudminiV3Protocol:     "socks5",
		CloudminiV3BandwidthMB:  "0",
		CloudminiV3SpeedMBps:    "0",
		CloudminiV3PollInterval: "10ms",
		CloudminiV3PollTimeout:  "1s",
	}
}
