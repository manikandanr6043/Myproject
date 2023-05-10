#!/bin/bash
set -ex

STACK_SUFFIX="${1:-$USER$USERNAME}"

RESOURCE_GROUP_NAME="${2:-RG-RTMZ-TDRV-DEVELOPMENT}"
# LOCATION="westeurope"
DEPLOYMENT_VERSION=$STACK_SUFFIX-$(date +%s)

## Register Providers
az provider register --wait --namespace Microsoft.App
az provider register --wait --namespace Microsoft.ContainerService
az provider register --wait --namespace Microsoft.Cdn

az extension add --name containerapp --upgrade

az account set --name trimble-connect-platform-sandbox

# az group create -n "$RESOURCE_GROUP_NAME" --location "$LOCATION" \
#   --tags product=CDE environment=development primary-contact-email=alexey.otavin@trimble.com primary-contact-name="Alexey Otavin" application=cde:trimble-drive team-contact-email=trimble-drive-dev-ug@trimble.com backup=no

## Deploy Template
az deployment group create \
    --resource-group "$RESOURCE_GROUP_NAME" \
    --name "$DEPLOYMENT_VERSION" \
    --no-prompt \
    --template-file ./iac/bicep/main.bicep \
    --parameters stackSuffix="$STACK_SUFFIX" \
    --query properties.outputs.fqdn

echo "Please wait a few minutes until endpoint is established..."
az deployment group wait --created -n "$DEPLOYMENT_VERSION" -g "$RESOURCE_GROUP_NAME"
