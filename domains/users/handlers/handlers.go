/*
Package handlers contains API handlers for user-related endpoints.
*/
package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/core/validation"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

// AddUserHandlers - adds user-related endpoints.
func AddUserHandlers(state schema.WithUsers) services.AddHandlersFunc {
	return func(baseGroup, secGroup *echo.Group) error {
		users := state.Users()

		baseGroup.POST("/login", handleAuthorize(users))

		baseGroup.POST("/wisher", handleUserRegister(users))

		return nil
	}
}

// bindAndValidateJSON binds request body to structure,
// validates it using `validate` tag and trows errors into gin context.
func bindAndValidateJSON(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return fmt.Errorf("error binding JSON: %w", err)
	}

	if err := validation.ValidateJSONCtx(c.Request().Context(), req); err != nil {
		return fmt.Errorf("error validating JSON: %w", err)
	}

	return nil
}
