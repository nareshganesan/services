package status

import (
	"github.com/gin-gonic/gin"
	"github.com/nareshganesan/services/shared"
	"net/http"
)

// APIUsecase handles the processing for  API status request
func APIUsecase(ctx *gin.Context) *shared.Response {
	data := make(map[string]interface{})
	data["message"] = "Index page for services app"
	data["code"] = http.StatusOK
	return shared.GetResponse(ctx, data)
}
