provider "openstack" {
  region = var.region
}
provider "metakube" {
  host = "https://metakube.apis.syseleven.de"
}
provider "sys11dbaas" {}
provider "environment" {}

locals {
  metakube_node_subnet_cidr = "192.168.1.0/24"
}
