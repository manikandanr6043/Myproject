package config

import "trimble.com/common/configuration"

type CommitProcessorConfig struct {
	Environment string                     `mapstructure:"environment"`
	MaxRetry    int32                      `mapstructure:"maxRetry"`
	Mongo       configuration.MongoConfig  `mapstructure:"mongo"`
	Kafka       configuration.KafkaConfig  `mapstructure:"kafka"`
	Cosmos      configuration.CosmosConfig `mapstructure:"cosmos"`
	Logger      configuration.LoggerConfig `mapstructure:"logger"`
}
