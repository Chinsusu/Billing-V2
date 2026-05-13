package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
)

const (
	credentialAuditActionRevealed = "credential.revealed"
	credentialAuditTargetType     = "service_credential"
	defaultRevealExpiresMessage   = "Credential is shown once. Store it securely before leaving this screen."
)

func (service *Service) ListServiceCredentials(ctx context.Context, filter ServiceCredentialFilter) ([]ServiceCredential, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	if service.credentials == nil {
		return nil, ErrCredentialStoreMissing
	}
	filter = filter.Normalize()
	if err := filter.Validate(); err != nil {
		return nil, err
	}
	return service.credentials.ListServiceCredentials(ctx, filter)
}

func (service *Service) RevealServiceCredential(ctx context.Context, input RevealServiceCredentialInput) (RevealServiceCredentialResult, error) {
	if err := service.readyCredentialReveal(); err != nil {
		return RevealServiceCredentialResult{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return RevealServiceCredentialResult{}, err
	}
	serviceRecord, err := service.store.GetServiceInstance(ctx, ServiceInstanceLookup{
		ID:          input.ServiceID,
		TenantID:    input.TenantID,
		BuyerUserID: input.BuyerUserID,
	})
	if err != nil {
		return RevealServiceCredentialResult{}, err
	}
	credential, err := service.credentials.GetServiceCredential(ctx, ServiceCredentialLookup{
		ID:        input.CredentialID,
		TenantID:  input.TenantID,
		ServiceID: serviceRecord.ID,
	})
	if err != nil {
		return RevealServiceCredentialResult{}, err
	}
	if credential.Status != CredentialStatusActive {
		return RevealServiceCredentialResult{}, ErrCredentialStatusInvalid
	}
	if err := service.enforceCredentialRevealRateLimit(ctx, input, serviceRecord); err != nil {
		return RevealServiceCredentialResult{}, err
	}
	plaintext, err := service.credentialCipher.Decrypt(credential.EncryptedPayload)
	if err != nil {
		return RevealServiceCredentialResult{}, fmt.Errorf("%w: %v", ErrCredentialDecryptFailed, err)
	}
	now := service.now().UTC()
	credential, err = service.credentials.MarkServiceCredentialRevealed(ctx, MarkServiceCredentialRevealedInput{
		ID:         credential.ID,
		TenantID:   credential.TenantID,
		ServiceID:  credential.ServiceID,
		ActorID:    input.ActorID,
		RevealedAt: now,
	})
	if err != nil {
		return RevealServiceCredentialResult{}, err
	}
	if err := service.appendCredentialRevealAudit(ctx, input, serviceRecord, credential); err != nil {
		return RevealServiceCredentialResult{}, err
	}
	return RevealServiceCredentialResult{
		Credential:           credential,
		Payload:              revealPayloadJSON(plaintext),
		RevealedAt:           now,
		RevealExpiresMessage: defaultRevealExpiresMessage,
	}, nil
}

func (service *Service) readyCredentialReveal() error {
	if err := service.ready(); err != nil {
		return err
	}
	if service.credentials == nil {
		return ErrCredentialStoreMissing
	}
	if service.credentialCipher == nil {
		return ErrCredentialCipherMissing
	}
	if service.credentialRevealLimits == nil {
		return ErrCredentialRevealLimiterMissing
	}
	return nil
}

func (service *Service) enforceCredentialRevealRateLimit(ctx context.Context, input RevealServiceCredentialInput, serviceRecord ServiceInstance) error {
	windowStart := service.now().UTC().Truncate(credentialRevealRateLimitWindow)
	counter, err := service.credentialRevealLimits.IncrementCredentialRevealRateLimit(ctx, CredentialRevealRateLimitInput{
		TenantID:    input.TenantID,
		ActorID:     input.ActorID,
		ServiceID:   serviceRecord.ID,
		WindowStart: windowStart,
	})
	if err != nil {
		return err
	}
	if counter.AttemptCount > credentialRevealRateLimitMaxAttempts {
		return ErrCredentialRevealRateLimited
	}
	return nil
}

func (service *Service) appendCredentialRevealAudit(
	ctx context.Context,
	input RevealServiceCredentialInput,
	serviceRecord ServiceInstance,
	credential ServiceCredential,
) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:   credential.TenantID,
		ActorID:    audit.ActorID(input.ActorID),
		ActorType:  audit.ActorTypeUser,
		Action:     credentialAuditActionRevealed,
		TargetType: credentialAuditTargetType,
		TargetID:   audit.TargetID(credential.ID),
		MetadataRedacted: orderAuditJSON(credentialRevealAuditMetadata{
			ServiceID:        serviceRecord.ID,
			ServiceDisplayID: serviceRecord.DisplayID,
			CredentialType:   credential.Type,
			Reason:           input.Reason,
		}),
		IPAddress:     input.ClientIP,
		UserAgent:     input.UserAgent,
		CorrelationID: audit.CorrelationID(credential.ID),
	})
	return err
}

type credentialRevealAuditMetadata struct {
	ServiceID        ServiceID      `json:"service_id"`
	ServiceDisplayID int64          `json:"service_display_id,omitempty"`
	CredentialType   CredentialType `json:"credential_type"`
	Reason           string         `json:"reason,omitempty"`
}

func revealPayloadJSON(plaintext string) json.RawMessage {
	trimmed := []byte(plaintext)
	if json.Valid(trimmed) {
		return append(json.RawMessage(nil), trimmed...)
	}
	body, err := json.Marshal(map[string]string{"value": plaintext})
	if err != nil {
		return json.RawMessage(`{"value":""}`)
	}
	return body
}
