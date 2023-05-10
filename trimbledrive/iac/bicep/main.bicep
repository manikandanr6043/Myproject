@description('Azure Location/Region')
param location string = resourceGroup().location

param env string = 'DEVELOPMENT'
param stackSuffix string  // int, tmp551, aotavin

param keyVaultName string = 'KV-RTMZ-TDRV-${toUpper(env)}' //KV-RTMZ-TDRV-DEVELOPMENT
param identityName string = 'DevContainerApp'
param acrName string = 'crrtmztdrv${toLower(env)}' // crrtmztdrvdevelopment
param containerAppsEnvName string = 'cae-rtmz-tdrv-${toLower(env)}' // cae-rtmz-tdrv-development
param logAnalyticsWorkspaceName string = 'LOG-RTMZ-TDRV-${toUpper(env)}' // LOG-RTMZ-TDRV-DEVELOPMENT
param appInsightsName string = 'appi-rtmz-tdrv-${toUpper(env)}' // tdrive-upload-fun-app-int
param cosmosAccountName string = 'cdbrtmztdrvdevcosmos' // cdbrtmztdrvdevcosmos
param storageAccountName string = 'srgrtmztdrvdevcontent' // srgrtmztdrvdevcontent
param fileShareNameCoverage string = 'coverage'
param hostingPlanName string = 'tdrive-dev-premium-plan'  // tdrive-dev-premium-plan
param frontDoorName string = 'cdn-rtmz-tdrv-tcdp'  // cdn-rtmz-tdrv-tcdp
param storageSystemTopicName string = '${storageAccountName}-863cdc00-f452-45b8-bb71-93587e219684'

param cosmosDatabaseName string = 'tdrive-${stackSuffix}'
param mongoDatabaseName string = 'tdrive-${stackSuffix}'
param blobContainerName string = 'tdrive-${stackSuffix}' // tdrive-int, tdrive-aotavin
param apiServiceName string = 'tdrive-api-${stackSuffix}'  // tdrive-api-int, tdrive-api-aotavin
param apiServiceImageTag string = '${apiServiceName}:latest'
param versionWorkerName string = 'tdrive-vworker-${stackSuffix}' // tdrive-vworker-int, tdrive-vworker-aotavin
param versionWorkerImageTag string = '${versionWorkerName}:latest'
param commitProcessorName string = 'tdrive-commit-${stackSuffix}' // tdrive-commit-int, tdrive-commit-aotavin
param commitProcessorImageTag string = '${commitProcessorName}:latest'
param functionAppName string = 'tdrive-upload-function-app-${stackSuffix}' // tdrive-upload-fun-app-int, tdrive-upload-function-app-aotavin

param kafkaVWorkerTopic string = 'db.tdrive-${stackSuffix}.latest' // 'db.tdrive-rkumar1.latest'
param kafkaVWorkerGroupId string = 'db.tdrive-${stackSuffix}.latest-vworker' // 'db.tdrive-aotavin.latest-vworker'
param kafkaVWorkerDLT string = 'db.tdrive-${stackSuffix}.latest.dlt' // 'db.tdrive-aotavin.latest.dlt'
param kafkaCommitProcessorTopic string = 'trimble.tdrive.commit_processor-${stackSuffix}' // 'trimble.tdrive.commit_processor-aotavin'
param kafkaCommitProcessorGroupId string = 'trimble.tdrive.commit_processor--${stackSuffix}-commit-processor' // 'trimble.tdrive.commit_processor-aotavin-commit-processor'
param kafkaControlTopic string = 'db.tdrive--${stackSuffix}.control' // 'db.tdrive-aotavin.control'
param kafkaThumbTopic string = 'trimble.tdrive.thumb_processor-${stackSuffix}' // 'trimble.tdrive.thumb_processor-aotavin'


@description('By default all apps scale down to 0. This setting ensures a minimum instance count of 1 for API.')
param keepApiServiceUp bool = true

var tags = {
  product: 'CDE'
  application: 'cde:trimble-drive'
  environment: env
  'primary-contact-email': 'alexey.otavin@trimble.com'
  'primary-contact-name': 'Alexey Otavin'
  'team-contact-email': 'trimble-drive-dev-ug@trimble.com'
  backup: 'no'
}

resource keyVault 'Microsoft.KeyVault/vaults@2019-09-01' existing = {
  name: keyVaultName
  // scope: resourceGroup('rg-contoso')   - if key vault is in a different resource group
}

module containerRegistryModule 'modules/registry.bicep' = {
  name: '${deployment().name}--registry'
  params: {
    location: location
    tags: tags
    acrName: acrName
  }
}

module identityModule 'modules/identity.bicep' = {
  name: '${deployment().name}--identity'
  params: {
    location: location
    tags: tags
    identityName: identityName
    keyVaultName: keyVaultName
  }
}

module storageModule 'modules/storage.bicep' = {
  name: '${deployment().name}--storage'
  params: {
    location: location
    tags: tags
    managedIdentityName: identityName
    storageAccountName: storageAccountName
    blobContainerName: blobContainerName
    fileShareName: fileShareNameCoverage
  }
}
module containerAppsEnvModule 'modules/capps-env.bicep' = {
  name: '${deployment().name}--containerAppsEnv'
  dependsOn: [
    storageModule
  ]
  params: {
    location: location
    tags: tags
    containerAppsEnvName: containerAppsEnvName
    logAnalyticsWorkspaceName: logAnalyticsWorkspaceName
    appInsightsName: appInsightsName
    storageAccountName: storageModule.outputs.name
  }
}

module cosmosModule 'modules/cosmos.bicep' = {
  name: '${deployment().name}--cosmos'
  params: {
    location: location
    tags: tags
    cosmosAccountName: cosmosAccountName
    cosmosDatabaseName: cosmosDatabaseName
    //managedIdentityName: identityName
  }
}

module versionsWorkerModule 'modules/container-apps/versions-worker.bicep' = {
  name: '${deployment().name}--versions-worker'
  dependsOn: [
    containerRegistryModule
    containerAppsEnvModule
  ]
  params: {
    location: location
    tags: tags
    containerRegistryName: containerRegistryModule.outputs.name
    containerAppsEnvName: containerAppsEnvModule.outputs.cappsEnvName
    containerName: versionWorkerName
    imageTag: versionWorkerImageTag
    mongoUri: keyVault.getSecret('MONGO-URI')
    mongoDatabaseName: mongoDatabaseName
    kafkaBrokers: keyVault.getSecret('CONFLUENT-CLOUD-BOOTSTRAP-SERVER')
    kafkaAPIKey: keyVault.getSecret('CONFLUENT-CLOUD-API-KEY')
    kafkaAPISecret: keyVault.getSecret('CONFLUENT-CLOUD-API-SECRET')
    kafkaTopic:kafkaVWorkerTopic
    kafkaGroupId: kafkaVWorkerGroupId
    kafkaDLT:kafkaVWorkerDLT
    kafkaControlTopic: kafkaControlTopic
  }
}

module commitProcessorModule 'modules/container-apps/commit-processor.bicep' = {
  name: '${deployment().name}--commit-processor'
  dependsOn: [
    containerRegistryModule
    containerAppsEnvModule
  ]
  params: {
    location: location
    tags: tags
    managedIdentityName: identityName
    containerRegistryName: containerRegistryModule.outputs.name
    containerAppsEnvName: containerAppsEnvModule.outputs.cappsEnvName
    containerName: commitProcessorName
    imageTag: commitProcessorImageTag
    mongoUri: keyVault.getSecret('MONGO-URI')
    mongoDatabaseName: mongoDatabaseName
    cosmosDatabaseName: cosmosDatabaseName
    kafkaBrokers: keyVault.getSecret('CONFLUENT-CLOUD-BOOTSTRAP-SERVER')
    kafkaAPIKey: keyVault.getSecret('CONFLUENT-CLOUD-API-KEY')
    kafkaAPISecret: keyVault.getSecret('CONFLUENT-CLOUD-API-SECRET')
    kafkaTopic: kafkaCommitProcessorTopic
    kafkaGroupId: kafkaCommitProcessorGroupId
    kafkaControlTopic: kafkaControlTopic
    azureClientId: keyVault.getSecret('CLIENT-ID')
  }
}

module apiModule 'modules/container-apps/api-service.bicep' = {
  name: '${deployment().name}--api-service'
  dependsOn: [
    containerRegistryModule
    containerAppsEnvModule
  ]
  params: {
    location: location
    tags: tags
    managedIdentityName: identityName
    containerRegistryName: containerRegistryModule.outputs.name
    containerAppsEnvName: containerAppsEnvModule.outputs.cappsEnvName
    containerName: apiServiceName
    imageTag: apiServiceImageTag
    minReplicas: keepApiServiceUp ? 1 : 0
    mongoUri: keyVault.getSecret('MONGO-URI')
    mongoDatabaseName: mongoDatabaseName
    cosmosDatabaseName: cosmosDatabaseName
    cosmosUrl: cosmosModule.outputs.documentEndpoint
    blobContainerName: blobContainerName
    storageAccountName: storageAccountName
    azureClientId: keyVault.getSecret('CLIENT-ID')
    frontdoorName: frontDoorName
  }
}

module funcApp 'modules/func-app.bicep' = {
  name: '${deployment().name}--funcapp'
  dependsOn: [
    containerAppsEnvModule
  ]
  params: {
    location: location
    tags: tags
    managedIdentityName: identityName
    storageSystemTopicName: storageSystemTopicName
    hostingPlanName: hostingPlanName
    functionAppName: functionAppName
    appInsightsInstrumentationKey: containerAppsEnvModule.outputs.appInsightsInstrumentationKey
    storageAccountName: storageAccountName
    blobContainerName: blobContainerName
    contentCdnEndpoint: apiModule.outputs.contentCdnEndpoint
    kafkaBrokers: keyVault.getSecret('CONFLUENT-CLOUD-BOOTSTRAP-SERVER')
    thumbTopic: kafkaThumbTopic
    commitTopic: kafkaCommitProcessorTopic
    cosmosDatabaseName: cosmosDatabaseName
    cosmosAccountName: cosmosAccountName
    keyVaultName: keyVaultName
  }
}

output fqdn string = apiModule.outputs.apiFqdn
output apiContainerName string = apiServiceName
output vworkerContainerName string = versionWorkerName
output commitContainerName string = commitProcessorName
