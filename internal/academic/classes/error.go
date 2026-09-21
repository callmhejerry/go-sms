package classes

import "github.com/callmhejerry/sms/internal/shared/apierror"

var (
	ErrClassNotFound = apierror.NotFound("class not found")

	ErrClassArmNotFound = apierror.NotFound("class arm not found")
)
