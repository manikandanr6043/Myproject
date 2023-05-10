package config

import "trimble.com/common/configuration"

type VersionWorkerConfig struct {
	Environment string                     `mapstructure:"environment"`
	Mongo       configuration.MongoConfig  `mapstructure:"mongo"`
	Kafka       configuration.KafkaConfig  `mapstructure:"kafka"`
	Logger      configuration.LoggerConfig `mapstructure:"logger"`
}
