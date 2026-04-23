package audit

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrAuditLogNotFound     = errors.New("audit log not found")
	ErrActorTypeMissing     = errors.New("audit actor type missing")
	ErrActorTypeInvalid     = errors.New("audit actor type invalid")
	ErrActorIDMissing       = errors.New("audit actor id missing")
	ErrActionMissing        = errors.New("audit action missing")
	ErrTargetTypeMissing    = errors.New("audit target type missing")
	ErrTargetIDMissing      = errors.New("audit target id missing")
	ErrCorrelationIDMissing = errors.New("audit correlation id missing")
	ErrCreatedTimeInvalid   = errors.New("audit created time invalid")
	ErrCreatedWindowInvalid = errors.New("audit created time window invalid")
	ErrServiceStoreMissing  = errors.New("audit service store missing")
)

type ID string
type ActorID string
type TargetID string
type CorrelationID string

func (id ID) Empty() bool       { return strings.TrimSpace(string(id)) == "" }
func (id ActorID) Empty() bool  { return strings.TrimSpace(string(id)) == "" }
func (id TargetID) Empty() bool { return strings.TrimSpace(string(id)) == "" }

type ActorType string

const (
	ActorTypeUser            ActorType = "user"
	ActorTypeSystem          ActorType = "system"
	ActorTypeWorker          ActorType = "worker"
	ActorTypeProviderWebhook ActorType = "provider_webhook"
)

func (actorType ActorType) Valid() bool {
	switch actorType {
	case ActorTypeUser, ActorTypeSystem, ActorTypeWorker, ActorTypeProviderWebhook:
		return true
	default:
		return false
	}
}

type Log struct {
	ID                     ID
	DisplayID              int64
	TenantID               tenant.ID
	ActorID                ActorID
	ActorType              ActorType
	Action                 string
	TargetType             string
	TargetID               TargetID
	BeforeSnapshotRedacted json.RawMessage
	AfterSnapshotRedacted  json.RawMessage
	MetadataRedacted       json.RawMessage
	IPAddress              string
	UserAgent              string
	CorrelationID          CorrelationID
	CreatedAt              time.Time
}

type AppendInput struct {
	TenantID               tenant.ID
	ActorID                ActorID
	ActorType              ActorType
	Action                 string
	TargetType             string
	TargetID               TargetID
	BeforeSnapshotRedacted json.RawMessage
	AfterSnapshotRedacted  json.RawMessage
	MetadataRedacted       json.RawMessage
	IPAddress              string
	UserAgent              string
	CorrelationID          CorrelationID
}

type Filter struct {
	TenantID    tenant.ID
	ActorID     ActorID
	ActorType   ActorType
	DisplayID   int64
	Action      string
	TargetType  string
	TargetID    TargetID
	CreatedFrom time.Time
	CreatedTo   time.Time
	Limit       int
}

type Lookup struct {
	ID       ID
	TenantID tenant.ID
}

func (input AppendInput) Normalize() AppendInput {
	output := input
	output.Action = strings.TrimSpace(output.Action)
	output.TargetType = strings.TrimSpace(output.TargetType)
	output.IPAddress = strings.TrimSpace(output.IPAddress)
	output.UserAgent = strings.TrimSpace(output.UserAgent)
	if len(output.MetadataRedacted) == 0 {
		output.MetadataRedacted = json.RawMessage(`{}`)
	}
	return output
}

func (input AppendInput) Validate() error {
	if input.ActorType == "" {
		return ErrActorTypeMissing
	}
	if !input.ActorType.Valid() {
		return ErrActorTypeInvalid
	}
	if input.ActorType == ActorTypeUser && input.ActorID == "" {
		return ErrActorIDMissing
	}
	if input.Action == "" {
		return ErrActionMissing
	}
	if input.TargetType == "" {
		return ErrTargetTypeMissing
	}
	if input.TargetID == "" {
		return ErrTargetIDMissing
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}

type Store interface {
	Append(ctx context.Context, input AppendInput) (Log, error)
	ListLogs(ctx context.Context, filter Filter) ([]Log, error)
	GetLog(ctx context.Context, lookup Lookup) (Log, error)
}
