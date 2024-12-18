package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func ValidateNationalID(nationalID string) bool {
	if len(nationalID) != 10 {
		return false
	}

	matched, _ := regexp.MatchString(`^\d{10}$`, nationalID)

	if !matched {
		return false
	}

	check, _ := strconv.Atoi(string(nationalID[9]))
	sum := 0
	for i := 0; i < 9; i++ {
		num, _ := strconv.Atoi(string(nationalID[i]))
		sum += num * (10 - i)
	}

	remaining := sum % 11
	return (remaining < 2 && remaining == check) || (remaining >= 2 && check == 11-remaining)
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type TranslateErrorDetail struct {
	ErrorList []ErrorDetail
	Message   string
	Error     error
}

func TranslateError(err error, obj interface{}) (TranslateErrorDetail, error) {
	var ve validator.ValidationErrors
	if err == nil {
		return TranslateErrorDetail{}, nil
	}

	if errors.As(err, &ve) {
		var errorDetails []ErrorDetail
		var errorMessages []string
		val := reflect.ValueOf(obj)

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		} else {
			return TranslateErrorDetail{}, fmt.Errorf("obj must be a pointer")
		}

		typ := val.Type()

		for _, fe := range ve {
			fieldName := fe.Field()

			field, ok := typ.FieldByName(fe.Field())
			if ok {
				jsonTag := field.Tag.Get("json")
				if jsonTag != "" && jsonTag != "-" {
					jsonParts := strings.Split(jsonTag, ",")
					if len(jsonParts) > 0 && jsonParts[0] != "" {
						fieldName = jsonParts[0]
					}
				} else {
					fieldName = toSnakeCase(fe.Field())
				}
			} else {
				fieldName = toSnakeCase(fe.Field())
			}

			var message string
			switch fe.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", fieldName)
			case "oneof":
				message = fmt.Sprintf("%s must be one of %s", fieldName, fe.Param())
			case "gtfield":
				message = fmt.Sprintf("%s must be greater than %s", fieldName, fe.Param())
			case "gte":
				message = fmt.Sprintf("%s must be greater than or equal to %s", fieldName, fe.Param())
			default:
				message = fmt.Sprintf("%s is not valid", fieldName)
			}

			errorDetails = append(errorDetails, ErrorDetail{
				Field:   fieldName,
				Message: message,
			})
		}
		joinedMessage := strings.Join(errorMessages, "; ")

		return TranslateErrorDetail{
			ErrorList: errorDetails,
			Message:   joinedMessage,
			Error:     errors.New(joinedMessage),
		}, nil
	}
	return TranslateErrorDetail{
		ErrorList: []ErrorDetail{
			{
				Field:   "error",
				Message: err.Error(),
			},
		},
		Message: err.Error(),
		Error:   err,
	}, nil
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result = append(result, '_', r+32)
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}
