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
    type = "SystemAssigned"
  }

  dns_prefix = local.resource_name

  key_vault_secrets_provider {
    secret_rotation_enabled = false
  }
}

resource "azurerm_role_assignment" "acraks" {
  principal_id = azurerm_kubernetes_cluster.aks.kubelet_identity[0].object_id
  role_definition_name = "AcrPull"
  scope = azurerm_container_registry.acr.id
  skip_service_principal_aad_check = true
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
    secret_permissions = ["Get", "List", "Set"]
  }
  access_policy {
    tenant_id = azurerm_kubernetes_cluster.aks.identity[0].tenant_id
    object_id = azurerm_kubernetes_cluster.aks.identity[0].principal_id
    secret_permissions = ["Get"]
  }
}

resource "azurerm_key_vault_secret" "languagekeysecret" {
  key_vault_id = azurerm_key_vault.akv.id
  name = "languagekey"
  value = azurerm_cognitive_account.textanalytics.primary_access_key
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
