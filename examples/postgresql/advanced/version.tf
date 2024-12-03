terraform {
  required_version = ">= 1.8"
  required_providers {
    sys11dbaas = {
      source  = "registry.terraform.io/syseleven/sys11dbaas"
      version = ">= 0.2.0"
    }
  }
}
