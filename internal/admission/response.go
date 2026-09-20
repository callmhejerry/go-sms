package admission

import "github.com/callmhejerry/sms/internal/shared/store"

type AcceptAdmissionResponse struct {
	Admission store.Admission `json:"admission"`
	Student   store.Student   `json:"student"`
}
