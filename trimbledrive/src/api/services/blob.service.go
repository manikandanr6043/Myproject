package services

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"trimble.com/common/api_error"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/config"
)

// BlobService -> struct for blob service
type BlobService struct {
	blobServiceClient *service.Client
	containerName     string
	containerUrl      string
}

// NewBlobService creates new BlobService
func NewBlobService(blobServiceClient *service.Client, tDriveConfig *config.TDriveConfig) *BlobService {
	return &BlobService{
		blobServiceClient: blobServiceClient,
		containerName:     tDriveConfig.StorageAccount.ContainerName,
		containerUrl:      tDriveConfig.StorageAccount.URL + tDriveConfig.StorageAccount.ContainerName + "/",
	}
}

// GenerateUploadUrl generated pre-signed upload url
func (b *BlobService) GenerateUploadUrl(ctx *requestcontext.RequestContext, storagePath string, fileName *string, expiryTime time.Time) (string, *api_error.ApiError) {
	uploadUrl, err := b.generateSASUrl(ctx, storagePath, fileName, expiryTime, true)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

// GenerateDownloadUrl generated pre-signed download url
func (b *BlobService) GenerateDownloadUrl(ctx *requestcontext.RequestContext, storagePath string, fileName string, expiryTime time.Time) (string, *api_error.ApiError) {
	uploadUrl, err := b.generateSASUrl(ctx, storagePath, &fileName, expiryTime, false)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

// GenerateSASToken generated SAS token based url for pre-signed access to blob
func (b *BlobService) generateSASUrl(ctx *requestcontext.RequestContext, storagePath string, fileName *string, expiryTime time.Time, upload bool) (string, *api_error.ApiError) {
	var blobPermissions sas.BlobPermissions
	if upload {
		blobPermissions = sas.BlobPermissions{Create: true, Write: true, Tag: false}
	} else {
		blobPermissions = sas.BlobPermissions{Read: true, Tag: false}
	}

	var correlationId = uuid.New().String()
	ctx.Logger().Info("Generating SASToken with cid=" + correlationId)
	blobSignatureValues := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    expiryTime,
		Permissions:   blobPermissions.String(),
		ContainerName: b.containerName,
		BlobName:      storagePath,
		CorrelationID: correlationId,
	}
	if fileName != nil {
		blobSignatureValues.ContentDisposition = "attachment; filename=" + *fileName
	}
	userDelegationCred, credentialErr := b.getUserDelegationCredential(ctx, expiryTime)
	if credentialErr != nil {
		return "", credentialErr
	}
	sasParams, err := blobSignatureValues.SignWithUserDelegation(userDelegationCred)
	if err != nil {
		ctx.Logger().Error("Error in GenerateSASToken", zap.Error(err))
		return "", api_error.InternalServerError
	}
	preSignedUrl := b.containerUrl + storagePath + "?" + sasParams.Encode()
	return preSignedUrl, nil
}

// getUserDelegationCredential generate UserDelegationCredential
func (b *BlobService) getUserDelegationCredential(ctx *requestcontext.RequestContext, expiryTime time.Time) (*service.UserDelegationCredential, *api_error.ApiError) {
	expiryTimeStr := expiryTime.Format(time.RFC3339)
	startTimeSTr := time.Now().UTC().Format(time.RFC3339)
	keyInfo := service.KeyInfo{
		Expiry: &expiryTimeStr,
		Start:  &startTimeSTr,
	}
	// If GetUserDelegationCredential is having performance impacts, we will then have to create the credential with higher expiry
	// and keep checking if credential expired and create new one if required
	userDelegationCred, credentialErr := b.blobServiceClient.GetUserDelegationCredential(ctx.DbCtx(), keyInfo, nil)
	if credentialErr != nil {
		ctx.Logger().Error("Error in GetUserDelegationCredential", zap.Error(credentialErr))
		return nil, api_error.InternalServerError
	}
	return userDelegationCred, nil
}
