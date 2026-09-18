package academic

import "time"

type CreateAcademicSessionRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsCurrent bool      `json:"is_current"`
}
