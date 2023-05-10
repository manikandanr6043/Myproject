package worker

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
	"trimble.com/common/constants"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/versions-worker/config"
	"trimble.com/tdrive/versions-worker/model"
	"trimble.com/tdrive/versions-worker/service"
)

type Worker struct {
	config      *config.VersionWorkerConfig
	mongoClient *mongo.Client
	reader      *kafka.Reader
	writer      *kafka.Writer
	logger      *zap.Logger
	service     *service.VersionProcessorService
}

func Work(config *config.VersionWorkerConfig,
	mongoClient *mongo.Client,
	reader *kafka.Reader,
	writer *kafka.Writer,
	logger *zap.Logger,
	service *service.VersionProcessorService) {

	// Create worker instance on invoke
	w := &Worker{
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
	go w.process(ctx, wg)

	<-sigchan

	// Graceful shutdown post termination signals
	cancel()
	wg.Wait()
	logger.Info("Shutting down the application")

	// Below is to trigger Fx App shutdown
	os.Exit(1)
}

func (w *Worker) process(ctx context.Context, wg *sync.WaitGroup) {

	// Close wait group
	defer wg.Done()

	// Close reader
	defer configuration.CloseReader(w.reader)

	// Close writer
	defer configuration.CloseWriter(w.writer)

	// Close mongo connection
	defer configuration.CloseMongoConnection(w.mongoClient)

	rCtx := requestcontext.NoOpRequestContext(w.logger)

	var events []model.Event
	seenResources := map[string]interface{}{}
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Stopping processing of messages.")
			return
		default:
			// Read message from kafka broker
			m, err := w.reader.FetchMessage(ctx)
			if err != nil {
				w.logger.Error("error while receiving message", zap.Error(err))
				continue
			}

			// Check if the message is intended for this consumer ("versions-worker")
			// If not, message should be ignored without committing the offset
			// so that it gets processed by any other intended consumer
			if isNotIntendedConsumer(m) {
				w.logger.Warn("Message is not intended for this consumer",
					zap.String("consumerName", *getHeader(constants.ConsumerNameHeader, m)))
				continue
			}

			retryCount := getRetryCount(m)
			w.logger.Debug("Received message",
				zap.Any("topic", m.Topic),
				zap.Any("partition", m.Partition),
				zap.Any("offset", m.Offset),
				zap.Any("key", string(m.Key)),
				zap.Any("message", string(m.Value)),
				zap.Any("retryCount", retryCount),
			)

			// Parse the incoming event
			var event model.Event
			unmarshalErr := json.Unmarshal(m.Value, &event)
			if unmarshalErr != nil {
				w.logger.Error("error in parsing event, publishing to DLT", zap.Error(unmarshalErr))
				// If unmarshalling of message itself fails, push it to DLT.
				failedMessage := kafka.Message{Topic: w.config.Kafka.DeadLetterTopic, Value: m.Value, Headers: generateReprocessMessageHeaders(retryCount)}
				w.writeMessages(failedMessage)
				continue
			}

			event.RetryCount = retryCount
			event.Message = m

			// Check if previous version for the same resource exists already in the batch.
			// If present, flush the batch to DB,then process the incoming event.
			if _, ok := seenResources[event.DocumentKey.Id]; ok {

				// Flush the batch of events to database
				w.flushBatch(rCtx, events)

				// Neutralize events list
				events = nil

				// Reset previous version set
				seenResources = map[string]interface{}{}

			}

			// generate full version document by applying update descriptors to the previous version
			// used in update cases
			spaceId := string(m.Key)
			generateFullDocumentErr := w.service.GenerateFullVersionDocument(rCtx, spaceId, &event)
			if generateFullDocumentErr != nil {
				w.logger.Error("error in generating actual event, publishing to DLT", zap.Any("failedEvent", event))
				// Publish to DLT.
				failedMessage := kafka.Message{Topic: w.config.Kafka.DeadLetterTopic, Value: m.Value, Headers: generateReprocessMessageHeaders(retryCount)}
				w.writeMessages(failedMessage)
				continue
			}

			events = append(events, event)
			seenResources[event.DocumentKey.Id] = struct{}{}

			// Start processing and pushing the data to database.
			// Messages are pushed to database on any one of two requirements met
			// 1. When the internal buffer size meets configured batch size.
			//    This case is very likely to happen when the topic is populous with messages (High traffic).
			// 2. When there is no more message in the pipeline.
			//    Meaning, current processing message is the last one in the pipeline.
			//    This helps in pushing the data to database during slow moving traffic,
			//    without expecting to meet batchSize limit, thereby reducing the wait time for replication.
			//    During slow moving traffic, data will be pushed as and when it reaches the consumer.
			//    Behaves like when batchSize=1.
			stats := w.reader.Stats()
			if len(events)%w.config.Kafka.BatchSize == 0 || stats.QueueLength == 0 {

				// Flush the batch of events to database
				w.flushBatch(rCtx, events)

				// Neutralize events list
				events = nil

				// Reset previous version set
				seenResources = map[string]interface{}{}

			}
		}

	}
}

// flushBatch Flush the data accumulated so far to database.
// Handle failures if any (retry mechanism).
func (w *Worker) flushBatch(rCtx *requestcontext.RequestContext, events []model.Event) {

	//Push the change documents to database
	failedEvents := w.service.PushToDatabase(rCtx, events)
	w.logger.Debug("Pushed to database", zap.Any("failedEvents", failedEvents))

	// original message reference of last item pushed to db
	lastProcessedMessage := events[len(events)-1].Message

	// Commit messages
	w.commitMessages(lastProcessedMessage)

	// If there are failed events,
	// publish those back to same topic with additional metadata like consumerName, retryCount.
	w.handleFailures(failedEvents...)

}

// handleFailures Handle failed events
// Publish the failed events as a message back to the same topic,
// with metadata like consumerName, incremented retryCount
func (w *Worker) handleFailures(failedEvents ...model.Event) {

	if failedEvents != nil {
		var failedMessages []kafka.Message
		for _, e := range failedEvents {
			w.logger.Debug("Publishing failed event to same topic", zap.Any("id", e.FullDocument["id"]),
				zap.String("version", repository.FormatVersion(e.FullDocument["majorVersion"].(int64), e.FullDocument["minorVersion"].(int64))),
				zap.Int64("retryCount", e.RetryCount))
			if failedKafkaMessage, err := w.convertEventToKafkaMessage(w.config.Kafka.Topic, "", e); err == nil {
				failedMessages = append(failedMessages, *failedKafkaMessage)
			}
		}
		w.writeMessages(failedMessages...)
	}

}

// convertEventToKafkaMessage Convert event struct to kafka message with given key, topic and payload
func (w *Worker) convertEventToKafkaMessage(topic string, key string, event model.Event) (*kafka.Message, error) {
	if marshal, err := json.Marshal(&event); err != nil {
		w.logger.Warn("error on marshalling event", zap.Error(err))
		return nil, err
	} else {
		return &kafka.Message{Topic: topic, Key: []byte(key), Value: marshal, Headers: generateReprocessMessageHeaders(event.RetryCount)}, nil
	}
}

// commitMessages Commit messages to kafka brokers
func (w *Worker) commitMessages(m ...kafka.Message) {
	if commitErr := w.reader.CommitMessages(context.TODO(), m...); commitErr != nil {
		w.logger.Fatal("failed to commit messages:", zap.Error(commitErr))
	}
}

// writeMessages Push messages to configured topic
func (w *Worker) writeMessages(messages ...kafka.Message) {

	err := w.writer.WriteMessages(context.TODO(), messages...)
	if err != nil {
		w.logger.Warn("Message write failed", zap.Any("failedWriteMessages", messages), zap.Error(err))
	}
}
