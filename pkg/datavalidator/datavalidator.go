package datavalidator

import (
	"errors"
	"reflect"
	"regexp"
)

// ValidatorFunc defines a function type for validation
type ValidatorFunc func(interface{}) error

// ValidationRule stores field validation logic
type ValidationRule struct {
	Field     string
	Validator ValidatorFunc
	ErrorMsg  string
}

// ValidateStruct applies validation rules on a struct
func ValidateStruct(data interface{}, rules []ValidationRule) map[string]string {
	errorsMap := make(map[string]string)
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, rule := range rules {
		field := val.FieldByName(rule.Field)
		if !field.IsValid() {
			errorsMap[rule.Field] = "Invalid field"
			continue
		}

		err := rule.Validator(field.Interface())
		if err != nil {
			errorsMap[rule.Field] = rule.ErrorMsg
		}
	}

	return errorsMap
}

// Predefined validation functions
func NotEmpty(value interface{}) error {
	if str, ok := value.(string); ok && str == "" {
		return errors.New("Field cannot be empty")
	}
	return nil
}

func IsEmail(value interface{}) error {
	return MatchRegex(value, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "Invalid email format")
}

func MinLength(length int) ValidatorFunc {
	return func(value interface{}) error {
		if str, ok := value.(string); ok && len(str) < length {
			return errors.New("Field must be at least " + string(length) + " characters long")
		}
		return nil
	}
}

func MinValue(min int) ValidatorFunc {
	return func(value interface{}) error {
		if num, ok := value.(int); ok && num < min {
			return errors.New("Value must be at least " + string(min))
		}
		return nil
	}
}

func MaxValue(max int) ValidatorFunc {
	return func(value interface{}) error {
		if num, ok := value.(int); ok && num > max {
			return errors.New("Value must not exceed " + string(max))
		}
		return nil
	}
}

func ValidateNestedStruct(field interface{}) error {
	subErrors := ValidateStruct(field, nil) // Nested validation
	if len(subErrors) > 0 {
		return errors.New("Nested struct validation failed")
	}
	return nil
}

// MatchRegex validates a string against a given regex pattern
func MatchRegex(value interface{}, pattern string, errorMsg string) error {
	if str, ok := value.(string); ok {
		re := regexp.MustCompile(pattern)
		if !re.MatchString(str) {
			return errors.New(errorMsg)
		}
	}
	return nil
}

// ID & Document Validators
func IsPAN(value interface{}) error {
	return MatchRegex(value, `^[A-Z]{5}[0-9]{4}[A-Z]{1}$`, "Invalid PAN format")
}

func IsAadhaar(value interface{}) error {
	return MatchRegex(value, `^[2-9]{1}[0-9]{11}$`, "Invalid Aadhaar number")
}

func IsPassport(value interface{}) error {
	return MatchRegex(value, `^[A-Z]{1}[0-9]{7}$`, "Invalid Passport format")
}

func IsCreditCard(value interface{}) error {
	return MatchRegex(value, `^\d{16}$`, "Invalid Credit Card number")
}

// IsPhoneNumber validates phone numbers (supports international formats)
func IsPhoneNumber(value interface{}) error {
	return MatchRegex(value, `^\+?[1-9]\d{1,14}$`, "Invalid phone number format")
}

// IsAlpha validates if a string contains only letters
func IsAlpha(value interface{}) error {
	return MatchRegex(value, `^[a-zA-Z]+$`, "Field must contain only letters")
}

// IsAlphanumeric validates if a string contains only letters and numbers
func IsAlphanumeric(value interface{}) error {
	return MatchRegex(value, `^[a-zA-Z0-9]+$`, "Field must contain only letters and numbers")
}
