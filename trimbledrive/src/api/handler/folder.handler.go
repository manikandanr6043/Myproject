package handler

import (
	"github.com/gin-gonic/gin"

	"trimble.com/common/constants"
	"trimble.com/common/requestcontext"

	"trimble.com/tdrive/api/model/folder"
	"trimble.com/tdrive/api/services"
	"trimble.com/tdrive/api/utils"
)

// FolderHandler -> struct for folder handler
type FolderHandler struct {
	service *services.FolderService
}

// NewFolderHandler creates new FolderHandler
func NewFolderHandler(service *services.FolderService) *FolderHandler {
	return &FolderHandler{
		service: service,
	}
}

// CreateFolder Creates a folder under a given parent
func (f *FolderHandler) CreateFolder(c *gin.Context) {
	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)

	// ------------- Path parameter "spaceId" -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	var request folder.CreateRequest
	// Validate Request Body
	if ok := utils.ValidateRequestBody(c, &request); !ok {
		// Return if validation failed
		return
	}
	// ------------- Query parameters -------------
	var queryParams folder.RequestParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	// Perform create folder
	response, createError := f.service.CreateFolder(ctx, spaceId, request, queryParams)
	if createError == nil {
		c.Header(constants.ETag, buildETag(response.Version))
		host := utils.GetHost(c.Request)
		location := "https://" + host + "/v1/spaces/" + (*response).SpaceId + "/folders/" + (*response).Id
		c.Header(constants.Location, location)
		c.Header(constants.ContentLocation, location)
		c.JSON(201, response)
	} else {
		utils.HandleApiError(c, createError)
	}
}

// GetFolder Perform get folder with the help of service functions
func (f *FolderHandler) GetFolder(c *gin.Context) {

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

	// ------------- Query parameters -------------
	var queryParams folder.DetailsParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	utils.ValidateAndSetDeleteQueryParam(c, queryParams.Deleted)

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	response, err := f.service.GetFolder(ctx, spaceId, folderId, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.Header(constants.ETag, buildETag(response.Version))
		c.JSON(200, response)
	}
}

// UpdateFolder Perform update folder with the help of service functions
func (f *FolderHandler) UpdateFolder(c *gin.Context) {

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

	// ------------- Header parameters -------------
	var headers folder.HeaderParams
	if ok := utils.ValidateHeaderParams(c, &headers); !ok {
		// Return if validation failed
		return
	}

	var request folder.UpdateRequest
	// Validate Request Body
	if ok := utils.ValidateRequestBody(c, &request); !ok {
		// Return if validation failed
		return
	}

	//------------- Query Params ---------------------
	var queryParams folder.UpdateParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		return
	}
	utils.ValidateAndSetDeleteQueryParam(c, queryParams.Deleted)

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}
	response, updateErr := f.service.UpdateFolder(ctx, spaceId, folderId, request, headers, queryParams)
	if updateErr != nil {
		utils.HandleApiError(c, updateErr)
	} else {
		c.Header(constants.ETag, buildETag(response.Version))
		c.JSON(200, response)
	}

}

// DeleteFolder Perform delete folder with the help of service functions
func (f *FolderHandler) DeleteFolder(c *gin.Context) {

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

	// ------------- Header parameters -------------
	var headers folder.HeaderParams
	if ok := utils.ValidateHeaderParams(c, &headers); !ok {
		// Return if validation failed
		return
	}
	//------------- Query Params ---------------------
	var queryParams folder.DeleteParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		return
	}

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	response, err := f.service.DeleteFolder(ctx, spaceId, folderId, headers, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.Header(constants.ETag, buildETag(response.Version))
		c.JSON(200, response)
	}
}

// GetChildrenUnderParentFolder Fetch all the children under given parent folder
func (f *FolderHandler) GetChildrenUnderParentFolder(c *gin.Context) {

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
	var queryParams folder.ListParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	utils.ValidateAndSetDeleteQueryParam(c, queryParams.Deleted)

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	response, err := f.service.GetChildren(ctx, spaceId, folderId, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.JSON(200, response)
	}
}

// buildETag Internal util function to build etag header
func buildETag(version string) string {
	return "W/" + version
}
