package handler

import (
	"github.com/gin-gonic/gin"

	"trimble.com/common/constants"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/model/file"
	"trimble.com/tdrive/api/services"
	"trimble.com/tdrive/api/utils"
)

// FileHandler -> struct for file handler
type FileHandler struct {
	service *services.FileService
}

// NewFileHandler creates new FileHandler
func NewFileHandler(service *services.FileService) *FileHandler {
	return &FileHandler{
		service: service,
	}
}

// GetFile Perform get file with the help of service functions
func (f *FileHandler) GetFile(c *gin.Context) {

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

	response, err := f.service.GetFileDetails(ctx, spaceId, fileId, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.Header(constants.ETag, buildETag(response.Version))
		c.JSON(200, response)
	}
}

// UpdateFile Perform update file with the help of service functions
func (f *FileHandler) UpdateFile(c *gin.Context) {

	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)

	// ------------- Path parameter -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	fileId, fileIdParamErr := utils.ValidatePathParam(c, "fileId")
	if paramErr != nil {
		utils.HandleApiError(c, fileIdParamErr)
		return
	}

	// ------------- Header parameters -------------
	var headers file.HeaderParams
	if ok := utils.ValidateHeaderParams(c, &headers); !ok {
		// Return if validation failed
		return
	}

	var request file.UpdateFileRequest
	// Validate Request Body
	if ok := utils.ValidateRequestBody(c, &request); !ok {
		// Return if validation failed
		return
	}

	// ------------- Query parameters -------------
	var queryParams file.UpdateParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}
	utils.ValidateAndSetDeleteQueryParam(c, queryParams.Deleted)

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	response, updateErr := f.service.UpdateFile(ctx, spaceId, fileId, request, headers, queryParams)
	if updateErr != nil {
		utils.HandleApiError(c, updateErr)
	} else {
		c.Header(constants.ETag, buildETag(response.Version))
		c.JSON(200, response)
	}

}

// DeleteFile Perform delete file with the help of service functions
func (f *FileHandler) DeleteFile(c *gin.Context) {

	// Get current request context from gin
	ctx := requestcontext.GetContextFromGin(c)

	// ------------- Path parameters -------------
	spaceId, paramErr := utils.ValidatePathParam(c, "spaceId")
	if paramErr != nil {
		utils.HandleApiError(c, paramErr)
		return
	}

	fileId, fileIdParamErr := utils.ValidatePathParam(c, "fileId")
	if paramErr != nil {
		utils.HandleApiError(c, fileIdParamErr)
		return
	}

	// ------------- Header parameters -------------
	var headers file.HeaderParams
	if ok := utils.ValidateHeaderParams(c, &headers); !ok {
		// Return if validation failed
		return
	}

	// ------------- Query parameters -------------
	var queryParams file.RequestParams
	if ok := utils.ValidateQueryParams(c, &queryParams); !ok {
		// Return if validation failed
		return
	}

	if fields := utils.GetFieldsQueryParam(queryParams.Fields); fields != nil {
		queryParams.Fields = fields
	}

	response, err := f.service.DeleteFile(ctx, spaceId, fileId, headers, queryParams)
	if err != nil {
		utils.HandleApiError(c, err)
	} else {
		c.Header(constants.ETag, buildETag(response.Version))
		c.JSON(200, response)
	}
}
