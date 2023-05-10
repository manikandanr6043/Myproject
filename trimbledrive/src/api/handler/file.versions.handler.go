package handler

import (
	"github.com/gin-gonic/gin"

	"trimble.com/common/constants"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/model/file"
	"trimble.com/tdrive/api/services"
	"trimble.com/tdrive/api/utils"
)

// FileVersionHandler -> struct for file version handler
type FileVersionHandler struct {
	service *services.FileVersionService
}

// NewFileVersionHandler creates new FileVersionHandler
func NewFileVersionHandler(service *services.FileVersionService) *FileVersionHandler {
	return &FileVersionHandler{
		service: service,
	}
}

// GetFileVersionDetails Perform get file version with the help of service functions
func (f *FileVersionHandler) GetFileVersionDetails(c *gin.Context) {

	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)
	// ------------- Path parameter -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	fileId, fileParamErr := utils.ValidatePathParam(c, "fileId")
	if paramErr != nil {
		utils.HandleApiError(c, fileParamErr)
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

	// ------------- Query parameters -------------
	var queryParams file.DetailsParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	utils.ValidateAndSetDeleteQueryParam(c, queryParams.Deleted)

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	response, err := f.service.GetFileVersionDetails(ctx, spaceId, fileId, *versionId, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.Header(constants.ETag, buildETag(response.Version))
		c.JSON(200, response)
	}
}

// ListFileVersions Fetch all the versions of given file
func (f *FileVersionHandler) ListFileVersions(c *gin.Context) {

	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)

	// ------------- Path parameters -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	folderId, folderParamErr := utils.ValidatePathParam(c, "fileId")
	if paramErr != nil {
		utils.HandleApiError(c, folderParamErr)
		return
	}

	// ------------- Query parameters -------------
	var queryParams file.ListVersionsParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}
	utils.ValidateAndSetDeleteQueryParam(c, queryParams.Deleted)

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	response, err := f.service.ListFileVersions(ctx, spaceId, folderId, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.JSON(200, response)
	}
}
