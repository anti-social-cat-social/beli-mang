package validation

import (
	"errors"
	localError "halosuster/pkg/error"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type validationError struct {
	Field    string `json:"field"`
	Messsage string `json:"message"`
}

func GenerateStructValidationError(err error) []validationError {
	var result []validationError

	for _, err := range err.(validator.ValidationErrors) {
		e := validationError{
			Field:    err.Field(),
			Messsage: err.Error(),
		}

		result = append(result, e)
	}

	return result
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "can not be empty!"
	case "max":
		return "should be less than or equal to " + ToCamelCase(fe.Param()) + "!"
	case "min":
		return "should be greater than or equal to " + ToCamelCase(fe.Param()) + "!"
	case "gte":
		return "should be greater than or equal to " + ToCamelCase(fe.Param()) + "!"
	case "gt":
		return "should be greater than " + ToCamelCase(fe.Param()) + "!"
	case "lte":
		return "should be less than or equal to " + ToCamelCase(fe.Param()) + "!"
	case "email":
		return "must be a valid email address!"
	case "eqfield":
		return "does not match with " + ToCamelCase(fe.Param()) + "!"
	case "ltfield":
		return "must be less than " + ToCamelCase(fe.Param()) + " field!"
	case "gtfield":
		return "must be greater than " + ToCamelCase(fe.Param()) + " field!"
	case "alpha":
		return "must be entirely alphabetic characters!"
	case "alphanum":
		return "must be entirely alpha-numeric characters!"
	case "numeric":
		return "must be an integer!"
	case "oneof":
		return "must be one of " + strings.Replace(fe.Param(), " ", ", ", -1)
	case "len":
		return "must have a length of " + ToCamelCase(fe.Param()) + "!"
	case "uuid":
		return "not a valid UUID!"
	case "url":
		return "must be a valid URL!"
	}
	return "something is wrong with this field!"
}

func FormatValidation(err error) *localError.GlobalError {
	var result []ErrorMsg

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			result = append(result, ErrorMsg{Field: ToCamelCase(fe.Field()), Message: getErrorMsg(fe)})
		}
	}
	if len(result) == 0 {
		localError.ErrBadRequest("request body cannot be empty!", nil)
	}

	return localError.ErrBadRequest(result, nil)
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToCamelCase(str string) string {
	return strings.ToLower(str[:1]) + str[1:]
}

func IsValidURL(url string) bool {
	pattern := `^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`
	match, _ := regexp.MatchString(pattern, url)
	return match
}

func ParseInt(value string) int {
	intValue, _ := strconv.Atoi(value)
	return intValue
}
