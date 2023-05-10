package middleware

import (
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"trimble.com/common/api_error"
	"trimble.com/common/requestcontext"

	"trimble.com/tdrive/api/utils"
)

type AuthMiddleware struct {
	JWKS *keyfunc.JWKS
}

// NewAuthMiddleware -> new instance of AuthMiddleware
func NewAuthMiddleware(jwks *keyfunc.JWKS) AuthMiddleware {
	return AuthMiddleware{
		JWKS: jwks,
	}
}

// Authenticate Validate JWT token for its authenticity against TID standards.
// Abort requests with status 401 if the token found invalid.
func (a *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := requestcontext.GetContextFromGin(c)
		authHeader := c.GetHeader("Authorization")
		bearer := strings.Split(authHeader, " ")

		// Behaviour of Go in Middleware functions:
		// If the authentication fails, calling Abort will ensure the remaining handlers
		// for this request are not called, but this method will be executed till the end.
		if len(bearer) != 2 {
			ctx.Logger().Error("Invalid Authorization Header")
			utils.HandleApiAbortError(c, api_error.UnAuthorized)
		} else {
			claims := jwt.MapClaims{}
			// Parse JWT with claims against TID provided JWKS configuration
			_, err := jwt.ParseWithClaims(bearer[1], claims, a.JWKS.Keyfunc)
			if err != nil {
				ctx.Logger().Error("Failed to parse the JWT.", zap.Error(err))
				utils.HandleApiAbortError(c, api_error.UnAuthorized)
			} else {
				if userId, exists := claims["sub"]; exists {
					ctx.SetUserId(userId.(string))
				}
				if appId, exists := claims["azp"]; exists {
					ctx.SetAppId(appId.(string))
				}
			}
		}
		c.Next()
	}

}
