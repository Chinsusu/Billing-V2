package notification

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

func TestTelegramDeliveryHandlerSendsSafeSummaryAndMarksSent(t *testing.T) {
	now := time.Date(2026, 5, 19, 16, 0, 0, 0, time.UTC)
	var requestPath string
	var requestBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath = r.URL.Path
		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		requestBody = string(body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	store := &fakeNotificationDeliveryStore{}
	handler, err := NewTelegramDeliveryHandler(store, TelegramConfig{
		BotToken:   "test-token",
		ChatID:     "test-chat",
		APIBaseURL: server.URL,
	}, server.Client())
	if err != nil {
		t.Fatalf("expected telegram handler: %v", err)
	}
	handler.Now = func() time.Time { return now }

	completion, err := handler.Handle(context.Background(), jobs.Job{
		ID:            "job-1",
		ReferenceType: DeliveryReferenceType,
		ReferenceID:   "44444444-4444-4444-4444-444444444444",
		CorrelationID: "corr-safe",
		PayloadJSON: telegramJobPayload(t, DeliveryJobPayload{
			NotificationID:        "44444444-4444-4444-4444-444444444444",
			NotificationDisplayID: 10001,
			EventType:             EventProvisioningManualReview,
			TemplateKey:           EventProvisioningManualReview,
			Channel:               ChannelTelegram,
			PayloadRedacted:       json.RawMessage(`{"provider_token":"secret","safe":"value"}`),
		}),
	})
	if err != nil {
		t.Fatalf("expected telegram delivery success: %v", err)
	}
	if completion.Status != jobs.StatusSucceeded || store.sentID != "44444444-4444-4444-4444-444444444444" ||
		!store.sentAt.Equal(now) {
		t.Fatalf("unexpected completion/store: %+v store=%+v", completion, store)
	}
	if !strings.Contains(requestPath, "/bottest-token/sendMessage") {
		t.Fatalf("expected telegram endpoint path, got %q", requestPath)
	}
	for _, want := range []string{"Billing notification", "Notification #10001", EventProvisioningManualReview, "corr-safe"} {
		if !strings.Contains(requestBody, want) {
			t.Fatalf("expected request body to contain %q, got %s", want, requestBody)
		}
	}
	for _, forbidden := range []string{"secret", "provider_token", "PayloadRedacted"} {
		if strings.Contains(requestBody, forbidden) {
			t.Fatalf("telegram request body leaked %q: %s", forbidden, requestBody)
		}
	}
}

func TestTelegramDeliveryHandlerMarksHTTPFailureRetryable(t *testing.T) {
	now := time.Date(2026, 5, 19, 16, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"ok":false,"description":"do not copy"}`))
	}))
	defer server.Close()

	store := &fakeNotificationDeliveryStore{}
	handler, err := NewTelegramDeliveryHandler(store, TelegramConfig{
		BotToken:   "test-token",
		ChatID:     "test-chat",
		APIBaseURL: server.URL,
	}, server.Client())
	if err != nil {
		t.Fatalf("expected telegram handler: %v", err)
	}
	handler.Now = func() time.Time { return now }

	completion, err := handler.Handle(context.Background(), jobs.Job{
		ID:            "job-1",
		ReferenceType: DeliveryReferenceType,
		ReferenceID:   "44444444-4444-4444-4444-444444444444",
		PayloadJSON: telegramJobPayload(t, DeliveryJobPayload{
			NotificationID:        "44444444-4444-4444-4444-444444444444",
			NotificationDisplayID: 10001,
			EventType:             EventProvisioningFailed,
			TemplateKey:           EventProvisioningFailed,
			Channel:               ChannelTelegram,
		}),
	})
	if err != nil {
		t.Fatalf("expected retryable completion, got error: %v", err)
	}
	if completion.Status != jobs.StatusFailedRetryable || completion.RetrySafety != jobs.RetrySafetySafeRetry {
		t.Fatalf("unexpected retryable completion: %+v", completion)
	}
	if store.failedID != "44444444-4444-4444-4444-444444444444" ||
		store.failedErrorCode != "telegram_http_500" ||
		strings.Contains(store.failedErrorMessage, "do not copy") {
		t.Fatalf("unexpected failed store state: %+v", store)
	}
}

func TestTelegramDeliveryHandlerRejectsNonTelegramJobTerminally(t *testing.T) {
	handler, err := NewTelegramDeliveryHandler(&fakeNotificationDeliveryStore{}, TelegramConfig{
		BotToken: "test-token",
		ChatID:   "test-chat",
	}, nil)
	if err != nil {
		t.Fatalf("expected telegram handler: %v", err)
	}
	completion, err := handler.Handle(context.Background(), jobs.Job{
		ID:            "job-1",
		ReferenceType: DeliveryReferenceType,
		ReferenceID:   "44444444-4444-4444-4444-444444444444",
		PayloadJSON: telegramJobPayload(t, DeliveryJobPayload{
			NotificationID:        "44444444-4444-4444-4444-444444444444",
			NotificationDisplayID: 10001,
			EventType:             EventTopupApproved,
			TemplateKey:           EventTopupApproved,
			Channel:               ChannelDashboard,
		}),
	})
	if err != nil {
		t.Fatalf("expected terminal completion, got error: %v", err)
	}
	if completion.Status != jobs.StatusFailedTerminal || completion.LastErrorCode != "notification_job_invalid" {
		t.Fatalf("unexpected non-telegram completion: %+v", completion)
	}
}

func TestTelegramSenderRequiresConfig(t *testing.T) {
	if _, err := NewTelegramSender(TelegramConfig{ChatID: "chat"}, nil); err == nil ||
		!strings.Contains(err.Error(), "TELEGRAM_BOT_TOKEN is required") {
		t.Fatalf("expected missing token error, got %v", err)
	}
	if _, err := NewTelegramSender(TelegramConfig{BotToken: "token"}, nil); err == nil ||
		!strings.Contains(err.Error(), "TELEGRAM_CHAT_ID is required") {
		t.Fatalf("expected missing chat error, got %v", err)
	}
}

func telegramJobPayload(t *testing.T, payload DeliveryJobPayload) json.RawMessage {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal telegram job payload: %v", err)
	}
	return body
}
