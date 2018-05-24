package shared

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	g "github.com/nareshganesan/services/globals"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// GetHeaderKey returns the value for the key from request header
// else empty string
func GetHeaderKey(ctx *gin.Context, key string) string {
	l := g.Gbl.Log
	if values, _ := ctx.Request.Header[key]; len(values) > 0 {
		return values[0]
	}
	l.WithFields(logrus.Fields{
		"key": key,
	}).Info("key not present in header")
	return ""
}

// SetHeaderKey sets the key and value to request header
func SetHeaderKey(ctx *gin.Context, key, value string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set(key, value)
		ctx.Next()
	}
}

// GetAuthorizationHeader helper to get authorization header from request
func GetAuthorizationHeader(ctx *gin.Context) string {
	l := g.Gbl.Log
	authHeader := GetHeaderKey(ctx, "Authorization")
	// extract the token
	if authHeader != "" {
		tokenVal := strings.Fields(authHeader)
		if len(tokenVal) > 1 {
			token := tokenVal[1]
			l.Error("Auhorization token is present in header!")
			return token
		} else if len(tokenVal) > 0 {
			l.Error("Auhorization token is broken! (no Bearer present)")
			return ""
		} else {
			l.Error("Auhorization token is not present")
			return ""
		}
	} else {
		l.Info("Authorization header / token is not present")
		return ""
	}
}

// ValidateRequest binds request to the fields of the supplied interface
// else returns invalid response and error
func ValidateRequest(ctx *gin.Context, obj interface{}) (*Response, error) {
	if err := ctx.ShouldBind(obj); err != nil {
		fmt.Println("ValidateRequest: Error")
		errMess := fmt.Sprintf("%s", err.Error())
		fmt.Println(err.Error())
		data := make(map[string]interface{})
		data["error"] = errMess
		data["message"] = "BadRequest"
		data["code"] = http.StatusBadRequest
		invResp := &Response{
			Ctx:  ctx,
			Data: data,
		}
		return invResp, err
	}
	return nil, nil
}

// GetRequestData binds the supplied interface with request json
// else returns invalid response and error
func GetRequestData(ctx *gin.Context, obj interface{}) (*Response, error) {
	if err := json.NewDecoder(ctx.Request.Body).Decode(obj); err != nil {
		fmt.Println("GetRequestData: Error")
		errMess := fmt.Sprintf("%s", err.Error())
		data := make(map[string]interface{})
		data["error"] = errMess
		data["message"] = "BadRequest"
		data["code"] = http.StatusBadRequest
		invResp := &Response{
			Ctx:  ctx,
			Data: data,
		}
		return invResp, err
	}
	return nil, nil
}
