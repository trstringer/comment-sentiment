provider "azurerm" {
  features {}
}

variable "dnszone_name" {}
variable "dnszone_resource_group_name" {}
variable "ip_address" {}

data "azurerm_dns_zone" "dnszone" {
  name = var.dnszone_name
  resource_group_name = var.dnszone_resource_group_name
}

resource "azurerm_dns_a_record" "arecord" {
  name = "a_record"
  zone_name = data.azurerm_dns_zone.dnszone.name
  resource_group_name = data.azurerm_dns_zone.dnszone.resource_group_name
  ttl = 300
  records = [var.ip_address]
}
