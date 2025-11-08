package validations

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-playground/validator/v10"
)

// IsFieldEqual validate a field that must satisfy value(this) == param
func IsFieldEqual(field reflect.Value, param string) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == param
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := strconv.ParseInt(param, 0, 64)

		return int64(field.Len()) == p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := strconv.ParseInt(param, 0, 64)

		return field.Int() == p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := strconv.ParseUint(param, 0, 64)

		return field.Uint() == p
	case reflect.Float32, reflect.Float64:
		p, _ := strconv.ParseFloat(param, 64)

		return field.Float() == p
	case reflect.Bool:
		p, _ := strconv.ParseBool(param)

		return field.Bool() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// SplitBySpaceWithQuote split string by space, ignore if the string is quoted
//
// Ex: aaaa bbb "cc    ccc" 'a c' => [aaaa, bbb, "cc    ccc", 'a c']
func SplitBySpaceWithQuote(value string) []string {
	r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'`)

	return r.FindAllString(value, -1)
}

// UnQuote unquote a string (support single and double quote)
//
// Ex: "this is double quoted string" => this is double quoted string
func UnQuote(value string) string {
	if value[0] == '\'' || value[0] == '"' {
		value = value[1:]
	}

	if value[len(value)-1] == '\'' || value[len(value)-1] == '"' {
		value = value[:len(value)-1]
	}

	return value
}

// TimeAfterNow validate a strfmt.DateTime field that must satisfy value(this) > current time
//
// Usage: `binding:"time_after_now"`
func TimeAfterNow(fl validator.FieldLevel) bool {
	field := fl.Field()

	if dateField, ok := field.Interface().(strfmt.DateTime); ok {
		return time.Time(dateField).After(time.Now())
	}

	return false
}

// TimeAfterField validate a strfmt.DateTime field that must satisfy value(this) > value(AnotherTimeField)
//
// Usage: `binding:"time_after_field=AnotherTimeField"`
func TimeAfterField(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := SplitBySpaceWithQuote(fl.Param())

	if len(param) != 1 {
		return false
	}

	var paramField reflect.Value
	if fl.Parent().Kind() == reflect.Ptr {
		paramField = fl.Parent().Elem().FieldByName(param[0])
	} else {
		paramField = fl.Parent().FieldByName(param[0])
	}

	if dateField, ok := field.Interface().(strfmt.DateTime); ok {
		if dateParamField, ok2 := paramField.Interface().(strfmt.DateTime); ok2 {
			isValid := time.Time(dateField).After(time.Time(dateParamField))

			return isValid
		}
	}

	return false
}

// MinIfFieldEqual validate a number field that must satisfy value(this) >= ANumber if ConditionalField == "ConditionValue"
//
// Usage: `binding:"min_if_field_eq=ANumber ConditionalField 'ConditionValue'"`
func MinIfFieldEqual(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := SplitBySpaceWithQuote(fl.Param())

	if len(param) != 3 {
		return false
	}

	var param1Field reflect.Value
	if fl.Parent().Kind() == reflect.Ptr {
		param1Field = fl.Parent().Elem().FieldByName(param[1])
	} else {
		param1Field = fl.Parent().FieldByName(param[1])
	}

	if !IsFieldEqual(param1Field, UnQuote(param[2])) {
		return true
	}

	iField := field.Interface()
	f64Param0, err := strconv.ParseFloat(param[0], 64)
	if err != nil {
		return false
	}

	return iField.(float64) >= f64Param0
}

// MaxIfFieldEqual validate a number field that must satisfy value(this) <= ANumber if ConditionalField == "ConditionValue"
//
// Usage: `binding:"max_if_field_eq=ANumber ConditionalField 'ConditionValue'"`
func MaxIfFieldEqual(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := SplitBySpaceWithQuote(fl.Param())

	if len(param) != 3 {
		return false
	}

	var param1Field reflect.Value
	if fl.Parent().Kind() == reflect.Ptr {
		param1Field = fl.Parent().Elem().FieldByName(param[1])
	} else {
		param1Field = fl.Parent().FieldByName(param[1])
	}

	if !IsFieldEqual(param1Field, UnQuote(param[2])) {
		return true
	}

	iField := field.Interface()
	f64Param0, err := strconv.ParseFloat(param[0], 64)
	if err != nil {
		return false
	}

	return iField.(float64) <= f64Param0
}

// LTEFieldIfFieldEqual validate a number field that must satisfy value(this) <= value(DestField) if ConditionalField == "ConditionValue"
//
// Usage: `binding:"lte_field_if_field_eq=DestField ConditionalField 'ConditionValue'"`
func LTEFieldIfFieldEqual(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := SplitBySpaceWithQuote(fl.Param())

	if len(param) != 3 {
		return false
	}

	var param0Field reflect.Value
	var param1Field reflect.Value
	if fl.Parent().Kind() == reflect.Ptr {
		param0Field = fl.Parent().Elem().FieldByName(param[0])
		param1Field = fl.Parent().Elem().FieldByName(param[1])
	} else {
		param0Field = fl.Parent().FieldByName(param[0])
		param1Field = fl.Parent().FieldByName(param[1])
	}

	if !IsFieldEqual(param1Field, UnQuote(param[2])) {
		return true
	}

	iField := field.Interface()
	iParam0Field := param0Field.Interface()

	return iField.(float64) <= iParam0Field.(float64)
}

// GTEFieldIfFieldEqual validate a number field that must satisfy value(this) >= value(DestField) if ConditionalField == "ConditionValue"
//
// Usage: `binding:"gte_field_if_field_eq=DestField ConditionalField 'ConditionValue'"`
func GTEFieldIfFieldEqual(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := SplitBySpaceWithQuote(fl.Param())

	if len(param) != 3 {
		return false
	}

	var param0Field reflect.Value
	var param1Field reflect.Value
	if fl.Parent().Kind() == reflect.Ptr {
		param0Field = fl.Parent().Elem().FieldByName(param[0])
		param1Field = fl.Parent().Elem().FieldByName(param[1])
	} else {
		param0Field = fl.Parent().FieldByName(param[0])
		param1Field = fl.Parent().FieldByName(param[1])
	}

	if !IsFieldEqual(param1Field, UnQuote(param[2])) {
		return true
	}

	iField := field.Interface()
	iParam0Field := param0Field.Interface()

	return iField.(float64) >= iParam0Field.(float64)
}

// MinFieldIfFieldEqual validate a number field that must satisfy min(value(this), value(DestField)) == value(DestField), if ConditionalField == "ConditionValue"
//
// Usage: `binding:"min_field_if_field_eq=DestField ConditionalField 'ConditionValue'"`
func MinFieldIfFieldEqual(fl validator.FieldLevel) bool {
	return GTEFieldIfFieldEqual(fl)
}

// MaxFieldIfFieldEqual validate a number field that must satisfy max(value(this), value(DestField)) == value(DestField), if ConditionalField == "ConditionValue"
//
// Usage: `binding:"max_field_if_field_eq=DestField ConditionalField 'ConditionValue'"`
func MaxFieldIfFieldEqual(fl validator.FieldLevel) bool {
	return LTEFieldIfFieldEqual(fl)
}
