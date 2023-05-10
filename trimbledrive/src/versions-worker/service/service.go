package service

import (
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"

	"trimble.com/common/constants"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/versions-worker/model"
)

// VersionProcessorService -> struct for versions processor service
type VersionProcessorService struct {
	repository *repository.ResourceVersionRepository
}

// NewVersionProcessorService creates new versions processor instance
func NewVersionProcessorService(repository *repository.ResourceVersionRepository) *VersionProcessorService {
	return &VersionProcessorService{
		repository: repository,
	}
}

// GenerateFullVersionDocument Create full version document by applying update descriptors to the previous version
func (m *VersionProcessorService) GenerateFullVersionDocument(ctx *requestcontext.RequestContext, spaceId string, event *model.Event) error {

	// If operationType is update and the document is not already merged,
	// Perform merge of current version with previous version
	if event.OperationType == "update" && event.FullDocument == nil {
		updatedFields := event.UpdateDescription.UpdatedFields
		previousVersionDbId, err := getPreviousVersionDbIdFromEvent(event)
		if err != nil {
			ctx.Logger().Warn("Error in generating previous version db identifier from event", zap.Error(err))
			return err
		}
		var versionToBeUpdated map[string]interface{}
		// Fetch version to be patched
		err = m.repository.GetResourceVersionByDbId(ctx, spaceId, *previousVersionDbId, &versionToBeUpdated)
		if err != nil {
			ctx.Logger().Warn("Error in fetching previous version", zap.Error(err))
			return err
		}
		// Modify only updated fields
		for k, v := range updatedFields {
			versionToBeUpdated[k] = v
		}
		event.FullDocument = versionToBeUpdated
	}
	return nil
}

// PushToDatabase process events and push change documents as versions to database
// Return failed events if any
func (m *VersionProcessorService) PushToDatabase(ctx *requestcontext.RequestContext, events []model.Event) []model.Event {

	// Formulate versions entity from the batch of events
	var versions []interface{}
	for _, event := range events {
		e := event.FullDocument

		// Few agreed fields manipulation that the system need to be aware of.

		// Generate resource version identifier "_id"
		if val, ok := e["id"]; ok {
			// Performed in patch cases where "_id" already exists and just needs to overridden with updated version
			e["_id"] = repository.GetDbVersionIdFromDbId(val.(string), int64(e["majorVersion"].(float64)), int64(e["minorVersion"].(float64)))
		} else {
			// Performed in new insert cases where "_id" have to formed using resourceId
			// Here resourceId is "_id" of latest document in the change document.
			// Hence, update the "_id" with resource version identifier
			// and add a field "id" to represent resource identifier.
			resourceId := e["_id"].(string)
			e["_id"] = repository.GetDbVersionIdFromDbId(resourceId, int64(e["majorVersion"].(float64)), int64(e["minorVersion"].(float64)))
			e["id"] = resourceId
		}

		// Convert to valid time fields if the incoming type is float64
		if val, ok := e["created"].(float64); ok {
			e["created"] = time.UnixMilli(int64(val))
		}
		if val, ok := e["modified"].(float64); ok {
			e["modified"] = time.UnixMilli(int64(val))
		}

		versions = append(versions, e)
	}
	// Perform bulk insert of documents
	failures := m.repository.BulkInsert(ctx, versions)

	// Collect the failed events if any
	var failedEvents []model.Event
	if failures != nil {
		for _, i := range failures {
			failedEvents = append(failedEvents, events[i])
		}
		return failedEvents
	}
	return nil
}

// getPreviousVersionDbIdFromEvent generate previous version db identifier from the incoming event
func getPreviousVersionDbIdFromEvent(event *model.Event) (*string, error) {
	if previousVersion, ok := event.UpdateDescription.UpdatedFields["previousVersion"].(string); ok {
		_, _, found := strings.Cut(previousVersion, constants.VersionSeparator)
		if !found {
			return nil, errors.New("previousVersion value is in invalid format")
		}
		return repository.GetDbVersionIdFromDbIdAndVersion(event.DocumentKey.Id, previousVersion), nil
	} else {
		return nil, errors.New("previousVersion not present in the event")
	}

}
