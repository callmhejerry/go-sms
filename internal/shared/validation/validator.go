package validation

import (
	"fmt"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"github.com/go-playground/validator/v10"
)

type AppValidator struct {
	validator *validator.Validate
}

func NewValdiator() *AppValidator {
	v := validator.New(validator.WithRequiredStructEnabled())

	validator := &AppValidator{
		validator: v,
	}

	v.RegisterValidation("date_of_birth", validator.dateOfBirthValidator)
	v.RegisterValidation("date", validator.date)

	return validator
}

// ValidateStruct validates a struct and returns a clean apierror.
func (appValidator *AppValidator) ValidateStruct(s any) error {
	validate := appValidator.validator

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return apierror.Validation("invalid request")
	}

	// Take the first error for simplicity (you can later return all of them)
	fieldErr := validationErrors[0]

	field := strings.ToLower(fieldErr.Field())
	tag := fieldErr.Tag()

	var message string
	switch tag {
	case "required":
		message = fmt.Sprintf("%s is required", field)
	case "email":
		message = fmt.Sprintf("%s must be a valid email", field)
	case "min":
		message = fmt.Sprintf("%s is too short", field)
	case "max":
		message = fmt.Sprintf("%s is too long", field)
	case "oneof":
		message = fmt.Sprintf("%s has an invalid value", field)
	default:
		message = fmt.Sprintf("%s is invalid", field)
	}

	return apierror.Validation(message)
}
