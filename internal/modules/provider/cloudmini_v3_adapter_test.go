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

func TestCloudminiV3AdapterRoutesSourcesToDistinctEndpoints(t *testing.T) {
	var sawSourceA, sawSourceB bool
	serverA := newCloudminiProvisionServer(t, "token-a", "ipv4_dc", "group-a", "proxy-a", "203.0.113.11", &sawSourceA)
	defer serverA.Close()
	serverB := newCloudminiProvisionServer(t, "token-b", "residential", "group-b", "proxy-b", "203.0.113.12", &sawSourceB)
	defer serverB.Close()

	adapter, err := NewCloudminiV3Adapter(CloudminiV3Config{
		CredentialCipher: &cloudminiV3TestCipher{},
		SourceEndpoints: map[SourceID]CloudminiV3EndpointConfig{
			"source-a": {
				BaseURL:  serverA.URL,
				APIToken: "token-a",
				Source:   CloudminiV3SourceConfig{Kind: "ipv4_dc", GroupID: "group-a", Protocol: "socks5"},
			},
			"source-b": {
				BaseURL:  serverB.URL,
				APIToken: "token-b",
				Source:   CloudminiV3SourceConfig{Kind: "residential", GroupID: "group-b", Protocol: "http"},
			},
		},
		PollInterval: time.Millisecond,
		PollTimeout:  50 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("create adapter: %v", err)
	}

	operationA := validOperation()
	operationA.SourceID = "source-a"
	resultA, err := adapter.Provision(context.Background(), operationA, ProvisionRequest{PlanKey: "proxy"})
	if err != nil {
		t.Fatalf("expected source A provision success: %v", err)
	}
	if resultA.ExternalResourceID != "proxy-a" || resultA.ServiceIdentifier != "203.0.113.11:1080" {
		t.Fatalf("unexpected source A result: %+v", resultA)
	}

	operationB := validOperation()
	operationB.SourceID = "source-b"
	resultB, err := adapter.Provision(context.Background(), operationB, ProvisionRequest{PlanKey: "proxy"})
	if err != nil {
		t.Fatalf("expected source B provision success: %v", err)
	}
	if resultB.ExternalResourceID != "proxy-b" || resultB.ServiceIdentifier != "203.0.113.12:1080" {
		t.Fatalf("unexpected source B result: %+v", resultB)
	}
	if !sawSourceA || !sawSourceB {
		t.Fatalf("expected both endpoints to receive provision calls: sourceA=%v sourceB=%v", sawSourceA, sawSourceB)
	}
}

func TestCloudminiV3AdapterMissingEndpointMappingDoesNotCallProvider(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		t.Fatalf("provider endpoint must not be called for missing source mapping")
	}))
	defer server.Close()

	adapter, err := NewCloudminiV3Adapter(CloudminiV3Config{
		CredentialCipher: &cloudminiV3TestCipher{},
		SourceEndpoints: map[SourceID]CloudminiV3EndpointConfig{
			"source-a": {
				BaseURL:  server.URL,
				APIToken: "token-a",
				Source:   CloudminiV3SourceConfig{Kind: "ipv4_dc", GroupID: "group-a", Protocol: "socks5"},
			},
		},
	})
	if err != nil {
		t.Fatalf("create adapter: %v", err)
	}

	operation := validOperation()
	operation.SourceID = "missing-source"
	result, err := adapter.Provision(context.Background(), operation, ProvisionRequest{PlanKey: "proxy"})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorConfigInvalid {
		t.Fatalf("expected config error, got result=%+v err=%v", result, err)
	}
	if called {
		t.Fatal("provider endpoint was called despite missing source mapping")
	}
}

func TestCloudminiV3AdapterRoutesProviderAccountEndpoint(t *testing.T) {
	var sawInventory bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiBearer(t, r, "account-token")
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/inventory/groups" || r.URL.Query().Get("kind") != "ipv4_dc" {
			t.Fatalf("unexpected request %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
		sawInventory = true
		writeCloudminiSuccess(t, w, http.StatusOK, []cloudminiV3GroupInventory{
			{ID: "group-account", Kind: "ipv4_dc", AllocatableUnits: 7},
		})
	}))
	defer server.Close()

	adapter, err := NewCloudminiV3Adapter(CloudminiV3Config{
		CredentialCipher: &cloudminiV3TestCipher{},
		AccountEndpoints: map[AccountID]CloudminiV3EndpointConfig{
			"account-a": {
				BaseURL:  server.URL,
				APIToken: "account-token",
				Source:   CloudminiV3SourceConfig{Kind: "ipv4_dc", GroupID: "group-account", Protocol: "socks5"},
			},
		},
	})
	if err != nil {
		t.Fatalf("create adapter: %v", err)
	}

	operation := validOperation()
	operation.SourceID = ""
	operation.ProviderAccountID = "account-a"
	result, err := adapter.CheckStock(context.Background(), operation, StockRequest{PlanKey: "proxy"})
	if err != nil {
		t.Fatalf("expected account endpoint stock success: %v", err)
	}
	if !sawInventory || result.StockStatus != StockStatusAvailable || result.CapacityCount != 7 {
		t.Fatalf("unexpected account endpoint result: saw=%v result=%+v", sawInventory, result)
	}
}

func TestCloudminiV3AdapterRejectsAccountFallbackWhenSourceMismatches(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		t.Fatalf("provider endpoint must not be called for source/account mismatch")
	}))
	defer server.Close()

	adapter, err := NewCloudminiV3Adapter(CloudminiV3Config{
		CredentialCipher: &cloudminiV3TestCipher{},
		AccountEndpoints: map[AccountID]CloudminiV3EndpointConfig{
			"account-a": {
				BaseURL:  server.URL,
				APIToken: "account-token",
				Source:   CloudminiV3SourceConfig{Kind: "ipv4_dc", GroupID: "group-account", Protocol: "socks5"},
			},
		},
	})
	if err != nil {
		t.Fatalf("create adapter: %v", err)
	}

	operation := validOperation()
	operation.SourceID = "source-without-endpoint"
	operation.ProviderAccountID = "account-a"
	result, err := adapter.CheckStock(context.Background(), operation, StockRequest{PlanKey: "proxy"})
	var adapterErr AdapterError
	if !errors.As(err, &adapterErr) || adapterErr.Code != ErrorConfigInvalid {
		t.Fatalf("expected config error, got result=%+v err=%v", result, err)
	}
	if called {
		t.Fatal("provider endpoint was called despite source/account mismatch")
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
		SourceConfigs: map[SourceID]CloudminiV3SourceConfig{
			validOperation().SourceID: source,
		},
		PollInterval: time.Millisecond,
		PollTimeout:  50 * time.Millisecond,
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
	assertCloudminiBearer(t, r, "test-token")
}

func assertCloudminiBearer(t *testing.T, r *http.Request, token string) {
	t.Helper()
	if got := r.Header.Get("Authorization"); got != "Bearer "+token {
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

func newCloudminiProvisionServer(t *testing.T, token string, kind string, groupID string, proxyID string, host string, sawCreate *bool) *httptest.Server {
	t.Helper()
	operationID := "op-" + proxyID
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCloudminiBearer(t, r, token)
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/proxies":
			*sawCreate = true
			var payload cloudminiV3CreateProxyRequest
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode request: %v", err)
			}
			if payload.Kind != kind || payload.GroupID != groupID {
				t.Fatalf("unexpected create payload: %+v", payload)
			}
			writeCloudminiSuccess(t, w, http.StatusAccepted, cloudminiV3ProxyMutationResponse{
				Resource:  cloudminiV3Resource{ID: proxyID, Kind: kind, Status: "provisioning"},
				Operation: cloudminiV3Operation{ID: operationID, State: cloudminiV3OperationAccepted},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/"+operationID:
			snapshot := cloudminiV3Proxy{
				ID:         proxyID,
				Kind:       kind,
				Status:     "running",
				Host:       host,
				OutboundIP: host,
				PortSocks:  1080,
				Username:   "proxy-user",
				Password:   "proxy-pass",
			}
			rawSnapshot, err := json.Marshal(snapshot)
			if err != nil {
				t.Fatalf("marshal snapshot: %v", err)
			}
			writeCloudminiSuccess(t, w, http.StatusOK, cloudminiV3Operation{
				ID:               operationID,
				ResourceID:       optionalString(proxyID),
				State:            cloudminiV3OperationSucceeded,
				ResourceSnapshot: rawSnapshot,
			})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
}
