provider "azurerm" {
  features {}
}

variable "resource_name" {}

variable "location" {
  default = "eastus"
}

variable "dnsname" {}

resource "azurerm_resource_group" "rg" {
  name = var.resource_name
  location = var.location
}

resource "azurerm_dns_zone" "dnszone" {
  resource_group_name = azurerm_resource_group.rg.name
  name = var.dnsname
}

output "dns_servers" {
  value = azurerm_dns_zone.dnszone.name_servers
}

output "dns_zone_name" {
  value = azurerm_dns_zone.dnszone.name
}

output "dns_resource_group_name" {
  value = azurerm_resource_group.rg.name
}
