package processor

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"trimble.com/common/configuration"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/commit-processor/config"
	"trimble.com/tdrive/commit-processor/model"
	"trimble.com/tdrive/commit-processor/service"
)

type Processor struct {
	config      *config.CommitProcessorConfig
	mongoClient *mongo.Client
	reader      *kafka.Reader
	writer      *kafka.Writer
	logger      *zap.Logger
	service     *service.CommitProcessorService
}

func CommitProcessor(config *config.CommitProcessorConfig,
	mongoClient *mongo.Client,
	reader *kafka.Reader,
	writer *kafka.Writer,
	logger *zap.Logger,
	service *service.CommitProcessorService) {

	// Create processor instance on invoke
	p := &Processor{
		config:      config,
		mongoClient: mongoClient,
		reader:      reader,
		writer:      writer,
		logger:      logger,
		service:     service,
	}

	// Wait group to make sure the needed process are gracefully stopped
	wg := new(sync.WaitGroup)

	// System notification signal setup for graceful shutdown
	// Context to trigger Done channel on termination
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go p.processUploads(ctx, wg)

	<-sigchan

	// Graceful shutdown post termination signals
	cancel()
	wg.Wait()
	logger.Info("Shutting down the application")

	// Below is to trigger Fx App shutdown
	os.Exit(1)
}

func (p *Processor) processUploads(ctx context.Context, wg *sync.WaitGroup) {

	// Close wait group
	defer wg.Done()

	// Close reader
	defer configuration.CloseReader(p.reader)

	// Close writer
	defer configuration.CloseWriter(p.writer)

	// Close mongo connection
	defer configuration.CloseMongoConnection(p.mongoClient)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Stopping processing of messages.")
			return
		default:
			// Read message from kafka broker
			m, err := p.reader.FetchMessage(ctx)
			if err != nil {
				p.logger.Error("Error while receiving message", zap.Error(err))
				continue
			}
			p.logger.Debug("Received message",
				zap.Any("topic", m.Topic),
				zap.Any("partition", m.Partition),
				zap.Any("offset", m.Offset),
				zap.Any("message", string(m.Value)))

			// Parse the incoming event
			var commitMessage model.CommitMessage
			unmarshalErr := json.Unmarshal(m.Value, &commitMessage)
			if unmarshalErr != nil {
				p.logger.Error("Error in parsing event", zap.Error(unmarshalErr))
				p.commitMessages(p.logger, m)
				continue
			}
			rCtx := requestcontext.NewRequestContext(commitMessage.Data.UploadId, p.logger)

			ok := p.service.CommitFileUpload(rCtx, commitMessage)
			rCtx.Logger().Debug("CommitFileUpload Completed", zap.Bool("commitStatus", ok))
			if !ok {
				// Log and commit if retry max count reached
				if commitMessage.RetryCount > p.config.MaxRetry {
					rCtx.Logger().Error("Reached maximum number of retries, hence skipping the event from processing", zap.Any("event", commitMessage))
				} else {
					// Increase retry count and post to same topic
					commitMessage.RetryCount = commitMessage.RetryCount + 1
					failedKafkaMessage, _ := convertEventToKafkaMessage(rCtx.Logger(), m.Topic, commitMessage)
					p.writeMessages(*failedKafkaMessage)
				}
			}
			// Commit messages
			p.commitMessages(rCtx.Logger(), m)
		}
	}
}

// convertEventToKafkaMessage Convert event struct to kafka message with given topic and payload
func convertEventToKafkaMessage(logger *zap.Logger, topic string, event model.CommitMessage) (*kafka.Message, error) {
	if marshal, err := json.Marshal(event); err != nil {
		logger.Warn("error on marshalling event", zap.Error(err))
		return nil, err
	} else {
		return &kafka.Message{Topic: topic, Value: marshal}, nil
	}
}

// commitMessages Commit message to kafka brokers
func (p *Processor) commitMessages(logger *zap.Logger, msg ...kafka.Message) {
	if commitErr := p.reader.CommitMessages(context.TODO(), msg...); commitErr != nil {
		logger.Fatal("Failed to commit message:", zap.Error(commitErr))
	}
}

// writeMessages Push messages to configured Topic
func (p *Processor) writeMessages(msg ...kafka.Message) {
	err := p.writer.WriteMessages(context.TODO(), msg...)
	if err != nil {
		p.logger.Warn("Message write failed", zap.Any("failedMessage", msg), zap.Error(err))
	}
}
