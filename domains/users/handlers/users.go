package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

type createUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name"`
}

func handleUserRegister(usr schema.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		req := new(createUser)

		if err := bindAndValidateJSON(c, req); err != nil {
			return fmt.Errorf("error registering user: %w", err)
		}

		err := usr.SaveUser(ctx, schema.User{
			Username: req.Username,
			Password: req.Password,
			FullName: req.FullName,
		})
		if err != nil {
			return fmt.Errorf("error registering user: %w", err)
		}

		return c.NoContent(http.StatusCreated)
	}
}
