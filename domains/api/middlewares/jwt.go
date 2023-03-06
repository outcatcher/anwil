package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/auth/dto"
)

const (
	authHeader = "Authorization"

	tokenPrefix = "Bearer "

	contextKeyUsername = "username"
)

// JWTAuth check JWT and loads user info into Gin context.
func JWTAuth(state dto.WithAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authHeader)

		if header == "" {
			log.Println("empty auth header")

			c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			c.Abort()

			return
		}

		tokenString := strings.TrimPrefix(header, tokenPrefix)

		claims, err := state.Authentication().ValidateToken(tokenString)
		if err != nil {
			log.Println("error in JWT:", err)

			c.String(http.StatusForbidden, http.StatusText(http.StatusForbidden))
			c.Abort()

			return
		}

		c.Set(contextKeyUsername, claims.Username)
	}
}
