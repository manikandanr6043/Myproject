param location string
param tags object
param storageAccountName string
param blobContainerName string
param fileShareName string
param managedIdentityName string

resource storageAccount 'Microsoft.Storage/storageAccounts@2022-05-01' = {
  name: storageAccountName
  location: location
  tags: tags
  kind: 'StorageV2'
  sku: {
    name: 'Standard_LRS'
  }
  properties: {
    dnsEndpointType: 'Standard'
    defaultToOAuthAuthentication: false
    publicNetworkAccess: 'Enabled'
    allowCrossTenantReplication: false
    minimumTlsVersion: 'TLS1_2'
    allowBlobPublicAccess: true
    allowSharedKeyAccess: true
    largeFileSharesState: 'Enabled'
    networkAcls: {
      resourceAccessRules: []
      bypass: 'AzureServices'
      virtualNetworkRules: []
      ipRules: [
        {
          value: '182.65.97.113'
          action: 'Allow'
        }
      ]
      defaultAction: 'Allow'
    }
    supportsHttpsTrafficOnly: true
    encryption: {
      requireInfrastructureEncryption: false
      services: {
        file: {
          keyType: 'Account'
          enabled: true
        }
        blob: {
          keyType: 'Account'
          enabled: true
        }
      }
      keySource: 'Microsoft.Storage'
    }
    accessTier: 'Hot'
  }

  resource blobService 'blobServices' = {
    name: 'default'

    properties: {
      cors: {
        corsRules: [
            {
                allowedOrigins: ['*']
                allowedMethods: ['GET','PUT','OPTIONS']
                maxAgeInSeconds: 0
                exposedHeaders: ['x-ms-*']
                allowedHeaders: ['x-ms-*']
            }
        ]
      }
      deleteRetentionPolicy: {
          allowPermanentDelete: false
          enabled: true
          days: 7
      }
      isVersioningEnabled: true
      changeFeed: {
          retentionInDays: 14
          enabled: true
      }
      restorePolicy: {
          enabled: false
      }
      containerDeleteRetentionPolicy: {
          enabled: false
      }
    }

    resource blobContainer 'containers' = {
      name: blobContainerName
      properties: {
        immutableStorageWithVersioning: {
          enabled: false
        }
        publicAccess: 'None'
      }
    }
  }

  resource fileService 'fileServices' = {
    name: 'default'
    properties: {
      shareDeleteRetentionPolicy: {
        enabled: false
        days: 0
      }
    }

    resource fileShares 'shares' = {
      name: fileShareName
      properties: {
        accessTier: 'TransactionOptimized'
        shareQuota: 1024
        enabledProtocols: 'SMB'
      }
    }
  }

  // TODO: does not work as adding a new rule for the TMP stack
  resource symbolicname 'managementPolicies' = {
    name: 'default'
    properties: {
      policy: {
        rules: [
          {
            enabled: true
            name: 'ThumbBlobExpiry-${blobContainerName}'
            type: 'Lifecycle'
            definition: {
              actions: {
                version: {
                  delete: {
                    daysAfterCreationGreaterThan: 14
                  }
                }
                baseBlob: {
                  delete: {
                    daysAfterCreationGreaterThan: 14
                  }
                }
              }
              filters: {
                blobTypes: [
                  'blockBlob'
                ]
                prefixMatch: [
                  '${blobContainerName}/tmp'
                ]
              }
            }
          }
        ]
      }
    }
  }
}

resource managedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' existing = {
  name: managedIdentityName
}

@description('This is the built-in Contributor role. See https://docs.microsoft.com/azure/role-based-access-control/built-in-roles#contributor')
resource contributorRoleDefinition 'Microsoft.Authorization/roleDefinitions@2018-01-01-preview' existing = {
  scope: subscription()
  name: 'b24988ac-6180-42a0-ab88-20f7382dd24c'
}



@description('This is the built-in Storage Blob Data Contributor role. See https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#storage-blob-data-contributor')
resource storageBlobDataContributorRoleDefinition 'Microsoft.Authorization/roleDefinitions@2018-01-01-preview' existing = {
  scope: subscription()
  name: 'ba92f5b4-2d11-453d-a403-e96b0029c9fe'
}

@description('This is the built-in Storage Blob Delegator role. See https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#storage-blob-delegator')
resource storageBlobDelegatorRoleDefinition 'Microsoft.Authorization/roleDefinitions@2018-01-01-preview' existing = {
  scope: subscription()
  name: 'db58b8e5-c6ad-4a2a-8342-4190687cbf4a'
}

resource roleAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: storageAccount
  name: guid(storageAccount.id, managedIdentity.id, contributorRoleDefinition.id)
  properties: {
    roleDefinitionId: contributorRoleDefinition.id
    principalId: managedIdentity.properties.principalId
    principalType: 'ServicePrincipal'
  }
}

resource role2Assignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: storageAccount
  name: guid(storageAccount.id, managedIdentity.id, storageBlobDataContributorRoleDefinition.id)
  properties: {
    roleDefinitionId: storageBlobDataContributorRoleDefinition.id
    principalId: managedIdentity.properties.principalId
    principalType: 'ServicePrincipal'
  }
}

resource role3Assignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
  scope: storageAccount
  name: guid(storageAccount.id, managedIdentity.id, storageBlobDelegatorRoleDefinition.id)
  properties: {
    roleDefinitionId: storageBlobDelegatorRoleDefinition.id
    principalId: managedIdentity.properties.principalId
    principalType: 'ServicePrincipal'
  }
}

output name string = storageAccount.name
