#!/bin/bash

terraform -chdir=./infra apply -auto-approve

RESOURCE_NAME=$(terraform -chdir=./infra output -raw resource_name)
OBJECT_ID=$(az aks show \
    -g "$RESOURCE_NAME" \
    -n "$RESOURCE_NAME" \
    --query "addonProfiles.azureKeyvaultSecretsProvider.identity.objectId" -o tsv)

terraform -chdir=./infra/secrets apply -auto-approve \
    -var="resource_name=$RESOURCE_NAME" \
    -var="object_id=$OBJECT_ID"
