package utils

import (
	"net/http"
	"strconv"
)

type CursorPaginationResponse struct {
	HasMore bool   `json:"has_more"`
	Next    string `json:"next"`
}

type OffsetPaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type OffsetPaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type CursorPaginationRequest struct {
	Next     *string `json:"next"`
	PageSize int     `json:"page_size"`
}

func ParseCursorPagination(r *http.Request) CursorPaginationRequest {
	query := r.URL.Query()

	var next *string
	pageSize := 20

	if value := query.Get("next"); value != "" {
		next = &value
	}

	if value := query.Get("page_size"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			if parsed > 100 {
				pageSize = 100
			} else if pageSize < 1 {
				pageSize = 20
			}
		}
	}
	return CursorPaginationRequest{
		Next:     next,
		PageSize: pageSize,
	}
}

func ParseOffsetPagination(r *http.Request) OffsetPaginationRequest {
	query := r.URL.Query()

	page := 1
	pageSize := 20

	if value := query.Get("page"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			page = parsed
		}
	}

	if value := query.Get("page_size"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			if parsed > 100 {
				pageSize = 100
			} else if pageSize < 1 {
				pageSize = 20
			}

		}
	}

	return OffsetPaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}
}
