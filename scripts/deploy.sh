#!/bin/bash

ACR=$(terraform -chdir=./infra output -raw acr_endpoint)
RESOURCE_NAME=$(terraform -chdir=./infra output -raw resource_name)
VERSION=$(./dist/comment-sentiment -v)
TENANT_ID=$(terraform -chdir=./infra output -raw tenant_id)
CLUSTER_USER_ID=$(az aks show \
    -g "$RESOURCE_NAME" \
    -n "$RESOURCE_NAME" \
    --query "addonProfiles.azureKeyvaultSecretsProvider.identity.clientId" -o tsv)

az aks get-credentials -g $RESOURCE_NAME -n $RESOURCE_NAME --overwrite-existing

helm repo update
helm dependency build ./charts/comment-sentiment
helm install \
    --set image.repository=$ACR/comment-sentiment \
    --set image.tag=$VERSION \
    --set keyvault.tenantID=$TENANT_ID \
    --set keyvault.name=$RESOURCE_NAME \
    --set keyvault.userID=$CLUSTER_USER_ID \
    comment-sentiment ./charts/comment-sentiment
