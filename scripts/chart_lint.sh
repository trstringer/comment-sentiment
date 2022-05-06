#!/bin/bash

ACR=$(terraform -chdir=./infra output -raw acr_endpoint)
VERSION=$(./dist/comment-sentiment -v)
TENANT_ID=$(terraform -chdir=./infra output -raw tenant_id)
RESOURCE_NAME=$(terraform -chdir=./infra output -raw resource_name)
CLUSTER_USER_ID=$(terraform -chdir=./infra output -raw cluster_identity_id)

helm lint \
    --set image.repository=$ACR/comment-sentiment \
    --set image.tag=$VERSION \
    --set keyvault.tenantID=$TENANT_ID \
    --set keyvault.name=$RESOURCE_NAME \
    --set keyvault.userID=$CLUSTER_USER_ID \
    ./charts/comment-sentiment
