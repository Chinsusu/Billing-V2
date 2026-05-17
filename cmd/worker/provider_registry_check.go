package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

type workerProviderRegistryCheckSummary struct {
	Mode                     string
	CloudminiV3Adapter       string
	CloudminiSourceMappings  int
	CloudminiAccountMappings int
}

func runProviderRegistryCheck(w io.Writer) error {
	if err := guardLocalWorkerEnvironment(); err != nil {
		return err
	}
	summary, err := checkWorkerProviderRegistry(readWorkerProviderEnv())
	if err != nil {
		return err
	}
	writeWorkerProviderRegistryCheckSummary(w, summary)
	return nil
}

func checkWorkerProviderRegistry(env workerProviderEnv) (workerProviderRegistryCheckSummary, error) {
	mode := strings.ToLower(strings.TrimSpace(env.Mode))
	if mode == "" {
		mode = providerModeFake
	}
	registry, err := buildWorkerProviderRegistry(env)
	if err != nil {
		return workerProviderRegistryCheckSummary{}, err
	}
	adapter, err := registry.Get(provider.TypeCloudminiV3)
	if err != nil {
		return workerProviderRegistryCheckSummary{}, err
	}
	summary := workerProviderRegistryCheckSummary{
		Mode:               mode,
		CloudminiV3Adapter: "fake",
	}
	if _, ok := adapter.(*provider.CloudminiV3Adapter); ok {
		summary.CloudminiV3Adapter = "real"
	}
	if mode == providerModeCloudminiV3 {
		config, err := cloudminiV3ConfigFromWorkerEnv(env)
		if err != nil {
			return workerProviderRegistryCheckSummary{}, err
		}
		summary.CloudminiSourceMappings = len(config.SourceConfigs) + len(config.SourceEndpoints)
		summary.CloudminiAccountMappings = len(config.AccountEndpoints)
	}
	return summary, nil
}

func writeWorkerProviderRegistryCheckSummary(w io.Writer, summary workerProviderRegistryCheckSummary) {
	fmt.Fprintln(w, "provider-registry-check result=PASS")
	fmt.Fprintf(w, "mode=%s\n", summary.Mode)
	fmt.Fprintf(w, "cloudmini_v3_adapter=%s\n", summary.CloudminiV3Adapter)
	fmt.Fprintf(w, "cloudmini_source_mappings=%d\n", summary.CloudminiSourceMappings)
	fmt.Fprintf(w, "cloudmini_account_mappings=%d\n", summary.CloudminiAccountMappings)
	fmt.Fprintln(w, "provider_api_called=no")
	fmt.Fprintln(w, "mutating_routes_called=no")
	fmt.Fprintln(w, "jobs_claimed=0")
	fmt.Fprintln(w, "secrets_printed=no")
}
