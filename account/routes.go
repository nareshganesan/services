package account

import (
	"github.com/gin-gonic/gin"
)

// RegisterAccount registers account entity related endpoints to router
func RegisterAccount(router *gin.RouterGroup) {
	router.POST("/login", Login)
	router.POST("/signup", Signup)
	router.POST("/update", UpdateAccount)
	router.POST("/delete", DeleteAccount)
	router.POST("/list", ListAccount)
}
