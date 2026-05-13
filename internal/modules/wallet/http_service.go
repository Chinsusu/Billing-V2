package wallet

import "context"

type HTTPService interface {
	ListWallets(ctx context.Context, filter WalletFilter) ([]Wallet, error)
	GetWallet(ctx context.Context, lookup WalletLookup) (Wallet, error)
	ListLedgerEntries(ctx context.Context, filter LedgerEntryFilter) ([]LedgerEntry, error)
	CreateWalletRefund(ctx context.Context, input CreateWalletRefundInput) (LedgerEntry, error)
	CreateWalletAdjustment(ctx context.Context, input CreateWalletAdjustmentInput) (LedgerEntry, error)
	CreateTopupRequest(ctx context.Context, input CreateTopupRequestInput) (TopupRequest, error)
	ListTopupRequests(ctx context.Context, filter TopupRequestFilter) ([]TopupRequest, error)
	GetTopupRequest(ctx context.Context, lookup TopupRequestLookup) (TopupRequest, error)
	ApproveTopupRequest(ctx context.Context, input ApproveTopupRequestInput) (TopupRequest, error)
	RejectTopupRequest(ctx context.Context, input RejectTopupRequestInput) (TopupRequest, error)
}
