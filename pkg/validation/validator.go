package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// ValidateDateFormat checks if the date string matches the supported friendly layouts
// Usage: binding:"date_format=dd/mm/yyyy"
func ValidateDateFormat(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	formatParam := fl.Param()

	var layout string
	switch formatParam {
	case "dd/mm/yyyy":
		layout = "02/01/2006"
	case "yyyy-mm-dd":
		layout = "2006-01-02"
	case "mm/dd/yyyy":
		layout = "01/02/2006"
	default:
		// Fallback to usage as direct Go layout if not found in map
		layout = formatParam
	}

	_, err := time.Parse(layout, dateStr)
	return err == nil
}


func ValidatePhone(fl validator.FieldLevel) bool {
	// Custom regex for generic phone if needed (though e164 is safer)
	var phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(fl.Field().String())
}

// RegisterTagName registers the "json" tag name as the field name in validator errors.
// This ensures error messages use the JSON field names (e.g. "email" instead of "Email").
func RegisterTagName(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// FormatError converts validator.ValidationErrors into a human-readable string.
func FormatError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range errs {
			messages = append(messages, formatField(e))
		}
		return strings.Join(messages, "; ")
	}
	return err.Error()
}

func formatField(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", e.Field())
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("Field '%s' must be at most %s characters", e.Field(), e.Param())
	case "date_format":
		return fmt.Sprintf("Field '%s' must match format %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("Field '%s' failed validation on '%s'", e.Field(), e.Tag())
	}
}
