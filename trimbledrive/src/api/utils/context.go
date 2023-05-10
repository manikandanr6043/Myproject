package utils

import "net/http"

// GetHost returns the service hostname as visible by the consumer.
// Detects if there is e.g. CDN like FrontDoor is present.
func GetHost(req *http.Request) string {
	host := req.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = req.Host
	}
	return host
}
