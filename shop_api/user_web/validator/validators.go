package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^1([38][0-9]|4[5-9]|5[0-3,5-9]|6[2567]|7[0-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}
