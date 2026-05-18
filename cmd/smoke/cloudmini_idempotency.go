package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/secrets"
)

const (
	cloudminiScenarioDuplicateCreate   = "duplicate-create"
	cloudminiScenarioTimeoutAfterSend  = "timeout-after-send"
	cloudminiEvidenceRateLimitExpected = "no-parallel-mutating-calls"
	cloudminiEvidenceSpendExpected     = "single-dev-resource"
)

type cloudminiEvidenceConfig struct {
	AppEnv         string
	Scenario       string
	PilotID        string
	RawOutputPath  string
	BaseURL        string
	APIToken       string
	SourceID       string
	Kind           string
	GroupID        string
	NodeID         string
	Protocol       string
	EncryptionKey  string
	PollInterval   time.Duration
	PollTimeout    time.Duration
	CleanupTimeout time.Duration
	MaxCreates     int
}

type cloudminiRawEvidence struct {
	PilotID   string                     `json:"pilot_id"`
	Scenario  string                     `json:"scenario"`
	Operation string                     `json:"operation_id"`
	Resources []cloudminiRawEvidenceItem `json:"resources"`
}

type cloudminiRawEvidenceItem struct {
	ResultIndex        int    `json:"result_index"`
	ExternalRequestID  string `json:"external_request_id,omitempty"`
	ExternalResourceID string `json:"external_resource_id,omitempty"`
	Status             string `json:"status"`
}

func runCloudminiIdempotencyEvidenceSmoke(timeout time.Duration) error {
	return runCloudminiIdempotencyEvidenceSmokeWithWriter(timeout, os.Stdout)
}

func runCloudminiIdempotencyEvidenceSmokeWithWriter(timeout time.Duration, out io.Writer) error {
	config, err := cloudminiEvidenceConfigFromEnv()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	adapter, err := newCloudminiEvidenceAdapter(config)
	if err != nil {
		return err
	}
	operation := cloudminiEvidenceOperation(config, 1)
	results, scenarioErr := runCloudminiEvidenceScenario(ctx, adapter, config, operation)
	resources := cloudminiUniqueResources(results)
	if len(resources) == 0 {
		if scenarioErr != nil {
			return scenarioErr
		}
		return fmt.Errorf("cloudmini evidence scenario returned no cleanup resource")
	}
	if err := writeCloudminiRawEvidence(config, operation, results); err != nil {
		return err
	}

	cleanupAdapter, err := newCloudminiEvidenceAdapterWithTimeout(config, config.CleanupTimeout)
	if err != nil {
		return err
	}
	cleanupResults, cleanupErr := cleanupCloudminiEvidenceResources(ctx, cleanupAdapter, config, operation, resources)
	duplicateSameResource := cloudminiDuplicateSameResource(config.Scenario, results, resources)
	if cleanupErr != nil {
		printCloudminiEvidenceSummary(out, config, operation, results, cleanupResults, resources, duplicateSameResource, "FAIL")
		return cleanupErr
	}
	if scenarioErr != nil {
		printCloudminiEvidenceSummary(out, config, operation, results, cleanupResults, resources, duplicateSameResource, "FAIL")
		return scenarioErr
	}
	if config.Scenario == cloudminiScenarioDuplicateCreate && !duplicateSameResource {
		printCloudminiEvidenceSummary(out, config, operation, results, cleanupResults, resources, duplicateSameResource, "FAIL")
		return fmt.Errorf("cloudmini duplicate-create did not prove one redacted resource reference")
	}
	printCloudminiEvidenceSummary(out, config, operation, results, cleanupResults, resources, duplicateSameResource, "PASS")
	return nil
}

func cloudminiEvidenceConfigFromEnv() (cloudminiEvidenceConfig, error) {
	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if appEnv == "" {
		return cloudminiEvidenceConfig{}, fmt.Errorf("APP_ENV is required")
	}
	switch appEnv {
	case "local", "dev", "staging", "sandbox":
	case "prod", "production":
		return cloudminiEvidenceConfig{}, fmt.Errorf("refusing to run cloudmini evidence with APP_ENV=%s", appEnv)
	default:
		return cloudminiEvidenceConfig{}, fmt.Errorf("APP_ENV must be local, dev, staging, or sandbox")
	}
	if os.Getenv("BILLING_CLOUDMINI_IDEMPOTENCY_EVIDENCE_APPROVED") != "yes" {
		return cloudminiEvidenceConfig{}, fmt.Errorf("BILLING_CLOUDMINI_IDEMPOTENCY_EVIDENCE_APPROVED=yes is required")
	}
	for _, key := range []string{
		"CLOUDMINI_SOURCE_ACCOUNT_OWNER",
		"CLOUDMINI_ENGINEERING_OWNER",
		"CLOUDMINI_OPS_OWNER",
		"CLOUDMINI_SECURITY_OWNER",
		"CLOUDMINI_CLEANUP_OWNER",
		"CLOUDMINI_FINANCE_QUOTA_OWNER",
		"CLOUDMINI_REVIEWER_SIGNOFF",
		"CLOUDMINI_PILOT_CLEANUP_DEADLINE",
		"CLOUDMINI_PILOT_STOP_CONDITION",
		"CLOUDMINI_PILOT_READONLY_EVIDENCE_REF",
		"CLOUDMINI_PILOT_CLEANUP_PROCEDURE_REF",
	} {
		if err := requireCloudminiEvidenceFilled(key); err != nil {
			return cloudminiEvidenceConfig{}, err
		}
	}
	if strings.TrimSpace(os.Getenv("CLOUDMINI_IDEMPOTENCY_PROVIDER_RATE_LIMIT")) != cloudminiEvidenceRateLimitExpected {
		return cloudminiEvidenceConfig{}, fmt.Errorf("CLOUDMINI_IDEMPOTENCY_PROVIDER_RATE_LIMIT must be %s", cloudminiEvidenceRateLimitExpected)
	}
	if strings.TrimSpace(os.Getenv("CLOUDMINI_IDEMPOTENCY_MAX_SPEND_EXPOSURE")) != cloudminiEvidenceSpendExpected {
		return cloudminiEvidenceConfig{}, fmt.Errorf("CLOUDMINI_IDEMPOTENCY_MAX_SPEND_EXPOSURE must be %s", cloudminiEvidenceSpendExpected)
	}
	maxActive := strings.TrimSpace(os.Getenv("CLOUDMINI_IDEMPOTENCY_MAX_ACTIVE_RESOURCES"))
	if maxActive != "1" {
		return cloudminiEvidenceConfig{}, fmt.Errorf("CLOUDMINI_IDEMPOTENCY_MAX_ACTIVE_RESOURCES must be 1")
	}

	scenario := strings.TrimSpace(os.Getenv("CLOUDMINI_IDEMPOTENCY_SCENARIO"))
	maxCreates := 0
	switch scenario {
	case cloudminiScenarioDuplicateCreate:
		maxCreates = 2
	case cloudminiScenarioTimeoutAfterSend:
		maxCreates = 1
	default:
		return cloudminiEvidenceConfig{}, fmt.Errorf("CLOUDMINI_IDEMPOTENCY_SCENARIO must be duplicate-create or timeout-after-send")
	}
	if strings.TrimSpace(os.Getenv("CLOUDMINI_IDEMPOTENCY_MAX_CREATE_ATTEMPTS")) != fmt.Sprint(maxCreates) {
		return cloudminiEvidenceConfig{}, fmt.Errorf("CLOUDMINI_IDEMPOTENCY_MAX_CREATE_ATTEMPTS must be %d for %s", maxCreates, scenario)
	}

	rawPath, err := validateCloudminiRawEvidencePath(os.Getenv("CLOUDMINI_IDEMPOTENCY_RAW_EVIDENCE_PATH"))
	if err != nil {
		return cloudminiEvidenceConfig{}, err
	}
	pollInterval, err := optionalCloudminiDuration("CLOUDMINI_V3_POLL_INTERVAL", 250*time.Millisecond)
	if err != nil {
		return cloudminiEvidenceConfig{}, err
	}
	pollTimeout, err := optionalCloudminiDuration("CLOUDMINI_V3_POLL_TIMEOUT", 30*time.Second)
	if err != nil {
		return cloudminiEvidenceConfig{}, err
	}
	cleanupTimeout := pollTimeout
	if cleanupTimeout < 30*time.Second {
		cleanupTimeout = 30 * time.Second
	}
	if value := strings.TrimSpace(os.Getenv("CLOUDMINI_V3_CLEANUP_POLL_TIMEOUT")); value != "" {
		cleanupTimeout, err = optionalCloudminiDuration("CLOUDMINI_V3_CLEANUP_POLL_TIMEOUT", cleanupTimeout)
		if err != nil {
			return cloudminiEvidenceConfig{}, err
		}
	}
	config := cloudminiEvidenceConfig{
		AppEnv:         appEnv,
		Scenario:       scenario,
		PilotID:        strings.TrimSpace(os.Getenv("CLOUDMINI_IDEMPOTENCY_PILOT_ID")),
		RawOutputPath:  rawPath,
		BaseURL:        strings.TrimSpace(os.Getenv("CLOUDMINI_V3_BASE_URL")),
		APIToken:       strings.TrimSpace(os.Getenv("CLOUDMINI_V3_API_TOKEN")),
		SourceID:       strings.TrimSpace(os.Getenv("CLOUDMINI_V3_SOURCE_ID")),
		Kind:           strings.TrimSpace(os.Getenv("CLOUDMINI_V3_KIND")),
		GroupID:        strings.TrimSpace(os.Getenv("CLOUDMINI_V3_GROUP_ID")),
		NodeID:         strings.TrimSpace(os.Getenv("CLOUDMINI_V3_NODE_ID")),
		Protocol:       strings.TrimSpace(os.Getenv("CLOUDMINI_V3_PROTOCOL")),
		EncryptionKey:  strings.TrimSpace(os.Getenv("ENCRYPTION_KEY")),
		PollInterval:   pollInterval,
		PollTimeout:    pollTimeout,
		CleanupTimeout: cleanupTimeout,
		MaxCreates:     maxCreates,
	}
	for key, value := range map[string]string{
		"CLOUDMINI_V3_BASE_URL":  config.BaseURL,
		"CLOUDMINI_V3_API_TOKEN": config.APIToken,
		"CLOUDMINI_V3_SOURCE_ID": config.SourceID,
		"CLOUDMINI_V3_KIND":      config.Kind,
		"CLOUDMINI_V3_GROUP_ID":  config.GroupID,
		"CLOUDMINI_V3_PROTOCOL":  config.Protocol,
		"ENCRYPTION_KEY":         config.EncryptionKey,
	} {
		if value == "" {
			return cloudminiEvidenceConfig{}, fmt.Errorf("%s is required", key)
		}
	}
	if config.PilotID == "" {
		config.PilotID = "cloudmini-" + scenario + "-" + time.Now().UTC().Format("20060102T150405Z")
	}
	return config, nil
}

func requireCloudminiEvidenceFilled(key string) error {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fmt.Errorf("%s is required", key)
	}
	switch strings.ToLower(value) {
	case "-", "todo", "tbd", "unknown", "none", "null", "placeholder":
		return fmt.Errorf("%s must be a real approved value, not a placeholder", key)
	default:
		return nil
	}
}

func optionalCloudminiDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", key)
	}
	return duration, nil
}

func validateCloudminiRawEvidencePath(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("CLOUDMINI_IDEMPOTENCY_RAW_EVIDENCE_PATH is required")
	}
	if !filepath.IsAbs(value) {
		return "", fmt.Errorf("CLOUDMINI_IDEMPOTENCY_RAW_EVIDENCE_PATH must be absolute")
	}
	repoRoot, err := os.Getwd()
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(repoRoot, value)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, "../") {
		return "", fmt.Errorf("CLOUDMINI_IDEMPOTENCY_RAW_EVIDENCE_PATH must be outside the repository")
	}
	return value, nil
}

func newCloudminiEvidenceAdapter(config cloudminiEvidenceConfig) (*provider.CloudminiV3Adapter, error) {
	return newCloudminiEvidenceAdapterWithTimeout(config, config.PollTimeout)
}

func newCloudminiEvidenceAdapterWithTimeout(config cloudminiEvidenceConfig, pollTimeout time.Duration) (*provider.CloudminiV3Adapter, error) {
	cipher, err := secrets.NewAESGCMCipher(config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("ENCRYPTION_KEY is required for cloudmini evidence")
	}
	return provider.NewCloudminiV3Adapter(provider.CloudminiV3Config{
		BaseURL:          config.BaseURL,
		APIToken:         config.APIToken,
		CredentialCipher: cipher,
		SourceConfigs: map[provider.SourceID]provider.CloudminiV3SourceConfig{
			provider.SourceID(config.SourceID): {
				Kind:     config.Kind,
				GroupID:  config.GroupID,
				NodeID:   config.NodeID,
				Protocol: config.Protocol,
			},
		},
		PollInterval: config.PollInterval,
		PollTimeout:  pollTimeout,
	})
}

func cloudminiEvidenceOperation(config cloudminiEvidenceConfig, attempt int) provider.OperationContext {
	operationID := provider.OperationID("billing-" + config.PilotID + "-" + config.Scenario)
	return provider.OperationContext{
		OperationID:        operationID,
		TenantID:           tenant.ID("cloudmini-evidence-tenant"),
		SourceID:           provider.SourceID(config.SourceID),
		ActorOrSystemID:    provider.ActorID("cloudmini-evidence"),
		IdempotencyKey:     provider.IdempotencyKey(string(operationID)),
		CorrelationID:      provider.CorrelationID("req-" + hashRedacted(string(operationID))),
		AttemptNumber:      attempt,
		RequestTimeout:     config.PollTimeout,
		CapabilitySnapshot: provider.DefaultCapabilityProfile(provider.TypeCloudminiV3),
	}
}

func runCloudminiEvidenceScenario(ctx context.Context, adapter *provider.CloudminiV3Adapter, config cloudminiEvidenceConfig, operation provider.OperationContext) ([]provider.OperationResult, error) {
	switch config.Scenario {
	case cloudminiScenarioDuplicateCreate:
		first, _ := adapter.Provision(ctx, operation, provider.ProvisionRequest{PlanKey: "proxy-static-10gb-monthly"})
		second, _ := adapter.Provision(ctx, operation, provider.ProvisionRequest{PlanKey: "proxy-static-10gb-monthly"})
		return []provider.OperationResult{first, second}, nil
	case cloudminiScenarioTimeoutAfterSend:
		result, _ := adapter.Provision(ctx, operation, provider.ProvisionRequest{PlanKey: "proxy-static-10gb-monthly"})
		if result.ErrorCode != provider.ErrorTimeoutRequestKnown {
			return []provider.OperationResult{result}, fmt.Errorf("timeout-after-send scenario expected request-known timeout, got %s", result.ErrorCode)
		}
		return []provider.OperationResult{result}, nil
	default:
		return nil, fmt.Errorf("unsupported cloudmini evidence scenario")
	}
}

func cloudminiUniqueResources(results []provider.OperationResult) []provider.ExternalResourceID {
	seen := map[provider.ExternalResourceID]bool{}
	var resources []provider.ExternalResourceID
	for _, result := range results {
		if result.ExternalResourceID == "" || seen[result.ExternalResourceID] {
			continue
		}
		seen[result.ExternalResourceID] = true
		resources = append(resources, result.ExternalResourceID)
	}
	return resources
}

func writeCloudminiRawEvidence(config cloudminiEvidenceConfig, operation provider.OperationContext, results []provider.OperationResult) error {
	raw := cloudminiRawEvidence{PilotID: config.PilotID, Scenario: config.Scenario, Operation: string(operation.OperationID)}
	for index, result := range results {
		raw.Resources = append(raw.Resources, cloudminiRawEvidenceItem{
			ResultIndex:        index + 1,
			ExternalRequestID:  string(result.ExternalRequestID),
			ExternalResourceID: string(result.ExternalResourceID),
			Status:             string(result.Status),
		})
	}
	file, err := os.OpenFile(config.RawOutputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("create raw evidence file outside repo: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(raw); err != nil {
		return fmt.Errorf("write raw evidence file outside repo: %w", err)
	}
	return nil
}

func cleanupCloudminiEvidenceResources(ctx context.Context, adapter *provider.CloudminiV3Adapter, config cloudminiEvidenceConfig, operation provider.OperationContext, resources []provider.ExternalResourceID) ([]provider.OperationResult, error) {
	var results []provider.OperationResult
	for index, resourceID := range resources {
		cleanupOperation := operation
		cleanupOperation.OperationID = provider.OperationID(fmt.Sprintf("%s-cleanup-%d", operation.OperationID, index+1))
		cleanupOperation.IdempotencyKey = provider.IdempotencyKey(fmt.Sprintf("%s-cleanup-%d", operation.IdempotencyKey, index+1))
		cleanupOperation.AttemptNumber = 1
		result, err := adapter.Terminate(ctx, cleanupOperation, provider.ResourceRequest{ExternalResourceID: resourceID, Reason: "cloudmini evidence cleanup"})
		results = append(results, result)
		if err != nil || result.Status != provider.OperationStatusSuccess {
			return results, fmt.Errorf("cleanup failed for redacted resource %s", hashRedacted(string(resourceID)))
		}
	}
	return results, nil
}

func cloudminiDuplicateSameResource(scenario string, results []provider.OperationResult, resources []provider.ExternalResourceID) bool {
	if scenario != cloudminiScenarioDuplicateCreate {
		return false
	}
	if len(results) != 2 || len(resources) != 1 || results[0].ExternalResourceID == "" {
		return false
	}
	if results[1].ExternalResourceID == "" {
		return results[1].Status != provider.OperationStatusSuccess && results[1].ErrorCode != ""
	}
	return results[1].ExternalResourceID == results[0].ExternalResourceID
}

func printCloudminiEvidenceSummary(out io.Writer, config cloudminiEvidenceConfig, operation provider.OperationContext, results []provider.OperationResult, cleanup []provider.OperationResult, resources []provider.ExternalResourceID, duplicateSameResource bool, result string) {
	fmt.Fprintf(out, "cloudmini_idempotency_evidence result=%s\n", result)
	fmt.Fprintf(out, "scenario=%s\n", config.Scenario)
	fmt.Fprintf(out, "pilot_environment=%s\n", config.AppEnv)
	fmt.Fprintf(out, "pilot_id=%s\n", config.PilotID)
	fmt.Fprintf(out, "operation_ref=%s\n", hashRedacted(string(operation.OperationID)))
	fmt.Fprintln(out, "approval_fields_present=yes")
	fmt.Fprintln(out, "owner_fields_present=yes")
	fmt.Fprintln(out, "raw_cleanup_reference_path_private=yes")
	fmt.Fprintf(out, "provider_kind=%s\n", config.Kind)
	fmt.Fprintf(out, "protocol=%s\n", config.Protocol)
	fmt.Fprintln(out, "mutating_routes_called=yes")
	fmt.Fprintf(out, "create_attempts=%d\n", len(results))
	fmt.Fprintf(out, "distinct_resource_count=%d\n", len(resources))
	fmt.Fprintf(out, "duplicate_same_resource=%t\n", duplicateSameResource)
	if config.Scenario == cloudminiScenarioTimeoutAfterSend && len(results) == 1 {
		timeoutManualReview := results[0].ErrorCode == provider.ErrorTimeoutRequestKnown && results[0].RetrySafety == provider.RetrySafetyManualReviewRequired
		fmt.Fprintf(out, "timeout_after_send_manual_review=%t\n", timeoutManualReview)
	}
	for index, item := range results {
		fmt.Fprintf(out, "create_%d_status=%s\n", index+1, item.Status)
		fmt.Fprintf(out, "create_%d_error_code=%s\n", index+1, stringOrNone(string(item.ErrorCode)))
		fmt.Fprintf(out, "create_%d_retry_safety=%s\n", index+1, stringOrNone(string(item.RetrySafety)))
		fmt.Fprintf(out, "create_%d_resource_ref=%s\n", index+1, hashRedacted(string(item.ExternalResourceID)))
	}
	fmt.Fprintf(out, "cleanup_attempts=%d\n", len(cleanup))
	for index, item := range cleanup {
		fmt.Fprintf(out, "cleanup_%d_status=%s\n", index+1, item.Status)
		fmt.Fprintf(out, "cleanup_%d_error_code=%s\n", index+1, stringOrNone(string(item.ErrorCode)))
	}
	fmt.Fprintln(out, "sensitive_values_printed=no")
	fmt.Fprintln(out, "raw_provider_ids_printed=no")
	fmt.Fprintln(out, "provider_payloads_printed=no")
}

func hashRedacted(value string) string {
	if strings.TrimSpace(value) == "" {
		return "none"
	}
	sum := sha256.Sum256([]byte(value))
	return "redacted:" + hex.EncodeToString(sum[:])[:12]
}

func stringOrNone(value string) string {
	if strings.TrimSpace(value) == "" {
		return "none"
	}
	return value
}
