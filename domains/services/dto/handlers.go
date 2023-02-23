package dto

import "github.com/gin-gonic/gin"

type AddHandlersFunc func(baseGroup, secGroup *gin.RouterGroup) error
