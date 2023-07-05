package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

const genderValidateName = "gender"

func genderValidator(fl validator.FieldLevel) bool {
	for _, value := range []string{"male", "female"} {
		if strings.Contains(fl.Field().String(), value) {
			return true
		}
	}
	return false
}
