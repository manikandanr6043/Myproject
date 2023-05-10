package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellofresh/health-go/v5"
	healthMongo "github.com/hellofresh/health-go/v5/checks/mongo"

	"trimble.com/tdrive/api/config"
)

// AppInfoHandler -> struct for app info handler
type AppInfoHandler struct {
	h *health.Health
}

// Version The version string is expected to be set on the compile/link time by the build workflow.
// The default value here represents a version for binaries built on developer machine (unofficial build).
var Version = "v0.0.0"

// NewAppInfoHandler creates new app info handler
func NewAppInfoHandler(config *config.TDriveConfig) *AppInfoHandler {
	// add checks on instance creation
	h, _ := health.New(
		health.WithSystemInfo(),
		health.WithComponent(health.Component{
			Name:    "TrimbleDrive",
			Version: Version}),
		health.WithChecks(health.Config{
			Name: "mongodb",
			Check: healthMongo.New(healthMongo.Config{
				DSN: config.Mongo.Uri,
			})},
		))
	return &AppInfoHandler{
		h: h,
	}
}

// Deep health check
func (a *AppInfoHandler) GetAppHealth(c *gin.Context) {
	a.h.HandlerFunc(c.Writer, c.Request)
}

// Lightweight shallow ping type of health check for load balancing
func (a *AppInfoHandler) Ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
