package provider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCloudminiV3AdapterProvisionEncryptsCredentialAndUsesIdempotency(t *testing.T) {
	cipher := &cloudminiV3TestCipher{}
	var sawCreate bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiAuth(t, r)
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/proxies":
			sawCreate = true
			if got := r.Header.Get("Idempotency-Key"); got != "tenant_a:order_1:item_1" {
				t.Fatalf("expected idempotency header, got %q", got)
			}
			var payload cloudminiV3CreateProxyRequest
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			if payload.Kind != "ipv4_dc" || payload.GroupID != "group-1" || payload.Protocol != "socks5" {
				t.Fatalf("unexpected create payload: %+v", payload)
			}
			writeCloudminiSuccess(t, w, http.StatusAccepted, cloudminiV3ProxyMutationResponse{
				Resource:  cloudminiV3Resource{ID: "proxy-1", Kind: "ipv4_dc", Status: "provisioning"},
				Operation: cloudminiV3Operation{ID: "op-1", State: cloudminiV3OperationAccepted},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-1":
			snapshot := cloudminiV3Proxy{
				ID:            "proxy-1",
				Kind:          "ipv4_dc",
				Status:        "running",
				Host:          "203.0.113.10",
				OutboundIP:    "203.0.113.10",
				PortSocks:     1080,
				Username:      "proxy-user",
				Password:      "proxy-pass",
				ConnectionURI: "socks5://proxy-user:proxy-pass@203.0.113.10:1080",
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
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := newTestCloudminiV3Adapter(t, server.URL, cipher, CloudminiV3SourceConfig{
		Kind:     "ipv4_dc",
		GroupID:  "group-1",
		Protocol: "socks5",
	})
	result, err := adapter.Provision(context.Background(), validOperation(), ProvisionRequest{PlanKey: "proxy"})
	if err != nil {
		t.Fatalf("expected provision success: %v", err)
	}
	if !sawCreate {
		t.Fatal("expected create request")
	}
	if result.Status != OperationStatusSuccess ||
		result.ExternalRequestID != "op-1" ||
		result.ExternalResourceID != "proxy-1" ||
		result.ServiceIdentifier != "203.0.113.10:1080" {
		t.Fatalf("unexpected result: %+v", result)
	}
	if result.Credential.Type != CredentialTypeProxyAuth || result.Credential.EncryptedPayload != "encrypted-proxy-payload" {
		t.Fatalf("unexpected credential envelope: %+v", result.Credential)
	}
	if strings.Contains(result.Credential.EncryptedPayload, "proxy-pass") {
		t.Fatalf("encrypted payload must not contain plaintext credential: %s", result.Credential.EncryptedPayload)
	}
	if len(cipher.plaintexts) != 1 || !strings.Contains(cipher.plaintexts[0], "proxy-pass") {
		t.Fatalf("expected plaintext to be passed only into cipher, got %#v", cipher.plaintexts)
	}
}

func TestCloudminiV3AdapterCheckStockMapsOutOfStock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiAuth(t, r)
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/inventory/groups" || r.URL.Query().Get("kind") != "residential" {
			t.Fatalf("unexpected request %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
		writeCloudminiSuccess(t, w, http.StatusOK, []cloudminiV3GroupInventory{
			{ID: "group-1", Kind: "residential", AllocatableUnits: 0},
		})
	}))
	defer server.Close()

	adapter := newTestCloudminiV3Adapter(t, server.URL, &cloudminiV3TestCipher{}, CloudminiV3SourceConfig{
		Kind:     "residential",
		GroupID:  "group-1",
		Protocol: "http",
	})
	result, err := adapter.CheckStock(context.Background(), validOperation(), StockRequest{PlanKey: "proxy"})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorOutOfStock {
		t.Fatalf("expected out-of-stock error, got result=%+v err=%v", result, err)
	}
	if result.StockStatus != StockStatusOutOfStock || result.Result.ErrorCode != ErrorOutOfStock {
		t.Fatalf("unexpected stock result: %+v", result)
	}
}

func TestCloudminiV3AdapterProvisionTimeoutAfterAcceptedRequiresManualReview(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiAuth(t, r)
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/proxies":
			writeCloudminiSuccess(t, w, http.StatusAccepted, cloudminiV3ProxyMutationResponse{
				Resource:  cloudminiV3Resource{ID: "proxy-1", Kind: "ipv4_dc", Status: "provisioning"},
				Operation: cloudminiV3Operation{ID: "op-1", State: cloudminiV3OperationAccepted},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-1":
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Operation{ID: "op-1", State: cloudminiV3OperationRunning})
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
	adapter.pollTimeout = 3 * time.Millisecond

	result, err := adapter.Provision(context.Background(), validOperation(), ProvisionRequest{PlanKey: "proxy"})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorTimeoutRequestKnown {
		t.Fatalf("expected request-known timeout, got result=%+v err=%v", result, err)
	}
	if result.Status != OperationStatusUnknown ||
		result.RetrySafety != RetrySafetyManualReviewRequired ||
		result.ExternalRequestID != "op-1" ||
		result.ExternalResourceID != "proxy-1" {
		t.Fatalf("unexpected timeout result: %+v", result)
	}
}

func TestCloudminiV3AdapterHealthMapsAuthDeniedDown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid or expired access"}`))
	}))
	defer server.Close()

	adapter := newTestCloudminiV3Adapter(t, server.URL, &cloudminiV3TestCipher{}, CloudminiV3SourceConfig{
		Kind:     "ipv4_dc",
		GroupID:  "group-1",
		Protocol: "socks5",
	})
	result, err := adapter.CheckHealth(context.Background(), validOperation(), HealthRequest{})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorAuthFailed {
		t.Fatalf("expected auth error, got result=%+v err=%v", result, err)
	}
	if result.HealthStatus != HealthStatusDown || result.Result.RetrySafety != RetrySafetyDoNotRetry {
		t.Fatalf("unexpected health result: %+v", result)
	}
}

func TestCloudminiV3AdapterChangeIPRequiresResidentialSource(t *testing.T) {
	adapter := newTestCloudminiV3Adapter(t, "http://cloudmini.example", &cloudminiV3TestCipher{}, CloudminiV3SourceConfig{
		Kind:     "ipv4_dc",
		GroupID:  "group-1",
		Protocol: "socks5",
	})

	result, err := adapter.ChangeIP(context.Background(), validOperation(), ResourceRequest{ExternalResourceID: "proxy-1"})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorCapabilityNotSupported {
		t.Fatalf("expected unsupported change-ip, got result=%+v err=%v", result, err)
	}
	if result.Status != OperationStatusCapabilityNotSupported {
		t.Fatalf("unexpected unsupported result: %+v", result)
	}
}

func TestCloudminiV3AdapterProvisionFailsClosedWithoutSourceMapping(t *testing.T) {
	adapter := newTestCloudminiV3Adapter(t, "http://cloudmini.example", &cloudminiV3TestCipher{}, CloudminiV3SourceConfig{})

	result, err := adapter.Provision(context.Background(), validOperation(), ProvisionRequest{PlanKey: "proxy"})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorConfigInvalid {
		t.Fatalf("expected config error, got result=%+v err=%v", result, err)
	}
	if result.RetrySafety != RetrySafetyDoNotRetry {
		t.Fatalf("expected do-not-retry config failure, got %+v", result)
	}
}

func newTestCloudminiV3Adapter(t *testing.T, baseURL string, cipher CredentialCipher, source CloudminiV3SourceConfig) *CloudminiV3Adapter {
	t.Helper()
	adapter, err := NewCloudminiV3Adapter(CloudminiV3Config{
		BaseURL:          baseURL,
		APIToken:         "test-token",
		CredentialCipher: cipher,
		DefaultSource:    source,
		PollInterval:     time.Millisecond,
		PollTimeout:      50 * time.Millisecond,
		Now: func() time.Time {
			return time.Date(2026, 5, 14, 1, 2, 3, 0, time.UTC)
		},
	})
	if err != nil {
		t.Fatalf("create adapter: %v", err)
	}
	return adapter
}

func assertCloudminiAuth(t *testing.T, r *http.Request) {
	t.Helper()
	if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
		t.Fatalf("expected bearer auth header, got %q", got)
	}
}

func writeCloudminiSuccess(t *testing.T, w http.ResponseWriter, status int, data interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
		"meta": map[string]interface{}{
			"timestamp": "2026-05-14T01:02:03Z",
		},
	}); err != nil {
		t.Fatalf("write response: %v", err)
	}
}

type cloudminiV3TestCipher struct {
	plaintexts []string
}

func (cipher *cloudminiV3TestCipher) Encrypt(plaintext string) (string, error) {
	cipher.plaintexts = append(cipher.plaintexts, plaintext)
	return "encrypted-proxy-payload", nil
}
