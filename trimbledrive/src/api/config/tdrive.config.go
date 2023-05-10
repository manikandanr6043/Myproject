package config

import "trimble.com/common/configuration"

type TDriveConfig struct {
	Environment    string                     `mapstructure:"environment"`
	Server         ServerConfig               `mapstructure:"server"`
	Mongo          configuration.MongoConfig  `mapstructure:"mongo"`
	StorageAccount StorageAccountConfig       `mapstructure:"storageAccount"`
	Cosmos         configuration.CosmosConfig `mapstructure:"cosmos"`
	TID            TIDConfig                  `mapstructure:"tid"`
	Logger         configuration.LoggerConfig `mapstructure:"logger"`
	Compress       CompressConfig             `mapstructure:"compress"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type StorageAccountConfig struct {
	Name          string `mapstructure:"name"`
	URL           string `mapstructure:"url"`
	ContainerName string `mapstructure:"containerName"`
}

type TIDConfig struct {
	JWKS string `mapstructure:"jwks"`
}

// CompressConfig holds the compress middleware configuration.
type CompressConfig struct {
	// ExcludedContentTypes defines the list of content types to compare the Content-Type header of the incoming requests and responses before compressing.
	// `application/grpc` is always excluded.
	ExcludedContentTypes []string `mapstructure:"excludeContentTypes"`
	// MinResponseBodyBytes defines the minimum amount of bytes a response body must have to be compressed.
	// Default: 1024.
	MinResponseBodyBytes int `mapstructure:"minResponseBodyBytes"`
}
