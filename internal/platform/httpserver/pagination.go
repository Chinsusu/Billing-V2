package httpserver

import (
	"errors"
	"net/http"
	"strconv"
)

const (
	DefaultPageLimit = 20
	MaxPageLimit     = 100
)

var (
	ErrPageLimitInvalid  = errors.New("page limit invalid")
	ErrPageLimitTooLarge = errors.New("page limit too large")
)

type Page struct {
	Limit      int     `json:"limit"`
	NextCursor *string `json:"next_cursor"`
}

type CursorPageRequest struct {
	Limit  int
	Cursor string
}

func ParseCursorPage(r *http.Request) (CursorPageRequest, error) {
	limit := DefaultPageLimit
	limitValue := r.URL.Query().Get("limit")
	if limitValue != "" {
		parsed, err := strconv.Atoi(limitValue)
		if err != nil || parsed <= 0 {
			return CursorPageRequest{}, ErrPageLimitInvalid
		}
		if parsed > MaxPageLimit {
			return CursorPageRequest{}, ErrPageLimitTooLarge
		}
		limit = parsed
	}
	return CursorPageRequest{
		Limit:  limit,
		Cursor: r.URL.Query().Get("cursor"),
	}, nil
}

func NewPage(limit int, nextCursor string) Page {
	page := Page{Limit: limit}
	if nextCursor != "" {
		page.NextCursor = &nextCursor
	}
	return page
}
