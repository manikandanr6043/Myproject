#!/bin/bash
set -ex

STACK_SUFFIX="${1:-$USER$USERNAME}"
FRONT_DOOR_PROFILE_NAME=cdn-rtmz-tdrv-tcdp
REGISTRY_NAME=crrtmztdrvdevelopment
FRONT_DOOR_STACK_NAME=tdrive-api-$STACK_SUFFIX
WAF_POLICY_NAME=wafrtmztdriveapi$STACK_SUFFIX
CONTAINER_APP_NAME=tdrive-api-$STACK_SUFFIX
RESOURCE_GROUP=RG-RTMZ-TDRV-DEVELOPMENT
SUBSCRIPTION_ID=16df61f4-08f2-4bf6-acc3-1d6b5cefde97
COSMOS_ACCOUNT_NAME=cdbrtmztdrvdevcosmos
COSMOS_DATABASE=tdrive-$STACK_SUFFIX
STORAGE_ACCOUNT=srgrtmztdrvdevcontent
STORAGE_CONTAINER=tdrive-$STACK_SUFFIX
FUNCTION_APP=tdrive-upload-function-app-$STACK_SUFFIX
THUMB_EVENT_GRID_SUBSCRIPTION_NAME=ThumbBlobTrigger-$STACK_SUFFIX
ACTIVATOR_EVENT_GRID_SUBSCRIPTION_NAME="activator-blob-trigger-$STACK_SUFFIX"
EVENT_GRID_SOURCE_STORAGE_RESOURCE="/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP/providers/Microsoft.Storage/storageaccounts/$STORAGE_ACCOUNT"
THUMB_PROCESSOR_TOPIC="trimble.tdrive.thumb_processor-$STACK_SUFFIX"
COMMIT_PROCESSOR_TOPIC="trimble.tdrive.commit_processor-$STACK_SUFFIX"
VWORKER_CONTAINER_APP_NAME=tdrive-vworker-$STACK_SUFFIX
KAFKA_CLUSTER_ID="lkc-gnxnjr"
KAFKA_TOPIC_PREFIX=db.tdrive-$STACK_SUFFIX.
COMMIT_PROCESSOR_CONTAINER_APP_NAME=tdrive-commit-$STACK_SUFFIX
COMMIT_PROCESSOR_DLT="trimble.tdrive.commit_processor.dlt-$STACK_SUFFIX"

az account set --name trimble-connect-platform-sandbox
az config set extension.use_dynamic_install=yes_without_prompt

# FrontDoor and WAF
az afd security-policy delete -g $RESOURCE_GROUP --profile-name "$FRONT_DOOR_PROFILE_NAME" --security-policy-name "$FRONT_DOOR_STACK_NAME" -y
az afd endpoint delete -g $RESOURCE_GROUP --profile-name "$FRONT_DOOR_PROFILE_NAME" --endpoint-name "$FRONT_DOOR_STACK_NAME" -y
az afd origin-group delete -g $RESOURCE_GROUP --profile-name "$FRONT_DOOR_PROFILE_NAME" --origin-group-name "$FRONT_DOOR_STACK_NAME" -y
az network front-door waf-policy delete -g $RESOURCE_GROUP -n "$WAF_POLICY_NAME"

# Containers
for item in $CONTAINER_APP_NAME $VWORKER_CONTAINER_APP_NAME $COMMIT_PROCESSOR_CONTAINER_APP_NAME
do
  az extension add --name containerapp --upgrade
  az containerapp delete -n "$item" -g $RESOURCE_GROUP --yes
  az acr repository delete -n "$REGISTRY_NAME" --repository "$item" --yes || true
done

# Delete Cosmos SQL API database
az cosmosdb sql database delete --account-name $COSMOS_ACCOUNT_NAME --resource-group $RESOURCE_GROUP --name "$COSMOS_DATABASE" --yes

# Delete Blob Storage container
az storage container delete --account-name $STORAGE_ACCOUNT --name "$STORAGE_CONTAINER" --auth-mode login

# Delete Event grid subscription
az eventgrid event-subscription delete --name "$THUMB_EVENT_GRID_SUBSCRIPTION_NAME" --source-resource-id $EVENT_GRID_SOURCE_STORAGE_RESOURCE

az eventgrid event-subscription delete --name "$ACTIVATOR_EVENT_GRID_SUBSCRIPTION_NAME" --source-resource-id $EVENT_GRID_SOURCE_STORAGE_RESOURCE

# Delete Function App
az functionapp delete --name "$FUNCTION_APP" --resource-group $RESOURCE_GROUP --keep-empty-plan

# Delete Kafka topics
confluent kafka topic delete "$THUMB_PROCESSOR_TOPIC" --cluster $KAFKA_CLUSTER_ID --force

confluent kafka topic delete "$COMMIT_PROCESSOR_TOPIC" --cluster $KAFKA_CLUSTER_ID --force

confluent kafka topic delete "$COMMIT_PROCESSOR_DLT" --cluster $KAFKA_CLUSTER_ID --force

for topic in $(confluent kafka topic list --cluster $KAFKA_CLUSTER_ID | grep "$KAFKA_TOPIC_PREFIX")
do
  confluent kafka topic delete "$topic" --force --cluster $KAFKA_CLUSTER_ID
done
