package router

import (
	"github.com/gin-gonic/gin"

	"trimble.com/tdrive/api/handler"
	"trimble.com/tdrive/api/middleware"
)

type FileSpaceRouter struct {
	handler        *handler.FileSpaceHandler
	authMiddleware middleware.AuthMiddleware
}

// NewFileSpaceRouter -> new instance of FileSpaceRouter
func NewFileSpaceRouter(filespaceHandler *handler.FileSpaceHandler, authMiddleware middleware.AuthMiddleware) *FileSpaceRouter {
	return &FileSpaceRouter{
		handler:        filespaceHandler,
		authMiddleware: authMiddleware,
	}
}

// Register registers FileSpace routes to router group
func (f *FileSpaceRouter) Register(routerGroup *gin.RouterGroup) {
	filespace := routerGroup.Group("spaces", f.authMiddleware.Authenticate())
	{
		filespace.POST("", f.handler.CreateFileSpace)

		filespace.GET("/:spaceId", f.handler.GetFileSpace)

		filespace.PATCH("/:spaceId", f.handler.UpdateFileSpace)

		filespace.DELETE("/:spaceId", f.handler.DeleteFileSpace)
	}
}
