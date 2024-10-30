package event

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"unicode"
)

func ValidateUsername() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		usernameToValidate := fl.Field().String()
		usernameExpression := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]+$`)
		return usernameExpression.MatchString(usernameToValidate)
	}
}

func ValidatePassword() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		passwordToValidate := fl.Field().String()

		var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool
		if len(passwordToValidate) >= 8 {
			hasMinLen = true
		}

		for _, char := range passwordToValidate {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}

		return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
	}
}

func ValidateAccountType() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		accountTypeToValidate := fl.Field().String()
		return accountTypeToValidate == string(TeacherAccountType) || accountTypeToValidate == string(StudentAccountType)
	}
}
