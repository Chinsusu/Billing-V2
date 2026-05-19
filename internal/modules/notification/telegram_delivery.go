package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
)

const defaultTelegramAPIBaseURL = "https://api.telegram.org"

var (
	ErrTelegramConfigInvalid = errors.New("telegram delivery config invalid")
	ErrTelegramSendFailed    = errors.New("telegram send failed")
)

type TelegramConfig struct {
	BotToken   string
	ChatID     string
	APIBaseURL string
}

type TelegramHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TelegramDeliveryHandler struct {
	Store  DeliveryStore
	Sender *TelegramSender
	Now    func() time.Time
}

type TelegramSender struct {
	Config TelegramConfig
	Client TelegramHTTPClient
}

type telegramSendError struct {
	code      string
	message   string
	retryable bool
}

func (err telegramSendError) Error() string {
	if err.message == "" {
		return ErrTelegramSendFailed.Error()
	}
	return err.message
}

func NewTelegramDeliveryRunner(store jobs.Store, deliveryStore DeliveryStore, workerID jobs.WorkerID, cfg TelegramConfig, client TelegramHTTPClient) (jobs.Runner, error) {
	handler, err := NewTelegramDeliveryHandler(deliveryStore, cfg, client)
	if err != nil {
		return jobs.Runner{}, err
	}
	return jobs.Runner{
		Store:     store,
		Handler:   handler,
		WorkerID:  workerID,
		BatchSize: 10,
		Types:     []jobs.Type{TelegramDeliveryJobType},
	}, nil
}

func NewTelegramDeliveryHandler(store DeliveryStore, cfg TelegramConfig, client TelegramHTTPClient) (*TelegramDeliveryHandler, error) {
	sender, err := NewTelegramSender(cfg, client)
	if err != nil {
		return nil, err
	}
	return &TelegramDeliveryHandler{Store: store, Sender: sender}, nil
}

func NewTelegramSender(cfg TelegramConfig, client TelegramHTTPClient) (*TelegramSender, error) {
	cfg = cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &TelegramSender{Config: cfg, Client: client}, nil
}

func (cfg TelegramConfig) Normalize() TelegramConfig {
	cfg.BotToken = strings.TrimSpace(cfg.BotToken)
	cfg.ChatID = strings.TrimSpace(cfg.ChatID)
	cfg.APIBaseURL = strings.TrimRight(strings.TrimSpace(cfg.APIBaseURL), "/")
	if cfg.APIBaseURL == "" {
		cfg.APIBaseURL = defaultTelegramAPIBaseURL
	}
	return cfg
}

func (cfg TelegramConfig) Validate() error {
	if cfg.BotToken == "" {
		return fmt.Errorf("%w: TELEGRAM_BOT_TOKEN is required", ErrTelegramConfigInvalid)
	}
	if cfg.ChatID == "" {
		return fmt.Errorf("%w: TELEGRAM_CHAT_ID is required", ErrTelegramConfigInvalid)
	}
	if _, err := url.ParseRequestURI(cfg.APIBaseURL); err != nil {
		return fmt.Errorf("%w: TELEGRAM_API_BASE_URL is invalid", ErrTelegramConfigInvalid)
	}
	return nil
}

func (handler *TelegramDeliveryHandler) Handle(ctx context.Context, job jobs.Job) (jobs.Completion, error) {
	if handler == nil || handler.Store == nil || handler.Sender == nil {
		return jobs.Completion{}, ErrStoreMissing
	}
	if job.ReferenceType != DeliveryReferenceType || job.ReferenceID == "" {
		return invalidTelegramJobCompletion(handler.now()), nil
	}
	var payload DeliveryJobPayload
	if err := json.Unmarshal(job.PayloadJSON, &payload); err != nil {
		return invalidTelegramJobCompletion(handler.now()), nil
	}
	if payload.NotificationID == "" || payload.Channel != ChannelTelegram {
		return invalidTelegramJobCompletion(handler.now()), nil
	}

	message := TelegramMessageFromJob(job, payload)
	if err := handler.Sender.Send(ctx, message); err != nil {
		sendErr := telegramErrorDetails(err)
		if _, markErr := handler.Store.MarkNotificationFailed(ctx, payload.NotificationID, handler.now(), sendErr.code, sendErr.message); markErr != nil {
			return jobs.Completion{}, markErr
		}
		status := jobs.StatusFailedTerminal
		retrySafety := jobs.RetrySafetyDoNotRetry
		if sendErr.retryable {
			status = jobs.StatusFailedRetryable
			retrySafety = jobs.RetrySafetySafeRetry
		}
		return jobs.Completion{
			Status:                   status,
			RetrySafety:              retrySafety,
			LastErrorCode:            sendErr.code,
			LastErrorMessageRedacted: sendErr.message,
			FinishedAt:               handler.now(),
		}, nil
	}

	if _, err := handler.Store.MarkNotificationSent(ctx, payload.NotificationID, handler.now()); err != nil {
		if errors.Is(err, ErrNotificationNotFound) {
			return jobs.Completion{
				Status:                   jobs.StatusFailedTerminal,
				RetrySafety:              jobs.RetrySafetyDoNotRetry,
				LastErrorCode:            "notification_not_found",
				LastErrorMessageRedacted: "notification record was not found",
				FinishedAt:               handler.now(),
			}, nil
		}
		return jobs.Completion{}, err
	}
	return jobs.Completion{Status: jobs.StatusSucceeded, FinishedAt: handler.now()}, nil
}

func TelegramMessageFromJob(job jobs.Job, payload DeliveryJobPayload) string {
	lines := []string{
		"Billing notification",
		fmt.Sprintf("Notification #%d", payload.NotificationDisplayID),
		"Event: " + safeTelegramField(payload.EventType),
		"Template: " + safeTelegramField(payload.TemplateKey),
		"Channel: telegram",
	}
	if job.CorrelationID != "" {
		lines = append(lines, "Correlation: "+safeTelegramField(string(job.CorrelationID)))
	}
	return strings.Join(lines, "\n")
}

func (sender *TelegramSender) Send(ctx context.Context, message string) error {
	if sender == nil || sender.Client == nil {
		return telegramSendError{code: "telegram_sender_missing", message: "telegram sender is not configured"}
	}
	if strings.TrimSpace(message) == "" {
		return telegramSendError{code: "telegram_message_empty", message: "telegram message is empty"}
	}
	body, err := json.Marshal(map[string]any{
		"chat_id":                  sender.Config.ChatID,
		"text":                     message,
		"disable_web_page_preview": true,
	})
	if err != nil {
		return telegramSendError{code: "telegram_request_invalid", message: "telegram request is invalid"}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sender.endpointURL(), bytes.NewReader(body))
	if err != nil {
		return telegramSendError{code: "telegram_request_invalid", message: "telegram request is invalid"}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := sender.Client.Do(req)
	if err != nil {
		return telegramSendError{code: "telegram_request_failed", message: "telegram request failed", retryable: true}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return telegramSendError{
			code:      fmt.Sprintf("telegram_http_%d", resp.StatusCode),
			message:   "telegram returned non-success status",
			retryable: resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500,
		}
	}
	return nil
}

func (sender *TelegramSender) endpointURL() string {
	return sender.Config.APIBaseURL + "/bot" + sender.Config.BotToken + "/sendMessage"
}

func invalidTelegramJobCompletion(now time.Time) jobs.Completion {
	return jobs.Completion{
		Status:                   jobs.StatusFailedTerminal,
		RetrySafety:              jobs.RetrySafetyDoNotRetry,
		LastErrorCode:            "notification_job_invalid",
		LastErrorMessageRedacted: "telegram notification delivery job is invalid",
		FinishedAt:               now,
	}
}

func telegramErrorDetails(err error) telegramSendError {
	var sendErr telegramSendError
	if errors.As(err, &sendErr) {
		return sendErr
	}
	return telegramSendError{code: "telegram_send_failed", message: "telegram send failed", retryable: true}
}

func safeTelegramField(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}

func (handler *TelegramDeliveryHandler) now() time.Time {
	if handler.Now == nil {
		return time.Now().UTC()
	}
	return handler.Now()
}
