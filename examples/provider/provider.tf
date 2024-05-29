terraform {
  required_providers {
    sys11dbaas = {
      source = "syseleven/sys11dbaas"
    }
  }
}

provider "sys11dbaas" {
  api_key      = var.api_key
  project      = var.project
  organization = var.org
}
