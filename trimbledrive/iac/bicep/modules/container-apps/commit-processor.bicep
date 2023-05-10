param location string
param tags object
param containerRegistryName string
param containerAppsEnvName string
param containerName string
param imageTag string
param mongoDatabaseName string
param kafkaTopic string //'trimble.tdrive.commit_processor-rkumar1'
param kafkaGroupId string // 'trimble.tdrive.commit_processor-rkumar1-commit-processor'
param kafkaControlTopic string // db.tdrive-rkumar1.control
param cosmosDatabaseName string
param managedIdentityName string

@secure()
param kafkaBrokers string

@secure()
param kafkaAPIKey string

@secure()
param kafkaAPISecret string

@secure()
param azureClientId string

@secure()
param mongoUri string

resource containerRegistry 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' existing = {
  name: containerRegistryName
}

resource cappsEnv 'Microsoft.App/managedEnvironments@2022-06-01-preview' existing = {
  name: containerAppsEnvName
}

resource managedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' existing = {
  name: managedIdentityName
}

resource commitProcessor 'Microsoft.App/containerApps@2022-06-01-preview' = {
  name: containerName
  location: location
  tags: tags
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${managedIdentity.id}': {}
    }
  }
  properties: {
    managedEnvironmentId: cappsEnv.id
    environmentId: cappsEnv.id
    template: {
      containers: [
        {
          name: containerName
          image: '${containerRegistry.properties.loginServer}/${imageTag}'
          env: [
            {
                name: 'GO_ENVIRONMENT'
                value: 'development'
            }
            {
                name: 'GO_MONGO_DBNAME'
                value: mongoDatabaseName
            }
            {
                name: 'GO_MONGO_URI'
                secretRef: 'mongo-uri'
            }
            {
                name: 'GO_KAFKA_BROKERS'
                value: kafkaBrokers
            }
            {
                name: 'GO_KAFKA_APIKEY'
                secretRef: 'confluent-cloud-api-key'
            }
            {
                name: 'GO_KAFKA_APISECRET'
                secretRef: 'confluent-cloud-api-secret'
            }
            {
                name: 'GO_KAFKA_GROUPID'
                value: kafkaGroupId
            }
            {
                name: 'GO_KAFKA_TOPIC'
                value: kafkaTopic
            }
            {
                name: 'GO_COSMOS_DBNAME'
                value: cosmosDatabaseName
            }
            {
                name: 'GO_CONTROL_TOPIC'
                value: kafkaControlTopic
            }
            {
                name: 'AZURE_CLIENT_ID'
                secretRef: 'client-id'
            }
          ]
          resources: {
            cpu:  json('0.5')
            memory: '1Gi'
          }
          volumeMounts: [
              {
                  volumeName: 'coverage'
                  mountPath: '/var/log/coverage'
              }
          ]
        }
      ]
      scale: {
        minReplicas: 1
      }
      volumes: [
        {
            name: 'coverage'
            storageType: 'AzureFile'
            storageName: 'coveragemount'
        }
      ]
    }
    configuration: {
      activeRevisionsMode: 'Single'
      registries: [
        {
            server: containerRegistry.properties.loginServer
            username: containerRegistry.name
            passwordSecretRef: 'container-registry-password'
        }
      ]
      secrets: [
        {
            name: 'mongo-uri'
            value: mongoUri
        }
        {
            name: 'confluent-cloud-api-key'
            value: kafkaAPIKey
        }
        {
            name: 'confluent-cloud-api-secret'
            value: kafkaAPISecret
        }
        {
            name: 'client-id'
            value: azureClientId
        }
        {
            name: 'container-registry-password'
            value: containerRegistry.listCredentials().passwords[0].value
        }
      ]
    }
  }
}
