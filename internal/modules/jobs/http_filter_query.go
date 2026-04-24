package jobs

import (
	"net/http"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func jobPositiveInt64Query(w http.ResponseWriter, r *http.Request, field string) (int64, bool, bool) {
	value, present, err := httpserver.ParseOptionalPositiveInt64Query(r, field)
	if err != nil {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField(field, "request.display_id_invalid", "Display id must be a positive number."),
		})
		return 0, present, false
	}
	return value, present, true
}
