package router

import (
	"github.com/gin-gonic/gin"

	"trimble.com/tdrive/api/handler"
	"trimble.com/tdrive/api/middleware"
)

type UploadsRouter struct {
	handler        *handler.UploadsHandler
	authMiddleware middleware.AuthMiddleware
}

// NewUploadsRouter -> new instance of UploadsRouter
func NewUploadsRouter(filesHandler *handler.UploadsHandler, authMiddleware middleware.AuthMiddleware) *UploadsRouter {
	return &UploadsRouter{
		handler:        filesHandler,
		authMiddleware: authMiddleware,
	}
}

// Register registers uploads routes to router group
func (u *UploadsRouter) Register(routerGroup *gin.RouterGroup) {
	uploads := routerGroup.Group("spaces/:spaceId/uploads", u.authMiddleware.Authenticate())
	{
		uploads.POST("", u.handler.FileUpload)

		uploads.GET("/:uploadId", u.handler.GetFileUploadDetails)

	}
}
