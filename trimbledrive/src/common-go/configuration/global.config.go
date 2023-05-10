package configuration

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ConfigClient[T any] struct {
}

// NewConfiguration is an exported method starts the viper
// (external lib) and sets the data read onto configuration struct.
func (c ConfigClient[T]) NewConfiguration() *T {
	var globalConfig *T
	viper := viper.NewWithOptions(viper.KeyDelimiter("_"))

	// We configure the app behavior using an environment variable GO_ENVIRONMENT.
	// Values allowed: "production" (default), "staging", "development"
	//
	// Order how configuration values are collected:
	// 1. Configuration is first read from `config.yaml`.
	// 2. Then depending on the GO_ENVIRONMENT value the settings will be read from `config.{environment}.yaml` file. Values from this env specific file it exists will override already read values.
	// 3. The last step is environment variables are read and if found they will override the values read from the configuration files. The `GO_` prefix is used and "_" and a key separator. Name must be all uppercase.
	//
	// Different configuration files are bundled in to the application distribution so it is possible to switch between execution modes
	// by changing the GO_ENVIRONMENT variable value and restarting the application.
	// Some common settings that differ in production from development include:
	// - Caching.
	// - Client-side resources are bundled, minified, and potentially served from a CDN.
	// - Diagnostic error pages disabled.
	// - Friendly error pages enabled.
	// - Production logging and monitoring enabled. For example, using Application Insights.

	appEnvironment := os.Getenv("GO_ENVIRONMENT")
	if appEnvironment == "" {
		appEnvironment = "production"
		os.Setenv("GO_ENVIRONMENT", appEnvironment)
		log.Printf("The GO_ENVIRONMENT environment variable is missing, default to `%s`.", appEnvironment)
	}

	log.Printf("GO_ENVIRONMENT=%s", appEnvironment)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO")
	viper.SetConfigType("yaml")

	viper.SetDefault("environment", "production")
	viper.SetDefault("server_port", 8080)
	viper.SetDefault("tid_jwks", "https://id.trimble.com/.well-known/jwks.json")
	viper.SetDefault("mongo_uri", "mongodb://localhost:27017")
	viper.SetDefault("mongo_dbname", "TrimbleDrive")

	// The configuration file will be searched in the "current" folder and in the folder where the executable is
	viper.AddConfigPath(".")
	if executableFilePath, err := os.Executable(); err == nil {
		exPath := filepath.Dir(executableFilePath)
		viper.AddConfigPath(exPath)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})

	// first read from generic config
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Couldn't load default configuration: %s", err)
	} else {
		viper.WatchConfig()
	}

	// second read from environment specific config
	viper.SetConfigFile("config." + appEnvironment + ".yaml")
	if err := viper.MergeInConfig(); err != nil {
		log.Printf("Couldn't load environment configuration: %s", err)
	} else {
		viper.WatchConfig()
	}

	if err := viper.Unmarshal(&globalConfig); err != nil {
		log.Fatalf("Marshal error on parsing configuration file: %s", err)
	}

	return globalConfig
}
