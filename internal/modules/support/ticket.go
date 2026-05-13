package support

import (
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type AddTicketNoteInput struct {
	Actor         identity.Actor
	TicketID      TicketID
	TenantID      tenant.ID
	Visibility    NoteVisibility
	BodyRedacted  string
	CorrelationID CorrelationID
}

func (id TicketID) Empty() bool {
	return trim(string(id)) == ""
}

func (id TicketNoteID) Empty() bool {
	return trim(string(id)) == ""
}

func (input CreateTicketInput) Normalize() CreateTicketInput {
	output := input
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.RequesterUserID = identity.UserID(trim(string(output.RequesterUserID)))
	output.Category = TicketCategory(trim(string(output.Category)))
	if output.Category == "" {
		output.Category = TicketCategoryOther
	}
	output.Priority = TicketPriority(trim(string(output.Priority)))
	if output.Priority == "" {
		output.Priority = TicketPriorityP3
	}
	output.Subject = trim(output.Subject)
	output.ReferenceType = ReferenceType(trim(string(output.ReferenceType)))
	output.ReferenceID = ReferenceID(trim(string(output.ReferenceID)))
	output.InitialNote = trim(output.InitialNote)
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	return output
}

func (input CreateTicketInput) Validate() error {
	if err := input.Actor.Validate(); err != nil {
		return err
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.RequesterUserID == "" {
		return identity.ErrActorIDMissing
	}
	if !input.Category.Valid() {
		return ErrTicketCategoryInvalid
	}
	if !input.Priority.Valid() {
		return ErrTicketPriorityInvalid
	}
	if input.Subject == "" {
		return ErrTicketSubjectMissing
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}

func (input CreateTicketNoteInput) Normalize() CreateTicketNoteInput {
	output := input
	output.TicketID = TicketID(trim(string(output.TicketID)))
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.AuthorID = identity.UserID(trim(string(output.AuthorID)))
	output.Visibility = NoteVisibility(trim(string(output.Visibility)))
	if output.Visibility == "" {
		output.Visibility = NoteVisibilityInternal
	}
	output.BodyRedacted = trim(output.BodyRedacted)
	return output
}

func (input CreateTicketNoteInput) Validate() error {
	if input.TicketID.Empty() {
		return ErrTicketIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.AuthorID == "" {
		return identity.ErrActorIDMissing
	}
	if !input.Visibility.Valid() {
		return ErrTicketNoteVisibilityInvalid
	}
	if input.BodyRedacted == "" {
		return ErrTicketNoteBodyMissing
	}
	return nil
}

func (input AddTicketNoteInput) Normalize() AddTicketNoteInput {
	output := input
	output.TicketID = TicketID(trim(string(output.TicketID)))
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.Visibility = NoteVisibility(trim(string(output.Visibility)))
	if output.Visibility == "" {
		output.Visibility = NoteVisibilityInternal
	}
	output.BodyRedacted = trim(output.BodyRedacted)
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	return output
}

func (input AddTicketNoteInput) Validate() error {
	if err := input.Actor.Validate(); err != nil {
		return err
	}
	if input.TicketID.Empty() {
		return ErrTicketIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if !input.Visibility.Valid() {
		return ErrTicketNoteVisibilityInvalid
	}
	if input.BodyRedacted == "" {
		return ErrTicketNoteBodyMissing
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}
