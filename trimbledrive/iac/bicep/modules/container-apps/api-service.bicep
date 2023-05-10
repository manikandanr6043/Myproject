param location string
param tags object
param containerRegistryName string
param containerAppsEnvName string
param containerName string
param imageTag string
param mongoDatabaseName string
param cosmosDatabaseName string
param cosmosUrl string
param blobContainerName string
param storageAccountName string
param frontdoorName string
param managedIdentityName string

@secure()
param mongoUri string

@secure()
param azureClientId string

param minReplicas int = 0

@description('The name of the SKU to use when creating the Front Door profile. If you use Private Link this must be set to `Premium_AzureFrontDoor`.')
@allowed([
  'Standard_AzureFrontDoor'
  'Premium_AzureFrontDoor'
])
param skuName string = 'Standard_AzureFrontDoor'

@allowed([
  'Detection'
  'Prevention'
])
@description('The mode that the WAF should be deployed using. In \'Prevention\' mode, the WAF will block requests it detects as malicious. In \'Detection\' mode, the WAF will not block requests and will simply log the request.')
param wafMode string = 'Prevention'

// @description('The list of managed rule sets to configure on the WAF.')
// param wafManagedRuleSets array = [
//   {
//     ruleSetType: 'Microsoft_DefaultRuleSet'
//     ruleSetVersion: '2.0'
//     ruleSetAction: 'Block'
//   }
//   {
//     ruleSetType: 'Microsoft_BotManagerRuleSet'
//     ruleSetVersion: '1.0'
//   }
// ]

resource containerRegistry 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' existing = {
  name: containerRegistryName
}

resource cappsEnv 'Microsoft.App/managedEnvironments@2022-06-01-preview' existing = {
  name: containerAppsEnvName
}

resource managedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' existing = {
  name: managedIdentityName
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2022-05-01' existing = {
  name: storageAccountName
}

resource apiService 'Microsoft.App/containerApps@2022-06-01-preview' = {
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
                name: 'GO_STORAGEACCOUNT_CONTAINERNAME'
                value: blobContainerName
            }
            {
                name: 'GO_STORAGEACCOUNT_URL'
                value: 'https://${frontDoorProfile::contentEndpoint.properties.hostName}/'
            }
            {
                name: 'GO_STORAGEACCOUNT_NAME'
                value: storageAccountName
            }
            {
                name: 'GO_COSMOS_DBNAME'
                value: cosmosDatabaseName
            }
            {
                name: 'GO_COSMOS_URL'
                value: cosmosUrl
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
          probes: [
            {
              type: 'startup'
              httpGet: {
                path: '/v1/app/health/ping'
                port: 8080
              }
              failureThreshold: 3
              initialDelaySeconds: 10
              periodSeconds: 10
              successThreshold: 1
              timeoutSeconds: 2
            }
            {
              type: 'liveness'
              httpGet: {
                path: '/v1/app/health/ping'
                port: 8080
              }
              failureThreshold: 3
              initialDelaySeconds: 10
              periodSeconds: 10
              successThreshold: 1
              timeoutSeconds: 2
            }
          ]
        }
      ]
      scale: {
        minReplicas: minReplicas
        maxReplicas: 10
        rules: [
            {
                name: 'azure-http-rule'
                http: {
                    metadata: {
                        concurrentRequests: '50'
                    }
                }
            }
        ]
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
      ingress: {
        external: true
        targetPort: 8080
        allowInsecure: false
      }
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

resource contentWafPolicy 'Microsoft.Network/frontdoorwebapplicationfirewallpolicies@2022-05-01' = {
  name: 'wafrtmztdrvcontent'
  location: 'Global'
  tags: tags
  sku: {
    name: skuName
  }
  properties: {
    policySettings: {
      enabledState: 'Enabled'
      mode: wafMode
      requestBodyCheck: 'Enabled'
    }
    // managedRules: {
    //   managedRuleSets: wafManagedRuleSets
    // }
  }
}

resource apiWafPolicy 'Microsoft.Network/frontdoorwebapplicationfirewallpolicies@2022-05-01' = {
  name: 'wafrtmz${replace(containerName,'-','')}'
  location: 'Global'
  tags: tags
  sku: {
    name: skuName
  }
  properties: {
    policySettings: {
      enabledState: 'Enabled'
      mode: wafMode
      requestBodyCheck: 'Enabled'
    }
    // managedRules: {
    //   managedRuleSets: wafManagedRuleSets
    // }
  }
}

resource frontDoorProfile 'Microsoft.Cdn/profiles@2021-06-01' = {
  name: frontdoorName
  location: 'Global'
  tags: tags
  sku: {
    name: skuName
  }
  properties: {
    originResponseTimeoutSeconds: 60
  }

  resource contentOriginGroup 'originGroups' = {
    name: 'content-group'
    properties: {
      loadBalancingSettings: {
        sampleSize: 4
        successfulSamplesRequired: 3
        additionalLatencyInMilliseconds: 50
      }
      sessionAffinityState: 'Disabled'
      healthProbeSettings: {
        probePath: '/'
        probeRequestType: 'HEAD'
        probeProtocol: 'Http'
        probeIntervalInSeconds: 100
      }
    }

    resource contentOrigin 'origins' = {
      name: 'content'
      properties: {
        hostName: replace(replace(storageAccount.properties.primaryEndpoints.blob, 'https://', ''), '/', '')
        httpPort: 80
        httpsPort: 443
        originHostHeader: replace(replace(storageAccount.properties.primaryEndpoints.blob, 'https://', ''), '/', '')
        priority: 1
        weight: 1000
        enabledState: 'Enabled'
        enforceCertificateNameCheck: true
      }
    }
  }

  resource apiOriginGroup 'originGroups' = {
    name: containerName
    properties: {
      loadBalancingSettings: {
        sampleSize: 4
        successfulSamplesRequired: 3
        additionalLatencyInMilliseconds: 50
      }
      sessionAffinityState: 'Disabled'
      healthProbeSettings: {
        probePath: '/v1/app/health/ping'
        probeRequestType: 'GET'
        probeProtocol: 'Https'
        probeIntervalInSeconds: 100
      }
    }

    resource apiOrigin 'origins' = {
      name: containerName
      properties: {
        hostName: replace(replace(apiService.properties.configuration.ingress.fqdn, 'https://', ''), '/', '')
        httpPort: 80
        httpsPort: 443
        originHostHeader: replace(replace(apiService.properties.configuration.ingress.fqdn, 'https://', ''), '/', '')
        priority: 1
        weight: 1000
        enabledState: 'Enabled'
        enforceCertificateNameCheck: true
      }
    }
  }

  resource contentEndpoint 'afdEndpoints' = {
    name: 'content'
    location: 'Global'
    properties: {
      enabledState: 'Enabled'
    }
    resource contentRoute 'routes' = {
      name: 'content-route'
      dependsOn: [
        frontDoorProfile::contentOriginGroup::contentOrigin // This explicit dependency is required to ensure that the origin group is not empty when the route is created.
      ]
      properties: {
        originGroup: {
          id: contentOriginGroup.id
        }
        supportedProtocols: ['Https']
        patternsToMatch: ['/*']
        forwardingProtocol: 'HttpsOnly'
        linkToDefaultDomain: 'Enabled'
        httpsRedirect: 'Disabled'
      }
    }
  }

  resource apiEndpoint 'afdEndpoints' = {
    name: containerName
    location: 'Global'
    properties: {
      enabledState: 'Enabled'
    }
    resource apiRoute 'routes' = {
      name: containerName
      dependsOn: [
        frontDoorProfile::apiOriginGroup::apiOrigin // This explicit dependency is required to ensure that the origin group is not empty when the route is created.
      ]
      properties: {
        originGroup: {
          id: apiOriginGroup.id
        }
        supportedProtocols: ['Https']
        patternsToMatch: ['/*']
        forwardingProtocol: 'HttpsOnly'
        linkToDefaultDomain: 'Enabled'
        httpsRedirect: 'Disabled'
        cacheConfiguration: {
          compressionSettings: {
              isCompressionEnabled: true
              contentTypesToCompress: [
                'application/eot'
                'application/font'
                'application/font-sfnt'
                'application/javascript'
                'application/json'
                'application/opentype'
                'application/otf'
                'application/pkcs7-mime'
                'application/truetype'
                'application/ttf'
                'application/vnd.ms-fontobject'
                'application/xhtml+xml'
                'application/xml'
                'application/xml+rss'
                'application/x-font-opentype'
                'application/x-font-truetype'
                'application/x-font-ttf'
                'application/x-httpd-cgi'
                'application/x-javascript'
                'application/x-mpegurl'
                'application/x-opentype'
                'application/x-otf'
                'application/x-perl'
                'application/x-ttf'
                'font/eot'
                'font/ttf'
                'font/otf'
                'font/opentype'
                'image/svg+xml'
                'text/css'
                'text/csv'
                'text/html'
                'text/javascript'
                'text/js'
                'text/plain'
                'text/richtext'
                'text/tab-separated-values'
                'text/xml'
                'text/x-script'
                'text/x-component'
                'text/x-java-source'
              ]
          }
          queryStringCachingBehavior: 'UseQueryString'
        }
      }
    }
  }

  resource contentSecurityPolicies 'securitypolicies' = {
    name: 'content'
    properties: {
      parameters: {
        wafPolicy: {
          id: contentWafPolicy.id
        }
        associations: [
          {
            domains: [
              {
                id: contentEndpoint.id
              }
            ]
            patternsToMatch: [
              '/*'
            ]
          }
        ]
        type: 'WebApplicationFirewall'
      }
    }
  }

  resource apiSecurityPolicies 'securitypolicies' = {
    name: containerName
    properties: {
      parameters: {
        wafPolicy: {
          id: apiWafPolicy.id
        }
        associations: [
          {
            domains: [
              {
                id: apiEndpoint.id
              }
            ]
            patternsToMatch: [
              '/*'
            ]
          }
        ]
        type: 'WebApplicationFirewall'
      }
    }
  }
}

output subdomain string = apiService.name
output containerFqdn string = apiService.properties.configuration.ingress.fqdn
output apiFqdn string = 'https://${frontDoorProfile::apiEndpoint.properties.hostName}'
output contentCdnEndpoint string = 'https://${frontDoorProfile::contentEndpoint.properties.hostName}/'

