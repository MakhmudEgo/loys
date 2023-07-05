package validators

import (
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation(genderValidateName, genderValidator)
	v.RegisterValidation(birthdayValidateName, birthdayValidator)
}
