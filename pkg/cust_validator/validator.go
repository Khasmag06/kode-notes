package cust_validator

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var symbols = "!@#$%^&*"

func PasswordValidate(fl validator.FieldLevel) bool {

	password := fl.Field().String()

	var hasLower, hasUpper, hasDigit bool

	for _, char := range password {
		switch {
		case unicode.IsLetter(char) && !unicode.Is(unicode.Latin, char):
			return false
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasDigit = true
		default:
			if !strings.Contains(symbols, string(char)) {
				return false
			}
		}
	}
	return hasUpper && hasLower && hasDigit
}

func GetErrorMsg(fe validator.FieldError) error {
	switch fe.Tag() {
	case "required":
		return fmt.Errorf("поле %s обязательно для заполнения", fe.Field())
	case "email":
		return fmt.Errorf("поле %s должно содержать корректный адрес электронной почты", fe.Field())
	case "password":
		return fmt.Errorf("поле %s обязательно должно содержать латинские буквы разного регистра, цифры и допускаются символы: %s", fe.Field(), symbols)
	case "min":
		return fmt.Errorf("поле %s должно быть не менее %s символов", fe.Field(), fe.Param())
	case "max":
		return fmt.Errorf("поле %s должно быть не более %s символов", fe.Field(), fe.Param())
	default:
		return fmt.Errorf("поле %s содержит некорректное значение", fe.Field())
	}
}
