package repository

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewLatestRepository),
	fx.Provide(NewResourceVersionRepository),
	fx.Provide(NewFilespaceRepository),
	fx.Provide(NewFileUploadRepository),
)
