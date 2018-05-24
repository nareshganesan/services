package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	g "github.com/nareshganesan/services/globals"
	"github.com/sirupsen/logrus"
	"time"
)

// LogrusMiddleware method for logging request info
func LogrusMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := g.GetGlobals().Log
		start := time.Now().UTC()
		path := ctx.Request.URL.Path

		ctx.Next()

		end := time.Now().UTC()
		latency := time.Since(start)

		requestID := ctx.GetString("RequestID")
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		comment := ctx.Errors.String()
		userAgent := ctx.Request.UserAgent()

		timeFormatted := end.Format(time.RFC3339)

		msg := fmt.Sprintf(
			"%s %s %s \"%s %s\" %d %s %s",
			clientIP,
			requestID,
			timeFormatted,
			method,
			path,
			statusCode,
			latency,
			userAgent,
		)

		l.WithFields(logrus.Fields{
			"request-id": requestID,
			"time":       timeFormatted,
			"method":     method,
			"path":       path,
			"latency":    latency,
			"ip":         clientIP,
			"comment":    comment,
			"status":     statusCode,
			"user-agent": userAgent,
		}).Info(msg)

	}
}
