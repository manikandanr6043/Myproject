param location string
param tags object

@description('Provide a globally unique name of your Azure Container Registry')
param acrName string = 'acr${uniqueString(resourceGroup().id)}'

@description('Provide a tier of your Azure Container Registry.')
param acrSku string = 'Basic'

// resource managedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' existing = {
//   name: 'tdrv-dev-cont-app'
//   scope: resourceGroup('RG-RTMZ-TDRV-INTERGATION')
// }

resource containerRegistry 'Microsoft.ContainerRegistry/registries@2023-01-01-preview' = {
  name: acrName
  location: location
  tags: tags
  sku: {
    name: acrSku
  }
  // identity: {
  //   type: 'UserAssigned'
  //   userAssignedIdentities: {
  //     '${managedIdentity.id}': {}
  //   }
  // }
  properties: {
    adminUserEnabled: true
    policies: {
      quarantinePolicy: {
          status: 'disabled'
      }
      trustPolicy: {
          type: 'Notary'
          status: 'disabled'
      }
      retentionPolicy: {
          days: 7
          status: 'disabled'
      }
      exportPolicy: {
          status: 'enabled'
      }
      azureADAuthenticationAsArmPolicy: {
          status: 'enabled'
      }
      softDeletePolicy: {
          retentionDays: 7
          status: 'disabled'
      }
    }
    encryption: {
      status: 'disabled'
    }
    dataEndpointEnabled: false
    publicNetworkAccess: 'Enabled'
    networkRuleBypassOptions: 'AzureServices'
    zoneRedundancy: 'Disabled'
    anonymousPullEnabled: false
  }
}

output name string = containerRegistry.name

@description('Output the login server property for later use')
output loginServer string = containerRegistry.properties.loginServer
