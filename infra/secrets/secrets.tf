provider "azurerm" {
  features {}
}

variable "resource_name" {}

variable "object_id" {}

data "azurerm_resource_group" "rg" {
  name = var.resource_name
}

data "azurerm_key_vault" "akv" {
  name = var.resource_name
  resource_group_name = data.azurerm_resource_group.rg.name
}

resource "azurerm_key_vault_access_policy" "secretsaccess" {
  key_vault_id = data.azurerm_key_vault.akv.id
  tenant_id = data.azurerm_key_vault.akv.tenant_id
  object_id = var.object_id

  secret_permissions = ["Get"]
}
