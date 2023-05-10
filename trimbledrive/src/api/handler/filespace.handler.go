package handler

import (
	"github.com/gin-gonic/gin"

	"trimble.com/common/constants"
	"trimble.com/common/requestcontext"

	"trimble.com/tdrive/api/model/filespace"
	"trimble.com/tdrive/api/services"
	"trimble.com/tdrive/api/utils"
)

// FileSpaceHandler -> struct for file space handler
type FileSpaceHandler struct {
	service *services.FilespaceService
}

// NewFileSpaceHandler creates new FileSpaceHandler
func NewFileSpaceHandler(service *services.FilespaceService) *FileSpaceHandler {
	return &FileSpaceHandler{
		service: service,
	}
}

// CreateFileSpace Perform create filespace with the help of service functions.
// validate create filespace request. Handle response
func (f *FileSpaceHandler) CreateFileSpace(c *gin.Context) {
	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)
	var request filespace.CreateRequest
	// Validate Request Body
	if ok := utils.ValidateRequestBody(c, &request); !ok {
		// Return if validation failed
		return
	}
	// Perform create filespace
	response, createError := f.service.CreateFilespace(ctx, request)
	if createError == nil {
		host := utils.GetHost(c.Request)
		location := "https://" + host + "/v1/spaces/" + (*(*response).Id)
		c.Header(constants.Location, location)
		c.Header(constants.ContentLocation, location)
		c.JSON(201, response)
	} else {
		utils.HandleApiError(c, createError)
	}
}

// GetFileSpace Perform get filespace with the help of service functions
func (f *FileSpaceHandler) GetFileSpace(c *gin.Context) {

	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)
	// ------------- Path parameter "spaceId" -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	// ------------- Query parameters -------------
	var queryParams filespace.DetailsParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	utils.ValidateAndSetDeleteQueryParam(c, queryParams.Deleted)
	filespaceResponse, err := f.service.GetFilespace(ctx, spaceId, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.JSON(200, filespaceResponse)
	}
}

// UpdateFileSpace Perform update filespace with the help of service functions
func (f *FileSpaceHandler) UpdateFileSpace(c *gin.Context) {

	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)

	// ------------- Path parameter "spaceId" -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	// ------------- Query parameters -------------
	var queryParams filespace.SpaceParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	var request filespace.UpdateRequest
	// Validate Request Body
	if ok := utils.ValidateRequestBody(c, &request); !ok {
		// Return if validation failed
		return
	}

	utils.ValidateAndSetDeleteQueryParam(c, queryParams.Deleted)

	filespaceResponse, updateErr := f.service.UpdateFilespace(ctx, spaceId, request, queryParams)
	if updateErr != nil {
		utils.HandleApiError(c, updateErr)
	} else {
		c.JSON(200, filespaceResponse)
	}

}

// DeleteFileSpace Perform delete filespace with the help of service functions
func (f *FileSpaceHandler) DeleteFileSpace(c *gin.Context) {

	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)

	// ------------- Path parameter "spaceId" -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	err := f.service.DeleteFilespace(ctx, spaceId)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.Status(204)
	}
}
