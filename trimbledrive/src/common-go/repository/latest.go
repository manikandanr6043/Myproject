package repository

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/requestcontext"
)

// LatestRepository -> struct for latest repository
type LatestRepository struct {
	collection *mongo.Collection
}

// Latest -> struct for latest mongodb entity
type Latest struct {
	Id              *string   `bson:"_id,omitempty"`
	SpaceId         string    `bson:"spaceId,omitempty"`
	MajorVersion    int64     `bson:"majorVersion,omitempty"`
	MinorVersion    int64     `bson:"minorVersion"`
	Name            string    `bson:"name,omitempty"`
	Type            string    `bson:"type,omitempty"`
	ParentFolderId  *string   `bson:"parentFolderId,omitempty"`
	Size            int64     `bson:"size,omitempty"`
	Contents        []Content `bson:"contents"`
	Deleted         bool      `bson:"deleted"`
	Created         time.Time `bson:"created,omitempty"`
	Modified        time.Time `bson:"modified,omitempty"`
	CreatedBy       string    `bson:"createdBy,omitempty"`
	ModifiedBy      string    `bson:"modifiedBy,omitempty"`
	CheckedOutBy    string    `bson:"checkedOutBy,omitempty"`
	CheckedOutOn    string    `bson:"checkedOutOn,omitempty"`
	CreatedOnClock  *int64    `bson:"createdOnClock,omitempty"`
	ModifiedOnClock *int64    `bson:"modifiedOnClock,omitempty"`
	PreviousVersion *string   `bson:"previousVersion,omitempty"`
}

type Content struct {
	Format   *string `bson:"format,omitempty"`
	Location string  `bson:"location"`
	Size     int64   `bson:"size"`
}

type AggregationResponse struct {
	Latest        Latest
	ParentFolders []Latest
}

// NewLatestRepository creates new latest repository instance
func NewLatestRepository(db *mongo.Database) *LatestRepository {
	return &LatestRepository{
		collection: db.Collection("latest"),
	}
}

// Create Creates a new latest entry in DB.
// Returns saved object and error if any.
func (r LatestRepository) Create(ctx *requestcontext.RequestContext, latest *Latest) *api_error.ApiError {
	// Generate ID if id is not present
	if latest.Id == nil {
		id := uuid.NewString()
		latest.Id = GetDbId(latest.SpaceId, latest.Type, &id)
	} else {
		latest.Id = GetDbId(latest.SpaceId, latest.Type, latest.Id)
	}

	// Set system default fields
	timestamp := time.Now()
	latest.MajorVersion = 1
	latest.MinorVersion = 0
	latest.Created = timestamp
	latest.Modified = timestamp
	latest.CreatedBy = ctx.UserId()
	latest.ModifiedBy = ctx.UserId()
	latest.ParentFolderId = GetDbId(latest.SpaceId, constants.Folder, latest.ParentFolderId)

	_, err := r.collection.InsertOne(ctx.DbCtx(), latest)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.Logger().Debug("File/Folder with same id or name already exists")
			return api_error.DuplicateKey
		} else {
			ctx.Logger().Error("Unexpected error on inserting file/folder", zap.Error(err))
			return api_error.InternalServerError
		}
	}
	ctx.Logger().Debug("Inserted object", zap.Any("objectId", latest.Id))
	latest.Id = GetVisibleId(latest.Id)
	latest.ParentFolderId = GetVisibleId(latest.ParentFolderId)
	return nil
}

// GetById Fetches latest entry with given identifier and type in DB.
// Returns the latest object and error if any.
func (r LatestRepository) GetById(ctx *requestcontext.RequestContext, spaceId string, objectType string, id string) (*Latest, *api_error.ApiError) {
	dbId := GetDbId(spaceId, objectType, &id)
	var latest *Latest
	err := r.collection.FindOne(ctx.DbCtx(), bson.D{{Key: "_id", Value: dbId}, {Key: "spaceId", Value: spaceId}, {Key: "type", Value: objectType}}).Decode(&latest)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.Logger().Debug("File or Folder not found", zap.String("objectId", id), zap.String("objectType", objectType))
			return nil, api_error.ResourceNotFound
		}
		ctx.Logger().Error("Unexpected error on get folder by id", zap.Error(err))
		return nil, api_error.InternalServerError
	}
	latest.Id = GetVisibleId(latest.Id)
	latest.ParentFolderId = GetVisibleId(latest.ParentFolderId)
	return latest, nil
}

// FindByParentFolderIdNameAndType Fetches the latest file/folder entry with given parentFolderId, name and type in DB.
// Returns the latest object and error if any.
func (r LatestRepository) FindByParentFolderIdNameAndType(ctx *requestcontext.RequestContext, spaceId string, parentFolderId string, name string, objectType string) (*Latest, *api_error.ApiError) {
	var latest *Latest
	dbParentFolderId := GetDbId(spaceId, constants.Folder, &parentFolderId)
	err := r.collection.FindOne(ctx.DbCtx(), bson.D{{Key: "spaceId", Value: spaceId}, {Key: "parentFolderId", Value: dbParentFolderId}, {Key: "name", Value: name}, {Key: "type", Value: objectType}}).Decode(&latest)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.Logger().Debug("File/Folder with given name not found for the parent ", zap.String("objectType", objectType), zap.String("name", name), zap.String("parentFolderId", parentFolderId))
			return nil, api_error.ResourceNotFound
		}
		ctx.Logger().Error("Unexpected error on findFileByParentFolderIdNameAndType", zap.Error(err))
		return nil, api_error.InternalServerError
	}
	latest.Id = GetVisibleId(latest.Id)
	latest.ParentFolderId = GetVisibleId(latest.ParentFolderId)
	return latest, nil
}

// UpdateNameOrParentFolder Updates latest entity with given name or parent folder id
// Returns updated object and error if any.
func (r LatestRepository) UpdateNameOrParentFolder(ctx *requestcontext.RequestContext, space *Filespace, objectType string, id string, name *string, parentFolderId *string, version *Version) (*Latest, *api_error.ApiError) {
	dbId := GetDbId(*space.ID, objectType, &id)
	// filter criteria to find the document
	filter := bson.D{{Key: "_id", Value: dbId}, {Key: "spaceId", Value: *space.ID}}
	if version != nil {
		filter = append(filter, bson.E{Key: "majorVersion", Value: version.MajorVersion})
		if version.MinorVersion == nil {
			ctx.Logger().Error("minorVersion is nil")
			return nil, api_error.InternalServerError
		}
		filter = append(filter, bson.E{Key: "minorVersion", Value: version.MinorVersion})
	}

	updateFields := bson.D{
		{Key: "modified", Value: time.Now()},
		{Key: "modifiedBy", Value: ctx.UserId()},
		{Key: "deleted", Value: false},
		{Key: "modifiedOnClock", Value: *space.Clock},
		{Key: "previousVersion", Value: bson.M{"$concat": bson.A{bson.M{"$toString": "$majorVersion"}, ".", bson.M{"$toString": "$minorVersion"}}}},
	}

	if name != nil {
		updateFields = append(updateFields, bson.E{Key: "name", Value: name})
	}

	if parentFolderId != nil {
		parentFolderId = GetDbId(*space.ID, constants.Folder, parentFolderId)
		updateFields = append(updateFields, bson.E{Key: "parentFolderId", Value: parentFolderId})
	}
	if objectType == constants.Folder {
		updateFields = append(updateFields, bson.E{Key: "majorVersion", Value: bson.D{{Key: "$sum", Value: bson.A{"$majorVersion", 1}}}})
	} else {
		updateFields = append(updateFields, bson.E{Key: "minorVersion", Value: bson.D{{Key: "$sum", Value: bson.A{"$minorVersion", 1}}}})
	}

	// update parameters to perform update operation
	updateStage := bson.D{{Key: "$set", Value: updateFields}}

	after := options.After
	updateOptions := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := r.collection.FindOneAndUpdate(ctx.DbCtx(), filter, mongo.Pipeline{updateStage}, updateOptions)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			if version != nil {
				ctx.Logger().Debug("Given resource or version is not found.")
				return nil, api_error.InvalidVersion
			} else {
				return nil, api_error.ResourceNotFound
			}
		} else if mongo.IsDuplicateKeyError(result.Err()) {
			ctx.Logger().Debug("File or folder with same name already exists")
			return nil, api_error.DuplicateResourceName
		} else {
			ctx.Logger().Error("Unexpected error on folder update", zap.Error(result.Err()))
			return nil, api_error.InternalServerError
		}
	}
	var latest *Latest
	err := result.Decode(&latest)
	if err != nil {
		return nil, api_error.InternalServerError
	}
	latest.Id = GetVisibleId(latest.Id)
	latest.ParentFolderId = GetVisibleId(latest.ParentFolderId)
	return latest, nil
}

// UpdateFileContent Updates latest file entity with given name on content update
// Returns updated object and error if any.
func (r LatestRepository) UpdateFileContent(ctx *requestcontext.RequestContext, space *Filespace, id string, parentFolderId string,
	size int64, contents []Content, name *string, version *Version) (*Latest, *api_error.ApiError) {
	dbFileId := GetDbId(*space.ID, constants.File, &id)
	dbParentFolderId := GetDbId(*space.ID, constants.Folder, &parentFolderId)
	// filter criteria to find the document
	// Deleted check is not added as idea is to increment version and undelete file
	filter := bson.D{{Key: "_id", Value: dbFileId}, {Key: "spaceId", Value: *space.ID}, {Key: "parentFolderId", Value: dbParentFolderId}}
	if version != nil {
		filter = append(filter, bson.E{Key: "majorVersion", Value: version.MajorVersion})
		if version.MinorVersion == nil {
			ctx.Logger().Error("minorVersion is nil")
			return nil, api_error.InternalServerError
		}
		filter = append(filter, bson.E{Key: "minorVersion", Value: version.MinorVersion})
	}
	updateFields := bson.D{
		{Key: "modified", Value: time.Now()},
		{Key: "modifiedBy", Value: ctx.UserId()},
		{Key: "deleted", Value: false},
		{Key: "modifiedOnClock", Value: *space.Clock},
		{Key: "size", Value: size},
		{Key: "contents", Value: contents},
		{Key: "previousVersion", Value: bson.M{"$concat": bson.A{bson.M{"$toString": "$majorVersion"}, ".", bson.M{"$toString": "$minorVersion"}}}},
		{Key: "minorVersion", Value: 0},
		{Key: "majorVersion", Value: bson.D{{Key: "$sum", Value: bson.A{"$majorVersion", 1}}}},
	}
	if name != nil {
		updateFields = append(updateFields, bson.E{Key: "name", Value: name})
	}

	// update parameters to perform update operation
	// On successful update, version will be incremented.
	updateStage := bson.D{
		{
			Key:   "$set",
			Value: updateFields,
		},
	}

	after := options.After
	updateOptions := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := r.collection.FindOneAndUpdate(ctx.DbCtx(), filter, mongo.Pipeline{updateStage}, updateOptions)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			ctx.Logger().Debug("Working on stale object")
			return nil, api_error.StaleObject
		} else if mongo.IsDuplicateKeyError(result.Err()) {
			ctx.Logger().Debug("File or folder with same name already exists")
			return nil, api_error.DuplicateResourceName
		} else {
			ctx.Logger().Error("Unexpected error on file update", zap.Error(result.Err()))
			return nil, api_error.InternalServerError
		}
	}
	var latest *Latest
	err := result.Decode(&latest)
	if err != nil {
		return nil, api_error.InternalServerError
	}
	return latest, nil
}

// DeleteById Mark the latest entry with given identifier as deleted in DB.
// Returns error if the delete operation is unsuccessful.
func (r LatestRepository) DeleteById(ctx *requestcontext.RequestContext, space *Filespace, objectType string, id string, version *Version) (*Latest, *api_error.ApiError) {
	dbId := GetDbId(*space.ID, objectType, &id)
	// filter criteria to find the document
	filter := bson.D{{Key: "_id", Value: dbId}, {Key: "spaceId", Value: *space.ID}}
	if version != nil {
		filter = append(filter, bson.E{Key: "majorVersion", Value: version.MajorVersion})
		if version.MinorVersion == nil {
			ctx.Logger().Error("minorVersion is nil")
			return nil, api_error.InternalServerError
		}
		filter = append(filter, bson.E{Key: "minorVersion", Value: version.MinorVersion})
	}
	// update parameters to perform update operation
	updateFields := bson.D{
		{Key: "deleted", Value: true},
		{Key: "modified", Value: time.Now()},
		{Key: "modifiedBy", Value: ctx.UserId()},
		{Key: "modifiedOnClock", Value: *space.Clock},
		{Key: "previousVersion", Value: bson.M{"$concat": bson.A{bson.M{"$toString": "$majorVersion"}, ".", bson.M{"$toString": "$minorVersion"}}}},
	}

	if objectType == constants.Folder {
		updateFields = append(updateFields, bson.E{Key: "majorVersion", Value: bson.D{{Key: "$sum", Value: bson.A{"$majorVersion", 1}}}})
	} else {
		updateFields = append(updateFields, bson.E{Key: "minorVersion", Value: bson.D{{Key: "$sum", Value: bson.A{"$minorVersion", 1}}}})
	}

	updateStage := bson.D{
		{
			Key:   "$set",
			Value: updateFields,
		},
	}

	after := options.After
	updateOptions := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := r.collection.FindOneAndUpdate(ctx.DbCtx(), filter, mongo.Pipeline{updateStage}, updateOptions)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			if version != nil {
				ctx.Logger().Debug("Given resource or version is not found.")
				return nil, api_error.InvalidVersion
			} else {
				return nil, api_error.ResourceNotFound
			}
		} else {
			ctx.Logger().Error("Unexpected error on delete", zap.Error(result.Err()))
			return nil, api_error.InternalServerError
		}
	}
	var latest *Latest
	err := result.Decode(&latest)
	if err != nil {
		return latest, api_error.InternalServerError
	}
	latest.Id = GetVisibleId(latest.Id)
	latest.ParentFolderId = GetVisibleId(latest.ParentFolderId)
	return latest, nil
}

func (r LatestRepository) CreateRootFolder(ctx *requestcontext.RequestContext, filespace *Filespace) *api_error.ApiError {
	// Create Root folder
	timestamp := time.Now()
	rootFolder := &Latest{
		Id:              GetDbId(*filespace.ID, constants.Folder, filespace.ID),
		SpaceId:         *filespace.ID,
		MajorVersion:    1,
		MinorVersion:    1,
		Name:            *filespace.ID,
		Type:            constants.Folder,
		Created:         timestamp,
		Modified:        timestamp,
		CreatedBy:       ctx.UserId(),
		ModifiedBy:      ctx.UserId(),
		CreatedOnClock:  filespace.Clock,
		ModifiedOnClock: filespace.Clock,
	}

	_, folderErr := r.collection.InsertOne(ctx.DbCtx(), rootFolder)
	if folderErr != nil {
		ctx.Logger().Error("Unexpected error on inserting root folder", zap.Error(folderErr))
		return api_error.InternalServerError
	}
	return nil
}

// GetEntityWithParentFoldersList Fetches the entity along with ParentFoldersList for a file or folder
// Returns ParentFoldersList error if any
func (r LatestRepository) GetEntityWithParentFoldersList(ctx *requestcontext.RequestContext, spaceId string, fileOrFolderId string, objectType string) (*AggregationResponse, *api_error.ApiError) {
	fileOrFolderDBId := GetDbId(spaceId, objectType, &fileOrFolderId)
	result, err := r.ParentFoldersListAggregatePipeline(ctx, spaceId, fileOrFolderDBId)
	if err != nil {
		ctx.Logger().Error("Unexpected Error in fetching the list of parent folders", zap.Error(err))
		return nil, api_error.InternalServerError
	}
	var results []AggregationResponse
	if err = result.All(ctx.DbCtx(), &results); err != nil {
		ctx.Logger().Error("Unexpected Error in fetching the list of parent folders", zap.Error(err))
		return nil, api_error.InternalServerError
	}
	if len(results) > 0 {
		aggregateResult := results[0]
		aggregateResult.Latest.Id = GetVisibleId(aggregateResult.Latest.Id)
		if aggregateResult.Latest.ParentFolderId != nil {
			aggregateResult.Latest.ParentFolderId = GetVisibleId(aggregateResult.Latest.ParentFolderId)
		}
		return &aggregateResult, nil
	}
	return nil, api_error.ResourceNotFound
}

func (r LatestRepository) ParentFoldersListAggregatePipeline(ctx *requestcontext.RequestContext, spaceId string, dbId *string) (*mongo.Cursor, error) {
	ctx.Logger().Debug("id space id", zap.Any("space", spaceId), zap.Any("id", *dbId))
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "spaceId", Value: spaceId}, {Key: "_id", Value: *dbId}}}}
	graphLookUpStage := bson.D{
		{Key: "$graphLookup",
			Value: bson.D{
				{Key: "from", Value: "latest"},
				{Key: "startWith", Value: "$parentFolderId"},
				{Key: "connectFromField", Value: "parentFolderId"},
				{Key: "connectToField", Value: "_id"},
				{Key: "as", Value: "folders"},
				{Key: "depthField", Value: "order"},
			},
		},
	}
	projectStage := bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "latest", Value: "$$ROOT"},
				{Key: "parentFolders",
					Value: bson.D{
						{Key: "$sortArray",
							Value: bson.D{
								{Key: "input", Value: "$folders"},
								{Key: "sortBy", Value: bson.D{{Key: "order", Value: 1}}},
							},
						},
					},
				},
			},
		},
	}
	unsetStage := bson.D{{Key: "$unset", Value: "parentFolders.order"}}
	projectStage2 := bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "latest.folders", Value: 0},
			},
		},
	}

	opts := options.Aggregate()
	result, err := r.collection.Aggregate(ctx.DbCtx(), mongo.Pipeline{matchStage, graphLookUpStage, projectStage, unsetStage, projectStage2}, opts)
	return result, err
}

// findItems Executes a find command and returns a list of the matching documents in the "latest" collection
//
// The filter parameter must be a document containing query operators and can be used to select which documents are
// included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all documents.
//
// The opts parameter can be used to specify options (Ex: Sort, Limit) for the operation (see the options.FindOptions documentation).
func (r LatestRepository) findItems(ctx *requestcontext.RequestContext, filter bson.D, options *options.FindOptions) ([]*Latest, *api_error.ApiError) {

	var results []*Latest
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

// GetChildren Fetch list of children under given folder (folderId)
func (r LatestRepository) GetChildren(
	ctx *requestcontext.RequestContext,
	spaceId string,
	folderId string,
	deleted *bool,
	cursor *Cursor,
	sort Sort,
	limit int64) ([]*Latest, *api_error.ApiError) {

	dbFolderId := GetDbId(spaceId, constants.Folder, &folderId)

	// filter criteria to find the document
	filter := bson.D{
		{Key: "spaceId", Value: spaceId},
		{Key: "parentFolderId", Value: dbFolderId},
	}

	if deleted == nil || !*deleted {
		filter = append(filter, bson.E{Key: "deleted", Value: false})
	}

	// Get sort direction and range operator to be used to move cursor/pages.
	order, rangeOp := getDBSortDirectionAndRangeOperator(sort)

	filter = applyCursorFilterForLatest(cursor, sort, rangeOp, filter)

	sortOption := bson.D{}
	for _, field := range sort.DBFields {
		sortOption = append(sortOption, bson.E{Key: field, Value: order})
	}
	// In case of sorting by non-unique field, add secondary sorting parameter with unique property ("createdOnClock").
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

func applyCursorFilterForLatest(cursor *Cursor, sort Sort, rangeOp string, filter bson.D) bson.D {
	if cursor == nil {
		return filter
	}
	// Filter criteria if cursor present.
	// In case of non-unique field, criteria will be combination of (non-unique field and unique field) OR (non-unique field) filters.
	// Here,"createdOnClock" represents a unique property to a "latest" document,
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
		for _, field := range sort.DBFields {
			filter = append(filter, bson.E{Key: field, Value: bson.D{{Key: rangeOp, Value: cursor.UniqueValue[field]}}})
		}
	}
	return filter
}

// getDBSortDirectionAndRangeOperator Get DB query equivalent of sort direction and range operators
// returns 1, $gt for asc, -1, $lt for desc
func getDBSortDirectionAndRangeOperator(sort Sort) (int, string) {
	if sort.Ascending {
		return 1, "$gt"
	} else {
		return -1, "$lt"
	}
}

// GetTopActiveParentFolderId fetches the TopActiveParentFolderId from parentFoldersList
func GetTopActiveParentFolderId(folders []Latest, rootId string) *string {
	for index, item := range folders {
		if item.Deleted {
			if index-1 >= 0 {
				return GetVisibleId(folders[index-1].Id)
			} else {
				return nil
			}
		}
	}
	if len(folders) > 0 {
		return GetVisibleId(folders[len(folders)-1].Id)
	} else {
		return &rootId
	}
}

func ConvertToVisibleId(latest *Latest) {
	latest.Id = GetVisibleId(latest.Id)
	latest.ParentFolderId = GetVisibleId(latest.ParentFolderId)
}
