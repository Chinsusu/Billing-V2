package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	SandboxCaseHealth      = "health_check"
	SandboxCaseQuote       = "quote_stock"
	SandboxCaseOrder       = "order_provision"
	SandboxCaseStatus      = "status_read"
	SandboxCaseCancel      = "cancel_terminate"
	SandboxCaseIdempotency = "idempotency_repeat"
)

type SandboxContractOptions struct {
	Adapter            Adapter
	Operation          OperationContext
	PlanKey            string
	Region             string
	ExternalResourceID ExternalResourceID
}

type SandboxContractResult struct {
	Name       string
	Status     OperationStatus
	Retry      RetrySafety
	ErrorCode  ErrorCode
	Err        error
	ObservedAt time.Time
}

func DefaultSandboxContractOptions(adapter Adapter) SandboxContractOptions {
	now := time.Date(2026, 4, 25, 0, 0, 0, 0, time.UTC)
	providerType := TypeManual
	if adapter != nil {
		providerType = adapter.ProviderType()
	}
	return SandboxContractOptions{
		Adapter: adapter,
		Operation: OperationContext{
			OperationID:       "sandbox-contract-op",
			TenantID:          tenant.ID("sandbox-tenant"),
			SourceID:          "sandbox-source",
			ProviderAccountID: "sandbox-account",
			ActorOrSystemID:   "sandbox-runner",
			IdempotencyKey:    "sandbox-tenant:sandbox-order:sandbox-item",
			CorrelationID:     "sandbox-contract-correlation",
			RequestTimeout:    15 * time.Second,
			Deadline:          now.Add(15 * time.Second),
			AttemptNumber:     1,
			CapabilitySnapshot: CapabilityProfile{
				SupportsHealthCheck:    true,
				SupportsLiveStockCheck: true,
				SupportsAutoProvision:  true,
				SupportsStatusSync:     true,
				SupportsTerminate:      true,
			},
			ProviderSourceSnapshot: SourceSnapshot{
				ProviderType: providerType,
				Name:         "Sandbox Contract Provider",
				Region:       "sandbox-region",
				PlanKey:      "sandbox-plan",
			},
		},
		PlanKey:            "sandbox-plan",
		Region:             "sandbox-region",
		ExternalResourceID: "sandbox-resource",
	}
}

func RunSandboxContract(ctx context.Context, options SandboxContractOptions) []SandboxContractResult {
	options = normalizeSandboxOptions(options)
	if err := validateSandboxOptions(options); err != nil {
		return []SandboxContractResult{{Name: "setup", Err: err}}
	}
	results := []SandboxContractResult{
		runHealthCase(ctx, options),
		runQuoteCase(ctx, options),
		runOrderCase(ctx, options),
		runStatusCase(ctx, options),
		runCancelCase(ctx, options),
		runIdempotencyCase(ctx, options),
	}
	return results
}

func SandboxContractPassed(results []SandboxContractResult) bool {
	for _, result := range results {
		if result.Err != nil {
			return false
		}
	}
	return true
}

func normalizeSandboxOptions(options SandboxContractOptions) SandboxContractOptions {
	if options.Operation.OperationID == "" && options.Adapter != nil {
		defaults := DefaultSandboxContractOptions(options.Adapter)
		if options.PlanKey == "" {
			options.PlanKey = defaults.PlanKey
		}
		if options.Region == "" {
			options.Region = defaults.Region
		}
		if options.ExternalResourceID == "" {
			options.ExternalResourceID = defaults.ExternalResourceID
		}
		options.Operation = defaults.Operation
	}
	return options
}

func validateSandboxOptions(options SandboxContractOptions) error {
	if options.Adapter == nil {
		return ErrAdapterNil
	}
	if err := options.Operation.Validate(); err != nil {
		return err
	}
	if options.PlanKey == "" {
		return fmt.Errorf("sandbox plan key missing")
	}
	if options.ExternalResourceID == "" {
		return fmt.Errorf("sandbox external resource id missing")
	}
	return nil
}

func runHealthCase(ctx context.Context, options SandboxContractOptions) SandboxContractResult {
	result, err := options.Adapter.CheckHealth(ctx, options.Operation, HealthRequest{})
	contract := contractResult(SandboxCaseHealth, result.Result, err)
	if err == nil && result.HealthStatus == "" {
		contract.Err = fmt.Errorf("health status missing")
	}
	return contract
}

func runQuoteCase(ctx context.Context, options SandboxContractOptions) SandboxContractResult {
	result, err := options.Adapter.CheckStock(ctx, options.Operation, StockRequest{PlanKey: options.PlanKey, Region: options.Region})
	contract := contractResult(SandboxCaseQuote, result.Result, err)
	if err == nil && result.StockStatus == "" {
		contract.Err = fmt.Errorf("stock status missing")
	}
	return contract
}

func runOrderCase(ctx context.Context, options SandboxContractOptions) SandboxContractResult {
	result, err := options.Adapter.Provision(ctx, options.Operation, ProvisionRequest{
		PlanKey:  options.PlanKey,
		Hostname: "sandbox-contract",
		InputSnapshot: map[string]string{
			"source": "sandbox-contract",
		},
	})
	return contractResult(SandboxCaseOrder, result, err)
}

func runStatusCase(ctx context.Context, options SandboxContractOptions) SandboxContractResult {
	result, err := options.Adapter.GetStatus(ctx, options.Operation, ResourceRequest{ExternalResourceID: options.ExternalResourceID})
	return contractResult(SandboxCaseStatus, result, err)
}

func runCancelCase(ctx context.Context, options SandboxContractOptions) SandboxContractResult {
	result, err := options.Adapter.Terminate(ctx, options.Operation, ResourceRequest{
		ExternalResourceID: options.ExternalResourceID,
		Reason:             "sandbox contract cleanup",
	})
	return contractResult(SandboxCaseCancel, result, err)
}

func runIdempotencyCase(ctx context.Context, options SandboxContractOptions) SandboxContractResult {
	first, firstErr := options.Adapter.Provision(ctx, options.Operation, ProvisionRequest{PlanKey: options.PlanKey, Hostname: "sandbox-contract"})
	if firstErr != nil {
		return contractResult(SandboxCaseIdempotency, first, firstErr)
	}
	second, secondErr := options.Adapter.Provision(ctx, options.Operation, ProvisionRequest{PlanKey: options.PlanKey, Hostname: "sandbox-contract"})
	contract := contractResult(SandboxCaseIdempotency, second, secondErr)
	if secondErr == nil && first.Status != second.Status {
		contract.Err = fmt.Errorf("idempotent repeat status changed from %s to %s", first.Status, second.Status)
	}
	return contract
}

func contractResult(name string, result OperationResult, err error) SandboxContractResult {
	contract := SandboxContractResult{
		Name:       name,
		Status:     result.Status,
		Retry:      result.RetrySafety,
		ErrorCode:  result.ErrorCode,
		Err:        err,
		ObservedAt: result.ObservedAt,
	}
	if err != nil {
		return contract
	}
	if result.Status == "" {
		contract.Err = fmt.Errorf("%s status missing", name)
		return contract
	}
	if result.ObservedAt.IsZero() {
		contract.Err = fmt.Errorf("%s observed_at missing", name)
		return contract
	}
	if containsSensitiveText(result.ErrorMessageRedacted) {
		contract.Err = fmt.Errorf("%s redacted message contains sensitive text", name)
	}
	return contract
}

func containsSensitiveText(value string) bool {
	normalized := strings.ToLower(value)
	for _, token := range []string{"secret", "token", "api_key", "authorization", "password"} {
		if strings.Contains(normalized, token) {
			return true
		}
	}
	return false
}
