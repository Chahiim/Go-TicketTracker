package validator

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// validates email fomrat
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) ValidData() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(field string, message string) {
	_, exists := v.Errors[field]
	if !exists {
		v.Errors[field] = message
	}
}

func (v *Validator) Check(ok bool, field string, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func IsValidEmail(email string) bool {
	return EmailRX.MatchString(email)
}

// Checks if a time.Time value is not the zero value
func IsValidDate(value time.Time) bool {
	return !value.IsZero()
}

func HasNumber(value string) bool {
	return regexp.MustCompile(`[0-9]`).MatchString(value)
}

func HasUpper(value string) bool {
	return regexp.MustCompile(`[A-Z]`).MatchString(value)
}

func HasSymbol(value string) bool {
	return regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(value)
}
