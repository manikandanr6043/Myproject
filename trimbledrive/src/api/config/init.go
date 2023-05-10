package config

import (
	"go.uber.org/fx"

	"trimble.com/common/configuration"
)

// Build service specific client
var client = configuration.ConfigClient[TDriveConfig]{}
var tDriveConfig = client.NewConfiguration()

func getTDriveConfig() *TDriveConfig {
	return tDriveConfig
}

// Build the required component specific client
var mongoConfigClient = configuration.MongoConfigClient{Config: tDriveConfig.Mongo}
var loggerConfigClient = configuration.LoggerConfigClient{Config: tDriveConfig.Logger, Environment: tDriveConfig.Environment}
var cosmosConfigClient = configuration.CosmosConfigClient{Config: tDriveConfig.Cosmos}

// Module Inject only needed dependencies
var Module = fx.Options(
	fx.Provide(getTDriveConfig),
	fx.Provide(mongoConfigClient.NewMongoClient, mongoConfigClient.NewDatabaseClient),
	fx.Provide(loggerConfigClient.NewLogger),
	fx.Provide(InitializeJWTConfig),
	fx.Provide(cosmosConfigClient.NewAzureCosmosConfig, cosmosConfigClient.NewCosmosDatabaseClient),
	fx.Provide(NewBlobConfig),
)
