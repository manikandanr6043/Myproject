package model

import "github.com/segmentio/kafka-go"

// Event -> defines event structure
type Event struct {
	DocumentKey struct {
		Id string `json:"_id"`
	}

	// FullDocument struct for mapping "latest" entity to "versions" entity coming in events
	FullDocument map[string]interface{} `json:"fullDocument"`

	// OperationType Field specifying type of operation (Ex: insert, update)
	OperationType string `json:"operationType"`

	// UpdateDescription Field specifying the details of update operation
	UpdateDescription struct {
		RemovedFields   []string               `json:"removedFields"`
		TruncatedArrays []string               `json:"truncatedArrays"`
		UpdatedFields   map[string]interface{} `json:"updatedFields"`
	} `json:"updateDescription"`

	// RetryCount Internal usage to support retry mechanism.
	RetryCount int64 `json:"-"`

	// Message Original message of an event for internal usage
	// A reference to the original message in the events data structure
	// Provides more flexibility in committing messages from the last event.
	Message kafka.Message `json:"-"`
}
