package validation

import "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func NewValdiator() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())

	validator := &Validator{
		validator: v,
	}

	v.RegisterValidation("date_of_birth", validator.dateOfBirthValidator)

	return validator
}
