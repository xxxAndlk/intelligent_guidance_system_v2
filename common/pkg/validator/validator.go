package validator

import (
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var (
	phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)
	uuidRegex  = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

type Validator struct {
	errors []string
}

func New() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, field+" is required")
	}
	return v
}

func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errors = append(v.errors, field+" must be at least "+string(rune(min))+" characters")
	}
	return v
}

func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.errors = append(v.errors, field+" must be at most "+string(rune(max))+" characters")
	}
	return v
}

func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		v.errors = append(v.errors, field+" must be a valid email address")
	}
	return v
}

func (v *Validator) Phone(field, value string) *Validator {
	if value == "" {
		return v
	}
	if !phoneRegex.MatchString(value) {
		v.errors = append(v.errors, field+" must be a valid phone number")
	}
	return v
}

func (v *Validator) URL(field, value string) *Validator {
	if value == "" {
		return v
	}
	_, err := url.ParseRequestURI(value)
	if err != nil {
		v.errors = append(v.errors, field+" must be a valid URL")
	}
	return v
}

func (v *Validator) UUID(field, value string) *Validator {
	if value == "" {
		return v
	}
	if !uuidRegex.MatchString(strings.ToLower(value)) {
		v.errors = append(v.errors, field+" must be a valid UUID")
	}
	return v
}

func (v *Validator) In(field, value string, allowed ...string) *Validator {
	if value == "" {
		return v
	}
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.errors = append(v.errors, field+" must be one of: "+strings.Join(allowed, ", "))
	return v
}

func (v *Validator) Match(field, value string, pattern *regexp.Regexp) *Validator {
	if value == "" {
		return v
	}
	if !pattern.MatchString(value) {
		v.errors = append(v.errors, field+" format is invalid")
	}
	return v
}

func (v *Validator) HasError() bool {
	return len(v.errors) > 0
}

func (v *Validator) Errors() []string {
	return v.errors
}

func (v *Validator) Error() string {
	return strings.Join(v.errors, "; ")
}

func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinLength(value string, min int) bool {
	return len(value) >= min
}

func MaxLength(value string, max int) bool {
	return len(value) <= max
}

func IsEmail(value string) bool {
	if value == "" {
		return false
	}
	_, err := mail.ParseAddress(value)
	return err == nil
}

func IsPhone(value string) bool {
	if value == "" {
		return false
	}
	return phoneRegex.MatchString(value)
}

func IsURL(value string) bool {
	if value == "" {
		return false
	}
	_, err := url.ParseRequestURI(value)
	return err == nil
}

func IsUUID(value string) bool {
	if value == "" {
		return false
	}
	return uuidRegex.MatchString(strings.ToLower(value))
}

func IsAlphanumeric(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func IsNumeric(value string) bool {
	for _, r := range value {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(value) > 0
}

func IsAlpha(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(value) > 0
}

func InSlice[T comparable](value T, allowed []T) bool {
	for _, a := range allowed {
		if value == a {
			return true
		}
	}
	return false
}

func TrimStrings(s string) string {
	return strings.TrimSpace(s)
}

func ToLowerStrings(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func ToUpperStrings(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}