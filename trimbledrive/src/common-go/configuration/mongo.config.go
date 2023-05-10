package configuration

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoConfig struct {
	Uri            string `mapstructure:"uri"`
	DbName         string `mapstructure:"dbname"`
	ConnectTimeOut int64  `mapstructure:"connectTimeOut"`
}
type MongoConfigClient struct {
	Config MongoConfig
}

func (m *MongoConfigClient) NewMongoClient() *mongo.Client {
	config := m.Config
	timeout := time.Duration(config.ConnectTimeOut) * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Uri))
	if err != nil {
		log.Fatal("Error creating mongo client", err)
	}
	log.Println("Mongo Client created")
	return client
}

func (m *MongoConfigClient) NewDatabaseClient(client *mongo.Client) *mongo.Database {
	config := m.Config
	log.Printf("Mongo DB Name : %s", config.DbName)
	return client.Database(config.DbName)
}

// CloseMongoConnection function to disconnect mongo connection
func CloseMongoConnection(mongoClient *mongo.Client) {
	if mongoErr := mongoClient.Disconnect(context.Background()); mongoErr != nil {
		log.Fatal("Error disconnecting mongo client: ", zap.Error(mongoErr))
	}
}
