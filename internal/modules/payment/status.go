package payment

type TransactionType string

const (
	TransactionTypeCharge     TransactionType = "charge"
	TransactionTypeRefund     TransactionType = "refund"
	TransactionTypeAdjustment TransactionType = "adjustment"
)

func (transactionType TransactionType) Valid() bool {
	switch transactionType {
	case TransactionTypeCharge, TransactionTypeRefund, TransactionTypeAdjustment:
		return true
	default:
		return false
	}
}

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusPosted  TransactionStatus = "posted"
	TransactionStatusFailed  TransactionStatus = "failed"
	TransactionStatusVoided  TransactionStatus = "voided"
)

func (status TransactionStatus) Valid() bool {
	switch status {
	case TransactionStatusPending, TransactionStatusPosted, TransactionStatusFailed, TransactionStatusVoided:
		return true
	default:
		return false
	}
}
