//go:build !debug

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Gracefully shutdown the server.
func (a *AppInfoHandler) Shutdown(c *gin.Context) {
	c.Status(http.StatusBadRequest)
}
