package audit

import (
	"errors"
	"testing"
)

func TestAppendInputNormalizeDefaultsMetadata(t *testing.T) {
	input := AppendInput{
		ActorType:     ActorTypeSystem,
		Action:        " tenant.create ",
		TargetType:    " tenant ",
		TargetID:      "11111111-1111-1111-1111-111111111111",
		CorrelationID: "22222222-2222-2222-2222-222222222222",
	}

	normalized := input.Normalize()

	if normalized.Action != "tenant.create" {
		t.Fatalf("expected trimmed action, got %q", normalized.Action)
	}
	if normalized.TargetType != "tenant" {
		t.Fatalf("expected trimmed target type, got %q", normalized.TargetType)
	}
	if string(normalized.MetadataRedacted) != "{}" {
		t.Fatalf("expected default metadata, got %s", normalized.MetadataRedacted)
	}
}

func TestAppendInputValidateRequiresActorIDForUser(t *testing.T) {
	input := AppendInput{
		ActorType:     ActorTypeUser,
		Action:        "tenant.update",
		TargetType:    "tenant",
		TargetID:      "11111111-1111-1111-1111-111111111111",
		CorrelationID: "22222222-2222-2222-2222-222222222222",
	}

	if err := input.Validate(); !errors.Is(err, ErrActorIDMissing) {
		t.Fatalf("expected actor id error, got %v", err)
	}
}

func TestAppendInputValidateRejectsMissingCorrelationID(t *testing.T) {
	input := AppendInput{
		ActorType:  ActorTypeSystem,
		Action:     "tenant.update",
		TargetType: "tenant",
		TargetID:   "11111111-1111-1111-1111-111111111111",
	}

	if err := input.Validate(); !errors.Is(err, ErrCorrelationIDMissing) {
		t.Fatalf("expected correlation id error, got %v", err)
	}
}
