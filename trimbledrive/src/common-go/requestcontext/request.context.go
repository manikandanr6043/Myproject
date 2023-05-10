package requestcontext

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"trimble.com/common/constants"
)

// RequestContext => struct for maintaining contextual information for a request
type RequestContext struct {
	userId      string
	appId       string
	requestId   string
	dbCtx       context.Context
	logger      *zap.Logger
	httpRequest *http.Request
}

// NewRequestContext Creates a new request context for a request
// Holds contextual information pertaining to a specific request
func NewRequestContext(requestId string, logger *zap.Logger) *RequestContext {
	ctx := &RequestContext{requestId: requestId}

	// Create contextual sub logger per request using global logger.
	ctxLogger := logger.With(
		zap.String(constants.RequestID, ctx.RequestId()),
	)
	ctx.setLogger(ctxLogger)
	return ctx
}

// NoOpRequestContext Creates a no operation request context
// Can be used in modules which doesn't need contextual/child log capabilities
// need only to log and interact with db layer
func NoOpRequestContext(logger *zap.Logger) *RequestContext {
	ctx := &RequestContext{logger: logger}
	return ctx
}

// GetContextFromGin Fetch current request context from incoming gin context
func GetContextFromGin(c *gin.Context) *RequestContext {
	ctx, _ := c.Get(constants.RequestCtx)
	return ctx.(*RequestContext)
}

// ------------- Getters and Setters -------------

func (r *RequestContext) UserId() string {
	return r.userId
}

func (r *RequestContext) AppId() string {
	return r.appId
}

func (r *RequestContext) RequestId() string {
	return r.requestId
}

func (r *RequestContext) Logger() *zap.Logger {
	return r.logger
}

func (r *RequestContext) setLogger(logger *zap.Logger) {
	r.logger = logger
}

func (r *RequestContext) SetUserId(userId string) {
	r.userId = userId
}

func (r *RequestContext) SetAppId(appId string) {
	r.appId = appId
}

func (r *RequestContext) SetRequestId(requestId string) {
	r.requestId = requestId
}

func (r *RequestContext) SetLogger(logger *zap.Logger) {
	r.logger = logger
}
func (r *RequestContext) DbCtx() context.Context {
	if r.dbCtx == nil {
		return context.TODO()
	}
	return r.dbCtx
}

func (r *RequestContext) SetDbCtx(dbCtx context.Context) {
	r.dbCtx = dbCtx
}

func (r *RequestContext) SetHttpRequest(httpRequest *http.Request) {
	r.httpRequest = httpRequest
}

func (r *RequestContext) HttpRequest() *http.Request {
	return r.httpRequest
}
