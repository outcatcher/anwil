package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleEcho(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
