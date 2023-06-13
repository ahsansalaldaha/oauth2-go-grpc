package utils

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// ValidateContainsNotHaveField - check if password does not contains the username field
func ValidateContainsNotHaveField(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Get the struct value
	structValue := reflect.ValueOf(fl.Parent().Interface())
	name := structValue.FieldByName("Name").String()

	return !strings.Contains(password, name)
}

// ValidateMinLength - check if min length is as per the requirement
func ValidateMinLength(minLength int) validator.Func {
	return func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return len(password) >= minLength
	}
}

// PasswordComplexity - represents the struct that defines the password complexity
type PasswordComplexity struct {
	ShouldHaveUppercase bool
	ShouldHaveLowercase bool
	ShouldHaveNumber    bool
	ShouldHaveSpecial   bool
}

// ValidateComplexity - check if password complexity matches the need
func ValidateComplexity(pc *PasswordComplexity) validator.Func {

	return func(fl validator.FieldLevel) bool {
		// Check the combination of uppercase letters, lowercase letters, numbers, and special characters
		var (
			hasUppercase bool
			hasLowercase bool
			hasNumber    bool
			hasSpecial   bool
		)

		password := fl.Field().String()
		for _, char := range password {
			switch {
			case unicode.IsUpper(char) && pc.ShouldHaveUppercase:
				hasUppercase = true
			case unicode.IsLower(char) && pc.ShouldHaveLowercase:
				hasLowercase = true
			case unicode.IsNumber(char) && pc.ShouldHaveNumber:
				hasNumber = true
			case pc.ShouldHaveSpecial && (unicode.IsPunct(char) || unicode.IsSymbol(char)):
				hasSpecial = true
			}
		}

		if (pc.ShouldHaveUppercase && !hasUppercase) ||
			(pc.ShouldHaveLowercase && !hasLowercase) ||
			(pc.ShouldHaveNumber && !hasNumber) ||
			(pc.ShouldHaveSpecial && !hasSpecial) {
			return false
		}

		return true
	}
}
