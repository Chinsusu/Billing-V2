package wallet

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func (handler *HTTPHandler) adminTopupRequestsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminTopupRequests, handler.options.AdminMiddleware),
	})
}

func (handler *HTTPHandler) resellerTopupRequestsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(handler.handleListAdminTopupRequests, handler.options.ResellerMiddleware),
	})
}

func (handler *HTTPHandler) adminTopupRequestRoute(w http.ResponseWriter, r *http.Request) {
	topupRequestID, action, ok := topupRequestPath(w, r, adminTopupRequestPrefix)
	if !ok {
		return
	}
	switch action {
	case "":
		dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
			http.MethodGet: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
				handler.handleGetAdminTopupRequest(w, r, topupRequestID)
			}, handler.options.AdminMiddleware),
		})
	case "approve":
		dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
			http.MethodPost: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
				handler.handleApproveAdminTopupRequest(w, r, topupRequestID)
			}, handler.options.AdminReviewMiddleware),
		})
	case "reject":
		dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
			http.MethodPost: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
				handler.handleRejectAdminTopupRequest(w, r, topupRequestID)
			}, handler.options.AdminReviewMiddleware),
		})
	default:
		writeWalletError(w, r, ErrTopupRequestIDMissing)
	}
}

func (handler *HTTPHandler) clientTopupRequestsRoute(w http.ResponseWriter, r *http.Request) {
	dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet:  handler.tenantRoute(handler.handleListClientTopupRequests, handler.options.ClientMiddleware),
		http.MethodPost: handler.tenantRoute(handler.handleCreateClientTopupRequest, handler.options.ClientMiddleware),
	})
}

func (handler *HTTPHandler) clientTopupRequestRoute(w http.ResponseWriter, r *http.Request) {
	topupRequestID, action, ok := topupRequestPath(w, r, clientTopupRequestPrefix)
	if !ok {
		return
	}
	if action != "" {
		writeWalletError(w, r, ErrTopupRequestIDMissing)
		return
	}
	dispatchWalletMethods(w, r, map[string]http.HandlerFunc{
		http.MethodGet: handler.tenantRoute(func(w http.ResponseWriter, r *http.Request) {
			handler.handleGetClientTopupRequest(w, r, topupRequestID)
		}, handler.options.ClientMiddleware),
	})
}

func (handler *HTTPHandler) handleCreateClientTopupRequest(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := clientTenantActor(w, r)
	if !ok {
		return
	}
	var body createTopupRequestBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField("body", "request.body_invalid", "Request body must be valid JSON."),
		})
		return
	}
	if _, err := handler.service.GetWallet(r.Context(), WalletLookup{
		ID: body.WalletID, TenantID: tenantID, OwnerType: OwnerTypeUser, OwnerID: UserOwnerID(actor.ID),
	}); err != nil {
		writeWalletError(w, r, err)
		return
	}
	request, err := handler.service.CreateTopupRequest(r.Context(), CreateTopupRequestInput{
		TenantID:         tenantID,
		WalletID:         body.WalletID,
		RequestedBy:      actor.ID,
		AmountMinor:      body.AmountMinor,
		Currency:         body.Currency,
		PaymentMethod:    body.PaymentMethod,
		PaymentReference: body.PaymentReference,
		IdempotencyKey:   IdempotencyKey(strings.TrimSpace(r.Header.Get("Idempotency-Key"))),
	})
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusCreated, newTopupRequestResponse(request))
}

func (handler *HTTPHandler) handleListClientTopupRequests(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := clientTenantActor(w, r)
	if !ok {
		return
	}
	filter, page, ok := topupRequestFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	filter.RequestedBy = actor.ID
	requests, err := handler.service.ListTopupRequests(r.Context(), filter)
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newTopupRequestResponses(requests), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetClientTopupRequest(w http.ResponseWriter, r *http.Request, topupRequestID TopupRequestID) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, actor, ok := clientTenantActor(w, r)
	if !ok {
		return
	}
	request, err := handler.service.GetTopupRequest(r.Context(), TopupRequestLookup{
		ID: topupRequestID, TenantID: tenantID, RequestedBy: actor.ID,
	})
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newTopupRequestResponse(request))
}

func (handler *HTTPHandler) handleListAdminTopupRequests(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromContext(w, r)
	if !ok {
		return
	}
	if _, ok := actorFromContext(w, r); !ok {
		return
	}
	filter, page, ok := topupRequestFilterFromRequest(w, r)
	if !ok {
		return
	}
	filter.TenantID = tenantID
	requests, err := handler.service.ListTopupRequests(r.Context(), filter)
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newTopupRequestResponses(requests), httpserver.NewPage(page.Limit, ""))
}

func (handler *HTTPHandler) handleGetAdminTopupRequest(w http.ResponseWriter, r *http.Request, topupRequestID TopupRequestID) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromContext(w, r)
	if !ok {
		return
	}
	if _, ok := actorFromContext(w, r); !ok {
		return
	}
	request, err := handler.service.GetTopupRequest(r.Context(), TopupRequestLookup{ID: topupRequestID, TenantID: tenantID})
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newTopupRequestResponse(request))
}

func (handler *HTTPHandler) handleApproveAdminTopupRequest(w http.ResponseWriter, r *http.Request, topupRequestID TopupRequestID) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromContext(w, r)
	if !ok {
		return
	}
	actor, ok := actorFromContext(w, r)
	if !ok {
		return
	}
	var body reviewTopupRequestBody
	if !decodeTopupJSON(w, r, &body) {
		return
	}
	request, err := handler.service.ApproveTopupRequest(r.Context(), ApproveTopupRequestInput{
		ID: topupRequestID, TenantID: tenantID, ReviewedBy: actor.ID, ReviewNote: body.ReviewNote,
	})
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newTopupRequestResponse(request))
}

func (handler *HTTPHandler) handleRejectAdminTopupRequest(w http.ResponseWriter, r *http.Request, topupRequestID TopupRequestID) {
	if !handler.ready(w, r) {
		return
	}
	tenantID, ok := tenantIDFromContext(w, r)
	if !ok {
		return
	}
	actor, ok := actorFromContext(w, r)
	if !ok {
		return
	}
	var body reviewTopupRequestBody
	if !decodeTopupJSON(w, r, &body) {
		return
	}
	request, err := handler.service.RejectTopupRequest(r.Context(), RejectTopupRequestInput{
		ID: topupRequestID, TenantID: tenantID, ReviewedBy: actor.ID, ReviewNote: body.ReviewNote,
	})
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	httpserver.WriteSuccess(w, r, http.StatusOK, newTopupRequestResponse(request))
}

func topupRequestFilterFromRequest(w http.ResponseWriter, r *http.Request) (TopupRequestFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return TopupRequestFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := TopupRequestFilter{Limit: page.Limit}
	query := r.URL.Query()
	if walletID := WalletID(strings.TrimSpace(query.Get("wallet_id"))); walletID != "" {
		filter.WalletID = walletID
	}
	if requestedBy := strings.TrimSpace(query.Get("requested_by")); requestedBy != "" {
		filter.RequestedBy = identity.UserID(requestedBy)
	}
	if displayID, present, ok := walletPositiveInt64Query(w, r, "display_id"); !ok {
		return TopupRequestFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	if method := PaymentMethod(strings.TrimSpace(query.Get("payment_method"))); method != "" {
		if !method.Valid() {
			writeWalletError(w, r, ErrPaymentMethodInvalid)
			return TopupRequestFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.PaymentMethod = method
	}
	if status := TopupStatus(strings.TrimSpace(query.Get("status"))); status != "" {
		if !status.Valid() {
			writeWalletError(w, r, ErrTopupStatusInvalid)
			return TopupRequestFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	amountMin, amountMax, ok := walletAmountRangeQuery(w, r)
	if !ok {
		return TopupRequestFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter.AmountMinMinor = amountMin
	filter.AmountMaxMinor = amountMax
	return filter, page, true
}

func topupRequestPath(w http.ResponseWriter, r *http.Request, prefix string) (TopupRequestID, string, bool) {
	value := strings.Trim(strings.TrimPrefix(r.URL.Path, prefix), "/")
	parts := strings.Split(value, "/")
	if len(parts) == 0 || parts[0] == "" || len(parts) > 2 {
		writeWalletError(w, r, ErrTopupRequestIDMissing)
		return "", "", false
	}
	if len(parts) == 2 {
		return TopupRequestID(parts[0]), parts[1], true
	}
	return TopupRequestID(parts[0]), "", true
}

func decodeTopupJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField("body", "request.body_invalid", "Request body must be valid JSON."),
		})
		return false
	}
	return true
}
