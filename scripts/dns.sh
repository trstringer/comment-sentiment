#!/bin/bash

DNS_NAME="yewhampshire.com"

terraform -chdir=./infra/dns apply -auto-approve \
    -var="resource_name=happyossdns" \
    -var="dnsname=$DNS_NAME"

echo "Navigate to domains.google.com and add the following name servers to $DNS_NAME"
terraform -chdir=./infra/dns output -json dns_servers
