package middleware

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// RequestID method encodes a request id into request header if it doesn't exists
func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := ctx.Request.Header.Get("X-Request-Id")

		// Create request id with UUID4
		if requestID == "" {
			uuid4 := uuid.NewV4()
			requestID = uuid4.String()
		}

		// Expose it for use in the application
		ctx.Set("RequestID", requestID)

		// Set X-Request-Id header
		ctx.Writer.Header().Set("X-Request-Id", requestID)
		ctx.Next()
	}
}
