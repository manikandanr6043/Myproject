#!/bin/bash
set -ex

STACK_SUFFIX="${1:-$USER$USERNAME}"
RESOURCE_GROUP=RG-RTMZ-TDRV-DEVELOPMENT
SUBSCRIPTION_ID=16df61f4-08f2-4bf6-acc3-1d6b5cefde97

KEY_VAULT=KV-RTMZ-TDRV-DEVELOPMENT
MONGO_DBNAME=tdrive-$STACK_SUFFIX
STORAGE_ACCOUNT=srgrtmztdrvdevcontent
STORAGE_CONTAINER=tdrive-$STACK_SUFFIX
FUNCTION_APP=tdrive-upload-function-app-$STACK_SUFFIX
KAFKA_CLUSTER_ID="lkc-gnxnjr"

DB_KAFKA_TOPIC=db.$MONGO_DBNAME.latest
DB_KAFKA_DLT=db.$MONGO_DBNAME.latest.dlt
THUMB_PROCESSOR_TOPIC="trimble.tdrive.thumb_processor-$STACK_SUFFIX"
COMMIT_PROCESSOR_TOPIC="trimble.tdrive.commit_processor-$STACK_SUFFIX"
CONTROL_TOPIC=db.$MONGO_DBNAME.control

THUMB_EVENT_GRID_SUBSCRIPTION_NAME="thumb-blob-trigger-$STACK_SUFFIX"
ACTIVATOR_EVENT_GRID_SUBSCRIPTION_NAME="activator-blob-trigger-$STACK_SUFFIX"
EVENT_TYPE_FILTER=Microsoft.Storage.BlobCreated
THUMB_EVENT_SUBJECT_PREFIX="/blobServices/default/containers/$STORAGE_CONTAINER/blobs/thumb"
ACTIVATOR_EVENT_SUBJECT_PREFIX="/blobServices/default/containers/$STORAGE_CONTAINER/blobs/orig"
EVENT_GRID_SOURCE_STORAGE_RESOURCE="/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP/providers/Microsoft.Storage/storageaccounts/$STORAGE_ACCOUNT"
FUNC_SUBSCRIPTION_ENDPOINT="/subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP/providers/Microsoft.Web/sites/$FUNCTION_APP/functions"
THUMB_FUNCTION_ENDPOINT="$FUNC_SUBSCRIPTION_ENDPOINT/thumbnail_observer"
ACTIVATOR_FUNCTION_ENDPOINT="$FUNC_SUBSCRIPTION_ENDPOINT/file_activator"


## Start DB Scripts
mongosh "$(az keyvault secret show --name MONGO-URI --vault-name $KEY_VAULT --query value | tr -d '"' | sed "s/\//\/$MONGO_DBNAME/3")" ./scripts/db-setup.js
## End DB Scripts


## Start Confluent Kafka
confluent kafka topic create "$DB_KAFKA_TOPIC" --cluster $KAFKA_CLUSTER_ID --partitions 1 --if-not-exists
confluent kafka topic create "$DB_KAFKA_DLT" --cluster $KAFKA_CLUSTER_ID --partitions 1 --if-not-exists
confluent kafka topic create "$CONTROL_TOPIC" --cluster $KAFKA_CLUSTER_ID --partitions 1 --if-not-exists

confluent kafka topic create "$COMMIT_PROCESSOR_TOPIC" --cluster $KAFKA_CLUSTER_ID --partitions 1 --if-not-exists

confluent kafka topic create "$THUMB_PROCESSOR_TOPIC" --cluster $KAFKA_CLUSTER_ID --partitions 1 --if-not-exists
## End Confluent Kafka


## Start container images
./scripts/build-push-images.sh "$STACK_SUFFIX"
## End container images


## Start container images
./scripts/deploy-bicep.sh "$STACK_SUFFIX"
## End container images


## Start update containers with latest images
./scripts/deploy-images.sh "$STACK_SUFFIX"
## End update containers with latest images


## Start publish Function App
cd ./src/functions/azure/UploadFunctionApp || exit

func azure functionapp publish "$FUNCTION_APP" --build remote --python


echo "creating event grid subscription for thumbnail_observer"
az eventgrid event-subscription create \
  --name "$THUMB_EVENT_GRID_SUBSCRIPTION_NAME" \
  --source-resource-id $EVENT_GRID_SOURCE_STORAGE_RESOURCE \
  --advanced-filter eventtype StringIn $EVENT_TYPE_FILTER \
  --subject-begins-with "$THUMB_EVENT_SUBJECT_PREFIX" \
  --endpoint "$THUMB_FUNCTION_ENDPOINT" \
  --endpoint-type azurefunction

echo "creating event grid subscription for file_activator"
az eventgrid event-subscription create \
  --name "$ACTIVATOR_EVENT_GRID_SUBSCRIPTION_NAME" \
  --source-resource-id $EVENT_GRID_SOURCE_STORAGE_RESOURCE \
  --advanced-filter eventtype StringIn $EVENT_TYPE_FILTER \
  --subject-begins-with "$ACTIVATOR_EVENT_SUBJECT_PREFIX" \
  --endpoint "$ACTIVATOR_FUNCTION_ENDPOINT" \
  --endpoint-type azurefunction
## End publish Function App
