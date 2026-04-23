package audit

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type logSummaryResponse struct {
	ID            ID            `json:"id"`
	DisplayID     int64         `json:"display_id"`
	TenantID      tenant.ID     `json:"tenant_id"`
	ActorID       ActorID       `json:"actor_id,omitempty"`
	ActorType     ActorType     `json:"actor_type"`
	Action        string        `json:"action"`
	TargetType    string        `json:"target_type"`
	TargetID      TargetID      `json:"target_id"`
	IPAddress     string        `json:"ip_address,omitempty"`
	CorrelationID CorrelationID `json:"correlation_id"`
	CreatedAt     time.Time     `json:"created_at"`
}

type logDetailResponse struct {
	ID                     ID              `json:"id"`
	DisplayID              int64           `json:"display_id"`
	TenantID               tenant.ID       `json:"tenant_id"`
	ActorID                ActorID         `json:"actor_id,omitempty"`
	ActorType              ActorType       `json:"actor_type"`
	Action                 string          `json:"action"`
	TargetType             string          `json:"target_type"`
	TargetID               TargetID        `json:"target_id"`
	BeforeSnapshotRedacted json.RawMessage `json:"before_snapshot_redacted,omitempty"`
	AfterSnapshotRedacted  json.RawMessage `json:"after_snapshot_redacted,omitempty"`
	MetadataRedacted       json.RawMessage `json:"metadata_redacted"`
	IPAddress              string          `json:"ip_address,omitempty"`
	UserAgent              string          `json:"user_agent,omitempty"`
	CorrelationID          CorrelationID   `json:"correlation_id"`
	CreatedAt              time.Time       `json:"created_at"`
}

func newLogSummaryResponse(record Log) logSummaryResponse {
	return logSummaryResponse{
		ID:            record.ID,
		DisplayID:     record.DisplayID,
		TenantID:      record.TenantID,
		ActorID:       record.ActorID,
		ActorType:     record.ActorType,
		Action:        record.Action,
		TargetType:    record.TargetType,
		TargetID:      record.TargetID,
		IPAddress:     record.IPAddress,
		CorrelationID: record.CorrelationID,
		CreatedAt:     record.CreatedAt,
	}
}

func newLogSummaryResponses(records []Log) []logSummaryResponse {
	responses := make([]logSummaryResponse, 0, len(records))
	for _, record := range records {
		responses = append(responses, newLogSummaryResponse(record))
	}
	return responses
}

func newLogDetailResponse(record Log) logDetailResponse {
	return logDetailResponse{
		ID:                     record.ID,
		DisplayID:              record.DisplayID,
		TenantID:               record.TenantID,
		ActorID:                record.ActorID,
		ActorType:              record.ActorType,
		Action:                 record.Action,
		TargetType:             record.TargetType,
		TargetID:               record.TargetID,
		BeforeSnapshotRedacted: record.BeforeSnapshotRedacted,
		AfterSnapshotRedacted:  record.AfterSnapshotRedacted,
		MetadataRedacted:       record.MetadataRedacted,
		IPAddress:              record.IPAddress,
		UserAgent:              record.UserAgent,
		CorrelationID:          record.CorrelationID,
		CreatedAt:              record.CreatedAt,
	}
}
