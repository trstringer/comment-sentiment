#!/bin/bash

if [[ -z "$GITHUB_APP_ID" ]]; then
    echo "GITHUB_APP_ID must be set to the registered GitHub app ID"
    exit 1
fi

ACR=$(terraform -chdir=./infra output -raw acr_endpoint)
RESOURCE_NAME=$(terraform -chdir=./infra output -raw resource_name)
VERSION=$(./dist/comment-sentiment -v)
TENANT_ID=$(terraform -chdir=./infra output -raw tenant_id)
CLUSTER_USER_ID=$(az aks show \
    -g "$RESOURCE_NAME" \
    -n "$RESOURCE_NAME" \
    --query "addonProfiles.azureKeyvaultSecretsProvider.identity.clientId" -o tsv)
LANGUAGE_ENDPOINT=$(terraform -chdir=./infra output -raw language_endpoint)

az aks get-credentials -g $RESOURCE_NAME -n $RESOURCE_NAME --overwrite-existing

helm repo update
helm dependency build ./charts/comment-sentiment
helm upgrade \
    --install \
    --set image.repository=$ACR/comment-sentiment \
    --set image.tag=$VERSION \
    --set keyvault.tenantID=$TENANT_ID \
    --set keyvault.name=$RESOURCE_NAME \
    --set keyvault.userID=$CLUSTER_USER_ID \
    --set languageEndpoint=$LANGUAGE_ENDPOINT \
    --set github.appID=$GITHUB_APP_ID \
    comment-sentiment ./charts/comment-sentiment
