//go:build !debug

package worker

import (
	"go.uber.org/zap"

	"trimble.com/tdrive/versions-worker/config"
)

func ShutDown(
	_ *config.VersionWorkerConfig,
	_ *zap.Logger) {

	// This is a helper file to simulate shutdown in non-debug mode

}
