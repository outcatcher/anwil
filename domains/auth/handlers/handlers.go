/*
Package handlers contains auth-related endpoints.
*/
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/auth/service/schema"
	services "github.com/outcatcher/anwil/domains/internals/services/schema"
)

// AddAuthHandlers adds authentication-related endpoints.
func AddAuthHandlers(_ schema.WithAuth) services.AddHandlersFunc {
	return func(baseGroup, _ *gin.RouterGroup) error {
		return nil
	}
}
