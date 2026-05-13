package notification

import (
	"encoding/json"
	"strings"
)

const redactedValue = "[redacted]"

var sensitiveKeyFragments = []string{
	"authorization",
	"api_key",
	"apikey",
	"cookie",
	"credential",
	"evidence_private",
	"password",
	"passwd",
	"private_key",
	"provider_response",
	"raw_response",
	"reset_link",
	"reset_token",
	"secret",
	"session",
	"token",
}

func RedactPayload(payload json.RawMessage) (json.RawMessage, error) {
	payload = defaultPayload(payload)
	var value any
	if err := json.Unmarshal(payload, &value); err != nil {
		return nil, ErrPayloadInvalid
	}
	redacted := redactValue(value)
	body, err := json.Marshal(redacted)
	if err != nil {
		return nil, ErrPayloadInvalid
	}
	return json.RawMessage(body), nil
}

func redactValue(value any) any {
	switch current := value.(type) {
	case map[string]any:
		output := make(map[string]any, len(current))
		for key, child := range current {
			if sensitivePayloadKey(key) {
				output[key] = redactedValue
				continue
			}
			output[key] = redactValue(child)
		}
		return output
	case []any:
		output := make([]any, 0, len(current))
		for _, child := range current {
			output = append(output, redactValue(child))
		}
		return output
	default:
		return value
	}
}

func sensitivePayloadKey(key string) bool {
	normalized := strings.ToLower(strings.TrimSpace(key))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	for _, fragment := range sensitiveKeyFragments {
		if strings.Contains(normalized, fragment) {
			return true
		}
	}
	return false
}
