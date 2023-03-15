package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/auth/service/schema"
	"github.com/outcatcher/anwil/domains/internals/logging"
)

const (
	authHeader = "Authorization"

	tokenPrefix = "Bearer "

	contextKeyUsername = "username"
)

// JWTAuth check JWT and loads user info into Gin context.
func JWTAuth(state schema.WithAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authHeader)

		logger := logging.LoggerFromCtx(c.Request.Context())

		if header == "" {
			logger.Println("empty auth header")

			c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			c.Abort()

			return
		}

		tokenString := strings.TrimPrefix(header, tokenPrefix)

		claims, err := state.Authentication().ValidateToken(tokenString)
		if err != nil {
			logger.Println("error in JWT:", err)

			c.String(http.StatusForbidden, http.StatusText(http.StatusForbidden))
			c.Abort()

			return
		}

		c.Set(contextKeyUsername, claims.Username)
	}
}
