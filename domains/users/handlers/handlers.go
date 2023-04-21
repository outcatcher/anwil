/*
Package handlers contains API handlers for user-related endpoints.
*/
package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/core/services"
	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/core/validation"
	usersSchema "github.com/outcatcher/anwil/domains/users/service/schema"
)

// AddUserHandlers - adds user-related endpoints.
func AddUserHandlers(state svcSchema.ProvidingServices) svcSchema.AddHandlersFunc {
	return func(baseGroup, secGroup *echo.Group) error {
		userService, err := services.GetServiceFromProvider[usersSchema.UserService](state, usersSchema.ServiceID)
		if err != nil {
			return fmt.Errorf("error adding user hanlders: %w", err)
		}

		baseGroup.POST("/login", handleAuthorize(userService))
		baseGroup.POST("/wisher", handleUserRegister(userService))

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
