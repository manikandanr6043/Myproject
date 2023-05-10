package configuration

import (
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"trimble.com/common/constants"
)

type LoggerConfig struct {
	Level    string          `mapstructure:"level"`
	Sampling *SamplingConfig `mapstructure:"sampling"`
}

type SamplingConfig struct {
	Initial    int `mapstructure:"initial"`
	Thereafter int `mapstructure:"thereafter"`
}

type LoggerConfigClient struct {
	Config      LoggerConfig
	Environment string
}

// NewLogger
// Global logging configuration
func (l *LoggerConfigClient) NewLogger() *zap.Logger {

	config := l.Config

	level := config.Level
	if level == "" {
		// If no level specified, set to default "info" level.
		level = zapcore.InfoLevel.String()
	}

	// parse supplied log level
	zapLoggerLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		// if error on parsing, switch to default level
		log.Println("Error parsing log level ", err)
		zapLoggerLevel = zapcore.InfoLevel
	}

	// Set sampling config
	sampling := &zap.SamplingConfig{
		Initial:    config.Sampling.Initial,
		Thereafter: config.Sampling.Thereafter,
	}

	// Explicitly set nil to ignore sampling if both values are 0
	if sampling.Initial == 0 && sampling.Thereafter == 0 {
		sampling = nil
	}

	isDevelopment := strings.EqualFold(l.Environment, "development")
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLoggerLevel),
		Sampling:         sampling,
		Development:      isDevelopment,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Configure time format
	zapConfig.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.UTC().Format(constants.TimeFormat))
	}

	// Create zap logger
	zapLogger, err := zapConfig.Build()
	if err != nil {
		log.Fatal("Error creating logger ", err)
	}
	return zapLogger

}
