#!/bin/bash

RESOURCE_NAME=$(terraform -chdir=./infra output -raw resource_name)
OBJECT_ID=$(az aks show \
    -g "$RESOURCE_NAME" \
    -n "$RESOURCE_NAME" \
    --query "addonProfiles.azureKeyvaultSecretsProvider.identity.objectId" -o tsv)

terraform -chdir=./infra/secrets destroy -auto-approve \
    -var="object_id=$OBJECT_ID" \
    -var="resource_name=$RESOURCE_NAME"

terraform -chdir=./infra destroy -auto-approve
