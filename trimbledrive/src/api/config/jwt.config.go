package config

import (
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
)

// InitializeJWTConfig An exported method to initialize jwks resource with desired config
func InitializeJWTConfig(tDriveConfig *TDriveConfig) *keyfunc.JWKS {

	// Create the keyfunc options. Use an error handler that logs. Refresh the JWKS when a JWT signed by an unknown KID
	// is found or at the specified interval. Timeout the initial JWKS refresh request after
	// 10 seconds. This timeout is also used to create the initial context. Context for keyfunc.Get.
	keyFuncOptions := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			log.Fatalf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(tDriveConfig.TID.JWKS, keyFuncOptions)
	if err != nil {
		log.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
	}
	return jwks
}
