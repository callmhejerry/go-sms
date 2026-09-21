package student

import (
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/callmhejerry/sms/internal/shared/utils"
)

type ListStudentOffsetResponse struct {
	Data       []store.Student                `json:"data"`
	Pagination utils.OffsetPaginationResponse `json:"pagination"`
}

type ListStudentCursorResponse struct {
	Data       []store.Student                `json:"data"`
	Pagination utils.CursorPaginationResponse `json:"pagination"`
}
