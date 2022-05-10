#!/bin/bash

if [[ -z "$KEYVAULT" ]]; then
    echo "KEYVAULT needs to be set to the initial key vault with app specs"
    exit 1
fi

PRIVATEKEY=$(az keyvault secret show \
    --vault-name $KEYVAULT \
    --name happyossprivatekey \
    --query value -o tsv)

terraform -chdir=./infra apply -auto-approve \
    -var="privatekey=$PRIVATEKEY"

RESOURCE_NAME=$(terraform -chdir=./infra output -raw resource_name)
OBJECT_ID=$(az aks show \
    -g "$RESOURCE_NAME" \
    -n "$RESOURCE_NAME" \
    --query "addonProfiles.azureKeyvaultSecretsProvider.identity.objectId" -o tsv)

terraform -chdir=./infra/secrets apply -auto-approve \
    -var="resource_name=$RESOURCE_NAME" \
    -var="object_id=$OBJECT_ID"
