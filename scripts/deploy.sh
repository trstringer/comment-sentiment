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
DNS_ZONE_NAME=$(terraform -chdir=./infra/dns output -raw dns_zone_name)
DNS_RESOURCE_GROUP=$(terraform -chdir=./infra/dns output -raw dns_resource_group_name)

az aks get-credentials -g $RESOURCE_NAME -n $RESOURCE_NAME --overwrite-existing

helm repo add jetstack https://charts.jetstack.io
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Install the dependencies first.
helm upgrade cert-manager jetstack/cert-manager \
    --install \
    --namespace cert-manager \
    --create-namespace \
    --version "v1.8.0" \
    --set installCRDs=true
helm upgrade contour bitnami/contour \
    --install \
    --namespace projectcontour \
    --create-namespace \
    --version "7.8.0"

helm upgrade \
    --install \
    --set image.repository=$ACR/comment-sentiment \
    --set image.tag=$VERSION \
    --set keyvault.tenantID=$TENANT_ID \
    --set keyvault.name=$RESOURCE_NAME \
    --set keyvault.userID=$CLUSTER_USER_ID \
    --set languageEndpoint=$LANGUAGE_ENDPOINT \
    --set github.appID=$GITHUB_APP_ID \
    --set fqdn=$DNS_ZONE_NAME \
    comment-sentiment ./charts/comment-sentiment

ENVOY_IP_ADDRESS=$(kubectl get svc -n projectcontour \
    -l "app.kubernetes.io/component=envoy" \
    -o jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}')

terraform -chdir=./infra/dns_records apply -auto-approve \
    -var="dnszone_name=$DNS_ZONE_NAME" \
    -var="dnszone_resource_group_name=$DNS_RESOURCE_GROUP" \
    -var="ip_address=$ENVOY_IP_ADDRESS"
