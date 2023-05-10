package utils

import (
	"github.com/gin-gonic/gin"

	"trimble.com/common/api_error"

	"trimble.com/tdrive/api/model"
)

func HandleApiError(c *gin.Context, e *api_error.ApiError) {
	c.JSON(e.StatusCode, model.New(e.ErrorCode, e.ErrorMessage))
}

func HandleApiAbortError(c *gin.Context, e *api_error.ApiError) {
	c.AbortWithStatusJSON(e.StatusCode, model.New(e.ErrorCode, e.ErrorMessage))
}
