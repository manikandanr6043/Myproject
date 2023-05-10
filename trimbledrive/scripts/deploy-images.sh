#!/bin/bash
set -ex

STACK_SUFFIX="${1:-$USER$USERNAME}"
RESOURCE_GROUP_NAME="${2:-RG-RTMZ-TDRV-DEVELOPMENT}"

REGISTRY_NAME=crrtmztdrvdevelopment
REGISTRY_LOGIN_SERVER=$REGISTRY_NAME.azurecr.io
IMAGE_NAME=tdrive-api-$STACK_SUFFIX
VWORKER_IMAGE_NAME=tdrive-vworker-$STACK_SUFFIX
COMMIT_PROCESSOR_IMAGE_NAME=tdrive-commit-$STACK_SUFFIX

az account set --name trimble-connect-platform-sandbox
az acr login --name $REGISTRY_NAME

IMAGE_DIGEST=$(docker inspect --format='{{index .RepoDigests 0}}' $REGISTRY_LOGIN_SERVER/"$IMAGE_NAME") &&
az containerapp update \
  --name "$IMAGE_NAME" \
  --resource-group "$RESOURCE_GROUP_NAME" \
  --image "$IMAGE_DIGEST" &

VWORKER_IMAGE_DIGEST=$(docker inspect --format='{{index .RepoDigests 0}}' $REGISTRY_LOGIN_SERVER/"$VWORKER_IMAGE_NAME") &&
az containerapp update \
  --name "$VWORKER_IMAGE_NAME" \
  --resource-group "$RESOURCE_GROUP_NAME" \
  --image "$VWORKER_IMAGE_DIGEST" &

COMMIT_PROCESSOR_IMAGE_DIGEST=$(docker inspect --format='{{index .RepoDigests 0}}' $REGISTRY_LOGIN_SERVER/"$COMMIT_PROCESSOR_IMAGE_NAME") &&
az containerapp update \
  --name "$COMMIT_PROCESSOR_IMAGE_NAME" \
  --resource-group "$RESOURCE_GROUP_NAME" \
  --image "$COMMIT_PROCESSOR_IMAGE_DIGEST"
