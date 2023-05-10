package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CorsMiddleware -> struct for requestId
type CorsMiddleware struct {
}

// NewCorsMiddleware -> new instance of CorsMiddleware
func NewCorsMiddleware() CorsMiddleware {
	return CorsMiddleware{}
}

// HandleCors -> preflight response
func (m CorsMiddleware) HandleCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			requestedHeaders := c.Request.Header.Get("access-control-request-headers")
			if requestedHeaders != "" {
				c.Header("Access-Control-Allow-Headers", requestedHeaders)
			} else {
				c.Header("Access-Control-Allow-Headers", "Authorization,Origin,Content-Type,Content-Length,Accept,Accept-Encoding,If-Match,If-None-Match")
			}
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD")
			c.Header("Access-Control-Max-Age", "1728000")
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
