package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleEcho(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func handleAPISpec(c *gin.Context) {
	c.HTML(http.StatusOK, "spec.html", map[string]interface{}{})
}
