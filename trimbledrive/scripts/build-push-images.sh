#!/bin/bash
set -ex

STACK_SUFFIX="${1:-$USER$USERNAME}"
REGISTRY_NAME=crrtmztdrvdevelopment
REGISTRY_LOGIN_SERVER=$REGISTRY_NAME.azurecr.io
IMAGE_NAME=tdrive-api-$STACK_SUFFIX
VWORKER_IMAGE_NAME=tdrive-vworker-$STACK_SUFFIX
COMMIT_PROCESSOR_IMAGE_NAME=tdrive-commit-$STACK_SUFFIX


az account set --name trimble-connect-platform-sandbox
az acr login --name $REGISTRY_NAME

docker build --target debug -f src/api/Dockerfile -t $REGISTRY_LOGIN_SERVER/"$IMAGE_NAME" --build-arg APP_VERSION=v0.0.0-"$STACK_SUFFIX" . &
docker build --target debug -f src/versions-worker/Dockerfile -t $REGISTRY_LOGIN_SERVER/"$VWORKER_IMAGE_NAME" --build-arg APP_VERSION=v0.0.0-"$STACK_SUFFIX" . &
docker build --target debug -f src/commit-processor/Dockerfile -t $REGISTRY_LOGIN_SERVER/"$COMMIT_PROCESSOR_IMAGE_NAME" .

docker push $REGISTRY_LOGIN_SERVER/"$IMAGE_NAME" &
docker push $REGISTRY_LOGIN_SERVER/"$VWORKER_IMAGE_NAME" &
docker push $REGISTRY_LOGIN_SERVER/"$COMMIT_PROCESSOR_IMAGE_NAME"

