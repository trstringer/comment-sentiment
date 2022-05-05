#!/bin/bash

ACR=$(terraform -chdir=./infra output -raw acr_endpoint)
RESOURCE_NAME=$(terraform -chdir=./infra output -raw resource_name)
VERSION=$(./dist/comment-sentiment -v)

az aks get-credentials -g $RESOURCE_NAME -n $RESOURCE_NAME --overwrite-existing

helm repo update
helm dependency build ./charts/comment-sentiment
helm install \
    --set image.repository=$ACR/comment-sentiment \
    --set image.tag=$VERSION \
    comment-sentiment ./charts/comment-sentiment
