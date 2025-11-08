package validations

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// StructValidation const
const (
	StructValidationTimeAfterNow                      = "time_after_now"
	StructValidationTimeAfterField                    = "time_after_field"
	StructValidationMinimumIfFieldEqual               = "min_if_field_eq"
	StructValidationMaximumIfFieldEqual               = "max_if_field_eq"
	StructValidationLessThanEqualFieldIfFieldEqual    = "lte_field_if_field_eq"
	StructValidationGreaterThanEqualFieldIfFieldEqual = "gte_field_if_field_eq"
	StructValidationMinimumFieldIfFieldEqual          = "min_field_if_field_eq"
	StructValidationMaximumFieldIfFieldEqual          = "max_field_if_field_eq"
)

// InitStructValidation init struct validation
func InitStructValidation() {
	structValidation := map[string]func(fl validator.FieldLevel) bool{
		StructValidationTimeAfterNow:                      TimeAfterNow,
		StructValidationTimeAfterField:                    TimeAfterField,
		StructValidationMinimumIfFieldEqual:               MinIfFieldEqual,
		StructValidationMaximumIfFieldEqual:               MaxIfFieldEqual,
		StructValidationLessThanEqualFieldIfFieldEqual:    LTEFieldIfFieldEqual,
		StructValidationGreaterThanEqualFieldIfFieldEqual: GTEFieldIfFieldEqual,
		StructValidationMinimumFieldIfFieldEqual:          MinFieldIfFieldEqual,
		StructValidationMaximumFieldIfFieldEqual:          MaxFieldIfFieldEqual,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for tag, validationFunc := range structValidation {
			err := v.RegisterValidation(tag, validationFunc)
			if err != nil {
				panic(fmt.Errorf("can not register validation function: %s", tag))
			}
		}
	}
}

func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

func IsSameDateMonth(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}
