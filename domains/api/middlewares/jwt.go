package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/internals/logging"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

const (
	authHeader = "Authorization"

	tokenPrefix = "Bearer "

	contextKeyUsername = "username"
)

// JWTAuth check JWT and loads user info into Gin context.
func JWTAuth(state schema.WithUsers) gin.HandlerFunc {
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

		claims, err := state.Users().ValidateUserToken(c.Request.Context(), tokenString)
		if err != nil {
			logger.Println("error in JWT:", err)

			c.String(http.StatusForbidden, http.StatusText(http.StatusForbidden))
			c.Abort()

			return
		}

		c.Set(contextKeyUsername, claims.Username)
	}
}
