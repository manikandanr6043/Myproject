package middleware

import (
	"go.uber.org/fx"
)

// Module Middleware exported
var Module = fx.Options(
	fx.Provide(NewRequestIdMiddleware),
	fx.Provide(NewAccessLogMiddleware),
	fx.Provide(NewCustomRecoveryMiddleware),
	fx.Provide(NewAuthMiddleware),
	fx.Provide(NewCorsMiddleware),
	fx.Provide(NewCacheControlMiddleware),
	fx.Provide(NewDecompressionMiddleware),
)
