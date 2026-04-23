package httpserver

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var ErrQueryIntegerInvalid = errors.New("query integer invalid")

func ParseOptionalPositiveInt64Query(r *http.Request, field string) (int64, bool, error) {
	return parseOptionalInt64Query(r, field, func(value int64) bool {
		return value > 0
	})
}

func ParseOptionalNonNegativeInt64Query(r *http.Request, field string) (int64, bool, error) {
	return parseOptionalInt64Query(r, field, func(value int64) bool {
		return value >= 0
	})
}

func parseOptionalInt64Query(r *http.Request, field string, valid func(int64) bool) (int64, bool, error) {
	value := strings.TrimSpace(r.URL.Query().Get(field))
	if value == "" {
		return 0, false, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || !valid(parsed) {
		return 0, true, ErrQueryIntegerInvalid
	}
	return parsed, true, nil
}
