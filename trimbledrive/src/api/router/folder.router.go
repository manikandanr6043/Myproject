package router

import (
	"github.com/gin-gonic/gin"

	"trimble.com/tdrive/api/handler"
	"trimble.com/tdrive/api/middleware"
)

type FolderRouter struct {
	handler         *handler.FolderHandler
	versionsHandler *handler.FolderVersionHandler
	authMiddleware  middleware.AuthMiddleware
}

// NewFolderRouter -> new instance of FolderRouter
func NewFolderRouter(folderHandler *handler.FolderHandler, versionsHandler *handler.FolderVersionHandler, authMiddleware middleware.AuthMiddleware) *FolderRouter {
	return &FolderRouter{
		handler:         folderHandler,
		versionsHandler: versionsHandler,
		authMiddleware:  authMiddleware,
	}
}

// Register registers folders routes to router group
func (f *FolderRouter) Register(routerGroup *gin.RouterGroup) {
	folder := routerGroup.Group("/spaces/:spaceId/folders", f.authMiddleware.Authenticate())
	{
		folder.POST("", f.handler.CreateFolder)

		folder.GET("/:folderId", f.handler.GetFolder)

		folder.PATCH("/:folderId", f.handler.UpdateFolder)

		folder.DELETE("/:folderId", f.handler.DeleteFolder)

		folder.GET("/:folderId/children", f.handler.GetChildrenUnderParentFolder)

		folder.GET("/:folderId/versions", f.versionsHandler.ListFolderVersions)

		folder.GET("/:folderId/versions/:version", f.versionsHandler.GetFolderVersionDetails)
	}
}
