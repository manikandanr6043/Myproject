package middleware

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"trimble.com/common/constants"
	"trimble.com/common/requestcontext"
)

// AccessLogMiddleware -> struct for access log middleware
type AccessLogMiddleware struct {
	logger *zap.Logger
}

// NewAccessLogMiddleware -> new instance of AccessLogMiddleware
func NewAccessLogMiddleware(logger *zap.Logger) AccessLogMiddleware {
	return AccessLogMiddleware{
		logger: logger,
	}
}

// HandleAccessLog -> It sets up the AccessLog middleware
func (m AccessLogMiddleware) HandleAccessLog() gin.HandlerFunc {
	return ginzap.GinzapWithConfig(m.logger, &ginzap.Config{
		Context: func(c *gin.Context) []zapcore.Field {
			var fields []zapcore.Field
			// log protocol
			fields = append(fields, zap.String("protocol", c.Request.Proto))
			// log content length
			fields = append(fields, zap.Int("contentLength", c.Writer.Size()))
			// Get current request context from gin
			ctx := requestcontext.GetContextFromGin(c)
			// log request ID
			fields = append(fields, zap.String(constants.RequestID, ctx.RequestId()))
			// log user ID
			fields = append(fields, zap.String(constants.UserID, ctx.UserId()))
			// log app ID
			fields = append(fields, zap.String(constants.AppID, ctx.AppId()))
			return fields
		},
	})
}
