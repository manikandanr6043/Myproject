package middleware

import (
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"trimble.com/common/api_error"
	"trimble.com/common/requestcontext"

	"trimble.com/tdrive/api/utils"
)

// CustomRecoveryMiddleware -> struct for custom recovery middleware
type CustomRecoveryMiddleware struct {
}

// NewCustomRecoveryMiddleware -> new instance of CustomRecoveryMiddleware
func NewCustomRecoveryMiddleware() CustomRecoveryMiddleware {
	return CustomRecoveryMiddleware{}
}

// HandlerCustomRecovery Recovery middleware recovers from any panics and writes a 500 if there was one.
// returns a gin.HandlerFunc (middleware) with a custom recovery handler
func (m CustomRecoveryMiddleware) HandlerCustomRecovery() gin.HandlerFunc {
	//No operation zap logger will be passed to custom recovery to make sure that we log the error using contextual logger.
	return ginzap.CustomRecoveryWithZap(zap.NewNop(), false, func(c *gin.Context, recovered interface{}) {
		ctx := requestcontext.GetContextFromGin(c)
		ctx.Logger().Error("Recovered from panic", zap.Error(recovered.(error)))
		utils.HandleApiAbortError(c, api_error.InternalServerError)
	})
}
