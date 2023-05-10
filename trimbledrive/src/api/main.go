package main

import (
	"go.uber.org/fx"

	"trimble.com/common/repository"

	"trimble.com/tdrive/api/config"
	"trimble.com/tdrive/api/handler"
	"trimble.com/tdrive/api/middleware"
	"trimble.com/tdrive/api/router"
	"trimble.com/tdrive/api/server"
	"trimble.com/tdrive/api/services"
)

func main() {
	app := fx.New(
		config.Module,
		repository.Module,
		services.Module,
		middleware.Module,
		router.Module,
		handler.Module,
		fx.Invoke(
			server.Initialize,
		),
	)

	app.Run()
}
