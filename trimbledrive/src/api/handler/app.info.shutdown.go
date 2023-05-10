//go:build debug

package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Gracefully shutdown the server.
func (a *AppInfoHandler) Shutdown(c *gin.Context) {
	c.Status(http.StatusAccepted)
	os.Exit(1)
}
