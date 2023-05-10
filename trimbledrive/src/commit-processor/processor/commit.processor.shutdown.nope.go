//go:build !debug

package processor

import (
	"go.uber.org/zap"

	"trimble.com/tdrive/commit-processor/config"
)

func ShutDown(
	_ *config.CommitProcessorConfig,
	_ *zap.Logger) {

	// This is a helper file to simulate shutdown in non-debug mode

}
