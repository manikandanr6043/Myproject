package worker

// Utility functionalities needed for worker
// Needed functionalities can be later extracted to common module.
import (
	"strconv"

	"github.com/segmentio/kafka-go"

	"trimble.com/common/constants"
)

// isNotIntendedConsumer Check consumerName of the messages
// Return true if it is not intended for this consumer or false if it's intended.
func isNotIntendedConsumer(m kafka.Message) bool {
	consumerName := getHeader(constants.ConsumerNameHeader, m)
	if consumerName != nil && *consumerName != constants.VersionWorkerConsumerName {
		return true
	}
	return false
}

// getRetryCount Get retry count from kafka message headers
func getRetryCount(m kafka.Message) int64 {
	retryCount := getHeader(constants.RetryCountHeader, m)
	if retryCount != nil {
		if value, err := strconv.ParseInt(*retryCount, 10, 64); err == nil {
			return value
		}
	}
	return 0
}

// generateReprocessMessageHeaders Generate kafka message headers with consumerName, retryCount.
func generateReprocessMessageHeaders(retryCount int64) []kafka.Header {
	count := retryCount + 1
	return []kafka.Header{
		{
			Key:   constants.ConsumerNameHeader,
			Value: []byte(constants.VersionWorkerConsumerName),
		}, {
			Key:   constants.RetryCountHeader,
			Value: []byte(strconv.FormatInt(count, 10)),
		},
	}
}

// getHeader Get value of the specified header
func getHeader(header string, m kafka.Message) *string {
	for _, h := range m.Headers {
		if h.Key == header {
			value := string(h.Value)
			return &value
		}
	}
	return nil
}
