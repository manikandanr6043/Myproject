package middleware

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"trimble.com/common/constants"
	"trimble.com/common/requestcontext"
)

// DecompressionMiddleware -> struct for requestId
type DecompressionMiddleware struct {
}

// NewDecompressionMiddleware -> new instance of CorsMiddleware
func NewDecompressionMiddleware() DecompressionMiddleware {
	return DecompressionMiddleware{}
}

// HandleDecompression -> decompress encoded requests
func (m DecompressionMiddleware) HandleDecompression() gin.HandlerFunc {
	return func(c *gin.Context) {
		encoding := c.Request.Header.Get(constants.ContentEncoding)
		if encoding != "" {
			ctx := requestcontext.GetContextFromGin(c)
			ctx.Logger().Info(constants.ContentEncoding + ": " + encoding)
			switch encoding {
			case "br":
				DecompressBrotli(c)
			case "gzip":
				DecompressGzip(c)
			}
		}

		c.Next()
	}
}

func DecompressBrotli(c *gin.Context) {
	if c.Request.Body == nil {
		return
	}
	r := brotli.NewReader(c.Request.Body)
	c.Request.Header.Del(constants.ContentEncoding)
	c.Request.Header.Del(constants.ContentLength)
	c.Request.Body = io.NopCloser(r)
}

func DecompressGzip(c *gin.Context) {
	if c.Request.Body == nil {
		return
	}
	r, err := gzip.NewReader(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Request.Header.Del(constants.ContentEncoding)
	c.Request.Header.Del(constants.ContentLength)
	c.Request.Body = r
}
