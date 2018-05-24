package status

import (
	"github.com/gin-gonic/gin"
)

// Index handler for index endpoint request
func Index(ctx *gin.Context) {
	resp := APIUsecase(ctx)
	resp.Send()
	return
}
