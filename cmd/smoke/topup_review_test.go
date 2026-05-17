package main

import "testing"

func TestVerifyTopupReviewFinalStateAcceptsOnlyWalletCreditDelta(t *testing.T) {
	before := topupReviewBaseline{
		WalletBalanceMinor: 3200,
		OrderCount:         3,
		ProviderJobCount:   2,
		ServiceCount:       1,
	}
	after := topupReviewBaseline{
		WalletBalanceMinor: before.WalletBalanceMinor + topupReviewApproveAmount,
		OrderCount:         before.OrderCount,
		ProviderJobCount:   before.ProviderJobCount,
		ServiceCount:       before.ServiceCount,
	}

	if err := verifyTopupReviewFinalState(before, after); err != nil {
		t.Fatalf("expected final state to pass: %v", err)
	}
}

func TestVerifyTopupReviewFinalStateRejectsProviderSideEffects(t *testing.T) {
	before := topupReviewBaseline{WalletBalanceMinor: 3200, OrderCount: 3, ProviderJobCount: 2, ServiceCount: 1}
	after := topupReviewBaseline{
		WalletBalanceMinor: before.WalletBalanceMinor + topupReviewApproveAmount,
		OrderCount:         before.OrderCount,
		ProviderJobCount:   before.ProviderJobCount + 1,
		ServiceCount:       before.ServiceCount,
	}

	if err := verifyTopupReviewFinalState(before, after); err == nil {
		t.Fatal("expected provider side effect to fail")
	}
}

func TestVerifyTopupReviewFinalStateRejectsUnexpectedWalletDelta(t *testing.T) {
	before := topupReviewBaseline{WalletBalanceMinor: 3200, OrderCount: 3, ProviderJobCount: 2, ServiceCount: 1}
	after := topupReviewBaseline{
		WalletBalanceMinor: before.WalletBalanceMinor + topupReviewApproveAmount + topupReviewRejectAmount,
		OrderCount:         before.OrderCount,
		ProviderJobCount:   before.ProviderJobCount,
		ServiceCount:       before.ServiceCount,
	}

	if err := verifyTopupReviewFinalState(before, after); err == nil {
		t.Fatal("expected unexpected wallet delta to fail")
	}
}
