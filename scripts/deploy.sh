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

ENVOY_IP_ADDRESS=$(kubectl get svc \
    -l "app.kubernetes.io/component=envoy" \
    -o jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}')

DNS_ZONE_NAME=$(terraform -chdir=./infra/dns output -raw dns_zone_name)
DNS_RESOURCE_GROUP=$(terraform -chdir=./infra/dns output -raw dns_resource_group_name)

terraform -chdir=./infra/dns_records apply -auto-approve \
    -var="dnszone_name=$DNS_ZONE_NAME" \
    -var="dnszone_resource_group_name=$DNS_RESOURCE_GROUP" \
    -var="ip_address=$ENVOY_IP_ADDRESS"
