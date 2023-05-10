package services

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/model/upload"
	"trimble.com/tdrive/api/utils"
)

// UploadService -> struct for upload service
type UploadService struct {
	fileSpaceService *FilespaceService
	folderService    *FolderService
	fileRepository   *repository.LatestRepository
	blobService      *BlobService
	uploadRepository *repository.FileUploadRepository
}

// NewUploadService creates new UploadService
func NewUploadService(filespaceService *FilespaceService, folderService *FolderService, latestRepository *repository.LatestRepository, blobService *BlobService, uploadRepository *repository.FileUploadRepository) *UploadService {
	return &UploadService{
		fileSpaceService: filespaceService,
		folderService:    folderService,
		fileRepository:   latestRepository,
		blobService:      blobService,
		uploadRepository: uploadRepository,
	}
}

// FileUpload perform file upload
func (u *UploadService) FileUpload(ctx *requestcontext.RequestContext, spaceId string, request upload.FileUploadRequest, headers upload.FileUploadHeaders, params upload.FileUploadParams) (*upload.FileUploadResponse, *api_error.ApiError) {
	if err := validateUploadContents(request.Contents); err != nil {
		return nil, err
	}
	ifMatch, ifMatchErr := utils.ValidateAndGetIfMatchHeader(headers.IfMatch)
	if ifMatchErr != nil {
		return nil, ifMatchErr
	}
	// Validate file space
	space, spaceErr := u.fileSpaceService.getAndValidateAccess(ctx, spaceId, nil)
	if spaceErr != nil {
		return nil, spaceErr
	}

	// Name validation
	// This will not be performed when `enforceSafeNames` set explicitly to false
	if space.EnforceSafeNames {
		if validNameErr := utils.ValidName(request.Name); validNameErr != nil {
			return nil, validNameErr
		}
	}

	// Validate parent folder
	_, folderErr := u.folderService.ValidateFolder(ctx, spaceId, request.ParentFolderId)
	if folderErr != nil {
		return nil, folderErr
	}
	// Validate file Id
	existingFile, validateErr := u.preCheckFileIdAndName(ctx, spaceId, request.Id, request.ParentFolderId, request.Name)
	if validateErr != nil {
		return nil, validateErr
	}
	// Validate IfNoneMatch And IfMatch
	preCondErr := preCheckIfNoneMatchAndIfMatch(existingFile, headers.IfNoneMatch, ifMatch)
	if preCondErr != nil {
		return nil, preCondErr
	}
	if request.Id == nil && existingFile != nil {
		request.Id = existingFile.Id
	}
	// Save entity in DB
	fileUploadEntity := createFileUploadEntity(spaceId, request, headers, ctx.UserId(), ifMatch)
	if createErr := u.uploadRepository.Insert(ctx, fileUploadEntity); createErr != nil {
		return nil, createErr
	}
	return u.convertFileUploadEntityToResponse(ctx, fileUploadEntity, request.Name, params.UrlExpiryInMins)
}

// GetFileUploadDetails Returns file upload details for the given uploadId
func (u *UploadService) GetFileUploadDetails(ctx *requestcontext.RequestContext, spaceId string, uploadId string, params upload.GetFileUploadDetailsParams) (*upload.FileUploadResponse, *api_error.ApiError) {
	// Validate file space
	_, spaceErr := u.fileSpaceService.getAndValidateAccess(ctx, spaceId, nil)
	if spaceErr != nil {
		return nil, spaceErr
	}
	// Get Upload from DB
	fileUploadEntity, err := u.uploadRepository.FindByUploadId(ctx, uploadId)
	if err != nil {
		return nil, err
	}
	if fileUploadEntity.CreatedBy != ctx.UserId() {
		return nil, api_error.PermissionDenied
	}
	return u.convertFileUploadEntityToResponse(ctx, fileUploadEntity, fileUploadEntity.Input.Name, params.UrlExpiryInMins)
}

// preCheckFileId validate the given fileId based on id, parentFolderId and Name
func (u *UploadService) preCheckFileIdAndName(ctx *requestcontext.RequestContext, spaceId string, fileId *string, parentFolderId string, fileName *string) (*repository.Latest, *api_error.ApiError) {
	var existingFile *repository.Latest
	var err *api_error.ApiError
	if fileId != nil {
		// check file with given Id exists, if yes validate parentId
		existingFile, err = u.fileRepository.GetById(ctx, spaceId, constants.File, *fileId)
		if err != nil && err != api_error.ResourceNotFound {
			return nil, err
		}
		if existingFile == nil && fileName == nil {
			// validate file name is present if no entry in DB with given id
			ctx.Logger().Debug("No file in DB with given id and name is not provided as well")
			return nil, api_error.InvalidUploadPayload
		}
		if existingFile != nil && existingFile.ParentFolderId != nil && *existingFile.ParentFolderId != parentFolderId {
			ctx.Logger().Debug("Existing file in DB has different parent")
			return nil, api_error.DuplicateKey
		}
	}
	// Fetch existingFile in DB, if not fetched already
	return u.preCheckFileByNameAndParent(ctx, spaceId, fileId, parentFolderId, fileName, existingFile)
}

// preCheckFileByName validate the given file exists based on parentFolderId and Name
func (u *UploadService) preCheckFileByNameAndParent(ctx *requestcontext.RequestContext, spaceId string, fileId *string, parentFolderId string, fileName *string, existingFile *repository.Latest) (*repository.Latest, *api_error.ApiError) {
	var err *api_error.ApiError
	// Fetch existingFile in DB, if not fetched already
	if existingFile == nil {
		existingFile, err = u.fileRepository.FindByParentFolderIdNameAndType(ctx, spaceId, parentFolderId, *fileName, constants.File)
		if err != nil && err != api_error.ResourceNotFound {
			return nil, err
		}
		if existingFile != nil && fileId != nil && fileId != existingFile.Id {
			ctx.Logger().Debug("Existing file in DB has different fileId")
			return nil, api_error.DuplicateKey
		}
	}

	return existingFile, nil
}

// convertFileUploadEntityToResponse converts repository.FileUpload into upload.FileUploadResponse
func (u *UploadService) convertFileUploadEntityToResponse(ctx *requestcontext.RequestContext, fileUploadEntity *repository.FileUpload, fileName *string, urlExpiryInMins *int) (*upload.FileUploadResponse, *api_error.ApiError) {
	var expiryTime time.Time
	if urlExpiryInMins == nil {
		expiryTime = time.Now().UTC().Add(time.Minute * constants.DefaultExpiryInMinutes)
	} else {
		expiryTime = time.Now().UTC().Add(time.Minute * time.Duration(*urlExpiryInMins))
	}
	contents := fileUploadEntity.Input.Contents
	responseContents := make([]upload.ContentDetails, len(contents))
	var i = 0
	for path, content := range contents {
		var uploadUrl *string
		if fileUploadEntity.Status == constants.FileUploadStatusUploadable && content.Status == constants.ContentUploadStatusUploadable {
			// Update File extension incase content is not src
			if fileName != nil && content.Format != nil {
				fileNameString := strings.TrimSuffix(*fileName, filepath.Ext(*fileName)) + "." + strings.ToLower(*content.Format)
				fileName = &fileNameString
			}
			url, err := u.blobService.GenerateUploadUrl(ctx, path, fileName, expiryTime)
			if err != nil {
				return nil, err
			}
			uploadUrl = &url
		}
		responseContents[i] = upload.ContentDetails{
			Format:     content.Format,
			Status:     content.Status,
			UploadMode: content.UploadMode,
			Url:        uploadUrl,
			Size:       content.Size,
			UpdatedOn:  content.UpdatedOn,
		}
		i++
	}
	uploadInput := upload.FileUploadInput{
		SpaceId:        fileUploadEntity.Input.SpaceId,
		Name:           fileUploadEntity.Input.Name,
		Contents:       responseContents,
		FileId:         fileUploadEntity.Input.FileId,
		ParentFolderId: fileUploadEntity.Input.ParentFolderId,
	}
	var uploadResult *upload.FileUploadResult
	if fileUploadEntity.Result != nil {
		uploadResult = &upload.FileUploadResult{
			Version: fileUploadEntity.Result.Version,
		}
	}
	response := &upload.FileUploadResponse{
		UploadId:    fileUploadEntity.UploadId,
		Input:       uploadInput,
		Result:      uploadResult,
		CreatedOn:   fileUploadEntity.CreatedOn,
		CreatedBy:   fileUploadEntity.CreatedBy,
		ModifiedOn:  fileUploadEntity.ModifiedOn,
		Status:      fileUploadEntity.Status,
		ErrorReason: fileUploadEntity.ErrorReason,
	}
	return response, nil
}

// validateUploadContents validates upload.FileUploadContent returns error if any
func validateUploadContents(contents *[]upload.FileUploadContent) *api_error.ApiError {
	if contents == nil {
		return nil
	}
	formatMap := make(map[string]bool, len(*contents))
	for _, content := range *contents {
		format := content.Format
		// validate for format duplication in request
		formatKey := constants.FormatSrc
		if format != nil {
			formatKey = *format
		}
		if formatMap[formatKey] {
			return api_error.DuplicateFormat
		} else {
			formatMap[formatKey] = true
		}
	}
	// Delete source from formatMap if present and validate the contents length
	if formatMap[constants.FormatSrc] {
		delete(formatMap, constants.FormatSrc)
	}
	if len(formatMap) > 5 {
		return api_error.FileContentsLimitExceeded
	}
	return nil
}

// preCheckIfNoneMatchAndIfMatch Pre check the given IfNoneMatch and IfMatch returns error if any
func preCheckIfNoneMatchAndIfMatch(existingFile *repository.Latest, ifNoneMatch *string, ifMatch *repository.Version) *api_error.ApiError {
	if existingFile != nil {
		if ifNoneMatch != nil && *ifNoneMatch == "*" {
			return api_error.FileWithSameNameExists
		}
		if ifMatch != nil && (existingFile.MajorVersion != ifMatch.MajorVersion || existingFile.MinorVersion != *ifMatch.MinorVersion) {
			return api_error.InvalidVersion
		}
	}
	return nil
}

// generateStoragePath Generates the file content storage path
func generateStoragePath(spaceId string, fileId string, uploadId string, format *string) string {
	var storagePathPrefix = constants.DefaultStorageDir
	if format == nil {
		srcFormat := constants.FormatSrc
		format = &srcFormat
	} else if *format == constants.FormatThumbnail {
		storagePathPrefix = constants.ThumbStorageDir
	}
	// Storage Path -> {storagePathPrefix}/{spaceId}/{fileId}/{uploadId}/{format}
	return storagePathPrefix + constants.StoragePathSeparator + spaceId + constants.StoragePathSeparator + fileId +
		constants.StoragePathSeparator + uploadId + constants.StoragePathSeparator + strings.ToLower(*format)
}

// createFileUploadEntity Creates and returns repository.FileUpload
func createFileUploadEntity(spaceId string, request upload.FileUploadRequest, headers upload.FileUploadHeaders, userId string, ifMatch *repository.Version) *repository.FileUpload {
	uploadId := uuid.NewString()
	fileId := request.Id
	if fileId == nil {
		guid := uuid.NewString()
		fileId = &guid
	}
	createdOn := time.Now().UTC()
	contents := make(map[string]repository.UploadContent)
	requestContents := request.Contents
	srcContentAdded := false
	if requestContents != nil {
		for _, content := range *requestContents {
			storagePath := generateStoragePath(spaceId, *fileId, uploadId, content.Format)
			contents[storagePath] = createUploadContent(content.Format, createdOn)
			if content.Format == nil {
				srcContentAdded = true
			}
		}
	}
	if !srcContentAdded {
		storagePath := generateStoragePath(spaceId, *fileId, uploadId, nil)
		contents[storagePath] = createUploadContent(nil, createdOn)
	}

	uploadInput := repository.UploadInput{
		SpaceId:        spaceId,
		Name:           request.Name,
		FileId:         *fileId,
		ParentFolderId: request.ParentFolderId,
		Contents:       contents,
		IfMatch:        ifMatch,
		IfNoneMatch:    headers.IfNoneMatch,
	}

	fileUploadEntity := &repository.FileUpload{
		UploadId:   uploadId,
		Status:     constants.FileUploadStatusUploadable,
		Input:      uploadInput,
		CreatedOn:  createdOn,
		CreatedBy:  userId,
		ModifiedOn: createdOn,
	}
	return fileUploadEntity
}

// createUploadContent creates and returns repository.UploadContent
func createUploadContent(format *string, updatedOn time.Time) repository.UploadContent {
	return repository.UploadContent{
		Format:     format,
		Status:     constants.ContentUploadStatusUploadable,
		UploadMode: constants.ContentUploadModeSinglepart,
		UpdatedOn:  updatedOn,
	}
}
