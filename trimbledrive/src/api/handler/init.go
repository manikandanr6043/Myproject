package handler

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewAppInfoHandler),
	fx.Provide(NewFileSpaceHandler),
	fx.Provide(NewFolderHandler),
	fx.Provide(NewUploadsHandler),
	fx.Provide(NewFileHandler),
	fx.Provide(NewFolderVersionHandler),
	fx.Provide(NewFileVersionHandler),
)
