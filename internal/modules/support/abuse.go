package support

import (
	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func (id AbuseCaseID) Empty() bool {
	return trim(string(id)) == ""
}

func (input CreateAbuseCaseInput) Normalize() CreateAbuseCaseInput {
	output := input
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.UserID = identity.UserID(trim(string(output.UserID)))
	output.ServiceID = order.ServiceID(trim(string(output.ServiceID)))
	output.ProviderSourceID = catalog.ProviderSourceID(trim(string(output.ProviderSourceID)))
	output.CaseType = AbuseCaseType(trim(string(output.CaseType)))
	output.Severity = AbuseSeverity(trim(string(output.Severity)))
	output.ReportSource = AbuseReportSource(trim(string(output.ReportSource)))
	output.EvidenceSummaryRedacted = trim(output.EvidenceSummaryRedacted)
	output.AssignedOwnerID = identity.UserID(trim(string(output.AssignedOwnerID)))
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	return output
}

func (input CreateAbuseCaseInput) Validate() error {
	if err := input.Actor.Validate(); err != nil {
		return err
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if !input.CaseType.Valid() {
		return ErrAbuseCaseTypeInvalid
	}
	if !input.Severity.Valid() {
		return ErrAbuseSeverityInvalid
	}
	if !input.ReportSource.Valid() {
		return ErrAbuseReportSourceInvalid
	}
	if input.EvidenceSummaryRedacted == "" {
		return ErrEvidenceSummaryMissing
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}

func (input SuspendServiceForAbuseInput) Normalize() SuspendServiceForAbuseInput {
	output := input
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.CaseID = AbuseCaseID(trim(string(output.CaseID)))
	output.ServiceID = order.ServiceID(trim(string(output.ServiceID)))
	output.FromStatus = order.ServiceStatus(trim(string(output.FromStatus)))
	output.Reason = trim(output.Reason)
	return output
}

func (input SuspendServiceForAbuseInput) Validate() error {
	if err := input.Actor.Validate(); err != nil {
		return err
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.CaseID.Empty() {
		return ErrAbuseCaseIDMissing
	}
	if input.ServiceID.Empty() {
		return order.ErrServiceIDMissing
	}
	if !input.FromStatus.Valid() {
		return order.ErrServiceStatusInvalid
	}
	if input.Reason == "" {
		return ErrAbuseReasonMissing
	}
	return nil
}

func (input MarkAbuseCaseSuspendedInput) Normalize() MarkAbuseCaseSuspendedInput {
	output := input
	output.ID = AbuseCaseID(trim(string(output.ID)))
	output.TenantID = tenant.ID(trim(string(output.TenantID)))
	output.ActionTaken = trim(output.ActionTaken)
	return output
}

func (input MarkAbuseCaseSuspendedInput) Validate() error {
	if input.ID.Empty() {
		return ErrAbuseCaseIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.ActionTaken == "" {
		return ErrAbuseReasonMissing
	}
	return nil
}

func supportAuditActorType(actor identity.Actor) audit.ActorType {
	if actor.IsSystem {
		return audit.ActorTypeSystem
	}
	return audit.ActorTypeUser
}
