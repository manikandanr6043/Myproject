package services

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewFilespaceService),
	fx.Provide(NewFolderService),
	fx.Provide(NewBlobService),
	fx.Provide(NewUploadService),
	fx.Provide(NewFileService),
	fx.Provide(NewFolderVersionService),
	fx.Provide(NewFileVersionService),
)
