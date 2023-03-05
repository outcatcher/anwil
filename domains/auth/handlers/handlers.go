/*
Package handlers contains auth-related endpoints.
*/
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/api/dto"
	services "github.com/outcatcher/anwil/domains/services/dto"
)

// AddAuthHandlers adds authentication-related endpoints.
func AddAuthHandlers(_ dto.State) services.AddHandlersFunc {
	return func(baseGroup, _ *gin.RouterGroup) error {
		return nil
	}
}
