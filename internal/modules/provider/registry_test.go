package provider

import (
	"errors"
	"testing"
)

func TestRegistryReturnsAdapterByProviderType(t *testing.T) {
	adapter := NewFakeAdapter(TypeProxmox)
	registry, err := NewRegistry(adapter)
	if err != nil {
		t.Fatalf("expected registry, got %v", err)
	}

	got, err := registry.Get(TypeProxmox)
	if err != nil {
		t.Fatalf("expected adapter, got %v", err)
	}
	if got != adapter {
		t.Fatal("expected registered adapter instance")
	}
}

func TestRegistryRejectsDuplicateType(t *testing.T) {
	registry, err := NewRegistry(NewFakeAdapter(TypeManual))
	if err != nil {
		t.Fatalf("expected registry, got %v", err)
	}

	err = registry.Register(NewFakeAdapter(TypeManual))
	if !errors.Is(err, ErrAdapterAlreadyRegister) {
		t.Fatalf("expected duplicate adapter error, got %v", err)
	}
}

func TestRegistryAdapterForOperationUsesSourceSnapshot(t *testing.T) {
	registry, err := NewFakeRegistry(TypeManual, TypeProxyUpstream)
	if err != nil {
		t.Fatalf("expected fake registry, got %v", err)
	}
	operation := validOperation()
	operation.ProviderSourceSnapshot.ProviderType = TypeProxyUpstream

	adapter, err := registry.AdapterForOperation(operation)
	if err != nil {
		t.Fatalf("expected adapter, got %v", err)
	}
	if adapter.ProviderType() != TypeProxyUpstream {
		t.Fatalf("expected proxy adapter, got %s", adapter.ProviderType())
	}
}

func TestFakeRegistryDefaultsProviderSet(t *testing.T) {
	registry, err := NewFakeRegistry()
	if err != nil {
		t.Fatalf("expected fake registry, got %v", err)
	}

	types := registry.Types()
	if len(types) != 8 {
		t.Fatalf("expected default providers, got %#v", types)
	}
}

func TestDefaultCapabilityProfileSeparatesManualAndProxy(t *testing.T) {
	manual := DefaultCapabilityProfile(TypeManual)
	if !manual.SupportsManualProvision || manual.SupportsAutoProvision {
		t.Fatalf("unexpected manual capabilities: %#v", manual)
	}

	proxy := DefaultCapabilityProfile(TypeProxyUpstream)
	if !proxy.Proxy.SupportsSOCKS5Protocol || proxy.VPS.SupportsVNCConsole {
		t.Fatalf("unexpected proxy capabilities: %#v", proxy)
	}
}
