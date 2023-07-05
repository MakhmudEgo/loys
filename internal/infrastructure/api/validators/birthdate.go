package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

const birthdayValidateName = "birthdate"

func birthdayValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.DateOnly, fl.Field().String())
	return err == nil
}
