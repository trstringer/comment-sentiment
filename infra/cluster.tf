provider "azurerm" {
  features {}
}

variable "resource_name" {
  default = "happyosscluster1"
}

variable "location" {
  default = "eastus"
}

resource "azurerm_resource_group" "rg" {
  name = var.resource_name
  location = var.location
}

resource "azurerm_kubernetes_cluster" "default" {
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
