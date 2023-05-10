package main

import (
	"log"

	"go.uber.org/fx"

	"trimble.com/common/repository"
	"trimble.com/tdrive/versions-worker/config"
	"trimble.com/tdrive/versions-worker/service"
	"trimble.com/tdrive/versions-worker/worker"
)

// Version The version string is expected to be set on the compile/link time by the build workflow.
// The default value here represents a version for binaries built on developer machine (unofficial build).
var Version = "v0.0.0"

func main() {

	log.Printf("Starting application with version : %s", Version)

	app := fx.New(
		config.Module,
		repository.Module,
		service.Module,
		fx.Invoke(
			worker.ShutDown,
			worker.Work,
		),
	)
	app.Run()
}
