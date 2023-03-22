/*
Package handlers contains API handlers for user-related endpoints.
*/
package handlers

import (
	"github.com/gin-gonic/gin"
	services "github.com/outcatcher/anwil/domains/internals/services/schema"
	"github.com/outcatcher/anwil/domains/internals/validation"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

// AddUserHandlers - adds user-related endpoints.
func AddUserHandlers(state schema.WithUsers) services.AddHandlersFunc {
	return func(baseGroup, secGroup *gin.RouterGroup) error {
		users := state.Users()

		baseGroup.POST("/login", handleAuthorize(users))

		baseGroup.POST("/wisher", handleUserRegister(users))

		return nil
	}
}

// bindAndValidateJSON binds request body to structure,
// validates it using `validate` tag and trows errors into gin context.
func bindAndValidateJSON(c *gin.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return c.Error(err)
	}

	if err := validation.ValidateJSONCtx(c.Request.Context(), req); err != nil {
		return c.Error(err)
	}

	return nil
}
