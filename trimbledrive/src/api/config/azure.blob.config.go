package config

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

// NewBlobConfig Creates Azure blob service client
func NewBlobConfig(tDriveConfig *TDriveConfig) *service.Client {
	// Authenticate using Azure AD
	cred, identityErr := azidentity.NewDefaultAzureCredential(nil)
	if identityErr != nil {
		log.Fatal("Error fetching azure credentials", identityErr)
	}
	serviceUrl := fmt.Sprintf("https://%s.blob.core.windows.net", tDriveConfig.StorageAccount.Name)
	blobClient, err := azblob.NewClient(serviceUrl, cred, nil)
	if err != nil {
		log.Fatal("Error fetching azure blob shared key credentials", err)
	}
	return blobClient.ServiceClient()
}
