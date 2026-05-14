package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	cloudminiV3OperationAccepted  cloudminiV3OperationState = "accepted"
	cloudminiV3OperationRunning   cloudminiV3OperationState = "running"
	cloudminiV3OperationSucceeded cloudminiV3OperationState = "succeeded"
	cloudminiV3OperationFailed    cloudminiV3OperationState = "failed"
	cloudminiV3OperationTimedOut  cloudminiV3OperationState = "timed_out"
	cloudminiV3OperationCancelled cloudminiV3OperationState = "cancelled"
)

type cloudminiV3HTTPDoer interface {
	Do(request *http.Request) (*http.Response, error)
}

type cloudminiV3Client struct {
	baseURL *url.URL
	token   string
	http    cloudminiV3HTTPDoer
}

type cloudminiV3Envelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   json.RawMessage `json:"error"`
}

type cloudminiV3APIError struct {
	Code      string          `json:"code"`
	Message   string          `json:"message"`
	Retryable bool            `json:"retryable"`
	Details   json.RawMessage `json:"details"`
}

type cloudminiV3Capabilities struct {
	Features map[string]bool `json:"features"`
}

type cloudminiV3GroupInventory struct {
	ID               string `json:"id"`
	Kind             string `json:"kind"`
	Name             string `json:"name"`
	SellState        string `json:"sell_state"`
	AllocatableUnits int    `json:"allocatable_units"`
}

type cloudminiV3OperationState string

type cloudminiV3Operation struct {
	ID               string                    `json:"id"`
	Action           string                    `json:"action"`
	ResourceID       *string                   `json:"resource_id"`
	ProxyKind        string                    `json:"proxy_kind"`
	State            cloudminiV3OperationState `json:"state"`
	DesiredStatus    *string                   `json:"desired_status"`
	ExternalRef      *string                   `json:"external_ref"`
	ErrorCode        *string                   `json:"error_code"`
	ErrorMessage     *string                   `json:"error_message"`
	ResourceSnapshot json.RawMessage           `json:"resource_snapshot"`
}

type cloudminiV3CreateProxyRequest struct {
	Kind             string  `json:"kind"`
	GroupID          string  `json:"group_id"`
	NodeID           *string `json:"node_id,omitempty"`
	Protocol         string  `json:"protocol"`
	BandwidthLimitMB int     `json:"bandwidth_limit_mb,omitempty"`
	SpeedLimitMBps   int     `json:"speed_limit_mbps,omitempty"`
	ExternalRef      *string `json:"external_ref,omitempty"`
}

type cloudminiV3ProxyMutationResponse struct {
	Resource  cloudminiV3Resource  `json:"resource"`
	Operation cloudminiV3Operation `json:"operation"`
}

type cloudminiV3Resource struct {
	ID            string  `json:"id"`
	Kind          string  `json:"kind"`
	GroupID       string  `json:"group_id"`
	NodeID        *string `json:"node_id"`
	Status        string  `json:"status"`
	DesiredStatus string  `json:"desired_status"`
}

type cloudminiV3Proxy struct {
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	NodeID        string `json:"node_id"`
	Status        string `json:"status"`
	Host          string `json:"host"`
	OutboundIP    string `json:"outbound_ip"`
	PortSocks     int    `json:"port_socks"`
	PortHTTP      int    `json:"port_http"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	ConnectionURI string `json:"connection_uri"`
}

func newCloudminiV3Client(baseURL string, token string, httpClient cloudminiV3HTTPDoer) (*cloudminiV3Client, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return nil, NewError(ErrorConfigInvalid, "cloudmini v3 base url is missing")
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, NewError(ErrorConfigInvalid, "cloudmini v3 base url is invalid")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, NewError(ErrorConfigInvalid, "cloudmini v3 base url scheme is invalid")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, NewError(ErrorAuthFailed, "cloudmini v3 api credential is missing")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &cloudminiV3Client{baseURL: parsed, token: token, http: httpClient}, nil
}

func (client *cloudminiV3Client) getCapabilities(ctx context.Context, operation OperationContext) (cloudminiV3Capabilities, error) {
	var response cloudminiV3Capabilities
	err := client.do(ctx, http.MethodGet, "/api/v3/capabilities", "", operation, nil, &response)
	return response, err
}

func (client *cloudminiV3Client) listGroupInventory(ctx context.Context, operation OperationContext, kind string) ([]cloudminiV3GroupInventory, error) {
	var response []cloudminiV3GroupInventory
	path := "/api/v3/inventory/groups?kind=" + url.QueryEscape(kind)
	err := client.do(ctx, http.MethodGet, path, "", operation, nil, &response)
	return response, err
}

func (client *cloudminiV3Client) createProxy(ctx context.Context, operation OperationContext, body cloudminiV3CreateProxyRequest) (cloudminiV3ProxyMutationResponse, error) {
	var response cloudminiV3ProxyMutationResponse
	err := client.do(ctx, http.MethodPost, "/api/v3/proxies", string(operation.IdempotencyKey), operation, body, &response)
	return response, err
}

func (client *cloudminiV3Client) getOperation(ctx context.Context, operation OperationContext, operationID string) (cloudminiV3Operation, error) {
	var response cloudminiV3Operation
	err := client.do(ctx, http.MethodGet, "/api/v3/operations/"+url.PathEscape(operationID), "", operation, nil, &response)
	return response, err
}

func (client *cloudminiV3Client) getProxy(ctx context.Context, operation OperationContext, resourceID string) (cloudminiV3Proxy, error) {
	var response cloudminiV3Proxy
	err := client.do(ctx, http.MethodGet, "/api/v3/proxies/"+url.PathEscape(resourceID), "", operation, nil, &response)
	return response, err
}

func (client *cloudminiV3Client) deleteProxy(ctx context.Context, operation OperationContext, resourceID string) (cloudminiV3ProxyMutationResponse, error) {
	var response cloudminiV3ProxyMutationResponse
	err := client.do(ctx, http.MethodDelete, "/api/v3/proxies/"+url.PathEscape(resourceID), string(operation.IdempotencyKey), operation, nil, &response)
	return response, err
}

func (client *cloudminiV3Client) proxyAction(ctx context.Context, operation OperationContext, resourceID string, action string) (cloudminiV3ProxyMutationResponse, error) {
	var response cloudminiV3ProxyMutationResponse
	path := "/api/v3/proxies/" + url.PathEscape(resourceID) + "/actions/" + url.PathEscape(action)
	err := client.do(ctx, http.MethodPost, path, string(operation.IdempotencyKey), operation, nil, &response)
	return response, err
}

func (client *cloudminiV3Client) do(
	ctx context.Context,
	method string,
	path string,
	idempotencyKey string,
	operation OperationContext,
	body interface{},
	output interface{},
) error {
	var payload io.Reader
	if body != nil {
		buffer := &bytes.Buffer{}
		if err := json.NewEncoder(buffer).Encode(body); err != nil {
			return AdapterError{Code: ErrorConfigInvalid, MessageRedacted: "cloudmini v3 request payload is invalid", Cause: err}
		}
		payload = buffer
	}
	requestURL := client.resolve(path)
	request, err := http.NewRequestWithContext(ctx, method, requestURL, payload)
	if err != nil {
		return AdapterError{Code: ErrorConfigInvalid, MessageRedacted: "cloudmini v3 request is invalid", Cause: err}
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+client.token)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if idempotencyKey != "" {
		request.Header.Set("Idempotency-Key", idempotencyKey)
	}
	if operation.CorrelationID != "" {
		request.Header.Set("X-Request-ID", string(operation.CorrelationID))
	}
	response, err := client.http.Do(request)
	if err != nil {
		return client.transportError(method, err)
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return AdapterError{Code: ErrorResponseInvalid, MessageRedacted: "cloudmini v3 response could not be read", Cause: err}
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return cloudminiV3ErrorFromStatus(response.StatusCode, bodyBytes)
	}
	var envelope cloudminiV3Envelope
	if err := json.Unmarshal(bodyBytes, &envelope); err != nil {
		return AdapterError{Code: ErrorResponseInvalid, MessageRedacted: "cloudmini v3 response envelope is invalid", Cause: err}
	}
	if !envelope.Success {
		return cloudminiV3ErrorFromEnvelope(http.StatusOK, envelope.Error)
	}
	if output == nil {
		return nil
	}
	if len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return AdapterError{Code: ErrorResponseInvalid, MessageRedacted: "cloudmini v3 response data is missing"}
	}
	if err := json.Unmarshal(envelope.Data, output); err != nil {
		return AdapterError{Code: ErrorResponseInvalid, MessageRedacted: "cloudmini v3 response data is invalid", Cause: err}
	}
	return nil
}

func (client *cloudminiV3Client) resolve(path string) string {
	resolved := *client.baseURL
	relativePath := path
	if queryIndex := strings.Index(relativePath, "?"); queryIndex >= 0 {
		resolved.RawQuery = relativePath[queryIndex+1:]
		relativePath = relativePath[:queryIndex]
	} else {
		resolved.RawQuery = ""
	}
	if strings.HasPrefix(relativePath, "/") {
		resolved.Path = strings.TrimRight(client.baseURL.Path, "/") + relativePath
	} else {
		resolved.Path = strings.TrimRight(client.baseURL.Path, "/") + "/" + relativePath
	}
	return resolved.String()
}

func (client *cloudminiV3Client) transportError(method string, err error) AdapterError {
	if method == http.MethodPost || method == http.MethodDelete || method == http.MethodPatch || method == http.MethodPut {
		return AdapterError{
			Code:            ErrorTimeoutUnknown,
			MessageRedacted: "cloudmini v3 mutation result is unknown",
			Safety:          RetrySafetyUnsafeRetry,
			Cause:           err,
		}
	}
	return AdapterError{
		Code:            ErrorTemporary,
		MessageRedacted: "cloudmini v3 read request failed",
		Safety:          RetrySafetySafeRetry,
		Cause:           err,
	}
}

func cloudminiV3ErrorFromStatus(statusCode int, body []byte) AdapterError {
	return cloudminiV3ErrorFromEnvelope(statusCode, extractCloudminiV3ErrorRaw(body))
}

func cloudminiV3ErrorFromEnvelope(statusCode int, raw json.RawMessage) AdapterError {
	apiErr := parseCloudminiV3APIError(raw)
	code := mapCloudminiV3ErrorCode(statusCode, apiErr.Code)
	return AdapterError{
		Code:            code,
		MessageRedacted: cloudminiV3RedactedMessage(code, apiErr.Code),
		Safety:          DefaultRetrySafety(code),
	}
}

func extractCloudminiV3ErrorRaw(body []byte) json.RawMessage {
	var envelope cloudminiV3Envelope
	if err := json.Unmarshal(body, &envelope); err == nil && len(envelope.Error) > 0 {
		return envelope.Error
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(body, &generic); err == nil && len(generic["error"]) > 0 {
		return generic["error"]
	}
	return nil
}

func parseCloudminiV3APIError(raw json.RawMessage) cloudminiV3APIError {
	if len(raw) == 0 || string(raw) == "null" {
		return cloudminiV3APIError{}
	}
	var apiErr cloudminiV3APIError
	if err := json.Unmarshal(raw, &apiErr); err == nil && apiErr.Code != "" {
		return apiErr
	}
	var message string
	if err := json.Unmarshal(raw, &message); err == nil {
		return cloudminiV3APIError{Message: message}
	}
	return cloudminiV3APIError{}
}

func mapCloudminiV3ErrorCode(statusCode int, providerCode string) ErrorCode {
	switch providerCode {
	case "CAPACITY_EXHAUSTED":
		return ErrorOutOfStock
	case "IDEMPOTENCY_CONFLICT", "INVALID_STATE_TRANSITION":
		return ErrorStateDrift
	case "RESERVATION_NOT_FOUND", "RESERVATION_EXPIRED", "RESERVATION_ALREADY_CONSUMED":
		return ErrorConfigInvalid
	case "OPERATION_NOT_FOUND", "PROXY_NOT_FOUND":
		return ErrorStateDrift
	case "INVALID_ACTION":
		return ErrorCapabilityNotSupported
	case "INVALID_INPUT":
		return ErrorConfigInvalid
	case "INTERNAL_ERROR":
		return ErrorTemporary
	}
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrorAuthFailed
	case http.StatusForbidden:
		return ErrorPermissionDenied
	case http.StatusTooManyRequests:
		return ErrorRateLimited
	case http.StatusConflict:
		return ErrorStateDrift
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return ErrorConfigInvalid
	case http.StatusNotFound:
		return ErrorStateDrift
	default:
		if statusCode >= 500 {
			return ErrorTemporary
		}
		return ErrorResponseInvalid
	}
}

func cloudminiV3RedactedMessage(code ErrorCode, providerCode string) string {
	if providerCode == "" {
		return fmt.Sprintf("cloudmini v3 request failed with %s", code)
	}
	return fmt.Sprintf("cloudmini v3 request failed with %s", providerCode)
}

func (operation cloudminiV3Operation) proxySnapshot() (cloudminiV3Proxy, error) {
	var proxy cloudminiV3Proxy
	if len(operation.ResourceSnapshot) == 0 || string(operation.ResourceSnapshot) == "null" {
		return cloudminiV3Proxy{}, NewError(ErrorCredentialMissing, "cloudmini v3 operation snapshot is missing")
	}
	if err := json.Unmarshal(operation.ResourceSnapshot, &proxy); err != nil {
		return cloudminiV3Proxy{}, AdapterError{Code: ErrorResponseInvalid, MessageRedacted: "cloudmini v3 operation snapshot is invalid", Cause: err}
	}
	return proxy, nil
}

func (proxy cloudminiV3Proxy) serviceIdentifier() string {
	host := strings.TrimSpace(proxy.Host)
	if host == "" {
		host = strings.TrimSpace(proxy.OutboundIP)
	}
	if host == "" {
		return strings.TrimSpace(proxy.ID)
	}
	if proxy.PortSocks > 0 {
		return host + ":" + strconv.Itoa(proxy.PortSocks)
	}
	if proxy.PortHTTP > 0 {
		return host + ":" + strconv.Itoa(proxy.PortHTTP)
	}
	return host
}

func (proxy cloudminiV3Proxy) maskedHint() string {
	identifier := proxy.serviceIdentifier()
	if identifier == "" {
		return "Proxy credential"
	}
	return "Proxy access for " + identifier
}
