package configuration

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// CosmosConfig struct for cosmos configuration properties
type CosmosConfig struct {
	Url    string `mapstructure:"url"`
	DbName string `mapstructure:"dbName"`
}

type CosmosConfigClient struct {
	Config CosmosConfig
}

// NewAzureCosmosConfig creates azure cosmos client
func (c *CosmosConfigClient) NewAzureCosmosConfig() *azcosmos.Client {
	// Authenticate using Azure AD
	cred, identityErr := azidentity.NewDefaultAzureCredential(nil)
	if identityErr != nil {
		log.Fatal("Error fetching azure credentials", identityErr)
	}
	// Create Cosmos DB client using Azure AD cred
	client, err := azcosmos.NewClient(c.Config.Url, cred, nil)
	if err != nil {
		log.Fatal("Error creating cosmos db client", err)
	}
	log.Println("Cosmos db client created!")
	return client
}

// NewCosmosDatabaseClient creates azure cosmos database client
func (c *CosmosConfigClient) NewCosmosDatabaseClient(cosmosClient *azcosmos.Client) *azcosmos.DatabaseClient {
	databaseClient, err := cosmosClient.NewDatabase(c.Config.DbName)
	if err != nil {
		log.Fatalf("Error creating cosmos database client: %s", err)
	}
	return databaseClient
}
