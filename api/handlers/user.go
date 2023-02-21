package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type credentials struct {
	// Username
	Username string `json:"username"`
	// SHA256 - HMAC encrypted password
	Password string `json:"password"`
}

func handleAuthorize(c *gin.Context) {
	creds := new(credentials)

	if err := c.Bind(creds); err != nil {
		log.Println(err)

		c.AbortWithStatus(http.StatusForbidden)

		return
	}
}
