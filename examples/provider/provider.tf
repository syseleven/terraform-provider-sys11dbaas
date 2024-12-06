terraform {
  required_providers {
    sys11dbaas = {
      source = "syseleven/sys11dbaas"
    }
  }
}

provider "sys11dbaas" {
  api_key      = "s11_prak_..."
  project      = "0123456789"
  organization = "0123-456-78-9"
}
