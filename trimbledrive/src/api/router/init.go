package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewAppInfoRouter),
	fx.Provide(NewFileSpaceRouter),
	fx.Provide(NewFolderRouter),
	fx.Provide(NewUploadsRouter),
	fx.Provide(NewFileRouter),
	fx.Provide(NewRouters),
)

type Routers []Router

// Router interface
type Router interface {
	Register(routerGroup *gin.RouterGroup)
}

// NewRouters Routers to be registered
func NewRouters(appInfoRouter *AppInfoRouter, filespaceRouter *FileSpaceRouter, folderRouter *FolderRouter,
	uploadsRouter *UploadsRouter, fileRouter *FileRouter) Routers {
	return Routers{
		appInfoRouter,
		filespaceRouter,
		folderRouter,
		uploadsRouter,
		fileRouter,
	}
}

// Register all the routers
func (r Routers) Register(routerGroup *gin.RouterGroup) {
	for _, route := range r {
		route.Register(routerGroup)
	}
}
