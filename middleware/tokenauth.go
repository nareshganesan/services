package middleware

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	g "github.com/nareshganesan/services/globals"
	"github.com/nareshganesan/services/shared"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// AuthMiddleware checks whether the user is authenticated or not
// based on Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := g.Gbl.Log
		// get authorization header
		token := shared.GetAuthorizationHeader(ctx)
		if token != "" {
			if claims := g.ParseJWT(token); claims != nil {
				l.Info("Token is valid")
				uid := (*claims)["uid"].(string)
				// uid := int(uidFloat)
				ctx.Set("uid", uid)
				ctx.Set("isAuthenticated", "true")
				ctx.Writer.Header().Set("Authorization", "Bearer "+token)
				l.WithFields(logrus.Fields{
					"uid":             uid,
					"isAuthenticated": strconv.FormatBool(true),
					"token":           token,
				}).Info("User Authenticated")
			} else {
				l.WithFields(logrus.Fields{
					"token": token,
				}).Info("Token is invalid")
				ctx.Set("isAuthenticated", "false")
				ctx.Set("authStatus", "is invalid")
			}
		} else {
			l.Info("Anonymous user")
			ctx.Set("isAuthenticated", "false")
			ctx.Set("authStatus", "cannot be empty")
		}
		ctx.Next()
	}
}

// AuthDecorator ensures that a request will be aborted with an error
// if the user is not authenticated
// Ref: https://github.com/demo-apps/go-gin-app/blob/master/middleware.auth.go
func AuthDecorator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// the user is not logged in
		isAuthenticatedInterface, _ := ctx.Get("isAuthenticated")
		isAuthenticated := isAuthenticatedInterface.(string)
		if isAuthenticated == "false" {
			authStatusInterface, _ := ctx.Get("authStatus")
			authStatus := authStatusInterface.(string)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization Header Token " + authStatus + ". Ex: Authorization: Bearer xyzasdf",
				"message": "Unauthorized request",
				"code":    http.StatusUnauthorized,
				"status":  http.StatusText(http.StatusUnauthorized),
			})
		} else {
			ctx.Next()
		}
	}
}
