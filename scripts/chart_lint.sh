#!/bin/bash

ACR=$(terraform -chdir=./infra output -raw acr_endpoint)
VERSION=$(./dist/comment-sentiment -v)

helm lint \
    --set image.repository=$ACR/comment-sentiment \
    --set image.tag=$VERSION \
    ./charts/comment-sentiment
