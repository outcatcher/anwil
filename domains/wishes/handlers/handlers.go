package handlers

import (
	"github.com/labstack/echo/v4"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/wishes/service/schema"
)

// AddHandlers - adds wishlist-related endpoints.
func AddHandlers(state schema.WithWishlist) services.AddHandlersFunc {
	return func(baseGroup, secGroup *echo.Group) error {
		users := state.Users()

		baseGroup.POST("/login", handleAuthorize(users))
		baseGroup.POST("/wisher", handleUserRegister(users))

		return nil
	}
}
