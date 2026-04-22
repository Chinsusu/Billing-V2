package provider

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrProviderTypeMissing    = errors.New("provider type missing")
	ErrAdapterMissing         = errors.New("provider adapter missing")
	ErrAdapterNil             = errors.New("provider adapter nil")
	ErrAdapterTypeMismatch    = errors.New("provider adapter type mismatch")
	ErrAdapterAlreadyRegister = errors.New("provider adapter already registered")
)

type Registry struct {
	adapters map[Type]Adapter
}

func NewRegistry(adapters ...Adapter) (*Registry, error) {
	registry := &Registry{adapters: make(map[Type]Adapter)}
	for _, adapter := range adapters {
		if err := registry.Register(adapter); err != nil {
			return nil, err
		}
	}
	return registry, nil
}

func (registry *Registry) Register(adapter Adapter) error {
	if adapter == nil {
		return ErrAdapterNil
	}
	providerType := adapter.ProviderType()
	if providerType == "" {
		return ErrProviderTypeMissing
	}
	if registry.adapters == nil {
		registry.adapters = make(map[Type]Adapter)
	}
	if existing := registry.adapters[providerType]; existing != nil {
		return fmt.Errorf("%w: %s", ErrAdapterAlreadyRegister, providerType)
	}
	registry.adapters[providerType] = adapter
	return nil
}

func (registry *Registry) Get(providerType Type) (Adapter, error) {
	if providerType == "" {
		return nil, ErrProviderTypeMissing
	}
	if registry == nil || registry.adapters == nil {
		return nil, fmt.Errorf("%w: %s", ErrAdapterMissing, providerType)
	}
	adapter := registry.adapters[providerType]
	if adapter == nil {
		return nil, fmt.Errorf("%w: %s", ErrAdapterMissing, providerType)
	}
	if adapter.ProviderType() != providerType {
		return nil, fmt.Errorf("%w: %s", ErrAdapterTypeMismatch, providerType)
	}
	return adapter, nil
}

func (registry *Registry) AdapterForOperation(operation OperationContext) (Adapter, error) {
	return registry.Get(operation.ProviderSourceSnapshot.ProviderType)
}

func (registry *Registry) Types() []Type {
	if registry == nil || len(registry.adapters) == 0 {
		return nil
	}
	types := make([]Type, 0, len(registry.adapters))
	for providerType := range registry.adapters {
		types = append(types, providerType)
	}
	sort.Slice(types, func(i, j int) bool { return types[i] < types[j] })
	return types
}

func NewFakeRegistry(providerTypes ...Type) (*Registry, error) {
	if len(providerTypes) == 0 {
		providerTypes = []Type{TypeManual, TypeProxmox, TypeOVH, TypeHetzner, TypeProxyUpstream, TypePreloadedProxyPool, TypeCustomAPI}
	}
	registry, err := NewRegistry()
	if err != nil {
		return nil, err
	}
	for _, providerType := range providerTypes {
		if err := registry.Register(NewFakeAdapter(providerType)); err != nil {
			return nil, err
		}
	}
	return registry, nil
}
