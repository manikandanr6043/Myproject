package configuration

import (
	"crypto/tls"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.uber.org/zap"
)

type KafkaConfig struct {
	Brokers         string `mapstructure:"brokers"`
	ApiKey          string `mapstructure:"apiKey"`
	ApiSecret       string `mapstructure:"apiSecret"`
	Topic           string `mapstructure:"topic"`
	DeadLetterTopic string `mapstructure:"deadLetterTopic"`
	GroupID         string `mapstructure:"groupId"`
	BatchSize       int    `mapstructure:"batchSize"`
}

type KafkaConfigClient struct {
	Config KafkaConfig
}

func (client *KafkaConfigClient) KafkaReader() *kafka.Reader {

	config := client.Config
	brokers := strings.Split(config.Brokers, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: config.GroupID,
		Topic:   config.Topic,
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			TLS:       &tls.Config{},
			SASLMechanism: plain.Mechanism{
				Username: config.ApiKey,
				Password: config.ApiSecret,
			},
		},
	})
	log.Printf("Listening to topic: %s", config.Topic)
	return reader
}

func (client *KafkaConfigClient) KafkaWriter() *kafka.Writer {
	config := client.Config
	brokers := strings.Split(config.Brokers, ",")
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.Murmur2Balancer{},
		Transport: &kafka.Transport{
			TLS: &tls.Config{},
			SASL: plain.Mechanism{
				Username: config.ApiKey,
				Password: config.ApiSecret,
			},
		},
	}
	return writer
}

// CloseReader function to close kafka reader
func CloseReader(reader *kafka.Reader) {
	err := reader.Close()
	if err != nil {
		log.Fatal("Error while closing the reader", zap.Error(err))
	}
}

// CloseWriter function to close kafka writer
func CloseWriter(writer *kafka.Writer) {
	err := writer.Close()
	if err != nil {
		log.Fatal("Error while closing the writer", zap.Error(err))
	}
}
