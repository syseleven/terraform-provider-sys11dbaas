terraform {
  required_version = ">= 1.8"
  required_providers {
    openstack = {
      source = "terraform-provider-openstack/openstack"
    }
    sys11dbaas = {
      source  = "registry.terraform.io/syseleven/sys11dbaas"
      version = ">= 0.3.2"
    }
  }
}
