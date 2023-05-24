param passwordLength int = 12
param utcValue string = utcNow()
var uniqueId = uniqueString(utcValue)
var randomPassword = substring(uniqueId, 0, passwordLength)
resource keyVaultSecret 'Microsoft.KeyVault/vaults/secrets@2021-04-01-preview' = {
  name: 'keyvaulttest1234567/mySecret9999'
  properties: {
    value: randomPassword
  }
}
output result string = utcValue
