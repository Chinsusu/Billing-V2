package order

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type serviceCredentialResponse struct {
	ID             CredentialID     `json:"id"`
	Type           CredentialType   `json:"credential_type"`
	MaskedHint     string           `json:"masked_hint"`
	Status         CredentialStatus `json:"status"`
	LastRevealedAt *time.Time       `json:"last_revealed_at,omitempty"`
}

type credentialRevealRequest struct {
	Reason string `json:"reason"`
}

type credentialRevealResponse struct {
	ID                   CredentialID     `json:"id"`
	Type                 CredentialType   `json:"credential_type"`
	MaskedHint           string           `json:"masked_hint"`
	Status               CredentialStatus `json:"status"`
	Payload              json.RawMessage  `json:"payload"`
	RevealedAt           time.Time        `json:"revealed_at"`
	RevealExpiresMessage string           `json:"reveal_expires_message"`
}

func newServiceCredentialResponse(credential ServiceCredential) serviceCredentialResponse {
	var lastRevealedAt *time.Time
	if !credential.LastRevealedAt.IsZero() {
		value := credential.LastRevealedAt
		lastRevealedAt = &value
	}
	return serviceCredentialResponse{
		ID:             credential.ID,
		Type:           credential.Type,
		MaskedHint:     credential.MaskedHint,
		Status:         credential.Status,
		LastRevealedAt: lastRevealedAt,
	}
}

func newServiceCredentialResponses(credentials []ServiceCredential) []serviceCredentialResponse {
	responses := make([]serviceCredentialResponse, 0, len(credentials))
	for _, credential := range credentials {
		responses = append(responses, newServiceCredentialResponse(credential))
	}
	return responses
}

func (request credentialRevealRequest) toInput(
	tenantID tenant.ID,
	serviceID ServiceID,
	credentialID CredentialID,
	actorID identity.UserID,
	clientIP string,
	userAgent string,
) RevealServiceCredentialInput {
	return RevealServiceCredentialInput{
		TenantID:     tenantID,
		ServiceID:    serviceID,
		CredentialID: credentialID,
		ActorID:      actorID,
		ClientIP:     clientIP,
		UserAgent:    userAgent,
		Reason:       request.Reason,
	}
}

func newCredentialRevealResponse(result RevealServiceCredentialResult) credentialRevealResponse {
	return credentialRevealResponse{
		ID:                   result.Credential.ID,
		Type:                 result.Credential.Type,
		MaskedHint:           result.Credential.MaskedHint,
		Status:               result.Credential.Status,
		Payload:              result.Payload,
		RevealedAt:           result.RevealedAt,
		RevealExpiresMessage: result.RevealExpiresMessage,
	}
}
