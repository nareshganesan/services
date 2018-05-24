package status

import (
	"github.com/gin-gonic/gin"
)

// RegisterStatus registers status related endpoints to router
func RegisterStatus(router *gin.RouterGroup) {
	router.GET("/", Index)
}
