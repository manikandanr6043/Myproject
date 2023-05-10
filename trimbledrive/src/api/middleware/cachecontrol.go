package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const CacheControlHeader = "Cache-Control"

// Code from https://github.com/joeig/gin-cachecontrol/blob/master/cachecontrol.go
// Config defines a cache-control configuration.
//
// References:
// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
type Config struct {
	MustRevalidate       bool
	NoCache              bool
	NoStore              bool
	NoTransform          bool
	Public               bool
	Private              bool
	ProxyRevalidate      bool
	MaxAge               *time.Duration
	SMaxAge              *time.Duration
	Immutable            bool
	StaleWhileRevalidate *time.Duration
	StaleIfError         *time.Duration
}

func (c *Config) buildCacheControl() string {
	var values []string

	if c.MustRevalidate {
		values = append(values, "must-revalidate")
	}

	if c.NoCache {
		values = append(values, "no-cache")
	}

	if c.NoStore {
		values = append(values, "no-store")
	}

	if c.NoTransform {
		values = append(values, "no-transform")
	}

	if c.Public {
		values = append(values, "public")
	}

	if c.Private {
		values = append(values, "private")
	}

	if c.ProxyRevalidate {
		values = append(values, "proxy-revalidate")
	}

	if c.MaxAge != nil {
		values = append(values, fmt.Sprintf("max-age=%.f", c.MaxAge.Seconds()))
	}

	if c.SMaxAge != nil {
		values = append(values, fmt.Sprintf("s-maxage=%.f", c.SMaxAge.Seconds()))
	}

	if c.Immutable {
		values = append(values, "immutable")
	}

	if c.StaleWhileRevalidate != nil {
		values = append(values, fmt.Sprintf("stale-while-revalidate=%.f", c.StaleWhileRevalidate.Seconds()))
	}

	if c.StaleIfError != nil {
		values = append(values, fmt.Sprintf("stale-if-error=%.f", c.StaleIfError.Seconds()))
	}

	return strings.Join(values, ", ")
}

// NoCachePrivatePreset is a cache-control configuration preset which advices the HTTP client not to cache private assets.
var NoCachePrivatePreset = Config{
	NoCache: true,
	Private: true,
}

// CacheAssetsForeverPreset is a cache-control configuration preset which advices the HTTP client
// and all caches in between to cache the object forever without revalidation.
// Technically, "forever" means 1 year, in order to comply with common CDN limits.
var CacheAssetsForeverPreset = Config{
	Public:    true,
	MaxAge:    Duration(8760 * time.Hour),
	Immutable: true,
}

// Duration is a helper function which returns a time.Duration pointer.
func Duration(duration time.Duration) *time.Duration {
	return &duration
}

// CacheControlMiddleware -> struct for requestId
type CacheControlMiddleware struct {
}

// NewCacheControlMiddleware -> new instance of CacheControlMiddleware
func NewCacheControlMiddleware() CacheControlMiddleware {
	return CacheControlMiddleware{}
}

// HandleCacheControl -> preflight response
func (m CacheControlMiddleware) HandleCacheControl() gin.HandlerFunc {
	value := NoCachePrivatePreset.buildCacheControl()

	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			header := c.Writer.Header()
			header.Set(CacheControlHeader, value)
		}

		c.Next()
	}
}
