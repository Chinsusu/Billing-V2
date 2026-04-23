package invoice

type Status string

const (
	StatusDraft         Status = "draft"
	StatusIssued        Status = "issued"
	StatusPaid          Status = "paid"
	StatusPartiallyPaid Status = "partially_paid"
	StatusOverdue       Status = "overdue"
	StatusVoided        Status = "voided"
)

func (status Status) Valid() bool {
	switch status {
	case StatusDraft, StatusIssued, StatusPaid, StatusPartiallyPaid, StatusOverdue, StatusVoided:
		return true
	default:
		return false
	}
}

func CanTransition(from Status, to Status) bool {
	if from == to {
		return from.Valid()
	}
	switch from {
	case StatusDraft:
		return to == StatusIssued || to == StatusVoided
	case StatusIssued:
		return to == StatusPartiallyPaid || to == StatusPaid || to == StatusOverdue || to == StatusVoided
	case StatusPartiallyPaid:
		return to == StatusPaid || to == StatusOverdue || to == StatusVoided
	case StatusOverdue:
		return to == StatusPartiallyPaid || to == StatusPaid || to == StatusVoided
	default:
		return false
	}
}
