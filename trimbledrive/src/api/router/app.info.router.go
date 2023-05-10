package router

import (
	"github.com/gin-gonic/gin"

	"trimble.com/tdrive/api/handler"
)

type AppInfoRouter struct {
	handler *handler.AppInfoHandler
}

// NewAppInfoRouter -> new instance of AppInfoRouter
func NewAppInfoRouter(appInfoHandler *handler.AppInfoHandler) *AppInfoRouter {
	return &AppInfoRouter{
		handler: appInfoHandler,
	}
}

// Register registers AppInfo routes to router group
func (a *AppInfoRouter) Register(routerGroup *gin.RouterGroup) {
	appInfo := routerGroup.Group("app")
	{
		appInfo.GET("/health", a.handler.GetAppHealth)
		appInfo.GET("/health/ping", a.handler.Ping)
		appInfo.HEAD("/health/ping", a.handler.Ping)
		appInfo.GET("/shutdown", a.handler.Shutdown)
	}
}
