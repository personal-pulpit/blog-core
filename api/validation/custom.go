package validation

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type CustomError struct {
	Field   string
	Value   string
	Message string
}

var (
	usernameValidaitor validator.Func = func(fld validator.FieldLevel) bool {
		username := fld.Field().String()
		re := regexp.MustCompile("^[a-zA-Z0-9_]{5,20}$")
		return re.MatchString(username)
	}
	phoneNumberValidatior validator.Func = func(fld validator.FieldLevel) bool {
		phonenumber := fld.Field().String()
		re := regexp.MustCompile(`^(\+98|0)?9\d{9}$`)
		return re.MatchString(phonenumber)
	}
	emailValidatior validator.Func = func(fld validator.FieldLevel) bool {
		email := fld.Field().String()
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		return re.MatchString(email)
	}
)

func InitValidations() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		val.RegisterValidation("usernamevalidaitor", usernameValidaitor)
		val.RegisterValidation("emailvalidatior", emailValidatior)
		val.RegisterValidation("phonenumbervalidaitor", phoneNumberValidatior)
	}
}

