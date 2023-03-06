package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	users "github.com/outcatcher/anwil/domains/users/dto"
)

type credentialsRequest struct {
	// Username
	Username string `json:"username"`
	// SHA256 - HMAC encrypted password
	Password string `json:"password"`
}

type jwtResponse struct {
	Token string `json:"token"`
}

func handleAuthorize(usr users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(credentialsRequest)

		if err := c.Bind(req); err != nil {
			log.Println(err)

			c.Abort()

			return
		}

		if req.Username == "" || req.Password == "" {
			c.String(http.StatusBadRequest, "missing credentials")
			c.Abort()

			return
		}

		user := users.User{
			Username: req.Username,
			Password: req.Password,
		}

		tok, err := usr.GetUserToken(c.Request.Context(), user)
		if err != nil {
			_ = c.Error(err)

			return
		}

		c.JSON(http.StatusOK, jwtResponse{Token: tok})
	}
}
