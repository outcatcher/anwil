package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	users "github.com/outcatcher/anwil/domains/users/dto"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

type credentialsRequest struct {
	// Username
	Username string `json:"username" validate:"required"`
	// SHA256 - HMAC encrypted password
	Password string `json:"password" validate:"required"`
}

type jwtResponse struct {
	Token string `json:"token"`
}

func handleAuthorize(usr schema.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(credentialsRequest)

		if err := bindAndValidateJSON(c, req); err != nil {
			return
		}

		user := users.User{
			Username: req.Username,
			Password: req.Password,
		}

		tok, err := usr.GenerateUserToken(c.Request.Context(), user)
		if err != nil {
			_ = c.Error(err)

			return
		}

		c.AbortWithStatusJSON(http.StatusOK, jwtResponse{Token: tok})
	}
}
