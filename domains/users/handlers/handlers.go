/*
Package handlers contains API handlers for user-related endpoints.
*/
package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/core/validation"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

// AddUserHandlers - adds user-related endpoints.
func AddUserHandlers(state schema.WithUsers) services.AddHandlersFunc {
	return func(baseGroup, secGroup fiber.Router) error {
		users := state.Users()

		baseGroup.Post("/login", handleAuthorize(users))

		baseGroup.Post("/wisher", handleUserRegister(users))

		return nil
	}
}

// bindAndValidateJSON binds request body to structure,
// validates it using `validate` tag and trows errors into gin context.
func bindAndValidateJSON(c *fiber.Ctx, req any) error {
	if err := c.BodyParser(req); err != nil {
		return fmt.Errorf("error binding request body: %w", err)
	}

	if err := validation.ValidateJSONCtx(c.Context(), req); err != nil {
		return fmt.Errorf("error validating request body: %w", err)
	}

	return nil
}
