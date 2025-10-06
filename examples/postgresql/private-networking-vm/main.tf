provider "openstack" {
  region = var.region
}
provider "sys11dbaas" {}

locals {
  virtual_machine_subnet_cidr = "192.168.1.0/24"
  external_network            = "ext-net"
}
