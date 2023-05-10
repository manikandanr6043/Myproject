#!/bin/bash

while getopts c:g:d: flag
do
    case "${flag}" in
        c) COSMOS_ACCOUNT_NAME=${OPTARG};;
        g) RESOURCE_GROUP_NAME=${OPTARG};;
        d) DB_NAME_PREFIX=${OPTARG};;
        *)
    esac
done

IFS=$'\n' read -rd '' -a dbsToDelete <<< "$(az cosmosdb sql database list -a "$COSMOS_ACCOUNT_NAME" -g "$RESOURCE_GROUP_NAME" --query "[?contains(name, '$DB_NAME_PREFIX')].{Name:name}" -o tsv)"
for dbToDelete in "${dbsToDelete[@]}"
do
    echo "Deleting cosmos DB $dbToDelete"
    az cosmosdb sql database delete --account-name "$COSMOS_ACCOUNT_NAME" -g "$RESOURCE_GROUP_NAME" --name "$dbToDelete" --yes
done
