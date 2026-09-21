package student

import "github.com/callmhejerry/sms/internal/shared/apierror"

var (
	ErrStudentNotFound = apierror.NotFound("Student with this id not found")
)
