package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type clientServiceRenewalSmokeResponse struct {
	Service            serviceInstanceSmokeResponse `json:"service"`
	Invoice            renewalInvoiceSmokeResponse  `json:"invoice"`
	PaymentTransaction renewalPaymentSmokeResponse  `json:"payment_transaction"`
	Ledger             renewalLedgerSmokeResponse   `json:"ledger"`
	AmountMinor        int64                        `json:"amount_minor"`
	Currency           string                       `json:"currency"`
	Renewed            bool                         `json:"renewed"`
}

type renewalInvoiceSmokeResponse struct {
	ID         string `json:"id"`
	DisplayID  int64  `json:"display_id"`
	Status     string `json:"status"`
	TotalMinor int64  `json:"total_minor"`
	Currency   string `json:"currency"`
}

type renewalPaymentSmokeResponse struct {
	ID        string `json:"id"`
	DisplayID int64  `json:"display_id"`
	Status    string `json:"status"`
}

type renewalLedgerSmokeResponse struct {
	ID        string `json:"id"`
	DisplayID int64  `json:"display_id"`
	WalletID  string `json:"wallet_id"`
	EntryType string `json:"entry_type"`
}

func runClientServiceRenewalSmoke(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	scenario billingMutationScenario,
	service serviceInstanceSmokeResponse,
) (clientServiceRenewalSmokeResponse, error) {
	headers := cloneHeaders(clientHeaders())
	headers["Idempotency-Key"] = scenario.renewalIdempotencyKey()
	record, err := doJSON[clientServiceRenewalSmokeResponse](ctx, client, http.MethodPost, baseURL, "/client/services/"+url.PathEscape(service.ID)+"/renew", headers, clientServiceRenewalBody{
		WalletID:   smokeWalletID,
		FromStatus: service.Status,
		Reason:     "Smoke renewal " + scenario.RunID,
	}, http.StatusOK)
	if err != nil {
		return clientServiceRenewalSmokeResponse{}, err
	}
	if err := validateClientServiceRenewalSmoke(service, record); err != nil {
		return clientServiceRenewalSmokeResponse{}, err
	}
	fmt.Printf(
		"billing mutation passed: service renewed %s (%d) invoice=%d transaction=%d ledger=%d\n",
		record.Service.ID,
		record.Service.DisplayID,
		record.Invoice.DisplayID,
		record.PaymentTransaction.DisplayID,
		record.Ledger.DisplayID,
	)
	return record, nil
}

func validateClientServiceRenewalSmoke(before serviceInstanceSmokeResponse, record clientServiceRenewalSmokeResponse) error {
	if !record.Renewed {
		return fmt.Errorf("expected service renewal response to be renewed")
	}
	if record.Service.ID != before.ID || record.Service.DisplayID != before.DisplayID {
		return fmt.Errorf("expected renewal service %s (%d), got %+v", before.ID, before.DisplayID, record.Service)
	}
	if record.Service.Status != "active" || record.Service.BillingStatus != "paid" {
		return fmt.Errorf("expected renewed active paid service, got %+v", record.Service)
	}
	if err := validateRenewedTermEnd(before.TermEnd, record.Service.TermEnd); err != nil {
		return err
	}
	if record.AmountMinor != smokeOrderAmount || record.Currency != smokeOrderCurrency {
		return fmt.Errorf("expected renewal amount %d %s, got %d %s", smokeOrderAmount, smokeOrderCurrency, record.AmountMinor, record.Currency)
	}
	if record.Invoice.ID == "" || record.Invoice.DisplayID <= 0 || record.Invoice.Status != "paid" || record.Invoice.TotalMinor != record.AmountMinor || record.Invoice.Currency != record.Currency {
		return fmt.Errorf("expected paid renewal invoice, got %+v", record.Invoice)
	}
	if record.PaymentTransaction.ID == "" || record.PaymentTransaction.DisplayID <= 0 || record.PaymentTransaction.Status != "posted" {
		return fmt.Errorf("expected posted renewal payment transaction, got %+v", record.PaymentTransaction)
	}
	if record.Ledger.ID == "" || record.Ledger.DisplayID <= 0 || record.Ledger.WalletID != smokeWalletID || record.Ledger.EntryType != "purchase" {
		return fmt.Errorf("expected renewal purchase ledger for wallet, got %+v", record.Ledger)
	}
	return nil
}

func validateRenewedTermEnd(beforeValue string, afterValue string) error {
	before, err := parseSmokeTime(beforeValue)
	if err != nil {
		return fmt.Errorf("parse service term_end before renewal: %w", err)
	}
	after, err := parseSmokeTime(afterValue)
	if err != nil {
		return fmt.Errorf("parse service term_end after renewal: %w", err)
	}
	if !after.After(before) {
		return fmt.Errorf("expected renewed term_end after previous term_end, before=%s after=%s", before.Format(time.RFC3339), after.Format(time.RFC3339))
	}
	return nil
}

func parseSmokeTime(value string) (time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed, nil
	}
	return time.Parse(time.RFC3339, value)
}
