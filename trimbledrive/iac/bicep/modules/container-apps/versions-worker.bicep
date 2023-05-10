param location string
param tags object
param containerRegistryName string
param containerAppsEnvName string
param containerName string
param imageTag string
param mongoDatabaseName string
param kafkaTopic string //'db.tdrive-rkumar1.latest'
param kafkaGroupId string // 'db.tdrive-rkumar1.latest-vworker'
param kafkaControlTopic string // db.tdrive-rkumar1.control
param kafkaDLT string //'db.tdrive-rkumar1.latest.dlt'

@secure()
param mongoUri string

@secure()
param kafkaBrokers string

@secure()
param kafkaAPIKey string

@secure()
param kafkaAPISecret string

resource containerRegistry 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' existing = {
  name: containerRegistryName
}

resource cappsEnv 'Microsoft.App/managedEnvironments@2022-06-01-preview' existing = {
  name: containerAppsEnvName
}

resource versionsWorker 'Microsoft.App/containerApps@2022-06-01-preview' = {
  name: containerName
  location: location
  tags: tags
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
                name: 'GO_KAFKA_DEADLETTERTOPIC'
                value: kafkaDLT
            }
            {
                name: 'GO_CONTROL_TOPIC'
                value: kafkaControlTopic
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
        // rules: [
        //   {
        //     name: 'kafka-rule'
        //     custom: {
        //       type: 'kafka'
        //       metadata: {
        //         bootstrapServers: 'kafka.svc:9092'
        //         consumerGroup: 'my-group'
        //         lagThreshold: '5'
        //         tls: 'enable'
        //         sasl: 'plaintext'
        //       }
        //       auth: [
        //         {
        //           secretRef: 'sb-root-connectionstring'
        //           triggerParameter: 'connection'
        //         }
        //       ]
        //     }
        //   }
        // ]
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
            name: 'container-registry-password'
            value: containerRegistry.listCredentials().passwords[0].value
        }
      ]
    }
  }
}
