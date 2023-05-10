package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"trimble.com/common/api_error"
	"trimble.com/common/requestcontext"
)

type FilespaceRepository struct {
	db               *mongo.Database
	collection       *mongo.Collection
	latestRepository *LatestRepository
}

func NewFilespaceRepository(db *mongo.Database, latestRepository *LatestRepository) *FilespaceRepository {
	return &FilespaceRepository{
		db:               db,
		collection:       db.Collection("filespace"),
		latestRepository: latestRepository,
	}
}

type Filespace struct {
	ID               *string    `bson:"_id"`
	Description      *string    `bson:"description,omitempty"`
	RootId           *string    `bson:"rootId"`
	Deleted          *bool      `bson:"deleted"`
	EnforceSafeNames bool       `bson:"enforceSafeNames"`
	Clock            *int64     `bson:"clock"`
	Created          *time.Time `bson:"created"`
	CreatedBy        *string    `bson:"createdBy"`
	Modified         *time.Time `bson:"modified"`
	ModifiedBy       *string    `bson:"modifiedBy"`
}

func NewFilespace(description *string, created time.Time, createdBy string) *Filespace {
	var clock int64 = 0
	var deleted = false
	return &Filespace{
		Description: description,
		Deleted:     &deleted,
		Clock:       &clock,
		Created:     &created,
		CreatedBy:   &createdBy,
		Modified:    &created,
		ModifiedBy:  &createdBy,
	}
}

// Insert Insert Filespace
func (r *FilespaceRepository) Insert(ctx *requestcontext.RequestContext, filespace *Filespace) *api_error.ApiError {
	ctx.Logger().Debug("Inserting filespace", zap.Any("fileSpace", filespace))

	// Start transaction session
	dbContext := context.TODO()
	session, sessionErr := r.db.Client().StartSession()
	if sessionErr != nil {
		ctx.Logger().Error("Unexpected error on starting transaction session", zap.Error(sessionErr))
		return api_error.InternalServerError
	}
	// Defer end session
	defer session.EndSession(dbContext)

	// Perform operations in transaction context
	_, err := session.WithTransaction(dbContext, func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Generate rootId
		rootFolderId := filespace.ID
		filespace.RootId = rootFolderId

		// Create FileSpace
		_, err := r.collection.InsertOne(sessCtx, filespace)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				ctx.Logger().Debug("File space with same id already exists")
				return nil, api_error.DuplicateKey
			} else {
				ctx.Logger().Error("Unexpected error on insert filespace", zap.Error(err))
				return nil, api_error.InternalServerError
			}
		}

		// Create Root folder
		ctx.SetDbCtx(sessCtx)
		rootFolderErr := r.latestRepository.CreateRootFolder(ctx, filespace)
		if rootFolderErr != nil {
			return nil, rootFolderErr
		}
		return nil, nil
	})
	if err != nil {
		ctx.Logger().Error("Unexpected error on performing inserts in transaction", zap.Error(err))
		return err.(*api_error.ApiError)
	}
	return nil
}

// GetById Get Filespace by Id
func (r *FilespaceRepository) GetById(ctx *requestcontext.RequestContext, filespaceId string) (*Filespace, *api_error.ApiError) {
	ctx.Logger().Debug("Fetching filespace with id", zap.String("fileSpaceId", filespaceId))
	var result *Filespace
	err := r.collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: filespaceId}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, api_error.SpaceNotFound
		}
		ctx.Logger().Error("Unexpected error on get filespace by id", zap.Error(err))
		return nil, api_error.InternalServerError
	}
	return result, nil
}

// UpdateDescription Update Filespace Description by Id
func (r *FilespaceRepository) UpdateDescription(ctx *requestcontext.RequestContext, filespace *Filespace) (*Filespace, *api_error.ApiError) {
	ctx.Logger().Debug("Updating filespace Description", zap.Any("fileSpace", filespace))
	filter := bson.D{{Key: "_id", Value: *filespace.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "description", Value: filespace.Description},
		{Key: "modified", Value: filespace.Modified},
		{Key: "modifiedBy", Value: filespace.ModifiedBy},
		{Key: "deleted", Value: false},
	}}, {
		Key:   "$inc",
		Value: bson.D{{Key: "clock", Value: 1}},
	}}
	filespace, err := r.findOneAndUpdate(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return filespace, nil

}

func (r *FilespaceRepository) findOneAndUpdate(ctx *requestcontext.RequestContext, filter bson.D, update bson.D) (*Filespace, *api_error.ApiError) {
	after := options.After
	updateOptions := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := r.collection.FindOneAndUpdate(ctx.DbCtx(), filter, update, updateOptions)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			ctx.Logger().Debug("File space is not found")
			return nil, api_error.SpaceNotFound
		} else {
			ctx.Logger().Error("Unexpected error in updating the space", zap.Error(result.Err()))
			return nil, api_error.InternalServerError
		}
	}
	var updatedFileSpace *Filespace
	err := result.Decode(&updatedFileSpace)
	if err != nil {
		return nil, api_error.InternalServerError
	}
	return updatedFileSpace, nil
}

// IncrementClock Increments Filespace clock
func (r *FilespaceRepository) IncrementClock(ctx *requestcontext.RequestContext, spaceId string) (*Filespace, *api_error.ApiError) {
	ctx.Logger().Debug("Updating filespace clock", zap.String("spaceId", spaceId))
	filter := bson.D{{Key: "_id", Value: spaceId}}
	update := bson.D{{
		Key:   "$inc",
		Value: bson.D{{Key: "clock", Value: 1}},
	}}
	updatedFileSpace, err := r.findOneAndUpdate(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return updatedFileSpace, nil
}

// UpdateDeleted Update Filespace deleted status by Id
func (r *FilespaceRepository) UpdateDeleted(ctx *requestcontext.RequestContext, filespace *Filespace) (*Filespace, *api_error.ApiError) {
	ctx.Logger().Debug("Deleting filespace: ", zap.Any("fileSpace", filespace))
	filter := bson.D{{Key: "_id", Value: *filespace.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "deleted", Value: filespace.Deleted},
		{Key: "modified", Value: filespace.Modified},
		{Key: "modifiedBy", Value: filespace.ModifiedBy},
	}}, {
		Key:   "$inc",
		Value: bson.D{{Key: "clock", Value: 1}},
	}}
	filespace, err := r.findOneAndUpdate(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return filespace, nil
}
