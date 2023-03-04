package dto

import "github.com/gin-gonic/gin"

// AddHandlersFunc - function adding handlers to secure or/and unsecure groups.
type AddHandlersFunc func(baseGroup, secGroup *gin.RouterGroup) error
