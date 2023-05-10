package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/model/upload"
	"trimble.com/tdrive/api/services"
	"trimble.com/tdrive/api/utils"
)

// UploadsHandler -> struct for uploads handler
type UploadsHandler struct {
	service *services.UploadService
}

// NewUploadsHandler creates new UploadsHandler
func NewUploadsHandler(service *services.UploadService) *UploadsHandler {
	return &UploadsHandler{
		service: service,
	}
}

// FileUpload Perform file upload with the help of service functions
func (u *UploadsHandler) FileUpload(c *gin.Context) {
	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)

	// ------------- Path parameter "spaceId" -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	var request upload.FileUploadRequest
	// Validate Request Body
	if ok := utils.ValidateRequestBody(c, &request); !ok {
		// Return if validation failed
		return
	}

	if request.Id == nil && request.Name == nil {
		utils.HandleApiError(c, api_error.InvalidUploadPayload)
		return
	}

	var queryParams upload.FileUploadParams
	// Validate query parameters
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	var headers upload.FileUploadHeaders
	// Validate query parameters
	if ok := utils.ValidateHeaderParams(c, &headers); !ok {
		// Return if validation failed
		return
	}
	fileUploadResponse, uploadErr := u.service.FileUpload(ctx, spaceId, request, headers, queryParams)
	if uploadErr != nil {
		utils.HandleApiError(c, uploadErr)
		return
	}
	c.Header(constants.Location, fmt.Sprintf("https://"+c.Request.Host+"/v1/spaces/%s/upload/uploads/%s", fileUploadResponse.Input.SpaceId, fileUploadResponse.UploadId))
	c.JSON(202, fileUploadResponse)
}

// GetFileUploadDetails Perform get file upload details help of service functions
func (u *UploadsHandler) GetFileUploadDetails(c *gin.Context) {
	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)

	// ------------- Path parameter "spaceId" -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	// ------------- Path parameter "uploadId" -------------
	uploadId, paramErr := utils.ValidatePathParam(c, "uploadId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	var queryParams upload.GetFileUploadDetailsParams
	// Validate query parameters
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}
	uploadDetails, err := u.service.GetFileUploadDetails(ctx, spaceId, uploadId, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
		return
	}
	c.JSON(200, uploadDetails)
}
