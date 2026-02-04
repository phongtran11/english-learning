package validation

import (
	"regexp"
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

// Custom regex for generic phone if needed (though e164 is safer)
var phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

func ValidatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}
