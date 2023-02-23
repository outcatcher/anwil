package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/outcatcher/anwil/domains/auth/dto"
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

func handleAuthorize(usr users.Service, authentication auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(credentialsRequest)

		if err := c.Bind(req); err != nil {
			log.Println(err)

			c.Abort()

			return
		}

		if req.Username == "" || req.Password == "" {
			c.String(http.StatusBadRequest, "missing credentialsRequest")
			c.Abort()

			return
		}

		user, err := usr.GetUser(c.Request.Context(), req.Username)
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusUnauthorized, "no user %s exists", req.Username)
			c.Abort()

			return
		}

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		err = authentication.ValidatePassword(req.Password, user.Password)
		if err != nil {
			c.String(http.StatusUnauthorized, "invalid password")
			c.Abort()

			return
		}

		tok, err := authentication.GenerateToken(&auth.Claims{Username: req.Username})
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		c.JSON(http.StatusOK, jwtResponse{Token: tok})
	}
}
