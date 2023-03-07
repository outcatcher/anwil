/*
Package handlers contains auth-related endpoints.
*/
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/auth/dto"
	services "github.com/outcatcher/anwil/domains/services/dto"
)

// AddAuthHandlers adds authentication-related endpoints.
func AddAuthHandlers(_ dto.WithAuth) services.AddHandlersFunc {
	return func(baseGroup, _ *gin.RouterGroup) error {
		return nil
	}
}
