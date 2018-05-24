package shared

import (
	// "strconv"
	"github.com/gin-gonic/gin"
	g "github.com/nareshganesan/services/globals"
	"net/http"
)

// Render one of HTML, JSON or XML based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present

// Response entity for services app
type Response struct {
	Ctx  *gin.Context
	Data map[string]interface{}
}

// Send helper to send response
// supported response types json, xml
func (r *Response) Send() {
	l := g.Gbl.Log
	statusCode := r.Data["code"].(int)
	r.Data["status"] = http.StatusText(statusCode)
	switch r.Ctx.Request.Header.Get("Accept") {
	case "application/json":
		l.Info("Render custom json response")
		// Respond with JSON
		r.Ctx.JSON(statusCode, r.Data)
	case "application/xml":
		l.Info("Render custom xml response")
		// Respond with XML
		r.Ctx.XML(statusCode, r.Data)
	default:
		l.Info("Render custom default json response")
		// Respond with JSON
		r.Ctx.JSON(statusCode, r.Data)
	}
	r.Ctx.Abort()
	return
}

// GetResponse returns response object give request context and data interface
func GetResponse(ctx *gin.Context, data map[string]interface{}) *Response {
	return &Response{
		Ctx:  ctx,
		Data: data,
	}
}
