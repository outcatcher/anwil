package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/outcatcher/anwil/domains/core/validation"
)

const applicationJSON = "application/json"

var errInvalidMIMEType = errors.New("invalid MIME type")

// RequireJSON forces usage of application/json content type in POST and PUT request to servers.
func RequireJSON(c *fiber.Ctx) error {
	typ := string(c.Request().Header.ContentType())

	if (c.Method() == http.MethodPost || c.Method() == http.MethodPut) &&
		!strings.HasPrefix(typ, applicationJSON) {
		err := fmt.Errorf(
			"%w: %w, %s expected",
			validation.ErrValidationFailed, errInvalidMIMEType, applicationJSON,
		)

		return err
	}

	return c.Next()
}
