package config

import (
	"go.uber.org/fx"

	"trimble.com/common/configuration"
)

// Build service specific client
var client = configuration.ConfigClient[CommitProcessorConfig]{}
var commitProcessorConfig = client.NewConfiguration()

func getCommitProcessorConfig() *CommitProcessorConfig {
	return commitProcessorConfig
}

// Build the required component specific client
var kafkaConfigClient = configuration.KafkaConfigClient{Config: commitProcessorConfig.Kafka}
var mongoConfigClient = configuration.MongoConfigClient{Config: commitProcessorConfig.Mongo}
var cosmosConfigClient = configuration.CosmosConfigClient{Config: commitProcessorConfig.Cosmos}
var loggerConfigClient = configuration.LoggerConfigClient{Config: commitProcessorConfig.Logger, Environment: commitProcessorConfig.Environment}

// Module Inject only needed dependencies
var Module = fx.Options(
	fx.Provide(getCommitProcessorConfig),
	fx.Provide(kafkaConfigClient.KafkaReader),
	fx.Provide(kafkaConfigClient.KafkaWriter),
	fx.Provide(mongoConfigClient.NewMongoClient, mongoConfigClient.NewDatabaseClient),
	fx.Provide(cosmosConfigClient.NewAzureCosmosConfig, cosmosConfigClient.NewCosmosDatabaseClient),
	fx.Provide(loggerConfigClient.NewLogger),
)
