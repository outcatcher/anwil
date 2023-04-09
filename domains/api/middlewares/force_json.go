package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/outcatcher/anwil/domains/core/validation"
)

const applicationJSON = "application/json"

// RequireJSON forces usage of application/json content type in POST and PUT request to servers.
func RequireJSON(c *fiber.Ctx) error {
	typ := string(c.Request().Header.ContentType())

	if (c.Method() == http.MethodPost || c.Method() == http.MethodPut) &&
		!strings.HasPrefix(typ, applicationJSON) {
		err := fmt.Errorf(
			"%w: invalid MIME type, %s expected",
			validation.ErrValidationFailed, applicationJSON,
		)

		return err
	}

	return nil
}
