package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	users "github.com/outcatcher/anwil/domains/users/dto"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

type credentialsRequest struct {
	// Username
	Username string `json:"username" validate:"required"`
	// User password
	Password string `json:"password" validate:"required"`
}

type jwtResponse struct {
	Token string `json:"token"`
}

func handleAuthorize(usr schema.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(credentialsRequest)

		if err := bindAndValidateJSON(c, req); err != nil {
			return fmt.Errorf("error authorizing user: %w", err)
		}

		user := users.User{
			Username: req.Username,
			Password: req.Password,
		}

		tok, err := usr.GenerateUserToken(c.Request().Context(), user)
		if err != nil {
			return fmt.Errorf("error authorizing user: %w", err)
		}

		return c.JSON(http.StatusOK, jwtResponse{Token: tok})
	}
}
