package validation

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	v10 "github.com/go-playground/validator/v10"
)

const (
	jwtHeaderKey = "jwt-header"
)

var (
	// Base regex for JWT is taken from go-playground validator.
	jwtHeaderRe = regexp.MustCompile(`^Bearer\s[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)

	errorMessagesHeader = map[string]string{
		"required":   "value of header '%s' is missing",
		jwtHeaderKey: fmt.Sprintf("value of header '%%s' not matching pattern %s", jwtHeaderRe),
	}

	headerValidate = v10.New()

	singleHeaderInit = sync.Once{}
)

func getHeaderTag(v any, fieldName string) string {
	typ := reflect.TypeOf(v)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fld, ok := typ.FieldByName(fieldName)
	if !ok {
		return ""
	}

	return fld.Tag.Get("header")
}

// ValidateHeaderCtx validates structure of filled header structure.
func ValidateHeaderCtx(ctx context.Context, v any) error {
	// making sure validation registration is done only once to avoid concurrent map writes
	// NOTE: maybe it makes sense to move this to `init` function
	singleHeaderInit.Do(func() {
		err := headerValidate.RegisterValidationCtx(jwtHeaderKey, func(ctx context.Context, fl v10.FieldLevel) bool {
			return jwtHeaderRe.MatchString(fl.Field().String())
		})
		if err != nil {
			panic(err)
		}
	})

	if err := headerValidate.StructCtx(ctx, v); err != nil {
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
			headerTag := getHeaderTag(v, fieldError.StructField())

			if headerTag == "" {
				fieldMessages[i] = fmt.Sprintf(
					"value of field '%s' is invalid: %s",
					fieldError.Field(), fieldError.Tag(),
				)

				continue
			}

			msg := errorMessagesHeader[fieldError.Tag()] // message for exact failure tag
			if msg == "" {
				msg = "value of header '%s' is invalid"
			}

			fieldMessages[i] = fmt.Sprintf(msg, headerTag)
		}

		return fmt.Errorf("%w:\n\t%s", ErrValidationFailed, strings.Join(fieldMessages, "\n\t"))
	}

	return nil
}
