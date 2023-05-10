package repository

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"trimble.com/common/api_error"
	"trimble.com/common/requestcontext"
)

// ResourceVersionRepository -> struct for versions repository
type ResourceVersionRepository struct {
	collection *mongo.Collection
}

// ResourceVersion -> struct for versions mongodb entity
type ResourceVersion struct {
	DocId           string    `bson:"_id,omitempty"`
	Id              string    `bson:"id,omitempty"`
	SpaceId         string    `bson:"spaceId,omitempty"`
	MajorVersion    int64     `bson:"majorVersion,omitempty"`
	MinorVersion    int64     `bson:"minorVersion,omitempty"`
	Name            string    `bson:"name,omitempty"`
	Type            string    `bson:"type,omitempty"`
	ParentFolderId  string    `bson:"parentFolderId,omitempty"`
	Size            int64     `bson:"size,omitempty"`
	Contents        []Content `bson:"contents"`
	Deleted         bool      `bson:"deleted"`
	Created         time.Time `bson:"created,omitempty"`
	Modified        time.Time `bson:"modified,omitempty"`
	CreatedBy       string    `bson:"createdBy,omitempty"`
	ModifiedBy      string    `bson:"modifiedBy,omitempty"`
	CreatedOnClock  *int64    `bson:"createdOnClock,omitempty"`
	ModifiedOnClock *int64    `bson:"modifiedOnClock,omitempty"`
	CheckedOutBy    string    `bson:"checkedOutBy,omitempty"`
	CheckedOutOn    string    `bson:"checkedOutOn,omitempty"`
}

type Version struct {
	MajorVersion int64  `json:"majorVersion"`
	MinorVersion *int64 `json:"minorVersion"`
}

// NewResourceVersionRepository creates new versions repository instance
func NewResourceVersionRepository(db *mongo.Database) *ResourceVersionRepository {
	return &ResourceVersionRepository{
		collection: db.Collection("versions"),
	}
}

// BulkInsert Executes InsertMany, an insert command to insert multiple documents into the collection.
// Return indexes of failed documents if any.
func (r ResourceVersionRepository) BulkInsert(ctx *requestcontext.RequestContext, inserts []interface{}) []int {

	manyOptions := &options.InsertManyOptions{}

	// This will enforce that writes will continue by skipping the document producing error and collecting those errors if any.
	// As per documentation, setting this might also boost performance.
	manyOptions.SetOrdered(false)

	// Insert multiple documents
	_, err := r.collection.InsertMany(ctx.DbCtx(), inserts, manyOptions)

	// Collect and return the failures if any.
	var failures []int
	if err != nil {
		if errors.As(err, &mongo.BulkWriteException{}) {
			bulkWriteErr := err.(mongo.BulkWriteException)
			for _, writeErr := range bulkWriteErr.WriteErrors {
				if mongo.IsDuplicateKeyError(writeErr) {
					ctx.Logger().Warn("Item with same id already exists", zap.Error(writeErr))
					continue
				} else {
					failures = append(failures, writeErr.Index)
				}
			}
		} else {
			ctx.Logger().Warn("Unexpected error in bulk insert", zap.Error(err))
			for i := range inserts {
				failures = append(failures, i)
			}
		}
		return failures
	}
	return nil
}

// GetResourceVersionByDbId Fetches the resource details of a specific version using db version identifier
// Returns the versions document and error if any
func (r ResourceVersionRepository) GetResourceVersionByDbId(ctx *requestcontext.RequestContext, spaceId string, id string, result any) error {
	return r.collection.FindOne(ctx.DbCtx(), bson.D{{Key: "spaceId", Value: spaceId}, {Key: "_id", Value: id}}).Decode(result)
}

// GetResourceVersionById Fetches the resource details of a specific version
// Returns the versions document and error if any
func (r ResourceVersionRepository) GetResourceVersionById(ctx *requestcontext.RequestContext, spaceId string, id string, version Version, objectType string) (*ResourceVersion, *api_error.ApiError) {
	var resourceVersion *ResourceVersion
	var minorVersion int64
	var err error
	if version.MinorVersion != nil {
		minorVersion = *version.MinorVersion
		dbVersionId := GetDbVersionId(spaceId, objectType, id, version.MajorVersion, *version.MinorVersion)
		err = r.GetResourceVersionByDbId(ctx, spaceId, *dbVersionId, &resourceVersion)
	} else {
		dbId := GetDbId(spaceId, objectType, &id)
		filter := bson.D{{Key: "spaceId", Value: spaceId}, {Key: "id", Value: dbId}, {Key: "majorVersion", Value: version.MajorVersion}}
		opts := &options.FindOneOptions{Sort: bson.D{{Key: "minorVersion", Value: -1}}, Collation: &options.Collation{Locale: "en"}}
		err = r.collection.FindOne(ctx.DbCtx(), filter, opts).Decode(&resourceVersion)
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.Logger().Debug("Resource version not found", zap.String("id", id), zap.Int64("majorVersion", version.MajorVersion), zap.Int64("minorVersion", minorVersion))
			return nil, api_error.VersionNotFound
		}
		ctx.Logger().Error("Unexpected error on get resource version by id", zap.Error(err))
		return nil, api_error.InternalServerError
	}
	resourceVersion.Id = *GetVisibleId(&resourceVersion.Id)
	resourceVersion.ParentFolderId = *GetVisibleId(&resourceVersion.ParentFolderId)
	return resourceVersion, nil
}

// findItems Executes a find command and returns a list of the matching documents in the "versions" collection
//
// The filter parameter must be a document containing query operators and can be used to select which documents are
// included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all documents.
//
// The opts parameter can be used to specify options (Ex: Sort, Limit) for the operation (see the options.FindOptions documentation).
func (r ResourceVersionRepository) findItems(ctx *requestcontext.RequestContext, filter bson.D, options *options.FindOptions) ([]*ResourceVersion, *api_error.ApiError) {

	var results []*ResourceVersion
	findResult, err := r.collection.Find(ctx.DbCtx(), filter, options)
	if err != nil {
		ctx.Logger().Error("Error on finding items in space", zap.Error(err))
		return nil, api_error.InternalServerError
	}
	if err = findResult.All(ctx.DbCtx(), &results); err != nil {
		ctx.Logger().Error("Error on mapping find result to entity", zap.Error(err))
		return nil, api_error.InternalServerError
	}

	return results, nil

}

// ListResourceVersions Fetch list of versions for a given resource
func (r ResourceVersionRepository) ListResourceVersions(
	ctx *requestcontext.RequestContext,
	spaceId string,
	resourceId string,
	resourceType string,
	fromVersion *Version,
	toVersion *Version,
	cursor *Cursor,
	sort Sort,
	limit int64) ([]*ResourceVersion, *api_error.ApiError) {

	dbId := GetDbId(spaceId, resourceType, &resourceId)

	// filter criteria to find the document
	filter := bson.D{
		{Key: "spaceId", Value: spaceId},
		{Key: "id", Value: dbId},
	}

	filter = applyFromVersionAndToVersionFilter(fromVersion, toVersion, filter)

	// Get sort direction and range operator to be used to move cursor/pages.
	order, rangeOp := GetDBSortDirectionAndRangeOperator(sort)

	filter = applyCursorFilterForResourceVersions(cursor, sort, rangeOp, filter)

	sortOption := bson.D{}
	for _, field := range sort.DBFields {
		sortOption = append(sortOption, bson.E{Key: field, Value: order})
	}
	// In case of sorting by non-unique field, add secondary sorting parameter with unique property ("modifiedOnClock").
	if !sort.UniquePageableProperty {
		sortOption = append(sortOption, bson.E{Key: "modifiedOnClock", Value: order})
	}

	findOptions := &options.FindOptions{Limit: &limit, Sort: sortOption, Collation: &options.Collation{Locale: "en"}}
	if results, err := r.findItems(ctx, filter, findOptions); err != nil {
		return nil, err
	} else {
		return results, nil
	}

}

func applyFromVersionAndToVersionFilter(fromVersion *Version, toVersion *Version, filter bson.D) bson.D {
	if fromVersion == nil && toVersion == nil {
		return filter
	}
	// If both from and to point to same version
	if fromVersion != nil && toVersion != nil && fromVersion.MajorVersion == toVersion.MajorVersion {
		versionFilters := bson.A{}
		// majorVersion = from majorVersion and minorVersion >= from minorVersion
		versionFilter1 := bson.D{{Key: "majorVersion", Value: fromVersion.MajorVersion}}
		if fromVersion.MinorVersion != nil {
			versionFilter1 = append(versionFilter1, bson.E{Key: "minorVersion", Value: bson.D{{Key: "$gte", Value: fromVersion.MinorVersion}}})
		}
		// majorVersion = to majorVersion and minorVersion <= from minorVersion
		versionFilter2 := bson.D{{Key: "majorVersion", Value: toVersion.MajorVersion}}
		if toVersion.MinorVersion != nil {
			versionFilter2 = append(versionFilter2, bson.E{Key: "minorVersion", Value: bson.D{{Key: "$lte", Value: toVersion.MinorVersion}}})
		}
		versionFilters = append(versionFilters, versionFilter1, versionFilter2)
		filter = append(filter, bson.E{
			Key:   "$and",
			Value: versionFilters,
		})
		return filter
	}
	// majorVersion = from majorVersion and minorVersion >= from minorVersion
	versionFilter1 := bson.D{}
	// majorVersion > from majorVersion and majorVersion < to majorVersion
	versionFilter2 := bson.D{}
	// majorVersion = to majorVersion and minorVersion <= to minorVersion
	versionFilter3 := bson.D{}

	if fromVersion != nil {
		versionFilter1 = append(versionFilter1, bson.E{Key: "majorVersion", Value: fromVersion.MajorVersion})
		versionFilter2 = append(versionFilter2, bson.E{Key: "majorVersion", Value: bson.D{{Key: "$gt", Value: fromVersion.MajorVersion}}})
		if fromVersion.MinorVersion != nil {
			versionFilter1 = append(versionFilter1, bson.E{Key: "minorVersion", Value: bson.D{{Key: "$gte", Value: fromVersion.MinorVersion}}})
		}
	}
	if toVersion != nil {
		versionFilter2 = append(versionFilter2, bson.E{Key: "majorVersion", Value: bson.D{{Key: "$lt", Value: toVersion.MajorVersion}}})
		versionFilter3 = append(versionFilter3, bson.E{Key: "majorVersion", Value: toVersion.MajorVersion})
		if toVersion.MinorVersion != nil {
			versionFilter3 = append(versionFilter3, bson.E{Key: "minorVersion", Value: bson.D{{Key: "$lte", Value: toVersion.MinorVersion}}})
		}
	}
	versionFilters := bson.A{}
	if len(versionFilter1) > 0 {
		versionFilters = append(versionFilters, versionFilter1)
	}
	versionFilters = append(versionFilters, versionFilter2)
	if len(versionFilter3) > 0 {
		versionFilters = append(versionFilters, versionFilter3)
	}
	combinedVersionFilter := bson.E{
		Key:   "$or",
		Value: versionFilters,
	}
	filter = append(filter, combinedVersionFilter)
	return filter
}

func applyCursorFilterForResourceVersions(cursor *Cursor, sort Sort, rangeOp string, filter bson.D) bson.D {
	if cursor == nil {
		return filter
	}
	// Filter criteria if cursor present.
	// In case of non-unique field, criteria will be combination of (non-unique field and unique field) OR (non-unique field) filters.
	// Here,"modifiedOnClock" represents a unique property to a "versions" document,
	// it will be added in the filter combination in case of non-unique cursor
	if cursor.NonUniqueValue != nil {
		cursorFilter1 := bson.D{{Key: "modifiedOnClock", Value: bson.D{{Key: rangeOp, Value: cursor.UniqueValue["modifiedOnClock"]}}}}
		cursorFilter2 := bson.D{}
		for _, field := range sort.DBFields {
			cursorFilter1 = append(cursorFilter1, bson.E{Key: field, Value: cursor.NonUniqueValue[field]})
			cursorFilter2 = append(cursorFilter2, bson.E{Key: field, Value: bson.D{{Key: rangeOp, Value: cursor.NonUniqueValue[field]}}})
		}
		nonUniqueCursor := bson.E{
			Key: "$or",
			Value: bson.A{
				cursorFilter1,
				cursorFilter2,
			},
		}
		filter = append(filter, nonUniqueCursor)
	} else {
		var cursorFilter1 bson.D
		cursorFilter2 := bson.D{}
		for _, field := range sort.DBFields {
			if field == "majorVersion" {
				cursorFilter1 = append(cursorFilter1, bson.E{Key: field, Value: cursor.UniqueValue[field]})
				cursorFilter2 = append(cursorFilter2, bson.E{Key: field, Value: bson.D{{Key: rangeOp, Value: cursor.UniqueValue[field]}}})
			} else if field == "minorVersion" {
				cursorFilter1 = append(cursorFilter1, bson.E{Key: field, Value: bson.D{{Key: rangeOp, Value: cursor.UniqueValue[field]}}})
			}
			if len(cursorFilter2) > 0 {
				filter = append(filter, bson.E{
					Key:   "$or",
					Value: bson.A{cursorFilter1, cursorFilter2},
				})
			} else {
				filter = append(filter, bson.E{Key: field, Value: bson.D{{Key: rangeOp, Value: cursor.UniqueValue[field]}}})
			}
		}
	}
	return filter
}

func ConvertResourceVersionToVisibleId(resourceVersion *ResourceVersion) {
	resourceVersion.Id = *GetVisibleId(&resourceVersion.Id)
	resourceVersion.ParentFolderId = *GetVisibleId(&resourceVersion.ParentFolderId)
}
