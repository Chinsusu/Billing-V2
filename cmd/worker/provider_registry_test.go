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
