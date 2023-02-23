package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireJSON forces usage of application/json content type in POST and PUT request to servers.
func RequireJSON(c *gin.Context) {
	typ := c.ContentType()

	if (c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut) &&
		typ != gin.MIMEJSON {
		c.String(http.StatusBadRequest, "invalid MIME type, expected %s", gin.MIMEJSON)
		c.Abort()

		return
	}

	c.Status(http.StatusOK) // temporary status
}
