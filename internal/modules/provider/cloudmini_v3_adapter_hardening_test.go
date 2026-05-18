package provider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCloudminiV3AdapterProvisionRequiresUsableProxyStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiAuth(t, r)
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/proxies":
			writeCloudminiSuccess(t, w, http.StatusAccepted, cloudminiV3ProxyMutationResponse{
				Resource:  cloudminiV3Resource{ID: "proxy-1", Kind: "ipv4_dc", Status: "creating"},
				Operation: cloudminiV3Operation{ID: "op-1", State: cloudminiV3OperationAccepted},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-1":
			snapshot := cloudminiV3Proxy{
				ID:         "proxy-1",
				Kind:       "ipv4_dc",
				Status:     "creating",
				Host:       "203.0.113.10",
				PortSocks:  1080,
				Username:   "proxy-user",
				Password:   "proxy-pass",
				OutboundIP: "203.0.113.10",
			}
			rawSnapshot, err := json.Marshal(snapshot)
			if err != nil {
				t.Fatalf("marshal snapshot: %v", err)
			}
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Operation{
				ID:               "op-1",
				ResourceID:       optionalString("proxy-1"),
				State:            cloudminiV3OperationSucceeded,
				ResourceSnapshot: rawSnapshot,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/proxies/proxy-1":
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Proxy{
				ID:         "proxy-1",
				Kind:       "ipv4_dc",
				Status:     "creating",
				Host:       "203.0.113.10",
				PortSocks:  1080,
				Username:   "proxy-user",
				Password:   "proxy-pass",
				OutboundIP: "203.0.113.10",
			})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := newTestCloudminiV3Adapter(t, server.URL, &cloudminiV3TestCipher{}, CloudminiV3SourceConfig{
		Kind:     "ipv4_dc",
		GroupID:  "group-1",
		Protocol: "socks5",
	})
	adapter.pollInterval = time.Millisecond
	adapter.pollTimeout = 75 * time.Millisecond
	result, err := adapter.Provision(context.Background(), validOperation(), ProvisionRequest{PlanKey: "proxy"})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorPartialSuccess {
		t.Fatalf("expected partial-success status error, got result=%+v err=%v", result, err)
	}
	if result.Status != OperationStatusPartialSuccess ||
		result.RetrySafety != RetrySafetyManualReviewRequired ||
		result.ExternalRequestID != "op-1" ||
		result.ExternalResourceID != "proxy-1" ||
		result.ProviderStatus != "creating" {
		t.Fatalf("unexpected not-usable result: %+v", result)
	}
	if result.Credential.HasEncryptedPayload() {
		t.Fatalf("credential must not be returned for not-usable proxy status: %+v", result.Credential)
	}
}

func TestCloudminiV3AdapterProvisionWaitsForUsableProxyStatus(t *testing.T) {
	var proxyReads int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiAuth(t, r)
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/proxies":
			writeCloudminiSuccess(t, w, http.StatusAccepted, cloudminiV3ProxyMutationResponse{
				Resource:  cloudminiV3Resource{ID: "proxy-1", Kind: "ipv4_dc", Status: "creating"},
				Operation: cloudminiV3Operation{ID: "op-1", State: cloudminiV3OperationAccepted},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-1":
			rawSnapshot, err := json.Marshal(cloudminiV3Proxy{
				ID:     "proxy-1",
				Kind:   "ipv4_dc",
				Status: "creating",
			})
			if err != nil {
				t.Fatalf("marshal snapshot: %v", err)
			}
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Operation{
				ID:               "op-1",
				ResourceID:       optionalString("proxy-1"),
				State:            cloudminiV3OperationSucceeded,
				ResourceSnapshot: rawSnapshot,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/proxies/proxy-1":
			proxyReads++
			status := "creating"
			if proxyReads > 1 {
				status = "running"
			}
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Proxy{
				ID:         "proxy-1",
				Kind:       "ipv4_dc",
				Status:     status,
				Host:       "203.0.113.10",
				PortSocks:  1080,
				Username:   "proxy-user",
				Password:   "proxy-pass",
				OutboundIP: "203.0.113.10",
			})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := newTestCloudminiV3Adapter(t, server.URL, &cloudminiV3TestCipher{}, CloudminiV3SourceConfig{
		Kind:     "ipv4_dc",
		GroupID:  "group-1",
		Protocol: "socks5",
	})
	result, err := adapter.Provision(context.Background(), validOperation(), ProvisionRequest{PlanKey: "proxy"})
	if err != nil {
		t.Fatalf("expected provision success after usable status: %v", err)
	}
	if proxyReads < 2 || result.Status != OperationStatusSuccess || result.ProviderStatus != "running" ||
		result.ExternalResourceID != "proxy-1" || !result.Credential.HasEncryptedPayload() {
		t.Fatalf("unexpected waited result: reads=%d result=%+v", proxyReads, result)
	}
}

func TestCloudminiV3AdapterProvisionRejectsUsableProxyWithoutCredential(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiAuth(t, r)
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/proxies":
			writeCloudminiSuccess(t, w, http.StatusAccepted, cloudminiV3ProxyMutationResponse{
				Resource:  cloudminiV3Resource{ID: "proxy-1", Kind: "ipv4_dc", Status: "creating"},
				Operation: cloudminiV3Operation{ID: "op-1", State: cloudminiV3OperationAccepted},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-1":
			rawSnapshot, err := json.Marshal(cloudminiV3Proxy{
				ID:     "proxy-1",
				Kind:   "ipv4_dc",
				Status: "creating",
			})
			if err != nil {
				t.Fatalf("marshal snapshot: %v", err)
			}
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Operation{
				ID:               "op-1",
				ResourceID:       optionalString("proxy-1"),
				State:            cloudminiV3OperationSucceeded,
				ResourceSnapshot: rawSnapshot,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/proxies/proxy-1":
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Proxy{
				ID:         "proxy-1",
				Kind:       "ipv4_dc",
				Status:     "running",
				Host:       "203.0.113.10",
				PortSocks:  1080,
				OutboundIP: "203.0.113.10",
			})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := newTestCloudminiV3Adapter(t, server.URL, &cloudminiV3TestCipher{}, CloudminiV3SourceConfig{
		Kind:     "ipv4_dc",
		GroupID:  "group-1",
		Protocol: "socks5",
	})
	result, err := adapter.Provision(context.Background(), validOperation(), ProvisionRequest{PlanKey: "proxy"})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorCredentialMissing {
		t.Fatalf("expected credential-missing error, got result=%+v err=%v", result, err)
	}
	if result.Credential.HasEncryptedPayload() {
		t.Fatalf("credential must not be returned when provider omits auth fields: %+v", result.Credential)
	}
}

func TestCloudminiV3AdapterTerminateUsesDeleteAndIdempotency(t *testing.T) {
	var sawDelete bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiAuth(t, r)
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v3/proxies/proxy-1":
			sawDelete = true
			if got := r.Header.Get("Idempotency-Key"); got != "tenant_a:order_1:item_1" {
				t.Fatalf("expected idempotency header, got %q", got)
			}
			writeCloudminiSuccess(t, w, http.StatusAccepted, cloudminiV3ProxyMutationResponse{
				Resource:  cloudminiV3Resource{ID: "proxy-1", Kind: "ipv4_dc", Status: "deleting"},
				Operation: cloudminiV3Operation{ID: "op-delete-1", State: cloudminiV3OperationAccepted},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-delete-1":
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Operation{
				ID:         "op-delete-1",
				ResourceID: optionalString("proxy-1"),
				State:      cloudminiV3OperationSucceeded,
			})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := newTestCloudminiV3Adapter(t, server.URL, &cloudminiV3TestCipher{}, CloudminiV3SourceConfig{
		Kind:     "ipv4_dc",
		GroupID:  "group-1",
		Protocol: "socks5",
	})
	result, err := adapter.Terminate(context.Background(), validOperation(), ResourceRequest{ExternalResourceID: "proxy-1"})
	if err != nil {
		t.Fatalf("expected terminate success: %v", err)
	}
	if !sawDelete ||
		result.Status != OperationStatusSuccess ||
		result.ExternalRequestID != "op-delete-1" ||
		result.ExternalResourceID != "proxy-1" ||
		result.ProviderStatus != string(cloudminiV3OperationSucceeded) {
		t.Fatalf("unexpected terminate result: sawDelete=%v result=%+v", sawDelete, result)
	}
}

func TestCloudminiV3AdapterTerminateTimeoutAfterAcceptedRequiresManualReview(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiAuth(t, r)
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v3/proxies/proxy-1":
			writeCloudminiSuccess(t, w, http.StatusAccepted, cloudminiV3ProxyMutationResponse{
				Resource:  cloudminiV3Resource{ID: "proxy-1", Kind: "ipv4_dc", Status: "deleting"},
				Operation: cloudminiV3Operation{ID: "op-delete-1", State: cloudminiV3OperationAccepted},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-delete-1":
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Operation{
				ID:         "op-delete-1",
				ResourceID: optionalString("proxy-1"),
				State:      cloudminiV3OperationRunning,
			})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := newTestCloudminiV3Adapter(t, server.URL, &cloudminiV3TestCipher{}, CloudminiV3SourceConfig{
		Kind:     "ipv4_dc",
		GroupID:  "group-1",
		Protocol: "socks5",
	})
	adapter.pollTimeout = 3 * adapter.pollInterval

	result, err := adapter.Terminate(context.Background(), validOperation(), ResourceRequest{ExternalResourceID: "proxy-1"})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorTimeoutRequestKnown {
		t.Fatalf("expected request-known timeout, got result=%+v err=%v", result, err)
	}
	if result.Status != OperationStatusUnknown ||
		result.RetrySafety != RetrySafetyManualReviewRequired ||
		result.ExternalRequestID != "op-delete-1" ||
		result.ExternalResourceID != "proxy-1" {
		t.Fatalf("unexpected terminate timeout result: %+v", result)
	}
}
