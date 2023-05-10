package handler

import (
	"github.com/gin-gonic/gin"

	"trimble.com/common/constants"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/model/folder"
	"trimble.com/tdrive/api/services"
	"trimble.com/tdrive/api/utils"
)

// FolderVersionHandler -> struct for folder versions handler
type FolderVersionHandler struct {
	service *services.FolderVersionService
}

// NewFolderVersionHandler creates new FolderVersionHandler
func NewFolderVersionHandler(service *services.FolderVersionService) *FolderVersionHandler {
	return &FolderVersionHandler{
		service: service,
	}
}

// GetFolderVersionDetails Perform get folder version details with the help of service functions
func (f *FolderVersionHandler) GetFolderVersionDetails(c *gin.Context) {

	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)
	// ------------- Path parameter -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	folderId, folderParamErr := utils.ValidatePathParam(c, "folderId")
	if paramErr != nil {
		utils.HandleApiError(c, folderParamErr)
		return
	}

	version, versionParamErr := utils.ValidatePathParam(c, "version")
	if versionParamErr != nil {
		utils.HandleApiError(c, versionParamErr)
		return
	}

	versionId, versionErr := utils.ValidateAndGetVersion(&version)
	if versionErr != nil {
		utils.HandleApiError(c, versionErr)
		return
	}

	var queryParams folder.DetailsParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	response, err := f.service.GetFolderVersionDetails(ctx, spaceId, folderId, *versionId, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.Header(constants.ETag, buildETag(response.Version))
		c.JSON(200, response)
	}
}

// ListFolderVersions Fetch all the versions of given folder
func (f *FolderVersionHandler) ListFolderVersions(c *gin.Context) {

	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)

	// ------------- Path parameters -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	folderId, folderParamErr := utils.ValidatePathParam(c, "folderId")
	if paramErr != nil {
		utils.HandleApiError(c, folderParamErr)
		return
	}

	// ------------- Query parameters -------------
	var queryParams folder.ListVersionsParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	utils.ValidateAndSetDeleteQueryParam(c, queryParams.Deleted)

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	response, err := f.service.ListFolderVersions(ctx, spaceId, folderId, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.JSON(200, response)
	}
}
