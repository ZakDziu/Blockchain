package types

import "github.com/gin-gonic/gin"

type Controller interface {
	RegisterRoutes(group *gin.RouterGroup)
}
