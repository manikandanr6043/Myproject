package router

import (
	"github.com/gin-gonic/gin"

	"trimble.com/tdrive/api/handler"
	"trimble.com/tdrive/api/middleware"
)

type FileRouter struct {
	handler         *handler.FileHandler
	versionsHandler *handler.FileVersionHandler
	authMiddleware  middleware.AuthMiddleware
}

// NewFileRouter -> new instance of FileRouter
func NewFileRouter(fileHandler *handler.FileHandler, versionsHandler *handler.FileVersionHandler, authMiddleware middleware.AuthMiddleware) *FileRouter {
	return &FileRouter{
		handler:         fileHandler,
		versionsHandler: versionsHandler,
		authMiddleware:  authMiddleware,
	}
}

// Register registers files routes to router group
func (f *FileRouter) Register(routerGroup *gin.RouterGroup) {
	file := routerGroup.Group("/spaces/:spaceId/files/:fileId", f.authMiddleware.Authenticate())
	{
		file.GET("", f.handler.GetFile)
		file.PATCH("", f.handler.UpdateFile)
		file.DELETE("", f.handler.DeleteFile)
		file.GET("/versions", f.versionsHandler.ListFileVersions)
		file.GET("/versions/:version", f.versionsHandler.GetFileVersionDetails)

	}
}
