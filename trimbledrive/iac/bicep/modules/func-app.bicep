param location string
param tags object
param hostingPlanName string
param functionAppName string
param storageAccountName string
param blobContainerName string
param contentCdnEndpoint string
param appInsightsInstrumentationKey string
param storageSystemTopicName string
param cosmosDatabaseName string
param cosmosAccountName string
param thumbTopic string
param commitTopic string
param managedIdentityName string
param keyVaultName string

@secure()
param kafkaBrokers string

resource managedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' existing = {
  name: managedIdentityName
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2022-05-01' existing = {
  name: storageAccountName
}

resource cosmosAccount 'Microsoft.DocumentDB/databaseAccounts@2022-08-15' existing = {
  name: cosmosAccountName
}

resource hostingPlan 'Microsoft.Web/serverfarms@2022-03-01' = {
  name: hostingPlanName
  location: location
  tags: tags
  sku: {
    name: 'EP2'
    tier: 'ElasticPremium'
    family: 'EP'
  }
  kind: 'elastic'
  properties: {
    elasticScaleEnabled: true
    maximumElasticWorkerCount: 20
    reserved: true
    isXenon: false
    hyperV: false
    targetWorkerCount: 0
    targetWorkerSizeId: 0
    zoneRedundant: false
  }
}

resource functionApp 'Microsoft.Web/sites@2022-03-01' = {
  name: functionAppName
  location: location
  tags: tags
  kind: 'functionapp,linux'
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${managedIdentity.id}': {}
    }
  }
  dependsOn: [
    storageAccount
  ]
  properties: {
    serverFarmId: hostingPlan.id
    clientAffinityEnabled: false
    reserved: true
    isXenon: false
    hyperV: false
    siteConfig: {
      numberOfWorkers: 1
      linuxFxVersion: 'Python|3.8'
      acrUseManagedIdentityCreds: false
      alwaysOn: false
      http20Enabled: true
      functionAppScaleLimit: 0
      minimumElasticInstanceCount: 1
      appSettings: [
        {
          name: 'APPINSIGHTS_INSTRUMENTATIONKEY'
          value: appInsightsInstrumentationKey
        }
        {
          name: 'APPLICATIONINSIGHTS_CONNECTION_STRING'
          value: 'InstrumentationKey=${appInsightsInstrumentationKey}'
        }
        {
          name: 'AzureWebJobsStorage'
          value: 'DefaultEndpointsProtocol=https;AccountName=${storageAccountName};EndpointSuffix=${environment().suffixes.storage};AccountKey=${storageAccount.listKeys().keys[0].value}'
        }
        {
          name: 'WEBSITE_CONTENTAZUREFILECONNECTIONSTRING'
          value: 'DefaultEndpointsProtocol=https;AccountName=${storageAccountName};EndpointSuffix=${environment().suffixes.storage};AccountKey=${storageAccount.listKeys().keys[0].value}'
        }
        {
          name: 'WEBSITE_CONTENTSHARE'
          value: '${toLower(functionAppName)}${uniqueString(resourceGroup().id, deployment().name)}'
        }
        {
          name: 'FUNCTIONS_EXTENSION_VERSION'
          value: '~4'
        }
        {
          name: 'FUNCTIONS_WORKER_RUNTIME'
          value: 'python'
        }
        {
          name: 'WEBSITE_NODE_DEFAULT_VERSION'
          value: '~14'
        }
        {
          name: 'WEBSITE_WEBDEPLOY_USE_SCM'
          value: 'true'
        }
        {
          name: 'BootstrapServer'
          value: kafkaBrokers
        }
        {
          name: 'ConfluentCloudUsername'
          value: '@Microsoft.KeyVault(VaultName=${keyVaultName};SecretName=CONFLUENT-CLOUD-API-KEY)'
        }
        {
          name: 'ConfluentCloudPassword'
          value: '@Microsoft.KeyVault(VaultName=${keyVaultName};SecretName=CONFLUENT-CLOUD-API-SECRET)'
        }
        {
          name: 'StorageAccountUrl'
          value: contentCdnEndpoint
        }
        {
          name: 'BlobContainer'
          value: blobContainerName
        }
        {
          name: 'CosmosAccountName'
          value: cosmosAccountName
        }
        {
          name: 'CosmosDbName'
          value: cosmosDatabaseName
        }
        {
          name: 'CosmosKey'
          value: cosmosAccount.listKeys().primaryMasterKey
        }
        {
          name: 'CosmosConnectionString'
          value: cosmosAccount.listConnectionStrings().connectionStrings[0].connectionString
        }
        {
          name: 'ThumbTopic'
          value: thumbTopic
        }
        {
          name: 'CommitTopic'
          value: commitTopic
        }
      ]
    }
    hostNamesDisabled: false
    containerSize: 0
    dailyMemoryTimeQuota: 0
    httpsOnly: false
    redundancyMode: 'None'
    storageAccountRequired: false
    keyVaultReferenceIdentity: managedIdentity.id
  }
}

resource systemTopic 'Microsoft.EventGrid/systemTopics@2021-12-01' = {
  name: storageSystemTopicName
  location: location
  tags: tags
  properties: {
    source: storageAccount.id
    topicType: 'Microsoft.Storage.StorageAccounts'
  }
}
