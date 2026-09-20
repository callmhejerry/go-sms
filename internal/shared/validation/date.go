package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func (v *AppValidator) date(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	_, err := time.Parse("2006-01-28", value)
	if err != nil {
		return false
	}

	return true
}
