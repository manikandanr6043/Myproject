//go:build debug

package worker

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"trimble.com/common/configuration"
	"trimble.com/tdrive/versions-worker/config"
)

func ShutDown(
	versionWorkerConfig *config.VersionWorkerConfig,
	logger *zap.Logger) {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	reader := shutdownKafkaReader(versionWorkerConfig)

	go func() {
		<-sigchan
		closeShutdownKafkaReader(logger, reader)
	}()

	go func() {
		defer closeShutdownKafkaReader(logger, reader)
		for {

			// Read message from kafka broker
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				logger.Error("error while receiving message", zap.Error(err))
				continue
			}
			logger.Fatal("Received shutdown message", zap.Any("message", string(m.Value)))
		}
	}()

}

// closeShutdownKafkaReader function to close kafka reader
func closeShutdownKafkaReader(logger *zap.Logger, shutdownKafkaReader *kafka.Reader) {
	err := shutdownKafkaReader.Close()
	if err != nil {
		logger.Fatal("error while closing the shutdown kafka reader", zap.Error(err))
	}
}

// shutdownKafkaReader Get Shutdown kafka topic reader
func shutdownKafkaReader(config *config.VersionWorkerConfig) *kafka.Reader {
	shutdownKafkaConfig := config.Kafka
	shutdownKafkaConfig.Topic = os.Getenv("GO_CONTROL_TOPIC")
	shutdownKafkaConfig.GroupID = shutdownKafkaConfig.Topic + "-vworker-shutdown"
	readerConfig := configuration.KafkaConfigClient{Config: shutdownKafkaConfig}
	return readerConfig.KafkaReader()
}
