package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func (v *Validator) dateOfBirthValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	dob, err := time.Parse("2006-01-02", value)
	if err != nil {
		return false
	}
	today := time.Now().UTC()

	// DOB cannot be in the future.
	if dob.After(today) {
		return false
	}

	return true
}

func IsValidDOB(dob time.Time) bool {
	if IsSameDay(time.Now(), dob) {
		return false
	}
	if time.Now().Before(dob) {
		return false
	}
	return true
}

func IsSameDay(this, that time.Time) bool {
	return (this.Day() == that.Day() && this.Month() == that.Month() && this.Year() == that.Year())
}
