package validation

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	v10 "github.com/go-playground/validator/v10"
)

var (
	errorMessagesJSON = map[string]string{
		"required": "value for field '%s' is missing",
	}

	// ErrValidationFailed - error for failed validations.
	ErrValidationFailed = errors.New("validation failed")

	jsonValidate = v10.New()
)

func getJSONTag(v any, fieldName string) string {
	typ := reflect.TypeOf(v)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fld, ok := typ.FieldByName(fieldName)
	if !ok {
		return ""
	}

	return fld.Tag.Get("json")
}

// ValidateJSONCtx validates structure.
//
// Error message contains JSON tags instead of structure field names.
func ValidateJSONCtx(ctx context.Context, v any) error {
	if err := jsonValidate.StructCtx(ctx, v); err != nil {
		errInvalidValidation := new(v10.InvalidValidationError)

		if errors.As(err, &errInvalidValidation) {
			return fmt.Errorf("invalid validation: %w", err)
		}

		errValidationErrors := new(v10.ValidationErrors)
		if !errors.As(err, errValidationErrors) {
			return fmt.Errorf("error during validation: %w", err)
		}

		fieldMessages := make([]string, len(*errValidationErrors))

		for i, fieldError := range *errValidationErrors {
			jsonTag := getJSONTag(v, fieldError.StructField())

			if jsonTag == "" {
				fieldMessages[i] = fmt.Sprintf(
					"value of field '%s' is invalid: %s",
					fieldError.Field(), fieldError.Tag(),
				)

				continue
			}

			msg := errorMessagesJSON[fieldError.Tag()] // message for exact failure tag
			if msg == "" {
				msg = "value of field '%s' is invalid"
			}

			fieldMessages[i] = fmt.Sprintf(msg, jsonTag)
		}

		return fmt.Errorf("%w:\n\t%s", ErrValidationFailed, strings.Join(fieldMessages, "\n\t"))
	}

	return nil
}
