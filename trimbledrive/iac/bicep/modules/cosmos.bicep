param location string
param tags object
param cosmosAccountName string
param cosmosDatabaseName string
//param managedIdentityName string

resource cosmosAccount 'Microsoft.DocumentDB/databaseAccounts@2022-08-15' = {
  name: cosmosAccountName
  kind: 'GlobalDocumentDB'
  location: location
  tags: tags
  properties: {
    consistencyPolicy: {
      defaultConsistencyLevel: 'Session'
    }
    locations: [
      {
        locationName: location
      }
    ]
    databaseAccountOfferType: 'Standard'
  }

  resource cosmosDatabase 'sqlDatabases' = {
    name: cosmosDatabaseName
    properties: {
      resource: {
        id: cosmosDatabaseName
      }
    }

    resource cosmosContainer 'containers' = {
      name: 'file_upload'
      properties: {
        resource: {
          id: 'file_upload'
          partitionKey: {
            paths: ['/id']
            kind: 'Hash'
          }
          defaultTtl: 1209600
        }
      }
    }
  }
}

// resource managedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' existing = {
//   name: managedIdentityName
// }

// resource roleAssignment 'Microsoft.Authorization/roleAssignments@2020-04-01-preview' = {
//   scope: '/'
//   name: guid(cosmosAccount.id, managedIdentity.id, '00000000-0000-0000-0000-000000000002')
//   properties: {
//     roleDefinitionId: '/subscriptions/${subscription().id}/resourceGroups/${resourceGroup().id}/providers/Microsoft.DocumentDB/databaseAccounts/${cosmosAccount.id}/sqlRoleDefinitions/00000000-0000-0000-0000-000000000002'
//     principalId: managedIdentity.properties.principalId
//     principalType: 'ServicePrincipal'
//   }
// }


output accountName string = cosmosAccount.name
output documentEndpoint string = cosmosAccount.properties.documentEndpoint
