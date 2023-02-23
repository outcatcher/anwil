/*
Package handlers contains API handlers for user-related endpoints.
*/
package handlers

import (
	"github.com/gin-gonic/gin"
	services "github.com/outcatcher/anwil/domains/services/dto"
	"github.com/outcatcher/anwil/domains/users/dto"
)

// AddUserHandlers - adds user-related endpoints.
func AddUserHandlers(state dto.WithUsers) services.AddHandlersFunc {
	return func(baseGroup, secGroup *gin.RouterGroup) error {
		baseGroup.POST("/wisher", handleUserRegister(state.Users()))

		return nil
	}
}
