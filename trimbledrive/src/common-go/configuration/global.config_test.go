package configuration

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Environment string           `mapstructure:"environment"`
	Server      TestServerConfig `mapstructure:"server"`
	TID         TestTIDConfig    `mapstructure:"tid"`
	Mongo       TestMongoConfig  `mapstructure:"mongo"`
}

type TestServerConfig struct {
	Port int `mapstructure:"port"`
}
type TestMongoConfig struct {
	Uri    string `mapstructure:"uri"`
	DbName string `mapstructure:"dbname"`
}

type TestTIDConfig struct {
	JWKS string `mapstructure:"jwks"`
}

func TestDefaultValuesAreSet(t *testing.T) {

	expectedEnv := "production"
	expectedPort := 8080
	expectedJwks := "https://id.trimble.com/.well-known/jwks.json"
	expectedMongoUri := "mongodb://localhost:27017"
	expectedMongoDbName := "TrimbleDrive"

	// make sure there are no env variables
	t.Setenv("GO_ENVIRONMENT", "")
	t.Setenv("GO_SERVER_PORT", "")
	t.Setenv("GO_TID_JWKS", "")
	t.Setenv("GO_MONGO_URI", "")
	t.Setenv("GO_MONGO_DBNAME", "")

	client := ConfigClient[TestConfig]{}
	config := client.NewConfiguration()

	assert.Equal(t, expectedEnv, config.Environment, "environment")
	assert.Equal(t, expectedPort, config.Server.Port, "port")
	assert.Equal(t, expectedJwks, config.TID.JWKS, "jwks")
	assert.Equal(t, expectedMongoUri, config.Mongo.Uri, "Mongo.Uri")
	assert.Equal(t, expectedMongoDbName, config.Mongo.DbName, "Mongo.DbName")
}

func TestValuesArePickedUpFromVariables(t *testing.T) {
	expectedEnv := "development"
	expectedPort := 123
	expectedJwks := "https://localhost/jwks.json"
	expectedMongoUri := "mongodb://remote:27017"
	expectedMongoDbName := "something"

	// set configuration via env variables
	t.Setenv("GO_ENVIRONMENT", expectedEnv)
	t.Setenv("GO_SERVER_PORT", strconv.Itoa(expectedPort))
	t.Setenv("GO_TID_JWKS", expectedJwks)
	t.Setenv("GO_MONGO_URI", expectedMongoUri)
	t.Setenv("GO_MONGO_DBNAME", expectedMongoDbName)

	client := ConfigClient[TestConfig]{}
	config := client.NewConfiguration()

	assert.Equal(t, expectedEnv, config.Environment, "environment")
	assert.Equal(t, expectedPort, config.Server.Port, "port")
	assert.Equal(t, expectedJwks, config.TID.JWKS, "jwks")
	assert.Equal(t, expectedMongoUri, config.Mongo.Uri, "Mongo.Uri")
	assert.Equal(t, expectedMongoDbName, config.Mongo.DbName, "Mongo.DbName")
}

func TestValuesArePickedUpFromConfigFile(t *testing.T) {
	expectedPort := 123
	t.Setenv("GO_SERVER.PORT", "")

	envConfigContent := []byte("---\nserver:\n  port: " + strconv.Itoa(expectedPort))
	if err := os.WriteFile("./config.yaml", envConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write config file: %s", err)
	}
	defer os.Remove("./config.yaml")

	client := ConfigClient[TestConfig]{}
	config := client.NewConfiguration()

	assert.Equal(t, expectedPort, config.Server.Port, "port")
}

func TestValuesArePickedUpFromEnvConfigFile(t *testing.T) {
	expectedEnv := "development"
	expectedPort := 124
	t.Setenv("GO_ENVIRONMENT", expectedEnv)
	t.Setenv("GO_SERVER_PORT", "")

	envConfigContent := []byte("---\nserver:\n  port: " + strconv.Itoa(expectedPort))
	if err := os.WriteFile("./config."+expectedEnv+".yaml", envConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write config file: %s", err)
	}
	defer os.Remove("./config." + expectedEnv + ".yaml")

	client := ConfigClient[TestConfig]{}
	config := client.NewConfiguration()

	assert.Equal(t, expectedPort, config.Server.Port, "port")
}

func TestEnvConfigFileOverridesGenericConfigFile(t *testing.T) {
	expectedEnv := "development"
	expectedPort := 124
	t.Setenv("GO_ENVIRONMENT", expectedEnv)
	t.Setenv("GO_SERVER_PORT", "")

	genericConfigContent := []byte("---\nserver:\n  port: " + strconv.Itoa(expectedPort+3))
	if err := os.WriteFile("./config.yaml", genericConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write config file: %s", err)
	}
	defer os.Remove("./config.yaml")
	envConfigContent := []byte("---\nserver:\n  port: " + strconv.Itoa(expectedPort))
	if err := os.WriteFile("./config."+expectedEnv+".yaml", envConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write config file: %s", err)
	}
	defer os.Remove("./config." + expectedEnv + ".yaml")

	client := ConfigClient[TestConfig]{}
	config := client.NewConfiguration()

	assert.Equal(t, expectedPort, config.Server.Port, "port")
}

func TestEnvVariablesOverrideConfigFiles(t *testing.T) {
	expectedEnv := "development"
	expectedPort := 124
	t.Setenv("GO_ENVIRONMENT", expectedEnv)
	t.Setenv("GO_SERVER_PORT", strconv.Itoa(expectedPort))

	genericConfigContent := []byte("---\nserver:\n  port: " + strconv.Itoa(expectedPort+3))
	if err := os.WriteFile("./config.yaml", genericConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write config file: %s", err)
	}
	defer os.Remove("./config.yaml")
	envConfigContent := []byte("---\nserver:\n  port: " + strconv.Itoa(expectedPort+4))
	if err := os.WriteFile("./config."+expectedEnv+".yaml", envConfigContent, 0644); err != nil {
		t.Fatalf("Failed to write config file: %s", err)
	}
	defer os.Remove("./config." + expectedEnv + ".yaml")

	client := ConfigClient[TestConfig]{}
	config := client.NewConfiguration()

	assert.Equal(t, expectedPort, config.Server.Port, "port")
}
