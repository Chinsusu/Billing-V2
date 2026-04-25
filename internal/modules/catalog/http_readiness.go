package catalog

import (
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func (handler *HTTPHandler) handleListProviderSourceReadiness(w http.ResponseWriter, r *http.Request) {
	if !handler.ready(w, r) {
		return
	}
	filter, page, ok := providerSourceReadinessFilterFromRequest(w, r)
	if !ok {
		return
	}
	records, err := handler.service.ListProviderSourceReadiness(r.Context(), filter)
	if err != nil {
		writeCatalogError(w, r, err)
		return
	}
	httpserver.WriteList(w, r, http.StatusOK, newProviderSourceReadinessResponses(records), httpserver.NewPage(page.Limit, ""))
}

func providerSourceReadinessFilterFromRequest(w http.ResponseWriter, r *http.Request) (ProviderSourceReadinessFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return ProviderSourceReadinessFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := ProviderSourceReadinessFilter{Limit: page.Limit}
	query := r.URL.Query()
	if planDisplayID, present, ok := catalogPositiveInt64Query(w, r, "plan_display_id"); !ok {
		return ProviderSourceReadinessFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.PlanDisplayID = planDisplayID
	}
	productType := ProductType(strings.TrimSpace(query.Get("product_type")))
	if productType != "" {
		if !productType.Valid() {
			writeFieldValidation(w, r, "product_type", "catalog.product_type_invalid", "Product type is invalid.")
			return ProviderSourceReadinessFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.ProductType = productType
	}
	status := PlanStatus(strings.TrimSpace(query.Get("status")))
	if status != "" {
		if !status.Valid() {
			writeFieldValidation(w, r, "status", "catalog.status_invalid", "Status is invalid.")
			return ProviderSourceReadinessFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.PlanStatus = status
	}
	if sourceDisplayID, present, ok := catalogPositiveInt64Query(w, r, "source_display_id"); !ok {
		return ProviderSourceReadinessFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.SourceDisplayID = sourceDisplayID
	}
	return filter, page, true
}

type providerSourceReadinessResponse struct {
	PlanDisplayID       int64                        `json:"plan_display_id"`
	PlanCode            string                       `json:"plan_code"`
	PlanName            string                       `json:"plan_name"`
	ProductType         ProductType                  `json:"product_type"`
	PlanStatus          PlanStatus                   `json:"plan_status"`
	PlanSourceDisplayID *int64                       `json:"plan_source_display_id,omitempty"`
	PlanSourceStatus    PlanSourceStatus             `json:"plan_source_status,omitempty"`
	SourceDisplayID     *int64                       `json:"source_display_id,omitempty"`
	SourceName          string                       `json:"source_name,omitempty"`
	SourceType          provider.Type                `json:"source_type,omitempty"`
	SourceStatus        ProviderSourceStatus         `json:"source_status,omitempty"`
	InventoryMode       InventoryMode                `json:"inventory_mode,omitempty"`
	State               ProviderSourceReadinessState `json:"state"`
	Reason              string                       `json:"reason"`
}

func newProviderSourceReadinessResponse(record ProviderSourceReadiness) providerSourceReadinessResponse {
	response := providerSourceReadinessResponse{
		PlanDisplayID: record.PlanDisplayID,
		PlanCode:      record.PlanCode,
		PlanName:      record.PlanName,
		ProductType:   record.ProductType,
		PlanStatus:    record.PlanStatus,
		State:         record.State,
		Reason:        record.Reason,
	}
	if record.PlanSourceDisplayID > 0 {
		response.PlanSourceDisplayID = int64Ptr(record.PlanSourceDisplayID)
		response.PlanSourceStatus = record.PlanSourceStatus
	}
	if record.SourceDisplayID > 0 {
		response.SourceDisplayID = int64Ptr(record.SourceDisplayID)
		response.SourceName = record.SourceName
		response.SourceType = record.SourceType
		response.SourceStatus = record.SourceStatus
		response.InventoryMode = record.InventoryMode
	}
	return response
}

func newProviderSourceReadinessResponses(records []ProviderSourceReadiness) []providerSourceReadinessResponse {
	responses := make([]providerSourceReadinessResponse, 0, len(records))
	for _, record := range records {
		responses = append(responses, newProviderSourceReadinessResponse(record))
	}
	return responses
}

func int64Ptr(value int64) *int64 {
	return &value
}
