provider "azurerm" {
  features {}
}

variable "resource_name" {
  default = "happyoss1"
}

variable "location" {
  default = "eastus"
}

resource "azurerm_resource_group" "rg" {
  name = var.resource_name
  location = var.location
}

resource "azurerm_cognitive_account" "textanalytics" {
  name = var.resource_name
  resource_group_name = azurerm_resource_group.rg.name
  location = var.location
  kind = "TextAnalytics"
  sku_name = "S"
  custom_subdomain_name = var.resource_name

  network_acls {
    default_action = "Allow"
    ip_rules = []
  }
}

resource "azurerm_container_registry" "acr" {
  name = var.resource_name
  resource_group_name = azurerm_resource_group.rg.name
  location = var.location
  sku = "Basic"
}

resource "azurerm_kubernetes_cluster" "aks" {
  name = var.resource_name
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

  dns_prefix = var.resource_name
}

resource "azurerm_role_assignment" "acraks" {
  principal_id = azurerm_kubernetes_cluster.aks.kubelet_identity[0].object_id
  role_definition_name = "AcrPull"
  scope = azurerm_container_registry.acr.id
  skip_service_principal_aad_check = true
}
