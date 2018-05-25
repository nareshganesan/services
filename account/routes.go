package account

import (
	"github.com/gin-gonic/gin"
	mw "github.com/nareshganesan/services/middleware"
)

// RegisterAccount registers account entity related endpoints to router
func RegisterAccount(router *gin.RouterGroup) {
	router.POST("/login", Login)
	router.POST("/signup", Signup)
	router.POST("/update", mw.AuthDecorator(), UpdateAccount)
	router.POST("/delete", mw.AuthDecorator(), DeleteAccount)
	router.POST("/list", mw.AuthDecorator(), ListAccount)
}
