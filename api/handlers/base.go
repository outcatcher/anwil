package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s server) handleEcho(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func (s server) handleAPISpec(c *gin.Context) {
	c.HTML(http.StatusOK, "spec.html", map[string]interface{}{})
}
