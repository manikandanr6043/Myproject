package config

import (
	"go.uber.org/fx"

	"trimble.com/common/configuration"
)

// Build service specific client
var client = configuration.ConfigClient[VersionWorkerConfig]{}
var versionWorkerConfig = client.NewConfiguration()

func getVersionWorkerConfig() *VersionWorkerConfig {
	return versionWorkerConfig
}

// Build the required component specific client
var kafkaConfigClient = configuration.KafkaConfigClient{Config: versionWorkerConfig.Kafka}
var mongoConfigClient = configuration.MongoConfigClient{Config: versionWorkerConfig.Mongo}
var loggerConfigClient = configuration.LoggerConfigClient{Config: versionWorkerConfig.Logger, Environment: versionWorkerConfig.Environment}

// Module Inject only needed dependencies
var Module = fx.Options(
	fx.Provide(getVersionWorkerConfig),
	fx.Provide(kafkaConfigClient.KafkaReader),
	fx.Provide(kafkaConfigClient.KafkaWriter),
	fx.Provide(mongoConfigClient.NewMongoClient, mongoConfigClient.NewDatabaseClient),
	fx.Provide(loggerConfigClient.NewLogger),
)
