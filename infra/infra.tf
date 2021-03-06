provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

locals {
  resource_name = "happyoss${random_integer.resourceid.result}"
}

resource "random_integer" "resourceid" {
  min = 10000
  max = 99999
}

variable "resource_name" {
  default = "happyoss1"
}

variable "location" {
  default = "eastus"
}

variable "privatekey" {
  sensitive = true
}

variable "webhooksecret" {
  sensitive = true
}

resource "azurerm_resource_group" "rg" {
  name = local.resource_name
  location = var.location
}

resource "azurerm_cognitive_account" "textanalytics" {
  name = local.resource_name
  resource_group_name = azurerm_resource_group.rg.name
  location = var.location
  kind = "TextAnalytics"
  sku_name = "S"
  custom_subdomain_name = local.resource_name

  network_acls {
    default_action = "Allow"
    ip_rules = []
  }
}

resource "azurerm_container_registry" "acr" {
  name = local.resource_name
  resource_group_name = azurerm_resource_group.rg.name
  location = var.location
  sku = "Basic"
}

resource "azurerm_user_assigned_identity" "aksuser" {
  name = local.resource_name
  resource_group_name = azurerm_resource_group.rg.name
  location = var.location
}

resource "azurerm_kubernetes_cluster" "aks" {
  name = local.resource_name
  resource_group_name = azurerm_resource_group.rg.name
  location = var.location

  default_node_pool {
    name = "default"
    vm_size = "Standard_B2s"
    node_count = 1
  }

  identity {
    type = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.aksuser.id]
  }

  dns_prefix = local.resource_name

  key_vault_secrets_provider {
    secret_rotation_enabled = false
  }
}

resource "azurerm_role_assignment" "acraks" {
  principal_id = azurerm_kubernetes_cluster.aks.kubelet_identity[0].object_id
  # principal_id = azurerm_user_assigned_identity.aksuser.principal_id
  role_definition_name = "AcrPull"
  scope = azurerm_container_registry.acr.id
  # skip_service_principal_aad_check = true
}

resource "azurerm_key_vault" "akv" {
  name = local.resource_name
  location = var.location
  resource_group_name = azurerm_resource_group.rg.name
  tenant_id = data.azurerm_client_config.current.tenant_id
  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    secret_permissions = ["Get", "List", "Set", "Delete", "Purge"]
  }
}

resource "azurerm_key_vault_secret" "languagekeysecret" {
  key_vault_id = azurerm_key_vault.akv.id
  name = "languagekey"
  value = azurerm_cognitive_account.textanalytics.primary_access_key
}

resource "azurerm_key_vault_secret" "privatekeysecret" {
  key_vault_id = azurerm_key_vault.akv.id
  name = "happyossprivatekey"
  value = var.privatekey
}

resource "azurerm_key_vault_secret" "webhooksecret" {
  key_vault_id = azurerm_key_vault.akv.id
  name = "happyosswebhooksecret"
  value = var.webhooksecret
}

output "acr_endpoint" {
  value = azurerm_container_registry.acr.login_server
}

output "language_endpoint" {
  value = azurerm_cognitive_account.textanalytics.endpoint
}

output "resource_name" {
  value = local.resource_name
}

output "tenant_id" {
  value = data.azurerm_client_config.current.tenant_id
}

output "cluster_identity_id" {
  value = azurerm_user_assigned_identity.aksuser.client_id
}
